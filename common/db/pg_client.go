package db

import (
	"context"
	"fmt"
	"net"
	"time"

	"deblasis.net/space-traffic-control/common/config"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/go-pg/pg/v10"
)

const (
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
	PoolSize     = 10
	MinIdleConns = 10
	MaxRetries   = 1
)

type postgresClient struct {
	Db *pg.DB
}

func (p postgresClient) GetConnection() *pg.DB {
	return p.Db
}

func (p postgresClient) Close() error {
	return p.Db.Close()
}

func NewPostgresClientFromConfig(config config.Config) PostgresClient {
	return NewPostgresClientFromPgOptions((config.Logger), GetPgConnectionOptions(config))
}

func NewPostgresClientFromPgOptions(logger log.Logger, pgOptions *pg.Options) PostgresClient {
	level.Debug(logger).Log("msg", fmt.Sprintf("Trying to connect to %v", pgOptions.Addr))
	db := pg.Connect(pgOptions)
	return postgresClient{
		Db: db,
	}
}

// NewPostgresClient returns a PostgresClient
func NewPostgresClient(db *pg.DB) PostgresClient {
	return postgresClient{
		Db: db,
	}
}

// GetPgConnectionOptions returns pg Options based on config
func GetPgConnectionOptions(config config.Config) *pg.Options {
	return &pg.Options{
		ApplicationName: "deblasis_spaceTrafficControl",
		Addr:            config.Db.Address,
		User:            config.Db.Username,
		Password:        config.Db.Password,
		Database:        config.Db.Database,
		ReadTimeout:     ReadTimeout,
		WriteTimeout:    WriteTimeout,
		PoolSize:        PoolSize,
		MinIdleConns:    MinIdleConns,
		MaxRetries:      MaxRetries,
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, 5*time.Second)
			if err != nil {
				return nil, err
			}
			return conn, conn.(*net.TCPConn).SetKeepAlive(true)
		},
	}
}
