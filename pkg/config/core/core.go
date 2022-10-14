// Package core in charge of initialising core configuration, DBs, repositories, services and handler.
package core

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/config/logger"
	"RD-Clone-API/pkg/context"
	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/routes"
	"RD-Clone-API/pkg/services"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var (
	errInvalidToken    = errors.New("invalid token")
	errParse           = errors.New("unable to parse token")
	errMalformedClaims = fmt.Errorf("malformed claims")
	errGetClaim        = fmt.Errorf("unable to get claim")
	errInvalidSigning  = fmt.Errorf("unexpected jwt signing method")
)

// Router initialises api and returns router to serve.
func Router() *echo.Echo {
	router := initialiseAPI()

	router.Use(doJWTFilter())
	router.Use(doLoggerFilter())

	return router
}

func doJWTFilter() echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		return strings.HasPrefix(c.Request().URL.RequestURI(), "/api/auth")
	}

	return middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:                skipper,
		ContinueOnIgnoredError: false,
		ContextKey:             "sub",
		TokenLookup:            "header:" + echo.HeaderAuthorization,
		AuthScheme:             "Bearer",
		ParseTokenFunc:         getParseTokenFunc(),
	})
}

func doLoggerFilter() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:  true,
		LogURI:     true,
		LogURIPath: true,
		LogStatus:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request", zap.String("URI", v.URI),
				zap.Any("latency", v.Latency),
				zap.String("method", v.Method),
				zap.Int("status", v.Status))
			return nil
		},
	})
}

func getParseTokenFunc() func(auth string, c echo.Context) (interface{}, error) {
	c := config.LoadConfig()
	signingKey := []byte(c.JWT.Key)

	return func(auth string, c echo.Context) (interface{}, error) {
		keyFunc := func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != "HS512" {
				return nil, errInvalidSigning
			}
			return signingKey, nil
		}

		token, err := jwt.Parse(auth, keyFunc)
		if err != nil {
			return nil, errParse
		}

		if !token.Valid {
			return nil, errInvalidToken
		}

		// check claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			err = errMalformedClaims
		}

		// func to get a claim by name
		getClaim := func(claim string) (string, error) {
			var str string
			if val, ok := claims[claim]; ok {
				if str, ok = val.(string); !ok {
					return str, errMalformedClaims
				}
			}
			if str == "" {
				err = errMalformedClaims
			}
			return str, err
		}

		username, err := getClaim("sub")
		if err != nil {
			return nil, errGetClaim
		}

		c.Set("user", username)

		return token, nil
	}
}

func initialiseAPI() *echo.Echo {
	c := config.LoadConfig()
	DBc := config.InitDatabase(c)
	router := echo.New()

	v := config.GetValidator()
	err := config.AddValidators(v.Validator)
	if err != nil {
		log.Fatalln("could not add validators")
	}

	router.Validator = v

	userRepository := db.NewUserRepository(DBc)
	tokenRepository := db.NewTokenRepository(DBc)
	refreshTokenRepository := db.NewRTRepository(DBc)

	refreshTokenService := services.NewRefreshTokenService(refreshTokenRepository)
	userService := services.NewUserService(userRepository, tokenRepository, refreshTokenService)

	userHandler := routes.NewUserHandler(userService)
	userHandler.Register(router, context.Handler)

	return router
}
