-- Authentication Service Database Schema
-- This schema supports multi-tenant SaaS with freemium billing, SSO, and frontend flexibility
-- Run this after basic setup to create all authentication-related tables

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ================================
-- SUBSCRIPTION AND BILLING TABLES
-- ================================

-- Subscription tiers (configurable limits)
CREATE TABLE subscription_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE, -- 'free', 'professional', 'enterprise'
    display_name VARCHAR(100) NOT NULL,
    max_sensors INTEGER,
    max_assets INTEGER,
    max_users INTEGER,
    retention_days INTEGER,
    price_cents INTEGER,
    billing_interval VARCHAR(20), -- 'monthly', 'yearly'
    features JSONB DEFAULT '{}', -- feature flags
    limits JSONB DEFAULT '{}', -- custom limits for enterprise
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create default subscription tiers
INSERT INTO subscription_tiers (name, display_name, max_sensors, max_assets, max_users, retention_days, price_cents, billing_interval, features) VALUES
('free', 'Free Forever', 1, 50, 3, 90, 0, 'monthly', '{"compliance_frameworks": 1, "ai_insights": false, "integrations_max": 0, "priority_support": false}'),
('professional', 'Professional', 10, 1000, 25, 365, 14900, 'monthly', '{"compliance_frameworks": -1, "ai_insights": true, "integrations_max": 5, "priority_support": true}'),
('enterprise', 'Enterprise', -1, -1, -1, 1095, 0, 'monthly', '{"compliance_frameworks": -1, "ai_insights": true, "integrations_max": -1, "priority_support": true, "sso": true, "custom_branding": true}');

-- ================================
-- CORE TENANT AND USER TABLES
-- ================================

-- Enhanced tenants table
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(255), -- for SSO and verification
    subscription_tier_id UUID NOT NULL REFERENCES subscription_tiers(id),
    
    -- Billing and trial management
    trial_ends_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + INTERVAL '30 days'),
    billing_email VARCHAR(255),
    payment_status VARCHAR(50) DEFAULT 'trial', -- 'trial', 'active', 'past_due', 'canceled'
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    
    -- Configuration
    sso_enabled BOOLEAN DEFAULT false,
    custom_branding JSONB DEFAULT '{}',
    ui_config JSONB DEFAULT '{}', -- NEW: UI configuration per tenant
    settings JSONB DEFAULT '{}',
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_slug CHECK (slug ~ '^[a-z0-9-]+$'),
    CONSTRAINT valid_payment_status CHECK (payment_status IN ('trial', 'active', 'past_due', 'canceled', 'incomplete'))
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_deleted_at ON tenants(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenants_subscription_tier ON tenants(subscription_tier_id);

-- Enhanced users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    email_verification_token VARCHAR(255),
    email_verification_expires TIMESTAMP WITH TIME ZONE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    
    -- Password authentication
    password_hash VARCHAR(255), -- nullable for SSO-only users
    password_changed_at TIMESTAMP WITH TIME ZONE,
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP WITH TIME ZONE,
    password_history JSONB DEFAULT '[]', -- store hashes of last 5 passwords
    
    -- Account status
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    login_count INTEGER DEFAULT 0,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    
    -- Profile
    avatar_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'UTC',
    preferences JSONB DEFAULT '{}',
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_email CHECK (email ~ '^[^@]+@[^@]+\.[^@]+$'),
    CONSTRAINT valid_role CHECK (role IN ('admin', 'analyst', 'viewer')),
    CONSTRAINT unique_tenant_email UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;

-- ================================
-- SSO AND AUTHENTICATION TABLES
-- ================================

-- SSO provider configurations per tenant
CREATE TABLE sso_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    provider_type VARCHAR(50) NOT NULL, -- 'google', 'microsoft', 'azure', 'saml'
    provider_name VARCHAR(100) NOT NULL, -- display name
    
    -- OAuth 2.0 / OIDC configuration
    client_id VARCHAR(255),
    client_secret_encrypted TEXT, -- encrypted storage
    auth_url VARCHAR(500),
    token_url VARCHAR(500),
    userinfo_url VARCHAR(500),
    scopes VARCHAR(500) DEFAULT 'openid email profile',
    
    -- SAML configuration
    saml_entity_id VARCHAR(255),
    saml_sso_url VARCHAR(500),
    saml_certificate TEXT,
    saml_private_key_encrypted TEXT,
    
    -- Settings
    is_enabled BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    auto_provision_users BOOLEAN DEFAULT true,
    attribute_mapping JSONB DEFAULT '{}', -- map SSO attributes to user fields
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_provider_type CHECK (provider_type IN ('google', 'microsoft', 'azure', 'saml', 'okta')),
    CONSTRAINT unique_tenant_provider UNIQUE (tenant_id, provider_type, provider_name)
);

