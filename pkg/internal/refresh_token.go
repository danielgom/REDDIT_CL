package internal

// RefreshTokenRequest in order to request a new JWT token with the refresh token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	Username     string `json:"username" validate:"required"`
}
