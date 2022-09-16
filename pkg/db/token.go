package db

import (
	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type tokenRepo struct {
	DB config.DBConn
}

// NewTokenRepository creates a new token repository instance.
func NewTokenRepository(conn *pgxpool.Pool) TokenRepository {
	return &tokenRepo{DB: conn}
}

// Save persists a new verification token to the DB.
func (r *tokenRepo) Save(ctx context.Context, token *model.VerificationToken) error {
	exec, err := r.DB.Exec(ctx, `INSERT INTO verification_token("id", "token", "expiry_date") VALUES ($1, $2, $3)`, token.User.ID,
		token.Token, token.ExpiryDate)
	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}

	if exec.RowsAffected() != 1 {
		return errors.NewInternalServerError("Database error", err)
	}
	return nil
}