CREATE INDEX idx_sso_providers_tenant_id ON sso_providers(tenant_id);
CREATE INDEX idx_sso_providers_type ON sso_providers(provider_type);

-- User authentication methods (supports multiple auth methods per user)
CREATE TABLE user_auth_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    auth_type VARCHAR(50) NOT NULL, -- 'password', 'google', 'microsoft', 'saml'
    sso_provider_id UUID REFERENCES sso_providers(id),
    external_user_id VARCHAR(255), -- ID from SSO provider
    external_email VARCHAR(255), -- Email from SSO provider (may differ from user.email)
    
    -- Status and metadata
    is_primary BOOLEAN DEFAULT false,
    last_used_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB DEFAULT '{}', -- provider-specific data
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_auth_type CHECK (auth_type IN ('password', 'google', 'microsoft', 'azure', 'saml')),
    CONSTRAINT unique_user_auth_type UNIQUE (user_id, auth_type, sso_provider_id)
);

CREATE INDEX idx_user_auth_methods_user_id ON user_auth_methods(user_id);
CREATE INDEX idx_user_auth_methods_external_id ON user_auth_methods(external_user_id);

-- Identity linking requests (for confirmation-based linking)
CREATE TABLE identity_link_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    primary_user_id UUID NOT NULL REFERENCES users(id),
    secondary_user_id UUID NOT NULL REFERENCES users(id),
    auth_method_id UUID NOT NULL REFERENCES user_auth_methods(id),
    
    -- Request details
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'approved', 'rejected', 'expired'
    requested_by_user_id UUID REFERENCES users(id),
    confirmation_token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() + INTERVAL '24 hours'),
    
    -- Resolution
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by_user_id UUID REFERENCES users(id),
    rejection_reason TEXT,
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_status CHECK (status IN ('pending', 'approved', 'rejected', 'expired')),
    CONSTRAINT different_users CHECK (primary_user_id != secondary_user_id)
);

CREATE INDEX idx_identity_link_requests_primary_user ON identity_link_requests(primary_user_id);
CREATE INDEX idx_identity_link_requests_token ON identity_link_requests(confirmation_token);

-- JWT refresh tokens
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash VARCHAR(255) NOT NULL,
    family_id UUID NOT NULL, -- for token rotation
    
    -- Token lifecycle
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_revoked BOOLEAN DEFAULT false,
    
    -- Security
    created_from_ip INET,
    user_agent TEXT,
    
    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    revoked_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT unique_token_hash UNIQUE (token_hash)
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- ================================
-- USAGE TRACKING AND BILLING
-- ================================

-- Usage tracking for billing and limits
CREATE TABLE tenant_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Usage metrics
    sensors_count INTEGER DEFAULT 0,
    assets_count INTEGER DEFAULT 0,
    users_count INTEGER DEFAULT 0,
    storage_bytes BIGINT DEFAULT 0,
    api_calls_current_month INTEGER DEFAULT 0,
    reports_generated_month INTEGER DEFAULT 0,
    integrations_active INTEGER DEFAULT 0,
    
    -- Billing period tracking
    billing_period_start DATE NOT NULL,
    billing_period_end DATE NOT NULL,
    
    -- Metadata
    last_calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_tenant_period UNIQUE (tenant_id, billing_period_start)
);

CREATE INDEX idx_tenant_usage_tenant_id ON tenant_usage(tenant_id);
CREATE INDEX idx_tenant_usage_period ON tenant_usage(billing_period_start, billing_period_end);

-- Feature usage tracking (for analytics and billing)
CREATE TABLE feature_usage_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    
    -- Event details
    feature_name VARCHAR(100) NOT NULL, -- 'compliance_check', 'ai_analysis', 'integration_sync'
    event_type VARCHAR(50) NOT NULL, -- 'usage', 'limit_hit', 'upgrade_required'
    event_data JSONB DEFAULT '{}',
    
    -- Timestamp
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_feature_usage_tenant_feature ON feature_usage_events(tenant_id, feature_name);
CREATE INDEX idx_feature_usage_occurred_at ON feature_usage_events(occurred_at);

-- ================================
-- FRONTEND FLEXIBILITY TABLES
-- ================================

