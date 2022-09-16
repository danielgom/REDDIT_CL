// Package db contains all repositories used by this API.
package db

import (
	"context"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	DB config.DBConn
}

// NewUserRepository creates a new user repository instance.
func NewUserRepository(conn *pgxpool.Pool) UserRepository {
	return &userRepo{DB: conn}
}

// FindByUsername finds a user by its username.
func (r *userRepo) FindByUsername(ctx context.Context, uName string) (*model.User, errors.CommonError) {
	query := `SELECT * FROM users WHERE username=$1`

	var user model.User

	findErr := r.DB.QueryRow(ctx, query, uName).Scan(&user.ID, &user.Username, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.Enabled)
	if findErr != nil {
		return nil, errors.NewInternalServerError("Database error", findErr)
	}

	return &user, nil
}

// Save persists a user to the DB.
func (r *userRepo) Save(ctx context.Context, user *model.User) (*model.User, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO users("username", "password", "email", "created_at", "updated_at", "enabled") 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, user.Username, user.Password, user.Email,
		user.CreatedAt, user.UpdatedAt, user.Enabled)

	saveErr := row.Scan(&user.ID)
	if saveErr != nil {
		return nil, errors.NewInternalServerError("Database error", saveErr)
	}

	return user, nil
}
