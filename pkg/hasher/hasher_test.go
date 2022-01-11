package hasher

import (
	"fmt"
	"testing"
)

var pass string = "password"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(pass)
	if err != nil || hash == pass {
		t.Fatalf("Failed to generate hash. Error %s", err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword(pass)
	if !CheckPasswordHash(pass, hash) {
		t.Fatalf("Hash is not equal")
	}
}

func ExampleHashPassword() {
	hash, _ := HashPassword(pass)
	fmt.Println(hash[:7])
	// Output: $2a$14$
}
