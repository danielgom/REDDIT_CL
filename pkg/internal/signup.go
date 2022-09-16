// Package internal has internal structs for api functionality
package internal

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
)

var (
	errInvalidPassword = errors.New("invalid password")
	errInvalidEmail    = errors.New("invalid email")
	errEmptyUsername   = errors.New("username should not be empty")
)

// RegisterRequest comes from the signup request.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Validate validates upcoming RegisterRequest.
func (r *RegisterRequest) Validate() error {
	if r.Username == "" {
		return errEmptyUsername
	}

	if !isValidPassword(r.Password) {
		return errInvalidPassword
	}

	if !isValidEmail(r.Email) {
		return errInvalidEmail
	}

	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(s string) bool {
	const minimumPasswordLength = 7

	var (
		hasMinLen  bool
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	if len(s) >= minimumPasswordLength {
		hasMinLen = true
	}

	hasNumber = strings.ContainsAny(s, "123456789")
	hasUpper = strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTVWXYZ")
	hasLower = strings.ContainsAny(s, "abcdefghijklmnopqrstvwxyz")

	for _, char := range s {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
