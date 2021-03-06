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
package db

import (
	"context"
	"fmt"
	"net"
	"time"

	"deblasis.net/space-traffic-control/common/config"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

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
