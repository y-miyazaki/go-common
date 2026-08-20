//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHTTPHandler_HandleMySQL(t *testing.T) {
	tests := []struct {
		setup      func(t *testing.T) *HTTPHandler
		name       string
		wantBody   string
		wantStatus int
	}{
		{
			name: "creates user and returns ok",
			setup: func(t *testing.T) *HTTPHandler {
				return NewHTTPHandler(newTestLogger(), newMySQLGormDB(t), nil, nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"message":"ok"`,
		},
		{
			name: "create table error returns ok without panic",
			setup: func(t *testing.T) *HTTPHandler {
				return NewHTTPHandler(newTestLogger(), newMySQLGormDBCreateTableError(t), nil, nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"message":"ok"`,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			h := tt.setup(t)
			status, body := invokeHandler(t, h.HandleMySQL)
			require.Equal(t, tt.wantStatus, status)
			require.Contains(t, body, tt.wantBody)
		})
	}
}

func TestHTTPHandler_HandleMySQL_Route(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{name: "post method returns not found", method: http.MethodPost, wantStatus: http.StatusNotFound},
	}

	h := &HTTPHandler{}
	router := gin.New()
	router.GET("/mysql", h.HandleMySQL)

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/mysql", http.NoBody)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
