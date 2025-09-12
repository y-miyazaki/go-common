package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTPHandler_SayHello(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance
	h := &HTTPHandler{}

	// Register the route
	router.GET("/hello", h.SayHello)

	// Create a test request
	req, err := http.NewRequest("GET", "/hello", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hello!")
	assert.Contains(t, w.Body.String(), "message")
}

func TestHTTPHandler_SayHello_POST(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create HTTPHandler instance
	h := &HTTPHandler{}

	// Register the route
	router.GET("/hello", h.SayHello)

	// Create a test request with POST method (should not match)
	req, err := http.NewRequest("POST", "/hello", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response - should be 404 Not Found for unmatched route
	assert.Equal(t, http.StatusNotFound, w.Code)
}
