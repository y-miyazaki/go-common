package repository

import "github.com/y-miyazaki/go-common/pkg/logger"

// BaseRepository struct.
type BaseRepository struct {
	Logger *logger.Logger
}

// NewBaseRepository returns BaseRepository instance.
func NewBaseRepository(logger *logger.Logger) *BaseRepository {
	return &BaseRepository{Logger: logger}
}
