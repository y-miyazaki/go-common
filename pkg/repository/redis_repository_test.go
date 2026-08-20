package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
	"go.uber.org/mock/gomock"
)

var errTestRedis = errors.New("redis connection error")

func redisStringResult(val string) *redis.StringCmd {
	cmd := redis.NewStringCmd(context.Background())
	cmd.SetVal(val)
	return cmd
}

func redisStringError(err error) *redis.StringCmd {
	cmd := redis.NewStringCmd(context.Background())
	cmd.SetErr(err)
	return cmd
}

func redisStatusResult(val string) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetVal(val)
	return cmd
}

func redisStatusError(err error) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetErr(err)
	return cmd
}

func redisIntResult(val int64) *redis.IntCmd {
	cmd := redis.NewIntCmd(context.Background())
	cmd.SetVal(val)
	return cmd
}


func redisBoolResult(val bool) *redis.BoolCmd {
	cmd := redis.NewBoolCmd(context.Background())
	cmd.SetVal(val)
	return cmd
}

func redisDurationResult(val time.Duration) *redis.DurationCmd {
	cmd := redis.NewDurationCmd(context.Background(), 0)
	cmd.SetVal(val)
	return cmd
}

func redisFloatResult(val float64) *redis.FloatCmd {
	cmd := redis.NewFloatCmd(context.Background())
	cmd.SetVal(val)
	return cmd
}

func redisSliceResult(val []any) *redis.SliceCmd {
	cmd := redis.NewSliceCmd(context.Background())
	cmd.SetVal(val)
	return cmd
}

func TestNewRedisRepository(t *testing.T) {
	t.Parallel()

	mockClient := &redis.Client{}
	repo := NewRedisRepository(mockClient)
	require.NotNil(t, repo)
	require.Equal(t, mockClient, repo.Client)
}

func TestRedisRepository_Get(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedValue := "test-value"
	mockClient.EXPECT().
		Get(gomock.Any(), "test-key").
		Return(redisStringResult(expectedValue))

	result, err := repo.Get(context.Background(), "test-key")

	require.NoError(t, err)
	require.Equal(t, expectedValue, result)
}

func TestRedisRepository_Get_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		Get(gomock.Any(), "test-key").
		Return(redisStringError(errTestRedis))

	result, err := repo.Get(context.Background(), "test-key")

	require.Error(t, err)
	require.Contains(t, err.Error(), "redis Get:")
	require.Empty(t, result)
	require.ErrorIs(t, err, errTestRedis)
}

func TestRedisRepository_Set(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		Set(gomock.Any(), "test-key", "test-value", time.Hour).
		Return(redisStatusResult("OK"))

	err := repo.Set(context.Background(), "test-key", "test-value", time.Hour)

	require.NoError(t, err)
}

func TestRedisRepository_Set_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		Set(gomock.Any(), "test-key", "test-value", time.Hour).
		Return(redisStatusError(errTestRedis))

	err := repo.Set(context.Background(), "test-key", "test-value", time.Hour)

	require.Error(t, err)
	require.Contains(t, err.Error(), "redis Set:")
	require.ErrorIs(t, err, errTestRedis)
}

func TestRedisRepository_Del(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedCount := int64(2)
	mockClient.EXPECT().
		Del(gomock.Any(), "key1", "key2").
		Return(redisIntResult(expectedCount))

	result, err := repo.Del(context.Background(), "key1", "key2")

	require.NoError(t, err)
	require.Equal(t, expectedCount, result)
}

func TestRedisRepository_Exists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedCount := int64(1)
	mockClient.EXPECT().
		Exists(gomock.Any(), "test-key").
		Return(redisIntResult(expectedCount))

	result, err := repo.Exists(context.Background(), "test-key")

	require.NoError(t, err)
	require.Equal(t, expectedCount, result)
}

func TestRedisRepository_Expire(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedResult := true
	mockClient.EXPECT().
		Expire(gomock.Any(), "test-key", time.Hour).
		Return(redisBoolResult(expectedResult))

	result, err := repo.Expire(context.Background(), "test-key", time.Hour)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestRedisRepository_TTL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedTTL := time.Hour * 2
	mockClient.EXPECT().
		TTL(gomock.Any(), "test-key").
		Return(redisDurationResult(expectedTTL))

	result, err := repo.TTL(context.Background(), "test-key")

	require.NoError(t, err)
	require.Equal(t, expectedTTL, result)
}

func TestRedisRepository_Incr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedValue := int64(5)
	mockClient.EXPECT().
		Incr(gomock.Any(), "counter").
		Return(redisIntResult(expectedValue))

	result, err := repo.Incr(context.Background(), "counter")

	require.NoError(t, err)
	require.Equal(t, expectedValue, result)
}

func TestRedisRepository_Ping(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		Ping(gomock.Any()).
		Return(redisStatusResult("PONG"))

	err := repo.Ping(context.Background())

	require.NoError(t, err)
}

func TestRedisRepository_Append(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedLength := int64(15)
	mockClient.EXPECT().
		Append(gomock.Any(), "test-key", "test-value").
		Return(redisIntResult(expectedLength))

	result, err := repo.Append(context.Background(), "test-key", "test-value")

	require.NoError(t, err)
	require.Equal(t, expectedLength, result)
}

