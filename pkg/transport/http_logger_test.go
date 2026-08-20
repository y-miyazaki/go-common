//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package transport

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"go.uber.org/zap"
)

func TestHTTPLoggerTypeConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  HTTPLoggerType
		want HTTPLoggerType
	}{
		{name: "external type", got: HTTPLoggerTypeExternal, want: HTTPLoggerType("external")},
		{name: "internal type", got: HTTPLoggerTypeInternal, want: HTTPLoggerType("internal")},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.want {
				t.Fatalf("HTTPLoggerType = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestNewTransportHTTPLogger(t *testing.T) {
	t.Parallel()

	l := logger.NewLogger(logrus.New())
	transport := NewTransportHTTPLogger(l, HTTPLoggerTypeExternal)
	if transport == nil {
		t.Fatal("NewTransportHTTPLogger() = nil, want non-nil")
	}
	if _, ok := transport.(*HTTPLogger); !ok {
		t.Fatalf("NewTransportHTTPLogger() type = %T, want *HTTPLogger", transport)
	}
}

func TestHTTPLogger_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name     string
		url      string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "successful request",
			url:      server.URL,
			wantCode: http.StatusOK,
		},
		{
			name:    "invalid url returns error",
			url:     "https://invalid-url",
			wantErr: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			l := logger.NewLogger(logrus.New())
			transport := NewTransportHTTPLogger(l, HTTPLoggerTypeExternal).(*HTTPLogger)

			req, err := http.NewRequest(http.MethodGet, tt.url, http.NoBody)
			if err != nil {
				t.Fatalf("http.NewRequest() error: %v", err)
			}

			resp, err := transport.RoundTrip(req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("RoundTrip(%q) error = nil, want error", tt.url)
				}
				return
			}
			if err != nil {
				t.Fatalf("RoundTrip(%q) unexpected error: %v", tt.url, err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("RoundTrip(%q) status = %d, want %d", tt.url, resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestNewTransportHTTPZapLogger(t *testing.T) {
	t.Parallel()

	zapConfig := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
	}
	l := logger.NewZapLogger(zapConfig)
	transport := NewTransportHTTPZapLogger(l, HTTPLoggerTypeExternal)
	if transport == nil {
		t.Fatal("NewTransportHTTPZapLogger() = nil, want non-nil")
	}
	if _, ok := transport.(*HTTPZapLogger); !ok {
		t.Fatalf("NewTransportHTTPZapLogger() type = %T, want *HTTPZapLogger", transport)
	}
}

func TestHTTPZapLogger_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name     string
		url      string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "successful request",
			url:      server.URL,
			wantCode: http.StatusOK,
		},
		{
			name:    "invalid url returns error",
			url:     "https://invalid-url",
			wantErr: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			zapConfig := &zap.Config{
				Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
				Encoding: "json",
			}
			l := logger.NewZapLogger(zapConfig)
			transport := NewTransportHTTPZapLogger(l, HTTPLoggerTypeExternal).(*HTTPZapLogger)

			req, err := http.NewRequest(http.MethodGet, tt.url, http.NoBody)
			if err != nil {
				t.Fatalf("http.NewRequest() error: %v", err)
			}

			resp, err := transport.RoundTrip(req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("RoundTrip(%q) error = nil, want error", tt.url)
				}
				return
			}
			if err != nil {
				t.Fatalf("RoundTrip(%q) unexpected error: %v", tt.url, err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("RoundTrip(%q) status = %d, want %d", tt.url, resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestNewTransportHTTPSlogLogger(t *testing.T) {
	t.Parallel()

	config := &logger.SlogConfig{
		Level:     logger.Level(slog.LevelInfo),
		Format:    "json",
		AddSource: false,
		Output:    io.Discard,
	}
	l := logger.NewSlogLogger(config)
	transport := NewTransportHTTPSlogLogger(l, HTTPLoggerTypeExternal)
	if transport == nil {
		t.Fatal("NewTransportHTTPSlogLogger() = nil, want non-nil")
	}
	if _, ok := transport.(*HTTPSlogLogger); !ok {
		t.Fatalf("NewTransportHTTPSlogLogger() type = %T, want *HTTPSlogLogger", transport)
	}
}

func TestHTTPSlogLogger_RoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	tests := []struct {
		name     string
		url      string
		wantErr  bool
		wantCode int
	}{
		{
			name:     "successful request",
			url:      server.URL,
			wantCode: http.StatusOK,
		},
		{
			name:    "invalid url returns error",
			url:     "https://invalid-url",
			wantErr: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			config := &logger.SlogConfig{
				Level:     logger.Level(slog.LevelInfo),
				Format:    "json",
				AddSource: false,
				Output:    io.Discard,
			}
			l := logger.NewSlogLogger(config)
			transport := NewTransportHTTPSlogLogger(l, HTTPLoggerTypeExternal).(*HTTPSlogLogger)

			req, err := http.NewRequest(http.MethodGet, tt.url, http.NoBody)
			if err != nil {
				t.Fatalf("http.NewRequest() error: %v", err)
			}

			resp, err := transport.RoundTrip(req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("RoundTrip(%q) error = nil, want error", tt.url)
				}
				return
			}
			if err != nil {
				t.Fatalf("RoundTrip(%q) unexpected error: %v", tt.url, err)
			}
			if resp.StatusCode != tt.wantCode {
				t.Fatalf("RoundTrip(%q) status = %d, want %d", tt.url, resp.StatusCode, tt.wantCode)
			}
		})
	}
}
