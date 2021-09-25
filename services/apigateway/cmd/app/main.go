package main

import (
	"context"
	"encoding/json"
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
	"deblasis.net/space-traffic-control/common/errs"
	authpb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	ccpb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	sspb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"github.com/etherlabsio/healthcheck/v2"

	// "deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"

	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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

	var (
		logger log.Logger = cfg.Logger
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

			authGw, err := newAuthSvcGateway(ctx, cfg,
				runtime.WithForwardResponseOption(httpHeaderRewriter(log.With(logger, "component", "httpHeaderRewriter"))),
				runtime.WithErrorHandler(noContentErrorHandler(log.With(logger, "component", "noContentErrorHandler"))),
			)
			if err != nil {
				panic(err)
			}
			//mux.Handle("/", authGw)

			ccGw, err := newCentralCommandSvcGateway(ctx, cfg,
				runtime.WithForwardResponseOption(httpHeaderRewriter(log.With(logger, "component", "httpHeaderRewriter"))),
				runtime.WithErrorHandler(noContentErrorHandler(log.With(logger, "component", "noContentErrorHandler"))),
			)
			if err != nil {
				panic(err)
			}

			ssGw, err := newShippingStationSvcGateway(ctx, cfg,
				runtime.WithForwardResponseOption(httpHeaderRewriter(log.With(logger, "component", "httpHeaderRewriter"))),
				runtime.WithErrorHandler(noContentErrorHandler(log.With(logger, "component", "noContentErrorHandler"))),
			)
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

				healthchecks := []healthcheck.Option{
					healthcheck.WithChecker(
						"AuthService", healthcheck.CheckerFunc(func(ctx context.Context) error {
							timeout := 2 * time.Second
							_, err := net.DialTimeout("tcp", cfg.APIGateway.AuthServiceGRPCEndpoint, timeout)
							return err
						})),
					healthcheck.WithChecker(
						"CentralCommandService", healthcheck.CheckerFunc(func(ctx context.Context) error {
							timeout := 2 * time.Second
							_, err := net.DialTimeout("tcp", cfg.APIGateway.CentralCommandServiceGRPCEndpoint, timeout)
							return err
						})),
					healthcheck.WithChecker(
						"ShippingStationService", healthcheck.CheckerFunc(func(ctx context.Context) error {
							timeout := 2 * time.Second
							_, err := net.DialTimeout("tcp", cfg.APIGateway.ShippingStationGRPCEndpoint, timeout)
							return err
						})),
				}

				router := LoggerMw(logger, allowCORS(NewRouter(authGw, ccGw, ssGw, healthchecks)), "grpc-gw")

				errc <- http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpServerPort), router)

			}()

			logger.Log("exit", <-errc)
		}

	}
}

