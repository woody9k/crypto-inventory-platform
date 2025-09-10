package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"saas-admin-service/internal/models"

	"github.com/gin-gonic/gin"
)

func ListTenants(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset := (page - 1) * limit

		query := `
			SELECT t.id, t.name, t.slug, t.domain, st.name as subscription_tier,
			       t.trial_ends_at, t.billing_email, t.payment_status, t.stripe_customer_id,
			       t.sso_enabled, t.created_at, t.updated_at, t.deleted_at,
			       CASE WHEN t.deleted_at IS NULL THEN true ELSE false END as is_active
			FROM tenants t
			LEFT JOIN subscription_tiers st ON t.subscription_tier_id = st.id
			ORDER BY t.created_at DESC
			LIMIT $1 OFFSET $2
		`

		rows, err := db.Query(query, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenants"})
			return
		}
		defer rows.Close()

		var tenants []models.Tenant
		for rows.Next() {
			var tenant models.Tenant
			err := rows.Scan(
				&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.Domain, &tenant.SubscriptionTier,
				&tenant.TrialEndsAt, &tenant.BillingEmail, &tenant.PaymentStatus, &tenant.StripeCustomerID,
				&tenant.SSOEnabled, &tenant.CreatedAt, &tenant.UpdatedAt, &tenant.DeletedAt, &tenant.IsActive,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan tenant"})
				return
			}
			tenants = append(tenants, tenant)
		}

		// Get total count
		var total int
		db.QueryRow("SELECT COUNT(*) FROM tenants").Scan(&total)

		c.JSON(http.StatusOK, gin.H{
			"tenants": tenants,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	}
}

func GetTenant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		var tenant models.Tenant
		query := `
			SELECT t.id, t.name, t.slug, t.domain, st.name as subscription_tier,
			       t.trial_ends_at, t.billing_email, t.payment_status, t.stripe_customer_id,
			       t.sso_enabled, t.created_at, t.updated_at, t.deleted_at,
			       CASE WHEN t.deleted_at IS NULL THEN true ELSE false END as is_active
			FROM tenants t
			LEFT JOIN subscription_tiers st ON t.subscription_tier_id = st.id
			WHERE t.id = $1
		`

		err := db.QueryRow(query, tenantID).Scan(
			&tenant.ID, &tenant.Name, &tenant.Slug, &tenant.Domain, &tenant.SubscriptionTier,
			&tenant.TrialEndsAt, &tenant.BillingEmail, &tenant.PaymentStatus, &tenant.StripeCustomerID,
			&tenant.SSOEnabled, &tenant.CreatedAt, &tenant.UpdatedAt, &tenant.DeletedAt, &tenant.IsActive,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tenant": tenant})
	}
}

func CreateTenant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string `json:"name" binding:"required"`
			Slug             string `json:"slug" binding:"required"`
			Domain           string `json:"domain"`
			SubscriptionTier string `json:"subscription_tier" binding:"required"`
			BillingEmail     string `json:"billing_email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get subscription tier ID
		var subscriptionTierID string
		err := db.QueryRow("SELECT id FROM subscription_tiers WHERE name = $1", req.SubscriptionTier).Scan(&subscriptionTierID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription tier"})
			return
		}

		// Create tenant
		query := `
			INSERT INTO tenants (name, slug, domain, subscription_tier_id, billing_email, trial_ends_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, created_at, updated_at
		`

		var tenantID string
		var createdAt, updatedAt time.Time
		trialEndsAt := time.Now().Add(30 * 24 * time.Hour) // 30 days trial

		err = db.QueryRow(query, req.Name, req.Slug, req.Domain, subscriptionTierID, req.BillingEmail, trialEndsAt).
			Scan(&tenantID, &createdAt, &updatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":   "Tenant created successfully",
			"tenant_id": tenantID,
		})
	}
}

func UpdateTenant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		var req struct {
			Name             *string `json:"name"`
			Domain           *string `json:"domain"`
			SubscriptionTier *string `json:"subscription_tier"`
			BillingEmail     *string `json:"billing_email"`
			PaymentStatus    *string `json:"payment_status"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Build dynamic update query
		updates := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.Name != nil {
			updates = append(updates, "name = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Name)
			argIndex++
		}
		if req.Domain != nil {
			updates = append(updates, "domain = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Domain)
			argIndex++
		}
		if req.BillingEmail != nil {
			updates = append(updates, "billing_email = $"+strconv.Itoa(argIndex))
			args = append(args, *req.BillingEmail)
			argIndex++
		}
		if req.PaymentStatus != nil {
			updates = append(updates, "payment_status = $"+strconv.Itoa(argIndex))
			args = append(args, *req.PaymentStatus)
			argIndex++
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		updates = append(updates, "updated_at = NOW()")
		args = append(args, tenantID)

		query := "UPDATE tenants SET " + strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

		_, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tenant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tenant updated successfully"})
	}
}

func DeleteTenant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		// Soft delete
		_, err := db.Exec("UPDATE tenants SET deleted_at = NOW() WHERE id = $1", tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tenant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
	}
}

func SuspendTenant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		_, err := db.Exec("UPDATE tenants SET payment_status = 'suspended', updated_at = NOW() WHERE id = $1", tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend tenant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tenant suspended successfully"})
	}
}

func ActivateTenant(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		_, err := db.Exec("UPDATE tenants SET payment_status = 'active', updated_at = NOW() WHERE id = $1", tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate tenant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tenant activated successfully"})
	}
}

func GetTenantStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		var stats models.TenantStats
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
				WHERE tenant_id = $1
				GROUP BY tenant_id
			) u ON t.id = u.tenant_id
			LEFT JOIN (
				SELECT 
					tenant_id,
					COUNT(*) as asset_count
				FROM network_assets 
				WHERE tenant_id = $1
				GROUP BY tenant_id
			) a ON t.id = a.tenant_id
			LEFT JOIN (
				SELECT 
					tenant_id,
					COUNT(*) as sensor_count
				FROM sensors 
				WHERE tenant_id = $1
				GROUP BY tenant_id
			) s ON t.id = s.tenant_id
			WHERE t.id = $1
		`

		err := db.QueryRow(query, tenantID).Scan(
			&stats.TenantID, &stats.TenantName, &stats.UserCount, &stats.AssetCount,
			&stats.SensorCount, &stats.LastActivity, &stats.StorageUsed, &stats.APIRequests,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenant stats"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"stats": stats})
	}
}
