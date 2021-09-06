package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"deblasis.net/space-traffic-control/common/config"

	"deblasis.net/space-traffic-control/services/apigateway/internal/api"
	api_v1 "deblasis.net/space-traffic-control/services/apigateway/internal/api/v1"

	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/gorilla/mux"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"

	auth_dbsvc_endpoints "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	auth_dbsvc_service "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"

	"github.com/go-kit/kit/log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	var (
		httpAddr   = net.JoinHostPort(cfg.ListenAddr, cfg.HttpServerPort)
		consulAddr = net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)

		retryMax     = cfg.APIGateway.RetryMax
		retryTimeout = cfg.APIGateway.RetryTimeoutMs * int(time.Millisecond)
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if len(consulAddr) > 0 {
			consulConfig.Address = consulAddr
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	// Transport domain.
	var (
		tracer          = stdopentracing.GlobalTracer() // no-op
		zipkinTracer, _ = stdzipkin.NewTracer(nil, stdzipkin.WithNoopTracer(true))
		ctx             = context.Background()
		r               = mux.NewRouter()
	)

	// Now we begin installing the routes. Each route corresponds to a single
	// method: sum, concat, uppercase, and count.

	{
		// Each method gets constructed with a factory. Factories take an
		// instance string, and return a specific endpoint. In the factory we
		// dial the instance string we get from Consul, and then leverage an
		// addsvc client package to construct a complete service. We can then
		// leverage the addsvc.Make{Sum,Concat}Endpoint constructors to convert
		// the complete service to specific endpoint.
		var (
			passingOnly = true
			endpoints   = auth_dbsvc_endpoints.EndpointSet
			instancer   = consulsd.NewInstancer(client, logger, auth_dbsvc_service.ServiceName, auth_dbsvc_service.Tags, passingOnly)
		)
		// {
		// 	factory := addsvcFactory(addendpoint.MakeSumEndpoint, tracer, zipkinTracer, logger)
		// 	endpointer := sd.NewEndpointer(instancer, factory, logger)
		// 	balancer := lb.NewRoundRobin(endpointer)
		// 	retry := lb.Retry(*retryMax, *retryTimeout, balancer)
		// 	endpoints.SumEndpoint = retry
		// }
		// {
		// 	factory := addsvcFactory(addendpoint.MakeConcatEndpoint, tracer, zipkinTracer, logger)
		// 	endpointer := sd.NewEndpointer(instancer, factory, logger)
		// 	balancer := lb.NewRoundRobin(endpointer)
		// 	retry := lb.Retry(*retryMax, *retryTimeout, balancer)
		// 	endpoints.ConcatEndpoint = retry
		// }

		var routes = api.Routes{
			api.Route{
				"Index",
				"GET",
				"/v1/",
				Index,
			},

			api.Route{
				"Login",
				http.MethodPost,
				"/v1/auth/login",
				api_v1.Login,
			},

			api.Route{
				"Signup",
				http.MethodPost,
				"/v1/user/signup",
				api_v1.Signup,
			},

			api.Route{
				"ShipRegister",
				http.MethodPost,
				"/v1/centcom/ship/register",
				api_v1.ShipRegister,
			},

			api.Route{
				"ShipsList",
				http.MethodGet,
				"/v1/centcom/ship/all",
				api_v1.ShipsList,
			},

			api.Route{
				"StationRegister",
				http.MethodPost,
				"/v1/centcom/station/register",
				api_v1.StationRegister,
			},

			api.Route{
				"StationsList",
				http.MethodGet,
				"/v1/centcom/station/all",
				api_v1.StationsList,
			},

			api.Route{
				"ShipLand",
				http.MethodPost,
				"/v1/shipping-station/land",
				api_v1.ShipLand,
			},

			api.Route{
				"ShipRequestLanding",
				http.MethodPost,
				"/v1/shipping-station/request-landing",
				api_v1.ShipRequestLanding,
			},
		}

		// Here we leverage the fact that addsvc comes with a constructor for an
		// HTTP handler, and just install it under a particular path prefix in
		// our router.

		//r.PathPrefix("/addsvc").Handler(http.StripPrefix("/addsvc", addtra nsport.NewHTTPHandler(endpoints, tracer, zipkinTracer, logger)))
	}

	// // stringsvc routes.
	// {
	// 	// addsvc had lots of nice importable Go packages we could leverage.
	// 	// With stringsvc we are not so fortunate, it just has some endpoints
	// 	// that we assume will exist. So we have to write that logic here. This
	// 	// is by design, so you can see two totally different methods of
	// 	// proxying to a remote service.

	// 	var (
	// 		tags        = []string{}
	// 		passingOnly = true
	// 		uppercase   endpoint.Endpoint
	// 		count       endpoint.Endpoint
	// 		instancer   = consulsd.NewInstancer(client, logger, "stringsvc", tags, passingOnly)
	// 	)
	// 	{
	// 		factory := stringsvcFactory(ctx, "GET", "/uppercase")
	// 		endpointer := sd.NewEndpointer(instancer, factory, logger)
	// 		balancer := lb.NewRoundRobin(endpointer)
	// 		retry := lb.Retry(*retryMax, *retryTimeout, balancer)
	// 		uppercase = retry
	// 	}
	// 	{
	// 		factory := stringsvcFactory(ctx, "GET", "/count")
	// 		endpointer := sd.NewEndpointer(instancer, factory, logger)
	// 		balancer := lb.NewRoundRobin(endpointer)
	// 		retry := lb.Retry(*retryMax, *retryTimeout, balancer)
	// 		count = retry
	// 	}

	// 	// We can use the transport/http.Server to act as our handler, all we
	// 	// have to do provide it with the encode and decode functions for our
	// 	// stringsvc methods.

	// 	r.Handle("/stringsvc/uppercase", httptransport.NewServer(uppercase, decodeUppercaseRequest, encodeJSONResponse))
	// 	r.Handle("/stringsvc/count", httptransport.NewServer(count, decodeCountRequest, encodeJSONResponse))
	// }

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errc <- http.ListenAndServe(*httpAddr, r)
	}()

	// Run!
	logger.Log("exit", <-errc)
}

