package storage

import (
	"github.com/go-redis/redis/v8"
	"context"
	"time"
)

type CacheProvider interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (c *RedisCache) Get(key string) (string, bool) {
	val, err := c.client.Get(context.Background(), key).Result()
	return val, err == nil
}

func (c *RedisCache) Set(key, value string) {
	c.client.Set(context.Background(), key, value, 30*time.Minute)
}