// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// The application represents for routing the endpoints
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"deblasis.net/space-traffic-control/common/config"
	consulreg "deblasis.net/space-traffic-control/common/consul"
	"deblasis.net/space-traffic-control/common/db"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/transport"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	grpcgokit "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
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
		httpAddr    = net.JoinHostPort(cfg.ListenAddr, cfg.HttpServerPort)
		grpcAddr    = net.JoinHostPort(cfg.ListenAddr, cfg.GrpcServerPort)
		auxGrpcAddr = net.JoinHostPort(cfg.ListenAddr, cfg.AuxGrpcServerPort)
	)

	level.Debug(cfg.Logger).Log("DB address", cfg.Db.Address)
	var (
		pgClient   = db.NewPostgresClientFromConfig(cfg)
		connection = pgClient.GetConnection()
	)
	defer connection.Close()

	var (
		g group.Group

		shipRepo    = repositories.NewShipRepository(connection, log.With(cfg.Logger, "component", "ShipRepository"))
		stationRepo = repositories.NewStationRepository(connection, log.With(cfg.Logger, "component", "StationRepository"))
		dockRepo    = repositories.NewDockRepository(connection, log.With(cfg.Logger, "component", "DockRepository"))
		svc         = service.NewCentralCommandDBService(shipRepo, stationRepo, dockRepo, log.With(cfg.Logger, "component", "CentralCommandDBService"))
		eps         = endpoints.NewEndpointSet(svc, log.With(cfg.Logger, "component", "EndpointSet"))

		httpHandler = transport.NewHTTPHandler(eps, log.With(cfg.Logger, "component", "HTTPHandler"))
		grpcServer  = transport.NewGRPCServer(eps, log.With(cfg.Logger, "component", "GRPCServer"))

		auxRepo       = repositories.NewAuxRepository(connection, log.With(cfg.Logger, "component", "AuxRepository"))
		auxsvc        = service.NewCentralCommandDBAuxService(auxRepo, log.With(cfg.Logger, "component", "CentralCommandDBAuxService"))
		auxeps        = endpoints.NewAuxEndpointSet(auxsvc, log.With(cfg.Logger, "component", "AuxEndpointSet"))
		auxGrpcServer = transport.NewAuxGrpcServer(auxeps, log.With(cfg.Logger, "component", "AUXGRPCServer"))
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
		//this acts like a feature flag, this endpoint is used only in integration-testing
		if auxGrpcAddr != "" {
			grpcListener, err := net.Listen("tcp", auxGrpcAddr)
			if err != nil {
				level.Error(cfg.Logger).Log("transport", "aux-gRPC", "during", "Listen", "err", err)
				os.Exit(1)
			}
			g.Add(func() error {
				level.Debug(cfg.Logger).Log("transport", "aux-gRPC", "addr", auxGrpcAddr)

				var baseServer *grpc.Server
				if cfg.SSL.ServerCert != "" && cfg.SSL.ServerKey != "" {
					creds, err := credentials.NewServerTLSFromFile(cfg.SSL.ServerCert, cfg.SSL.ServerKey)
					if err != nil {
						level.Error(cfg.Logger).Log("serviceName", service.AuxServiceName, "certificates", creds, "err", err)
						os.Exit(1)
					}
					level.Info(cfg.Logger).Log("serviceName", service.AuxServiceName, "protocol", "GRPC", "exposed", cfg.AuxGrpcServerPort, "certFile", cfg.SSL.ServerCert, "keyFile", cfg.SSL.ServerKey)
					baseServer = grpc.NewServer(grpc.UnaryInterceptor(grpcgokit.Interceptor), grpc.Creds(creds))
				} else {
					baseServer = grpc.NewServer(grpc.UnaryInterceptor(grpcgokit.Interceptor))
				}
				pb.RegisterCentralCommandDBAuxServiceServer(baseServer, auxGrpcServer)

				grpc_health_v1.RegisterHealthServer(baseServer, &healthcheck.HealthSvcImpl{})

				return baseServer.Serve(grpcListener)
			}, func(error) {
				grpcListener.Close()
			})
		}
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
