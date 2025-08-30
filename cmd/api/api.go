package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/theluminousartemis/caching-proxy/internal/cache"
)

type config struct {
	Port       int
	Origin     string
	redisCfg   redisConfig
	clearCache bool
}

type redisConfig struct {
	addr     string
	password string
	db       int
}

type application struct {
	config config
	logger *slog.Logger
	cache  cache.Cache
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/health", app.health)
	r.Get("/", app.getContent)
	return r
}

func (app *application) start(r http.Handler) error {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.Port),
		Handler: r,
	}
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.logger.Info("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Info("Starting proxy server", "port", app.config.Port, "origin", app.config.Origin)
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}
	app.logger.Info("server has stopped", "addr", app.config.Port)
	return nil

}
