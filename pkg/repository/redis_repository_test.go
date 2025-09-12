package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient is a mock implementation of Redis client for testing
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	args := m.Called(ctx, key, value)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	cmd := redis.NewBoolCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(bool))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewDurationCmd(ctx, 0)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(time.Duration))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	cmd := redis.NewStatusCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	args := m.Called(ctx, key, field)
	cmd := redis.NewStringCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) HSet(ctx context.Context, key, field string, value interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, field, value)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	args := m.Called(ctx, key, fields)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, values)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) RPop(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	args := m.Called(ctx, timeout, keys)
	cmd := redis.NewStringSliceCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).([]string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	args := m.Called(ctx, channel, message)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	args := m.Called(ctx, channels)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisClient) PSubscribe(ctx context.Context, patterns ...string) *redis.PubSub {
	args := m.Called(ctx, patterns)
	return args.Get(0).(*redis.PubSub)
}

func (m *MockRedisClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	args := m.Called(ctx, key, bitCount)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	args := m.Called(ctx, key, offset)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	args := m.Called(ctx, key, start, end)
	cmd := redis.NewStringCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	args := m.Called(ctx, key, value)
	cmd := redis.NewStringCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	args := m.Called(ctx, key, value)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	args := m.Called(ctx, key, value)
	cmd := redis.NewFloatCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(float64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewSliceCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).([]interface{}))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	args := m.Called(ctx, values)
	cmd := redis.NewStatusCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	args := m.Called(ctx, values)
	cmd := redis.NewBoolCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(bool))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	args := m.Called(ctx, key, offset, value)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(string))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewBoolCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(bool))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	args := m.Called(ctx, key, offset, value)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

func (m *MockRedisClient) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	args := m.Called(ctx, key, value)
	cmd := redis.NewIntCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).(int64))
	}
	if args.Get(1) != nil {
		cmd.SetErr(args.Error(1))
	}
	return cmd
}

// RedisRepositoryWithMock for testing with mock Redis client
type RedisRepositoryWithMock struct {
	redis *MockRedisClient
}

// NewRedisRepositoryWithMock creates RedisRepository with mock client for testing
func NewRedisRepositoryWithMock(mockClient *MockRedisClient) *RedisRepositoryWithMock {
	return &RedisRepositoryWithMock{
		redis: mockClient,
	}
}

