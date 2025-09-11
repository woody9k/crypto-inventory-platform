# Database Schema Documentation

This document describes the database schema, sample data, and data management for the Crypto Inventory Management System.

## 🗄️ Database Overview

- **Database**: PostgreSQL 15
- **Name**: `crypto_inventory`
- **User**: `crypto_user`
- **Schema**: Multi-tenant with complete data isolation

## 📋 Schema Structure

### Core Tables

#### Tenants
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(255),
    subscription_tier VARCHAR(50) DEFAULT 'basic',
    payment_status VARCHAR(50) DEFAULT 'active',
    billing_email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Users (Tenant Users)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'viewer',
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);
```

#### Platform Users (Platform Admins)
```sql
CREATE TABLE platform_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role_id UUID REFERENCES platform_roles(id),
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Asset Management

#### Network Assets
```sql
CREATE TABLE network_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    hostname VARCHAR(255),
    ip_address INET,
    port INTEGER,
    asset_type VARCHAR(50),
    operating_system VARCHAR(255),
    environment VARCHAR(50),
    business_unit VARCHAR(100),
    owner_email VARCHAR(255),
    description TEXT,
    tags JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Certificates
```sql
CREATE TABLE certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    serial_number VARCHAR(255),
    subject_dn TEXT,
    issuer_dn TEXT,
    common_name VARCHAR(255),
    subject_alternative_names TEXT[],
    signature_algorithm VARCHAR(100),
    public_key_algorithm VARCHAR(100),
    public_key_size INTEGER,
    not_before TIMESTAMP WITH TIME ZONE,
    not_after TIMESTAMP WITH TIME ZONE,
    fingerprint_sha1 VARCHAR(40),
    fingerprint_sha256 VARCHAR(64),
    is_self_signed BOOLEAN DEFAULT FALSE,
    is_ca_certificate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Crypto Implementations
```sql
CREATE TABLE crypto_implementations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id UUID REFERENCES network_assets(id) ON DELETE CASCADE,
    protocol VARCHAR(50),
    protocol_version VARCHAR(20),
    cipher_suite VARCHAR(100),
    key_exchange_algorithm VARCHAR(100),
    signature_algorithm VARCHAR(100),
    symmetric_encryption VARCHAR(100),
    hash_algorithm VARCHAR(100),
    key_size INTEGER,
    certificate_id UUID REFERENCES certificates(id),
    discovery_method VARCHAR(50),
    confidence_score DECIMAL(3,2),
    risk_score INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Sensor Management

#### Sensors
```sql
CREATE TABLE sensors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sensor_type VARCHAR(50),
    deployment_location VARCHAR(255),
    ip_address INET,
    hostname VARCHAR(255),
    version VARCHAR(50),
    status VARCHAR(50) DEFAULT 'inactive',
    last_heartbeat_at TIMESTAMP WITH TIME ZONE,
    api_key_hash VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Discovery Events
```sql
CREATE TABLE discovery_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    sensor_id UUID REFERENCES sensors(id) ON DELETE CASCADE,
    asset_id UUID REFERENCES network_assets(id) ON DELETE CASCADE,
    event_type VARCHAR(100),
    event_data JSONB,
    confidence_score DECIMAL(3,2),
    risk_score INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Compliance & Reporting

#### Compliance Frameworks
```sql
CREATE TABLE compliance_frameworks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    version VARCHAR(50),
    description TEXT,
    requirements JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Compliance Assessments
```sql
CREATE TABLE compliance_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    framework_id UUID REFERENCES compliance_frameworks(id),
    assessment_name VARCHAR(255),
    scope_filter JSONB,
    overall_score DECIMAL(5,2),
    total_checks INTEGER,
    passed_checks INTEGER,
    failed_checks INTEGER,
    not_applicable_checks INTEGER,
    assessment_results JSONB,
    assessed_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Audit & Logging

#### Audit Logs
```sql
CREATE SCHEMA audit;

CREATE TABLE audit.audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## 🌱 Sample Data

### Default Tenants

#### Demo Corporation
- **Tenant ID**: `550e8400-e29b-41d4-a716-446655440000`
- **Slug**: `demo-corp`
- **Users**:
  - `admin@democorp.com` (admin)
  - `analyst@democorp.com` (analyst)
  - `viewer@democorp.com` (viewer)

#### Debban Corporation (Rich Sample Data)
- **Tenant ID**: `550e8400-e29b-41d4-a716-446655440200`
- **Slug**: `debban-corp`
- **Users**:
  - `brian@debban.com` (admin)
  - `sarah@debban.com` (analyst)
  - `mike@debban.com` (viewer)

### Sample Data Features

#### Network Assets (10 total for Debban Corp)
- Production web servers (2)
- API Gateway
- Database servers (2)
- Redis cache
- Monitoring server
- Legacy application (with security issues)
- Development & staging servers

