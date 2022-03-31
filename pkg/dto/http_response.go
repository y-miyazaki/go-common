package dto

// HTTPBaseErrorResponse struct.
type HTTPBaseErrorResponse struct {
	Error *HTTPErrorResponse `json:"error"`
}

// HTTPErrorResponse struct.
type HTTPErrorResponse struct {
	Message interface{} `json:"message"`
	Details interface{} `json:"details"`
}
