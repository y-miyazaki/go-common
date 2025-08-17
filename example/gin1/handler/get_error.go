package handler

import (
	"github.com/y-miyazaki/go-common/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HandleError1 responds with an internal server error for testing purposes.
func (h *HTTPHandler) HandleError1(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, map[string]string{
		"test": "testmessage",
	}, errors.New("error test1"))
}

// HandleError2 responds with an internal server error using structured response.
func (h *HTTPHandler) HandleError2(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, dto.HTTPErrorResponse{
		Message: map[string]string{
			"test": "testmessage",
		},
	}, errors.New("error test2"))
}
