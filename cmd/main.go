package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/IvanovDmytroA/lets-go-chat/internal/server"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	server.Start()
}
