package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (handler *Handler) GetReady(c *gin.Context) {
	if handler.repo.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "database down"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
