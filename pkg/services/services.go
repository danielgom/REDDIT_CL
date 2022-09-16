package services

import (
	"RD-Clone-API/pkg/internal"
	"context"
)

// UserService contains all the business logic for our user API.
type UserService interface {
	SignUp(context.Context, *internal.RegisterRequest) error
}
