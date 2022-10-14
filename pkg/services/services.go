package services

import (
	"context"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/util/errors"
)

// UserService contains all the business logic for the user.
type UserService interface {
	SignUp(context.Context, *internal.RegisterRequest) (*internal.RegisterResponse, errors.CommonError)
	Get(context.Context, string) (*internal.UserResponse, errors.CommonError)
	VerifyAccount(context.Context, string) errors.CommonError
	Login(context.Context, *internal.LoginRequest) (*internal.LoginResponse, errors.CommonError)
	RefreshToken(ctx context.Context, request *internal.RefreshTokenRequest) (*internal.RefreshTokenResponse, errors.CommonError)
}

// RefreshTokenService contains all the business logic for the RefreshToken service.
type RefreshTokenService interface {
	Create(context.Context) (string, errors.CommonError)
	Validate(context.Context, string) errors.CommonError
}

// SubredditService contains all business logic for the Subreddit service.
type SubredditService interface {
	Create(context.Context, *internal.NewSubreddit, string) (*internal.SubredditResponse, errors.CommonError)
	Get(context.Context, int) (*internal.SubredditResponse, errors.CommonError)
	GetAll(context.Context) ([]*internal.SubredditResponse, errors.CommonError)
}
