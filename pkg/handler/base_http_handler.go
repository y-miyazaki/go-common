package handler

import (
	"archive/zip"
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/pkg/context"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
)

// BaseHTTPHandler base handler struct.
type BaseHTTPHandler struct {
	Logger *infrastructure.Logger
}

// ResponseCSV responses csv data.
func (h *BaseHTTPHandler) ResponseCSV(c *gin.Context, statusCode int, fileName string, data []byte) {
	c.Writer.Header().Set("Content-Description", "File Transfer")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)
	c.Data(statusCode, "text/csv", data)
}

// ResponseZIP responses zip data.
func (h *BaseHTTPHandler) ResponseZIP(c *gin.Context, statusCode int, fileName string, mapContentFile map[string]bytes.Buffer) error {
	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)

	zipW := zip.NewWriter(c.Writer)
	defer func() {
		err := zipW.Close()
		if err != nil {
			h.Logger.WithError(err).Errorf("can't close zip.Writer.(%s)", fileName)
		}
	}()
	for key, contentFile := range mapContentFile {
		f, err := zipW.Create(key)
		if err != nil {
			return err
		}
		_, err = f.Write(contentFile.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

// ResponseStatusBadRequest returns 400 error.
func (h *BaseHTTPHandler) ResponseStatusBadRequest(c *gin.Context, messages interface{}) {
	c.JSON(http.StatusBadRequest, messages)
}

// ResponseStatusForbidden returns 403 error.
func (h *BaseHTTPHandler) ResponseStatusForbidden(c *gin.Context, messages interface{}) {
	c.JSON(http.StatusForbidden, messages)
}

// ResponseStatusNotFound returns 404 error.
func (h *BaseHTTPHandler) ResponseStatusNotFound(c *gin.Context, messages interface{}) {
	c.JSON(http.StatusNotFound, messages)
}

// ResponseStatusInternalServerError returns 500 error.
func (h *BaseHTTPHandler) ResponseStatusInternalServerError(c *gin.Context, messages interface{}, err error) {
	context.SetGinContextError(c, err)
	context.SetGinContextErrorMessage(c, messages)
	c.JSON(http.StatusInternalServerError, messages)
}
