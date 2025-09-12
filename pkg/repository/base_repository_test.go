package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseRepository(t *testing.T) {
	repo := NewBaseRepository()
	assert.NotNil(t, repo)
	assert.IsType(t, &BaseRepository{}, repo)
}
