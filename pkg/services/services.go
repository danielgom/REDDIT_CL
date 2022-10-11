package services

import (
	"context"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/util/errors"
)

// UserService contains all the business logic for the user.
type UserService interface {
	SignUp(context.Context, *internal.RegisterRequest) (*internal.RegisterResponse, errors.CommonError)
	VerifyAccount(context.Context, string) errors.CommonError
	Login(context.Context, *internal.LoginRequest) (*internal.LoginResponse, errors.CommonError)
	RefreshToken(ctx context.Context, request *internal.RefreshTokenRequest) (*internal.RefreshTokenResponse, errors.CommonError)
}

// RefreshTokenService contains all the business logic for the RefreshToken.
type RefreshTokenService interface {
	Create(context.Context) (string, errors.CommonError)
	Validate(context.Context, string) errors.CommonError
}
