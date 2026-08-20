//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/repository"
	"gorm.io/gorm"
)

func TestNewHTTPHandler(t *testing.T) {
	t.Parallel()

	mockLogger := &logger.Logger{}
	mysqlDB := &gorm.DB{}
	postgresDB := &gorm.DB{}
	awsS3Repo := &repository.AWSS3Repository{}
	redisRepo := &repository.RedisRepository{}

	tests := []struct {
		h    *HTTPHandler
		name string
	}{
		{
			name: "constructor wires dependencies",
			h:    NewHTTPHandler(mockLogger, mysqlDB, postgresDB, awsS3Repo, redisRepo),
		},
		{
			name: "nil dependencies are accepted",
			h:    NewHTTPHandler(nil, nil, nil, nil, nil),
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.NotNil(t, tt.h)
			require.NotNil(t, tt.h.BaseHTTPHandler)
		})
	}
}

func TestHTTPHandler_Struct(t *testing.T) {
	t.Parallel()

	h := &HTTPHandler{}
	require.NotNil(t, h)
	require.Nil(t, h.BaseHTTPHandler)
	require.Nil(t, h.mysqlDB)
	require.Nil(t, h.postgresDB)
	require.Nil(t, h.awsS3Repository)
	require.Nil(t, h.redisRepository)
}
