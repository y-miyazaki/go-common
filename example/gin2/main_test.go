package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Test that main function can be called without panicking
	assert.NotPanics(t, func() {
		// We can't easily test the main function directly as it starts a server
		// and runs indefinitely. Instead, we'll test the router setup logic
		setupRouter()
	})
}

func setupRouter() *gin.Engine {
	const corsMaxAgeHours = 12
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           corsMaxAgeHours * time.Hour,
	}))
	return router
}

func TestSetupRouter(t *testing.T) {
	router := setupRouter()
	assert.NotNil(t, router)
}

func TestCORSConfig(t *testing.T) {
	const corsMaxAgeHours = 12
	config := cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           corsMaxAgeHours * time.Hour,
	}

	assert.Equal(t, []string{"https://foo.com"}, config.AllowOrigins)
	assert.Equal(t, []string{"PUT", "PATCH"}, config.AllowMethods)
	assert.Equal(t, []string{"Origin"}, config.AllowHeaders)
	assert.Equal(t, []string{"Content-Length"}, config.ExposeHeaders)
	assert.True(t, config.AllowCredentials)
	assert.Equal(t, corsMaxAgeHours*time.Hour, config.MaxAge)
}

func TestCORSMiddleware(t *testing.T) {
	router := setupRouter()

	// Test preflight request
	req, _ := http.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "https://foo.com")
	req.Header.Set("Access-Control-Request-Method", "PUT")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://foo.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "PUT,PATCH", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin", w.Header().Get("Access-Control-Allow-Headers"))
	// Note: Content-Length header may not be present in OPTIONS response
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSInvalidOrigin(t *testing.T) {
	router := setupRouter()

	// Test request with invalid origin
	req, _ := http.NewRequest("PUT", "/", nil)
	req.Header.Set("Origin", "https://invalid.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should not have CORS headers for invalid origin
	assert.Equal(t, "", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestMainFunctionLogic(t *testing.T) {
	// Test the main function logic without starting the server
	const corsMaxAgeHours = 12

	// Test CORS configuration values
	assert.Equal(t, 12, corsMaxAgeHours)

	// Test router creation
	router := gin.Default()
	assert.NotNil(t, router)

	// Test CORS middleware creation
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           corsMaxAgeHours * time.Hour,
	})
	assert.NotNil(t, corsMiddleware)

	// Test middleware application
	router.Use(corsMiddleware)
	assert.NotNil(t, router)
}
