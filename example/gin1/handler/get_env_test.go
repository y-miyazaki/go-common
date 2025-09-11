// pragma: allowlist-secret
package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/y-miyazaki/go-common/pkg/handler"
)

func TestHTTPHandler_GetEnv(t *testing.T) {
	// Set up test environment variables
	originalPassword := os.Getenv("APP_DATABASE_MASTER_PASSWORD")
	originalAddr := os.Getenv("REDIS_ADDR")

	defer func() {
		// Restore original environment variables
		os.Setenv("APP_DATABASE_MASTER_PASSWORD", originalPassword)
		os.Setenv("REDIS_ADDR", originalAddr)
	}()

	// Set test values
	testPassword := "test_password_123" // pragma: allowlist-secret
	testAddr := "localhost:6379"
	os.Setenv("APP_DATABASE_MASTER_PASSWORD", testPassword)
	os.Setenv("REDIS_ADDR", testAddr)

	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance with proper BaseHTTPHandler
	h := &HTTPHandler{
		BaseHTTPHandler: &handler.BaseHTTPHandler{},
	}

	// Register the route
	router.GET("/env", h.HandleEnv)

	// Create a test request
	req, err := http.NewRequest("GET", "/env", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hello!")
	assert.Contains(t, w.Body.String(), testPassword)
	assert.Contains(t, w.Body.String(), testAddr)
	assert.Contains(t, w.Body.String(), "message")
	assert.Contains(t, w.Body.String(), "password")
	assert.Contains(t, w.Body.String(), "addr")
}

func TestHTTPHandler_GetEnv_EmptyEnvVars(t *testing.T) {
	// Clear environment variables
	originalPassword := os.Getenv("APP_DATABASE_MASTER_PASSWORD")
	originalAddr := os.Getenv("REDIS_ADDR")

	defer func() {
		// Restore original environment variables
		os.Setenv("APP_DATABASE_MASTER_PASSWORD", originalPassword)
		os.Setenv("REDIS_ADDR", originalAddr)
	}()

	// Clear test values
	os.Unsetenv("APP_DATABASE_MASTER_PASSWORD")
	os.Unsetenv("REDIS_ADDR")

	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance
	h := &HTTPHandler{}

	// Register the route
	router.GET("/env", h.HandleEnv)

	// Create a test request
	req, err := http.NewRequest("GET", "/env", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hello!")
	assert.Contains(t, w.Body.String(), "message")
	assert.Contains(t, w.Body.String(), "password")
	assert.Contains(t, w.Body.String(), "addr")
	// Empty environment variables should result in empty strings
	assert.Contains(t, w.Body.String(), "\"password\":\"\"")
	assert.Contains(t, w.Body.String(), "\"addr\":\"\"")
}

func TestHTTPHandler_GetEnv_POST(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance
	h := &HTTPHandler{}

	// Register the route
	router.GET("/env", h.HandleEnv)

	// Create a test request with POST method (should not match)
	req, err := http.NewRequest("POST", "/env", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response - should be 404 Not Found for unmatched route
	assert.Equal(t, http.StatusNotFound, w.Code)
}
