// The application represents for routing the endpoints
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"deblasis.net/space-traffic-control/common/bootstrap"
	"deblasis.net/space-traffic-control/common/config"
	consulreg "deblasis.net/space-traffic-control/common/consul"
	"deblasis.net/space-traffic-control/common/db"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	grpcgokit "github.com/go-kit/kit/transport/grpc"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	var (
		httpAddr = net.JoinHostPort(cfg.ListenAddr, cfg.HttpServerPort)
		grpcAddr = net.JoinHostPort(cfg.ListenAddr, cfg.GrpcServerPort)
	)

	level.Debug(cfg.Logger).Log("DB address", cfg.Db.Address)
	var (
		pgClient   = db.NewPostgresClientFromConfig(cfg)
		connection = pgClient.GetConnection()
	)
	defer connection.Close()

	zipkinTracer, tracer := bootstrap.SetupTracers(cfg, service.ServiceName)

	// var ints, chars metrics.Counter
	// {
	// 	// Business-level metrics.
	// 	ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
	// 		Namespace: service.Namespace,
	// 		Subsystem: strings.Split(service.ServiceName, ".")[2],
	// 		Name:      "integers_summed", //TODO
	// 		Help:      "Total count of integers summed via the Sum method.",
	// 	}, []string{})
	// 	chars = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
	// 		Namespace: service.Namespace,
	// 		Subsystem: strings.Split(service.ServiceName, ".")[2],
	// 		Name:      "characters_concatenated", //TODO
	// 		Help:      "Total count of characters concatenated via the Concat method.",
	// 	}, []string{})
	// }

	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: service.Namespace,
			Subsystem: strings.Split(service.ServiceName, "-")[2],
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	var (
		g group.Group

		shipRepo    = repositories.NewShipRepository(connection, log.With(cfg.Logger, "component", "ShipRepository"))
		stationRepo = repositories.NewStationRepository(connection, log.With(cfg.Logger, "component", "StationRepository"))
		svc         = service.NewCentralCommandDBService(shipRepo, stationRepo, log.With(cfg.Logger, "component", "CentralCommandDBService"))
		eps         = endpoints.NewEndpointSet(svc, log.With(cfg.Logger, "component", "EndpointSet"), duration, tracer, zipkinTracer)

		httpHandler = transport.NewHTTPHandler(eps, log.With(cfg.Logger, "component", "HTTPHandler"))
		grpcServer  = transport.NewGRPCServer(eps, log.With(cfg.Logger, "component", "GRPCServer"))
	)

	// consul
	{
		if cfg.Consul.Host != "" && cfg.Consul.Port != "" {
			consulAddres := net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)
			grpcPort, _ := strconv.Atoi(cfg.GrpcServerPort)
			metricsPort, _ := strconv.Atoi(cfg.HttpServerPort)
			tags := []string{service.Namespace, service.ServiceName, "centralCommandDBService", "DBService"}
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
			pb.RegisterCentralCommandDBServiceServer(baseServer, grpcServer)

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

// func addsvcFactory(makeEndpoint func(addservice.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) sd.Factory {
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
