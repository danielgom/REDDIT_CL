package internal

import "time"

// RefreshTokenRequest in order to request a new JWT token with the refresh token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	Username     string `json:"username" validate:"required"`
}

type RefreshTokenResponse struct {
	Username     string    `json:"username"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
