// The application represents for routing the endpoints
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"deblasis.net/space-traffic-control/common/bootstrap"
	"deblasis.net/space-traffic-control/common/config"
	consulreg "deblasis.net/space-traffic-control/common/consul"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	dbe "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	dbs "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	dbt "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/transport"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/transport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	grpcgokit "github.com/go-kit/kit/transport/grpc"
	"github.com/hashicorp/consul/api"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	var (
		httpAddr   = net.JoinHostPort(cfg.ListenAddr, cfg.HttpServerPort)
		grpcAddr   = net.JoinHostPort(cfg.ListenAddr, cfg.GrpcServerPort)
		consulAddr = net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)
	)

	zipkinTracer, tracer := bootstrap.SetupTracers(cfg, service.ServiceName)

	// var ints, chars metrics.Counter
	// {
	// 	// Business-level metrics.
	// 	ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
	// 		Namespace: service.Namespace,
	// 		Subsystem: service.ServiceName,
	// 		Name:      "integers_summed", //TODO
	// 		Help:      "Total count of integers summed via the Sum method.",
	// 	}, []string{})
	// 	chars = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
	// 		Namespace: service.Namespace,
	// 		Subsystem: service.ServiceName,
	// 		Name:      "characters_concatenated", //TODO
	// 		Help:      "Total count of characters concatenated via the Concat method.",
	// 	}, []string{})
	// }

	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: service.Namespace,
			Subsystem: strings.Split(service.ServiceName, ".")[2],
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	//dependencies
	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if len(consulAddr) > 0 {
			consulConfig.Address = consulAddr
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			cfg.Logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	// Transport domain.
	// var (
	// 	ctx = context.Background()
	// 	r   = mux.NewRouter()
	// )
	// Each method gets constructed with a factory. Factories take an
	// instance string, and return a specific endpoint. In the factory we
	// dial the instance string we get from Consul, and then leverage an
	// addsvc client package to construct a complete service. We can then
	// leverage the addsvc.Make{Sum,Concat}Endpoint constructors to convert
	// the complete service to specific endpoint.

	var (
		logger       = cfg.Logger
		retryMax     = cfg.APIGateway.RetryMax
		retryTimeout = cfg.APIGateway.RetryTimeoutMs * int(time.Millisecond)
		tags         = []string{"authDBService"}
		passingOnly  = true
		db_endpoints = dbe.EndpointSet{}

		instancer = consulsd.NewInstancer(client, logger, dbs.ServiceName, tags, passingOnly)
	)
	instancesChannel := make(chan sd.Event)

	done := make(chan bool, 1)
	go func(ok chan bool) {
		for evt := range instancesChannel {
			for _, i := range evt.Instances {
				logger.Log("received_instance", i)
			}
			if len(evt.Instances) > 0 {
				instancer.Stop()
				instancer = consulsd.NewInstancer(client, logger, dbs.ServiceName, tags, passingOnly)
				ok <- true
				return
			}
			logger.Log("msg", "waiting for instances")
			time.Sleep(time.Second * 1)
		}
	}(done)
	instancer.Register(instancesChannel)
	<-done

	{
		factory := userManagerServiceFactory(dbe.MakeCreateUserEndpoint, cfg, tracer, zipkinTracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, time.Duration(retryTimeout), balancer)
		db_endpoints.CreateUserEndpoint = retry
	}
	{
		factory := userManagerServiceFactory(dbe.MakeGetUserByUsernameEndpoint, cfg, tracer, zipkinTracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, time.Duration(retryTimeout), balancer)
		db_endpoints.GetUserByUsernameEndpoint = retry
	}

	// Here we leverage the fact that addsvc comes with a constructor for an
	// HTTP handler, and just install it under a particular path prefix in
	// our router.

	logger.Log("retryMax", retryMax, "retryTimeout", retryTimeout)

	//r.PathPrefix("/addsvc").Handler(http.StripPrefix("/addsvc", addtransport.NewHTTPHandler(db_endpoints, tracer, zipkinTracer, logger)))

	var (
		g   group.Group
		svc = service.NewAuthService(log.With(cfg.Logger, "component", "AuthService"), cfg.JWT, db_endpoints)

		// svc = service.NewUserManager(log.With(cfg.Logger, "component", "UserManager"))
		eps = endpoints.NewEndpointSet(svc, log.With(cfg.Logger, "component", "EndpointSet"), duration, tracer, zipkinTracer)

		httpHandler = transport.NewHTTPHandler(eps, log.With(cfg.Logger, "component", "HTTPHandler"))
		grpcServer  = transport.NewGRPCServer(eps, log.With(cfg.Logger, "component", "GRPCServer"))
	)
	fmt.Printf("svc %v", svc)

	// consul
	{
		if cfg.Consul.Host != "" && cfg.Consul.Port != "" {
			consulAddres := net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)
			grpcPort, _ := strconv.Atoi(cfg.GrpcServerPort)
			metricsPort, _ := strconv.Atoi(cfg.HttpServerPort)
			tags := []string{service.Namespace, service.ServiceName}
			consulReg := consulreg.NewConsulRegister(consulAddres, service.ServiceName, grpcPort, metricsPort, tags, cfg.Logger, cfg.BindOnLocalhost)
			svcRegistar, err := consulReg.NewConsulGRPCRegister()
			defer svcRegistar.Deregister()
			if err != nil {
				level.Error(cfg.Logger).Log(
					"consulAddres", consulAddres,
					"serviceName", service.ServiceName,
					"grpcPort", grpcPort,
					"metricsPort", metricsPort,
					"tags", tags,
					"err", err,
				)
			}
			svcRegistar.Register()
		}
	}
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			level.Error(cfg.Logger).Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			if cfg.SSL.ServerCert != "" && cfg.SSL.ServerKey != "" {
				level.Debug(cfg.Logger).Log("transport", "HTTP", "addr", httpAddr, "TLS", "enabled")
				return http.ServeTLS(httpListener, httpHandler, cfg.SSL.ServerCert, cfg.SSL.ServerKey)
			} else {
				level.Debug(cfg.Logger).Log("transport", "HTTP", "addr", httpAddr, "TLS", "disabled")
				return http.Serve(httpListener, httpHandler)
			}
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			level.Error(cfg.Logger).Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Debug(cfg.Logger).Log("transport", "gRPC", "addr", grpcAddr)

			var baseServer *grpc.Server
			if cfg.SSL.ServerCert != "" && cfg.SSL.ServerKey != "" {
				creds, err := credentials.NewServerTLSFromFile(cfg.SSL.ServerCert, cfg.SSL.ServerKey)
				if err != nil {
					level.Error(cfg.Logger).Log("serviceName", service.ServiceName, "certificates", creds, "err", err)
					os.Exit(1)
				}
				level.Info(cfg.Logger).Log("serviceName", service.ServiceName, "protocol", "GRPC", "exposed", cfg.GrpcServerPort, "certFile", cfg.SSL.ServerCert, "keyFile", cfg.SSL.ServerKey)
				baseServer = grpc.NewServer(grpc.UnaryInterceptor(grpcgokit.Interceptor), grpc.Creds(creds))
			} else {
				baseServer = grpc.NewServer(grpc.UnaryInterceptor(grpcgokit.Interceptor))
			}
			pb.RegisterAuthServiceServer(baseServer, grpcServer)

			grpc_health_v1.RegisterHealthServer(baseServer, &healthcheck.HealthSvcImpl{})

			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	level.Info(cfg.Logger).Log("exit", g.Run())
}

