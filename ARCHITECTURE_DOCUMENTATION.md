# Architecture Documentation

*Last Updated: 2025-01-09*
*Platform: Crypto Inventory Management System*

## 🎯 System Overview

The Crypto Inventory Management System is a multi-tenant SaaS platform designed for managing cryptocurrency network assets, sensors, and compliance monitoring. The system features a clear separation between tenant-level and platform-level functionality, with dedicated interfaces and services for each.

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                            │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐              ┌─────────────────┐          │
│  │   TENANT UI     │              │  SAAS ADMIN UI  │          │
│  │   (React SPA)   │              │  (HTML/JS SPA)  │          │
│  │   Port: 3001    │              │   Port: 3002    │          │
│  │                 │              │                 │          │
│  │ • Asset Mgmt    │              │ • Tenant Mgmt   │          │
│  │ • Sensor Mgmt   │              │ • User Mgmt     │          │
│  │ • Reports       │              │ • Statistics    │          │
│  │ • Tenant Roles  │              │ • Monitoring    │          │
│  └─────────────────┘              └─────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API GATEWAY LAYER                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │  Auth Service   │  │SaaS Admin Svc   │  │Inventory Service│  │
│  │   Port: 8081    │  │   Port: 8084    │  │   Port: 8082    │  │
│  │                 │  │                 │  │                 │  │
│  │ • JWT Auth      │  │ • Platform Mgmt │  │ • Asset Mgmt    │  │
│  │ • User Mgmt     │  │ • Tenant Mgmt   │  │ • Sensor Mgmt   │  │
│  │ • Tenant Auth   │  │ • Statistics    │  │ • Reports       │  │
│  │ • SSO Support   │  │ • Monitoring    │  │ • Compliance    │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      DATA LAYER                                │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   PostgreSQL    │  │     Redis       │  │    InfluxDB     │  │
│  │   Port: 5432    │  │   Port: 6379    │  │   Port: 8086    │  │
│  │                 │  │                 │  │                 │  │
│  │ • User Data     │  │ • Sessions      │  │ • Metrics       │  │
│  │ • Tenant Data   │  │ • Caching       │  │ • Time Series   │  │
│  │ • Asset Data    │  │ • Rate Limiting │  │ • Sensor Data   │  │
│  │ • RBAC Data     │  │ • Queues        │  │ • Analytics     │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 🔧 Service Architecture

### 1. Authentication Service

**Technology Stack**: Go + Gin + JWT + Argon2id
**Port**: 8081
**Purpose**: Tenant authentication and user management

**Key Components**:
- JWT token generation and validation
- Password hashing with Argon2id
- Multi-tenant user management
- SSO provider integration
- Session management

**Database Tables**:
- `users` - Tenant user accounts
- `tenants` - Tenant organizations
- `subscription_tiers` - Billing tiers
- `user_auth_methods` - Authentication methods
- `sso_providers` - SSO configuration

### 2. Inventory Service

**Technology Stack**: Go + Gin + PostgreSQL
**Port**: 8082
**Purpose**: Asset and sensor management

**Key Components**:
- Network asset management
- Sensor registration and monitoring
- Compliance tracking
- Report generation
- Real-time data processing

**Database Tables**:
- `network_assets` - Network devices and systems
- `sensors` - Monitoring sensors
- `certificates` - SSL/TLS certificates
- `compliance_assessments` - Compliance data
- `reports` - Generated reports

### 3. SaaS Admin Service

**Technology Stack**: Go + Gin + JWT + PostgreSQL
**Port**: 8084
**Purpose**: Platform administration and tenant management

**Key Components**:
- Platform-wide tenant management
- Cross-tenant user management
- Platform statistics and monitoring
- System health monitoring
- Platform-level RBAC

**Database Tables**:
- `platform_users` - Platform administrators
- `platform_roles` - Platform-level roles
- `platform_permissions` - Platform permissions
- `platform_role_permissions` - Role-permission mappings

### 4. Frontend Services

#### Tenant UI (React SPA)
**Technology Stack**: React 18 + TypeScript + Vite + TailwindCSS
**Port**: 3001
**Purpose**: Tenant-level asset and user management

**Key Features**:
- Asset management interface
- Sensor monitoring dashboard
- Report generation and viewing
- Tenant role management
- Real-time data updates

#### SaaS Admin UI (HTML/JS SPA)
**Technology Stack**: Vanilla HTML/JS + TailwindCSS + Axios
**Port**: 3002
**Purpose**: Platform administration interface

**Key Features**:
- Platform statistics dashboard
- Tenant management interface
- User management system
- System monitoring
- Real-time platform metrics

## 🗄️ Database Architecture

### Multi-Tenant Data Model

The system uses a **shared database, shared schema** approach with tenant isolation through application-level logic.

#### Tenant Isolation Strategy

1. **Tenant ID Filtering**: All tenant queries include `tenant_id` filters
2. **Row-Level Security**: Database-level tenant isolation (future enhancement)
3. **Application-Level Validation**: Service-level tenant access validation

#### Core Tables

**Authentication & Users**:
```sql
-- Core tenant management
tenants (id, name, slug, domain, subscription_tier_id, billing_email, payment_status)
subscription_tiers (id, name, description, price, features)
users (id, tenant_id, email, password_hash, first_name, last_name, role, is_active)

-- Platform administration
platform_users (id, email, password_hash, first_name, last_name, role_id, is_active)
platform_roles (id, name, display_name, description, is_system_role)
platform_permissions (id, name, resource, action, description)
platform_role_permissions (role_id, permission_id)
```

