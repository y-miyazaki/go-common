package infrastructure

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	client := NewRedis(options)
	assert.NotNil(t, client)
	assert.IsType(t, &redis.Client{}, client)
}
