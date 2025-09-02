# AI Agent Handoff Guide

## Document Purpose

This guide provides comprehensive context for AI agents to continue development of the crypto inventory SaaS platform. It includes all necessary information about the project scope, technical decisions, architecture, and implementation strategy.

## Project Summary

### Business Context
We are building a **SaaS platform for cryptographic inventory and compliance management**. The platform helps organizations:

1. **Discover and inventory** all cryptographic implementations across their networks
2. **Analyze compliance** against frameworks like PCI-DSS, NIST, FIPS
3. **Generate reports** for executive and compliance teams
4. **Monitor in real-time** for changes in crypto configurations
5. **Prepare for post-quantum** cryptography migration

### Target Market
- **Primary**: Enterprise organizations (1000+ employees) with complex network infrastructure
- **Secondary**: Mid-market companies, financial services, government contractors, healthcare
- **Business Model**: SaaS subscription based on endpoints monitored, with tiered pricing

### Core Value Proposition
- **Dual Discovery**: Active network probing + passive traffic analysis sensors
- **Real-time Inventory**: Continuous monitoring vs. point-in-time snapshots
- **Enterprise Features**: Multi-tenancy, SSO, RBAC, ITAM integration
- **Security First**: Isolated instances and strict data protection
- **Compliance Automation**: Framework-specific analysis and automated reporting

## Architecture Overview

### System Components
```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Web Frontend  │  │  Mobile/CLI     │  │ 3rd Party APIs  │
│    (React)      │  │    Clients      │  │   (ITAM/SSO)    │
└─────────┬───────┘  └─────────┬───────┘  └─────────┬───────┘
          │                    │                    │
          └────────────────────┼────────────────────┘
                               │
                    ┌─────────────────┐
                    │   API Gateway   │
                    │     (NGINX)     │
                    └─────────┬───────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
    ┌─────────▼───────┐ ┌─────▼─────┐ ┌─────▼─────┐
    │ Auth Service    │ │Inventory  │ │Compliance │
    │     (Go)        │ │Service(Go)│ │Engine(Go) │
    └─────────────────┘ └───────────┘ └───────────┘
              │               │               │
              └───────────────┼───────────────┘
                              │
                    ┌─────────▼───────┐
                    │  Data Layer     │
                    │ PostgreSQL +    │
                    │ InfluxDB +      │
                    │ Redis + NATS    │
                    └─────────────────┘
```

### Service Boundaries
1. **Authentication Service**: User auth, SSO, RBAC, JWT token management
2. **Inventory Service**: Asset discovery, crypto implementation tracking, certificates
3. **Compliance Engine**: Framework analysis, gap assessment, scoring
4. **Report Service**: PDF/Excel generation, scheduled reports, dashboards
5. **Sensor Manager**: Network sensor coordination, data ingestion, validation
6. **Integration Service**: ITAM system connectors and data synchronization
7. **Network Sensor**: Cross-platform Go binary for passive traffic analysis

### Technology Stack
- **Backend**: Go (Gin framework) for all services
- **AI/ML**: Python + TensorFlow/PyTorch for AI analysis service
- **Network Sensor**: Go with cross-platform deployment (Windows/Linux/containers)
- **Frontend**: React + TypeScript + Ant Design
- **Databases**: PostgreSQL (relational), InfluxDB (time-series), Redis (cache)
- **Infrastructure**: Docker + Kubernetes + Terraform
- **Messaging**: NATS for inter-service communication
- **Monitoring**: Prometheus + Grafana + ELK stack

## Key Technical Decisions

### Why Go for Backend & Sensors?
- **Performance**: Compiled language with excellent concurrency
- **Cross-Platform Deployment**: Single codebase compiles to Windows, Linux, macOS, ARM
- **Flexible Sensor Deployment**: Native binaries, Windows Services, Linux systemd, containers
- **Talent Pool**: Large developer community for easy handoff
- **Network Programming**: Excellent standard library for network operations and packet analysis
- **Memory Efficiency**: Low footprint ideal for agent deployment across enterprise networks

