package db

import (
	"context"

	"RD-Clone-API/pkg/model"
)

// UserRepository serves as a middleware to call our userDB.
type UserRepository interface {
	FindByUsername(context.Context, string) (*model.User, error)
	Save(context.Context, *model.User) (*model.User, error)
}

// TokenRepository serves as a middleware to call our tokenDB.
type TokenRepository interface {
	Save(context.Context, *model.VerificationToken) error
}
