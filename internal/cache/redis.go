package cache

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	rdb *redis.Client
}

func (rc *redisCache) Set(ctx context.Context, value any) error {
	json, err := json.Marshal(value)
	if err != nil {
		slog.Error("Failed to marshal value", "error", err)
		return err
	}
	err = rc.rdb.Set(ctx, key, json, TTL).Err()
	if err != nil {
		slog.Error("Failed to set value in cache", "error", err)
		return err
	}
	return nil
}

func (rc *redisCache) Get(ctx context.Context, key string) (any, error) {
	return rc.rdb.Get(ctx, key).Result()
}
