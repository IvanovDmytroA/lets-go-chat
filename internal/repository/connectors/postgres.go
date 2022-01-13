package repository

import (
	"log"

	"github.com/uptrace/bun"
)

// Postgres worker
type PostgresWorker struct {
	Pool
}

// Init Postgres database
func (p *PostgresWorker) Init(db *bun.DB) {
	p.Pool.DB = db
	_, err := p.DB.Exec("create database gochat")
	if err != nil {
		log.Println("Database already exists")
	}

	createTables(p)
}

// Returns Postgres connection pool
func (p *PostgresWorker) Get() *bun.DB {
	return p.DB
}

func createTables(p *PostgresWorker) {
	_, err := p.DB.Exec(`create table if not exists users
		(
			id varchar primary key,
			user_name varchar,
			password varchar
		);
		create unique index if not exists users_id
			on users (id);
		create unique index if not exists users_username
			on users (user_name);
		create table if not exists messages
		(
			id int primary key,
			user_id varchar,
			text varchar,
			constraint fk_user
      			foreign key(user_id) 
	  				references users(id)
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
