// Package internal has internal structs for api functionality
package internal

import (
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"RD-Clone-API/pkg/model"
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

// RegisterResponse is the struct for a successful signUp.
type RegisterResponse struct {
	Username  string
	Email     string
	CreatedAt time.Time
	Enabled   bool
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

// BuildRegisterResponse builds the output of the signUp response when is not error -ed.
func BuildRegisterResponse(user *model.User) *RegisterResponse {
	var response RegisterResponse

	response.Username = user.Username
	response.Email = user.Email
	response.CreatedAt = user.CreatedAt
	response.Enabled = user.Enabled

	return &response
}
