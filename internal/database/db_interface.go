package database

import "database/sql"

// DBClient adalah kontrak database agar bisa diganti driver kapan saja
type DBClient interface {
	GetDB() *sql.DB
	HealthCheck() error
	Close() error
}
