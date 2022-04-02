package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/y-miyazaki/go-common/pkg/dto"
)

type test struct {
	message string
}

// GetError1 handler
func (h *HTTPHandler) GetError1(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, map[string]string{
		"test": "testmessage",
	}, errors.New("error test1."))
}

// GetError2 handler
func (h *HTTPHandler) GetError2(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, dto.HTTPErrorResponse{
		Message: map[string]string{
			"test": "testmessage",
		},
	}, errors.New("error test2."))
}
