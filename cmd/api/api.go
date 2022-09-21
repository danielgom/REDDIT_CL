// Package main where the execution of the program lives.
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/routes"
	"RD-Clone-API/pkg/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const defaultServerTimeout = time.Second * 5

func main() {
	c := config.LoadConfig()
	r := echo.New()

	DBc := config.InitDatabase(c)
	userRepository := db.NewUserRepository(DBc)
	tokenRepository := db.NewTokenRepository(DBc)

	refreshTokenService := services.NewRefreshTokenService()
	userService := services.NewUserService(userRepository, tokenRepository, refreshTokenService)
	userHandler := routes.NewUserHandler(userService)
	userHandler.Register(r)

	r.Use(middleware.Logger())
	r.Server = &http.Server{
		ReadTimeout:       defaultServerTimeout,
		WriteTimeout:      defaultServerTimeout,
		IdleTimeout:       defaultServerTimeout,
		ReadHeaderTimeout: defaultServerTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Start(c.Port); err != nil {
			log.Fatalln("error serving", err.Error())
		}
	}()

	<-ctx.Done()

	stop()

	log.Println("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), defaultServerTimeout)
	if err := r.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalln("server forced to shutdown: ", err)
	}
}
