package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"saas-admin-service/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ListPlatformUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT id, email, first_name, last_name, role, is_active, email_verified,
			       last_login_at, created_at, updated_at
			FROM platform_users
			ORDER BY created_at DESC
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch platform users"})
			return
		}
		defer rows.Close()

		var users []models.PlatformUser
		for rows.Next() {
			var user models.PlatformUser
			err := rows.Scan(
				&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role,
				&user.IsActive, &user.EmailVerified, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan platform user"})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func GetPlatformUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		var user models.PlatformUser
		query := `
			SELECT id, email, first_name, last_name, role, is_active, email_verified,
			       last_login_at, created_at, updated_at
			FROM platform_users
			WHERE id = $1
		`

		err := db.QueryRow(query, userID).Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role,
			&user.IsActive, &user.EmailVerified, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Platform user not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch platform user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func CreatePlatformUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email     string `json:"email" binding:"required,email"`
			Password  string `json:"password" binding:"required,min=8"`
			FirstName string `json:"first_name" binding:"required"`
			LastName  string `json:"last_name" binding:"required"`
			Role      string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Create platform user
		query := `
			INSERT INTO platform_users (email, password_hash, first_name, last_name, role, is_active, email_verified)
			VALUES ($1, $2, $3, $4, $5, true, false)
			RETURNING id, created_at, updated_at
		`

		var userID string
		var createdAt, updatedAt time.Time
		err = db.QueryRow(query, req.Email, string(hashedPassword), req.FirstName, req.LastName, req.Role).
			Scan(&userID, &createdAt, &updatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create platform user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Platform user created successfully",
			"user_id": userID,
		})
	}
}

func UpdatePlatformUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		var req struct {
			FirstName *string `json:"first_name"`
			LastName  *string `json:"last_name"`
			Role      *string `json:"role"`
			IsActive  *bool   `json:"is_active"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Build dynamic update query
		updates := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.FirstName != nil {
			updates = append(updates, "first_name = $"+strconv.Itoa(argIndex))
			args = append(args, *req.FirstName)
			argIndex++
		}
		if req.LastName != nil {
			updates = append(updates, "last_name = $"+strconv.Itoa(argIndex))
			args = append(args, *req.LastName)
			argIndex++
		}
		if req.Role != nil {
			updates = append(updates, "role = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Role)
			argIndex++
		}
		if req.IsActive != nil {
			updates = append(updates, "is_active = $"+strconv.Itoa(argIndex))
			args = append(args, *req.IsActive)
			argIndex++
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		updates = append(updates, "updated_at = NOW()")
		args = append(args, userID)

		query := "UPDATE platform_users SET " + strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

		_, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update platform user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Platform user updated successfully"})
	}
}

func DeletePlatformUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		_, err := db.Exec("DELETE FROM platform_users WHERE id = $1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete platform user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Platform user deleted successfully"})
	}
}
