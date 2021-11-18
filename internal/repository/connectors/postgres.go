package repository

import (
	"database/sql"
	"log"
)

// Postgres worker
type PostgresWorker struct {
	Pool
}

// Init Postgres database
func (p *PostgresWorker) Init(db *sql.DB) {
	p.Pool.DB = db
	_, err := p.DB.Exec("create database gochat")
	if err != nil {
		log.Println("Database already exists")
	}

	createTables(p)
}

// Returns Postgres connection pool
func (p *PostgresWorker) Get() *sql.DB {
	return p.DB
}

func createTables(p *PostgresWorker) {
	_, err := p.DB.Exec(`create table if not exists users
		(
			id serial primary key,
			username varchar,
			password varchar
		);
		create unique index if not exists users_id
			on users (id);
		create unique index if not exists users_username
			on users (username);
		`)

	if err != nil {
		log.Fatal(err)
	}
}
