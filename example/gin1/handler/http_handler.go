package handler

import (
	"github.com/y-miyazaki/go-common/pkg/handler"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// HTTPHandler struct.
type HTTPHandler struct {
	*handler.BaseHTTPHandler
}

// NewHTTPHandler returns HTTPHandler struct.
func NewHTTPHandler(logger *logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		BaseHTTPHandler: &handler.BaseHTTPHandler{
			Logger: logger,
		},
	}
}