func TestRedisRepository_BitCount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedCount := int64(5)
	bitCount := &redis.BitCount{}
	mockClient.EXPECT().
		BitCount(gomock.Any(), "test-key", bitCount).
		Return(redisIntResult(expectedCount))

	result, err := repo.BitCount(context.Background(), "test-key", bitCount)

	require.NoError(t, err)
	require.Equal(t, expectedCount, result)
}

func TestRedisRepository_Decr(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		Decr(gomock.Any(), "counter").
		Return(redisIntResult(4))

	err := repo.Decr(context.Background(), "counter")

	require.NoError(t, err)
}

func TestRedisRepository_DecrBy(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		DecrBy(gomock.Any(), "counter", int64(3)).
		Return(redisIntResult(2))

	err := repo.DecrBy(context.Background(), "counter", 3)

	require.NoError(t, err)
}

func TestRedisRepository_GetBit(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedBit := int64(1)
	mockClient.EXPECT().
		GetBit(gomock.Any(), "test-key", int64(5)).
		Return(redisIntResult(expectedBit))

	result, err := repo.GetBit(context.Background(), "test-key", 5)

	require.NoError(t, err)
	require.Equal(t, expectedBit, result)
}

func TestRedisRepository_GetRange(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedRange := "ello"
	mockClient.EXPECT().
		GetRange(gomock.Any(), "test-key", int64(1), int64(4)).
		Return(redisStringResult(expectedRange))

	result, err := repo.GetRange(context.Background(), "test-key", 1, 4)

	require.NoError(t, err)
	require.Equal(t, expectedRange, result)
}

func TestRedisRepository_GetSet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	oldValue := "old-value"
	mockClient.EXPECT().
		GetSet(gomock.Any(), "test-key", "new-value").
		Return(redisStringResult(oldValue))

	result, err := repo.GetSet(context.Background(), "test-key", "new-value")

	require.NoError(t, err)
	require.Equal(t, oldValue, result)
}

func TestRedisRepository_IncrBy(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		IncrBy(gomock.Any(), "counter", int64(5)).
		Return(redisIntResult(10))

	err := repo.IncrBy(context.Background(), "counter", 5)

	require.NoError(t, err)
}

func TestRedisRepository_IncrByfloat(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		IncrByFloat(gomock.Any(), "counter", float64(2.5)).
		Return(redisFloatResult(7.5))

	err := repo.IncrByfloat(context.Background(), "counter", 2.5)

	require.NoError(t, err)
}

func TestRedisRepository_MGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedValues := []any{"value1", "value2", nil}
	mockClient.EXPECT().
		MGet(gomock.Any(), "key1", "key2", "key3").
		Return(redisSliceResult(expectedValues))

	result, err := repo.MGet(context.Background(), "key1", "key2", "key3")

	require.NoError(t, err)
	require.Equal(t, expectedValues, result)
}

func TestRedisRepository_MSet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		MSet(gomock.Any(), "key1", "value1", "key2", "value2").
		Return(redisStatusResult("OK"))

	result, err := repo.MSet(context.Background(), "key1", "value1", "key2", "value2")

	require.NoError(t, err)
	require.Equal(t, "OK", result)
}

func TestRedisRepository_MSetNX(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedResult := true
	mockClient.EXPECT().
		MSetNX(gomock.Any(), "key1", "value1", "key2", "value2").
		Return(redisBoolResult(expectedResult))

	result, err := repo.MSetNX(context.Background(), "key1", "value1", "key2", "value2")

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestRedisRepository_SetBit(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		SetBit(gomock.Any(), "test-key", int64(5), 1).
		Return(redisIntResult(0))

	err := repo.SetBit(context.Background(), "test-key", 5, 1)

	require.NoError(t, err)
}

func TestRedisRepository_SetEX(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	mockClient.EXPECT().
		SetEx(gomock.Any(), "test-key", "test-value", time.Hour).
		Return(redisStatusResult("OK"))

	err := repo.SetEX(context.Background(), "test-key", "test-value", time.Hour)

	require.NoError(t, err)
}

func TestRedisRepository_SetNX(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedResult := true
	mockClient.EXPECT().
		SetNX(gomock.Any(), "test-key", "test-value", time.Hour).
		Return(redisBoolResult(expectedResult))

	err := repo.SetNX(context.Background(), "test-key", "test-value", time.Hour)

	require.NoError(t, err)
}

func TestRedisRepository_SetRange(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedLength := int64(10)
	mockClient.EXPECT().
		SetRange(gomock.Any(), "test-key", int64(5), "test").
		Return(redisIntResult(expectedLength))

	result, err := repo.SetRange(context.Background(), "test-key", 5, "test")

	require.NoError(t, err)
	require.Equal(t, expectedLength, result)
}

func TestRedisRepository_StrLen(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := NewRedisRepositoryWithInterface(mockClient)

	expectedLength := int64(10)
	mockClient.EXPECT().
		StrLen(gomock.Any(), "test-key").
		Return(redisIntResult(expectedLength))

	result, err := repo.StrLen(context.Background(), "test-key")

	require.NoError(t, err)
	require.Equal(t, expectedLength, result)
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
