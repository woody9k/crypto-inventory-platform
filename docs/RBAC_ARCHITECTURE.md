# Multi-Tenant RBAC (Role-Based Access Control) Architecture

## Overview

This document outlines the comprehensive Role-Based Access Control (RBAC) system for the Crypto Inventory Platform, designed to support both SaaS-level administration and tenant-level management in a multi-tenant environment.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Role Hierarchy](#role-hierarchy)
3. [Permission System](#permission-system)
4. [Database Schema](#database-schema)
5. [API Design](#api-design)
6. [Frontend Components](#frontend-components)
7. [Security Considerations](#security-considerations)
8. [Implementation Phases](#implementation-phases)
9. [Migration Strategy](#migration-strategy)

## Architecture Overview

### Multi-Level RBAC Design

The RBAC system operates at two distinct levels:

1. **Platform Level (SaaS Administrators)**
   - Manage the entire platform
   - Oversee all tenants and their data
   - Handle billing, support, and platform operations

2. **Tenant Level (Customer Administrators)**
   - Manage their own tenant's resources
   - Control user access within their organization
   - Configure tenant-specific settings

### Key Principles

- **Principle of Least Privilege**: Users receive minimum permissions necessary
- **Separation of Concerns**: Clear distinction between platform and tenant administration
- **Audit Trail**: All permission changes and access attempts are logged
- **Scalability**: System supports thousands of tenants and users
- **Flexibility**: Custom roles and permissions can be created per tenant

## Role Hierarchy

### Platform Level Roles

| Role | Display Name | Description | Access Level |
|------|-------------|-------------|--------------|
| `super_admin` | Super Administrator | Full platform access, all tenants | Global |
| `platform_admin` | Platform Administrator | Platform management, limited tenant access | Platform + Limited Tenant |
| `support_admin` | Support Administrator | Customer support, read-only platform access | Read-Only Platform |

### Tenant Level Roles

| Role | Display Name | Description | Access Level |
|------|-------------|-------------|--------------|
| `tenant_owner` | Tenant Owner | Full tenant control, billing, user management | Full Tenant |
| `tenant_admin` | Tenant Administrator | Tenant management, user management | Full Tenant (No Billing) |
| `security_admin` | Security Administrator | Security settings, compliance, reports | Security + Compliance |
| `analyst` | Security Analyst | Data analysis, reporting, monitoring | Analysis + Reporting |
| `viewer` | Viewer | Read-only access to tenant data | Read-Only |
| `api_user` | API User | API-only access for integrations | API Access |

## Permission System

### Permission Structure

Each permission follows the pattern: `{resource}.{action}` with optional scope.

- **Resource**: The entity being accessed (e.g., `assets`, `sensors`, `users`)
- **Action**: The operation being performed (e.g., `create`, `read`, `update`, `delete`, `manage`)
- **Scope**: The extent of access (e.g., `tenant`, `own`, `team`)

### Platform Permissions

```typescript
const PLATFORM_PERMISSIONS = {
  // Tenant Management
  'tenants.create': { resource: 'tenants', action: 'create' },
  'tenants.read': { resource: 'tenants', action: 'read' },
  'tenants.update': { resource: 'tenants', action: 'update' },
  'tenants.delete': { resource: 'tenants', action: 'delete' },
  'tenants.manage': { resource: 'tenants', action: 'manage' },

  // Platform User Management
  'platform_users.create': { resource: 'platform_users', action: 'create' },
  'platform_users.read': { resource: 'platform_users', action: 'read' },
  'platform_users.update': { resource: 'platform_users', action: 'update' },
  'platform_users.delete': { resource: 'platform_users', action: 'delete' },

  // Platform Settings
  'platform.settings': { resource: 'platform', action: 'manage' },
  'platform.billing': { resource: 'platform', action: 'manage' },
  'platform.analytics': { resource: 'platform', action: 'read' },

  // Support Access
  'support.tenants': { resource: 'tenants', action: 'read' },
  'support.users': { resource: 'users', action: 'read' }
};
```

### Tenant Permissions

```typescript
const TENANT_PERMISSIONS = {
  // Asset Management
  'assets.create': { resource: 'assets', action: 'create', scope: 'tenant' },
  'assets.read': { resource: 'assets', action: 'read', scope: 'tenant' },
  'assets.update': { resource: 'assets', action: 'update', scope: 'tenant' },
  'assets.delete': { resource: 'assets', action: 'delete', scope: 'tenant' },
  'assets.manage': { resource: 'assets', action: 'manage', scope: 'tenant' },

  // Sensor Management
  'sensors.create': { resource: 'sensors', action: 'create', scope: 'tenant' },
  'sensors.read': { resource: 'sensors', action: 'read', scope: 'tenant' },
  'sensors.update': { resource: 'sensors', action: 'update', scope: 'tenant' },
  'sensors.delete': { resource: 'sensors', action: 'delete', scope: 'tenant' },
  'sensors.manage': { resource: 'sensors', action: 'manage', scope: 'tenant' },

  // Report Management
  'reports.create': { resource: 'reports', action: 'create', scope: 'tenant' },
  'reports.read': { resource: 'reports', action: 'read', scope: 'tenant' },
  'reports.update': { resource: 'reports', action: 'update', scope: 'tenant' },
  'reports.delete': { resource: 'reports', action: 'delete', scope: 'tenant' },
  'reports.manage': { resource: 'reports', action: 'manage', scope: 'tenant' },

  // User Management
  'users.create': { resource: 'users', action: 'create', scope: 'tenant' },
  'users.read': { resource: 'users', action: 'read', scope: 'tenant' },
  'users.update': { resource: 'users', action: 'update', scope: 'tenant' },
  'users.delete': { resource: 'users', action: 'delete', scope: 'tenant' },
  'users.manage': { resource: 'users', action: 'manage', scope: 'tenant' },

  // Settings Management
  'settings.read': { resource: 'settings', action: 'read', scope: 'tenant' },
  'settings.update': { resource: 'settings', action: 'update', scope: 'tenant' },
  'settings.manage': { resource: 'settings', action: 'manage', scope: 'tenant' },

  // Billing Management
  'billing.read': { resource: 'billing', action: 'read', scope: 'tenant' },
  'billing.update': { resource: 'billing', action: 'update', scope: 'tenant' },

  // Compliance Management
  'compliance.read': { resource: 'compliance', action: 'read', scope: 'tenant' },
  'compliance.update': { resource: 'compliance', action: 'update', scope: 'tenant' },
  'compliance.manage': { resource: 'compliance', action: 'manage', scope: 'tenant' }
};
```

## Database Schema

### Core Tables

#### Platform Level

```sql
-- Platform roles (SaaS administrators)
CREATE TABLE platform_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Platform permissions
CREATE TABLE platform_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Platform role-permission mappings
CREATE TABLE platform_role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES platform_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES platform_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

-- Platform users (SaaS administrators)
CREATE TABLE platform_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role_id UUID NOT NULL REFERENCES platform_roles(id),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

#### Tenant Level

```sql
-- Tenant roles (Customer administrators)
CREATE TABLE tenant_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Tenant permissions
CREATE TABLE tenant_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    scope VARCHAR(50) DEFAULT 'tenant',
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tenant role-permission mappings
CREATE TABLE tenant_role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES tenant_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES tenant_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

-- User role assignments (replaces simple role field in users table)
CREATE TABLE user_tenant_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES tenant_roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(user_id, tenant_id, role_id)
);
```

### Audit and Compliance

```sql
-- Permission audit logs
CREATE TABLE permission_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    tenant_id UUID REFERENCES tenants(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    permission_required VARCHAR(100),
    permission_granted BOOLEAN,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Helper Functions

```sql
-- Function to check if a user has a specific permission
CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_tenant_id UUID,
    p_permission_name VARCHAR(100)
) RETURNS BOOLEAN;

-- Function to get all permissions for a user in a tenant
CREATE OR REPLACE FUNCTION get_user_permissions(
    p_user_id UUID,
    p_tenant_id UUID
) RETURNS TABLE(permission_name VARCHAR(100), resource VARCHAR(50), action VARCHAR(50), scope VARCHAR(50));
```

## API Design

### Platform Administration APIs

```typescript
// Platform User Management
GET    /api/v1/platform/users              // List platform users
POST   /api/v1/platform/users              // Create platform user
GET    /api/v1/platform/users/:id          // Get platform user
PUT    /api/v1/platform/users/:id          // Update platform user
DELETE /api/v1/platform/users/:id          // Delete platform user

// Tenant Management
GET    /api/v1/platform/tenants            // List all tenants
POST   /api/v1/platform/tenants            // Create new tenant
GET    /api/v1/platform/tenants/:id        // Get tenant details
PUT    /api/v1/platform/tenants/:id        // Update tenant
DELETE /api/v1/platform/tenants/:id        // Delete tenant

// Platform Analytics
GET    /api/v1/platform/analytics/overview // Platform overview stats
GET    /api/v1/platform/analytics/tenants  // Tenant usage analytics
GET    /api/v1/platform/analytics/users    // User activity analytics
```

### Tenant Administration APIs

```typescript
// Tenant User Management
GET    /api/v1/tenant/users                // List tenant users
POST   /api/v1/tenant/users                // Create tenant user
GET    /api/v1/tenant/users/:id            // Get tenant user
PUT    /api/v1/tenant/users/:id            // Update tenant user
DELETE /api/v1/tenant/users/:id            // Delete tenant user

// Role Management
GET    /api/v1/tenant/roles                // List tenant roles
POST   /api/v1/tenant/roles                // Create custom role
GET    /api/v1/tenant/roles/:id            // Get role details
PUT    /api/v1/tenant/roles/:id            // Update role
DELETE /api/v1/tenant/roles/:id            // Delete role

// Permission Management
GET    /api/v1/tenant/permissions          // List available permissions
GET    /api/v1/tenant/users/:id/permissions // Get user permissions
PUT    /api/v1/tenant/users/:id/permissions // Update user permissions
```

### Permission Checking Middleware

```go
// Enhanced middleware for permission checking
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("userID")
        tenantID := c.GetString("tenantID")
        
        hasPermission, err := checkUserPermission(userID, tenantID, permission)
        if err != nil || !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Insufficient permissions",
                "required_permission": permission,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## Frontend Components

### Platform Administration Interface

```typescript
// Platform Admin Dashboard
interface PlatformAdminDashboard {
  // Tenant overview with key metrics
  tenantOverview: TenantOverview;
  
  // Platform-wide analytics
  platformAnalytics: PlatformAnalytics;
  
  // Recent activity across all tenants
  recentActivity: ActivityLog[];
}

// Tenant Management
interface TenantManagement {
  // List all tenants with search/filter
  tenantList: Tenant[];
  
  // Tenant details and settings
  tenantDetails: TenantDetails;
  
  // Tenant user management
  tenantUsers: User[];
}
```

### Tenant Administration Interface

```typescript
// User Management
interface UserManagement {
  // List tenant users with roles
  users: UserWithRole[];
  
  // Role assignment interface
  roleAssignment: RoleAssignment;
  
  // Permission management
  permissionManagement: PermissionManagement;
}

// Role Management
interface RoleManagement {
  // System and custom roles
  roles: Role[];
  
  // Permission matrix
  permissionMatrix: PermissionMatrix;
  
  // Role creation/editing
  roleEditor: RoleEditor;
}
```

### Permission-Aware Components

```typescript
// Conditional rendering based on permissions
const PermissionGate: React.FC<{
  permission: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}> = ({ permission, children, fallback = null }) => {
  const { hasPermission } = usePermissions();
  
  if (!hasPermission(permission)) {
    return <>{fallback}</>;
  }
  
  return <>{children}</>;
};

// Usage example
<PermissionGate permission="sensors.create">
  <CreateSensorButton />
</PermissionGate>
```

## Security Considerations

### Authentication & Authorization

1. **JWT Token Enhancement**
   - Include role and permission information in tokens
   - Implement token refresh with permission updates
   - Short-lived access tokens (15 minutes)

2. **Permission Caching**
   - Cache user permissions in Redis
   - Invalidate cache on role/permission changes
   - TTL-based cache expiration

3. **Audit Logging**
   - Log all permission checks
   - Track role changes and assignments
   - Monitor suspicious access patterns

### Data Isolation

1. **Tenant Isolation**
   - All queries must include tenant_id filter
   - Row-level security policies
   - Cross-tenant data access prevention

2. **Resource-Level Permissions**
   - Fine-grained access control per resource
   - Owner-based permissions
   - Team-based access controls

### API Security

1. **Rate Limiting**
   - Per-user and per-tenant rate limits
   - Different limits for different roles
   - API key management for integrations

2. **Input Validation**
   - Strict validation of role/permission assignments
   - Prevention of privilege escalation
   - SQL injection prevention

## Implementation Phases

### Phase 1: Database Schema (Week 1)
- [ ] Create RBAC database tables
- [ ] Implement helper functions
- [ ] Create seed data for default roles
- [ ] Set up audit logging

### Phase 2: Backend Services (Week 2)
- [ ] Enhance authentication service
- [ ] Implement permission checking middleware
- [ ] Create platform administration APIs
- [ ] Add tenant role management APIs

### Phase 3: Frontend Components (Week 3)
- [ ] Build platform admin interface
- [ ] Create tenant user management UI
- [ ] Implement role assignment interface
- [ ] Add permission-aware components

### Phase 4: Integration & Testing (Week 4)
- [ ] Integrate RBAC with existing services
- [ ] Comprehensive testing of all roles
- [ ] Performance optimization
- [ ] Security audit

### Phase 5: Migration & Deployment (Week 5)
- [ ] Migrate existing users to new role system
- [ ] Deploy to staging environment
- [ ] User acceptance testing
- [ ] Production deployment

## Migration Strategy

### Existing User Migration

1. **Current Role Mapping**
   ```sql
   -- Map existing roles to new system
   INSERT INTO tenant_roles (tenant_id, name, display_name, is_system_role)
   SELECT DISTINCT 
     tenant_id,
     CASE 
       WHEN role = 'admin' THEN 'tenant_admin'
       WHEN role = 'analyst' THEN 'analyst'
       WHEN role = 'viewer' THEN 'viewer'
     END,
     CASE 
       WHEN role = 'admin' THEN 'Tenant Administrator'
       WHEN role = 'analyst' THEN 'Security Analyst'
       WHEN role = 'viewer' THEN 'Viewer'
     END,
     true
   FROM users;
   ```

2. **User Role Assignment**
   ```sql
   -- Assign users to their new roles
   INSERT INTO user_tenant_roles (user_id, tenant_id, role_id)
   SELECT 
     u.id,
     u.tenant_id,
     tr.id
   FROM users u
   JOIN tenant_roles tr ON tr.tenant_id = u.tenant_id
   WHERE tr.name = CASE 
     WHEN u.role = 'admin' THEN 'tenant_admin'
     WHEN u.role = 'analyst' THEN 'analyst'
     WHEN u.role = 'viewer' THEN 'viewer'
   END;
   ```

### Backward Compatibility

1. **API Compatibility**
   - Maintain existing API endpoints
   - Add new RBAC endpoints alongside
   - Gradual migration of frontend components

2. **Database Compatibility**
   - Keep existing `role` field in users table
   - Populate from new RBAC system
   - Remove after full migration

## Success Metrics

### Security Metrics
- Zero unauthorized access incidents
- 100% audit trail coverage
- Permission check response time < 50ms

### Usability Metrics
- Role assignment completion time < 2 minutes
- User permission confusion rate < 5%
- Admin task completion time improvement > 30%

### Performance Metrics
- Permission check overhead < 10ms
- Database query performance maintained
- Cache hit rate > 95%

## Conclusion

This RBAC architecture provides a comprehensive, scalable, and secure foundation for multi-tenant administration. The two-level design ensures clear separation between platform and tenant management while maintaining flexibility for custom roles and permissions.

The implementation follows security best practices and provides a clear migration path from the existing simple role system to the full RBAC implementation.
