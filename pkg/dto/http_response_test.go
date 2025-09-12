package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPBaseErrorResponse(t *testing.T) {
	response := &HTTPBaseErrorResponse{
		Error: &HTTPErrorResponse{
			Message: "test error",
		},
	}

	data, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "test error")
}

func TestHTTPErrorResponse(t *testing.T) {
	errorResponse := &HTTPErrorResponse{
		Message: map[string]string{"code": "400", "description": "Bad Request"},
	}

	data, err := json.Marshal(errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Bad Request")
}

func TestHTTPErrorResponseStringMessage(t *testing.T) {
	errorResponse := &HTTPErrorResponse{
		Message: "simple error message",
	}

	data, err := json.Marshal(errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "simple error message")
}
