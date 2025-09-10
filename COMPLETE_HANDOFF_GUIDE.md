# Complete Platform Handoff Guide

*Last Updated: 2025-01-09*
*Platform: Crypto Inventory Management System with SaaS Admin Separation*

## 🎯 Executive Summary

This document provides a comprehensive handoff guide for the Crypto Inventory Management System, a multi-tenant SaaS platform for managing cryptocurrency network assets, sensors, and compliance monitoring. The platform features a complete separation between tenant-level and platform-level administration, with dedicated interfaces and services for each.

## 📋 Quick Start Guide

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)
- PostgreSQL client (optional, for database access)

### Starting the Platform

1. **Start Backend Services**:
```bash
cd /home/bwoodward/CodeProjects/X
docker-compose up -d
```

2. **Start Tenant UI**:
```bash
cd web-ui
npm run dev
# Access at http://localhost:3001
```

3. **Start SaaS Admin UI**:
```bash
cd saas-admin-ui
python3 -m http.server 3002
# Access at http://localhost:3002/simple.html
```

### Default Credentials

**Tenant Users**:
- Email: `demo@example.com`
- Password: `demo123`

**Platform Admins**:
- Email: `admin@crypto-inventory.com`
- Password: `admin123`

## 🏗️ System Architecture

### Service Overview

| Service | Port | Purpose | Technology |
|---------|------|---------|------------|
| **Auth Service** | 8081 | Tenant authentication | Go + Gin + JWT |
| **Inventory Service** | 8082 | Asset management | Go + Gin + PostgreSQL |
| **SaaS Admin Service** | 8084 | Platform administration | Go + Gin + JWT |
| **Tenant UI** | 3001 | Tenant interface | React + TypeScript + Vite |
| **SaaS Admin UI** | 3002 | Platform admin interface | HTML/JS + TailwindCSS |
| **PostgreSQL** | 5432 | Primary database | PostgreSQL 15 |
| **Redis** | 6379 | Caching & sessions | Redis 7 |
| **InfluxDB** | 8086 | Time series data | InfluxDB 2.7 |

### Data Flow

```
┌─────────────────┐    ┌─────────────────┐
│   Tenant UI     │    │  SaaS Admin UI  │
│   (:3001)       │    │   (:3002)       │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────────────────────────────┐
│           Backend Services              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │  Auth   │ │Inventory│ │SaaS Admin│   │
│  │ :8081   │ │ :8082   │ │ :8084   │   │
│  └─────────┘ └─────────┘ └─────────┘   │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│            Data Layer                   │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │PostgreSQL│ │  Redis  │ │InfluxDB │   │
│  │ :5432    │ │ :6379   │ │ :8086   │   │
│  └─────────┘ └─────────┘ └─────────┘   │
└─────────────────────────────────────────┘
```

## 📁 Project Structure

```
/
├── services/                    # Backend services
│   ├── auth-service/           # Tenant authentication (port 8081)
│   ├── inventory-service/      # Asset management (port 8082)
│   └── saas-admin-service/     # Platform administration (port 8084)
├── web-ui/                     # Tenant React application (port 3001)
├── saas-admin-ui/              # Platform admin interface (port 3002)
├── scripts/                    # Database scripts and migrations
│   └── database/
│       ├── 001_auth_schema.sql
│       ├── migrations.sql
│       └── seed.sql
├── docs/                       # Documentation
│   ├── API_DOCUMENTATION.md
│   ├── ARCHITECTURE_DOCUMENTATION.md
│   └── SAAS_ADMIN_SEPARATION_HANDOFF.md
├── docker-compose.yml          # Development environment
└── README.md                   # Project overview
```

## 🔧 Development Workflow

### Backend Development

1. **Service Structure**:
   - Each service follows Go project layout
   - `cmd/main.go` - Entry point
   - `internal/` - Private application code
   - `go.mod` - Dependencies

2. **Database Changes**:
   - Update SQL scripts in `/scripts/database/`
   - Run migrations: `docker exec crypto-postgres psql -U crypto_user -d crypto_inventory -f /scripts/migrations.sql`

3. **Testing Services**:
   ```bash
   # Check service health
   curl http://localhost:8081/health
   curl http://localhost:8082/health
   curl http://localhost:8084/health
   ```

### Frontend Development

1. **Tenant UI** (React):
   ```bash
   cd web-ui
   npm install
   npm run dev
   ```

2. **SaaS Admin UI** (HTML/JS):
   ```bash
   cd saas-admin-ui
   python3 -m http.server 3002
   ```

## 🗄️ Database Schema

### Core Tables

**Authentication & Users**:
- `users` - Tenant user accounts
- `tenants` - Tenant organizations
- `subscription_tiers` - Billing tiers
- `platform_users` - Platform administrators

**Asset Management**:
- `network_assets` - Network devices and systems
- `crypto_implementations` - Crypto implementations
- `sensors` - Monitoring sensors
- `certificates` - SSL/TLS certificates

