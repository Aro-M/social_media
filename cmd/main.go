package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"social_media/internal/config"
	"social_media/internal/repository/postgres"
	"social_media/internal/router/handler"
	"social_media/internal/service"
	"social_media/pkg/logger"
)

func main() {
	if err := config.Init(config.DotEnv); err != nil {
		panic(err)
	}

	log := logger.New(config.LogLevel())

	log.WithFields(logrus.Fields{
		"version": config.Version(),
		"port":    config.Port(),
	}).Info("starting social_media")

	repo, err := postgres.NewRepository()
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	defer repo.Close()

	svc := service.New(repo, log)

	e := echo.New()
	e.Use(middleware.Recover())
	h := handler.New(svc)
	h.Register(e)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", config.Port())); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout())
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("server forced to shutdown")
	}

	log.Info("server stopped")
}
