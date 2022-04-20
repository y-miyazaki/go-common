package service

import "github.com/y-miyazaki/go-common/pkg/logger"

// BaseService struct.
type BaseService struct {
	Logger *logger.Logger
}

// NewBaseService returns BaseService instance.
func NewBaseService(l *logger.Logger) *BaseService {
	return &BaseService{Logger: l}
}
