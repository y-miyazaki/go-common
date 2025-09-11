package transport

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y-miyazaki/go-common/pkg/logger"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewTransportHTTPLogger(t *testing.T) {
	l := logger.NewLogger(logrus.New())
	transport := NewTransportHTTPLogger(l, HTTPLoggerTypeExternal)
	assert.NotNil(t, transport)
	assert.IsType(t, &HTTPLogger{}, transport)
}

func TestHTTPLogger_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	l := logger.NewLogger(logrus.New())
	transport := NewTransportHTTPLogger(l, HTTPLoggerTypeExternal).(*HTTPLogger)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := transport.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPLogger_RoundTrip_Error(t *testing.T) {
	l := logger.NewLogger(logrus.New())
	transport := NewTransportHTTPLogger(l, HTTPLoggerTypeExternal).(*HTTPLogger)

	req, _ := http.NewRequest("GET", "http://invalid-url", nil)
	_, err := transport.RoundTrip(req)

	assert.Error(t, err)
}

func TestHTTPLoggerType(t *testing.T) {
	assert.Equal(t, HTTPLoggerType("external"), HTTPLoggerTypeExternal)
	assert.Equal(t, HTTPLoggerType("internal"), HTTPLoggerTypeInternal)
}

func TestNewTransportHTTPZapLogger(t *testing.T) {
	zapConfig := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
	}
	l := logger.NewZapLogger(zapConfig)
	transport := NewTransportHTTPZapLogger(l, HTTPLoggerTypeExternal)
	assert.NotNil(t, transport)
	assert.IsType(t, &HTTPZapLogger{}, transport)
}

func TestHTTPZapLogger_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	zapConfig := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
	}
	l := logger.NewZapLogger(zapConfig)
	transport := NewTransportHTTPZapLogger(l, HTTPLoggerTypeExternal).(*HTTPZapLogger)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := transport.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPZapLogger_RoundTrip_Error(t *testing.T) {
	zapConfig := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
	}
	l := logger.NewZapLogger(zapConfig)
	transport := NewTransportHTTPZapLogger(l, HTTPLoggerTypeExternal).(*HTTPZapLogger)

	req, _ := http.NewRequest("GET", "http://invalid-url", nil)
	_, err := transport.RoundTrip(req)

	assert.Error(t, err)
}

func TestNewTransportHTTPSlogLogger(t *testing.T) {
	config := &logger.SlogConfig{
		Level:     logger.Level(slog.LevelInfo),
		Format:    "json",
		AddSource: false,
		Output:    io.Discard,
	}
	l := logger.NewSlogLogger(config)
	transport := NewTransportHTTPSlogLogger(l, HTTPLoggerTypeExternal)
	assert.NotNil(t, transport)
	assert.IsType(t, &HTTPSlogLogger{}, transport)
}

func TestHTTPSlogLogger_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	config := &logger.SlogConfig{
		Level:     logger.Level(slog.LevelInfo),
		Format:    "json",
		AddSource: false,
		Output:    io.Discard,
	}
	l := logger.NewSlogLogger(config)
	transport := NewTransportHTTPSlogLogger(l, HTTPLoggerTypeExternal).(*HTTPSlogLogger)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := transport.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPSlogLogger_RoundTrip_Error(t *testing.T) {
	config := &logger.SlogConfig{
		Level:     logger.Level(slog.LevelInfo),
		Format:    "json",
		AddSource: false,
		Output:    io.Discard,
	}
	l := logger.NewSlogLogger(config)
	transport := NewTransportHTTPSlogLogger(l, HTTPLoggerTypeExternal).(*HTTPSlogLogger)

	req, _ := http.NewRequest("GET", "http://invalid-url", nil)
	_, err := transport.RoundTrip(req)

	assert.Error(t, err)
}