#### Certificates (3 total for Debban Corp)
- Valid wildcard cert (*.debban.com)
- **Expiring database cert** (3 days - urgent!)
- Self-signed legacy cert (security issue)

#### Crypto Implementations (4 total for Debban Corp)
- Excellent TLS 1.3 (web servers)
- Good TLS 1.2 (API gateway, database)
- **Weak TLS 1.0** (legacy app - high risk)

#### Compliance Assessment
- PCI DSS assessment with 78.5% score
- Critical findings on legacy systems
- Certificate expiry warnings

#### Sensors (4 total for Debban Corp)
- 3 active sensors (datacenter, office, cloud)
- 1 inactive staging sensor

## 🔧 Data Management

### Loading Sample Data

```bash
# Load comprehensive demo data for brian@debban.com
docker-compose exec postgres psql -U crypto_user -d crypto_inventory -f /scripts/database/seed_brian_debban.sql

# Load basic demo data (admin@democorp.com)
docker-compose exec postgres psql -U crypto_user -d crypto_inventory -f /scripts/database/seed.sql
```

### Database Maintenance

```bash
# Connect to database
docker-compose exec postgres psql -U crypto_user -d crypto_inventory

# Backup database
docker-compose exec postgres pg_dump -U crypto_user crypto_inventory > backup_$(date +%Y%m%d).sql

# Restore database
docker-compose exec -T postgres psql -U crypto_user -d crypto_inventory < backup_file.sql

# Vacuum and analyze
VACUUM ANALYZE;
```

### Common Queries

#### Check Tenant Data
```sql
-- List all tenants
SELECT id, name, slug, subscription_tier, created_at FROM tenants;

-- Count users per tenant
SELECT t.name, COUNT(u.id) as user_count 
FROM tenants t 
LEFT JOIN users u ON t.id = u.tenant_id 
GROUP BY t.id, t.name;

-- List assets per tenant
SELECT t.name, COUNT(a.id) as asset_count 
FROM tenants t 
LEFT JOIN network_assets a ON t.id = a.tenant_id 
GROUP BY t.id, t.name;
```

#### Security Analysis
```sql
-- Find expiring certificates
SELECT c.common_name, c.not_after, 
       EXTRACT(DAYS FROM (c.not_after - NOW())) as days_until_expiry
FROM certificates c 
WHERE c.not_after < NOW() + INTERVAL '30 days'
ORDER BY c.not_after;

-- Find high-risk crypto implementations
SELECT a.hostname, ci.protocol, ci.protocol_version, ci.risk_score
FROM crypto_implementations ci
JOIN network_assets a ON ci.asset_id = a.id
WHERE ci.risk_score > 50
ORDER BY ci.risk_score DESC;

-- Find self-signed certificates
SELECT c.common_name, c.subject_dn, a.hostname
FROM certificates c
JOIN network_assets a ON c.tenant_id = a.tenant_id
WHERE c.is_self_signed = true;
```

#### Compliance Reporting
```sql
-- Latest compliance assessment per tenant
SELECT t.name, ca.assessment_name, ca.overall_score, ca.created_at
FROM tenants t
JOIN compliance_assessments ca ON t.id = ca.tenant_id
WHERE ca.created_at = (
    SELECT MAX(created_at) 
    FROM compliance_assessments ca2 
    WHERE ca2.tenant_id = ca.tenant_id
);

-- Failed compliance checks
SELECT t.name, ca.assessment_name, ca.failed_checks, ca.total_checks,
       ROUND((ca.failed_checks::decimal / ca.total_checks) * 100, 2) as failure_rate
FROM tenants t
JOIN compliance_assessments ca ON t.id = ca.tenant_id
WHERE ca.failed_checks > 0
ORDER BY failure_rate DESC;
```

## 🔐 Security Considerations

### Data Isolation
- All tenant data is isolated by `tenant_id`
- Foreign key constraints ensure data integrity
- Cascade deletes prevent orphaned data

### Password Security
- All passwords hashed with bcrypt (cost 10)
- Default password for demo users: `admin123`
- Platform admin password needs hash update

### Audit Trail
- All user actions logged in `audit.audit_logs`
- Includes IP address, user agent, success/failure
- Immutable audit records

## 📊 Performance Optimization

### Indexes
```sql
-- Tenant isolation indexes
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_network_assets_tenant_id ON network_assets(tenant_id);
CREATE INDEX idx_certificates_tenant_id ON certificates(tenant_id);

-- Time-based queries
CREATE INDEX idx_certificates_not_after ON certificates(not_after);
CREATE INDEX idx_discovery_events_created_at ON discovery_events(created_at);
CREATE INDEX idx_audit_logs_created_at ON audit.audit_logs(created_at);

-- Search indexes
CREATE INDEX idx_network_assets_hostname ON network_assets(hostname);
CREATE INDEX idx_certificates_common_name ON certificates(common_name);
```

### Partitioning (Future)
Consider partitioning large tables by tenant_id or date for better performance.

---

*Last updated: 2025-01-09*
