// Package handler provides HTTP request handlers for the Gin web framework.
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"gorm.io/gorm"
)

func TestHTTPHandler_HandleEnv(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test environment variables
	originalPassword := os.Getenv("APP_DATABASE_MASTER_PASSWORD")
	originalAddr := os.Getenv("REDIS_ADDR")
	defer func() {
		os.Setenv("APP_DATABASE_MASTER_PASSWORD", originalPassword)
		os.Setenv("REDIS_ADDR", originalAddr)
	}()

	tests := []struct {
		name         string
		password     string
		addr         string
		expectedCode int
	}{
		{
			name:         "success with env vars",
			password:     "test_password", // pragma: allowlist-secret
			addr:         "localhost:6379",
			expectedCode: http.StatusOK,
		},
		{
			name:         "success with empty env vars",
			password:     "",
			addr:         "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("APP_DATABASE_MASTER_PASSWORD", tt.password)
			os.Setenv("REDIS_ADDR", tt.addr)

			// Create handler
			mockLogger := &logger.Logger{}
			h := NewHTTPHandler(mockLogger, nil, nil, nil, nil)

			// Create test request
			req, _ := http.NewRequest(http.MethodGet, "/env", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			h.HandleEnv(c)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)

			// Parse response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Hello!", response["message"])
			assert.Equal(t, tt.password, response["password"])
			assert.Equal(t, tt.addr, response["addr"])
		})
	}
}

func TestHTTPHandler_NewHTTPHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &logger.Logger{}
	mysqlDB := &gorm.DB{}
	postgresDB := &gorm.DB{}
	awsS3Repo := &repository.AWSS3Repository{}
	redisRepo := &repository.RedisRepository{}

	h := NewHTTPHandler(mockLogger, mysqlDB, postgresDB, awsS3Repo, redisRepo)

	assert.NotNil(t, h)
	assert.Equal(t, mockLogger, h.Logger)
	assert.Equal(t, mysqlDB, h.mysqlDB)
	assert.Equal(t, postgresDB, h.postgresDB)
	assert.Equal(t, awsS3Repo, h.awsS3Repository)
	assert.Equal(t, redisRepo, h.redisRepository)
}
