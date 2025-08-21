package cache

import "github.com/redis/go-redis/v9"

type redisCache struct {
	rdb *redis.Client
}

func (rc *redisCache) Set(key string, value any) error {
	return nil
}

func (rc *redisCache) Get(key string) (any, error) {
	return nil, nil
}
