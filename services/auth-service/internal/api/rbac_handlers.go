package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/democorp/crypto-inventory/services/auth-service/internal/models"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/rbac"
)

// RBACHandlers contains handlers for RBAC operations
type RBACHandlers struct {
	rbacService *rbac.RBACService
}

// NewRBACHandlers creates new RBAC handlers
func NewRBACHandlers(rbacService *rbac.RBACService) *RBACHandlers {
	return &RBACHandlers{
		rbacService: rbacService,
	}
}

// =================================================================
// Tenant Role Management Handlers
// =================================================================

// GetTenantRoles gets all roles for a tenant
func (h *RBACHandlers) GetTenantRoles(c *gin.Context) {
	tenantIDStr := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	roles, err := h.rbacService.GetTenantRoles(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get tenant roles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}

// GetTenantPermissions gets all available permissions
func (h *RBACHandlers) GetTenantPermissions(c *gin.Context) {
	permissions, err := h.rbacService.GetTenantPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get permissions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
	})
}

// GetPermissionMatrix gets the permission matrix for a role
func (h *RBACHandlers) GetPermissionMatrix(c *gin.Context) {
	roleIDStr := c.Param("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	matrix, err := h.rbacService.GetPermissionMatrix(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get permission matrix",
		})
		return
	}

	c.JSON(http.StatusOK, matrix)
}

// UpdateRolePermissions updates permissions for a role
func (h *RBACHandlers) UpdateRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	var req struct {
		PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = h.rbacService.UpdateRolePermissions(roleID, req.PermissionIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update role permissions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role permissions updated successfully",
	})
}

// =================================================================
// User Role Assignment Handlers
// =================================================================

// AssignRole assigns a role to a user
func (h *RBACHandlers) AssignRole(c *gin.Context) {
	var req models.RoleAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get assigned_by from context (current user)
	assignedByStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	assignedBy, err := uuid.Parse(assignedByStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	req.AssignedBy = assignedBy

	err = h.rbacService.AssignRole(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role assigned successfully",
	})
}

// RemoveRole removes a role from a user
func (h *RBACHandlers) RemoveRole(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	tenantIDStr := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	roleIDStr := c.Param("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role ID",
		})
		return
	}

	err = h.rbacService.RemoveRole(userID, tenantID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role removed successfully",
	})
}

// GetUserRoles gets all roles for a user in a tenant
func (h *RBACHandlers) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	tenantIDStr := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	roles, err := h.rbacService.GetUserRoles(userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user roles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}

// GetUserPermissions gets all permissions for a user in a tenant
func (h *RBACHandlers) GetUserPermissions(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	tenantIDStr := c.Param("tenantId")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user permissions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
	})
}

// =================================================================
// Platform Administration Handlers
// =================================================================

// GetPlatformUsers gets all platform users
func (h *RBACHandlers) GetPlatformUsers(c *gin.Context) {
	users, err := h.rbacService.GetPlatformUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get platform users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetPlatformRoles gets all platform roles
func (h *RBACHandlers) GetPlatformRoles(c *gin.Context) {
	roles, err := h.rbacService.GetPlatformRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get platform roles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}

// =================================================================
// Permission Check Handlers
// =================================================================

// CheckPermission checks if a user has a specific permission
func (h *RBACHandlers) CheckPermission(c *gin.Context) {
	var req models.PermissionCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	hasPermission, err := h.rbacService.CheckPermission(req.UserID, req.TenantID, req.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check permission",
		})
		return
	}

	// Log the permission check
	h.rbacService.LogPermissionCheck(&req, hasPermission)

	c.JSON(http.StatusOK, gin.H{
		"granted": hasPermission,
	})
}

// =================================================================
// Audit and Monitoring Handlers
// =================================================================

// GetAuditLogs gets permission audit logs
func (h *RBACHandlers) GetAuditLogs(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	userIDStr := c.Query("user_id")
	tenantIDStr := c.Query("tenant_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT id, user_id, tenant_id, action, resource_type, resource_id, 
		       permission_required, permission_granted, ip_address, user_agent, created_at
		FROM permission_audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			query += " AND user_id = $" + strconv.Itoa(argIndex)
			args = append(args, userID)
			argIndex++
		}
	}

	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err == nil {
			query += " AND tenant_id = $" + strconv.Itoa(argIndex)
			args = append(args, tenantID)
			argIndex++
		}
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := h.rbacService.GetDB().Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get audit logs",
		})
		return
	}
	defer rows.Close()

	var logs []models.PermissionAuditLog
	for rows.Next() {
		var log models.PermissionAuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.TenantID, &log.Action, &log.ResourceType,
			&log.ResourceID, &log.PermissionRequired, &log.PermissionGranted,
			&log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan audit log",
			})
			return
		}
		logs = append(logs, log)
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"page":  page,
		"limit": limit,
	})
}
