//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPBaseErrorResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		response   *HTTPBaseErrorResponse
		wantSubstr string
	}{
		{
			name: "marshals nested error message",
			response: &HTTPBaseErrorResponse{
				Error: &HTTPErrorResponse{Message: "test error"},
			},
			wantSubstr: "test error",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(tt.response)
			require.NoError(t, err)
			require.Contains(t, string(data), tt.wantSubstr)
		})
	}
}

func TestHTTPErrorResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		response   *HTTPErrorResponse
		wantSubstr string
	}{
		{
			name: "marshals map message",
			response: &HTTPErrorResponse{
				Message: map[string]string{"code": "400", "description": "Bad Request"},
			},
			wantSubstr: "Bad Request",
		},
		{
			name:       "marshals string message",
			response:   &HTTPErrorResponse{Message: "simple error message"},
			wantSubstr: "simple error message",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(tt.response)
			require.NoError(t, err)
			require.Contains(t, string(data), tt.wantSubstr)
		})
	}
}
