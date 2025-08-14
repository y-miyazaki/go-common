package repository

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// RedisRepositoryInterface interface.
// nolint:iface,revive,unused
type RedisRepositoryInterface interface {
	Append(c context.Context, key, value string) error
	BitCount(c context.Context, key string, bitCount *redis.BitCount) (int64, error)
	Decr(c context.Context, key string) error
	DecrBy(c context.Context, key string, value int64) error
	Get(c context.Context, key string) (string, error)
	GetBit(c context.Context, key string, offset int64) (int64, error)
	GetRange(c context.Context, key string, start, end int64) (string, error)
	GetSet(c context.Context, key string, value interface{}) (string, error)
	Incr(c context.Context, key string) error
	IncrBy(c context.Context, key string, value int64) error
	IncrByfloat(c context.Context, key string, value float64) error
	MGet(c context.Context, keys ...string) ([]interface{}, error)
	MSet(c context.Context, values ...interface{}) (string, error)
	MSetNX(c context.Context, values ...interface{}) (bool, error)
	Set(c context.Context, key string, value interface{}, expiration time.Duration) error
	SetBit(c context.Context, key string, offset int64, value int) error
	SetEX(c context.Context, key string, value interface{}, expiration time.Duration) error
	SetNX(c context.Context, key string, value interface{}, expiration time.Duration) error
	SetRange(c context.Context, key string, offset int64, value string) (int64, error)
	StrLen(c context.Context, key string) (int64, error)
}

// RedisRepository struct.
type RedisRepository struct {
	redis *redis.Client
}

// NewRedisRepository returns RedisRepository instance.
func NewRedisRepository(r *redis.Client) *RedisRepository {
	return &RedisRepository{
		redis: r,
	}
}

// Append Redis `APPEND key value` command.
func (r *RedisRepository) Append(c context.Context, key, value string) error {
	if err := r.redis.Append(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis Append: %w", err)
	}
	return nil
}

// BitCount Redis `BITCOUNT key` command.
func (r *RedisRepository) BitCount(c context.Context, key string, bitCount *redis.BitCount) (int64, error) {
	res, err := r.redis.BitCount(c, key, bitCount).Result()
	if err != nil {
		return 0, fmt.Errorf("redis BitCount: %w", err)
	}
	return res, nil
}

// Decr Redis `DECR key` command.
func (r *RedisRepository) Decr(c context.Context, key string) error {
	if err := r.redis.Decr(c, key).Err(); err != nil {
		return fmt.Errorf("redis Decr: %w", err)
	}
	return nil
}

// DecrBy Redis `DECRBY key value` command.
func (r *RedisRepository) DecrBy(c context.Context, key string, value int64) error {
	if err := r.redis.DecrBy(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis DecrBy: %w", err)
	}
	return nil
}

// Get Redis `GET key` command.
func (r *RedisRepository) Get(c context.Context, key string) (string, error) {
	res, err := r.redis.Get(c, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis Get: %w", err)
	}
	return res, nil
}

// GetBit Redis `GETBIT key start end` command.
func (r *RedisRepository) GetBit(c context.Context, key string, offset int64) (int64, error) {
	res, err := r.redis.GetBit(c, key, offset).Result()
	if err != nil {
		return 0, fmt.Errorf("redis GetBit: %w", err)
	}
	return res, nil
}

// GetRange Redis `GETRANGE key start end` command.
func (r *RedisRepository) GetRange(c context.Context, key string, start, end int64) (string, error) {
	res, err := r.redis.GetRange(c, key, start, end).Result()
	if err != nil {
		return "", fmt.Errorf("redis GetRange: %w", err)
	}
	return res, nil
}

// GetSet Redis `GETSET key` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) GetSet(c context.Context, key string, value interface{}) (string, error) {
	res, err := r.redis.GetSet(c, key, value).Result()
	if err != nil {
		return "", fmt.Errorf("redis GetSet: %w", err)
	}
	return res, nil
}

// Incr Redis `INCR key` command.
func (r *RedisRepository) Incr(c context.Context, key string) error {
	if err := r.redis.Incr(c, key).Err(); err != nil {
		return fmt.Errorf("redis Incr: %w", err)
	}
	return nil
}

// IncrBy Redis `INCRBY key value` command.
func (r *RedisRepository) IncrBy(c context.Context, key string, value int64) error {
	if err := r.redis.IncrBy(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis IncrBy: %w", err)
	}
	return nil
}

// IncrByfloat Redis `INCRBYFLOAT key value` command.
func (r *RedisRepository) IncrByfloat(c context.Context, key string, value float64) error {
	if err := r.redis.IncrByFloat(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis IncrByFloat: %w", err)
	}
	return nil
}

// MGet Redis `MGET keys...` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) MGet(c context.Context, keys ...string) ([]interface{}, error) {
	res, err := r.redis.MGet(c, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis MGet: %w", err)
	}
	return res, nil
}

// MSet Redis `MSET key value key2 value2...` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) MSet(c context.Context, values ...interface{}) (string, error) {
	res, err := r.redis.MSet(c, values...).Result()
	if err != nil {
		return "", fmt.Errorf("redis MSet: %w", err)
	}
	return res, nil
}

// MSetNX Redis `MSETNX key value key2 value2...` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) MSetNX(c context.Context, values ...interface{}) (bool, error) {
	res, err := r.redis.MSetNX(c, values...).Result()
	if err != nil {
		return false, fmt.Errorf("redis MSetNX: %w", err)
	}
	return res, nil
}

// Set Redis `SET key value [expiration]` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) Set(c context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.redis.Set(c, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis Set: %w", err)
	}
	return nil
}

// SetBit Redis `SETBIT key value offset value` command.
func (r *RedisRepository) SetBit(c context.Context, key string, offset int64, value int) error {
	if err := r.redis.SetBit(c, key, offset, value).Err(); err != nil {
		return fmt.Errorf("redis SetBit: %w", err)
	}
	return nil
}

// SetEX Redis `SETEX key value expiration` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) SetEX(c context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.redis.SetEX(c, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis SetEX: %w", err)
	}
	return nil
}

// SetNX Redis `SETNX key value [expiration]` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) SetNX(c context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.redis.SetNX(c, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis SetNX: %w", err)
	}
	return nil
}

// SetRange Redis `SETRANGE key start end` command.
func (r *RedisRepository) SetRange(c context.Context, key string, offset int64, value string) (int64, error) {
	res, err := r.redis.SetRange(c, key, offset, value).Result()
	if err != nil {
		return 0, fmt.Errorf("redis SetRange: %w", err)
	}
	return res, nil
}

// StrLen Redis `STRLEN key` command.
func (r *RedisRepository) StrLen(c context.Context, key string) (int64, error) {
	res, err := r.redis.StrLen(c, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis StrLen: %w", err)
	}
	return res, nil
}
