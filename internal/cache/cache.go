package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	key string        = "products"
	TTL time.Duration = 5 * time.Minute
)

type Cache interface {
	Set(ctx context.Context, value any) error
	Get(ctx context.Context, key string) (any, error)
}

func NewRedisConfig(rdb *redis.Client) Cache {
	return &redisCache{rdb: rdb}
}
