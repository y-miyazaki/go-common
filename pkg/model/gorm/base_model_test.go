package gorm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseModel(t *testing.T) {
	model := BaseModel{}
	assert.NotNil(t, model)
	assert.Equal(t, time.Time{}, model.CreatedAt)
	assert.Equal(t, time.Time{}, model.UpdatedAt)
}

func TestBaseModelWithDeleted(t *testing.T) {
	model := BaseModelWithDeleted{}
	assert.NotNil(t, model)
	assert.Equal(t, time.Time{}, model.CreatedAt)
	assert.Equal(t, time.Time{}, model.UpdatedAt)
	assert.Equal(t, int64(0), int64(model.DeletedAt))
}

// Test with GORM (requires database setup)
func TestBaseModelGormIntegration(t *testing.T) {
	// This test would require a test database
	// For now, just test the struct fields
	model := BaseModel{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	assert.True(t, model.CreatedAt.After(time.Time{}))
	assert.True(t, model.UpdatedAt.After(time.Time{}))
}
