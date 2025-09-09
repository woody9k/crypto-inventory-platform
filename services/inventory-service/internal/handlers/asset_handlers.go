package handlers

import (
	"inventory-service/internal/models"
	"inventory-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService *services.AssetService
}

func NewAssetHandler(assetService *services.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

// GetAssets handles GET /api/v1/assets
func (h *AssetHandler) GetAssets(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	// Parse filters
	var filters models.AssetFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	// Set default pagination
	if filters.Page == 0 {
		filters.Page = 1
	}
	if filters.PageSize == 0 {
		filters.PageSize = 20
	}

	assets, total, err := h.assetService.GetAssets(tenantUUID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assets", "details": err.Error()})
		return
	}

	// Calculate pagination metadata
	totalPages := (total + filters.PageSize - 1) / filters.PageSize
	hasNext := filters.Page < totalPages
	hasPrev := filters.Page > 1

	response := gin.H{
		"assets": assets,
		"pagination": gin.H{
			"page":        filters.Page,
			"page_size":   filters.PageSize,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetAssetByID handles GET /api/v1/assets/:id
func (h *AssetHandler) GetAssetByID(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}

	asset, err := h.assetService.GetAssetByID(tenantUUID, assetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"asset": asset})
}

// GetAssetCrypto handles GET /api/v1/assets/:id/crypto
func (h *AssetHandler) GetAssetCrypto(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset ID"})
		return
	}

	cryptoImpls, err := h.assetService.GetCryptoImplementations(tenantUUID, assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve crypto implementations", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"crypto_implementations": cryptoImpls})
}

// SearchAssets handles GET /api/v1/assets/search
func (h *AssetHandler) SearchAssets(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Parse optional parameters
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	filters := models.AssetFilters{
		Search:   query,
		PageSize: limit,
		Page:     1,
	}

	assets, total, err := h.assetService.GetAssets(tenantUUID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"assets":  assets,
		"total":   total,
		"showing": len(assets),
	})
}

// GetRiskSummary handles GET /api/v1/risk/summary
func (h *AssetHandler) GetRiskSummary(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	summary, err := h.assetService.GetRiskSummary(tenantUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get risk summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"risk_summary": summary})
}

// Health check handler
func (h *AssetHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "inventory-service",
	})
}
