// Package handler provides HTTP handler helpers used by HTTP handlers across the project.
package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"

	"go-common/pkg/context"
	"go-common/pkg/logger"

	"github.com/gin-gonic/gin"
)

// BaseHTTPHandler base handler struct.
type BaseHTTPHandler struct {
	Logger *logger.Logger
}

// ResponseCSV responses csv data.
func (h *BaseHTTPHandler) ResponseCSV(c *gin.Context, statusCode int, fileName string, data []byte) {
	// reference receiver to satisfy linters when receiver is unused
	_ = h
	c.Writer.Header().Set("Content-Description", "File Transfer")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)
	c.Data(statusCode, "text/csv", data)
}

// ResponseZIP responses zip data.
func (h *BaseHTTPHandler) ResponseZIP(c *gin.Context, statusCode int, fileName string, mapContentFile map[string]bytes.Buffer) error {
	// reference receiver to satisfy linters when receiver is unused
	_ = h
	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)

	// ensure status code is actually used
	c.Writer.WriteHeader(statusCode)

	zipW := zip.NewWriter(c.Writer)
	defer func() {
		err := zipW.Close()
		if err != nil {
			h.Logger.WithError(err).Errorf("can't close zip.Writer.(%s)", fileName)
		}
	}()
	// copying bytes.Buffer from map values is intentional to avoid changing public API
	// nolint:gocritic
	for key := range mapContentFile {
		contentFile := mapContentFile[key]
		f, err := zipW.Create(key)
		if err != nil {
			return fmt.Errorf("zip create: %w", err)
		}
		_, err = f.Write(contentFile.Bytes())
		if err != nil {
			return fmt.Errorf("zip write: %w", err)
		}
	}
	return nil
}

// ResponseStatusBadRequest returns 400 error.
func (h *BaseHTTPHandler) ResponseStatusBadRequest(c *gin.Context, messages any) {
	_ = h
	c.JSON(http.StatusBadRequest, messages)
}

// ResponseStatusForbidden returns 403 error.
func (h *BaseHTTPHandler) ResponseStatusForbidden(c *gin.Context, messages any) {
	_ = h
	c.JSON(http.StatusForbidden, messages)
}

// ResponseStatusNotFound returns 404 error.
func (h *BaseHTTPHandler) ResponseStatusNotFound(c *gin.Context, messages any) {
	_ = h
	c.JSON(http.StatusNotFound, messages)
}

// ResponseStatusInternalServerError returns 500 error.
func (h *BaseHTTPHandler) ResponseStatusInternalServerError(c *gin.Context, messages any, err error) {
	_ = h
	context.SetGinContextError(c, err)
	context.SetGinContextErrorMessage(c, messages)
	c.JSON(http.StatusInternalServerError, messages)
}
