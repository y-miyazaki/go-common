// Package context provides Gin context utility functions.
package context

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// for gin.Context Key
const (
	contextKeyError        string = "error"
	contextKeyErrorMessage string = "errormessage"
)

var (
	// ErrCannotGetError indicates that the error value in context cannot be retrieved
	ErrCannotGetError = errors.New("can't get error")
	// ErrCannotGetMessage indicates that the message value in context cannot be retrieved
	ErrCannotGetMessage = errors.New("can't get message")
)

// GetGinContextError gets error.
func GetGinContextError(c *gin.Context) (err, err2 error) {
	if tmp, exists := c.Get(contextKeyError); exists {
		if ginErr, ok := tmp.(error); ok {
			return ginErr, nil
		}
		return nil, ErrCannotGetError
	}
	return nil, nil
}

// GetGinContextErrorMessage gets error message.
func GetGinContextErrorMessage(c *gin.Context) (string, error) {
	if tmp, exists := c.Get(contextKeyErrorMessage); exists {
		if message, ok := tmp.(string); ok {
			return message, nil
		}
		return "", ErrCannotGetMessage
	}
	return "", nil
}

// SetGinContextError sets error.
func SetGinContextError(c *gin.Context, err error) {
	c.Set(contextKeyError, err)
}

// SetGinContextErrorMessage sets the error message in Gin context.
func SetGinContextErrorMessage(c *gin.Context, message any) {
	c.Set(contextKeyErrorMessage, message)
}
