package bootstrap

import (
	"os"

	"github.com/go-redis/redis"
)

// NewRedisClient creates a new Redis client
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})
}
