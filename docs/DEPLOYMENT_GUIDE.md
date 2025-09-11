# Deployment Guide

This guide covers building, deploying, and troubleshooting the Crypto Inventory Management System.

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose**: Latest version
- **System Dependencies** (Ubuntu/Debian):
  ```bash
  sudo apt-get update
  sudo apt-get install -y build-essential libpcap-dev python3 python3-pip python3-venv python3.12-venv
  ```

### Build & Deploy

1. **Clone and build**:
   ```bash
   git clone <repository-url>
   cd crypto-inventory-platform
   make build-all
   ```

2. **Start services**:
   ```bash
   docker-compose up -d
   ```

3. **Verify deployment**:
   ```bash
   docker-compose ps
   curl -s http://localhost:8081/health
   curl -s http://localhost:8082/health
   curl -s http://localhost:8084/health
   ```

4. **Load sample data** (optional):
   ```bash
   docker-compose exec postgres psql -U crypto_user -d crypto_inventory -f /scripts/database/seed_brian_debban.sql
   ```

## 🔧 Build Process

### Makefile Targets

| Target | Description | Dependencies |
|--------|-------------|--------------|
| `build-all` | Build all components | All individual build targets |
| `build-services` | Build Go services | Go 1.21+, libpcap-dev |
| `build-sensor` | Build network sensor | Go 1.21+, libpcap-dev, CGO |
| `build-frontend` | Build React UI | Node.js 18+ |
| `build-ai-service` | Build AI service | Python 3.12+, venv |

### Build Troubleshooting

#### Go Services Build Issues

**Error**: `undefined: pcapErrorNotActivated`
```bash
# Solution: Install libpcap development headers
sudo apt-get install -y libpcap-dev build-essential
CGO_ENABLED=1 go build
```

**Error**: `handler.SubmitDiscoveries undefined`
```bash
# Solution: Missing handler methods - check services/sensor-manager/internal/handlers/outbound.go
# Ensure all required methods are implemented
```

#### Frontend Build Issues

**Error**: `npm ci` lockfile mismatch
```bash
# Solution: Use npm install instead
npm install
# Or update package-lock.json
npm ci --legacy-peer-deps
```

#### AI Service Build Issues

**Error**: `pip: not found`
```bash
# Solution: Install Python and pip
sudo apt-get install -y python3 python3-pip python3-venv python3.12-venv
```

**Error**: `externally-managed-environment`
```bash
# Solution: Use virtual environment (handled by Makefile)
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## 🐳 Docker Services

### Service Architecture

| Service | Port | Health Check | Dependencies |
|---------|------|--------------|--------------|
| **postgres** | 5432 | TCP connection | - |
| **redis** | 6379 | TCP connection | postgres |
| **influxdb** | 8086 | HTTP :8086/health | postgres |
| **nats** | 4222 | TCP connection | postgres |
| **auth-service** | 8081 | GET /health | postgres, redis |
| **inventory-service** | 8082 | GET /health | postgres, redis |
| **report-generator** | 8083 | GET /health | postgres, redis |
| **saas-admin-service** | 8084 | GET /health | postgres, redis |
| **sensor-manager** | 8085 | GET /health | postgres, redis |
| **ai-analysis-service** | 8087 | GET /health | postgres, redis |
| **compliance-engine** | 8088 | GET /health | postgres, redis |
| **api-gateway** | 8080 | GET /health | All services |
| **web-ui** | 3000 | HTTP :3000 | api-gateway |
| **saas-admin-ui** | 3002 | HTTP :3002 | api-gateway |

### Container Management

```bash
# View all containers
docker-compose ps

# View logs
docker-compose logs -f [service-name]

# Restart specific service
docker-compose restart [service-name]

# Scale services
docker-compose up -d --scale inventory-service=2

# Clean rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 🔍 Health Monitoring

### Service Health Checks

```bash
# Check all services
curl -s http://localhost:8080/health || echo "API Gateway down"
curl -s http://localhost:8081/health || echo "Auth Service down"
curl -s http://localhost:8082/health || echo "Inventory Service down"
curl -s http://localhost:8084/health || echo "SaaS Admin Service down"
curl -s http://localhost:8085/health || echo "Sensor Manager down"

# Check database
docker-compose exec postgres pg_isready -U crypto_user

# Check Redis
docker-compose exec redis redis-cli ping
```

### Log Monitoring

