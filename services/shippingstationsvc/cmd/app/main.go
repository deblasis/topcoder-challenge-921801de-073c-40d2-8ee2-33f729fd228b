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
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	consulreg "deblasis.net/space-traffic-control/common/consul"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	cce "deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	ccs "deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
	cct "deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/transport"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/internal/acl"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/service"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/transport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hashicorp/consul/api"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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


	var (
		logger       = cfg.Logger
		retryMax     = cfg.APIGateway.RetryMax
		retryTimeout = cfg.APIGateway.RetryTimeoutMs * int(time.Millisecond)
		tags         = []string{""}
		passingOnly  = true
		cc_endpoints = cce.EndpointSet{}
		instancer    sd.Instancer
	)
	if cfg.BindOnLocalhost {
		instancer = sd.FixedInstancer{"localhost:9482"} //TODO from config
	} else {
		instancer = consulsd.NewInstancer(client, logger, ccs.ServiceName, tags, passingOnly)
	}
	instancesChannel := make(chan sd.Event)
	go func() {
		for event := range instancesChannel {
			if len(event.Instances) > 0 {
				logger.Log("received_instances", strings.Join(event.Instances, ","))
				return
			}
		}
	}()

	{
		factory := centralCommandServiceFactory(cce.MakeGetNextAvailableDockingStationEndpoint, cfg, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, time.Duration(retryTimeout), balancer)
		cc_endpoints.GetNextAvailableDockingStationEndpoint = retry
	}
	{
		factory := centralCommandServiceFactory(cce.MakeRegisterShipLandingEndpoint, cfg, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, time.Duration(retryTimeout), balancer)
		cc_endpoints.RegisterShipLandingEndpoint = retry
	}

	

	logger.Log("retryMax", retryMax, "retryTimeout", retryTimeout)

	var (
		g group.Group

		jwtHandler                = auth.NewJwtHandler(log.With(cfg.Logger, "component", "JwtHandler"), cfg.JWT)
		grpcServerAuthInterceptor = auth.NewAuthServerInterceptor(log.With(cfg.Logger, "component", "AuthServerInterceptor"), jwtHandler, acl.AclRules())

		svc         = service.NewShippingStationService(log.With(cfg.Logger, "component", "CentralCommandService"), cfg.JWT, cc_endpoints)
		eps         = endpoints.NewEndpointSet(svc, log.With(cfg.Logger, "component", "EndpointSet"))
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
				baseServer = grpc.NewServer(
					grpc.UnaryInterceptor(grpcServerAuthInterceptor.Unary()), grpc.Creds(creds),
				)
			} else {
				baseServer = grpc.NewServer(
					grpc.UnaryInterceptor(grpcServerAuthInterceptor.Unary()),
				)
			}
			pb.RegisterShippingStationServiceServer(baseServer, grpcServer)
			reflection.Register(baseServer)

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

func centralCommandServiceFactory(makeEndpoint func(ccs.CentralCommandService) endpoint.Endpoint, cfg config.Config, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		var (
			conn *grpc.ClientConn
			err  error
		)

		if cfg.BindOnLocalhost {
			conn, err = grpc.Dial("localhost:"+strings.Split(instance, ":")[1], grpc.WithInsecure())
		} else {
			conn, err = grpc.Dial(instance, grpc.WithInsecure())
		}

		if err != nil {
			return nil, nil, err
		}
		service := cct.NewGRPCClient(conn, logger)
		endpoint := makeEndpoint(service)
		return endpoint, conn, nil
	}
}
