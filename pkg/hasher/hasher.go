package hasher

import "golang.org/x/crypto/bcrypt"

const hashCost int = 14

// HashPassword generates the bcrypt hash of the password.
// Returns a hashed password and nil or returns error when a provided string is empty.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password and bcrypt hash.
// Returns true if hashed password have the same hash as provided bcrypt hash, otherwise returns false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
