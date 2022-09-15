package main

import (
	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/routes"
	"RD-Clone-API/pkg/services"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	c := config.LoadConfig()
	r := echo.New()

	DBconn := config.InitDatabase(c)
	userRepository := db.NewUserRepository(DBconn)
	tokenRepository := db.NewTokenRepository(DBconn)

	userService := services.NewUserService(userRepository, tokenRepository)
	userHandler := routes.NewUserHandler(userService)
	userHandler.Register(r)

	srv := http.Server{
		Addr:              c.Port,
		Handler:           r,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 3,
		WriteTimeout:      time.Second * 5,
		IdleTimeout:       time.Second * 5,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("error serving", err.Error())
		}
	}()

	<-ctx.Done()

	stop()

	log.Println("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	log.Println("server exiting")

}
