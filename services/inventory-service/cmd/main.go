package main

import (
	"inventory-service/internal/config"
	"inventory-service/internal/database"
	"inventory-service/internal/handlers"
	"inventory-service/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services
	assetService := services.NewAssetService(db)

	// Initialize handlers
	assetHandler := handlers.NewAssetHandler(assetService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint (no auth required)
	r.GET("/health", assetHandler.Health)

	// API routes with JWT middleware
	api := r.Group("/api/v1")
	api.Use(handlers.JWTMiddleware(cfg))
	{
		// Asset endpoints
		api.GET("/assets", assetHandler.GetAssets)
		api.GET("/assets/search", assetHandler.SearchAssets)
		api.GET("/assets/:id", assetHandler.GetAssetByID)
		api.GET("/assets/:id/crypto", assetHandler.GetAssetCrypto)

		// Risk endpoints
		api.GET("/risk/summary", assetHandler.GetRiskSummary)
	}

	// Start server
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: r,
	}

	log.Printf("🚀 Inventory Service starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("📊 Ready to serve crypto asset inventory")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
