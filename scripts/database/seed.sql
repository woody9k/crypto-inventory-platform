-- =================================================================
-- Crypto Inventory Platform - Development Seed Data
-- =================================================================

-- =================================================================
-- Compliance Frameworks
-- =================================================================

-- PCI DSS Framework
INSERT INTO compliance_frameworks (id, name, version, description, organization, effective_date, rules, active) VALUES (
    uuid_generate_v4(),
    'PCI DSS',
    '4.0',
    'Payment Card Industry Data Security Standard',
    'PCI Security Standards Council',
    '2022-03-31',
    '{
        "requirements": [
            {
                "id": "4.1.1",
                "title": "Strong Cryptography and Security Protocols",
                "description": "Use strong cryptography and security protocols to safeguard sensitive cardholder data during transmission over open, public networks",
                "checks": [
                    {"cipher_suites": {"prohibited": ["TLS_RSA_WITH_RC4_128_MD5", "TLS_RSA_WITH_RC4_128_SHA"]}},
                    {"tls_versions": {"minimum": "1.2", "prohibited": ["SSLv2", "SSLv3", "TLSv1.0", "TLSv1.1"]}},
                    {"key_sizes": {"minimum_rsa": 2048, "minimum_ecc": 224}}
                ]
            },
            {
                "id": "4.2.1",
                "title": "Never Send Unprotected PANs",
                "description": "Never send unprotected PANs by end-user messaging technologies",
                "checks": [
                    {"protocols": {"required_encryption": ["email", "instant_messaging", "SMS", "chat"]}}
                ]
            }
        ]
    }',
    true
) ON CONFLICT (name, version) DO NOTHING;

-- NIST Cybersecurity Framework
INSERT INTO compliance_frameworks (id, name, version, description, organization, effective_date, rules, active) VALUES (
    uuid_generate_v4(),
    'NIST CSF',
    '1.1',
    'NIST Cybersecurity Framework',
    'National Institute of Standards and Technology',
    '2018-04-16',
    '{
        "categories": [
            {
                "id": "PR.DS-2",
                "title": "Data-in-transit is protected",
                "description": "Information is protected in transit",
                "checks": [
                    {"encryption_in_transit": {"required": true}},
                    {"weak_protocols": {"prohibited": ["HTTP", "FTP", "Telnet", "SNMP v1/v2"]}},
                    {"certificate_validation": {"required": true}}
                ]
            },
            {
                "id": "PR.DS-1",
                "title": "Data-at-rest is protected",
                "description": "Information is protected at rest",
                "checks": [
                    {"encryption_at_rest": {"required": true}},
                    {"key_management": {"rotation_required": true}}
                ]
            }
        ]
    }',
    true
) ON CONFLICT (name, version) DO NOTHING;

-- =================================================================
-- Demo Tenant and Users
-- =================================================================

-- Create demo tenant
INSERT INTO tenants (id, name, slug, subscription_tier, max_endpoints, max_users) VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    'Demo Corporation',
    'demo-corp',
    'enterprise',
    1000,
    50
) ON CONFLICT (slug) DO NOTHING;

-- Create demo users
INSERT INTO users (id, tenant_id, email, first_name, last_name, password_hash, role, email_verified, active) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440000',
    'admin@democorp.com',
    'Admin',
    'User',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye/D7zrZI/PCMZ6qO8PQ8DbZOF5.XzEQm', -- password: admin123
    'admin',
    true,
    true
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440000',
    'analyst@democorp.com',
    'Security',
    'Analyst',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye/D7zrZI/PCMZ6qO8PQ8DbZOF5.XzEQm', -- password: admin123
    'analyst',
    true,
    true
),
(
    '550e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440000',
    'viewer@democorp.com',
    'Read Only',
    'User',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye/D7zrZI/PCMZ6qO8PQ8DbZOF5.XzEQm', -- password: admin123
    'viewer',
    true,
    true
) ON CONFLICT (tenant_id, email) DO NOTHING;

-- =================================================================
-- Demo Network Assets
-- =================================================================

