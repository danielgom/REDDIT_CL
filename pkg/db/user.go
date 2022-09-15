package db

import (
	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	DB config.DBConn
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{DB: conn}
}

func (r *UserRepo) FindByUsername(ctx context.Context, uName string) (*model.User, error) {
	query := `SELECT * FROM users WHERE username=$1`

	var user *model.User

	err := r.DB.QueryRow(ctx, query, uName).Scan(&user.Id, &user.Username, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.Enabled)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r *UserRepo) Save(ctx context.Context, user *model.User) (*model.User, error) {
	row := r.DB.QueryRow(ctx, `INSERT INTO users("username", "password", "email", "created_at", "updated_at", "enabled") 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, user.Username, user.Password, user.Email,
		user.CreatedAt, user.UpdatedAt, user.Enabled)

	err := row.Scan(&user.Id)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not save user %s", err.Error()))
	}

	return user, nil
}
