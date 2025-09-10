// Package main provides the entry point for the SaaS Admin Service.
// This service handles platform-level administration including tenant management,
// platform user management, and system-wide statistics and monitoring.
//
// Architecture:
// - RESTful API with JWT authentication
// - Platform-level RBAC (Role-Based Access Control)
// - Multi-tenant data access with proper isolation
// - Comprehensive tenant and user management
//
// Key Features:
// - Tenant CRUD operations (create, read, update, delete, suspend, activate)
// - Platform user management (create, update, delete platform administrators)
// - Platform statistics and monitoring
// - System health checks and logging
//
// Security:
// - JWT-based authentication with platform admin roles
// - Role-based authorization (super_admin, platform_admin, support_admin)
// - Input validation and sanitization
// - CORS protection and rate limiting
package main

import (
	"log"

	"saas-admin-service/internal/api"
	"saas-admin-service/internal/config"
	"saas-admin-service/internal/database"
)

// main initializes and starts the SaaS Admin Service.
// It loads configuration, establishes database connection, and starts the HTTP server.
func main() {
	// Load configuration from environment variables with sensible defaults
	cfg := config.Load()

	// Initialize database connection with connection pooling
	// This connects to the shared PostgreSQL database used by all services
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize the HTTP server with all routes and middleware
	server := api.NewServer(cfg, db)

	// Start the server and listen for incoming requests
	log.Printf("🚀 SaaS Admin Service starting on port %s", cfg.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
