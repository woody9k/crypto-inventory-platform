package handlers

import (
	"github.com/gin-gonic/gin"
)

// Handler contains all the handler functions
type Handler struct {
	// Add any dependencies here (database, services, etc.)
}

// NewHandler creates a new handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// Health handles health check requests
func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "sensor-manager",
	})
}