```bash
# Follow all logs
docker-compose logs -f

# Follow specific service
docker-compose logs -f auth-service
docker-compose logs -f inventory-service
docker-compose logs -f api-gateway

# View recent logs
docker-compose logs --tail=100 auth-service
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Services Won't Start

**Symptoms**: Containers exit immediately or fail to start
```bash
# Check logs
docker-compose logs [service-name]

# Check port conflicts
netstat -tulpn | grep :8080
netstat -tulpn | grep :5432

# Check disk space
df -h
```

**Solutions**:
- Free up disk space
- Stop conflicting services
- Check Docker daemon status: `systemctl status docker`

#### 2. Database Connection Issues

**Symptoms**: Services can't connect to PostgreSQL
```bash
# Check database status
docker-compose exec postgres pg_isready -U crypto_user

# Check database logs
docker-compose logs postgres

# Reset database
docker-compose down
docker volume rm crypto-inventory_postgres_data
docker-compose up -d postgres
```

#### 3. Frontend Not Loading

**Symptoms**: Blank page or connection errors
```bash
# Check API gateway
curl http://localhost:8080/health

# Check frontend container
docker-compose logs web-ui

# Check nginx config
docker-compose exec api-gateway nginx -t
```

#### 4. Authentication Failures

**Symptoms**: Login fails or tokens invalid
```bash
# Check auth service
curl http://localhost:8081/health

# Check database users
docker-compose exec postgres psql -U crypto_user -d crypto_inventory -c "SELECT email, is_active FROM users;"

# Check JWT secrets in environment
docker-compose exec auth-service env | grep JWT
```

#### 5. Build Failures

**Symptoms**: `make build-*` commands fail
```bash
# Check system dependencies
dpkg -l | grep libpcap-dev
dpkg -l | grep build-essential
python3 --version
node --version

# Clean and rebuild
make clean
make build-all
```

### Performance Issues

#### High Memory Usage
```bash
# Check container memory
docker stats

# Check specific service
docker-compose exec auth-service top
docker-compose exec inventory-service top
```

#### Slow Database Queries
```bash
# Connect to database
docker-compose exec postgres psql -U crypto_user -d crypto_inventory

# Check slow queries
SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;

# Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size 
FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## 🔐 Security Considerations

### Production Deployment

1. **Environment Variables**:
   ```bash
   # Set strong JWT secrets
   export JWT_SECRET="your-strong-secret-here"
   export JWT_REFRESH_SECRET="your-refresh-secret-here"
   
   # Set database passwords
   export POSTGRES_PASSWORD="strong-db-password"
   export REDIS_PASSWORD="strong-redis-password"
   ```

2. **Network Security**:
   - Use reverse proxy (nginx/traefik)
   - Enable HTTPS with Let's Encrypt
   - Configure firewall rules
   - Use Docker networks for service isolation

3. **Database Security**:
   - Use managed PostgreSQL service
   - Enable encryption at rest
   - Regular backups
   - Monitor access logs

### Monitoring & Alerting

```bash
# Set up health check monitoring
*/5 * * * * curl -f http://localhost:8080/health || echo "API Gateway down" | mail -s "Service Alert" admin@company.com

# Monitor disk space
*/10 * * * * df -h | awk '$5 > 80 {print $0}' | mail -s "Disk Space Alert" admin@company.com
```

## 📊 Maintenance

### Regular Tasks

1. **Database Maintenance**:
   ```bash
   # Backup database
   docker-compose exec postgres pg_dump -U crypto_user crypto_inventory > backup_$(date +%Y%m%d).sql
   
   # Vacuum database
   docker-compose exec postgres psql -U crypto_user -d crypto_inventory -c "VACUUM ANALYZE;"
   ```

2. **Log Rotation**:
   ```bash
   # Configure Docker log rotation
   # Add to docker-compose.yml:
   logging:
     driver: "json-file"
     options:
       max-size: "10m"
       max-file: "3"
   ```

3. **Security Updates**:
   ```bash
   # Update base images
   docker-compose pull
   docker-compose up -d
   
   # Update system packages
   sudo apt-get update && sudo apt-get upgrade
   ```

## 🆘 Support

### Debug Information

When reporting issues, include:

```bash
# System information
uname -a
docker --version
docker-compose --version

# Service status
docker-compose ps
docker-compose logs --tail=50

# Resource usage
docker stats --no-stream
df -h
free -h
```

### Emergency Recovery

```bash
# Complete reset
docker-compose down -v
docker system prune -a
git clean -fd
make build-all
docker-compose up -d
```

---

*Last updated: 2025-01-09*
