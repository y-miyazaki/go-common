package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseService(t *testing.T) {
	service := NewBaseService()
	assert.NotNil(t, service)
	assert.IsType(t, &BaseService{}, service)
}

func TestBaseService(t *testing.T) {
	service := &BaseService{}
	assert.NotNil(t, service)
}
