// Package api provides the HTTP server and routing for the SaaS Admin Service.
// It handles all platform administration endpoints including tenant management,
// user management, statistics, and system monitoring.
package api

import (
	"database/sql"
	"log"
	"net/http"

	"saas-admin-service/internal/config"
	"saas-admin-service/internal/handlers"
	"saas-admin-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server instance with configuration,
// database connection, and router setup.
type Server struct {
	config *config.Config // Service configuration (port, database URL, JWT secret)
	db     *sql.DB        // Database connection for platform data access
	router *gin.Engine    // Gin router with all endpoints and middleware
}

// NewServer creates and initializes a new HTTP server instance.
// It sets up all routes, middleware, and handlers for the SaaS admin service.
//
// Parameters:
//   - cfg: Service configuration including port, database URL, and JWT secret
//   - db: Database connection for accessing platform and tenant data
//
// Returns:
//   - *Server: Configured server instance ready to start
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	server := &Server{
		config: cfg,
		db:     db,
	}

	// Set up all routes, middleware, and handlers
	server.setupRouter()
	return server
}

// setupRouter configures all HTTP routes, middleware, and handlers for the SaaS admin service.
// It sets up authentication, authorization, CORS, and all API endpoints.
func (s *Server) setupRouter() {
	// Set Gin mode based on environment
	if s.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router with basic middleware
	s.router = gin.New()
	s.router.Use(gin.Logger())   // Request logging
	s.router.Use(gin.Recovery()) // Panic recovery

	// CORS middleware for cross-origin requests
	// This allows the SaaS admin UI to make requests from different origins
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint for service monitoring
	// Used by load balancers and monitoring systems to check service health
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":   "saas-admin-service",
			"status":    "healthy",
			"timestamp": gin.H{},
		})
	})

	// API routes - All endpoints are prefixed with /api/v1
	api := s.router.Group("/api/v1")
	{
		// Authentication routes - Handle platform admin login and token refresh
		// These endpoints are shared with the tenant auth service but use platform user tables
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login(s.db, s.config.JWTSecret))          // Platform admin login
			auth.POST("/refresh", handlers.RefreshToken(s.db, s.config.JWTSecret)) // Token refresh
		}

		// Protected routes - All admin endpoints require authentication and platform admin role
		// These routes are protected by JWT authentication and platform admin authorization
		protected := api.Group("/admin")
		protected.Use(middleware.AuthMiddleware(s.config.JWTSecret)) // JWT authentication
		protected.Use(middleware.PlatformAdminMiddleware())          // Platform admin role check
		{
			// Tenant management endpoints - Full CRUD operations for tenant organizations
			// These endpoints allow platform admins to manage all tenants in the system
			tenants := protected.Group("/tenants")
			{
				// Read-only: super_admin, platform_admin, support_admin
				readTenants := tenants.Group("")
				readTenants.Use(middleware.Authorize("super_admin", "platform_admin", "support_admin"))
				{
					readTenants.GET("", handlers.ListTenants(s.db))                  // List all tenants with pagination
					readTenants.GET("/:id", handlers.GetTenant(s.db))                // Get specific tenant details
					readTenants.GET("/:id/stats", handlers.GetTenantStats(s.db))     // Get tenant statistics
					readTenants.GET("/:id/billing", handlers.GetTenantBilling(s.db)) // Billing overview for a tenant
				}

				// Create/Update/Suspend/Activate: super_admin, platform_admin
				manageTenants := tenants.Group("")
				manageTenants.Use(middleware.Authorize("super_admin", "platform_admin"))
				{
					manageTenants.POST("", handlers.CreateTenant(s.db))
					manageTenants.PUT("/:id", handlers.UpdateTenant(s.db))
					manageTenants.POST("/:id/suspend", handlers.SuspendTenant(s.db))
					manageTenants.POST("/:id/activate", handlers.ActivateTenant(s.db))
					manageTenants.PUT("/:id/billing", handlers.UpdateTenantBilling(s.db)) // Update plan or cancel/reactivate
				}

				// Delete: super_admin only
				deleteTenants := tenants.Group("")
				deleteTenants.Use(middleware.Authorize("super_admin"))
				{
					deleteTenants.DELETE("/:id", handlers.DeleteTenant(s.db))
				}
			}

			// Platform user management endpoints - Manage SaaS administrators
			// These endpoints allow management of platform-level users (not tenant users)
			users := protected.Group("/users")
			{
				// Read: super_admin, platform_admin, support_admin
				readUsers := users.Group("")
				readUsers.Use(middleware.Authorize("super_admin", "platform_admin", "support_admin"))
				{
					readUsers.GET("", handlers.ListPlatformUsers(s.db))
					readUsers.GET("/:id", handlers.GetPlatformUser(s.db))
				}

				// Create/Update: super_admin, platform_admin
				manageUsers := users.Group("")
				manageUsers.Use(middleware.Authorize("super_admin", "platform_admin"))
				{
					manageUsers.POST("", handlers.CreatePlatformUser(s.db))
					manageUsers.PUT("/:id", handlers.UpdatePlatformUser(s.db))
				}

				// Delete: super_admin only
				deleteUsers := users.Group("")
				deleteUsers.Use(middleware.Authorize("super_admin"))
				{
					deleteUsers.DELETE("/:id", handlers.DeletePlatformUser(s.db))
				}
			}

			// Billing endpoints (admin-wide)
			billing := protected.Group("/billing")
			{
				// Invoices listing: super_admin, platform_admin (read-only)
				listInvoices := billing.Group("")
				listInvoices.Use(middleware.Authorize("super_admin", "platform_admin"))
				{
					listInvoices.GET("/invoices", handlers.ListInvoices(s.db))
				}

				// Credits issuance: super_admin only
				credits := billing.Group("")
				credits.Use(middleware.Authorize("super_admin"))
				{
					credits.POST("/credits", func(c *gin.Context) {
						c.JSON(http.StatusAccepted, gin.H{"message": "Credit issuance scheduled"})
					})
				}

				// Provider-agnostic webhook (no auth, separate route group)
			}

			// Platform roles and permissions endpoints - RBAC management
			// These endpoints manage platform-level roles and permissions
			roles := protected.Group("/roles")
			{
				// Read: super_admin, platform_admin, support_admin
				readRoles := roles.Group("")
				readRoles.Use(middleware.Authorize("super_admin", "platform_admin", "support_admin"))
				{
					readRoles.GET("", handlers.ListPlatformRoles(s.db))
					readRoles.GET("/:id", handlers.GetPlatformRole(s.db))
				}

				// Create/Update/Delete: super_admin only
				writeRoles := roles.Group("")
				writeRoles.Use(middleware.Authorize("super_admin"))
				{
					writeRoles.POST("", handlers.CreatePlatformRole(s.db))
					writeRoles.PUT("/:id", handlers.UpdatePlatformRole(s.db))
					writeRoles.DELETE("/:id", handlers.DeletePlatformRole(s.db))
				}
			}

			permissions := protected.Group("/permissions")
			{
				// Read: super_admin, platform_admin, support_admin
				permissions.Use(middleware.Authorize("super_admin", "platform_admin", "support_admin"))
				permissions.GET("", handlers.ListPlatformPermissions(s.db))   // List all platform permissions
				permissions.GET("/:id", handlers.GetPlatformPermission(s.db)) // Get specific platform permission
			}

			// Platform statistics endpoints - System-wide metrics and analytics
			// These endpoints provide insights into platform usage and performance
			stats := protected.Group("/stats")
			{
				stats.Use(middleware.Authorize("super_admin", "platform_admin", "support_admin"))
				stats.GET("/platform", handlers.GetPlatformStats(s.db)) // Platform-wide statistics
				stats.GET("/tenants", handlers.GetTenantsStats(s.db))   // All tenants statistics
			}

			// System monitoring endpoints - Health checks and system monitoring
			// These endpoints provide system health and monitoring capabilities
			monitoring := protected.Group("/monitoring")
			{
				monitoring.Use(middleware.Authorize("super_admin", "platform_admin", "support_admin"))
				monitoring.GET("/health", handlers.GetSystemHealth(s.db)) // System health check
				monitoring.GET("/logs", handlers.GetSystemLogs(s.db))     // System logs (future)
			}
		}

		// Provider-agnostic webhook (public path; provider handles verification)
		billingWebhook := api.Group("/admin/billing/webhook")
		{
			billingWebhook.POST(":provider", handlers.BillingWebhook(s.db))
		}
	}
}

// Start begins listening for HTTP requests on the configured port.
// This method blocks until the server is stopped or encounters an error.
//
// Returns:
//   - error: Any error that occurs while starting the server
func (s *Server) Start() error {
	log.Printf("🚀 SaaS Admin Service listening on port %s", s.config.Port)
	return s.router.Run(":" + s.config.Port)
}
