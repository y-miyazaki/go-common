package infrastructure

import redis "github.com/go-redis/redis/v8"

// NewRedis returns redis client.
func NewRedis(o *redis.Options) *redis.Client {
	return redis.NewClient(o)
}
