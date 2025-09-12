-- =================================================================
-- Sample Data for brian@debban.com
-- =================================================================
-- This script creates comprehensive demo data for the brian@debban.com user

-- =================================================================
-- Create Tenant for brian@debban.com
-- =================================================================

-- Create tenant
INSERT INTO tenants (id, name, slug, domain, subscription_tier, payment_status, billing_email, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    'Debban Corporation',
    'debban-corp',
    'debban.com',
    'enterprise',
    'active',
    'billing@debban.com',
    NOW() - INTERVAL '6 months',
    NOW()
) ON CONFLICT (slug) DO NOTHING;

-- Create user
INSERT INTO users (id, tenant_id, email, first_name, last_name, password_hash, role, email_verified, is_active, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440201',
    '550e8400-e29b-41d4-a716-446655440200',
    'brian@debban.com',
    'Brian',
    'Debban',
    '$argon2id$v=19$m=65536,t=3,p=2$8Ll1hG8Y7AO+m8hQxRuozA$nO6gHyQ3JAccN5XWX5gjdnxTx+XutgIHYGCuXpJ2LKQ', -- password: Password123!
    'admin',
    true,
    true,
    NOW() - INTERVAL '6 months',
    NOW()
) ON CONFLICT (tenant_id, email) DO NOTHING;

-- Create additional users for the tenant
INSERT INTO users (id, tenant_id, email, first_name, last_name, password_hash, role, email_verified, is_active, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440202',
    '550e8400-e29b-41d4-a716-446655440200',
    'sarah@debban.com',
    'Sarah',
    'Johnson',
    '$argon2id$v=19$m=65536,t=3,p=2$8Ll1hG8Y7AO+m8hQxRuozA$nO6gHyQ3JAccN5XWX5gjdnxTx+XutgIHYGCuXpJ2LKQ', -- password: Password123!
    'analyst',
    true,
    true,
    NOW() - INTERVAL '5 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440203',
    '550e8400-e29b-41d4-a716-446655440200',
    'mike@debban.com',
    'Mike',
    'Chen',
    '$argon2id$v=19$m=65536,t=3,p=2$8Ll1hG8Y7AO+m8hQxRuozA$nO6gHyQ3JAccN5XWX5gjdnxTx+XutgIHYGCuXpJ2LKQ', -- password: Password123!
    'viewer',
    true,
    true,
    NOW() - INTERVAL '4 months',
    NOW()
) ON CONFLICT (tenant_id, email) DO NOTHING;

-- =================================================================
-- Network Assets for Debban Corporation
-- =================================================================

-- Production web servers
INSERT INTO network_assets (tenant_id, hostname, ip_address, port, asset_type, operating_system, environment, business_unit, owner_email, description, tags, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    'web-prod-01.debban.com',
    '10.10.1.10',
    443,
    'server',
    'Ubuntu 22.04 LTS',
    'production',
    'Engineering',
    'brian@debban.com',
    'Primary production web server - customer portal',
    '{"service": "web", "critical": true, "public_facing": true, "load_balancer": "primary"}',
    NOW() - INTERVAL '5 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'web-prod-02.debban.com',
    '10.10.1.11',
    443,
    'server',
    'Ubuntu 22.04 LTS',
    'production',
    'Engineering',
    'brian@debban.com',
    'Secondary production web server - customer portal',
    '{"service": "web", "critical": true, "public_facing": true, "load_balancer": "secondary"}',
    NOW() - INTERVAL '5 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'api-gateway.debban.com',
    '10.10.2.10',
    443,
    'service',
    'Alpine Linux',
    'production',
    'Engineering',
    'brian@debban.com',
    'API Gateway - microservices entry point',
    '{"service": "api", "critical": true, "public_facing": true, "microservices": true}',
    NOW() - INTERVAL '4 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'db-primary.debban.internal',
    '10.10.3.10',
    5432,
    'server',
    'PostgreSQL on Ubuntu 22.04',
    'production',
    'Database Team',
    'dba@debban.com',
    'Primary database server - customer data',
    '{"service": "database", "critical": true, "encrypted": true, "backup_enabled": true}',
    NOW() - INTERVAL '6 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'db-replica.debban.internal',
    '10.10.3.11',
    5432,
    'server',
    'PostgreSQL on Ubuntu 22.04',
    'production',
    'Database Team',
    'dba@debban.com',
    'Read replica database server',
    '{"service": "database", "critical": true, "encrypted": true, "replica": true}',
    NOW() - INTERVAL '5 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'redis-cache.debban.internal',
    '10.10.4.10',
    6379,
    'service',
    'Redis on Ubuntu 22.04',
    'production',
    'Engineering',
    'brian@debban.com',
    'Redis cache cluster - session storage',
    '{"service": "cache", "critical": true, "clustered": true}',
    NOW() - INTERVAL '4 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'monitoring.debban.internal',
    '10.10.5.10',
    3000,
    'server',
    'Ubuntu 22.04',
    'production',
    'DevOps',
    'devops@debban.com',
    'Monitoring and observability stack',
    '{"service": "monitoring", "critical": false, "internal": true}',
    NOW() - INTERVAL '3 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'legacy-app.debban.internal',
    '10.10.6.10',
    8080,
    'server',
    'Windows Server 2019',
    'production',
    'Legacy Systems',
    'legacy@debban.com',
    'Legacy application server - scheduled for migration',
    '{"service": "legacy", "critical": false, "migration_planned": true, "end_of_life": "2024-12-31"}',
    NOW() - INTERVAL '2 years',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'dev-server-01.debban.internal',
    '10.20.1.10',
    443,
    'server',
    'Ubuntu 22.04',
    'development',
    'Engineering',
    'brian@debban.com',
    'Development environment server',
    '{"service": "web", "environment": "development", "testing": true}',
    NOW() - INTERVAL '3 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    'staging.debban.com',
    '10.20.2.10',
    443,
    'server',
    'Ubuntu 22.04',
    'staging',
    'Engineering',
    'brian@debban.com',
    'Staging environment - pre-production testing',
    '{"service": "web", "environment": "staging", "testing": true}',
    NOW() - INTERVAL '2 months',
    NOW()
);

