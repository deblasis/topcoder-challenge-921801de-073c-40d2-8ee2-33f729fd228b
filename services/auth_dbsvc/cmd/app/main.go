// The application represents for routing the endpoints
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/repositories"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service/db"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service/endpoints"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/transport"

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

	httpAddr := net.JoinHostPort("localhost", cfg.HttpServerPort)
	// grpcAddr := net.JoinHostPort("localhost", cfg.GrpcServerPort)

	level.Debug(cfg.Logger).Log("DB address", cfg.DbConfig.Address)

	pgClient := db.NewPostgresClientFromConfig(cfg)
	connection := pgClient.GetConnection()
	defer connection.Close()

	repo := repositories.NewUserRepository(connection, log.With(cfg.Logger, "component", "UserRepository"))

	svc := service.NewUserManager(repo, log.With(cfg.Logger, "component", "UserManager"))
	eps := endpoints.NewEndpointSet(svc, log.With(cfg.Logger, "component", "EndpointSet"))
	httpHandler := transport.NewHTTPHandler(eps, log.With(cfg.Logger, "component", "HTTPHandler"))

	var g group.Group

	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			level.Error(cfg.Logger).Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			level.Debug(cfg.Logger).Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})

	}
	{
		// This function just sits and waits for ctrl-C.
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
