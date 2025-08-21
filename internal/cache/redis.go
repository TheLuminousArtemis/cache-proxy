package cache

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	rdb *redis.Client
}

func (rc *redisCache) Set(ctx context.Context, key string, value any) error {
	err := rc.rdb.Set(ctx, key, value, TTL).Err()
	if err != nil {
		slog.Error("Failed to set value in cache", "error", err)
		return err
	}
	return nil
}

func (rc *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := rc.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (rc *redisCache) Del(ctx context.Context) error {
	err := rc.rdb.FlushDB(ctx).Err()
	if err != nil {
		return err
	}
	return nil
}
