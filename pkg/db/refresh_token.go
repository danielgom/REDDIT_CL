package db

import (
	"context"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type refreshTokenRepo struct {
	DB config.DBConn
}

// NewRTRepository creates a new refresh token repository instance.
func NewRTRepository(conn *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepo{DB: conn}
}

// Save persists a new refresh token to the DB.
func (r *refreshTokenRepo) Save(ctx context.Context, token *model.RefreshToken) errors.CommonError {
	exec, err := r.DB.Exec(ctx, `INSERT INTO refresh_token("token", "expires_at") VALUES($1, $2)`,
		token.Token, token.ExpiresAt)

	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}

	if exec.RowsAffected() != 1 {
		return errors.NewInternalServerError("Database error", err)
	}
	return nil
}

// FindByToken looks for a refresh token in the DB.
func (r *refreshTokenRepo) FindByToken(ctx context.Context, token string) (*model.RefreshToken, errors.CommonError) {
	var rT model.RefreshToken

	err := r.DB.QueryRow(ctx, `SELECT expires_at FROM refresh_token WHERE token=$1`, token).Scan(&rT.ExpiresAt)

	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}

	return &rT, nil
}
