package main

import (
	"sync"

	"github.com/IvanovDmytroA/lets-go-chat/internal/server"
)

func main() {
	once := sync.Once{}
	once.Do(func() {
		server.Start()
	})
}
