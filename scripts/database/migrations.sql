-- =================================================================
-- Crypto Inventory Platform - Database Migrations
-- Version: 1.0.0
-- =================================================================

-- =================================================================
-- Tenant Management
-- =================================================================

-- Tenants and users tables are created by the authentication schema
-- Skipping to avoid conflicts...

-- =================================================================
-- Network Discovery
-- =================================================================

-- Network assets table
CREATE TABLE IF NOT EXISTS network_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    hostname VARCHAR(255),
    ip_address INET,
    port INTEGER,
    asset_type asset_type NOT NULL,
    operating_system VARCHAR(100),
    environment environment_type,
    business_unit VARCHAR(100),
    owner_email VARCHAR(255),
    description TEXT,
    tags JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    first_discovered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_port CHECK (port IS NULL OR (port BETWEEN 1 AND 65535)),
    CONSTRAINT valid_owner_email CHECK (owner_email IS NULL OR owner_email ~ '^[^@]+@[^@]+\.[^@]+$'),
    CONSTRAINT asset_identifier CHECK (hostname IS NOT NULL OR ip_address IS NOT NULL)
);

-- Certificates table
CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    serial_number VARCHAR(255),
    subject_dn TEXT NOT NULL,
    issuer_dn TEXT NOT NULL,
    common_name VARCHAR(255),
    subject_alternative_names TEXT[],
    signature_algorithm VARCHAR(100),
    public_key_algorithm VARCHAR(100),
    public_key_size INTEGER,
    not_before TIMESTAMP WITH TIME ZONE,
    not_after TIMESTAMP WITH TIME ZONE,
    fingerprint_sha1 VARCHAR(40),
    fingerprint_sha256 VARCHAR(64) NOT NULL,
    certificate_pem TEXT,
    is_self_signed BOOLEAN DEFAULT FALSE,
    is_ca_certificate BOOLEAN DEFAULT FALSE,
    key_usage TEXT[],
    extended_key_usage TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_fingerprint_sha1 CHECK (fingerprint_sha1 IS NULL OR fingerprint_sha1 ~ '^[a-fA-F0-9]{40}$'),
    CONSTRAINT valid_fingerprint_sha256 CHECK (fingerprint_sha256 ~ '^[a-fA-F0-9]{64}$'),
    CONSTRAINT valid_key_size CHECK (public_key_size IS NULL OR public_key_size > 0),
    CONSTRAINT unique_fingerprint_per_tenant UNIQUE (tenant_id, fingerprint_sha256)
);

-- Crypto implementations table
CREATE TABLE IF NOT EXISTS crypto_implementations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES network_assets(id) ON DELETE CASCADE,
    protocol protocol_type NOT NULL,
    protocol_version VARCHAR(20),
    cipher_suite VARCHAR(255),
    key_exchange_algorithm VARCHAR(100),
    signature_algorithm VARCHAR(100),
    symmetric_encryption VARCHAR(100),
    hash_algorithm VARCHAR(100),
    key_size INTEGER,
    certificate_id UUID REFERENCES certificates(id),
    discovery_method discovery_method NOT NULL,
    confidence_score DECIMAL(3,2) DEFAULT 1.0,
    source_sensor_id UUID, -- References sensors table (added later)
    raw_data JSONB,
    risk_score INTEGER DEFAULT 0, -- 0-100 risk score
    compliance_status JSONB DEFAULT '{}', -- Framework compliance results
    first_discovered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_verified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_confidence CHECK (confidence_score BETWEEN 0.0 AND 1.0),
    CONSTRAINT valid_risk_score CHECK (risk_score BETWEEN 0 AND 100),
    CONSTRAINT valid_key_size CHECK (key_size IS NULL OR key_size > 0)
);

-- =================================================================
-- Sensor Management
-- =================================================================

-- Sensors table
CREATE TABLE IF NOT EXISTS sensors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sensor_type sensor_type NOT NULL,
    deployment_location VARCHAR(255),
    ip_address INET,
    hostname VARCHAR(255),
    version VARCHAR(50),
    configuration JSONB DEFAULT '{}',
    status sensor_status DEFAULT 'inactive',
    last_heartbeat_at TIMESTAMP WITH TIME ZONE,
    registration_token VARCHAR(255),
    token_expires_at TIMESTAMP WITH TIME ZONE,
    api_key_hash VARCHAR(255), -- For API authentication
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT unique_name_per_tenant UNIQUE (tenant_id, name),
    CONSTRAINT token_expiry_check CHECK (
        (registration_token IS NULL AND token_expires_at IS NULL) OR
        (registration_token IS NOT NULL AND token_expires_at IS NOT NULL)
    )
);

-- Add foreign key reference back to crypto_implementations
ALTER TABLE crypto_implementations 
    ADD CONSTRAINT fk_source_sensor 
    FOREIGN KEY (source_sensor_id) 
    REFERENCES sensors(id) ON DELETE SET NULL;

-- =================================================================
-- Compliance Framework
-- =================================================================

-- Compliance frameworks table
CREATE TABLE IF NOT EXISTS compliance_frameworks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    description TEXT,
    organization VARCHAR(255),
    effective_date DATE,
    rules JSONB NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_framework_version UNIQUE (name, version)
);

