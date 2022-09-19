package db

import (
	"context"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
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
func (r *tokenRepo) Save(ctx context.Context, token *model.VerificationToken) errors.CommonError {
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

func (r *tokenRepo) FindByToken(ctx context.Context, token string) (*model.VerificationToken, errors.CommonError) {
	var verToken model.VerificationToken
	var user model.User

	row := r.DB.QueryRow(ctx, `SELECT * FROM verification_token t JOIN users u on u.id = t.id WHERE t.token=$1`, token)
	err := row.Scan(&verToken.ID, &verToken.Token, &verToken.ExpiryDate, &user.ID, &user.Username,
		&user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Enabled)

	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}

	verToken.User = &user

	return &verToken, nil
}