-- Web servers
INSERT INTO network_assets (tenant_id, hostname, ip_address, port, asset_type, operating_system, environment, business_unit, owner_email, description, tags) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    'web-prod-01.democorp.com',
    '10.1.1.10',
    443,
    'server',
    'Ubuntu 22.04',
    'production',
    'IT Infrastructure',
    'admin@democorp.com',
    'Primary production web server',
    '{"service": "web", "critical": true, "public_facing": true}'
),
(
    '550e8400-e29b-41d4-a716-446655440000',
    'web-prod-02.democorp.com',
    '10.1.1.11',
    443,
    'server',
    'Ubuntu 22.04',
    'production',
    'IT Infrastructure',
    'admin@democorp.com',
    'Secondary production web server',
    '{"service": "web", "critical": true, "public_facing": true}'
),
(
    '550e8400-e29b-41d4-a716-446655440000',
    'api-gateway.democorp.com',
    '10.1.2.10',
    443,
    'service',
    'Alpine Linux',
    'production',
    'IT Infrastructure',
    'admin@democorp.com',
    'API Gateway service',
    '{"service": "api", "critical": true, "public_facing": true}'
),
(
    '550e8400-e29b-41d4-a716-446655440000',
    'db-primary.democorp.internal',
    '10.1.3.10',
    5432,
    'server',
    'PostgreSQL on Ubuntu',
    'production',
    'Database Team',
    'dba@democorp.com',
    'Primary database server',
    '{"service": "database", "critical": true, "encrypted": true}'
),
(
    '550e8400-e29b-41d4-a716-446655440000',
    'legacy-ftp.democorp.internal',
    '10.1.4.10',
    21,
    'server',
    'Windows Server 2016',
    'production',
    'Legacy Systems',
    'legacy@democorp.com',
    'Legacy FTP server (scheduled for decommission)',
    '{"service": "ftp", "legacy": true, "decommission_planned": true}'
);

-- =================================================================
-- Demo Certificates
-- =================================================================

-- Valid certificate
INSERT INTO certificates (tenant_id, serial_number, subject_dn, issuer_dn, common_name, subject_alternative_names, signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, fingerprint_sha1, fingerprint_sha256, is_self_signed, is_ca_certificate) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    '03:e8:6a:c3:cf:75:8d:c4:8b:5e:2d:17:2a:9b:1c:82:4a:5f',
    'CN=*.democorp.com,O=Demo Corporation,L=San Francisco,ST=California,C=US',
    'CN=Let''s Encrypt Authority X3,O=Let''s Encrypt,C=US',
    '*.democorp.com',
    ARRAY['democorp.com', 'www.democorp.com', 'api.democorp.com'],
    'SHA256withRSA',
    'RSA',
    2048,
    NOW() - INTERVAL '30 days',
    NOW() + INTERVAL '60 days',
    '2f:5d:9b:8c:7a:3e:4f:1d:6c:8a:9b:7e:5f:2d:4c:8b:1a:9e:7f:3c',
    'a1:b2:c3:d4:e5:f6:7a:8b:9c:0d:1e:2f:3a:4b:5c:6d:7e:8f:9a:0b:1c:2d:3e:4f:5a:6b:7c:8d:9e:0f:1a:2b',
    false,
    false
);

-- Expiring certificate
INSERT INTO certificates (tenant_id, serial_number, subject_dn, issuer_dn, common_name, signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, fingerprint_sha1, fingerprint_sha256, is_self_signed, is_ca_certificate) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    '04:a1:7b:2c:3d:4e:5f:6a:7b:8c:9d:0e:1f:2a:3b:4c:5d:6e',
    'CN=db-primary.democorp.internal,O=Demo Corporation,C=US',
    'CN=Demo Corporation Internal CA,O=Demo Corporation,C=US',
    'db-primary.democorp.internal',
    'SHA256withRSA',
    'RSA',
    2048,
    NOW() - INTERVAL '360 days',
    NOW() + INTERVAL '5 days', -- Expiring soon!
    '3a:6b:7c:8d:9e:0f:1a:2b:3c:4d:5e:6f:7a:8b:9c:0d:1e:2f:3a:4b',
    'b2:c3:d4:e5:f6:7a:8b:9c:0d:1e:2f:3a:4b:5c:6d:7e:8f:9a:0b:1c:2d:3e:4f:5a:6b:7c:8d:9e:0f:1a:2b:3c',
    false,
    false
);

