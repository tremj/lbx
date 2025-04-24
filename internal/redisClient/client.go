package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	rc *redis.Client
}

func NewClient(addr string) *RedisClient {
	rc := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	return &RedisClient{rc: rc}
}

func (c *RedisClient) Set(key, value string) error {
	return c.rc.Set(ctx, key, value, 0).Err()
}

func (c *RedisClient) Get(key string) (string, error) {
	return c.rc.Get(ctx, key).Result()
}

func (c *RedisClient) Del(key string) error {
	return c.rc.Del(ctx, key).Err()
}
