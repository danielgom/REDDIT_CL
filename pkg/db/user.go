// Package db contains all repositories used by this API.
package db

import (
	"context"
	"strings"

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
	return r.findUser(ctx, `SELECT * FROM users WHERE username=$1`, uName)
}

// FindByEmail finds a user by its email.
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, errors.CommonError) {
	return r.findUser(ctx, `SELECT * FROM users WHERE email=$1`, email)
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

func (r *userRepo) Update(ctx context.Context, user *model.User) errors.CommonError {
	_, err := r.DB.Exec(ctx, `UPDATE users SET password=$2, email=$3, updated_at=$4, enabled=$5 WHERE username=$1`,
		user.Username, user.Password, user.Email, user.UpdatedAt, user.Enabled)

	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}

	return nil
}

func (r *userRepo) findUser(ctx context.Context, query string, args ...any) (*model.User, errors.CommonError) {
	var user model.User

	findErr := r.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Email,
		&user.CreatedAt, &user.UpdatedAt, &user.Enabled)
	if findErr != nil {
		if strings.Contains(findErr.Error(), errNotFound) {
			return nil, errors.NewNotFoundError("user not found")
		}
		return nil, errors.NewInternalServerError("Database error", findErr)
	}

	return &user, nil
}
