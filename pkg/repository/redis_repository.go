// Package repository provides data access layer implementations for various services
// including Redis, AWS services, and other infrastructure components.
package repository

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// RedisClientInterface defines the interface for Redis client operations
type RedisClientInterface interface {
	// Append Redis `APPEND key value` command.
	Append(_ context.Context, _ string, _ string) *redis.IntCmd
	// BitCount Redis `BITCOUNT key` command.
	BitCount(_ context.Context, _ string, _ *redis.BitCount) *redis.IntCmd
	// Decr Redis `DECR key` command.
	Decr(_ context.Context, _ string) *redis.IntCmd
	// DecrBy Redis `DECRBY key value` command.
	DecrBy(_ context.Context, _ string, _ int64) *redis.IntCmd
	// Del Redis `DEL key [key ...]` command.
	Del(_ context.Context, _ ...string) *redis.IntCmd
	// Exists Redis `EXISTS key [key ...]` command.
	Exists(_ context.Context, _ ...string) *redis.IntCmd
	// Get Redis `GET key` command.
	Get(_ context.Context, _ string) *redis.StringCmd
	// GetBit Redis `GETBIT key start end` command.
	GetBit(_ context.Context, _ string, _ int64) *redis.IntCmd
	// GetRange Redis `GETRANGE key start end` command.
	GetRange(_ context.Context, _ string, _ int64, _ int64) *redis.StringCmd
	// GetSet Redis `GETSET key` command.
	GetSet(_ context.Context, _ string, _ any) *redis.StringCmd
	// Incr Redis `INCR key` command.
	Incr(_ context.Context, _ string) *redis.IntCmd
	// IncrBy Redis `INCRBY key value` command.
	IncrBy(_ context.Context, _ string, _ int64) *redis.IntCmd
	// IncrByFloat Redis `INCRBYFLOAT key value` command.
	IncrByFloat(_ context.Context, _ string, _ float64) *redis.FloatCmd
	// MGet Redis `MGET keys...` command.
	MGet(_ context.Context, _ ...string) *redis.SliceCmd
	// MSet Redis `MSET key value key2 value2...` command.
	MSet(_ context.Context, _ ...any) *redis.StatusCmd
	// MSetNX Redis `MSETNX key value key2 value2...` command.
	MSetNX(_ context.Context, _ ...any) *redis.BoolCmd
	// Ping Redis `PING` command.
	Ping(_ context.Context) *redis.StatusCmd
	// Set Redis `SET key value [expiration]` command.
	Set(_ context.Context, _ string, _ any, _ time.Duration) *redis.StatusCmd
	// SetBit Redis `SETBIT key value offset value` command.
	SetBit(_ context.Context, _ string, _ int64, _ int) *redis.IntCmd
	// SetEX Redis `SETEX key value expiration` command.
	SetEX(_ context.Context, _ string, _ any, _ time.Duration) *redis.StatusCmd
	// SetNX Redis `SETNX key value [expiration]` command.
	SetNX(_ context.Context, _ string, _ any, _ time.Duration) *redis.BoolCmd
	// SetRange Redis `SETRANGE key start end` command.
	SetRange(_ context.Context, _ string, _ int64, _ string) *redis.IntCmd
	// TTL Redis `TTL key` command.
	TTL(_ context.Context, _ string) *redis.DurationCmd
	// Expire Redis `EXPIRE key seconds` command.
	Expire(_ context.Context, _ string, _ time.Duration) *redis.BoolCmd
	// StrLen Redis `STRLEN key` command.
	StrLen(_ context.Context, _ string) *redis.IntCmd
}

// RedisRepository struct.
type RedisRepository struct {
	Client RedisClientInterface
}

// NewRedisRepository returns RedisRepository instance.
func NewRedisRepository(r *redis.Client) *RedisRepository {
	return &RedisRepository{
		Client: r,
	}
}

// NewRedisRepositoryWithInterface creates RedisRepository with interface for testing
func NewRedisRepositoryWithInterface(r RedisClientInterface) *RedisRepository {
	return &RedisRepository{
		Client: r,
	}
}

// Append Redis `APPEND key value` command.
func (r *RedisRepository) Append(c context.Context, key, value string) (int64, error) {
	res, err := r.Client.Append(c, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("redis Append: %w", err)
	}
	return res, nil
}

// BitCount Redis `BITCOUNT key` command.
func (r *RedisRepository) BitCount(c context.Context, key string, bitCount *redis.BitCount) (int64, error) {
	res, err := r.Client.BitCount(c, key, bitCount).Result()
	if err != nil {
		return 0, fmt.Errorf("redis BitCount: %w", err)
	}
	return res, nil
}

// Decr Redis `DECR key` command.
func (r *RedisRepository) Decr(c context.Context, key string) error {
	if err := r.Client.Decr(c, key).Err(); err != nil {
		return fmt.Errorf("redis Decr: %w", err)
	}
	return nil
}

// DecrBy Redis `DECRBY key value` command.
func (r *RedisRepository) DecrBy(c context.Context, key string, value int64) error {
	if err := r.Client.DecrBy(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis DecrBy: %w", err)
	}
	return nil
}

// Del Redis `DEL key [key ...]` command.
func (r *RedisRepository) Del(c context.Context, keys ...string) (int64, error) {
	res, err := r.Client.Del(c, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis Del: %w", err)
	}
	return res, nil
}