// func addsvcFactory(makeEndpoint func(auth_dbsvc_service.UserManager) endpoint.Endpoint, tracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) sd.Factory {
// 	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
// 		// We could just as easily use the HTTP or Thrift client package to make
// 		// the connection to addsvc. We've chosen gRPC arbitrarily. Note that
// 		// the transport is an implementation detail: it doesn't leak out of
// 		// this function. Nice!

// 		conn, err := grpc.Dial(instance, grpc.WithInsecure())
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		service := addtransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)
// 		endpoint := makeEndpoint(service)

// 		// Notice that the addsvc gRPC client converts the connection to a
// 		// complete addsvc, and we just throw away everything except the method
// 		// we're interested in. A smarter factory would mux multiple methods
// 		// over the same connection. But that would require more work to manage
// 		// the returned io.Closer, e.g. reference counting. Since this is for
// 		// the purposes of demonstration, we'll just keep it simple.

// 		return endpoint, conn, nil
// 	}
// }

// func stringsvcFactory(ctx context.Context, method, path string) sd.Factory {
// 	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
// 		if !strings.HasPrefix(instance, "http") {
// 			instance = "http://" + instance
// 		}
// 		tgt, err := url.Parse(instance)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		tgt.Path = path

// 		// Since stringsvc doesn't have any kind of package we can import, or
// 		// any formal spec, we are forced to just assert where the endpoints
// 		// live, and write our own code to encode and decode requests and
// 		// responses. Ideally, if you write the service, you will want to
// 		// provide stronger guarantees to your clients.

// 		var (
// 			enc httptransport.EncodeRequestFunc
// 			dec httptransport.DecodeResponseFunc
// 		)
// 		switch path {
// 		case "/uppercase":
// 			enc, dec = encodeJSONRequest, decodeUppercaseResponse
// 		case "/count":
// 			enc, dec = encodeJSONRequest, decodeCountResponse
// 		default:
// 			return nil, nil, fmt.Errorf("unknown stringsvc path %q", path)
// 		}

// 		return httptransport.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
// 	}
// }

func encodeJSONRequest(_ context.Context, req *http.Request, request interface{}) error {
	// Both uppercase and count requests are encoded in the same way:
	// simple JSON serialization to the request body.
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// I've just copied these functions from stringsvc3/transport.go, inlining the
// struct definitions.

// func decodeUppercaseResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
// 	var response struct {
// 		V   string `json:"v"`
// 		Err string `json:"err,omitempty"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return nil, err
// 	}
// 	return response, nil
// }

// func decodeCountResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
// 	var response struct {
// 		V int `json:"v"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return nil, err
// 	}
// 	return response, nil
// }

// func decodeUppercaseRequest(ctx context.Context, req *http.Request) (interface{}, error) {
// 	var request struct {
// 		S string `json:"s"`
// 	}
// 	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
// 		return nil, err
// 	}
// 	return request, nil
// }

// func decodeCountRequest(ctx context.Context, req *http.Request) (interface{}, error) {
// 	var request struct {
// 		S string `json:"s"`
// 	}
// 	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
// 		return nil, err
// 	}
// 	return request, nil
// }

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