// newAuthSvcGateway returns a new gateway server which translates HTTP into gRPC.
func newAuthSvcGateway(ctx context.Context, cfg config.Config, opts ...runtime.ServeMuxOption) (http.Handler, error) {

	//conn, err := dial(ctx, fmt.Sprintf("%v.service.consul:%d", auth_service.ServiceName, auth_service.GrpcServerPort))
	conn, err := dial(ctx, cfg.APIGateway.AuthServiceGRPCEndpoint)
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
func newCentralCommandSvcGateway(ctx context.Context, cfg config.Config, opts ...runtime.ServeMuxOption) (http.Handler, error) {

	//conn, err := dial(ctx, fmt.Sprintf("%v.service.consul:%d", cc_service.ServiceName, cc_service.GrpcServerPort))
	conn, err := dial(ctx, cfg.APIGateway.CentralCommandServiceGRPCEndpoint)
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

func newShippingStationSvcGateway(ctx context.Context, cfg config.Config, opts ...runtime.ServeMuxOption) (http.Handler, error) {

	//conn, err := dial(ctx, fmt.Sprintf("%v.service.consul:%d", cc_service.ServiceName, cc_service.GrpcServerPort))
	conn, err := dial(ctx, cfg.APIGateway.ShippingStationGRPCEndpoint)
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
		sspb.RegisterShippingStationServiceHandler,
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

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(authGw http.Handler, ccGw http.Handler, ssGw http.Handler, healthchecks []healthcheck.Option) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	//TODO index route, maybe swagger?
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

		Route{
			"ShipLand",
			strings.ToUpper("Post"),
			"/shipping-station/land",
			ssGw.ServeHTTP,
		},

		Route{
			"ShipRequestLanding",
			strings.ToUpper("Post"),
			"/shipping-station/request-landing",
			ssGw.ServeHTTP,
		},
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

	if len(healthchecks) > 0 {
		router.Handle("/health", healthcheck.Handler(healthchecks...))
	}
	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func noContentErrorHandler(logger log.Logger) func(ctx context.Context, sm *runtime.ServeMux, m runtime.Marshaler, rw http.ResponseWriter, r *http.Request, e error) {
	return func(ctx context.Context, sm *runtime.ServeMux, m runtime.Marshaler, rw http.ResponseWriter, r *http.Request, e error) {
		//TODO refactor

		logger.Log("msg", "checking metadata")
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			logger.Log("msg", "no md received from context")
		}

		if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
			logger.Log("msg", "x-http-code is "+vals[0])
			code, err := strconv.Atoi(vals[0])
			if err != nil {
				logger.Log("err", err)
				panic(err)
			}

			noContent := false

			if vals = md.HeaderMD.Get("x-no-content"); len(vals) > 0 {
				logger.Log("msg", "x-no-content exists")
				noContent, err = strconv.ParseBool(vals[0])
				if err != nil {
					logger.Log("err", err)
					panic(err)
				}
			}

			delete(md.HeaderMD, "x-http-code")
			delete(rw.Header(), "Grpc-Metadata-X-Http-Code")

			rw.WriteHeader(code)
			if noContent {
				rw.Header().Add("x-stc-error", strings.ReplaceAll(e.Error(), "\n", "--"))
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(&errs.Err{
				Message: e.Error(),
			})

			return
		}
		logger.Log("msg", "executing DefaultHTTPErrorHandler")
		runtime.DefaultHTTPErrorHandler(ctx, sm, m, rw, r, e)
	}
}

type errorBody struct {
	Error string `json:"error,omitempty"`
}

func httpHeaderRewriter(logger log.Logger) func(c context.Context, rw http.ResponseWriter, m proto.Message) error {
	return func(ctx context.Context, rw http.ResponseWriter, m proto.Message) error {
		//TODO refactor

		var (
			code      int
			noContent bool
			err       error
		)
		logger.Log(
			"msg", "checking metadata",
		)
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			logger.Log("msg", "no md received from context")
		}

		// set http status code
		if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
			logger.Log("msg", "x-http-code is "+vals[0])
			code, err = strconv.Atoi(vals[0])
			if err != nil {
				logger.Log("err", err)
				panic(err)
			}
		}
		if vals := md.HeaderMD.Get("x-no-content"); len(vals) > 0 {
			logger.Log("msg", "x-no-content exists")
			noContent, err = strconv.ParseBool(vals[0])
			if err != nil {
				logger.Log("err", err)
				panic(err)
			}
		}
		// delete the headers to not expose any grpc-metadata in http response
		if noContent {
			delete(rw.Header(), "Grpc-Metadata-Content-Type")
			delete(rw.Header(), "Grpc-Metadata-X-No-Content")
		} else {
			delete(rw.Header(), "Grpc-Metadata-X-Stc-Error")
		}
		if code != 0 {
			delete(md.HeaderMD, "x-http-code")
			delete(rw.Header(), "Grpc-Metadata-X-Http-Code")
			rw.WriteHeader(code)
		}

		return nil
	}
}
