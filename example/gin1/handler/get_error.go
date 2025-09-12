package handler

import (
	"errors"

	"github.com/y-miyazaki/go-common/pkg/dto"

	"github.com/gin-gonic/gin"
)

var (
	// ErrTest1 is an error for testing purposes
	ErrTest1 = errors.New("error test1")
	// ErrTest2 is an error for testing purposes
	ErrTest2 = errors.New("error test2")
)

// HandleError1 responds with an internal server error for testing purposes.
func (h *HTTPHandler) HandleError1(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, map[string]string{
		"test": "testmessage",
	}, ErrTest1)
}

// HandleError2 responds with an internal server error using structured response.
func (h *HTTPHandler) HandleError2(c *gin.Context) {
	h.ResponseStatusInternalServerError(c, dto.HTTPErrorResponse{
		Message: map[string]string{
			"test": "testmessage",
		},
	}, ErrTest2)
}
