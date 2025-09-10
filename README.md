# Crypto Inventory Management System

A comprehensive multi-tenant SaaS platform for managing cryptocurrency network assets, sensors, and compliance monitoring with complete separation between tenant and platform administration.

## 🎯 Overview

The Crypto Inventory Management System is a production-ready platform that provides:

- **Multi-tenant Architecture**: Complete tenant isolation with shared infrastructure
- **Asset Management**: Comprehensive network asset and crypto implementation tracking
- **Sensor Monitoring**: Real-time network sensor deployment and monitoring
- **RBAC System**: Role-based access control at both tenant and platform levels
- **SaaS Admin Console**: Dedicated platform administration interface
- **Compliance Tracking**: Built-in compliance assessment and reporting

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)

### Installation

1. **Clone the repository**:
```bash
git clone <repository-url>
cd crypto-inventory-platform
```

2. **Start the platform**:
```bash
# Start all backend services
docker-compose up -d

# Start tenant UI (in one terminal)
cd web-ui
npm install
npm run dev

# Start SaaS admin UI (in another terminal)
cd saas-admin-ui
python3 -m http.server 3002
```

3. **Access the platform**:
- **Tenant UI**: http://localhost:3001
- **SaaS Admin Console**: http://localhost:3002/simple.html

### Default Credentials

**Tenant Users**:
- Email: `demo@example.com`
- Password: `demo123`

**Platform Admins**:
- Email: `admin@crypto-inventory.com`
- Password: `admin123`

## 🏗️ Architecture

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

### Key Features

#### Tenant Management
- **Asset Management**: Track network devices, servers, and applications
- **Sensor Deployment**: Deploy and monitor network sensors
- **User Management**: Manage tenant users and roles
- **Compliance Tracking**: Monitor compliance with various frameworks
- **Reporting**: Generate comprehensive reports and analytics

#### Platform Administration
- **Tenant Management**: Create, update, suspend, and activate tenants
- **User Management**: Manage platform administrators
- **Statistics**: Platform-wide metrics and analytics
- **Monitoring**: System health and performance monitoring
- **RBAC**: Platform-level role-based access control

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
│   ├── SAAS_ADMIN_SEPARATION_HANDOFF.md
│   └── COMPLETE_HANDOFF_GUIDE.md
├── docker-compose.yml          # Development environment
└── README.md                   # This file
```

## 🔧 Development

### Backend Development

Each service follows Go project layout:

```
services/auth-service/
├── cmd/main.go                 # Entry point
├── internal/
│   ├── api/                   # HTTP handlers
│   ├── config/                # Configuration
│   ├── database/              # Database connection
│   ├── middleware/            # HTTP middleware
│   └── models/                # Data models
├── go.mod                     # Dependencies
└── Dockerfile.dev             # Development Dockerfile
```

### Frontend Development

#### Tenant UI (React)
```bash
cd web-ui
npm install
npm run dev
```

#### SaaS Admin UI (HTML/JS)
```bash
cd saas-admin-ui
python3 -m http.server 3002
```

### Database Management

```bash
# Connect to database
docker exec -it crypto-postgres psql -U crypto_user -d crypto_inventory

# Run migrations
docker exec crypto-postgres psql -U crypto_user -d crypto_inventory -f /scripts/migrations.sql

# View tables
\dt
```

## 🔐 Security

### Authentication & Authorization

- **JWT-based Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Granular permissions at tenant and platform levels
- **Multi-tenant Isolation**: Complete data isolation between tenants
- **Password Security**: Argon2id hashing with salt
- **CORS Protection**: Configured for specific origins

### Role Hierarchy

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

For complete API documentation, see [API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md).

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

1. **Environment Variables**: Set production database URLs and JWT secrets
2. **Database**: Use managed PostgreSQL service with read replicas
3. **Frontend**: Build and serve with nginx/apache
4. **Security**: Configure HTTPS and proper CORS origins
5. **Monitoring**: Set up health checks and logging

## 🔍 Monitoring & Troubleshooting

### Health Checks

```bash
# Check all services
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8084/health
```

### Common Issues

1. **Service Won't Start**: Check Docker logs and port availability
2. **Authentication Failures**: Verify JWT secrets and user data
3. **Frontend Issues**: Check browser console and API accessibility

### Debug Commands

```bash
# View service logs
docker-compose logs -f auth-service
docker-compose logs -f inventory-service
docker-compose logs -f saas-admin-service

# Check database
docker exec crypto-postgres psql -U crypto_user -d crypto_inventory -c "SELECT * FROM users LIMIT 5;"
```

## 📚 Documentation

- **[API Documentation](docs/API_DOCUMENTATION.md)**: Complete API reference
- **[Architecture Documentation](docs/ARCHITECTURE_DOCUMENTATION.md)**: System architecture details
- **[SaaS Admin Separation](docs/SAAS_ADMIN_SEPARATION_HANDOFF.md)**: Platform admin implementation
- **[Complete Handoff Guide](docs/COMPLETE_HANDOFF_GUIDE.md)**: Comprehensive handoff documentation

## 🎯 Roadmap

### Immediate Priorities
- [ ] Unit and integration testing
- [ ] Performance optimization
- [ ] Security audit
- [ ] Documentation updates

### Future Enhancements
- [ ] Advanced analytics and reporting
- [ ] Billing and subscription management
- [ ] Multi-factor authentication
- [ ] API rate limiting
- [ ] Audit logging system

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Check the documentation in the `docs/` directory
- Review the troubleshooting section above
- Create an issue in the repository

---

*Last updated: 2025-01-09*

# Crypto Inventory - Web UI Updates (Dashboard, Roles, Reports, Sensors)

This update adds tenant dashboard data wiring, RBAC integration, report management, and sensor registration flows.

## What’s new
- Dashboard
  - Overview stats backed by inventory APIs via `dashboardApi`
  - Recent Activity from auth-service `/audit/logs` with time range (24h/7d/30d)
  - Expiring Certificates panel from inventory-service
  - Top Risks panel from risk summary
  - Quick actions navigate to Assets, Reports, Roles
- Roles & Permissions
  - `RoleManagement` loads tenant roles and permissions via RBAC endpoints
  - Permission toggles persist to backend
- Reports
  - `ReportsPage` lists reports, generates new ones, polls status, supports delete
- Sensors
  - `SensorRegistrationPage` lists/creates/deletes pending registration keys using sensor-manager

## Local services (defaults)
- auth-service: `http://localhost:8081` (configured in `web-ui/src/services/api.ts`)
- inventory-service: `http://localhost:8082`
- report-generator: `http://localhost:8083`
- sensor-manager: `http://localhost:8080`

Ensure these are running or adjust base URLs inside the respective `*Api` files as needed.

## Testing
- Added Vitest + Testing Library setup
- Tests:
  - `src/pages/__tests__/DashboardPage.test.tsx`
  - `src/pages/__tests__/ReportsPage.test.tsx`
  - `src/pages/__tests__/SensorRegistrationPage.test.tsx`
  - `src/components/rbac/__tests__/RoleManagement.test.tsx`

Run tests:

```
pnpm test
```

or

```
npm run test
```