**RBAC System**:
- `tenant_roles` - Tenant-level roles
- `tenant_permissions` - Tenant-level permissions
- `platform_roles` - Platform-level roles
- `platform_permissions` - Platform-level permissions

### Database Access

```bash
# Connect to database
docker exec -it crypto-postgres psql -U crypto_user -d crypto_inventory

# View all tables
\dt

# View specific table
SELECT * FROM tenants LIMIT 5;
```

## 🔐 Security & Authentication

### Authentication Flow

1. **Tenant Users**:
   - Login via Auth Service (port 8081)
   - JWT token with tenant context
   - Access to tenant-specific data only

2. **Platform Admins**:
   - Login via SaaS Admin Service (port 8084)
   - JWT token with platform context
   - Access to platform-wide data

### Role-Based Access Control

**Tenant Roles**:
- `tenant_owner` - Full tenant access
- `tenant_admin` - Tenant management
- `security_admin` - Security settings
- `analyst` - Data analysis

**Platform Roles**:
- `super_admin` - Full platform access
- `platform_admin` - Platform management
- `support_admin` - Support and monitoring

## 📊 API Documentation

### Key Endpoints

**Authentication Service** (`:8081`):
- `POST /api/v1/auth/login` - Tenant user login
- `POST /api/v1/auth/register` - Tenant user registration
- `POST /api/v1/auth/refresh` - Token refresh

**Inventory Service** (`:8082`):
- `GET /api/v1/assets` - List assets
- `POST /api/v1/assets` - Create asset
- `GET /api/v1/assets/:id` - Get asset details

**SaaS Admin Service** (`:8084`):
- `POST /api/v1/auth/login` - Platform admin login
- `GET /api/v1/admin/tenants` - List all tenants
- `GET /api/v1/admin/stats/platform` - Platform statistics

For complete API documentation, see `API_DOCUMENTATION.md`.

## 🚀 Deployment

### Development Environment

The current setup uses Docker Compose for local development:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Considerations

1. **Environment Variables**:
   - Set production database URLs
   - Configure JWT secrets
   - Set up proper CORS origins

2. **Database**:
   - Use managed PostgreSQL service
   - Set up read replicas
   - Configure backups

3. **Frontend**:
   - Build and serve with nginx/apache
   - Set up CDN for static assets
   - Configure HTTPS

## 🔍 Monitoring & Troubleshooting

### Health Checks

```bash
# Check all services
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8084/health
```

### Common Issues

1. **Service Won't Start**:
   - Check Docker logs: `docker-compose logs <service-name>`
   - Verify port availability
   - Check database connection

2. **Authentication Failures**:
   - Verify JWT secrets match
   - Check user exists in database
   - Validate password hashes

3. **Frontend Issues**:
   - Check browser console for errors
   - Verify API endpoints are accessible
   - Check CORS configuration

### Debug Commands

```bash
# View service logs
docker-compose logs -f auth-service
docker-compose logs -f inventory-service
docker-compose logs -f saas-admin-service

# Check database
docker exec crypto-postgres psql -U crypto_user -d crypto_inventory -c "SELECT * FROM users LIMIT 5;"

# Test API endpoints
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"demo123"}'
```

## 📚 Documentation References

- **API Documentation**: `API_DOCUMENTATION.md`
- **Architecture Documentation**: `ARCHITECTURE_DOCUMENTATION.md`
- **SaaS Admin Separation**: `SAAS_ADMIN_SEPARATION_HANDOFF.md`
- **Main Platform Handoff**: `MAIN_PLATFORM_HANDOFF.md`

## 🆘 Support & Maintenance

### Regular Maintenance

1. **Database Backups**: Set up automated backups
2. **Log Rotation**: Configure log rotation for services
3. **Security Updates**: Keep dependencies updated
4. **Monitoring**: Set up alerts for service health

### Development Team Handoff

1. **Code Ownership**: All services are well-documented and commented
2. **Database Schema**: Complete schema with migrations
3. **API Documentation**: Comprehensive endpoint documentation
4. **Frontend Components**: React components with TypeScript
5. **Testing**: Manual testing completed, unit tests recommended

## 🎯 Next Steps

### Immediate Priorities

1. **Unit Testing**: Add comprehensive test suites
2. **Integration Testing**: End-to-end testing
3. **Performance Testing**: Load testing and optimization
4. **Security Audit**: Comprehensive security review

### Future Enhancements

1. **Advanced Analytics**: Platform-wide usage analytics
2. **Billing Integration**: Stripe/payment processing
3. **Multi-Factor Authentication**: Enhanced security
4. **API Rate Limiting**: Platform-wide rate limiting
5. **Audit Logging**: Comprehensive audit trails

---

*This handoff guide should be updated as the platform evolves. Last updated: 2025-01-09*