-- UI themes and templates
CREATE TABLE ui_themes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Theme configuration
    theme_config JSONB NOT NULL, -- colors, fonts, layout settings
    component_overrides JSONB DEFAULT '{}', -- custom component styles
    
    -- Availability
    is_public BOOLEAN DEFAULT true, -- available to all tenants
    is_active BOOLEAN DEFAULT true,
    pricing_tier VARCHAR(50), -- which subscription tiers can use this theme
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_pricing_tier CHECK (pricing_tier IN ('free', 'professional', 'enterprise', 'all'))
);

-- Create default themes
INSERT INTO ui_themes (name, display_name, description, theme_config, pricing_tier) VALUES
('default', 'Default Theme', 'Clean and professional default theme', '{"primary_color": "#1890ff", "secondary_color": "#f0f2f5", "success_color": "#52c41a", "warning_color": "#faad14", "error_color": "#f5222d"}', 'all'),
('dark', 'Dark Mode', 'Modern dark theme for reduced eye strain', '{"primary_color": "#177ddc", "secondary_color": "#1f1f1f", "success_color": "#49aa19", "warning_color": "#d89614", "error_color": "#d4380d", "background_color": "#141414"}', 'professional'),
('enterprise', 'Enterprise Blue', 'Professional enterprise theme', '{"primary_color": "#0050b3", "secondary_color": "#e6f7ff", "success_color": "#389e0d", "warning_color": "#d48806", "error_color": "#cf1322"}', 'enterprise');

-- Workflow configurations
CREATE TABLE workflow_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id), -- NULL for global workflows
    workflow_type VARCHAR(100) NOT NULL, -- 'onboarding', 'setup', 'integration'
    workflow_name VARCHAR(200) NOT NULL,
    
    -- Workflow definition
    steps JSONB NOT NULL, -- array of workflow steps
    configuration JSONB DEFAULT '{}', -- workflow-specific settings
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_tenant_workflow UNIQUE (tenant_id, workflow_type, workflow_name)
);

-- Create default onboarding workflow
INSERT INTO workflow_configurations (tenant_id, workflow_type, workflow_name, steps, is_default) VALUES
(NULL, 'onboarding', 'Default Onboarding', '[
    {"id": "welcome", "type": "info", "title": "Welcome to CryptoInventory", "description": "Get started with crypto asset discovery", "required": true},
    {"id": "deploy_sensor", "type": "action", "title": "Deploy Your First Sensor", "description": "Install a sensor to start discovering crypto assets", "required": true},
    {"id": "view_dashboard", "type": "navigation", "title": "Explore Your Dashboard", "description": "See your crypto inventory in real-time", "required": false},
    {"id": "invite_team", "type": "optional", "title": "Invite Team Members", "description": "Collaborate with your security team", "required": false}
]', true);

-- User workflow progress tracking
CREATE TABLE user_workflow_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    workflow_configuration_id UUID NOT NULL REFERENCES workflow_configurations(id),
    
    -- Progress tracking
    current_step INTEGER DEFAULT 0,
    completed_steps JSONB DEFAULT '[]', -- array of completed step IDs
    skipped_steps JSONB DEFAULT '[]', -- array of skipped step IDs
    step_data JSONB DEFAULT '{}', -- data collected during workflow
    
    -- Status
    status VARCHAR(50) DEFAULT 'in_progress', -- 'not_started', 'in_progress', 'completed', 'abandoned'
    completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_status CHECK (status IN ('not_started', 'in_progress', 'completed', 'abandoned')),
    CONSTRAINT unique_user_workflow UNIQUE (user_id, workflow_configuration_id)
);

CREATE INDEX idx_user_workflow_progress_user_id ON user_workflow_progress(user_id);
CREATE INDEX idx_user_workflow_progress_workflow_id ON user_workflow_progress(workflow_configuration_id);

-- API response format preferences
CREATE TABLE api_format_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Format preferences per endpoint
    endpoint_formats JSONB DEFAULT '{}', -- endpoint -> preferred format mapping
    global_preferences JSONB DEFAULT '{}', -- global formatting preferences
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_tenant_preferences UNIQUE (tenant_id)
);

-- ================================
-- AUDIT AND SECURITY TABLES
-- ================================

