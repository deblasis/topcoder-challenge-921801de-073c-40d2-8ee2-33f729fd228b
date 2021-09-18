package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"github.com/go-kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
)

//this guy connects to the centralcommanddb and runs a query every x seconds, logging the output, that's it

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	level.Debug(cfg.Logger).Log("DB address", cfg.Db.Address)
	var (
		pgClient   = db.NewPostgresClientFromConfig(cfg)
		connection = pgClient.GetConnection()
	)
	defer connection.Close()

	cfg.Logger = level.NewFilter(cfg.Logger, level.AllowInfo())

	stats := &stats{
		Status: "OK",
	}

	r := mux.NewRouter()

	ctx, cancel := context.WithCancel(context.Background())
	cancelInterrupt := make(chan struct{})
	var (
		g group.Group
	)
	{
		g.Add(func() error {
			cockoo := time.NewTicker(2 * time.Second) //TODO cfg
			level.Info(cfg.Logger).Log("msg", "⏳ checking how many ships left... 🔍")

			return func() error {
				for {
					select {
					case <-cockoo.C:

						var ships_have_left int

						level.Debug(cfg.Logger).Log("msg", "⏳ checking how many ships left... 🔍")

						_, err := connection.WithContext(ctx).QueryOne(
							pg.Scan(&ships_have_left), "select ships_have_left()",
						)
						if err != nil {
							level.Error(cfg.Logger).Log("err", err)
							return err
						}
						if ships_have_left > 0 {

							stats.LastTimeShipsLeft = time.Now().UTC()
							stats.LastTimeNumberShipsLeft = ships_have_left

							level.Info(cfg.Logger).Log("msg", fmt.Sprintf("🚀 %v ships left", ships_have_left))
						} else {
							level.Debug(cfg.Logger).Log("msg", fmt.Sprintf("🚀 %v ships left", ships_have_left))
						}
					case <-cancelInterrupt:
						cockoo.Stop()
						cancel()
						return nil
					}
				}
				return nil
			}()

		}, func(e error) {
			level.Warn(cfg.Logger).Log("cancelling", e)
		})
	}
	{
		g.Add(func() error {
			r.Handle("/health", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
				rw.Header().Set("Content-Type", "application/json")

				json.NewEncoder(rw).Encode(&stats)

			}))
			return http.ListenAndServe(":9500", r)
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	{

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

type stats struct {
	Status                  string    `json:"status,omitempty"`
	LastTimeShipsLeft       time.Time `json:"last_time_ships_left,omitempty"`
	LastTimeNumberShipsLeft int       `json:"number_ships_left_last_time,omitempty"`
}
