//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// The application represents for run migrations

package main

import (
	"flag"
	"os"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"github.com/go-kit/kit/log/level"
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
