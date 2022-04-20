package repository

import "github.com/y-miyazaki/go-common/pkg/logger"

// BaseRepository struct.
type BaseRepository struct {
	Logger *logger.Logger
}

// NewBaseRepository returns BaseRepository instance.
func NewBaseRepository(l *logger.Logger) *BaseRepository {
	return &BaseRepository{Logger: l}
}
