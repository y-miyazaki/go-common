package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// RedisRepositoryInterface interface.
type RedisRepositoryInterface interface {
}

// RedisRepository struct.
type RedisRepository struct {
	logger *logger.Logger
	redis  *redis.Client
}

// NewRedisRepository returns RedisRepository instance.
func NewRedisRepository(logger *logger.Logger, r *redis.Client) *RedisRepository {
	return &RedisRepository{
		logger: logger,
		redis:  r,
	}
}

// Append Redis `APPEND key value` command.
func (r *RedisRepository) Append(c context.Context, key, value string) error {
	return r.redis.Append(c, key, value).Err()
}

// BitCount Redis `BITCOUNT key` command.
func (r *RedisRepository) BitCount(c context.Context, key string, bitCount *redis.BitCount) (int64, error) {
	return r.redis.BitCount(c, key, bitCount).Result()
}

// Decr Redis `DECR key` command.
func (r *RedisRepository) Decr(c context.Context, key string) error {
	return r.redis.Decr(c, key).Err()
}

// DecrBy Redis `DECRBY key value` command.
func (r *RedisRepository) DecrBy(c context.Context, key string, value int64) error {
	return r.redis.DecrBy(c, key, value).Err()
}

// Get Redis `GET key` command.
func (r *RedisRepository) Get(c context.Context, key string) (string, error) {
	return r.redis.Get(c, key).Result()
}

// GetBit Redis `GETBIT key start end` command.
func (r *RedisRepository) GetBit(c context.Context, key string, offset int64) (int64, error) {
	return r.redis.GetBit(c, key, offset).Result()
}

// GetRange Redis `GETRANGE key start end` command.
func (r *RedisRepository) GetRange(c context.Context, key string, start, end int64) (string, error) {
	return r.redis.GetRange(c, key, start, end).Result()
}

// GetSet Redis `GETSET key` command.
func (r *RedisRepository) GetSet(c context.Context, key string, value interface{}) (string, error) {
	return r.redis.GetSet(c, key, value).Result()
}

// Incr Redis `INCR key` command.
func (r *RedisRepository) Incr(c context.Context, key string) error {
	return r.redis.Incr(c, key).Err()
}

// IncrBy Redis `INCRBY key value` command.
func (r *RedisRepository) IncrBy(c context.Context, key string, value int64) error {
	return r.redis.IncrBy(c, key, value).Err()
}

// IncrByfloat Redis `INCRBYFLOAT key value` command.
func (r *RedisRepository) IncrByfloat(c context.Context, key string, value float64) error {
	return r.redis.IncrByFloat(c, key, value).Err()
}

// MGet Redis `MGET keys...` command.
func (r *RedisRepository) MGet(c context.Context, keys ...string) ([]interface{}, error) {
	return r.redis.MGet(c, keys...).Result()
}

// MSet Redis `MSET key value key2 value2...` command.
func (r *RedisRepository) MSet(c context.Context, values ...interface{}) (string, error) {
	return r.redis.MSet(c, values...).Result()
}

// MSetNX Redis `MSETNX key value key2 value2...` command.
func (r *RedisRepository) MSetNX(c context.Context, values ...interface{}) (bool, error) {
	return r.redis.MSetNX(c, values...).Result()
}

// Set Redis `SET key value [expiration]` command.
func (r *RedisRepository) Set(c context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redis.Set(c, key, value, expiration).Err()
}

// SetBit Redis `SETBIT key value offset value` command.
func (r *RedisRepository) SetBit(c context.Context, key string, offset int64, value int) error {
	return r.redis.SetBit(c, key, offset, value).Err()
}

// SetEX Redis `SETEX key value expiration` command.
func (r *RedisRepository) SetEX(c context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redis.SetEX(c, key, value, expiration).Err()
}

// SetNX Redis `SETNX key value [expiration]` command.
func (r *RedisRepository) SetNX(c context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redis.SetNX(c, key, value, expiration).Err()
}

// SetRange Redis `SETRANGE key start end` command.
func (r *RedisRepository) SetRange(c context.Context, key string, offset int64, value string) (int64, error) {
	return r.redis.SetRange(c, key, offset, value).Result()
}

// StrLen Redis `STRLEN key` command.
func (r *RedisRepository) StrLen(c context.Context, key string) (int64, error) {
	return r.redis.StrLen(c, key).Result()
}
