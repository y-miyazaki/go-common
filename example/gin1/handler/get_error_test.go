package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y-miyazaki/go-common/pkg/handler"
	"github.com/y-miyazaki/go-common/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHTTPHandler_HandleError1(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance with mock logger
	mockLogger := logger.NewLogger(logrus.New())
	baseHandler := &handler.BaseHTTPHandler{
		Logger: mockLogger,
	}
	h := &HTTPHandler{
		BaseHTTPHandler: baseHandler,
	}

	// Register the route
	router.GET("/error1", h.HandleError1)

	// Create a test request
	req, err := http.NewRequest("GET", "/error1", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Note: The error message might be in the response body or logged
	// Just check that we get an internal server error
	assert.Contains(t, w.Body.String(), "test")
}

func TestHTTPHandler_HandleError2(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance with mock logger
	mockLogger := logger.NewLogger(logrus.New())
	baseHandler := &handler.BaseHTTPHandler{
		Logger: mockLogger,
	}
	h := &HTTPHandler{
		BaseHTTPHandler: baseHandler,
	}

	// Register the route
	router.GET("/error2", h.HandleError2)

	// Create a test request
	req, err := http.NewRequest("GET", "/error2", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Note: The error message might be in the response body or logged
	// Just check that we get an internal server error
	assert.Contains(t, w.Body.String(), "test")
}

func TestHTTPHandler_HandleError1_POST(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance with mock logger
	mockLogger := logger.NewLogger(logrus.New())
	baseHandler := &handler.BaseHTTPHandler{
		Logger: mockLogger,
	}
	h := &HTTPHandler{
		BaseHTTPHandler: baseHandler,
	}

	// Register the route
	router.GET("/error1", h.HandleError1)

	// Create a test request with POST method (should not match)
	req, err := http.NewRequest("POST", "/error1", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response - should be 404 Not Found for unmatched route
	assert.Equal(t, http.StatusNotFound, w.Code)
}
