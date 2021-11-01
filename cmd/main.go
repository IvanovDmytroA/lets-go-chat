package main

import (
	"fmt"

	"github.com/IvanovDmytroA/lets-go-chat/pkg/hasher"
)

func main() {
	checkPasswordHash("password")
}

func checkPasswordHash(password string) {
	hash, err := hasher.HashPassword(password)

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash: ", hash)

	match := hasher.CheckPasswordHash(password, hash)

	fmt.Println("Match: ", match)
}
