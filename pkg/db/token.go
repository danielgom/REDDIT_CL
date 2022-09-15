package db

import (
	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TokenRepo struct {
	DB config.DBConn
}

func NewTokenRepository(conn *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{DB: conn}
}

func (r *TokenRepo) Save(ctx context.Context, token *model.VerificationToken) error {
	exec, err := r.DB.Exec(ctx, `INSERT INTO verification_token("id", "token", "expiry_date") VALUES ($1, $2, $3)`, token.User.Id,
		token.Token, token.ExpiryDate)
	if err != nil {
		return err
	}

	if exec.RowsAffected() != 1 {
		return errors.New("could not save the token")
	}
	return nil
}
