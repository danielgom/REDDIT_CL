package security

import (
	"fmt"
	"time"

	"RD-Clone-API/pkg/config"
	"github.com/golang-jwt/jwt"
)

// GenerateTokenWithExp generates a JWT with an expiration of 1 hour (exp time comes from the config).
func GenerateTokenWithExp(username string) (string, time.Time, error) {
	jwtConfig := config.LoadConfig().JWT

	currentTime := time.Now().Local()
	expirationDate := currentTime.Add(time.Second * time.Duration(jwtConfig.Expiration))

	claims := jwt.StandardClaims{
		ExpiresAt: expirationDate.Unix(),
		IssuedAt:  currentTime.Unix(),
		Issuer:    "GO-Reddit-CL",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtConfig.Key))

	if err != nil {
		return "", time.Time{}, fmt.Errorf("could not generate JWT %w please try again", err)
	}

	return signedToken, expirationDate, nil
}
