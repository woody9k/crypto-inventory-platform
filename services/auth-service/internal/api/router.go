package api

import (
	"database/sql"
	"net/http"

	"github.com/democorp/crypto-inventory/services/auth-service/internal/config"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter(cfg *config.Config, db *sql.DB, redis *redis.Client) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logging())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORSOrigins
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{
		"Origin", "Content-Length", "Content-Type", "Authorization",
		"X-Requested-With", "X-Request-ID",
	}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck(db, redis))

	// API routes
	v1 := router.Group("/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", loginHandler(cfg, db, redis))
			auth.POST("/logout", middleware.RequireAuth(cfg), logoutHandler(redis))
			auth.POST("/refresh", refreshTokenHandler(cfg, redis))
			auth.GET("/me", middleware.RequireAuth(cfg), getMeHandler(db))
		}

		// SSO routes
		sso := v1.Group("/sso")
		{
			sso.POST("/initiate", initiateSSOHandler())
			sso.POST("/callback", ssoCallbackHandler())
		}

		// User management routes (admin only)
		users := v1.Group("/users")
		users.Use(middleware.RequireAuth(cfg))
		users.Use(middleware.RequireRole("admin"))
		{
			users.GET("", listUsersHandler(db))
			users.POST("", createUserHandler(db))
			users.GET("/:id", getUserHandler(db))
			users.PUT("/:id", updateUserHandler(db))
			users.DELETE("/:id", deleteUserHandler(db))
		}

		// Tenant management routes (admin only)
		tenants := v1.Group("/tenants")
		tenants.Use(middleware.RequireAuth(cfg))
		tenants.Use(middleware.RequireRole("admin"))
		{
			tenants.GET("", listTenantsHandler(db))
			tenants.POST("", createTenantHandler(db))
			tenants.GET("/:id", getTenantHandler(db))
			tenants.PUT("/:id", updateTenantHandler(db))
		}
	}

	return router
}

// healthCheck returns a simple health status
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "auth-service",
		"timestamp": gin.H{},
	})
}

// readinessCheck checks if the service is ready to handle requests
func readinessCheck(db *sql.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "database connection failed",
			})
			return
		}

		// Check Redis connection
		if err := redis.Ping(c.Request.Context()).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "redis connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": "auth-service",
		})
	}
}

// Placeholder handlers - these will be implemented in separate files
func loginHandler(cfg *config.Config, db *sql.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Login endpoint - to be implemented"})
	}
}

func logoutHandler(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Logout endpoint - to be implemented"})
	}
}

func refreshTokenHandler(cfg *config.Config, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Refresh token endpoint - to be implemented"})
	}
}

func getMeHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Get me endpoint - to be implemented"})
	}
}

func initiateSSOHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "SSO initiate endpoint - to be implemented"})
	}
}

func ssoCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "SSO callback endpoint - to be implemented"})
	}
}

func listUsersHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "List users endpoint - to be implemented"})
	}
}

func createUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Create user endpoint - to be implemented"})
	}
}

func getUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Get user endpoint - to be implemented"})
	}
}

func updateUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Update user endpoint - to be implemented"})
	}
}

func deleteUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Delete user endpoint - to be implemented"})
	}
}

func listTenantsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "List tenants endpoint - to be implemented"})
	}
}

func createTenantHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Create tenant endpoint - to be implemented"})
	}
}

func getTenantHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Get tenant endpoint - to be implemented"})
	}
}

func updateTenantHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Update tenant endpoint - to be implemented"})
	}
}
