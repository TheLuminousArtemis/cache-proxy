package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	TTL time.Duration = 5 * time.Minute
)

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context) error
}

func NewRedisConfig(rdb *redis.Client) Cache {
	return &redisCache{rdb: rdb}
}
