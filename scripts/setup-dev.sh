#!/bin/bash

# Development Environment Setup Script
# This script sets up the complete development environment for the Crypto Inventory Platform

set -e

echo "🚀 Setting up Crypto Inventory Platform Development Environment"
echo "=============================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    local missing_deps=()
    
    # Check Docker
    if ! command_exists docker; then
        missing_deps+=("docker")
    fi
    
    # Check Docker Compose
    if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
        missing_deps+=("docker-compose")
    fi
    
    # Check Go
    if ! command_exists go; then
        missing_deps+=("go (version 1.21+)")
    else
        go_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
        if [[ $(echo "$go_version 1.21" | tr " " "\n" | sort -V | head -n1) != "1.21" ]]; then
            log_warning "Go version $go_version detected. Recommended: 1.21+"
        fi
    fi
    
    # Check Node.js
    if ! command_exists node; then
        missing_deps+=("node.js (version 18+)")
    else
        node_version=$(node --version | sed 's/v//')
        if [[ $(echo "$node_version 18.0.0" | tr " " "\n" | sort -V | head -n1) != "18.0.0" ]]; then
            log_warning "Node.js version $node_version detected. Recommended: 18+"
        fi
    fi
    
    # Check Python
    if ! command_exists python3; then
        missing_deps+=("python3 (version 3.11+)")
    else
        python_version=$(python3 --version | grep -oE '[0-9]+\.[0-9]+')
        if [[ $(echo "$python_version 3.11" | tr " " "\n" | sort -V | head -n1) != "3.11" ]]; then
            log_warning "Python version $python_version detected. Recommended: 3.11+"
        fi
    fi
    
    # Check Git
    if ! command_exists git; then
        missing_deps+=("git")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing required dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo ""
        echo "Please install the missing dependencies and run this script again."
        echo "Installation guides:"
        echo "  - Docker: https://docs.docker.com/get-docker/"
        echo "  - Go: https://golang.org/doc/install"
        echo "  - Node.js: https://nodejs.org/en/download/"
        echo "  - Python: https://www.python.org/downloads/"
        exit 1
    fi
    
    log_success "All prerequisites are installed"
}

# Initialize Go modules
setup_go_services() {
    log_info "Setting up Go services..."
    
    services=("auth-service" "inventory-service" "compliance-engine" "report-generator" "sensor-manager")
    
    for service in "${services[@]}"; do
        if [ -d "services/$service" ]; then
            log_info "Downloading dependencies for $service..."
            cd "services/$service"
            go mod download
            go mod tidy
            cd ../..
        fi
    done
    
    # Setup sensor
    if [ -d "sensor" ]; then
        log_info "Downloading dependencies for network sensor..."
        cd sensor
        go mod download
        go mod tidy
        cd ..
    fi
    
    log_success "Go services initialized"
}

# Setup Python AI service
setup_ai_service() {
    log_info "Setting up AI Analysis Service..."
    
    if [ -d "services/ai-analysis-service" ]; then
        cd "services/ai-analysis-service"
        
        # Create virtual environment if it doesn't exist
        if [ ! -d "venv" ]; then
            log_info "Creating Python virtual environment..."
            python3 -m venv venv
        fi
        
        # Activate virtual environment and install dependencies
        log_info "Installing Python dependencies..."
        source venv/bin/activate
        pip install --upgrade pip
        pip install -r requirements.txt
        deactivate
        
        cd ../..
        log_success "AI Analysis Service initialized"
    else
        log_warning "AI Analysis Service directory not found, skipping..."
    fi
}

# Setup frontend
setup_frontend() {
    log_info "Setting up Web Frontend..."
    
    if [ -d "web-ui" ]; then
        cd "web-ui"
        
        # Check if package.json exists
        if [ ! -f "package.json" ]; then
            log_info "Initializing frontend project..."
            
            # Create basic package.json
            cat > package.json << 'EOF'
{
  "name": "crypto-inventory-web-ui",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "@testing-library/jest-dom": "^5.16.4",
    "@testing-library/react": "^13.3.0",
    "@testing-library/user-event": "^13.5.0",
    "@types/jest": "^27.5.2",
    "@types/node": "^16.11.47",
    "@types/react": "^18.0.15",
    "@types/react-dom": "^18.0.6",
    "antd": "^5.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.3.0",
    "react-scripts": "5.0.1",
    "typescript": "^4.7.4",
    "web-vitals": "^2.1.4"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test",
    "eject": "react-scripts eject",
    "lint": "eslint src --ext .ts,.tsx",
    "format": "prettier --write src/**/*.{ts,tsx}"
  },
  "eslintConfig": {
    "extends": [
      "react-app",
      "react-app/jest"
    ]
  },
  "browserslist": {
    "production": [
      ">0.2%",
      "not dead",
      "not op_mini all"
    ],
    "development": [
      "last 1 chrome version",
      "last 1 firefox version",
      "last 1 safari version"
    ]
  },
  "devDependencies": {
    "@typescript-eslint/eslint-plugin": "^5.30.7",
    "@typescript-eslint/parser": "^5.30.7",
    "eslint": "^8.20.0",
    "prettier": "^2.7.1"
  }
}
EOF
        fi
        
        log_info "Installing frontend dependencies..."
        npm install
        
        cd ..
        log_success "Web Frontend initialized"
    else
        log_warning "Web UI directory not found, skipping..."
    fi
}