### Why React + TypeScript?
- **Enterprise Adoption**: Widely used in enterprise applications
- **Component Ecosystem**: Rich library of enterprise-grade UI components
- **Type Safety**: TypeScript reduces runtime errors
- **Talent Pool**: Largest frontend developer community

### Multi-Tenancy Strategy
- **Kubernetes Namespaces**: Logical isolation per tenant
- **Database Isolation**: Separate PostgreSQL instances for enterprise customers
- **Network Policies**: Tenant traffic isolation within cluster
- **Resource Quotas**: CPU/memory limits per tenant

### Security Approach
- **Zero Trust**: Never trust, always verify
- **Encryption Everywhere**: TLS 1.3 for all communications, encryption at rest
- **Isolated Deployment**: Prefer dedicated instances over shared tenancy
- **RBAC**: Granular role-based access control
- **Audit Everything**: Comprehensive security event logging

## Data Models

### Core Entities
1. **Tenants**: Multi-tenant isolation, subscription management
2. **Users**: Authentication, roles, SSO integration
3. **Network Assets**: Discovered servers, endpoints, services
4. **Crypto Implementations**: TLS configs, cipher suites, key algorithms
5. **Certificates**: X.509 certificate tracking with expiration monitoring
6. **Sensors**: Network monitoring agents and their configurations
7. **Compliance Assessments**: Framework analysis results and scoring
8. **Reports**: Generated compliance and inventory reports

### Database Schema Highlights
```sql
-- Multi-tenant isolation on all tables
CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    subscription_tier VARCHAR(50) DEFAULT 'basic'
);

-- Core crypto discovery
CREATE TABLE crypto_implementations (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    asset_id UUID NOT NULL REFERENCES network_assets(id),
    protocol VARCHAR(50) NOT NULL, -- 'TLS', 'SSH', 'IPSec'
    cipher_suite VARCHAR(255),
    key_size INTEGER,
    confidence_score DECIMAL(3,2), -- 0.0 to 1.0
    discovery_method VARCHAR(50) -- 'passive', 'active', 'manual'
);
```

### Time-Series Data (InfluxDB)
- **Sensor Metrics**: CPU, memory, network packets analyzed
- **Discovery Events**: Timestamped crypto discoveries
- **Platform Metrics**: API response times, queue depths
- **Audit Logs**: Security events with full context

## API Design

### RESTful Conventions
- **Base URL**: `https://api.cryptoinventory.com/v1/{service}/{resource}`
- **Authentication**: JWT Bearer tokens + API keys
- **Rate Limiting**: Per-tenant quotas (1000/hour default, 10000/hour enterprise)
- **Versioning**: URL-based versioning for backward compatibility

### Key Endpoints
```
# Authentication
POST /v1/auth/login
POST /v1/auth/sso/initiate
POST /v1/auth/refresh

# Inventory Management
GET  /v1/inventory/assets?search=hostname&environment=prod
POST /v1/inventory/assets
GET  /v1/inventory/crypto-implementations?weak_crypto=true
GET  /v1/inventory/certificates?expiring_days=30

# Compliance Analysis
GET  /v1/compliance/frameworks
POST /v1/compliance/assessments
GET  /v1/compliance/assessments/{id}

# Sensor Management
POST /v1/sensors/register
POST /v1/sensors/{id}/heartbeat
POST /v1/sensors/data/ingest

# Reporting
POST /v1/reports/generate
GET  /v1/reports/{id}/download
```

## Implementation Strategy

### Development Phases

#### Phase 1: MVP Core (Weeks 1-8)
- [ ] Basic multi-tenant authentication system
- [ ] Asset discovery and crypto implementation tracking
- [ ] Cross-platform network sensor (Windows/Linux native binaries)
- [ ] Simple web interface for inventory viewing
- [ ] PostgreSQL schema and data models
- [ ] Basic anomaly detection (simple ML models)