-- Self-signed certificate (security issue)
INSERT INTO certificates (tenant_id, serial_number, subject_dn, issuer_dn, common_name, signature_algorithm, public_key_algorithm, public_key_size, not_before, not_after, fingerprint_sha1, fingerprint_sha256, is_self_signed, is_ca_certificate) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    '01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12',
    'CN=legacy-ftp.democorp.internal,O=Demo Corporation,C=US',
    'CN=legacy-ftp.democorp.internal,O=Demo Corporation,C=US', -- Same as subject (self-signed)
    'legacy-ftp.democorp.internal',
    'SHA1withRSA', -- Weak signature algorithm
    'RSA',
    1024, -- Weak key size
    NOW() - INTERVAL '1000 days',
    NOW() + INTERVAL '100 days',
    '4b:5c:6d:7e:8f:9a:0b:1c:2d:3e:4f:5a:6b:7c:8d:9e:0f:1a:2b:3c',
    'c3:d4:e5:f6:7a:8b:9c:0d:1e:2f:3a:4b:5c:6d:7e:8f:9a:0b:1c:2d:3e:4f:5a:6b:7c:8d:9e:0f:1a:2b:3c:4d',
    true,
    false
);

-- =================================================================
-- Demo Crypto Implementations
-- =================================================================

-- Good TLS implementation
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    (SELECT id FROM network_assets WHERE hostname = 'web-prod-01.democorp.com'),
    'TLS',
    '1.3',
    'TLS_AES_256_GCM_SHA384',
    'X25519',
    'RSA-PSS',
    'AES-256-GCM',
    'SHA384',
    2048,
    (SELECT id FROM certificates WHERE common_name = '*.democorp.com'),
    'passive',
    0.95,
    10
);

-- Weak TLS implementation
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    (SELECT id FROM network_assets WHERE hostname = 'legacy-ftp.democorp.internal'),
    'TLS',
    '1.0', -- Weak version
    'TLS_RSA_WITH_RC4_128_MD5', -- Weak cipher
    'RSA',
    'RSA',
    'RC4-128', -- Weak encryption
    'MD5', -- Weak hash
    1024, -- Weak key size
    (SELECT id FROM certificates WHERE common_name = 'legacy-ftp.democorp.internal'),
    'active',
    0.99,
    85 -- High risk
);

-- Database encryption
INSERT INTO crypto_implementations (tenant_id, asset_id, protocol, protocol_version, cipher_suite, key_exchange_algorithm, signature_algorithm, symmetric_encryption, hash_algorithm, key_size, certificate_id, discovery_method, confidence_score, risk_score) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    (SELECT id FROM network_assets WHERE hostname = 'db-primary.democorp.internal'),
    'TLS',
    '1.2',
    'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384',
    'ECDHE',
    'RSA',
    'AES-256-GCM',
    'SHA384',
    2048,
    (SELECT id FROM certificates WHERE common_name = 'db-primary.democorp.internal'),
    'passive',
    0.92,
    25
);

-- =================================================================
-- Demo Sensors
-- =================================================================

INSERT INTO sensors (id, tenant_id, name, sensor_type, deployment_location, ip_address, hostname, version, status, last_heartbeat_at, api_key_hash) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440100',
    '550e8400-e29b-41d4-a716-446655440000',
    'datacenter-sensor-01',
    'network',
    'Primary Datacenter - Rack A1',
    '10.1.0.100',
    'sensor-dc1-a1.democorp.internal',
    '1.0.0',
    'active',
    NOW() - INTERVAL '2 minutes',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvw'
),
(
    '550e8400-e29b-41d4-a716-446655440101',
    '550e8400-e29b-41d4-a716-446655440000',
    'office-sensor-hq',
    'network',
    'Corporate HQ - Floor 5',
    '192.168.1.200',
    'sensor-hq-5f.democorp.internal',
    '1.0.0',
    'active',
    NOW() - INTERVAL '5 minutes',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvx'
),
(
    '550e8400-e29b-41d4-a716-446655440102',
    '550e8400-e29b-41d4-a716-446655440000',
    'cloud-monitor-aws',
    'cloud',
    'AWS US-West-2',
    NULL,
    'cloud-sensor-aws.democorp.com',
    '1.0.0',
    'inactive',
    NOW() - INTERVAL '2 hours',
    '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvz'
);

