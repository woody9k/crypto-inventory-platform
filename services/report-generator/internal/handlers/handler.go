// Package handlers provides HTTP handlers for the report generator service.
// It includes handlers for report generation, management, and demo data endpoints.
package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Handler provides report generation handlers and manages report state.
// In a production environment, this would interface with a database
// instead of using in-memory storage.
type Handler struct {
	// In-memory storage for demo purposes
	// In production, this would be replaced with database operations
	reports map[string]*Report
}

// Report represents a generated report with metadata and content.
// This structure tracks the report lifecycle from generation to completion.
type Report struct {
	ID          string      `json:"id"`           // Unique report identifier
	Title       string      `json:"title"`        // Human-readable report title
	Type        string      `json:"type"`         // Report type (crypto_summary, compliance_status, etc.)
	Status      string      `json:"status"`       // Current status (generating, completed, failed)
	CreatedAt   time.Time   `json:"created_at"`   // Report creation timestamp
	CompletedAt *time.Time  `json:"completed_at,omitempty"` // Completion timestamp (nil if not completed)
	Data        interface{} `json:"data,omitempty"`         // Report content/data
	DownloadURL string      `json:"download_url,omitempty"` // URL for downloading the report
}

// ReportTemplate represents a predefined report template that users can generate.
// Templates define the structure and parameters for different report types.
type ReportTemplate struct {
	ID          string `json:"id"`          // Unique template identifier
	Name        string `json:"name"`        // Human-readable template name
	Description string `json:"description"` // Template description
	Type        string `json:"type"`        // Template type (summary, compliance, etc.)
	Category    string `json:"category"`    // Template category (crypto, compliance, network, security)
}

// NewHandler creates a new handler instance with initialized storage.
// This function sets up the handler with empty in-memory storage for demo purposes.
func NewHandler() *Handler {
	return &Handler{
		reports: make(map[string]*Report),
	}
}

// Health returns the health status of the report generator service.
// This endpoint is used by load balancers and monitoring systems to check service availability.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy", "service": "report-generator"})
}
