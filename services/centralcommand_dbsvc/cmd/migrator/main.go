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
// The application represents for run migrations

package main

import (
	"flag"
	"os"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"github.com/go-kit/log/level"
	"github.com/go-pg/migrations/v8"
	"github.com/pkg/errors"
)

// User can define another path of migrations
var migrationDir = flag.String("dir", "./scripts/migrations/", "directory with migrations")

// 	true  - perform db init
// 	false - left empty db
var doInit = flag.Bool("init", true, "perform db init (for empty db)")

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	// prepare pg connection
	pgClient := db.NewPostgresClientFromConfig(cfg)
	connection := pgClient.GetConnection()
	defer connection.Close()

	level.Debug(cfg.Logger).Log("database", cfg.Db.Database)

	migrationCollection := migrations.NewCollection()
	if *doInit {
		// perform the DB
		_, _, err := migrationCollection.Run(connection, "init")
		if err != nil {
			level.Error(cfg.Logger).Log("err", errors.Wrap(err, "Could not init migrations"))
			os.Exit(1)
		}
	}

	// scan the dir for files with .sql extension and adds  migrations to the collection
	err = migrationCollection.DiscoverSQLMigrations(*migrationDir)
	if err != nil {
		level.Error(cfg.Logger).Log("err", errors.Wrap(err, "Failed to read migrations"))
		os.Exit(1)
	}

	_, _, err = migrationCollection.Run(connection, "up")
	if err != nil {
		level.Error(cfg.Logger).Log("err", errors.Wrap(err, "Could not migrate"))
		os.Exit(1)
	}
	level.Info(cfg.Logger).Log("msg", "migrated successfully!")
}
