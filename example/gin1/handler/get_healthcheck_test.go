//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHTTPHandler_HealthCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantBody   string
		wantStatus int
	}{
		{name: "returns ok message", wantStatus: http.StatusOK, wantBody: `"message":"ok"`},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &HTTPHandler{}
			status, body := invokeHandler(t, h.HealthCheck)
			require.Equal(t, tt.wantStatus, status)
			require.Contains(t, body, tt.wantBody)
		})
	}
}

func TestHTTPHandler_HealthCheck_Route(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{name: "post method returns not found", method: http.MethodPost, wantStatus: http.StatusNotFound},
	}

	h := &HTTPHandler{}
	router := gin.New()
	router.GET("/health", h.HealthCheck)

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", http.NoBody)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
