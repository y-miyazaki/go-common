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
}

// NewHTTPHandler returns HTTPHandler struct.
func NewHTTPHandler(logger *logger.Logger, mysqlDB, postgresDB *gorm.DB, awsS3Repository *repository.AWSS3Repository) *HTTPHandler {
	return &HTTPHandler{
		BaseHTTPHandler: &handler.BaseHTTPHandler{
			Logger: logger,
		},
		mysqlDB:         mysqlDB,
		postgresDB:      postgresDB,
		awsS3Repository: awsS3Repository,
	}
}
