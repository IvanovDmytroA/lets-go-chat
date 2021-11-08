package main

import (
	"log"
	"net/http"

	"github.com/IvanovDmytroA/lets-go-chat/internal/handler"
)

func main() {
	http.HandleFunc("/v1", handler.PageViewHandler)
	http.HandleFunc("/v1/user", handler.CreateUser)
	http.HandleFunc("/v1/user/login", handler.LoginUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
