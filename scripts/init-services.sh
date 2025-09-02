#!/bin/bash

# Initialize all Go services with basic structure
SERVICES=("inventory-service" "compliance-engine" "report-generator" "sensor-manager")

for service in "${SERVICES[@]}"; do
    echo "Initializing $service..."
    
    cd "services/$service"
    
    # Initialize Go module
    go mod init "github.com/democorp/crypto-inventory/services/$service"
    
    # Create directory structure
    mkdir -p cmd internal/{api,config,database,middleware,models} pkg
    
    # Create basic main.go
    cat > cmd/main.go << EOF
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()
	
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "$service",
		})
	})

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
		log.Printf("$service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down $service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("$service forced to shutdown: %v", err)
	}

	log.Println("$service stopped")
}
EOF

    # Create basic Dockerfile
    cat > Dockerfile.dev << EOF
FROM golang:1.21-alpine AS base

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=base /app/main .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
EOF

    # Add basic dependencies
    go get github.com/gin-gonic/gin
    
    cd ../..
    
    echo "$service initialized!"
done

echo "All services initialized successfully!"
