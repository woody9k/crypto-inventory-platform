-- =================================================================
-- RBAC Seed Data
-- =================================================================
-- This file creates default tenant roles and assigns them to existing users

-- =================================================================
-- Create Default Tenant Roles for Each Tenant
-- =================================================================

-- Create system roles for each existing tenant
INSERT INTO tenant_roles (tenant_id, name, display_name, description, is_system_role)
SELECT 
    t.id,
    'tenant_owner',
    'Tenant Owner',
    'Full tenant control, billing, user management',
    true
FROM tenants t
ON CONFLICT (tenant_id, name) DO NOTHING;

INSERT INTO tenant_roles (tenant_id, name, display_name, description, is_system_role)
SELECT 
    t.id,
    'tenant_admin',
    'Tenant Administrator',
    'Tenant management, user management',
    true
FROM tenants t
ON CONFLICT (tenant_id, name) DO NOTHING;

INSERT INTO tenant_roles (tenant_id, name, display_name, description, is_system_role)
SELECT 
    t.id,
    'security_admin',
    'Security Administrator',
    'Security settings, compliance, reports',
    true
FROM tenants t
ON CONFLICT (tenant_id, name) DO NOTHING;

INSERT INTO tenant_roles (tenant_id, name, display_name, description, is_system_role)
SELECT 
    t.id,
    'analyst',
    'Security Analyst',
    'Data analysis, reporting, monitoring',
    true
FROM tenants t
ON CONFLICT (tenant_id, name) DO NOTHING;

INSERT INTO tenant_roles (tenant_id, name, display_name, description, is_system_role)
SELECT 
    t.id,
    'viewer',
    'Viewer',
    'Read-only access to tenant data',
    true
FROM tenants t
ON CONFLICT (tenant_id, name) DO NOTHING;

INSERT INTO tenant_roles (tenant_id, name, display_name, description, is_system_role)
SELECT 
    t.id,
    'api_user',
    'API User',
    'API-only access for integrations',
    true
FROM tenants t
ON CONFLICT (tenant_id, name) DO NOTHING;

-- =================================================================
-- Assign Permissions to Tenant Roles
-- =================================================================

-- Tenant Owner gets all permissions
INSERT INTO tenant_role_permissions (role_id, permission_id)
SELECT tr.id, tp.id
FROM tenant_roles tr
JOIN tenant_permissions tp ON true
WHERE tr.name = 'tenant_owner'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Tenant Admin gets most permissions except billing
INSERT INTO tenant_role_permissions (role_id, permission_id)
SELECT tr.id, tp.id
FROM tenant_roles tr
JOIN tenant_permissions tp ON true
WHERE tr.name = 'tenant_admin'
  AND tp.name NOT LIKE 'billing.%'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Security Admin gets security, compliance, and reporting permissions
INSERT INTO tenant_role_permissions (role_id, permission_id)
SELECT tr.id, tp.id
FROM tenant_roles tr
JOIN tenant_permissions tp ON true
WHERE tr.name = 'security_admin'
  AND (tp.resource IN ('assets', 'sensors', 'reports', 'compliance') 
       OR tp.name LIKE 'compliance.%'
       OR tp.name LIKE 'reports.%')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Analyst gets read permissions for most resources and create/update for reports
INSERT INTO tenant_role_permissions (role_id, permission_id)
SELECT tr.id, tp.id
FROM tenant_roles tr
JOIN tenant_permissions tp ON true
WHERE tr.name = 'analyst'
  AND (tp.action = 'read' 
       OR (tp.resource = 'reports' AND tp.action IN ('create', 'update')))
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Viewer gets only read permissions
INSERT INTO tenant_role_permissions (role_id, permission_id)
SELECT tr.id, tp.id
FROM tenant_roles tr
JOIN tenant_permissions tp ON true
WHERE tr.name = 'viewer'
  AND tp.action = 'read'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- API User gets read permissions for assets and sensors
INSERT INTO tenant_role_permissions (role_id, permission_id)
SELECT tr.id, tp.id
FROM tenant_roles tr
JOIN tenant_permissions tp ON true
WHERE tr.name = 'api_user'
  AND tp.action = 'read'
  AND tp.resource IN ('assets', 'sensors')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- =================================================================
-- Migrate Existing Users to New Role System
-- =================================================================

-- Map existing users to their new roles based on current role field
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
END
ON CONFLICT (user_id, tenant_id, role_id) DO NOTHING;

-- =================================================================
-- Create Default Platform Administrator
-- =================================================================

-- Create a default super admin platform user
-- Password: 'admin123' (hashed with Argon2id)
INSERT INTO platform_users (id, email, password_hash, first_name, last_name, role_id, is_active, email_verified)
SELECT 
    '00000000-0000-0000-0000-000000000001',
    'admin@crypto-inventory.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye/D7zrZI/PCMZ6qO8PQ8DbZOF5.XzEQm',
    'Platform',
    'Administrator',
    pr.id,
    true,
    true
FROM platform_roles pr
WHERE pr.name = 'super_admin'
ON CONFLICT (email) DO NOTHING;

-- =================================================================
-- Seed Complete
-- =================================================================
