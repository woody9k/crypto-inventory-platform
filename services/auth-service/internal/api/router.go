package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/democorp/crypto-inventory/services/auth-service/internal/auth"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/config"
	"github.com/democorp/crypto-inventory/services/auth-service/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter(cfg *config.Config, db *sql.DB, redis *redis.Client) *gin.Engine {
	router := gin.New()

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry, 7*24*time.Hour) // 7 days refresh expiry

	// Initialize auth service
	authService := auth.NewAuthService(db, redis, jwtService)

	// Initialize handlers
	authHandlers := NewAuthHandlers(authService, cfg)

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

	// API routes with versioning
	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			// Basic authentication
			auth.POST("/register", authHandlers.Register)
			auth.POST("/verify-email", authHandlers.VerifyEmail)
			auth.POST("/login", authHandlers.Login)
			auth.POST("/logout", middleware.RequireAuth(cfg, jwtService), authHandlers.Logout)
			auth.POST("/refresh", authHandlers.RefreshToken)
			auth.POST("/forgot-password", authHandlers.ForgotPassword)
			auth.POST("/reset-password", authHandlers.ResetPassword)

			// Current user management
			auth.GET("/me", middleware.RequireAuth(cfg, jwtService), authHandlers.GetMe)
			auth.PUT("/me", middleware.RequireAuth(cfg, jwtService), authHandlers.UpdateMe)
			auth.POST("/change-password", middleware.RequireAuth(cfg, jwtService), authHandlers.ChangePassword)

			// Flexible authentication flow (frontend-agnostic)
			auth.POST("/initiate", authInitiateHandler(db))
			auth.GET("/methods", authMethodsHandler(db))
			auth.POST("/authenticate", authenticateHandler(cfg, db, redis))
			auth.POST("/complete", authCompleteHandler(cfg, db, redis))
		}

		// SSO routes (nested under auth)
		authSSO := auth.Group("/sso")
		{
			authSSO.GET("/:provider/authorize", ssoAuthorizeHandler(cfg, db))
			authSSO.GET("/:provider/callback", ssoCallbackHandler(cfg, db, redis))
			authSSO.POST("/link", middleware.RequireAuth(cfg, jwtService), ssoLinkHandler(db))
			authSSO.DELETE("/unlink", middleware.RequireAuth(cfg, jwtService), ssoUnlinkHandler(db))
			authSSO.GET("/providers", ssoProvidersHandler(db))
		}

		// User management routes (admin only)
		users := v1.Group("/users")
		users.Use(middleware.RequireAuth(cfg, jwtService))
		users.Use(middleware.RequireRole("admin"))
		{
			users.GET("", listUsersHandler(db))
			users.POST("", createUserHandler(db))
			users.GET("/:id", getUserHandler(db))
			users.PUT("/:id", updateUserHandler(db))
			users.DELETE("/:id", deleteUserHandler(db))
		}

		// Current tenant management routes
		tenant := v1.Group("/tenant")
		tenant.Use(middleware.RequireAuth(cfg, jwtService))
		{
			tenant.GET("", getCurrentTenantHandler(db))
			tenant.PUT("", middleware.RequireRole("admin"), updateCurrentTenantHandler(db))
			tenant.GET("/usage", getTenantUsageHandler(db))
			tenant.GET("/billing", getTenantBillingHandler(db))
			tenant.POST("/upgrade", middleware.RequireRole("admin"), upgradeTenantHandler(db))
		}

		// Frontend flexibility: UI configuration routes
		ui := v1.Group("/ui")
		ui.Use(middleware.RequireAuth(cfg, jwtService))
		{
			ui.GET("/config", getUIConfigHandler(db))
			ui.GET("/config/tenant", getTenantUIConfigHandler(db))
			ui.PUT("/config/tenant", middleware.RequireRole("admin"), updateTenantUIConfigHandler(db))
			ui.GET("/themes", getUIThemesHandler(db))
			ui.PUT("/branding", middleware.RequireRole("admin"), updateBrandingHandler(db))
		}

		// SSO configuration routes (admin only)
		ssoConfig := v1.Group("/tenant/sso")
		ssoConfig.Use(middleware.RequireAuth(cfg, jwtService))
		ssoConfig.Use(middleware.RequireRole("admin"))
		{
			ssoConfig.GET("/providers", listSSOProvidersHandler(db))
			ssoConfig.POST("/providers", createSSOProviderHandler(db))
			ssoConfig.PUT("/providers/:id", updateSSOProviderHandler(db))
			ssoConfig.DELETE("/providers/:id", deleteSSOProviderHandler(db))
			ssoConfig.POST("/providers/:id/test", testSSOProviderHandler(cfg, db))
		}

		// Billing and subscription routes
		billing := v1.Group("/billing")
		billing.Use(middleware.RequireAuth(cfg, jwtService))
		{
			billing.GET("/tiers", getSubscriptionTiersHandler(db))
			billing.GET("/usage/current", getCurrentUsageHandler(db))
			billing.GET("/usage/history", getUsageHistoryHandler(db))
			billing.POST("/check-limits", checkLimitsHandler(db))
		}

		// Feature availability routes
		features := v1.Group("/features")
		features.Use(middleware.RequireAuth(cfg, jwtService))
		{
			features.GET("/availability", getFeatureAvailabilityHandler(db))
		}

		// Frontend flexibility: Workflow management routes
		workflows := v1.Group("/workflows")
		workflows.Use(middleware.RequireAuth(cfg, jwtService))
		{
			workflows.GET("/onboarding", getOnboardingWorkflowHandler(db))
			workflows.POST("/onboarding/:step", completeOnboardingStepHandler(db))
			workflows.GET("/onboarding/progress", getOnboardingProgressHandler(db))
			workflows.POST("/onboarding/skip/:step", skipOnboardingStepHandler(db))
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
