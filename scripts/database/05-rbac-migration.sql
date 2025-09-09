-- =================================================================
-- RBAC (Role-Based Access Control) Migration
-- Version: 1.0.0
-- Date: 2024-01-XX
-- =================================================================
-- This migration adds comprehensive RBAC system to support both
-- platform-level (SaaS) and tenant-level administration

-- =================================================================
-- Platform-Level Roles and Permissions
-- =================================================================

-- Platform roles (SaaS administrators)
CREATE TABLE IF NOT EXISTS platform_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Platform permissions
CREATE TABLE IF NOT EXISTS platform_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL, -- 'tenants', 'users', 'platform', 'billing'
    action VARCHAR(50) NOT NULL,   -- 'create', 'read', 'update', 'delete', 'manage'
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Platform role-permission mappings
CREATE TABLE IF NOT EXISTS platform_role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES platform_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES platform_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

-- Platform users (SaaS administrators)
CREATE TABLE IF NOT EXISTS platform_users (
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

-- =================================================================
-- Tenant-Level Roles and Permissions
-- =================================================================

-- Tenant roles (Customer administrators)
CREATE TABLE IF NOT EXISTS tenant_roles (
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
CREATE TABLE IF NOT EXISTS tenant_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL, -- 'assets', 'sensors', 'reports', 'users', 'settings'
    action VARCHAR(50) NOT NULL,   -- 'create', 'read', 'update', 'delete', 'manage'
    scope VARCHAR(50) DEFAULT 'tenant', -- 'tenant', 'own', 'team'
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tenant role-permission mappings
CREATE TABLE IF NOT EXISTS tenant_role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES tenant_roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES tenant_permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(role_id, permission_id)
);

-- User role assignments (replaces simple role field in users table)
CREATE TABLE IF NOT EXISTS user_tenant_roles (
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

-- =================================================================
-- Resource-Level Permissions
-- =================================================================

-- Resource ownership and access control
CREATE TABLE IF NOT EXISTS resource_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_type VARCHAR(50) NOT NULL, -- 'network_asset', 'sensor', 'report'
    resource_id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permissions JSONB DEFAULT '{}', -- Custom permissions per resource
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =================================================================
-- Audit and Compliance
-- =================================================================

-- Permission audit logs
CREATE TABLE IF NOT EXISTS permission_audit_logs (
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

-- =================================================================
-- Indexes for Performance
-- =================================================================

-- Platform-level indexes
CREATE INDEX IF NOT EXISTS idx_platform_users_email ON platform_users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_platform_users_role ON platform_users(role_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_platform_role_permissions_role ON platform_role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_platform_role_permissions_permission ON platform_role_permissions(permission_id);

-- Tenant-level indexes
CREATE INDEX IF NOT EXISTS idx_tenant_roles_tenant ON tenant_roles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_role_permissions_role ON tenant_role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_tenant_role_permissions_permission ON tenant_role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_user_tenant_roles_user ON user_tenant_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_tenant_roles_tenant ON user_tenant_roles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_tenant_roles_role ON user_tenant_roles(role_id);

-- Resource permission indexes
CREATE INDEX IF NOT EXISTS idx_resource_permissions_type_id ON resource_permissions(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_resource_permissions_tenant ON resource_permissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_resource_permissions_owner ON resource_permissions(owner_id);

-- Audit log indexes
CREATE INDEX IF NOT EXISTS idx_permission_audit_user ON permission_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_permission_audit_tenant ON permission_audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_permission_audit_created ON permission_audit_logs(created_at);

-- =================================================================
-- Seed Data for Default Roles and Permissions
-- =================================================================

-- Platform roles
INSERT INTO platform_roles (name, display_name, description, is_system_role) VALUES
('super_admin', 'Super Administrator', 'Full platform access, can manage all tenants and platform settings', true),
('platform_admin', 'Platform Administrator', 'Platform management with limited tenant access', true),
('support_admin', 'Support Administrator', 'Customer support access, read-only platform data', true)
ON CONFLICT (name) DO NOTHING;

-- Platform permissions
INSERT INTO platform_permissions (name, resource, action, description) VALUES
-- Tenant management
('tenants.create', 'tenants', 'create', 'Create new tenants'),
('tenants.read', 'tenants', 'read', 'View tenant information'),
('tenants.update', 'tenants', 'update', 'Update tenant settings'),
('tenants.delete', 'tenants', 'delete', 'Delete tenants'),
('tenants.manage', 'tenants', 'manage', 'Full tenant management'),

-- Platform user management
('platform_users.create', 'platform_users', 'create', 'Create platform users'),
('platform_users.read', 'platform_users', 'read', 'View platform users'),
('platform_users.update', 'platform_users', 'update', 'Update platform users'),
('platform_users.delete', 'platform_users', 'delete', 'Delete platform users'),

-- Platform settings
('platform.settings', 'platform', 'manage', 'Manage platform settings'),
('platform.billing', 'platform', 'manage', 'Manage platform billing'),
('platform.analytics', 'platform', 'read', 'View platform analytics'),

-- Support access
('support.tenants', 'tenants', 'read', 'View tenant data for support'),
('support.users', 'users', 'read', 'View user data for support')
ON CONFLICT (name) DO NOTHING;

-- Assign permissions to platform roles
-- Super Admin gets all permissions
INSERT INTO platform_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM platform_roles r, platform_permissions p
WHERE r.name = 'super_admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Platform Admin gets most permissions except super admin functions
INSERT INTO platform_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM platform_roles r, platform_permissions p
WHERE r.name = 'platform_admin' 
  AND p.name NOT IN ('platform_users.delete', 'tenants.delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Support Admin gets read-only access
INSERT INTO platform_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM platform_roles r, platform_permissions p
WHERE r.name = 'support_admin' 
  AND p.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Tenant permissions
INSERT INTO tenant_permissions (name, resource, action, scope, description) VALUES
-- Asset management
('assets.create', 'assets', 'create', 'tenant', 'Create network assets'),
('assets.read', 'assets', 'read', 'tenant', 'View network assets'),
('assets.update', 'assets', 'update', 'tenant', 'Update network assets'),
('assets.delete', 'assets', 'delete', 'tenant', 'Delete network assets'),
('assets.manage', 'assets', 'manage', 'tenant', 'Full asset management'),

-- Sensor management
('sensors.create', 'sensors', 'create', 'tenant', 'Create sensors'),
('sensors.read', 'sensors', 'read', 'tenant', 'View sensors'),
('sensors.update', 'sensors', 'update', 'tenant', 'Update sensors'),
('sensors.delete', 'sensors', 'delete', 'tenant', 'Delete sensors'),
('sensors.manage', 'sensors', 'manage', 'tenant', 'Full sensor management'),

-- Report management
('reports.create', 'reports', 'create', 'tenant', 'Create reports'),
('reports.read', 'reports', 'read', 'tenant', 'View reports'),
('reports.update', 'reports', 'update', 'tenant', 'Update reports'),
('reports.delete', 'reports', 'delete', 'tenant', 'Delete reports'),
('reports.manage', 'reports', 'manage', 'tenant', 'Full report management'),

-- User management
('users.create', 'users', 'create', 'tenant', 'Create tenant users'),
('users.read', 'users', 'read', 'tenant', 'View tenant users'),
('users.update', 'users', 'update', 'tenant', 'Update tenant users'),
('users.delete', 'users', 'delete', 'tenant', 'Delete tenant users'),
('users.manage', 'users', 'manage', 'tenant', 'Full user management'),

-- Settings management
('settings.read', 'settings', 'read', 'tenant', 'View tenant settings'),
('settings.update', 'settings', 'update', 'tenant', 'Update tenant settings'),
('settings.manage', 'settings', 'manage', 'tenant', 'Full settings management'),

-- Billing management
('billing.read', 'billing', 'read', 'tenant', 'View billing information'),
('billing.update', 'billing', 'update', 'tenant', 'Update billing settings'),

-- Compliance management
('compliance.read', 'compliance', 'read', 'tenant', 'View compliance data'),
('compliance.update', 'compliance', 'update', 'tenant', 'Update compliance settings'),
('compliance.manage', 'compliance', 'manage', 'tenant', 'Full compliance management')
ON CONFLICT (name) DO NOTHING;

-- =================================================================
-- Functions for Permission Checking
-- =================================================================

-- Function to check if a user has a specific permission
CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_tenant_id UUID,
    p_permission_name VARCHAR(100)
) RETURNS BOOLEAN AS $$
DECLARE
    has_permission BOOLEAN := FALSE;
BEGIN
    -- Check if user has the permission through their tenant role
    SELECT EXISTS(
        SELECT 1
        FROM user_tenant_roles utr
        JOIN tenant_role_permissions trp ON utr.role_id = trp.role_id
        JOIN tenant_permissions tp ON trp.permission_id = tp.id
        WHERE utr.user_id = p_user_id
          AND utr.tenant_id = p_tenant_id
          AND utr.is_active = true
          AND (utr.expires_at IS NULL OR utr.expires_at > NOW())
          AND tp.name = p_permission_name
    ) INTO has_permission;
    
    RETURN has_permission;
END;
$$ LANGUAGE plpgsql;

-- Function to get all permissions for a user in a tenant
CREATE OR REPLACE FUNCTION get_user_permissions(
    p_user_id UUID,
    p_tenant_id UUID
) RETURNS TABLE(permission_name VARCHAR(100), resource VARCHAR(50), action VARCHAR(50), scope VARCHAR(50)) AS $$
BEGIN
    RETURN QUERY
    SELECT tp.name, tp.resource, tp.action, tp.scope
    FROM user_tenant_roles utr
    JOIN tenant_role_permissions trp ON utr.role_id = trp.role_id
    JOIN tenant_permissions tp ON trp.permission_id = tp.id
    WHERE utr.user_id = p_user_id
      AND utr.tenant_id = p_tenant_id
      AND utr.is_active = true
      AND (utr.expires_at IS NULL OR utr.expires_at > NOW());
END;
$$ LANGUAGE plpgsql;

-- Function to check platform permissions
CREATE OR REPLACE FUNCTION platform_user_has_permission(
    p_user_id UUID,
    p_permission_name VARCHAR(100)
) RETURNS BOOLEAN AS $$
DECLARE
    has_permission BOOLEAN := FALSE;
BEGIN
    -- Check if platform user has the permission through their role
    SELECT EXISTS(
        SELECT 1
        FROM platform_users pu
        JOIN platform_role_permissions prp ON pu.role_id = prp.role_id
        JOIN platform_permissions pp ON prp.permission_id = pp.id
        WHERE pu.id = p_user_id
          AND pu.is_active = true
          AND pu.deleted_at IS NULL
          AND pp.name = p_permission_name
    ) INTO has_permission;
    
    RETURN has_permission;
END;
$$ LANGUAGE plpgsql;

-- =================================================================
-- Views for Easy Permission Management
-- =================================================================

-- View of all user permissions in a tenant
CREATE OR REPLACE VIEW user_tenant_permissions AS
SELECT 
    u.id as user_id,
    u.email,
    u.first_name,
    u.last_name,
    t.id as tenant_id,
    t.name as tenant_name,
    tr.name as role_name,
    tr.display_name as role_display_name,
    tp.name as permission_name,
    tp.resource,
    tp.action,
    tp.scope,
    utr.assigned_at,
    utr.expires_at,
    utr.is_active
FROM users u
JOIN user_tenant_roles utr ON u.id = utr.user_id
JOIN tenants t ON utr.tenant_id = t.id
JOIN tenant_roles tr ON utr.role_id = tr.id
JOIN tenant_role_permissions trp ON tr.id = trp.role_id
JOIN tenant_permissions tp ON trp.permission_id = tp.id
WHERE u.deleted_at IS NULL
  AND utr.is_active = true
  AND (utr.expires_at IS NULL OR utr.expires_at > NOW());

-- View of platform administrators
CREATE OR REPLACE VIEW platform_administrators AS
SELECT 
    pu.id,
    pu.email,
    pu.first_name,
    pu.last_name,
    pr.name as role_name,
    pr.display_name as role_display_name,
    pu.is_active,
    pu.last_login_at,
    pu.created_at
FROM platform_users pu
JOIN platform_roles pr ON pu.role_id = pr.id
WHERE pu.deleted_at IS NULL;

-- =================================================================
-- Triggers for Automatic Updates
-- =================================================================

-- Update timestamps for RBAC tables
CREATE TRIGGER update_platform_roles_updated_at BEFORE UPDATE ON platform_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenant_roles_updated_at BEFORE UPDATE ON tenant_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resource_permissions_updated_at BEFORE UPDATE ON resource_permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =================================================================
-- Migration Complete
-- =================================================================
