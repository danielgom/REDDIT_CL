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
	"RD-Clone-API/pkg/config/core"
	"RD-Clone-API/pkg/config/logger"
)

const defaultServerTimeout = time.Second * 5

func main() {
	logger.Initialise()
	r := core.Router()

	r.Server = &http.Server{
		ReadTimeout:       defaultServerTimeout,
		WriteTimeout:      defaultServerTimeout,
		IdleTimeout:       defaultServerTimeout,
		ReadHeaderTimeout: defaultServerTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Start(config.LoadConfig().Port); err != nil {
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
