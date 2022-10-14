package db

import (
	"context"

	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
)

const errNotFound = "no rows in result set"

// UserRepository serves as a middleware to call our users table.
type UserRepository interface {
	FindByUsername(context.Context, string) (*model.User, errors.CommonError)
	FindByEmail(context.Context, string) (*model.User, errors.CommonError)
	Save(context.Context, *model.User) (*model.User, errors.CommonError)
	Update(context.Context, *model.User) errors.CommonError
}

// TokenRepository serves as a middleware to call our verification_token table.
type TokenRepository interface {
	Save(context.Context, *model.VerificationToken) errors.CommonError
	FindByToken(context.Context, string) (*model.VerificationToken, errors.CommonError)
}

// RefreshTokenRepository serves as a middleware to call our refresh_token table.
type RefreshTokenRepository interface {
	Save(context.Context, *model.RefreshToken) errors.CommonError
	FindByToken(context.Context, string) (*model.RefreshToken, errors.CommonError)
}

// SubredditRepository serves as a middleware to call our subreddit table.
type SubredditRepository interface {
	Save(context.Context, *model.Subreddit) (*model.Subreddit, errors.CommonError)
	FindByID(context.Context, int) (*model.Subreddit, errors.CommonError)
	FindAll(ctx context.Context) ([]*model.Subreddit, errors.CommonError)
}
