package helpers

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	// Hashing the login req pass
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {

		// This will tell you if the hash is too short, has a bad prefix, or just mismatched.
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			fmt.Println("DEBUG: bcrypt comparison failed - Passwords do not match.")
		} else {
			// This catches issues like invalid hash format (too short, bad prefix)
			fmt.Printf("DEBUG: bcrypt comparison failed with structural error: %v\n", err)
		}
		return false
	}

	// If err is nil, the passwords match.
	return true
}