-- Compliance assessments table
CREATE TABLE IF NOT EXISTS compliance_assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    framework_id UUID NOT NULL REFERENCES compliance_frameworks(id),
    assessment_name VARCHAR(255) NOT NULL,
    scope_filter JSONB,
    overall_score DECIMAL(5,2),
    total_checks INTEGER DEFAULT 0,
    passed_checks INTEGER DEFAULT 0,
    failed_checks INTEGER DEFAULT 0,
    not_applicable_checks INTEGER DEFAULT 0,
    assessment_results JSONB,
    ai_insights JSONB DEFAULT '{}',
    assessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    assessed_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_score CHECK (overall_score IS NULL OR overall_score BETWEEN 0 AND 100),
    CONSTRAINT valid_checks CHECK (
        total_checks = passed_checks + failed_checks + not_applicable_checks
    )
);

-- =================================================================
-- Reporting
-- =================================================================

-- Reports table
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    report_type report_type NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    template_id UUID,
    parameters JSONB,
    status report_status DEFAULT 'pending',
    file_path VARCHAR(500),
    file_format file_format,
    file_size_bytes BIGINT,
    generated_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    ai_generated BOOLEAN DEFAULT FALSE,
    requested_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_file_size CHECK (file_size_bytes IS NULL OR file_size_bytes > 0)
);

-- =================================================================
-- AI/ML Tables
-- =================================================================

-- AI models table
CREATE TABLE IF NOT EXISTS ai_models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    model_type VARCHAR(100) NOT NULL, -- 'anomaly_detection', 'risk_scoring', 'nlp', etc.
    description TEXT,
    file_path VARCHAR(500),
    hyperparameters JSONB DEFAULT '{}',
    metrics JSONB DEFAULT '{}',
    active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_active_model_type UNIQUE (model_type, active) 
        DEFERRABLE INITIALLY DEFERRED
);

-- AI analysis results table
CREATE TABLE IF NOT EXISTS ai_analysis_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    model_id UUID NOT NULL REFERENCES ai_models(id),
    target_type VARCHAR(50) NOT NULL, -- 'asset', 'crypto_implementation', 'certificate'
    target_id UUID NOT NULL,
    analysis_type VARCHAR(100) NOT NULL,
    confidence_score DECIMAL(3,2),
    results JSONB NOT NULL,
    anomaly_detected BOOLEAN DEFAULT FALSE,
    risk_level VARCHAR(20), -- 'low', 'medium', 'high', 'critical'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT valid_confidence CHECK (confidence_score IS NULL OR confidence_score BETWEEN 0.0 AND 1.0),
    CONSTRAINT valid_risk_level CHECK (risk_level IS NULL OR risk_level IN ('low', 'medium', 'high', 'critical'))
);

-- =================================================================
-- Audit and Logging
-- =================================================================

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit.audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =================================================================
-- Indexes for Performance
-- =================================================================

-- Tenant isolation indexes
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_network_assets_tenant_id ON network_assets(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_tenant_id ON crypto_implementations(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_certificates_tenant_id ON certificates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_sensors_tenant_id ON sensors(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_compliance_assessments_tenant_id ON compliance_assessments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_reports_tenant_id ON reports(tenant_id);

-- Search and lookup indexes
-- User indexes are created by the authentication schema
CREATE INDEX IF NOT EXISTS idx_network_assets_ip_port ON network_assets(ip_address, port) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_network_assets_hostname ON network_assets(hostname) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_certificates_fingerprint_sha256 ON certificates(fingerprint_sha256);
CREATE INDEX IF NOT EXISTS idx_certificates_common_name ON certificates(common_name);
CREATE INDEX IF NOT EXISTS idx_certificates_expiry ON certificates(not_after);
CREATE INDEX IF NOT EXISTS idx_sensors_status ON sensors(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sensors_last_heartbeat ON sensors(last_heartbeat_at) WHERE deleted_at IS NULL;

-- Crypto implementations indexes
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_asset_id ON crypto_implementations(asset_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_protocol ON crypto_implementations(protocol) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_cipher_suite ON crypto_implementations(cipher_suite) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_discovery ON crypto_implementations(discovery_method, first_discovered_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_certificate_id ON crypto_implementations(certificate_id) WHERE deleted_at IS NULL;

-- Time-based indexes
CREATE INDEX IF NOT EXISTS idx_compliance_assessments_date ON compliance_assessments(assessed_at);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);
CREATE INDEX IF NOT EXISTS idx_ai_analysis_results_created_at ON ai_analysis_results(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit.audit_logs(created_at);

-- JSON indexes for tags and metadata
CREATE INDEX IF NOT EXISTS idx_network_assets_tags ON network_assets USING GIN(tags) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_raw_data ON crypto_implementations USING GIN(raw_data) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sensors_configuration ON sensors USING GIN(configuration) WHERE deleted_at IS NULL;

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_crypto_implementations_tenant_protocol ON crypto_implementations(tenant_id, protocol) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_certificates_tenant_expiry ON certificates(tenant_id, not_after);
CREATE INDEX IF NOT EXISTS idx_ai_analysis_results_tenant_type ON ai_analysis_results(tenant_id, target_type, analysis_type);

-- =================================================================
-- Triggers for Automatic Updates
-- =================================================================

-- Update timestamps for core platform tables (tenants/users triggers are in auth schema)

CREATE TRIGGER update_network_assets_updated_at BEFORE UPDATE ON network_assets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certificates_updated_at BEFORE UPDATE ON certificates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_crypto_implementations_updated_at BEFORE UPDATE ON crypto_implementations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sensors_updated_at BEFORE UPDATE ON sensors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_compliance_frameworks_updated_at BEFORE UPDATE ON compliance_frameworks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_compliance_assessments_updated_at BEFORE UPDATE ON compliance_assessments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reports_updated_at BEFORE UPDATE ON reports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ai_models_updated_at BEFORE UPDATE ON ai_models
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
