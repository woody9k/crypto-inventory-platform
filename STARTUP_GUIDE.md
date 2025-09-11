# 🚀 Crypto Inventory Platform - Startup Guide

## Prerequisites

Before starting the platform, ensure you have:

- **Docker** (version 20.10+)
- **Docker Compose** (version 2.0+)
- **Git** (for cloning updates)
- **8GB+ RAM** (recommended for smooth operation)
- **10GB+ free disk space**

## Quick Start (Recommended)

### 1. Start All Services
```bash
# Navigate to the project directory
cd /path/to/crypto-inventory-platform

# Start all services in the background
docker-compose up -d

# Verify all services are running
docker-compose ps
```

### 2. Wait for Services to Initialize
```bash
# Wait about 30-60 seconds for all services to start
# Check logs if needed
docker-compose logs -f
```

### 3. Access the Platform
- **Main Application**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Grafana Dashboard**: http://localhost:3001 (admin/admin123)
- **Database Admin**: http://localhost:8090

## Step-by-Step Startup Process

### Step 1: Verify Environment
```bash
# Check Docker is running
docker --version
docker-compose --version

# Check available resources
docker system df
```

### Step 2: Start Infrastructure Services First
```bash
# Start databases and message queue first
docker-compose up -d postgres redis influxdb nats

# Wait 10 seconds for databases to initialize
sleep 10

# Check database health
docker-compose logs postgres
docker-compose logs redis
```

### Step 3: Start Application Services
```bash
# Start all application services
docker-compose up -d

# Monitor startup logs
docker-compose logs -f --tail=50
```

### Step 4: Verify All Services
```bash
# Check service status
docker-compose ps

# Test API gateway
curl http://localhost:8080/health

# Test main application
curl http://localhost:3000
```

## Service Architecture

### Core Services
- **API Gateway** (port 8080): Routes all frontend requests
- **Web UI** (port 3000): React frontend application
- **Auth Service** (port 8081): Authentication and authorization
- **Inventory Service** (port 8082): Asset management
- **Report Generator** (port 8083): Report generation and downloads

### Supporting Services
- **PostgreSQL** (port 5432): Main database
- **Redis** (port 6379): Caching and sessions
- **InfluxDB** (port 8086): Time-series data
- **NATS** (port 4222): Message queue
- **Grafana** (port 3001): Monitoring dashboard
- **Adminer** (port 8090): Database administration

## Troubleshooting

### Common Issues

#### 1. Services Won't Start
```bash
# Check Docker daemon
sudo systemctl status docker

# Check available ports
netstat -tulpn | grep -E ':(3000|8080|5432|6379)'

# Check Docker logs
docker-compose logs [service-name]
```

#### 2. Database Connection Issues
```bash
# Restart databases
docker-compose restart postgres redis influxdb

# Check database logs
docker-compose logs postgres
docker-compose logs redis
```

#### 3. Frontend Not Loading
```bash
# Check API gateway
curl http://localhost:8080/health

# Check web UI service
docker-compose logs web-ui

# Restart web UI
docker-compose restart web-ui
```

#### 4. Reports Not Working
```bash
# Check report service
docker-compose logs report-generator

# Test report API
curl http://localhost:8080/api/v1/reports/templates

# Restart report service
docker-compose restart report-generator
```

### Health Checks

#### Check All Services
```bash
# Quick health check
curl http://localhost:8080/health

# Individual service checks
curl http://localhost:8081/health  # Auth service
curl http://localhost:8082/health  # Inventory service
curl http://localhost:8083/health  # Report service
```

#### Check Database Connections
```bash
# PostgreSQL
docker-compose exec postgres psql -U postgres -c "SELECT 1;"

# Redis
docker-compose exec redis redis-cli ping

# InfluxDB
curl http://localhost:8086/health
```

## Development Mode

### Start with Live Reload
```bash
# Start only infrastructure
docker-compose up -d postgres redis influxdb nats

# Start services with live reload (if configured)
docker-compose -f docker-compose.dev.yml up
```

### Frontend Development
```bash
# Start only backend services
docker-compose up -d postgres redis influxdb nats auth-service inventory-service report-generator api-gateway

# Start frontend in development mode
cd web-ui
npm install
npm run dev
```

## Production Deployment

### Using GCP Scripts
```bash
# Start balanced performance setup ($72/month for 4h/day)
./scripts/start-balanced.sh

# Start high performance setup ($270/month for 4h/day)
./scripts/start-production.sh

# Stop services to save costs
./scripts/stop-balanced.sh
./scripts/stop-production.sh
```

## Monitoring and Logs

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f [service-name]

# Last 100 lines
docker-compose logs --tail=100 [service-name]
```

### Monitor Resources
```bash
# Docker resource usage
docker stats

# System resources
htop
```

### Grafana Dashboard
- URL: http://localhost:3001
- Username: admin
- Password: admin123

## Stopping Services

### Graceful Shutdown
```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v

# Stop and remove everything
docker-compose down --rmi all --volumes --remove-orphans
```

### Clean Up
```bash
# Remove unused containers and images
docker system prune -f

# Remove unused volumes
docker volume prune -f

# Complete cleanup (WARNING: removes all Docker data)
docker system prune -a --volumes -f
```

## Default Credentials

### Application Users
- **Platform Admin**: admin@platform.com / admin123
- **Tenant Admin**: admin@democorp.com / admin123
- **Regular User**: user@democorp.com / user123

### Database Access
- **PostgreSQL**: postgres / postgres
- **Redis**: No password required
- **InfluxDB**: admin / admin123

### Monitoring
- **Grafana**: admin / admin123

## API Testing

### Test Core Endpoints
```bash
# Health check
curl http://localhost:8080/health

# Reports templates
curl http://localhost:8080/api/v1/reports/templates

# Generate a report
curl -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{"type": "crypto_summary", "title": "Test Report", "format": "pdf"}'

# Download report (replace {id} with actual report ID)
curl "http://localhost:8080/api/v1/reports/{id}/download?format=pdf"
```

## Next Steps

1. **Access the application** at http://localhost:3000
2. **Login** with the default credentials
3. **Explore the features**:
   - Dashboard with system overview
   - Asset management
   - Report generation and downloads
   - User and role management
   - Sensor monitoring

4. **Generate your first report**:
   - Go to Reports page
   - Click "Generate Report"
   - Select a report type
   - Download in your preferred format

## Support

If you encounter issues:

1. Check the troubleshooting section above
2. Review service logs: `docker-compose logs [service-name]`
3. Verify all services are running: `docker-compose ps`
4. Check system resources: `docker stats`
5. Restart services if needed: `docker-compose restart [service-name]`

For additional help, refer to:
- `docs/REPORTS_SYSTEM.md` - Report system documentation
- `docs/DEPLOYMENT_COST_ANALYSIS.md` - Deployment and cost analysis
- `README.md` - Main project documentation
