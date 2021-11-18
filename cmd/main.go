package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IvanovDmytroA/lets-go-chat/internal/configuration"
	handler "github.com/IvanovDmytroA/lets-go-chat/internal/handler/transport"
	"github.com/IvanovDmytroA/lets-go-chat/internal/repository"
	connectors "github.com/IvanovDmytroA/lets-go-chat/internal/repository/connectors"
	"github.com/gorilla/mux"
)

func main() {
	port := initServer()
	env, err := configuration.InitEnv()
	if err != nil {
		log.Fatal("Failed to init environment configuration")
	}

	initDb(env)

	r := mux.NewRouter()
	r.HandleFunc("/v1", handler.PageViewHandler)
	r.HandleFunc("/v1/user", handler.CreateUser)
	r.HandleFunc("/v1/user/login", handler.LoginUser)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initServer() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func initDb(e *configuration.Env) {
	switch e.DataBase.Type {
	case "postgres":
		var worker connectors.Worker = &connectors.PostgresWorker{}
		var dbUrl = os.Getenv("DATABASE_URL")
		if dbUrl == "" {
			dbUrl = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				e.DataBase.Host, e.DataBase.Port, e.DataBase.User, e.DataBase.Password, e.DataBase.Name)
		}

		dbc, err := sql.Open(e.DataBase.Type, dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		worker.Init(dbc)
		repository.InitUserRepository(&worker)
	}
}
