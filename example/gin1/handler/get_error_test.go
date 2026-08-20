//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/dto"
	"github.com/y-miyazaki/go-common/pkg/handler"
)

func TestHTTPHandler_HandleError(t *testing.T) {
	tests := []struct {
		handler    func(*HTTPHandler, *gin.Context)
		name       string
		wantBody   string
		wantStatus int
	}{
		{
			name:       "handle error1 returns internal server error",
			handler:    func(h *HTTPHandler, c *gin.Context) { h.HandleError1(c) },
			wantStatus: http.StatusInternalServerError,
			wantBody:   "test",
		},
		{
			name:       "handle error2 returns internal server error",
			handler:    func(h *HTTPHandler, c *gin.Context) { h.HandleError2(c) },
			wantStatus: http.StatusInternalServerError,
			wantBody:   "test",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPHandler{
				BaseHTTPHandler: &handler.BaseHTTPHandler{Logger: newTestLogger()},
			}
			status, body := invokeHandler(t, func(c *gin.Context) { tt.handler(h, c) })
			require.Equal(t, tt.wantStatus, status)
			require.Contains(t, body, tt.wantBody)
		})
	}
}

func TestHTTPHandler_HandleError_Route(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		method     string
		wantStatus int
	}{
		{name: "post to error1 returns not found", path: "/error1", method: http.MethodPost, wantStatus: http.StatusNotFound},
	}

	h := &HTTPHandler{BaseHTTPHandler: &handler.BaseHTTPHandler{Logger: newTestLogger()}}
	router := gin.New()
	router.GET("/error1", h.HandleError1)
	router.GET("/error2", h.HandleError2)

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, http.NoBody)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandlerErrorSentinels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "err test1 message", err: ErrTest1, want: "error test1"},
		{name: "err test2 message", err: ErrTest2, want: "error test2"},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.err.Error())
		})
	}
}

// Ensure dto import is used by error handlers at compile time.
var _ = dto.HTTPErrorResponse{}
