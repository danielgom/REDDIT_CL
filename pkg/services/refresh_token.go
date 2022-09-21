package services

import (
	"context"
	"time"

	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/google/uuid"
)

type refreshTokenSvc struct {
	rTDB db.RefreshTokenRepository
}

// NewRefreshTokenService returns a new refresh token service instance.
func NewRefreshTokenService(rTDB db.RefreshTokenRepository) RefreshTokenService {
	return &refreshTokenSvc{rTDB: rTDB}
}

// Create creates a new refresh token.
func (r *refreshTokenSvc) Create(ctx context.Context) (string, errors.CommonError) {
	const refreshTokenValidHours = 24

	token := uuid.New().String()

	refreshT := &model.RefreshToken{
		Token:     token,
		ExpiresAt: time.Now().Local().Add(time.Hour * refreshTokenValidHours),
	}

	saveErr := r.rTDB.Save(ctx, refreshT)
	if saveErr != nil {
		return "", saveErr
	}

	return token, nil
}

// Validate checks whether the current refresh token is valid and if it has not yet expired.
func (r *refreshTokenSvc) Validate(ctx context.Context, token string) errors.CommonError {
	refreshToken, commonError := r.rTDB.FindByToken(ctx, token)
	if commonError != nil {
		return commonError
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return errors.NewUnauthorisedError("refresh token expired")
	}

	return nil
}