#### Phase 2: Compliance Engine + AI Foundation (Weeks 9-12)
- [ ] PCI-DSS compliance framework implementation
- [ ] NIST cybersecurity framework support
- [ ] Compliance assessment engine with AI risk scoring
- [ ] Gap analysis and remediation recommendations
- [ ] Basic PDF report generation with AI insights
- [ ] AI analysis service infrastructure

#### Phase 3: Enterprise Features + Integration Hub (Weeks 13-16)
- [ ] SSO integration (SAML/OIDC)
- [ ] Role-based access control
- [ ] Advanced sensor management (container deployment, policy management)
- [ ] Integration Hub with ITAM connectors (ServiceNow, Lansweeper)
- [ ] Real-time dashboards with AI-powered insights
- [ ] Natural language report generation
- [ ] Predictive compliance analysis

#### Phase 4: Scale, Polish + Advanced AI (Weeks 17-20)
- [ ] Kubernetes deployment automation with DaemonSet sensor deployment
- [ ] Advanced monitoring and alerting with AI anomaly detection
- [ ] Performance optimization and sensor edge AI capabilities
- [ ] Security hardening and AI model security
- [ ] Load testing and scalability validation
- [ ] MLOps pipeline for continuous model improvement

### Folder Structure
```
crypto-inventory-platform/
├── services/
│   ├── auth-service/           # Authentication & authorization
│   ├── inventory-service/      # Asset & crypto inventory
│   ├── compliance-engine/      # Framework compliance analysis
│   ├── report-generator/       # PDF/Excel report generation
│   ├── sensor-manager/         # Sensor coordination
│   ├── integration-service/    # ITAM system connectors
│   └── ai-analysis-service/    # AI/ML analysis and inference
├── sensor/                     # Network sensor Go binary
├── web-ui/                     # React frontend application
├── infrastructure/             # Terraform infrastructure modules
├── k8s/                        # Kubernetes deployment manifests
├── docs/                       # Architecture and API documentation
├── scripts/                    # Deployment and utility scripts
├── tests/                      # Integration and E2E tests
└── docker-compose.yml          # Local development environment
```

### Development Workflow
1. **Feature Branch**: Create branch from `develop`
2. **Development**: Write code with tests
3. **Code Review**: Pull request review process
4. **Staging Deploy**: Automatic deployment to staging on merge to `develop`
5. **Production Deploy**: Manual deployment from `main` branch

## Security Requirements

### Compliance Standards
- **SOC 2 Type II**: Security controls and monitoring
- **GDPR/CCPA**: Data privacy and user rights
- **ISO 27001**: Information security management
- **NIST Cybersecurity Framework**: Security controls implementation

### Security Controls
1. **Encryption**: TLS 1.3 minimum, AES-256 encryption at rest
2. **Authentication**: Multi-factor authentication required for admin access
3. **Authorization**: Principle of least privilege, granular RBAC
4. **Network Security**: VPN access for sensors, network segmentation
5. **Monitoring**: Real-time security event monitoring and alerting
6. **Incident Response**: Documented procedures for security incidents

### Data Protection
- **Data Classification**: Public, internal, confidential, restricted
- **Retention Policies**: Automated data lifecycle management
- **Backup Strategy**: Encrypted backups with cross-region replication
- **Data Anonymization**: Remove PII from analytics and logs

## Deployment Strategy

### Environment Progression
1. **Local Development**: Docker Compose for full stack
2. **Staging**: Kubernetes cluster with production-like configuration
3. **Production**: Multi-AZ Kubernetes with high availability

### Infrastructure as Code
- **Terraform**: Cloud resource provisioning
- **Helm Charts**: Kubernetes application deployment
- **GitOps**: ArgoCD for automated deployments
- **Monitoring**: Prometheus + Grafana stack

