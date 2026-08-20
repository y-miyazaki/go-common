package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBaseRepository(t *testing.T) {
	t.Parallel()

	repo := NewBaseRepository()
	require.NotNil(t, repo)
	require.IsType(t, &BaseRepository{}, repo)
}
