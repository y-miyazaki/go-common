package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRedis handler
func (h *HTTPHandler) GetRedis(c *gin.Context) {

	// Set key
	err := h.redisRepository.Set(c, "a", 1, 0)
	if err != nil {
		h.Logger.WithError(err).Errorf("can't set redis key")
	}

	// Get key
	value, err := h.redisRepository.Get(c, "a")
	if err != nil {
		h.Logger.WithError(err).Errorf("can't get redis key")
	}
	h.Logger.Infof("value = %s", value)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