### Multi-Cloud Strategy
- **Primary**: AWS (EKS, RDS, ElastiCache)
- **Backup**: Azure or GCP for disaster recovery
- **Abstraction**: Terraform providers for cloud portability

## Integration Points

### External Systems
1. **SSO Providers**: Okta, Azure AD, Google Workspace
2. **ITAM Systems**: ServiceNow, Lansweeper, Device42
3. **SIEM Platforms**: Splunk, QRadar, Sentinel
4. **Ticketing**: Jira, ServiceNow, PagerDuty

### API Integration Strategy
- **Webhook Support**: Real-time event notifications
- **Rate Limiting**: Protect against abusive usage
- **API Versioning**: Maintain backward compatibility
- **SDK Generation**: Auto-generated client libraries

## Performance Requirements

### Scalability Targets
- **Concurrent Users**: 1000+ simultaneous users
- **Assets Monitored**: 100,000+ network assets per tenant
- **Data Ingestion**: 10,000+ crypto discoveries per minute
- **Query Performance**: <2 second response for complex inventory queries

### Monitoring KPIs
- **Uptime**: 99.9% availability SLA
- **Response Time**: 95th percentile <500ms for API calls
- **Discovery Accuracy**: >95% true positive rate
- **Sensor Performance**: <1% CPU usage impact on monitored systems

## Troubleshooting Common Issues

### Development Environment
```bash
# Services not starting
docker-compose down && docker-compose up -d

# Database connection issues
docker-compose exec postgres pg_isready -U crypto_user

# Frontend build errors
cd web-ui && rm -rf node_modules && npm install

# Go module issues
go mod tidy && go mod vendor
```

### Production Issues
```bash
# Check pod status
kubectl get pods -n crypto-inventory-production

# View service logs
kubectl logs -f deployment/inventory-service -n crypto-inventory-production

# Database connection issues
kubectl exec -it postgres-0 -- pg_isready

# Certificate expiration
kubectl get certificates -A
```

## Next Steps for AI Agent

### Immediate Actions
1. **Review Architecture**: Understand the overall system design and component interactions
2. **Set Up Environment**: Start with local Docker Compose development environment
3. **Create Project Structure**: Initialize the folder structure and basic services
4. **Implement Auth Service**: Begin with authentication and tenant management
5. **Cross-Platform Sensor**: Develop Go-based sensor with Windows/Linux deployment
6. **Database Schema**: Create PostgreSQL schema with migration scripts

### Priority Implementation Order
1. **Authentication & Multi-tenancy** (Critical foundation)
2. **Cross-Platform Network Sensor** (Core differentiator - flexible deployment)
3. **Basic Asset Discovery** (Core value proposition)
4. **Executive Dashboard UI** (Impressive demo interface)
5. **Integration Hub Foundation** (Enterprise differentiation)
6. **Basic AI Analysis** (Anomaly detection and risk scoring)
7. **Compliance Framework** (Key business value)

### Development Principles
- **Test-Driven Development**: Write tests before implementation
- **Security First**: Implement security controls from the beginning
- **Documentation**: Update docs as code is written
- **Incremental Value**: Each phase should deliver working functionality
- **Performance Awareness**: Monitor and optimize from early stages

### Questions to Ask
- Should we start with a specific compliance framework (PCI-DSS vs NIST)?
- What's the preferred cloud provider for initial deployment?
- Are there specific enterprise customers to target for beta testing?
- What's the timeline for MVP delivery?
- Should we prioritize Windows or Linux sensor deployment first?
- Which ITAM integration should we build first (ServiceNow is most common)?
- What level of AI sophistication for the initial release (simple anomaly detection vs advanced NLP)?
- Should sensors include edge AI capabilities or rely entirely on cloud analysis?
- How sophisticated should the initial UI be (basic dashboard vs full executive interface)?

---

**This document provides complete context for continuing development. The architecture is designed for handoff-ready code with comprehensive documentation, standard patterns, and scalable infrastructure.**
