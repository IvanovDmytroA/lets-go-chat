package main

import (
	"log"
	"net/http"
	"os"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1", handler.PageViewHandler)
	r.HandleFunc("/v1/user", handler.CreateUser)
	r.HandleFunc("/v1/user/login", handler.LoginUser)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
