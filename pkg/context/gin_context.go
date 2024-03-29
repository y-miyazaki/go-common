package context

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// for gin.Context Key
const (
	contextKeyError        string = "error"
	contextKeyErrorMessage string = "errormessage"
)

// GetGinContextError gets error.
func GetGinContextError(c *gin.Context) (err, err2 error) {
	if tmp, exists := c.Get(contextKeyError); exists {
		if err, ok := tmp.(error); ok {
			return err, nil
		}
		return nil, errors.New("can't get error")
	}
	return nil, nil
}

// GetGinContextErrorMessage gets error message.
func GetGinContextErrorMessage(c *gin.Context) (string, error) {
	if tmp, exists := c.Get(contextKeyErrorMessage); exists {
		if message, ok := tmp.(string); ok {
			return message, nil
		}
		return "", errors.New("can't get message")
	}
	return "", nil
}

// SetGinContextError sets error.
func SetGinContextError(c *gin.Context, err error) {
	c.Set(contextKeyError, err)
}

// SetGinContextErrorMessage sets error message.
func SetGinContextErrorMessage(c *gin.Context, message interface{}) {
	c.Set(contextKeyErrorMessage, message)
}
