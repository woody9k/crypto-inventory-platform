-- =================================================================
-- Crypto Inventory Platform - Database Initialization
-- =================================================================

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Set timezone
SET timezone = 'UTC';

-- Create custom types
CREATE TYPE subscription_tier AS ENUM ('basic', 'professional', 'enterprise');
CREATE TYPE user_role AS ENUM ('admin', 'analyst', 'viewer');
CREATE TYPE asset_type AS ENUM ('server', 'endpoint', 'service', 'appliance');
CREATE TYPE environment_type AS ENUM ('production', 'staging', 'development', 'test');
CREATE TYPE discovery_method AS ENUM ('passive', 'active', 'manual', 'integration');
CREATE TYPE sensor_type AS ENUM ('network', 'endpoint', 'cloud', 'api');
CREATE TYPE sensor_status AS ENUM ('active', 'inactive', 'error', 'maintenance');
CREATE TYPE protocol_type AS ENUM ('TLS', 'SSH', 'IPSec', 'VPN', 'Database', 'API');
CREATE TYPE report_type AS ENUM ('compliance', 'inventory', 'risk', 'certificate');
CREATE TYPE report_status AS ENUM ('pending', 'generating', 'completed', 'failed', 'expired');
CREATE TYPE file_format AS ENUM ('PDF', 'Excel', 'CSV', 'JSON');

-- =================================================================
-- Utility Functions
-- =================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Function to generate tenant slug
CREATE OR REPLACE FUNCTION generate_tenant_slug(tenant_name TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN lower(regexp_replace(trim(tenant_name), '[^a-zA-Z0-9]+', '-', 'g'));
END;
$$ language 'plpgsql';

-- =================================================================
-- Grant permissions to application user
-- =================================================================

GRANT USAGE ON SCHEMA public TO crypto_user;
GRANT CREATE ON SCHEMA public TO crypto_user;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO crypto_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO crypto_user;

-- Grant permissions on future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON SEQUENCES TO crypto_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crypto_user;

-- Create application-specific schemas
CREATE SCHEMA IF NOT EXISTS compliance;
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE SCHEMA IF NOT EXISTS audit;

GRANT USAGE ON SCHEMA compliance TO crypto_user;
GRANT USAGE ON SCHEMA analytics TO crypto_user;
GRANT USAGE ON SCHEMA audit TO crypto_user;
GRANT CREATE ON SCHEMA compliance TO crypto_user;
GRANT CREATE ON SCHEMA analytics TO crypto_user;
GRANT CREATE ON SCHEMA audit TO crypto_user;