**Asset Management**:
```sql
-- Network assets
network_assets (id, tenant_id, name, type, ip_address, status, metadata)
crypto_implementations (id, tenant_id, name, version, type, status)
certificates (id, tenant_id, name, type, status, expiry_date)

-- Sensors and monitoring
sensors (id, tenant_id, name, type, status, last_heartbeat)
sensor_data (id, sensor_id, metric_name, value, timestamp)
```

**RBAC System**:
```sql
-- Tenant-level RBAC
tenant_roles (id, tenant_id, name, display_name, description, is_system_role)
tenant_permissions (id, tenant_id, name, resource, action, description)
tenant_role_permissions (role_id, permission_id)
user_roles (user_id, role_id)
```

## 🔐 Security Architecture

### Authentication Flow

1. **Tenant Users**:
   ```
   User → Tenant UI → Auth Service → JWT Token → Tenant Services
   ```

2. **Platform Admins**:
   ```
   Admin → SaaS Admin UI → SaaS Admin Service → JWT Token → Platform Services
   ```

### Authorization Model

**Role-Based Access Control (RBAC)**:
- **Tenant Roles**: `tenant_owner`, `tenant_admin`, `security_admin`, `analyst`
- **Platform Roles**: `super_admin`, `platform_admin`, `support_admin`

**Permission Model**:
- **Resource-Action Permissions**: `resource.action` (e.g., `assets.read`, `users.create`)
- **Scope-Based Access**: Tenant-scoped vs Platform-scoped permissions
- **Hierarchical Roles**: Role inheritance and permission aggregation

### Security Measures

1. **Password Security**: Argon2id hashing with salt
2. **JWT Security**: Short-lived access tokens with refresh tokens
3. **CORS Protection**: Configured for specific origins
4. **Input Validation**: Comprehensive request validation
5. **SQL Injection Prevention**: Parameterized queries
6. **Rate Limiting**: API rate limiting with Redis

## 📊 Data Flow Architecture

### Tenant Data Flow

```
Tenant User → Tenant UI → Auth Service → Inventory Service → PostgreSQL
     ↓              ↓           ↓              ↓
   JWT Token    React Query   JWT Validation   Data Query
     ↓              ↓           ↓              ↓
   Local Storage  State Mgmt   Role Check     Tenant Filter
```

### Platform Admin Data Flow

```
Platform Admin → SaaS Admin UI → SaaS Admin Service → PostgreSQL
       ↓              ↓              ↓
   JWT Token      Axios HTTP     JWT Validation
       ↓              ↓              ↓
   Local Storage   API Calls     Role Check
```

### Real-Time Data Flow

```
Sensors → InfluxDB → Inventory Service → WebSocket → Tenant UI
   ↓         ↓            ↓              ↓
Sensor Data  Time Series  Data Processing  Real-time Updates
```

## 🚀 Deployment Architecture

### Development Environment

```yaml
# docker-compose.yml
services:
  postgres:     # Database
  redis:        # Caching & Sessions
  influxdb:     # Time Series Data
  nats:         # Message Queue
  
  auth-service:        # Tenant Authentication
  inventory-service:   # Asset Management
  saas-admin-service:  # Platform Administration
  
  # Frontend (manual start)
  tenant-ui:    # npm run dev (port 3001)
  saas-admin-ui: # python3 -m http.server (port 3002)
```

### Production Environment

**Recommended Setup**:
- **Load Balancer**: Nginx or HAProxy
- **Application Servers**: Docker containers with orchestration
- **Database**: PostgreSQL with read replicas
- **Caching**: Redis cluster
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)

## 🔄 Integration Patterns

### Service-to-Service Communication

1. **Synchronous**: HTTP REST APIs
2. **Asynchronous**: NATS message queue
3. **Data Sharing**: Shared PostgreSQL database
4. **Caching**: Redis for session and data caching

### Frontend Integration

1. **API Communication**: Axios HTTP client
2. **State Management**: React Query for server state
3. **Real-time Updates**: WebSocket connections
4. **Authentication**: JWT token management

## 📈 Scalability Considerations

### Horizontal Scaling

1. **Stateless Services**: All services are stateless
2. **Database Sharding**: Future tenant-based sharding
3. **Load Balancing**: Multiple service instances
4. **CDN Integration**: Static asset delivery

### Performance Optimization

1. **Database Indexing**: Optimized queries with proper indexes
2. **Caching Strategy**: Multi-level caching (Redis, CDN)
3. **Connection Pooling**: Database connection optimization
4. **Query Optimization**: Efficient database queries

## 🔍 Monitoring & Observability

### Health Checks

- **Service Health**: `/health` endpoints for all services
- **Database Health**: Connection and query performance monitoring
- **Dependency Health**: External service availability

### Metrics Collection

- **Application Metrics**: Request rates, response times, error rates
- **Business Metrics**: User activity, tenant usage, feature adoption
- **Infrastructure Metrics**: CPU, memory, disk, network usage

### Logging Strategy

- **Structured Logging**: JSON-formatted logs
- **Log Aggregation**: Centralized log collection
- **Log Analysis**: Real-time log monitoring and alerting

## 🛠️ Development Workflow

### Code Organization

```
/
├── services/           # Backend services
│   ├── auth-service/
│   ├── inventory-service/
│   └── saas-admin-service/
├── web-ui/            # Tenant React application
├── saas-admin-ui/     # Platform admin interface
├── scripts/           # Database scripts and migrations
├── docs/              # Documentation
└── docker-compose.yml # Development environment
```

### Development Process

1. **Local Development**: Docker Compose for services
2. **Frontend Development**: Vite dev servers
3. **Database Migrations**: SQL scripts in `/scripts/database/`
4. **Testing**: Unit and integration tests
5. **Documentation**: API docs and architecture docs

---

*This architecture documentation should be updated as the system evolves. Last updated: 2025-01-09*