# Create development databases
setup_databases() {
    log_info "Setting up development databases..."
    
    # Start only database services
    docker-compose up -d postgres redis influxdb
    
    # Wait for databases to be ready
    log_info "Waiting for databases to be ready..."
    sleep 15
    
    # Check if databases are responding
    if docker-compose exec -T postgres pg_isready -U crypto_user -d crypto_inventory >/dev/null 2>&1; then
        log_success "PostgreSQL is ready"
    else
        log_warning "PostgreSQL might not be fully ready yet"
    fi
    
    if docker-compose exec -T redis redis-cli ping >/dev/null 2>&1; then
        log_success "Redis is ready"
    else
        log_warning "Redis might not be fully ready yet"
    fi
    
    log_success "Development databases are running"
}

# Create development configuration
create_dev_config() {
    log_info "Creating development configuration..."
    
    # Create .env file for development
    if [ ! -f ".env" ]; then
        cat > .env << 'EOF'
# Development Environment Configuration
NODE_ENV=development
ENV=development
LOG_LEVEL=debug

# Database Configuration
DATABASE_URL=postgres://crypto_user:crypto_pass_dev@localhost:5432/crypto_inventory?sslmode=disable
REDIS_URL=redis://:redis_pass_dev@localhost:6379/0

# InfluxDB Configuration
INFLUXDB_URL=http://localhost:8086
INFLUXDB_TOKEN=dev-token-1234567890
INFLUXDB_ORG=crypto-inventory
INFLUXDB_BUCKET=metrics

# NATS Configuration
NATS_URL=nats://nats_user:nats_pass_dev@localhost:4222

# JWT Configuration
JWT_SECRET=dev-secret-key-change-in-production
JWT_EXPIRY=24h

# Frontend Configuration
REACT_APP_API_URL=http://localhost:8080
REACT_APP_WS_URL=ws://localhost:8080
REACT_APP_ENV=development

# CORS Configuration
CORS_ORIGINS=http://localhost:3000
EOF
        log_success "Created .env file for development"
    else
        log_info ".env file already exists, skipping..."
    fi
}

# Build development binaries
build_services() {
    log_info "Building development binaries..."
    
    # Create bin directory
    mkdir -p bin
    
    # Build Go services
    services=("auth-service" "inventory-service" "compliance-engine" "report-generator" "sensor-manager")
    
    for service in "${services[@]}"; do
        if [ -d "services/$service" ]; then
            log_info "Building $service..."
            cd "services/$service"
            go build -o "../../bin/$service" ./cmd/main.go
            cd ../..
        fi
    done
    
    # Build sensor
    if [ -d "sensor" ]; then
        log_info "Building network sensor..."
        cd sensor
        go build -o "../bin/crypto-sensor" ./cmd/main.go
        cd ..
    fi
    
    log_success "Development binaries built"
}

# Verify setup
verify_setup() {
    log_info "Verifying development setup..."
    
    local errors=0
    
    # Check if binaries exist
    binaries=("auth-service" "inventory-service" "compliance-engine" "report-generator" "sensor-manager" "crypto-sensor")
    for binary in "${binaries[@]}"; do
        if [ ! -f "bin/$binary" ]; then
            log_error "Binary bin/$binary not found"
            ((errors++))
        fi
    done
    
    # Check if databases are running
    if ! docker-compose ps | grep -q "postgres.*Up"; then
        log_error "PostgreSQL container is not running"
        ((errors++))
    fi
    
    if ! docker-compose ps | grep -q "redis.*Up"; then
        log_error "Redis container is not running"
        ((errors++))
    fi
    
    if [ $errors -eq 0 ]; then
        log_success "Development environment setup completed successfully!"
        echo ""
        echo "🎉 Next Steps:"
        echo "  1. Start all services: make start"
        echo "  2. View logs: make logs"
        echo "  3. Run tests: make test"
        echo "  4. Access Web UI: http://localhost:3000"
        echo "  5. Access API: http://localhost:8080"
        echo "  6. Access Database Admin: http://localhost:8090"
        echo ""
        echo "📖 Documentation:"
        echo "  - Architecture: ./architecture_docs/"
        echo "  - Contributing: ./CONTRIBUTING.md"
        echo "  - API Docs: http://localhost:8080/docs (when running)"
    else
        log_error "Setup completed with $errors errors. Please review and fix the issues above."
        exit 1
    fi
}

# Main setup function
main() {
    echo "Starting development environment setup..."
    echo "This script will:"
    echo "  ✓ Check prerequisites"
    echo "  ✓ Initialize Go services"
    echo "  ✓ Setup Python AI service"
    echo "  ✓ Setup React frontend"
    echo "  ✓ Start development databases"
    echo "  ✓ Create development configuration"
    echo "  ✓ Build development binaries"
    echo "  ✓ Verify setup"
    echo ""
    
    read -p "Continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Setup cancelled."
        exit 0
    fi
    
    check_prerequisites
    setup_go_services
    setup_ai_service
    setup_frontend
    create_dev_config
    setup_databases
    build_services
    verify_setup
}

# Run main function
main "$@"
