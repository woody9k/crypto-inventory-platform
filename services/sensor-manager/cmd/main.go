package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/democorp/crypto-inventory/services/sensor-manager/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize handlers
	handler := handlers.NewHandler()

	// Initialize router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", handler.Health)

	// API routes
	api := router.Group("/api/v1")
	{
		// Registration management
		api.POST("/sensors/pending", handler.CreatePendingSensor)
		api.GET("/sensors/pending", handler.GetPendingSensors)
		api.DELETE("/sensors/pending/:key", handler.DeletePendingSensor)

		// Admin settings
		api.GET("/admin/settings", handler.GetAdminSettings)
		api.PUT("/admin/settings", handler.UpdateAdminSettings)

		// Sensor registration
		api.POST("/sensors/register", handler.RegisterSensor)

		// Sensor-specific routes
		sensors := api.Group("/sensors/:sensor_id")
		{
			// Outbound-only communication endpoints
			sensors.POST("/heartbeat", handler.Heartbeat)
			sensors.GET("/commands", handler.PollCommands)
			sensors.POST("/commands/:command_id/ack", handler.AcknowledgeCommand)
			sensors.GET("/webhook-config", handler.GetWebhookConfig)

			// Discovery submission
			sensors.POST("/discoveries", handler.SubmitDiscoveries)

			// Air-gapped export
			sensors.POST("/exports", handler.SubmitAirGappedExport)

			// Legacy endpoints (for backward compatibility)
			sensors.POST("/health", handler.ReportHealth)
			sensors.GET("/config", handler.GetSensorConfig)
		}
	}

	// Setup server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server
	go func() {
		log.Printf("🚀 Sensor Manager starting on port %s", port)
		log.Printf("📡 Ready to manage network sensors")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down sensor-manager...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("sensor-manager forced to shutdown: %v", err)
	}

	log.Println("sensor-manager stopped")
}
