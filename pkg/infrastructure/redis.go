package infrastructure

import redis "github.com/redis/go-redis/v9"

// NewRedis returns redis client.
func NewRedis(o *redis.Options) *redis.Client {
	return redis.NewClient(o)
}
