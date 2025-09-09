# 🚀 Remote Development Server Handoff Guide
*Created: 2025-01-09*
*Purpose: Smooth transition to new remote development server*

## 📋 Pre-Migration Checklist

### ✅ **Current Status (Local Development)**
- **Platform**: 90% complete with fully functional authentication and frontend
- **Services**: Auth service working, Inventory service implemented, all build issues resolved
- **Frontend**: Complete React TypeScript application with professional UI
- **Database**: Enhanced PostgreSQL schema with multi-tenant support
- **Infrastructure**: Docker Compose setup with all services defined

### 🔧 **Critical Fixes Applied This Session**
1. **Go Version Compatibility**: Fixed `go.mod` files in compliance-engine, sensor-manager, report-generator
2. **Python Dependencies**: Removed incompatible `pickle5` package from AI service
3. **Build Issues**: All services now build successfully without errors

## 🖥️ Remote Server Setup Requirements

### **System Requirements**
- **OS**: Linux (Ubuntu 20.04+ recommended)
- **Docker**: Docker Engine 20.10+ and Docker Compose 2.0+
- **Node.js**: 18.x or 20.x (for frontend development)
- **Go**: 1.21+ (for backend development)
- **Python**: 3.11+ (for AI service)
- **Memory**: Minimum 8GB RAM (16GB recommended)
- **Storage**: 50GB+ available space

### **Required Ports**
- **8080**: AI Analysis Service
- **8081**: Auth Service
- **8082**: Inventory Service
- **8083**: Compliance Engine
- **8084**: Report Generator
- **8085**: Sensor Manager
- **3000**: Frontend Development Server
- **5432**: PostgreSQL
- **6379**: Redis
- **8086**: InfluxDB
- **4222**: NATS

## 📦 Migration Steps

### **Step 1: Clone Repository**
```bash
# Clone the repository
git clone <repository-url> crypto-inventory
cd crypto-inventory

# Verify you're on the latest commit
git log --oneline -5
# Should see: "feat: Fix build compatibility issues for remote development"
```

### **Step 2: Environment Setup**
```bash
# Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Node.js 18.x
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install Go 1.21
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Python 3.11
sudo apt update
sudo apt install -y python3.11 python3.11-pip python3.11-venv
```

### **Step 3: Start Backend Services**
```bash
# Navigate to project root
cd crypto-inventory

# Start all backend services
docker-compose up -d

# Verify services are running
docker-compose ps
# All services should show "Up" status

# Check service health
curl http://localhost:8081/health  # Auth service
curl http://localhost:8082/health  # Inventory service
```

### **Step 4: Start Frontend Development**
```bash
# Navigate to frontend directory
cd web-ui

# Install dependencies
npm install

# Start development server
npm run dev

# Access at: http://your-server-ip:3000
```

## 🧪 Verification Tests

### **Backend API Tests**
```bash
# Test user registration
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "SecurePassword123!", 
       "first_name": "Test", "last_name": "User", "tenant_name": "Test Company"}'

# Test user login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "SecurePassword123!"}'

# Test inventory service
curl http://localhost:8082/api/v1/assets
```

### **Frontend Tests**
1. **Access**: Navigate to `http://your-server-ip:3000`
2. **Registration**: Create a new account
3. **Login**: Log in with created credentials
4. **Dashboard**: Verify dashboard loads with user info
5. **Assets Page**: Navigate to `/assets` and verify asset table loads
6. **Theme Toggle**: Test light/dark theme switching

## 🔧 Troubleshooting Guide

### **Common Issues & Solutions**

#### **Docker Services Won't Start**
```bash
# Check Docker status
sudo systemctl status docker

# Restart Docker
sudo systemctl restart docker

# Check for port conflicts
sudo netstat -tulpn | grep :8081
```

#### **Frontend Build Errors**
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Check Node.js version
node --version  # Should be 18.x or 20.x
```

#### **Database Connection Issues**
```bash
# Check PostgreSQL container
docker logs crypto-postgres

# Reset database
docker-compose down
docker volume rm crypto-inventory_postgres_data
docker-compose up -d
```

#### **Go Build Errors**
```bash
# Check Go version
go version  # Should be 1.21+

# Clean Go modules
cd services/auth-service
go mod tidy
go mod vendor
```

## 📁 Key Files & Directories

### **Critical Configuration Files**
- `docker-compose.yml` - All service definitions
- `web-ui/package.json` - Frontend dependencies
- `services/*/go.mod` - Go module definitions (all fixed to Go 1.21)
- `services/ai-analysis-service/requirements.txt` - Python dependencies (pickle5 removed)

### **Database Files**
- `scripts/database/001_auth_schema.sql` - Enhanced authentication schema
- `scripts/database/migrations.sql` - Additional schema migrations
- `scripts/database/seed.sql` - Sample data

### **Frontend Files**
- `web-ui/src/` - React TypeScript source code
- `web-ui/vite.config.ts` - Vite build configuration
- `web-ui/tailwind.config.js` - TailwindCSS configuration

## 🚀 Next Development Priorities

### **Phase 2: Enhanced Inventory Management (Next Session)**
1. **Real-time Asset Discovery**: Integrate with network sensors
2. **Advanced Risk Analysis**: Implement ML-based risk scoring
3. **Compliance Automation**: Build compliance framework engine
4. **Report Generation**: Create PDF/Excel report templates

### **Phase 3: Production Readiness**
1. **Security Hardening**: Implement rate limiting, CORS, security headers
2. **Monitoring**: Add Prometheus metrics and Grafana dashboards
3. **Testing**: Comprehensive unit and integration test suite
4. **Documentation**: API documentation with OpenAPI specs

## 📞 Support Information

### **If You Encounter Issues**
1. **Check Logs**: `docker-compose logs [service-name]`
2. **Verify Dependencies**: Ensure all required software is installed
3. **Check Ports**: Ensure no port conflicts
4. **Review This Guide**: Most common issues are covered above

### **Current Working State**
- ✅ **Authentication**: Fully functional with JWT tokens
- ✅ **Frontend**: Complete React app with professional UI
- ✅ **Database**: Enhanced schema with multi-tenant support
- ✅ **Build System**: All services build without errors
- ✅ **Development Environment**: Docker Compose + Vite setup

---

**🎯 Success Criteria**: After following this guide, you should have a fully functional crypto inventory platform running on your remote server with working authentication, frontend, and all backend services building successfully.

**📝 Note**: This handoff document assumes you're starting fresh on a new server. If you're migrating from an existing setup, adapt the steps accordingly.
