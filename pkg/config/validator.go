// Package config is the configuration package of all our tools
package config

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// CustomValidator sets the validator for struct fields.
type CustomValidator struct {
	Validator *validator.Validate
}

// GetValidator gets a new validator instance.
func GetValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}

// AddValidators registers all validators for struct fields.
func AddValidators(v *validator.Validate) error {
	return v.RegisterValidation("password", ValidatePassword)
}

// Validate validates all custom/non-custom validators of struct fields.
func (c *CustomValidator) Validate(i any) error {
	if err := c.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// ValidatePassword is a custom password validator.
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	hasNumber = strings.ContainsAny(password, "123456789")
	hasUpper = strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTVWXYZ")
	hasLower = strings.ContainsAny(password, "abcdefghijklmnopqrstvwxyz")

	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