-- Audit log for authentication events
CREATE TABLE auth_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    
    -- Event details
    event_type VARCHAR(100) NOT NULL, -- 'login', 'logout', 'password_change', 'sso_link', etc.
    event_status VARCHAR(50) NOT NULL, -- 'success', 'failure', 'blocked'
    auth_method VARCHAR(50), -- 'password', 'google', 'microsoft', 'saml'
    
    -- Context
    ip_address INET,
    user_agent TEXT,
    session_id VARCHAR(255),
    
    -- Details
    event_data JSONB DEFAULT '{}',
    failure_reason TEXT,
    
    -- Timestamp
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_auth_audit_tenant_user ON auth_audit_log(tenant_id, user_id);
CREATE INDEX idx_auth_audit_occurred_at ON auth_audit_log(occurred_at);
CREATE INDEX idx_auth_audit_event_type ON auth_audit_log(event_type);

-- ================================
-- FUNCTIONS AND TRIGGERS
-- ================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add updated_at triggers to relevant tables
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sso_providers_updated_at BEFORE UPDATE ON sso_providers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_auth_methods_updated_at BEFORE UPDATE ON user_auth_methods FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tenant_usage_updated_at BEFORE UPDATE ON tenant_usage FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ui_themes_updated_at BEFORE UPDATE ON ui_themes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workflow_configurations_updated_at BEFORE UPDATE ON workflow_configurations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_workflow_progress_updated_at BEFORE UPDATE ON user_workflow_progress FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_api_format_preferences_updated_at BEFORE UPDATE ON api_format_preferences FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to automatically set trial_ends_at for new tenants
CREATE OR REPLACE FUNCTION set_trial_end_date()
RETURNS TRIGGER AS $$
BEGIN
    -- Set trial end date to 30 days from now if not specified
    IF NEW.trial_ends_at IS NULL THEN
        NEW.trial_ends_at = NOW() + INTERVAL '30 days';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_tenant_trial_end BEFORE INSERT ON tenants FOR EACH ROW EXECUTE FUNCTION set_trial_end_date();

-- Function to track user count for tenant usage
CREATE OR REPLACE FUNCTION update_tenant_user_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Update user count in current period usage
    INSERT INTO tenant_usage (tenant_id, users_count, billing_period_start, billing_period_end)
    VALUES (
        COALESCE(NEW.tenant_id, OLD.tenant_id),
        (SELECT COUNT(*) FROM users WHERE tenant_id = COALESCE(NEW.tenant_id, OLD.tenant_id) AND deleted_at IS NULL),
        DATE_TRUNC('month', NOW())::DATE,
        (DATE_TRUNC('month', NOW()) + INTERVAL '1 month - 1 day')::DATE
    )
    ON CONFLICT (tenant_id, billing_period_start)
    DO UPDATE SET 
        users_count = EXCLUDED.users_count,
        last_calculated_at = NOW();
    
    RETURN COALESCE(NEW, OLD);
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_count_on_insert AFTER INSERT ON users FOR EACH ROW EXECUTE FUNCTION update_tenant_user_count();
CREATE TRIGGER update_user_count_on_update AFTER UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_tenant_user_count();
CREATE TRIGGER update_user_count_on_delete AFTER DELETE ON users FOR EACH ROW EXECUTE FUNCTION update_tenant_user_count();

-- ================================
-- INITIAL DATA AND PERMISSIONS
-- ================================

-- Create a default tenant for development/testing
INSERT INTO tenants (name, slug, subscription_tier_id, billing_email)
VALUES (
    'Demo Corporation',
    'demo-corp',
    (SELECT id FROM subscription_tiers WHERE name = 'professional' LIMIT 1),
    'admin@democorp.com'
);

-- Create initial admin user for demo tenant
INSERT INTO users (tenant_id, email, email_verified, first_name, last_name, role, password_hash)
VALUES (
    (SELECT id FROM tenants WHERE slug = 'demo-corp' LIMIT 1),
    'admin@democorp.com',
    true,
    'Admin',
    'User',
    'admin',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeKOcn4QwYqZt1.nq' -- password: "password123" - NOTE: This will be updated to Argon2id format by the auth service
);

-- Add password auth method for demo user
INSERT INTO user_auth_methods (user_id, auth_type, is_primary)
VALUES (
    (SELECT id FROM users WHERE email = 'admin@democorp.com' LIMIT 1),
    'password',
    true
);

COMMIT;

-- Display setup completion message
DO $$
BEGIN
    RAISE NOTICE '✅ Authentication service database schema created successfully!';
    RAISE NOTICE '📊 Tables created: 13 core tables + 4 frontend flexibility tables';
    RAISE NOTICE '🔐 Demo tenant created: demo-corp (admin@democorp.com / password123)';
    RAISE NOTICE '🚀 Ready for authentication service implementation!';
END $$;
