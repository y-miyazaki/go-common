package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// GetError500 handler
func (h *HTTPHandler) GetError500(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, "error...", errors.New("500 error test."))
}
