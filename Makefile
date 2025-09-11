# Crypto Inventory Platform - Makefile

.PHONY: help build-services build-sensor build-frontend build-all test test-unit test-integration test-e2e start stop restart logs clean install-deps

# Default target
help: ## Show this help message
	@echo "Crypto Inventory Platform - Development Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development Environment
start: ## Start all services with docker-compose
	docker-compose up -d

stop: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs from all services
	docker-compose logs -f

clean: ## Clean up containers, volumes, and images
	docker-compose down -v --remove-orphans
	docker system prune -f

# Build Commands
build-services: ## Build all Go backend services
	@echo "Building Go services..."
	cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go
	cd services/inventory-service && go build -o ../../bin/inventory-service ./cmd/main.go
	cd services/compliance-engine && go build -o ../../bin/compliance-engine ./cmd/main.go
	cd services/report-generator && go build -o ../../bin/report-generator ./cmd/main.go
	cd services/sensor-manager && go build -o ../../bin/sensor-manager ./cmd/main.go
	@echo "Go services built successfully!"

build-sensor: ## Build cross-platform network sensor
	@echo "Building network sensor..."
	cd sensor && go build -o ../bin/crypto-sensor ./cmd/main.go
	@echo "Building cross-platform binaries..."
	cd sensor && GOOS=windows GOARCH=amd64 go build -o ../bin/crypto-sensor-windows-amd64.exe ./cmd/main.go
	cd sensor && GOOS=linux GOARCH=amd64 go build -o ../bin/crypto-sensor-linux-amd64 ./cmd/main.go
	cd sensor && GOOS=darwin GOARCH=amd64 go build -o ../bin/crypto-sensor-darwin-amd64 ./cmd/main.go
	cd sensor && GOOS=linux GOARCH=arm64 go build -o ../bin/crypto-sensor-linux-arm64 ./cmd/main.go
	@echo "Network sensor built successfully!"

build-frontend: ## Build React frontend
	@echo "Building frontend..."
	cd web-ui && npm ci && npm run build
	@echo "Frontend built successfully!"

build-ai-service: ## Build AI analysis service
	@echo "Building AI service..."
	cd services/ai-analysis-service && pip install -r requirements.txt
	@echo "AI service dependencies installed!"

build-all: build-services build-sensor build-frontend build-ai-service ## Build all components

# Test Commands
test-unit: ## Run unit tests for all services
	@echo "Running unit tests..."
	cd services/auth-service && go test ./...
	cd services/inventory-service && go test ./...
	cd services/compliance-engine && go test ./...
	cd services/report-generator && go test ./...
	cd services/sensor-manager && go test ./...
	cd sensor && go test ./...
	@echo "Unit tests completed!"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	cd tests/integration && go test ./...
	@echo "Integration tests completed!"

test-e2e: ## Run end-to-end tests
	@echo "Running E2E tests..."
	cd tests/e2e && npm test
	@echo "E2E tests completed!"

test-load: ## Run load tests
	@echo "Running load tests..."
	cd tests && k6 run load-test.js
	@echo "Load tests completed!"

test: test-unit test-integration ## Run unit and integration tests

# Database Commands
db-migrate: ## Run database migrations
	# Apply migrations script as mounted in docker-compose (03-migrations.sql)
	docker-compose exec postgres psql -U crypto_user -d crypto_inventory -f /docker-entrypoint-initdb.d/03-migrations.sql

db-seed: ## Seed database with test data
	# Apply seed script as mounted in docker-compose (04-seed.sql)
	docker-compose exec postgres psql -U crypto_user -d crypto_inventory -f /docker-entrypoint-initdb.d/04-seed.sql

db-reset: ## Reset database (WARNING: destroys all data)
	docker-compose down -v
	# Start infrastructure first (align with Startup Guide)
	docker-compose up -d postgres redis influxdb nats
	sleep 10
	$(MAKE) db-migrate
	$(MAKE) db-seed

# Infrastructure convenience targets (align with Startup Guide)
infra-up: ## Start core infrastructure services (postgres, redis, influxdb, nats)
	docker-compose up -d postgres redis influxdb nats

infra-down: ## Stop core infrastructure services
	docker-compose stop postgres redis influxdb nats

# Development Setup
install-deps: ## Install development dependencies
	@echo "Installing Go dependencies..."
	cd services/auth-service && go mod tidy
	cd services/inventory-service && go mod tidy
	cd services/compliance-engine && go mod tidy
	cd services/report-generator && go mod tidy
	cd services/sensor-manager && go mod tidy
	cd sensor && go mod tidy
	@echo "Installing frontend dependencies..."
	cd web-ui && npm install
	@echo "Installing AI service dependencies..."
	cd services/ai-analysis-service && pip install -r requirements.txt
	@echo "Dependencies installed!"

# Code Quality
lint: ## Run linters for all code
	@echo "Running Go linters..."
	golangci-lint run ./services/...
	golangci-lint run ./sensor/...
	@echo "Running frontend linter..."
	cd web-ui && npm run lint
	@echo "Running Python linter..."
	cd services/ai-analysis-service && flake8 .

format: ## Format all code
	@echo "Formatting Go code..."
	gofmt -s -w ./services/
	gofmt -s -w ./sensor/
	@echo "Formatting frontend code..."
	cd web-ui && npm run format
	@echo "Formatting Python code..."
	cd services/ai-analysis-service && black .

# Security
security-scan: ## Run security scans
	@echo "Running Go security scan..."
	gosec ./services/...
	gosec ./sensor/...
	@echo "Running npm audit..."
	cd web-ui && npm audit
	@echo "Running Python security scan..."
	cd services/ai-analysis-service && safety check

# Docker Commands
docker-build: ## Build all Docker images
	docker-compose build

docker-pull: ## Pull latest base images
	docker-compose pull

# Sensor Deployment
sensor-package: build-sensor ## Package sensor for distribution
	@echo "Creating sensor distribution packages..."
	mkdir -p dist/sensor
	cp bin/crypto-sensor-* dist/sensor/
	cp sensor/README.md dist/sensor/
	cp sensor/config.example.yaml dist/sensor/config.yaml
	cd dist && tar -czf crypto-sensor-release.tar.gz sensor/
	@echo "Sensor packages created in dist/ directory"

# Documentation
docs-serve: ## Serve documentation locally
	@echo "Starting documentation server..."
	cd docs && python3 -m http.server 8000

# Monitoring
monitor: ## Show system status
	@echo "=== Docker Services ==="
	docker-compose ps
	@echo ""
	@echo "=== System Resources ==="
	docker stats --no-stream
	@echo ""
	@echo "=== Service Health ==="
	curl -s http://localhost:8081/health || echo "Auth service: DOWN"
	curl -s http://localhost:8082/health || echo "Inventory service: DOWN"
	curl -s http://localhost:8083/health || echo "Compliance service: DOWN"
