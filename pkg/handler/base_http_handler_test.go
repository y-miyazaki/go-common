package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/y-miyazaki/go-common/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBaseHTTPHandler_ResponseCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &BaseHTTPHandler{
		Logger: logger.NewLogger(logrus.New()),
	}

	data := []byte("name,age\nJohn,30")
	handler.ResponseCSV(c, http.StatusOK, "test.csv", data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment;filename=test.csv", w.Header().Get("Content-Disposition"))
	assert.Equal(t, data, w.Body.Bytes())
}

func TestBaseHTTPHandler_ResponseZIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &BaseHTTPHandler{
		Logger: logger.NewLogger(logrus.New()),
	}

	mapContentFile := map[string]bytes.Buffer{
		"test1.txt": *bytes.NewBufferString("content1"),
		"test2.txt": *bytes.NewBufferString("content2"),
	}

	err := handler.ResponseZIP(c, http.StatusOK, "test.zip", mapContentFile)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment;filename=test.zip", w.Header().Get("Content-Disposition"))
	assert.True(t, len(w.Body.Bytes()) > 0)
}

func TestBaseHTTPHandler_ResponseStatusBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &BaseHTTPHandler{
		Logger: logger.NewLogger(logrus.New()),
	}

	messages := map[string]string{"error": "bad request"}
	handler.ResponseStatusBadRequest(c, messages)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "bad request")
}

func TestBaseHTTPHandler_ResponseStatusForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &BaseHTTPHandler{
		Logger: logger.NewLogger(logrus.New()),
	}

	messages := map[string]string{"error": "forbidden"}
	handler.ResponseStatusForbidden(c, messages)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "forbidden")
}

func TestBaseHTTPHandler_ResponseStatusNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &BaseHTTPHandler{
		Logger: logger.NewLogger(logrus.New()),
	}

	messages := map[string]string{"error": "not found"}
	handler.ResponseStatusNotFound(c, messages)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestBaseHTTPHandler_ResponseStatusInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &BaseHTTPHandler{
		Logger: logger.NewLogger(logrus.New()),
	}

	messages := map[string]string{"error": "internal server error"}
	err := assert.AnError
	handler.ResponseStatusInternalServerError(c, messages, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
}
