// Package core in charge of initialising core configuration, DBs, repositories, services and handler.
//
//nolint:wrapcheck // Should not wrap validation errors
package core

import (
	"strings"
	"unicode"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/config/logger"
	"RD-Clone-API/pkg/context"
	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/routes"
	"RD-Clone-API/pkg/services"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type customValidator struct {
	validator *validator.Validate
}

// Router initialises api and returns router to serve.
func Router() *echo.Echo {
	router := initialiseAPI()

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	}))

	return router
}

func initialiseAPI() *echo.Echo {
	c := config.LoadConfig()
	DBc := config.InitDatabase(c)
	router := echo.New()
	router.Validator = &customValidator{validator: validator.New()}

	userRepository := db.NewUserRepository(DBc)
	tokenRepository := db.NewTokenRepository(DBc)
	refreshTokenRepository := db.NewRTRepository(DBc)

	refreshTokenService := services.NewRefreshTokenService(refreshTokenRepository)
	userService := services.NewUserService(userRepository, tokenRepository, refreshTokenService)

	userHandler := routes.NewUserHandler(userService)
	userHandler.Register(router, context.Handler)

	return router
}

func (c *customValidator) Validate(i any) error {
	if err := c.validator.RegisterValidation("password", validatePassword); err != nil {
		return err
	}
	if err := c.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	hasNumber = strings.ContainsAny(password, "123456789")
	hasUpper = strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTVWXYZ")
	hasLower = strings.ContainsAny(password, "abcdefghijklmnopqrstvwxyz")

	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
