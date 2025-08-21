package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/theluminousartemis/caching-proxy/internal/models"
)

var (
	key string        = "products"
	TTL time.Duration = 5 * time.Minute
)

type Cache interface {
	Set(ctx context.Context, value any) error
	Get(ctx context.Context, key string) (*models.Products, error)
	Del(ctx context.Context, key string) error
}

func NewRedisConfig(rdb *redis.Client) Cache {
	return &redisCache{rdb: rdb}
}
