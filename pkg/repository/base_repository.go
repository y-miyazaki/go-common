package repository

import (
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
)

// BaseRepository struct.
type BaseRepository struct {
	Logger *infrastructure.Logger
}

// NewBaseRepository returns BaseRepository instance.
func NewBaseRepository(logger *infrastructure.Logger) *BaseRepository {
	return &BaseRepository{Logger: logger}
}