-- Update crypto implementations with sensor references
UPDATE crypto_implementations 
SET source_sensor_id = '550e8400-e29b-41d4-a716-446655440100'
WHERE asset_id IN (
    SELECT id FROM network_assets 
    WHERE hostname IN ('web-prod-01.democorp.com', 'db-primary.democorp.internal')
);

UPDATE crypto_implementations 
SET source_sensor_id = '550e8400-e29b-41d4-a716-446655440101'
WHERE asset_id IN (
    SELECT id FROM network_assets 
    WHERE hostname = 'legacy-ftp.democorp.internal'
);

-- =================================================================
-- Demo AI Models
-- =================================================================

INSERT INTO ai_models (name, version, model_type, description, active, hyperparameters, metrics) VALUES 
(
    'Crypto Anomaly Detector',
    '1.0.0',
    'anomaly_detection',
    'Machine learning model for detecting unusual cryptographic configurations',
    true,
    '{"algorithm": "isolation_forest", "contamination": 0.1, "n_estimators": 100}',
    '{"accuracy": 0.92, "precision": 0.89, "recall": 0.94, "f1_score": 0.91}'
),
(
    'Risk Scoring Engine',
    '1.2.1',
    'risk_scoring',
    'AI model for calculating cryptographic risk scores based on multiple factors',
    true,
    '{"algorithm": "gradient_boosting", "learning_rate": 0.1, "max_depth": 6, "n_estimators": 200}',
    '{"mse": 0.15, "r2_score": 0.87, "mae": 0.23}'
),
(
    'Report Generator NLP',
    '0.9.0',
    'nlp',
    'Natural language processing model for generating executive reports',
    false,
    '{"model": "transformer", "max_length": 1024, "temperature": 0.7}',
    '{"bleu_score": 0.78, "rouge_l": 0.82, "perplexity": 12.4}'
);

-- =================================================================
-- Demo Compliance Assessment
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
    assessed_by
) VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    (SELECT id FROM compliance_frameworks WHERE name = 'PCI DSS' LIMIT 1),
    'Q4 2023 PCI DSS Assessment',
    '{"environments": ["production"], "asset_types": ["server", "service"]}',
    73.5,
    12,
    8,
    3,
    1,
    '{
        "summary": "Mixed compliance status with critical issues on legacy systems",
        "critical_findings": [
            "Legacy FTP server using weak TLS 1.0 with RC4 cipher",
            "Database certificate expiring in 5 days"
        ],
        "recommendations": [
            "Upgrade legacy FTP server to TLS 1.2 minimum",
            "Renew database certificate immediately",
            "Implement automated certificate renewal"
        ]
    }',
    '550e8400-e29b-41d4-a716-446655440001'
);

-- =================================================================
-- Demo Audit Logs
-- =================================================================

INSERT INTO audit.audit_logs (tenant_id, user_id, action, resource_type, resource_id, ip_address, user_agent, success) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    '550e8400-e29b-41d4-a716-446655440001',
    'login',
    'user',
    '550e8400-e29b-41d4-a716-446655440001',
    '192.168.1.100',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    true
),
(
    '550e8400-e29b-41d4-a716-446655440000',
    '550e8400-e29b-41d4-a716-446655440001',
    'view_compliance_assessment',
    'compliance_assessment',
    (SELECT id FROM compliance_assessments LIMIT 1),
    '192.168.1.100',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    true
),
(
    '550e8400-e29b-41d4-a716-446655440000',
    '550e8400-e29b-41d4-a716-446655440002',
    'generate_report',
    'report',
    NULL,
    '192.168.1.105',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    true
);

-- =================================================================
-- Success Message
-- =================================================================

DO $$
BEGIN
    RAISE NOTICE 'Database seeded successfully with demo data!';
    RAISE NOTICE 'Demo tenant: demo-corp';
    RAISE NOTICE 'Demo users: admin@democorp.com, analyst@democorp.com, viewer@democorp.com';
    RAISE NOTICE 'Password for all demo users: admin123';
    RAISE NOTICE 'Created % network assets', (SELECT COUNT(*) FROM network_assets);
    RAISE NOTICE 'Created % crypto implementations', (SELECT COUNT(*) FROM crypto_implementations);
    RAISE NOTICE 'Created % certificates', (SELECT COUNT(*) FROM certificates);
    RAISE NOTICE 'Created % sensors', (SELECT COUNT(*) FROM sensors);
END $$;
