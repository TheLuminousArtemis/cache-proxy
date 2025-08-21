package cache

import "github.com/redis/go-redis/v9"

type Cache interface {
	Set(key string, value any) error
	Get(key string) (any, error)
}

func NewRedisConfig(rdb *redis.Client) Cache {
	return &redisCache{rdb: rdb}
}
