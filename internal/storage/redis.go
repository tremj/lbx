package storage

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func Get(ctx context.Context, name string) (string, error) {
	return redisClient.Get(ctx, name).Result()
}

func SaveConfig(ctx context.Context, name string, data []byte) error {
	return redisClient.Set(ctx, name, data, 0).Err()
}

func DeleteConfig(ctx context.Context, name string) error {
	return redisClient.Del(ctx, name).Err()
}

func GetKeys(ctx context.Context) ([]string, error) {
	return redisClient.Keys(ctx, "*").Result()
}
