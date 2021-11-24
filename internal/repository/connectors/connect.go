package repository

import (
	"github.com/uptrace/bun"
)

// Database worker interface
type Worker interface {
	Init(*bun.DB)
	Get() *bun.DB
}

// Database connection pool
type Pool struct {
	DB *bun.DB
}
