// pragma: allowlist-secret
//
//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/handler"
)

func TestHTTPHandler_HandleEnv(t *testing.T) {
	originalPassword := os.Getenv("APP_DATABASE_MASTER_PASSWORD")
	originalAddr := os.Getenv("REDIS_ADDR")
	defer func() {
		os.Setenv("APP_DATABASE_MASTER_PASSWORD", originalPassword)
		os.Setenv("REDIS_ADDR", originalAddr)
	}()

	tests := []struct {
		name     string
		password string
		addr     string
	}{
		{
			name:     "returns env vars in response",
			password: "test_password", // pragma: allowlist-secret
			addr:     "localhost:6379",
		},
		{
			name:     "returns empty strings when unset",
			password: "",
			addr:     "",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_DATABASE_MASTER_PASSWORD", tt.password)
			os.Setenv("REDIS_ADDR", tt.addr)

			h := &HTTPHandler{BaseHTTPHandler: &handler.BaseHTTPHandler{}}
			status, body := invokeHandler(t, h.HandleEnv)

			require.Equal(t, http.StatusOK, status)
			require.Contains(t, body, "Hello!")
			require.Contains(t, body, tt.password)
			require.Contains(t, body, tt.addr)
		})
	}
}

func TestHTTPHandler_HandleEnv_Route(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{name: "post method returns not found", method: http.MethodPost, wantStatus: http.StatusNotFound},
	}

	h := &HTTPHandler{}
	router := gin.New()
	router.GET("/env", h.HandleEnv)

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/env", http.NoBody)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