-- =================================================================
-- Certificates for Debban Corporation
-- =================================================================

-- Valid wildcard certificate
INSERT INTO certificates (tenant_id, serial_number, subject_dn, issuer_dn, common_name, subject_alternative_names, signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, fingerprint_sha1, fingerprint_sha256, is_self_signed, is_ca_certificate, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    '03:a1:b2:c3:d4:e5:f6:a7:b8:c9:d0:e1:f2:a3:b4:c5:d6:e7',
    'CN=*.debban.com,O=Debban Corporation,L=San Francisco,ST=California,C=US',
    'CN=Let''s Encrypt Authority X3,O=Let''s Encrypt,C=US',
    '*.debban.com',
    ARRAY['debban.com', 'www.debban.com', 'api.debban.com', 'app.debban.com'],
    'SHA256withRSA',
    'RSA',
    2048,
    NOW() - INTERVAL '45 days',
    NOW() + INTERVAL '45 days',
    'a1b2c3d4e5f6789012345678901234567890abcd',
    'b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567890ab',
    false,
    false,
    NOW() - INTERVAL '45 days',
    NOW()
);

-- Expiring certificate (urgent attention needed)
INSERT INTO certificates (tenant_id, serial_number, subject_dn, issuer_dn, common_name, subject_alternative_names, signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, fingerprint_sha1, fingerprint_sha256, is_self_signed, is_ca_certificate, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    '04:b2:c3:d4:e5:f6:a7:b8:c9:d0:e1:f2:a3:b4:c5:d6:e7:f8',
    'CN=db-primary.debban.internal,O=Debban Corporation,C=US',
    'CN=Debban Corporation Internal CA,O=Debban Corporation,C=US',
    'db-primary.debban.internal',
    ARRAY['db-primary.debban.internal', 'db-master.debban.internal'],
    'SHA256withRSA',
    'RSA',
    2048,
    NOW() - INTERVAL '365 days',
    NOW() + INTERVAL '3 days', -- Expiring in 3 days!
    'c2d3e4f5a6b789012345678901234567890abcde',
    'd3e4f5a6b789012345678901234567890abcdef1234567890abcdef1234567890abc',
    false,
    false,
    NOW() - INTERVAL '365 days',
    NOW()
);

-- Self-signed certificate (security issue)
INSERT INTO certificates (tenant_id, serial_number, subject_dn, issuer_dn, common_name, signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, fingerprint_sha1, fingerprint_sha256, is_self_signed, is_ca_certificate, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    '01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12',
    'CN=legacy-app.debban.internal,O=Debban Corporation,C=US',
    'CN=legacy-app.debban.internal,O=Debban Corporation,C=US', -- Same as subject (self-signed)
    'legacy-app.debban.internal',
    'SHA1withRSA', -- Weak signature algorithm
    'RSA',
    1024, -- Weak key size
    NOW() - INTERVAL '1000 days',
    NOW() + INTERVAL '100 days',
    'd4e5f6a7b8901234567890123456789012345678',
    'e5f6a7b890123456789012345678901234567890abcdef1234567890abcdef1234567890abcd',
    true,
    false,
    NOW() - INTERVAL '1000 days',
    NOW()
);

