// Package core in charge of initialising core configuration, DBs, repositories, services and handler.
package core

import (
	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/config/logger"
	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/routes"
	"RD-Clone-API/pkg/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

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

	userRepository := db.NewUserRepository(DBc)
	tokenRepository := db.NewTokenRepository(DBc)
	refreshTokenRepository := db.NewRTRepository(DBc)

	refreshTokenService := services.NewRefreshTokenService(refreshTokenRepository)
	userService := services.NewUserService(userRepository, tokenRepository, refreshTokenService)

	userHandler := routes.NewUserHandler(userService)
	userHandler.Register(router)

	return router
}
