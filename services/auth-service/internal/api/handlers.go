package api

import (
	"net/http"

	"github.com/democorp/crypto-inventory/services/auth-service/internal/auth"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/config"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandlers contains all authentication-related handlers
type AuthHandlers struct {
	authService *auth.AuthService
	config      *config.Config
}

// NewAuthHandlers creates a new instance of auth handlers
func NewAuthHandlers(authService *auth.AuthService, cfg *config.Config) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		config:      cfg,
	}
}

// Register handles user registration
func (h *AuthHandlers) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		switch err {
		case auth.ErrEmailExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to register user",
				"details": err.Error(),
			})
		}
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login handles user authentication
func (h *AuthHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
		case auth.ErrUserInactive:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User account is inactive",
			})
		case auth.ErrEmailNotVerified:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Email not verified",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to authenticate user",
				"details": err.Error(),
			})
		}
		return
	}

	// Remove password hash from response
	authResponse.User.PasswordHash = ""

	c.JSON(http.StatusOK, authResponse)
}

// Logout handles user logout
func (h *AuthHandlers) Logout(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.authService.Logout(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to logout",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// RefreshToken handles token refresh
func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	authResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err {
		case auth.ErrInvalidToken, auth.ErrExpiredToken:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired refresh token",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to refresh token",
				"details": err.Error(),
			})
		}
		return
	}

	// Remove password hash from response
	authResponse.User.PasswordHash = ""

	c.JSON(http.StatusOK, authResponse)
}

// GetMe returns current user information
func (h *AuthHandlers) GetMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get user",
				"details": err.Error(),
			})
		}
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateMe handles user profile updates
func (h *AuthHandlers) UpdateMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implement user update logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "User update not yet implemented",
		"user_id": userID,
		"request": req,
	})
}

// ChangePassword handles password changes
func (h *AuthHandlers) ChangePassword(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implement password change logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Password change not yet implemented",
		"user_id": userID,
	})
}

// ForgotPassword handles password reset requests
func (h *AuthHandlers) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implement forgot password logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Forgot password not yet implemented",
		"email":   req.Email,
	})
}

// ResetPassword handles password reset with token
func (h *AuthHandlers) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implement password reset logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Password reset not yet implemented",
		"token":   req.Token,
	})
}

// VerifyEmail handles email verification
func (h *AuthHandlers) VerifyEmail(c *gin.Context) {
	// TODO: Implement email verification logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Email verification not yet implemented",
	})
}
