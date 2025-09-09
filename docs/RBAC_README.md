# RBAC System Implementation Guide

## 🎯 Overview

The Role-Based Access Control (RBAC) system provides comprehensive permission management for the crypto inventory management platform. This system supports both platform-level (SaaS) administrators and tenant-level (customer) administrators with granular permission control.

## ✅ Implementation Status

- **Backend**: ✅ Complete with Go services, middleware, and API endpoints
- **Database**: ✅ Complete with 12+ tables and comprehensive schema
- **Frontend**: ✅ Complete with React components and role management UI
- **Documentation**: ✅ Complete with architecture docs and implementation guides

## 🏗️ Architecture

### Backend Components

#### Database Schema
- **12+ Tables**: Roles, permissions, user assignments, audit logs
- **Platform Roles**: SaaS-level administration (super_admin, platform_admin, support_admin)
- **Tenant Roles**: Customer-level administration (tenant_owner, tenant_admin, security_admin, analyst, viewer, api_user)
- **Permissions**: Granular resource/action/scope based permissions
- **Audit Logging**: Complete trail of all permission checks and role changes

#### Go Services
- **RBAC Service**: Core permission checking and role management logic
- **Middleware**: Authentication and permission enforcement
- **API Handlers**: RESTful endpoints for role and user management
- **Database Integration**: PostgreSQL with optimized queries and indexing

#### API Endpoints
```bash
# Role Management
GET    /api/v1/tenant/:tenantId/roles                    # List tenant roles
PUT    /api/v1/tenant/:tenantId/roles/:roleId/permissions # Update role permissions

# User Management  
GET    /api/v1/tenant/:tenantId/users/:userId/roles      # List user roles
POST   /api/v1/tenant/:tenantId/users/:userId/roles      # Assign role to user
DELETE /api/v1/tenant/:tenantId/users/:userId/roles/:roleId # Remove role

# Permissions
GET    /api/v1/permissions                               # List all permissions
POST   /api/v1/permissions/check                         # Check user permission

# Platform Administration
GET    /api/v1/platform/users                            # List platform users
GET    /api/v1/platform/roles                            # List platform roles

# Audit & Monitoring
GET    /api/v1/audit/logs                                # View audit logs
```

### Frontend Components

#### Pages
- **Role Management Page** (`/roles`): Complete RBAC interface with tabs for:
  - **Roles Tab**: Role management and permission matrix
  - **Users Tab**: User role assignments and management
  - **Permissions Tab**: Permission overview by category
  - **Audit Tab**: Permission check logs and audit trail

#### Components
- **PermissionGate**: Context-based permission management with conditional rendering
- **RoleManagement**: Role creation, editing, and permission assignment
- **UserManagement**: User lifecycle management and role assignments
- **PermissionProvider**: React context for permission state management

## 🚀 Getting Started

### 1. Start the Services
```bash
# Start all services
docker-compose up -d

# Start frontend separately
cd web-ui && npm run dev
```

### 2. Access the System
- **Main Application**: http://localhost:3000
- **Role Management**: http://localhost:3000/roles
- **API Documentation**: Available at backend service endpoints

### 3. Default Roles and Permissions

#### Platform Roles
- **super_admin**: Full platform access
- **platform_admin**: Platform management capabilities
- **support_admin**: Customer support and monitoring

#### Tenant Roles
- **tenant_owner**: Full tenant management
- **tenant_admin**: Tenant administration
- **security_admin**: Security and compliance
- **analyst**: Data analysis and reporting
- **viewer**: Read-only access
- **api_user**: API access only

## 🔧 Configuration

### Database Setup
The RBAC schema is automatically created when the PostgreSQL container starts. The migration scripts are located in:
- `scripts/database/05-rbac-migration.sql` - Schema creation
- `scripts/database/06-rbac-seed.sql` - Default data seeding

### Environment Variables
```bash
# Database
POSTGRES_DB=crypto_inventory
POSTGRES_USER=crypto_user
POSTGRES_PASSWORD=your_password

# Auth Service
AUTH_SERVICE_PORT=8081
JWT_SECRET=your_jwt_secret
```

## 📊 Permission System

### Permission Structure
Permissions follow the format: `resource.action.scope`

#### Resources
- **assets**: Network asset management
- **sensors**: Sensor management
- **reports**: Report generation and viewing
- **users**: User management
- **settings**: Tenant settings

#### Actions
- **create**: Create new resources
- **read**: View resources
- **update**: Modify resources
- **delete**: Remove resources
- **manage**: Full management capabilities

#### Scopes
- **platform**: Platform-wide access
- **tenant**: Tenant-specific access

### Example Permissions
```bash
assets.read.tenant          # View tenant assets
sensors.manage.tenant       # Full sensor management
users.create.tenant         # Create tenant users
settings.update.tenant      # Update tenant settings
```

## 🔒 Security Features

### Authentication
- JWT-based token authentication
- Token validation on all protected endpoints
- Secure password hashing with bcrypt

### Authorization
- Granular permission checking
- Role-based access control
- Tenant isolation for multi-tenancy
- Permission inheritance and hierarchy

### Audit & Monitoring
- Complete audit trail of all permission checks
- Role assignment and removal logging
- User activity tracking
- Security event monitoring

## 🧪 Testing

### Backend Testing
```bash
# Test API endpoints
curl -H "Authorization: Bearer <token>" http://localhost:8081/api/v1/permissions

# Test permission checking
curl -X POST http://localhost:8081/api/v1/permissions/check \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"permission": "assets.read", "resource": "assets", "action": "read"}'
```

### Frontend Testing
1. Navigate to http://localhost:3000/roles
2. Test role management interface
3. Verify permission-based conditional rendering
4. Check user role assignment functionality

## 📝 Development Notes

### Adding New Permissions
1. Add permission to database seed script
2. Update frontend permission lists
3. Add permission checks to relevant components
4. Update API endpoints as needed

### Adding New Roles
1. Add role to database seed script
2. Define role permissions
3. Update frontend role management
4. Test role assignment functionality

### Customizing UI
- All components support dark mode
- Responsive design for mobile/desktop
- Customizable permission matrix
- Extensible role management interface

## 🐛 Troubleshooting

### Common Issues
1. **Blank Roles Page**: Check React Router configuration and component imports
2. **Permission Denied**: Verify user has correct role and permissions
3. **API Errors**: Check authentication token and endpoint availability
4. **Database Issues**: Verify PostgreSQL is running and schema is created

### Debug Mode
Enable debug logging by setting environment variables:
```bash
DEBUG=true
LOG_LEVEL=debug
```

## 📚 Additional Resources

- [RBAC Architecture Documentation](./RBAC_ARCHITECTURE.md)
- [API Documentation](./API_DOCUMENTATION.md)
- [Database Schema](./DATABASE_SCHEMA.md)
- [Frontend Component Guide](./FRONTEND_COMPONENTS.md)

## 🤝 Contributing

When adding new features to the RBAC system:
1. Update database schema if needed
2. Add corresponding Go models and services
3. Update frontend components
4. Add comprehensive tests
5. Update documentation

## 📄 License

This RBAC system is part of the Crypto Inventory Management Platform and follows the same licensing terms.
