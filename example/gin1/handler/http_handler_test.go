package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/repository"

	"gorm.io/gorm"
)

func TestNewHTTPHandler(t *testing.T) {
	// Create mock dependencies
	mockLogger := &logger.Logger{}
	mockMySQLDB := &gorm.DB{}
	mockPostgresDB := &gorm.DB{}
	mockS3Repo := &repository.AWSS3Repository{}
	mockRedisRepo := &repository.RedisRepository{}

	// Create HTTPHandler using constructor
	handler := NewHTTPHandler(mockLogger, mockMySQLDB, mockPostgresDB, mockS3Repo, mockRedisRepo)

	// Assert that handler is not nil
	assert.NotNil(t, handler)

	// Assert that BaseHTTPHandler is properly initialized
	assert.NotNil(t, handler.BaseHTTPHandler)
	assert.Equal(t, mockLogger, handler.BaseHTTPHandler.Logger)

	// Assert that dependencies are properly set
	assert.Equal(t, mockMySQLDB, handler.mysqlDB)
	assert.Equal(t, mockPostgresDB, handler.postgresDB)
	assert.Equal(t, mockS3Repo, handler.awsS3Repository)
	assert.Equal(t, mockRedisRepo, handler.redisRepository)
}

func TestNewHTTPHandler_NilLogger(t *testing.T) {
	// Create HTTPHandler with nil logger
	handler := NewHTTPHandler(nil, nil, nil, nil, nil)

	// Assert that handler is created but BaseHTTPHandler has nil logger
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.BaseHTTPHandler)
	assert.Nil(t, handler.BaseHTTPHandler.Logger)

	// Assert that other dependencies are nil
	assert.Nil(t, handler.mysqlDB)
	assert.Nil(t, handler.postgresDB)
	assert.Nil(t, handler.awsS3Repository)
	assert.Nil(t, handler.redisRepository)
}

func TestHTTPHandler_Struct(t *testing.T) {
	// Test struct initialization
	h := &HTTPHandler{}

	// Assert that struct is properly initialized with zero values
	assert.NotNil(t, h)
	assert.Nil(t, h.BaseHTTPHandler)
	assert.Nil(t, h.mysqlDB)
	assert.Nil(t, h.postgresDB)
	assert.Nil(t, h.awsS3Repository)
	assert.Nil(t, h.redisRepository)
}
