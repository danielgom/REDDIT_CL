package services

import (
	"context"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/util/errors"
)

// UserService contains all the business logic for our user API.
type UserService interface {
	SignUp(context.Context, *internal.RegisterRequest) (*internal.RegisterResponse, errors.CommonError)
}
