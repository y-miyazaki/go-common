package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinCors_DefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.NotNil(t, config)
	assert.Contains(t, config.AllowMethods, "GET")
	assert.Contains(t, config.AllowHeaders, "Origin")
	assert.True(t, config.AllowCredentials)
	assert.Equal(t, 86400*time.Second, config.MaxAge)
}

func TestGinCors_AllowAllOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config := &GinCorsConfig{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST"},
		AllowHeaders:    []string{"Content-Type"},
	}
	middleware := GinCors(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Origin", "http://example.com")

	middleware(c)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestGinCors_SpecificOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config := &GinCorsConfig{
		AllowOrigins: []string{"http://example.com"},
		AllowMethods: []string{"GET", "POST"},
	}
	middleware := GinCors(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Origin", "http://example.com")

	middleware(c)

	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestGinCors_InvalidOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config := &GinCorsConfig{
		AllowOrigins: []string{"http://example.com"},
	}
	middleware := GinCors(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Origin", "http://invalid.com")

	middleware(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGinCors_OptionsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config := &GinCorsConfig{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type"},
		MaxAge:          3600 * time.Second,
	}
	middleware := GinCors(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("OPTIONS", "/test", nil)
	c.Request.Header.Set("Origin", "http://example.com")

	middleware(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "GET,POST,OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "3600", w.Header().Get("Access-Control-Max-Age"))
}

func TestNormalize(t *testing.T) {
	input := []string{" GET ", "post", "GET", " Post "}
	expected := []string{"get", "post"}
	result := normalize(input)
	assert.Equal(t, expected, result)
}

func TestConvert(t *testing.T) {
	input := []string{"get", "post"}
	result := convert(input, func(s string) string { return strings.ToUpper(s) })
	expected := []string{"GET", "POST"}
	assert.Equal(t, expected, result)
}
