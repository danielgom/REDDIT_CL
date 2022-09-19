package db

import (
	"context"

	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
)

// UserRepository serves as a middleware to call our userDB.
type UserRepository interface {
	FindByUsername(context.Context, string) (*model.User, errors.CommonError)
	FindByEmail(context.Context, string) (*model.User, errors.CommonError)
	Save(context.Context, *model.User) (*model.User, errors.CommonError)
	Update(context.Context, *model.User) errors.CommonError
}

// TokenRepository serves as a middleware to call our tokenDB.
type TokenRepository interface {
	Save(context.Context, *model.VerificationToken) errors.CommonError
	FindByToken(context.Context, string) (*model.VerificationToken, errors.CommonError)
}