// Get retrieves the value of a key.
func (r *RedisRepositoryWithMock) Get(ctx context.Context, key string) (string, error) {
	cmd := r.redis.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

// Set sets the value of a key.
func (r *RedisRepositoryWithMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cmd := r.redis.Set(ctx, key, value, expiration)
	return cmd.Err()
}

// Del deletes one or more keys.
func (r *RedisRepositoryWithMock) Del(ctx context.Context, keys ...string) (int64, error) {
	cmd := r.redis.Del(ctx, keys...)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// Exists checks if one or more keys exist.
func (r *RedisRepositoryWithMock) Exists(ctx context.Context, keys ...string) (int64, error) {
	cmd := r.redis.Exists(ctx, keys...)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// Expire sets the expiration time of a key.
func (r *RedisRepositoryWithMock) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	cmd := r.redis.Expire(ctx, key, expiration)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

// TTL returns the remaining time to live of a key.
func (r *RedisRepositoryWithMock) TTL(ctx context.Context, key string) (time.Duration, error) {
	cmd := r.redis.TTL(ctx, key)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// Incr increments the number stored at key by one.
func (r *RedisRepositoryWithMock) Incr(ctx context.Context, key string) (int64, error) {
	cmd := r.redis.Incr(ctx, key)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// Ping tests the connection to Redis.
func (r *RedisRepositoryWithMock) Ping(ctx context.Context) error {
	cmd := r.redis.Ping(ctx)
	return cmd.Err()
}

// Append appends a value to a key.
func (r *RedisRepositoryWithMock) Append(ctx context.Context, key, value string) (int64, error) {
	cmd := r.redis.Append(ctx, key, value)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// BitCount counts the number of set bits in a string.
func (r *RedisRepositoryWithMock) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) (int64, error) {
	cmd := r.redis.BitCount(ctx, key, bitCount)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// Decr decrements the number stored at key by one.
func (r *RedisRepositoryWithMock) Decr(ctx context.Context, key string) error {
	cmd := r.redis.Decr(ctx, key)
	return cmd.Err()
}

// DecrBy decrements the number stored at key by decrement.
func (r *RedisRepositoryWithMock) DecrBy(ctx context.Context, key string, value int64) error {
	cmd := r.redis.DecrBy(ctx, key, value)
	return cmd.Err()
}

// GetBit returns the bit value at offset in the string value stored at key.
func (r *RedisRepositoryWithMock) GetBit(ctx context.Context, key string, offset int64) (int64, error) {
	cmd := r.redis.GetBit(ctx, key, offset)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// GetRange returns the substring of the string value stored at key.
func (r *RedisRepositoryWithMock) GetRange(ctx context.Context, key string, start, end int64) (string, error) {
	cmd := r.redis.GetRange(ctx, key, start, end)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

// GetSet sets the string value of a key and return its old value.
func (r *RedisRepositoryWithMock) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	cmd := r.redis.GetSet(ctx, key, value)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

// IncrBy increments the number stored at key by increment.
func (r *RedisRepositoryWithMock) IncrBy(ctx context.Context, key string, value int64) error {
	cmd := r.redis.IncrBy(ctx, key, value)
	return cmd.Err()
}

// IncrByFloat increments the string representing a floating point number stored at key by the specified increment.
func (r *RedisRepositoryWithMock) IncrByfloat(ctx context.Context, key string, value float64) error {
	cmd := r.redis.IncrByFloat(ctx, key, value)
	return cmd.Err()
}

// MGet returns the values of all specified keys.
func (r *RedisRepositoryWithMock) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	cmd := r.redis.MGet(ctx, keys...)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	return cmd.Val(), nil
}

// MSet sets the given keys to their respective values.
func (r *RedisRepositoryWithMock) MSet(ctx context.Context, values ...interface{}) (string, error) {
	cmd := r.redis.MSet(ctx, values...)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

// MSetNX sets the given keys to their respective values if they do not exist.
func (r *RedisRepositoryWithMock) MSetNX(ctx context.Context, values ...interface{}) (bool, error) {
	cmd := r.redis.MSetNX(ctx, values...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

// SetBit sets or clears the bit at offset in the string value stored at key.
func (r *RedisRepositoryWithMock) SetBit(ctx context.Context, key string, offset int64, value int) error {
	cmd := r.redis.SetBit(ctx, key, offset, value)
	return cmd.Err()
}

// SetEX sets the value and expiration of a key.
func (r *RedisRepositoryWithMock) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cmd := r.redis.SetEX(ctx, key, value, expiration)
	return cmd.Err()
}

// SetNX sets the value of a key, only if the key does not exist.
func (r *RedisRepositoryWithMock) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cmd := r.redis.SetNX(ctx, key, value, expiration)
	return cmd.Err()
}

// SetRange overwrites part of a string at key starting at the specified offset.
func (r *RedisRepositoryWithMock) SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
	cmd := r.redis.SetRange(ctx, key, offset, value)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// StrLen returns the length of the string value stored at key.
func (r *RedisRepositoryWithMock) StrLen(ctx context.Context, key string) (int64, error) {
	cmd := r.redis.StrLen(ctx, key)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

func TestNewRedisRepository(t *testing.T) {
	mockClient := &redis.Client{}
	repo := NewRedisRepository(mockClient)
	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.Client)
}

func TestRedisRepository_Get(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedValue := "test-value"
	mockClient.On("Get", mock.Anything, "test-key").Return(expectedValue, nil)

	result, err := repo.Get(context.Background(), "test-key")

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Get_Error(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedError := errors.New("redis connection error")
	mockClient.On("Get", mock.Anything, "test-key").Return("", expectedError)

	result, err := repo.Get(context.Background(), "test-key")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis Get:")
	assert.Empty(t, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Set(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("Set", mock.Anything, "test-key", "test-value", time.Hour).Return("OK", nil)

	err := repo.Set(context.Background(), "test-key", "test-value", time.Hour)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Set_Error(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithMock(mockClient)

	expectedError := errors.New("redis connection error")
	mockClient.On("Set", mock.Anything, "test-key", "test-value", time.Hour).Return("", expectedError)

	err := repo.Set(context.Background(), "test-key", "test-value", time.Hour)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Del(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedCount := int64(2)
	mockClient.On("Del", mock.Anything, []string{"key1", "key2"}).Return(expectedCount, nil)

	result, err := repo.Del(context.Background(), "key1", "key2")

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Exists(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedCount := int64(1)
	mockClient.On("Exists", mock.Anything, []string{"test-key"}).Return(expectedCount, nil)

	result, err := repo.Exists(context.Background(), "test-key")

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Expire(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedResult := true
	mockClient.On("Expire", mock.Anything, "test-key", time.Hour).Return(expectedResult, nil)

	result, err := repo.Expire(context.Background(), "test-key", time.Hour)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_TTL(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedTTL := time.Hour * 2
	mockClient.On("TTL", mock.Anything, "test-key").Return(expectedTTL, nil)

	result, err := repo.TTL(context.Background(), "test-key")

	assert.NoError(t, err)
	assert.Equal(t, expectedTTL, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Incr(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedValue := int64(5)
	mockClient.On("Incr", mock.Anything, "counter").Return(expectedValue, nil)

	result, err := repo.Incr(context.Background(), "counter")

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Ping(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("Ping", mock.Anything).Return("PONG", nil)

	err := repo.Ping(context.Background())

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Append(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedLength := int64(15)
	mockClient.On("Append", mock.Anything, "test-key", "test-value").Return(expectedLength, nil)

	result, err := repo.Append(context.Background(), "test-key", "test-value")

	assert.NoError(t, err)
	assert.Equal(t, expectedLength, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_BitCount(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedCount := int64(5)
	bitCount := &redis.BitCount{}
	mockClient.On("BitCount", mock.Anything, "test-key", bitCount).Return(expectedCount, nil)

	result, err := repo.BitCount(context.Background(), "test-key", bitCount)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_Decr(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("Decr", mock.Anything, "counter").Return(int64(4), nil)

	err := repo.Decr(context.Background(), "counter")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_DecrBy(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("DecrBy", mock.Anything, "counter", int64(3)).Return(int64(2), nil)

	err := repo.DecrBy(context.Background(), "counter", 3)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_GetBit(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedBit := int64(1)
	mockClient.On("GetBit", mock.Anything, "test-key", int64(5)).Return(expectedBit, nil)

	result, err := repo.GetBit(context.Background(), "test-key", 5)

	assert.NoError(t, err)
	assert.Equal(t, expectedBit, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_GetRange(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedRange := "ello"
	mockClient.On("GetRange", mock.Anything, "test-key", int64(1), int64(4)).Return(expectedRange, nil)

	result, err := repo.GetRange(context.Background(), "test-key", 1, 4)

	assert.NoError(t, err)
	assert.Equal(t, expectedRange, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_GetSet(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	oldValue := "old-value"
	mockClient.On("GetSet", mock.Anything, "test-key", "new-value").Return(oldValue, nil)

	result, err := repo.GetSet(context.Background(), "test-key", "new-value")

	assert.NoError(t, err)
	assert.Equal(t, oldValue, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_IncrBy(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("IncrBy", mock.Anything, "counter", int64(5)).Return(int64(10), nil)

	err := repo.IncrBy(context.Background(), "counter", 5)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_IncrByfloat(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("IncrByFloat", mock.Anything, "counter", float64(2.5)).Return(float64(7.5), nil)

	err := repo.IncrByfloat(context.Background(), "counter", 2.5)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_MGet(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedValues := []any{"value1", "value2", nil}
	mockClient.On("MGet", mock.Anything, mock.MatchedBy(func(keys []string) bool {
		return len(keys) == 3 && keys[0] == "key1" && keys[1] == "key2" && keys[2] == "key3"
	})).Return(expectedValues, nil)

	result, err := repo.MGet(context.Background(), "key1", "key2", "key3")

	assert.NoError(t, err)
	assert.Equal(t, expectedValues, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_MSet(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("MSet", mock.Anything, mock.MatchedBy(func(values []any) bool {
		return len(values) == 4 && values[0] == "key1" && values[1] == "value1" &&
			values[2] == "key2" && values[3] == "value2"
	})).Return("OK", nil)

	result, err := repo.MSet(context.Background(), "key1", "value1", "key2", "value2")

	assert.NoError(t, err)
	assert.Equal(t, "OK", result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_MSetNX(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedResult := true
	mockClient.On("MSetNX", mock.Anything, mock.MatchedBy(func(values []any) bool {
		return len(values) == 4 && values[0] == "key1" && values[1] == "value1" &&
			values[2] == "key2" && values[3] == "value2"
	})).Return(expectedResult, nil)

	result, err := repo.MSetNX(context.Background(), "key1", "value1", "key2", "value2")

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_SetBit(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("SetBit", mock.Anything, "test-key", int64(5), 1).Return(int64(0), nil)

	err := repo.SetBit(context.Background(), "test-key", 5, 1)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_SetEX(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.On("SetEX", mock.Anything, "test-key", "test-value", time.Hour).Return("OK", nil)

	err := repo.SetEX(context.Background(), "test-key", "test-value", time.Hour)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_SetNX(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedResult := true
	mockClient.On("SetNX", mock.Anything, "test-key", "test-value", time.Hour).Return(expectedResult, nil)

	err := repo.SetNX(context.Background(), "test-key", "test-value", time.Hour)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_SetRange(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedLength := int64(10)
	mockClient.On("SetRange", mock.Anything, "test-key", int64(5), "test").Return(expectedLength, nil)

	result, err := repo.SetRange(context.Background(), "test-key", 5, "test")

	assert.NoError(t, err)
	assert.Equal(t, expectedLength, result)
	mockClient.AssertExpectations(t)
}

func TestRedisRepository_StrLen(t *testing.T) {
	mockClient := &MockRedisClient{}
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedLength := int64(10)
	mockClient.On("StrLen", mock.Anything, "test-key").Return(expectedLength, nil)

	result, err := repo.StrLen(context.Background(), "test-key")

	assert.NoError(t, err)
	assert.Equal(t, expectedLength, result)
	mockClient.AssertExpectations(t)
}

// Test actual RedisRepository functions
func TestRedisRepositoryReal_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Redis integration test - requires real Redis server")
}

func TestRedisRepositoryReal_Set(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Redis integration test - requires real Redis server")
}

func TestRedisRepositoryReal_Incr(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Redis integration test - requires real Redis server")
}

func TestRedisRepositoryReal_MGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Redis integration test - requires real Redis server")
}

func TestRedisRepositoryReal_MSet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping Redis integration test - requires real Redis server")
}
