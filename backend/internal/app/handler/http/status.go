package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Bounds the /ready database check so an exhausted pool fails the probe instead
// of hanging it.
const readinessCheckTimeout = 2 * time.Second

func (handler *Handler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (handler *Handler) GetReady(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), readinessCheckTimeout)
	defer cancel()

	if err := handler.repo.Ping(ctx); err != nil {
		// Log on the request context: ctx is already cancelled if we timed out.
		slog.WarnContext(c.Request.Context(), "Readiness check failed", "error", err)

		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "database down"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
