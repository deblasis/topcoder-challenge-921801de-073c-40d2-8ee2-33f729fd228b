// The application represents for run migrations

package main

import (
	"flag"
	"os"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service/db"
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
			cfg.Logger.Fatal(errors.Wrap(err, "Could not init migrations"))
		}
	}

	// scan the dir for files with .sql extension and adds  migrations to the collection
	err = migrationCollection.DiscoverSQLMigrations(*migrationDir)
	if err != nil {
		cfg.Logger.Fatal(errors.Wrap(err, "Failed to read migrations"))
	}

	_, _, err = migrationCollection.Run(connection, "up")
	if err != nil {
		cfg.Logger.Fatal(errors.Wrap(err, "Could not migrate"))
	}
	cfg.Logger.Info("migrated successfully!")
}
