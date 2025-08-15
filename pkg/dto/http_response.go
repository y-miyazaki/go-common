// Package dto defines data transfer objects for HTTP responses.
package dto

// HTTPBaseErrorResponse struct.
type HTTPBaseErrorResponse struct {
	Error *HTTPErrorResponse `json:"error"`
}

// HTTPErrorResponse represents HTTP error response structure.
type HTTPErrorResponse struct {
	Message any `json:"message"`
}