// Exists Redis `EXISTS key [key ...]` command.
func (r *RedisRepository) Exists(c context.Context, keys ...string) (int64, error) {
	res, err := r.Client.Exists(c, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("redis Exists: %w", err)
	}
	return res, nil
}

// Get Redis `GET key` command.
func (r *RedisRepository) Get(c context.Context, key string) (string, error) {
	res, err := r.Client.Get(c, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis Get: %w", err)
	}
	return res, nil
}

// GetBit Redis `GETBIT key start end` command.
func (r *RedisRepository) GetBit(c context.Context, key string, offset int64) (int64, error) {
	res, err := r.Client.GetBit(c, key, offset).Result()
	if err != nil {
		return 0, fmt.Errorf("redis GetBit: %w", err)
	}
	return res, nil
}

// GetRange Redis `GETRANGE key start end` command.
func (r *RedisRepository) GetRange(c context.Context, key string, start, end int64) (string, error) {
	res, err := r.Client.GetRange(c, key, start, end).Result()
	if err != nil {
		return "", fmt.Errorf("redis GetRange: %w", err)
	}
	return res, nil
}

// GetSet Redis `GETSET key` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) GetSet(c context.Context, key string, value any) (string, error) {
	res, err := r.Client.GetSet(c, key, value).Result()
	if err != nil {
		return "", fmt.Errorf("redis GetSet: %w", err)
	}
	return res, nil
}

// Incr Redis `INCR key` command.
func (r *RedisRepository) Incr(c context.Context, key string) (int64, error) {
	res, err := r.Client.Incr(c, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis Incr: %w", err)
	}
	return res, nil
}

// IncrBy Redis `INCRBY key value` command.
func (r *RedisRepository) IncrBy(c context.Context, key string, value int64) error {
	if err := r.Client.IncrBy(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis IncrBy: %w", err)
	}
	return nil
}

// IncrByfloat Redis `INCRBYFLOAT key value` command.
func (r *RedisRepository) IncrByfloat(c context.Context, key string, value float64) error {
	if err := r.Client.IncrByFloat(c, key, value).Err(); err != nil {
		return fmt.Errorf("redis IncrByFloat: %w", err)
	}
	return nil
}

// MGet Redis `MGET keys...` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) MGet(c context.Context, keys ...string) ([]any, error) {
	res, err := r.Client.MGet(c, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis MGet: %w", err)
	}
	return res, nil
}

// MSet Redis `MSET key value key2 value2...` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) MSet(c context.Context, values ...any) (string, error) {
	res, err := r.Client.MSet(c, values...).Result()
	if err != nil {
		return "", fmt.Errorf("redis MSet: %w", err)
	}
	return res, nil
}

// MSetNX Redis `MSETNX key value key2 value2...` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) MSetNX(c context.Context, values ...any) (bool, error) {
	res, err := r.Client.MSetNX(c, values...).Result()
	if err != nil {
		return false, fmt.Errorf("redis MSetNX: %w", err)
	}
	return res, nil
}

// Ping Redis `PING` command.
func (r *RedisRepository) Ping(c context.Context) error {
	if err := r.Client.Ping(c).Err(); err != nil {
		return fmt.Errorf("redis Ping: %w", err)
	}
	return nil
}

// Set Redis `SET key value [expiration]` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) Set(c context.Context, key string, value any, expiration time.Duration) error {
	if err := r.Client.Set(c, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis Set: %w", err)
	}
	return nil
}

// SetBit Redis `SETBIT key value offset value` command.
func (r *RedisRepository) SetBit(c context.Context, key string, offset int64, value int) error {
	if err := r.Client.SetBit(c, key, offset, value).Err(); err != nil {
		return fmt.Errorf("redis SetBit: %w", err)
	}
	return nil
}

// SetEX Redis `SETEX key value expiration` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) SetEX(c context.Context, key string, value any, expiration time.Duration) error {
	if err := r.Client.SetEX(c, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis SetEX: %w", err)
	}
	return nil
}

// SetNX Redis `SETNX key value [expiration]` command.
// nolint:revive // keep interface{} for Go 1.16 compatibility
func (r *RedisRepository) SetNX(c context.Context, key string, value any, expiration time.Duration) error {
	if err := r.Client.SetNX(c, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis SetNX: %w", err)
	}
	return nil
}

// SetRange Redis `SETRANGE key start end` command.
func (r *RedisRepository) SetRange(c context.Context, key string, offset int64, value string) (int64, error) {
	res, err := r.Client.SetRange(c, key, offset, value).Result()
	if err != nil {
		return 0, fmt.Errorf("redis SetRange: %w", err)
	}
	return res, nil
}

// TTL Redis `TTL key` command.
func (r *RedisRepository) TTL(c context.Context, key string) (time.Duration, error) {
	res, err := r.Client.TTL(c, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis TTL: %w", err)
	}
	return res, nil
}

// Expire Redis `EXPIRE key seconds` command.
func (r *RedisRepository) Expire(c context.Context, key string, expiration time.Duration) (bool, error) {
	res, err := r.Client.Expire(c, key, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("redis Expire: %w", err)
	}
	return res, nil
}

// StrLen Redis `STRLEN key` command.
func (r *RedisRepository) StrLen(c context.Context, key string) (int64, error) {
	res, err := r.Client.StrLen(c, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis StrLen: %w", err)
	}
	return res, nil
}
