package cache

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/theluminousartemis/caching-proxy/internal/models"
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

func (rc *redisCache) Get(ctx context.Context, key string) (*models.Products, error) {
	cacheResp := rc.rdb.Get(ctx, key)
	if cacheResp.Err() != nil {
		slog.Error("Failed to get value from cache", "error", cacheResp.Err())
		return nil, cacheResp.Err()
	}
	var products models.Products
	if err := json.Unmarshal([]byte(cacheResp.Val()), &products); err != nil {
		slog.Error("Failed to unmarshal value", "error", err)
		return nil, err
	}
	return &products, nil
}

func (rc *redisCache) Del(ctx context.Context, key string) error {
	return rc.rdb.Del(ctx, key).Err()
}