-- =================================================================
-- Crypto Implementations for Debban Corporation
-- =================================================================

-- Excellent TLS 1.3 implementation
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    (SELECT id FROM network_assets WHERE hostname = 'web-prod-01.debban.com'),
    'TLS',
    '1.3',
    'TLS_AES_256_GCM_SHA384',
    'X25519',
    'RSA-PSS',
    'AES-256-GCM',
    'SHA384',
    2048,
    (SELECT id FROM certificates WHERE common_name = '*.debban.com'),
    'passive',
    0.98,
    5, -- Very low risk
    NOW() - INTERVAL '2 months',
    NOW()
);

-- Good TLS 1.2 implementation
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    (SELECT id FROM network_assets WHERE hostname = 'api-gateway.debban.com'),
    'TLS',
    '1.2',
    'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384',
    'ECDHE',
    'RSA',
    'AES-256-GCM',
    'SHA384',
    2048,
    (SELECT id FROM certificates WHERE common_name = '*.debban.com'),
    'passive',
    0.95,
    15, -- Low risk
    NOW() - INTERVAL '1 month',
    NOW()
);

-- Database encryption
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    (SELECT id FROM network_assets WHERE hostname = 'db-primary.debban.internal'),
    'TLS',
    '1.2',
    'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384',
    'ECDHE',
    'RSA',
    'AES-256-GCM',
    'SHA384',
    2048,
    (SELECT id FROM certificates WHERE common_name = 'db-primary.debban.internal'),
    'passive',
    0.92,
    20, -- Low risk
    NOW() - INTERVAL '3 months',
    NOW()
);

-- Weak TLS implementation (legacy app)
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    (SELECT id FROM network_assets WHERE hostname = 'legacy-app.debban.internal'),
    'TLS',
    '1.0', -- Very weak version
    'TLS_RSA_WITH_RC4_128_MD5', -- Very weak cipher
    'RSA',
    'RSA',
    'RC4-128', -- Weak encryption
    'MD5', -- Weak hash
    1024, -- Weak key size
    (SELECT id FROM certificates WHERE common_name = 'legacy-app.debban.internal'),
    'active',
    0.99,
    90, -- Very high risk
    NOW() - INTERVAL '1 year',
    NOW()
);

-- =================================================================
-- Sensors for Debban Corporation
-- =================================================================

INSERT INTO sensors (id, tenant_id, name, sensor_type, deployment_location, ip_address, hostname, version, status, last_heartbeat_at, api_key_hash, created_at, updated_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440300',
    '550e8400-e29b-41d4-a716-446655440200',
    'datacenter-sensor-01',
    'network',
    'Primary Datacenter - Rack A1',
    '10.10.0.100',
    'sensor-dc1-a1.debban.internal',
    '1.2.0',
    'active',
    NOW() - INTERVAL '1 minute',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvw',
    NOW() - INTERVAL '5 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440301',
    '550e8400-e29b-41d4-a716-446655440200',
    'office-sensor-hq',
    'network',
    'Corporate HQ - Floor 12',
    '192.168.10.200',
    'sensor-hq-12f.debban.internal',
    '1.2.0',
    'active',
    NOW() - INTERVAL '3 minutes',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvx',
    NOW() - INTERVAL '4 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440302',
    '550e8400-e29b-41d4-a716-446655440200',
    'cloud-monitor-aws',
    'cloud',
    'AWS US-West-2',
    NULL,
    'cloud-sensor-aws.debban.com',
    '1.1.5',
    'active',
    NOW() - INTERVAL '2 minutes',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvy',
    NOW() - INTERVAL '3 months',
    NOW()
),
(
    '550e8400-e29b-41d4-a716-446655440303',
    '550e8400-e29b-41d4-a716-446655440200',
    'staging-sensor',
    'network',
    'Staging Environment',
    '10.20.0.50',
    'sensor-staging.debban.internal',
    '1.2.0',
    'inactive',
    NOW() - INTERVAL '1 hour',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvz',
    NOW() - INTERVAL '2 months',
    NOW()
);

-- =================================================================
-- Compliance Assessment for Debban Corporation
-- =================================================================

