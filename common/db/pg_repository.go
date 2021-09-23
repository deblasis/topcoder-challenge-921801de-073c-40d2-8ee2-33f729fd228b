package db

import pg "github.com/go-pg/pg/v10"

// PostgresClient declare methods for Pg
type PostgresClient interface {

	// GetConnection returns DB connection
	GetConnection() *pg.DB

	// Close closes the database client, releasing any open resources
	Close() error
}
