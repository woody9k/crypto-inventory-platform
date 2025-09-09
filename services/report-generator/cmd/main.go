// Package main provides the entry point for the report generator service.
// This service handles the generation of various compliance and security reports
// including crypto summaries, compliance status, network topology, and risk assessments.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/democorp/crypto-inventory/services/report-generator/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main initializes and starts the report generator service.
// It sets up HTTP routes, CORS middleware, and graceful shutdown handling.
func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize handlers for report generation and management
	handler := handlers.NewHandler()

	// Initialize Gin router with default middleware
	router := gin.Default()

	// Configure CORS middleware for cross-origin requests
	// This allows the web UI to communicate with the report generator service
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // In production, specify exact origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint for service monitoring
	router.GET("/health", handler.Health)

	// API routes for report generation and management
	api := router.Group("/api/v1")
	{
		// Report generation and management endpoints
		api.POST("/reports/generate", handler.GenerateReport)     // Generate new reports
		api.GET("/reports/templates", handler.GetReportTemplates) // Get available templates
		api.GET("/reports/:id", handler.GetReport)                // Get specific report
		api.GET("/reports", handler.ListReports)                  // List all reports
		api.DELETE("/reports/:id", handler.DeleteReport)          // Delete report
		
		// Quick demo endpoints for immediate data access
		// These provide sample data for demonstration purposes
		api.GET("/reports/demo/crypto-summary", handler.GetCryptoSummary)
		api.GET("/reports/demo/compliance-status", handler.GetComplianceStatus)
		api.GET("/reports/demo/network-topology", handler.GetNetworkTopology)
	}

	// Setup HTTP server with configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Default port for report generator service
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine to allow graceful shutdown handling
	go func() {
		log.Printf("🚀 Report Generator starting on port %s", port)
		log.Printf("📊 Ready to generate reports")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down report-generator...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("report-generator forced to shutdown: %v", err)
	}

	log.Println("report-generator stopped")
}