INSERT INTO compliance_assessments (
    tenant_id, 
    framework_id, 
    assessment_name, 
    scope_filter,
    overall_score,
    total_checks,
    passed_checks,
    failed_checks,
    not_applicable_checks,
    assessment_results,
    assessed_by,
    created_at,
    updated_at
) VALUES (
    '550e8400-e29b-41d4-a716-446655440200',
    (SELECT id FROM compliance_frameworks WHERE name = 'PCI DSS' LIMIT 1),
    'Q4 2024 PCI DSS Assessment - Debban Corp',
    '{"environments": ["production"], "asset_types": ["server", "service"], "business_units": ["Engineering", "Database Team"]}',
    78.5,
    15,
    11,
    3,
    1,
    '{
        "summary": "Good overall compliance with some critical issues on legacy systems",
        "critical_findings": [
            "Legacy application server using weak TLS 1.0 with RC4 cipher",
            "Database certificate expiring in 3 days - immediate action required",
            "Self-signed certificate on legacy application"
        ],
        "high_priority_findings": [
            "Missing certificate monitoring automation",
            "Legacy application scheduled for migration but still in production"
        ],
        "recommendations": [
            "Immediately renew database certificate",
            "Accelerate legacy application migration timeline",
            "Implement automated certificate monitoring and renewal",
            "Upgrade legacy application to TLS 1.2 minimum before migration"
        ],
        "next_assessment": "2025-04-01"
    }',
    '550e8400-e29b-41d4-a716-446655440201',
    NOW() - INTERVAL '1 week',
    NOW()
);

-- =================================================================
-- Recent Discovery Events
-- =================================================================

INSERT INTO discovery_events (tenant_id, sensor_id, asset_id, event_type, event_data, confidence_score, risk_score, created_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440300',
    (SELECT id FROM network_assets WHERE hostname = 'web-prod-01.debban.com'),
    'crypto_discovery',
    '{"protocol": "TLS", "version": "1.3", "cipher": "TLS_AES_256_GCM_SHA384", "port": 443}',
    0.98,
    5,
    NOW() - INTERVAL '2 hours'
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440300',
    (SELECT id FROM network_assets WHERE hostname = 'legacy-app.debban.internal'),
    'security_issue',
    '{"issue": "weak_tls", "protocol": "TLS", "version": "1.0", "cipher": "TLS_RSA_WITH_RC4_128_MD5", "risk_level": "critical"}',
    0.99,
    90,
    NOW() - INTERVAL '1 hour'
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440301',
    (SELECT id FROM network_assets WHERE hostname = 'db-primary.debban.internal'),
    'certificate_expiry_warning',
    '{"certificate": "db-primary.debban.internal", "expires_in_days": 3, "action_required": "immediate"}',
    1.0,
    75,
    NOW() - INTERVAL '30 minutes'
);

-- =================================================================
-- Audit Logs for Debban Corporation
-- =================================================================

INSERT INTO audit.audit_logs (tenant_id, user_id, action, resource_type, resource_id, ip_address, user_agent, success, created_at) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440201',
    'login',
    'user',
    '550e8400-e29b-41d4-a716-446655440201',
    '192.168.10.100',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    true,
    NOW() - INTERVAL '2 hours'
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440201',
    'view_dashboard',
    'dashboard',
    NULL,
    '192.168.10.100',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    true,
    NOW() - INTERVAL '1 hour'
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440202',
    'generate_compliance_report',
    'report',
    NULL,
    '192.168.10.105',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    true,
    NOW() - INTERVAL '45 minutes'
),
(
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440201',
    'view_certificate_details',
    'certificate',
    (SELECT id FROM certificates WHERE common_name = 'db-primary.debban.internal'),
    '192.168.10.100',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    true,
    NOW() - INTERVAL '30 minutes'
);

-- =================================================================
-- Success Message
-- =================================================================

DO $$
BEGIN
    RAISE NOTICE 'Sample data created successfully for brian@debban.com!';
    RAISE NOTICE 'Tenant: Debban Corporation (debban-corp)';
    RAISE NOTICE 'Users: brian@debban.com, sarah@debban.com, mike@debban.com';
    RAISE NOTICE 'Password for all users: admin123';
    RAISE NOTICE 'Created % network assets', (SELECT COUNT(*) FROM network_assets WHERE tenant_id = '550e8400-e29b-41d4-a716-446655440200');
    RAISE NOTICE 'Created % crypto implementations', (SELECT COUNT(*) FROM crypto_implementations WHERE tenant_id = '550e8400-e29b-41d4-a716-446655440200');
    RAISE NOTICE 'Created % certificates', (SELECT COUNT(*) FROM certificates WHERE tenant_id = '550e8400-e29b-41d4-a716-446655440200');
    RAISE NOTICE 'Created % sensors', (SELECT COUNT(*) FROM sensors WHERE tenant_id = '550e8400-e29b-41d4-a716-446655440200');
    RAISE NOTICE 'Created % compliance assessments', (SELECT COUNT(*) FROM compliance_assessments WHERE tenant_id = '550e8400-e29b-41d4-a716-446655440200');
END $$;
