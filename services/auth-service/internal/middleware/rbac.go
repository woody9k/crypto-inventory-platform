package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/democorp/crypto-inventory/services/auth-service/internal/models"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/rbac"
)

// RequirePermission creates middleware that requires a specific permission
func RequirePermission(rbacService *rbac.RBACService, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found",
			})
			c.Abort()
			return
		}

		tenantIDStr, exists := c.Get("tenantID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant ID not found",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid tenant ID",
			})
			c.Abort()
			return
		}

		// Check permission
		hasPermission, err := rbacService.CheckPermission(userID, tenantID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check permission",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient permissions",
				"required_permission": permission,
			})
			c.Abort()
			return
		}

		// Log permission check for audit
		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		rbacService.LogPermissionCheck(&models.PermissionCheckRequest{
			UserID:     userID,
			TenantID:   tenantID,
			Permission: permission,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
		}, true)

		c.Next()
	}
}

// RequireAnyPermission creates middleware that requires any of the specified permissions
func RequireAnyPermission(rbacService *rbac.RBACService, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found",
			})
			c.Abort()
			return
		}

		tenantIDStr, exists := c.Get("tenantID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant ID not found",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid tenant ID",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		hasAnyPermission := false
		for _, permission := range permissions {
			hasPermission, err := rbacService.CheckPermission(userID, tenantID, permission)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to check permission",
				})
				c.Abort()
				return
			}
			if hasPermission {
				hasAnyPermission = true
				break
			}
		}

		if !hasAnyPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":                "Insufficient permissions",
				"required_permissions": permissions,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePlatformPermission creates middleware that requires a platform permission
func RequirePlatformPermission(rbacService *rbac.RBACService, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		// Check platform permission
		hasPermission, err := rbacService.CheckPlatformPermission(userID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check platform permission",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient platform permissions",
				"required_permission": permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireResourcePermission creates middleware that requires permission for a specific resource
func RequireResourcePermission(rbacService *rbac.RBACService, resourceType, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User ID not found",
			})
			c.Abort()
			return
		}

		tenantIDStr, exists := c.Get("tenantID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Tenant ID not found",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid tenant ID",
			})
			c.Abort()
			return
		}

		// Construct permission name
		permission := strings.ToLower(resourceType) + "." + strings.ToLower(action)

		// Check permission
		hasPermission, err := rbacService.CheckPermission(userID, tenantID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check permission",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient permissions",
				"required_permission": permission,
				"resource_type":       resourceType,
				"action":              action,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireScope creates middleware that requires a specific scope (tenant, own, team)
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware can be used in combination with permission middleware
		// to enforce scope-based access control
		c.Set("required_scope", scope)
		c.Next()
	}
}

// PlatformAdminOnly creates middleware that requires platform admin access
func PlatformAdminOnly(rbacService *rbac.RBACService) gin.HandlerFunc {
	return RequirePlatformPermission(rbacService, "platform.settings")
}

// TenantAdminOnly creates middleware that requires tenant admin access
func TenantAdminOnly(rbacService *rbac.RBACService) gin.HandlerFunc {
	return RequireAnyPermission(rbacService, "users.manage", "settings.manage")
}

// SuperAdminOnly creates middleware that requires super admin access
func SuperAdminOnly(rbacService *rbac.RBACService) gin.HandlerFunc {
	return RequirePlatformPermission(rbacService, "tenants.manage")
}
