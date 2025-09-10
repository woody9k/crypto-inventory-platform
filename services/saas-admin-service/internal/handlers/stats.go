package handlers

import (
	"database/sql"
	"net/http"

	"saas-admin-service/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPlatformStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats models.PlatformStats

		// Get total tenants
		db.QueryRow("SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL").Scan(&stats.TotalTenants)

		// Get active tenants
		db.QueryRow("SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL AND payment_status = 'active'").Scan(&stats.ActiveTenants)

		// Get total users
		db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&stats.TotalUsers)

		// Get total assets
		db.QueryRow("SELECT COUNT(*) FROM network_assets WHERE deleted_at IS NULL").Scan(&stats.TotalAssets)

		// Get total sensors
		db.QueryRow("SELECT COUNT(*) FROM sensors WHERE deleted_at IS NULL").Scan(&stats.TotalSensors)

		c.JSON(http.StatusOK, gin.H{"stats": stats})
	}
}

func GetTenantsStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT 
				t.id as tenant_id,
				t.name as tenant_name,
				COALESCE(user_count, 0) as user_count,
				COALESCE(asset_count, 0) as asset_count,
				COALESCE(sensor_count, 0) as sensor_count,
				COALESCE(last_activity, t.created_at) as last_activity,
				COALESCE(storage_used, 0) as storage_used,
				COALESCE(api_requests, 0) as api_requests
			FROM tenants t
			LEFT JOIN (
				SELECT 
					tenant_id,
					COUNT(*) as user_count,
					MAX(last_login_at) as last_activity
				FROM users 
				GROUP BY tenant_id
			) u ON t.id = u.tenant_id
			LEFT JOIN (
				SELECT 
					tenant_id,
					COUNT(*) as asset_count
				FROM network_assets 
				GROUP BY tenant_id
			) a ON t.id = a.tenant_id
			LEFT JOIN (
				SELECT 
					tenant_id,
					COUNT(*) as sensor_count
				FROM sensors 
				GROUP BY tenant_id
			) s ON t.id = s.tenant_id
			WHERE t.deleted_at IS NULL
			ORDER BY t.created_at DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenants stats"})
			return
		}
		defer rows.Close()

		var stats []models.TenantStats
		for rows.Next() {
			var stat models.TenantStats
			err := rows.Scan(
				&stat.TenantID, &stat.TenantName, &stat.UserCount, &stat.AssetCount,
				&stat.SensorCount, &stat.LastActivity, &stat.StorageUsed, &stat.APIRequests,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan tenant stats"})
				return
			}
			stats = append(stats, stat)
		}

		c.JSON(http.StatusOK, gin.H{"tenants_stats": stats})
	}
}

func GetSystemHealth(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connectivity
		var dbStatus string
		err := db.QueryRow("SELECT 'healthy'").Scan(&dbStatus)
		if err != nil {
			dbStatus = "unhealthy"
		}

		// Get service status (simplified)
		health := gin.H{
			"database":  dbStatus,
			"timestamp": gin.H{},
		}

		status := http.StatusOK
		if dbStatus != "healthy" {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, health)
	}
}

func GetSystemLogs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for system logs
		c.JSON(http.StatusOK, gin.H{
			"logs":    []gin.H{},
			"message": "System logs endpoint - to be implemented",
		})
	}
}
