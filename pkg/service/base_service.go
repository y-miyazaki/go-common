package service

import (
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
)

// BaseService struct.
type BaseService struct {
	Logger *infrastructure.Logger
}

// NewBaseService returns BaseService instance.
func NewBaseService(logger *infrastructure.Logger) *BaseService {
	return &BaseService{Logger: logger}
}
