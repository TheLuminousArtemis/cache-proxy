package main

import (
	"fmt"
	"log/slog"
	"net/http"

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
	r.Get("/products", app.getProducts)
	return r
}

func (app *application) start(r http.Handler) error {
	app.logger.Info("Starting proxy server", "port", app.config.Port, "origin", app.config.Origin)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), r)
}
