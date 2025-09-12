package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y-miyazaki/go-common/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGinHTTPLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("X-Trace-ID", "test-trace-id")

	l := logger.NewLogger(logrus.New())
	middleware := GinHTTPLogger(l, "X-Trace-ID", "X-Forwarded-For")

	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGinHTTPLogger_ErrorStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	w.WriteHeader(http.StatusInternalServerError)

	l := logger.NewLogger(logrus.New())
	middleware := GinHTTPLogger(l, "", "")

	middleware(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestClientIP_WithHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("X-Forwarded-For", "192.168.1.1")

	ip := clientIP(c, "X-Forwarded-For")
	assert.Equal(t, "192.168.1.1", ip)
}

func TestClientIP_WithoutHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "127.0.0.1:12345"

	ip := clientIP(c, "X-Forwarded-For")
	assert.Equal(t, "127.0.0.1", ip)
}
