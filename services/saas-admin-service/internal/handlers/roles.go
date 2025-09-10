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

func ListPlatformRoles(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT id, name, display_name, description, is_system_role, created_at, updated_at
			FROM platform_roles
			ORDER BY created_at ASC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch platform roles"})
			return
		}
		defer rows.Close()

		var roles []models.PlatformRole
		for rows.Next() {
			var role models.PlatformRole
			err := rows.Scan(
				&role.ID, &role.Name, &role.DisplayName, &role.Description,
				&role.IsSystemRole, &role.CreatedAt, &role.UpdatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan platform role"})
				return
			}
			roles = append(roles, role)
		}

		c.JSON(http.StatusOK, gin.H{"roles": roles})
	}
}

func GetPlatformRole(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Param("id")

		var role models.PlatformRole
		query := `
			SELECT id, name, display_name, description, is_system_role, created_at, updated_at
			FROM platform_roles
			WHERE id = $1
		`

		err := db.QueryRow(query, roleID).Scan(
			&role.ID, &role.Name, &role.DisplayName, &role.Description,
			&role.IsSystemRole, &role.CreatedAt, &role.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Platform role not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch platform role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"role": role})
	}
}

func CreatePlatformRole(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string `json:"name" binding:"required"`
			DisplayName string `json:"display_name" binding:"required"`
			Description string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `
			INSERT INTO platform_roles (name, display_name, description, is_system_role)
			VALUES ($1, $2, $3, false)
			RETURNING id, created_at, updated_at
		`

		var roleID string
		var createdAt, updatedAt time.Time
		err := db.QueryRow(query, req.Name, req.DisplayName, req.Description).
			Scan(&roleID, &createdAt, &updatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create platform role"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Platform role created successfully",
			"role_id": roleID,
		})
	}
}

func UpdatePlatformRole(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Param("id")

		var req struct {
			DisplayName *string `json:"display_name"`
			Description *string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Build dynamic update query
		updates := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.DisplayName != nil {
			updates = append(updates, "display_name = $"+strconv.Itoa(argIndex))
			args = append(args, *req.DisplayName)
			argIndex++
		}
		if req.Description != nil {
			updates = append(updates, "description = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Description)
			argIndex++
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		updates = append(updates, "updated_at = NOW()")
		args = append(args, roleID)

		query := "UPDATE platform_roles SET " + strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

		_, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update platform role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Platform role updated successfully"})
	}
}

func DeletePlatformRole(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Param("id")

		// Check if it's a system role
		var isSystemRole bool
		err := db.QueryRow("SELECT is_system_role FROM platform_roles WHERE id = $1", roleID).Scan(&isSystemRole)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Platform role not found"})
			return
		}

		if isSystemRole {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete system roles"})
			return
		}

		_, err = db.Exec("DELETE FROM platform_roles WHERE id = $1", roleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete platform role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Platform role deleted successfully"})
	}
}

func ListPlatformPermissions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT id, name, resource, action, description, created_at
			FROM platform_permissions
			ORDER BY resource, action
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch platform permissions"})
			return
		}
		defer rows.Close()

		var permissions []models.PlatformPermission
		for rows.Next() {
			var permission models.PlatformPermission
			err := rows.Scan(
				&permission.ID, &permission.Name, &permission.Resource, &permission.Action,
				&permission.Description, &permission.CreatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan platform permission"})
				return
			}
			permissions = append(permissions, permission)
		}

		c.JSON(http.StatusOK, gin.H{"permissions": permissions})
	}
}

func GetPlatformPermission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionID := c.Param("id")

		var permission models.PlatformPermission
		query := `
			SELECT id, name, resource, action, description, created_at
			FROM platform_permissions
			WHERE id = $1
		`

		err := db.QueryRow(query, permissionID).Scan(
			&permission.ID, &permission.Name, &permission.Resource, &permission.Action,
			&permission.Description, &permission.CreatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Platform permission not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch platform permission"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"permission": permission})
	}
}
