package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/theluminousartemis/caching-proxy/internal/cache"
	"github.com/theluminousartemis/caching-proxy/internal/env"
)

var (
	port   string = "port"
	origin string = "origin"
)

func main() {
	//flags
	portFlag := flag.Int(port, 9090, "port where proxy server would start on")
	originFlag := flag.String(origin, "http://dummyjson.com/products", "origin server where requests are forwarded to")
	flag.Parse()

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	cfg := config{
		Port:   *portFlag,
		Origin: *originFlag,
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
		},
	}

	//redis
	logger.Info("Connecting to redis", "address", cfg.redisCfg.addr, "password", cfg.redisCfg.password, "db", cfg.redisCfg.db)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.redisCfg.addr,
		Password: cfg.redisCfg.password,
		DB:       cfg.redisCfg.db,
	})
	status, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("Failed to connect to redis", "error", err)
		log.Fatal(err)
	}
	logger.Info("Connected to redis", "status", status)
	cache := cache.NewRedisConfig(rdb)

	//application
	app := &application{
		config: cfg,
		logger: logger,
		cache:  cache,
	}

	if err := app.start(app.mount()); err != nil {
		app.logger.Error("Failed to start proxy server", "error", err)
		log.Fatal(err)
	}
}
