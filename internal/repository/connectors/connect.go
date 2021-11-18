package repository

import (
	"database/sql"
)

// Database worker interface
type Worker interface {
	Init(*sql.DB)
	Get() *sql.DB
}

// Database connection pool
type Pool struct {
	DB *sql.DB
}