func userManagerServiceFactory(makeEndpoint func(dbs.UserManager, log.Logger) endpoint.Endpoint, cfg config.Config, tracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		// We could just as easily use the HTTP or Thrift client package to make
		// the connection to addsvc. We've chosen gRPC arbitrarily. Note that
		// the transport is an implementation detail: it doesn't leak out of
		// this function. Nice!

		var (
			conn *grpc.ClientConn
			err  error
		)

		// if cfg.SSL.ServerCert != "" && cfg.SSL.ServerKey != "" {
		// 	creds, err := credentials.NewServerTLSFromFile(cfg.SSL.ServerCert, cfg.SSL.ServerKey)
		// 	if err != nil {
		// 		level.Error(cfg.Logger).Log("client", "userManagerServiceFactory", "certificates", creds, "err", err)
		// 		os.Exit(1)
		// 	}
		// 	level.Info(cfg.Logger).Log("client", "userManagerServiceFactory", "protocol", "GRPC", "certFile", cfg.SSL.ServerCert, "keyFile", cfg.SSL.ServerKey)

		// 	conn, err = grpc.Dial(instance, grpc.WithTransportCredentials(creds))
		// 	if err != nil {
		// 		return nil, nil, err
		// 	}
		// } else {
		conn, err = grpc.Dial(instance, grpc.WithInsecure())
		//}

		if err != nil {
			return nil, nil, err
		}
		service := dbt.NewGRPCClient(conn, tracer, zipkinTracer, logger)
		endpoint := makeEndpoint(service, logger)
		level.Debug(logger).Log(
			"method", "userManagerServiceFactory",
			"instance", instance,
			"conn", conn,
		)

		// Notice that the addsvc gRPC client converts the connection to a
		// complete addsvc, and we just throw away everything except the method
		// we're interested in. A smarter factory would mux multiple methods
		// over the same connection. But that would require more work to manage
		// the returned io.Closer, e.g. reference counting. Since this is for
		// the purposes of demonstration, we'll just keep it simple.

		return endpoint, conn, nil
	}
}
