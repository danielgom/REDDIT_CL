package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const passwordHashDifficulty = 10

// HashPassword produces a encrypted password.
func HashPassword(password string) (string, error) {
	hPass, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashDifficulty)
	if err != nil {
		return "", fmt.Errorf("could encrypt password %w", err)
	}
	return string(hPass), nil
}

// CheckPasswordHash checks whether the password provided is the correct one.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
