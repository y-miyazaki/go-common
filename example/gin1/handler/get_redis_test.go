//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"github.com/y-miyazaki/go-common/pkg/repository/mocks"
	"go.uber.org/mock/gomock"
)

var errTestRedisSet = errors.New("redis set failed")

func TestHTTPHandler_HandleRedis(t *testing.T) {
	tests := []struct {
		setup      func(t *testing.T, ctrl *gomock.Controller) *HTTPHandler
		name       string
		wantBody   string
		wantStatus int
	}{
		{
			name: "set and get key returns ok",
			setup: func(t *testing.T, ctrl *gomock.Controller) *HTTPHandler {
				mockRedis := mocks.NewMockRedisClientInterface(ctrl)
				setCmd := redis.NewStatusCmd(context.Background())
				setCmd.SetVal("OK")
				getCmd := redis.NewStringCmd(context.Background())
				getCmd.SetVal("1")

				mockRedis.EXPECT().
					Set(gomock.Any(), "a", 1, time.Duration(0)).
					Return(setCmd)
				mockRedis.EXPECT().
					Get(gomock.Any(), "a").
					Return(getCmd)

				repo := repository.NewRedisRepositoryWithInterface(mockRedis)
				return NewHTTPHandler(newTestLogger(), nil, nil, nil, repo)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"message":"ok"`,
		},
		{
			name: "redis errors still return ok response",
			setup: func(t *testing.T, ctrl *gomock.Controller) *HTTPHandler {
				mockRedis := mocks.NewMockRedisClientInterface(ctrl)
				setCmd := redis.NewStatusCmd(context.Background())
				setCmd.SetErr(errTestRedisSet)
				getCmd := redis.NewStringCmd(context.Background())
				getCmd.SetErr(errTestRedisSet)

				mockRedis.EXPECT().
					Set(gomock.Any(), "a", 1, time.Duration(0)).
					Return(setCmd)
				mockRedis.EXPECT().
					Get(gomock.Any(), "a").
					Return(getCmd)

				repo := repository.NewRedisRepositoryWithInterface(mockRedis)
				return NewHTTPHandler(newTestLogger(), nil, nil, nil, repo)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"message":"ok"`,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			h := tt.setup(t, ctrl)
			status, body := invokeHandler(t, h.HandleRedis)
			require.Equal(t, tt.wantStatus, status)
			require.Contains(t, body, tt.wantBody)
		})
	}
}

func TestHTTPHandler_HandleRedis_Route(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{name: "post method returns not found", method: http.MethodPost, wantStatus: http.StatusNotFound},
	}

	h := &HTTPHandler{}
	router := gin.New()
	router.GET("/redis", h.HandleRedis)

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/redis", http.NoBody)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
