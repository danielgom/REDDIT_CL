package security

import (
	"fmt"
	"time"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/model"
	"github.com/golang-jwt/jwt"
)

// GenerateTokenWithExp generates a JWT with an expiration of 1 hour (exp time comes from the config).
func GenerateTokenWithExp(user *model.User) (string, error) {
	jwtConfig := config.LoadConfig().JWT

	currentTime := time.Now().Local()

	claims := jwt.StandardClaims{
		ExpiresAt: currentTime.Add(time.Second * time.Duration(jwtConfig.Expiration)).Unix(),
		Id:        "",
		IssuedAt:  currentTime.Unix(),
		Issuer:    "GO-Reddit-CL",
		Subject:   user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtConfig.Key))

	if err != nil {
		return "", fmt.Errorf("could not generate JWT %w please try again", err)
	}

	return signedToken, nil
}
