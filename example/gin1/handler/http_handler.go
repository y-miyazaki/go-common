package handler

import (
	"github.com/y-miyazaki/go-common/pkg/handler"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/repository"

	"gorm.io/gorm"
)

// HTTPHandler struct.
type HTTPHandler struct {
	*handler.BaseHTTPHandler
	mysqlDB         *gorm.DB
	postgresDB      *gorm.DB
	awsS3Repository *repository.AWSS3Repository
	redisRepository *repository.RedisRepository
}

// NewHTTPHandler returns HTTPHandler struct.
func NewHTTPHandler(l *logger.Logger, mysqlDB, postgresDB *gorm.DB, awsS3Repository *repository.AWSS3Repository, redisRepository *repository.RedisRepository) *HTTPHandler {
	return &HTTPHandler{
		BaseHTTPHandler: &handler.BaseHTTPHandler{
			Logger: l,
		},
		mysqlDB:         mysqlDB,
		postgresDB:      postgresDB,
		awsS3Repository: awsS3Repository,
		redisRepository: redisRepository,
	}
}
