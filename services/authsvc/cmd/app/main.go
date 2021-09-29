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
	"deblasis.net/space-traffic-control/common/cache"
	"deblasis.net/space-traffic-control/common/config"
	consulreg "deblasis.net/space-traffic-control/common/consul"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	dbe "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	dbs "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	dbt "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/transport"
	"deblasis.net/space-traffic-control/services/authsvc/internal/acl"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/transport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	consul "github.com/hashicorp/consul/api"
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
		httpAddr   = net.JoinHostPort(cfg.ListenAddr, cfg.HttpServerPort)
		grpcAddr   = net.JoinHostPort(cfg.ListenAddr, cfg.GrpcServerPort)
		consulAddr = net.JoinHostPort(cfg.Consul.Host, cfg.Consul.Port)
	)

	//dependencies
	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := consul.DefaultConfig()
		if len(consulAddr) > 0 {
			consulConfig.Address = consulAddr
		}
		consulClient, err := consul.NewClient(consulConfig)
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
		tags         = []string{"authDBService"}
		passingOnly  = true
		db_endpoints = dbe.EndpointSet{}

		instancer = consulsd.NewInstancer(client, logger, dbs.ServiceName, tags, passingOnly)
	)
	instancesChannel := make(chan sd.Event)

	go func() {
		for event := range instancesChannel {
			if len(event.Instances) > 0 {
				logger.Log("received_instances", strings.Join(event.Instances, ","))
				return
			}
		}
	}()

	instancer.Register(instancesChannel)

	{
		factory := authDBServiceFactory(dbe.MakeCreateUserEndpoint, cfg, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, time.Duration(retryTimeout), balancer)
		db_endpoints.CreateUserEndpoint = retry
	}
	{
		factory := authDBServiceFactory(dbe.MakeGetUserByUsernameEndpoint, cfg, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, time.Duration(retryTimeout), balancer)
		db_endpoints.GetUserByUsernameEndpoint = retry
	}

	var (
		g group.Group

		jwtHandler                = auth.NewJwtHandler(log.With(cfg.Logger, "component", "JwtHandler"), cfg.JWT)
		grpcServerAuthInterceptor = auth.NewAuthServerInterceptor(log.With(cfg.Logger, "component", "AuthServerInterceptor"), jwtHandler, acl.AclRules())

		memTokensCache = cache.NewMemTokensCache()
		svc            = service.NewAuthService(log.With(cfg.Logger, "component", "AuthService"), cfg.JWT, db_endpoints, memTokensCache, jwtHandler)

		eps = endpoints.NewEndpointSet(svc, log.With(cfg.Logger, "component", "EndpointSet"))

		httpHandler = transport.NewHTTPHandler(eps, log.With(cfg.Logger, "component", "HTTPHandler"))
		grpcServer  = transport.NewGRPCServer(log.With(cfg.Logger, "component", "GRPCServer"), eps)
	)

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

func authDBServiceFactory(makeEndpoint func(dbs.AuthDBService) endpoint.Endpoint, cfg config.Config, logger log.Logger) sd.Factory {
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
		service := dbt.NewGRPCClient(conn, logger)
		endpoint := makeEndpoint(service)
		level.Debug(logger).Log(
			"method", "authDBServiceFactory",
			"instance", instance,
			"conn", conn,
		)
		return endpoint, conn, nil
	}
}
