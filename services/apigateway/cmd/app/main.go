package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"deblasis.net/space-traffic-control/common/config"
	consulreg "deblasis.net/space-traffic-control/common/consul"
	authpb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	ccpb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"

	// "deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"

	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	auth_service "deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	cc_service "deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	ServiceName = "deblasis-v1-APIGateway"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	var (
		httpAddr = net.JoinHostPort(cfg.ListenAddr, cfg.HttpServerPort)
		// consulAddr = net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)

		// retryMax     = cfg.APIGateway.RetryMax
		// retryTimeout = cfg.APIGateway.RetryTimeoutMs * int(time.Millisecond)
	)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Service discovery domain. In this example we use Consul.
	// var client consulsd.Client
	// {
	// 	consulConfig := capi.DefaultConfig()
	// 	if len(consulAddr) > 0 {
	// 		consulConfig.Address = consulAddr
	// 	}
	// 	consulClient, err := capi.NewClient(consulConfig)
	// 	if err != nil {
	// 		logger.Log("err", err)
	// 		os.Exit(1)
	// 	}
	// 	client = consulsd.NewClient(consulClient)

	// }

	// Transport domain.
	var (
		// tracer          = stdopentracing.GlobalTracer() // no-op
		// zipkinTracer, _ = stdzipkin.NewTracer(nil, stdzipkin.WithNoopTracer(true))
		ctx = context.Background()
	)

	// var duration metrics.Histogram
	// {
	// 	// Endpoint-level metrics.
	// 	duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	// 		Namespace: service.Namespace,
	// 		Subsystem: strings.Split(service.ServiceName, ".")[2],
	// 		Name:      "request_duration_seconds",
	// 		Help:      "Request duration in seconds.",
	// 	}, []string{"method", "success"})
	// }
	// http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	// Now we begin installing the routes. Each route corresponds to a single
	// method: sum, concat, uppercase, and count.

	{

		{
			if cfg.Consul.Host != "" && cfg.Consul.Port != "" {
				consulAddres := net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)
				httpPort, _ := strconv.Atoi(cfg.HttpServerPort)
				metricsPort, _ := strconv.Atoi(cfg.HttpServerPort)
				tags := []string{service.Namespace, ServiceName}
				consulReg := consulreg.NewConsulRegister(consulAddres, ServiceName, httpPort, metricsPort, tags, cfg.Logger, cfg.BindOnLocalhost)
				svcRegistar, err := consulReg.NewConsulHTTPRegister()
				defer svcRegistar.Deregister()
				if err != nil {
					level.Error(cfg.Logger).Log(
						"consulAddres", consulAddres,
						"serviceName", ServiceName,
						"metricsPort", metricsPort,
						"tags", tags,
						"err", err,
					)
				}
				svcRegistar.Register()
			}
		}
		{

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			//mux := http.NewServeMux()

			authGw, err := newAuthSvcGateway(ctx, []runtime.ServeMuxOption{})
			if err != nil {
				panic(err)
			}
			//mux.Handle("/", authGw)

			ccGw, err := newCentralCommandSvcGateway(ctx, []runtime.ServeMuxOption{})
			if err != nil {
				panic(err)
			}

			//mux.Handle("/centcom", ccGw)

			// s := &http.Server{
			// 	Addr:    httpAddr,
			// 	Handler: allowCORS(mux),
			// }

			//mux.HandleFunc("/openapiv2/", openAPIServer(opts.OpenAPIDir))
			//mux.HandleFunc("/healthz", healthzServer(conn))

			// gw, err := newGateway(ctx, conn, opts.Mux)
			// if err != nil {
			// 	return err
			// }
			// mux.Handle("/", gw)

			// Interrupt handler.
			errc := make(chan error)
			go func() {
				c := make(chan os.Signal)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				errc <- fmt.Errorf("%s", <-c)
			}()

			go func() {
				logger.Log("transport", "HTTP", "addr", httpAddr)

				router := LoggerMw(logger, allowCORS(NewRouter(authGw, ccGw)), "grpc-gw")

				errc <- http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpServerPort), router)

			}()

			logger.Log("exit", <-errc)
		}

	}
}

// newAuthSvcGateway returns a new gateway server which translates HTTP into gRPC.
func newAuthSvcGateway(ctx context.Context, opts []runtime.ServeMuxOption) (http.Handler, error) {

	conn, err := dial(ctx, fmt.Sprintf("%v.service.consul:%d", auth_service.ServiceName, auth_service.GrpcServerPort))
	if err != nil {
		panic(err)
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			glog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := runtime.NewServeMux(opts...)

	for _, f := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		authpb.RegisterAuthServiceHandler,
	} {
		if err := f(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

// newAuthSvcGateway returns a new gateway server which translates HTTP into gRPC.
func newCentralCommandSvcGateway(ctx context.Context, opts []runtime.ServeMuxOption) (http.Handler, error) {

	conn, err := dial(ctx, fmt.Sprintf("%v.service.consul:%d", cc_service.ServiceName, cc_service.GrpcServerPort))
	if err != nil {
		panic(err)
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			glog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := runtime.NewServeMux(opts...)

	for _, f := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		ccpb.RegisterCentralCommandServiceHandler,
	} {
		if err := f(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

func LoggerMw(logger log.Logger, inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		logger.Log(
			"method", r.Method,
			"uri", r.RequestURI,
			"name", name,
			"duration", time.Since(start),
		)
	})
}

func dial(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, grpc.WithInsecure())
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
}

// healthzServer returns a simple health handler which returns ok.
func healthzServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(authGw http.Handler, ccGw http.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			Index,
		},

		Route{
			"Login",
			strings.ToUpper("Post"),
			"/auth/login",
			authGw.ServeHTTP,
		},

		Route{
			"Signup",
			strings.ToUpper("Post"),
			"/user/signup",
			authGw.ServeHTTP,
		},

		Route{
			"ShipRegister",
			strings.ToUpper("Post"),
			"/centcom/ship/register",
			ccGw.ServeHTTP,
		},

		Route{
			"ShipsList",
			strings.ToUpper("Get"),
			"/centcom/ship/all",
			ccGw.ServeHTTP,
		},

		Route{
			"StationRegister",
			strings.ToUpper("Post"),
			"/centcom/station/register",
			ccGw.ServeHTTP,
		},

		Route{
			"StationsList",
			strings.ToUpper("Get"),
			"/centcom/station/all",
			ccGw.ServeHTTP,
		},

		// Route{
		// 	"ShipLand",
		// 	strings.ToUpper("Post"),
		// 	"/shipping-station/land",
		// 	ShipLand,
		// },

		// Route{
		// 	"ShipRequestLanding",
		// 	strings.ToUpper("Post"),
		// 	"/shipping-station/request-landing",
		// 	ShipRequestLanding,
		// },
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		//handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
