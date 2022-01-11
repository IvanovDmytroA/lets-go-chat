package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestGet(t *testing.T) {
	var worker Worker = &PostgresWorker{}
	var dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", "postgres", "pass", "localhost", 5432)
	}

	dbc := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	db := bun.NewDB(dbc, pgdialect.New())
	worker.Init(db)

	con := worker.Get()

	if con == nil {
		t.Fatalf("Failed to obtain new connection")
	}
}
