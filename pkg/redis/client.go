package redis

import (
	"github.com/go-redis/redis"
)

// NewClient creates a new Redis client
func NewClient(address string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: address,
	})
}
