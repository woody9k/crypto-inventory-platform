#!/bin/bash

# GitHub Issues Creation Script
# This script creates the initial development issues for the Crypto Inventory Platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if GitHub CLI is installed
check_gh_cli() {
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI (gh) is not installed"
        echo "Please install GitHub CLI: https://cli.github.com/"
        exit 1
    fi
    
    # Check if authenticated
    if ! gh auth status &> /dev/null; then
        log_error "GitHub CLI is not authenticated"
        echo "Please run: gh auth login"
        exit 1
    fi
    
    log_success "GitHub CLI is ready"
}

# Create milestones
create_milestones() {
    log_info "Creating project milestones..."
    
    milestones=(
        "Phase 1: Foundation|MVP Core - Authentication, Sensors, Basic Discovery|2024-03-01"
        "Phase 2: Intelligence|AI Analysis and Compliance Framework|2024-04-01"
        "Phase 3: Enterprise|Integration Hub and Advanced Features|2024-05-01"
        "Phase 4: Scale|Production Ready and Polish|2024-06-01"
    )
    
    for milestone_data in "${milestones[@]}"; do
        IFS='|' read -r title description due_date <<< "$milestone_data"
        
        log_info "Creating milestone: $title"
        gh api repos/:owner/:repo/milestones \
            --method POST \
            --field title="$title" \
            --field description="$description" \
            --field due_on="$due_date" \
            --field state="open" || log_warning "Failed to create milestone: $title"
    done
    
    log_success "Milestones created"
}

# Create labels
create_labels() {
    log_info "Creating project labels..."
    
    labels=(
        "mvp|Critical for MVP|e11d21"
        "authentication|Authentication Service|0052cc"
        "inventory|Inventory Service|1d76db"
        "compliance|Compliance Engine|5319e7"
        "reports|Report Generator|f9d0c4"
        "sensors|Sensor Manager|c2e0c6"
        "integration|Integration Service|0e8a16"
        "ai|AI Analysis Service|8b5cf6"
        "frontend|Web Frontend|fef2c0"
        "database|Database Related|bfd4f2"
        "infrastructure|Infrastructure & DevOps|d4c5f9"
        "security|Security Related|b60205"
        "performance|Performance Related|0e8a16"
        "documentation|Documentation|c5def5"
        "critical|Critical Priority|b60205"
        "high|High Priority|d93f0b"
        "medium|Medium Priority|fbca04"
        "low|Low Priority|0e8a16"
    )
    
    for label_data in "${labels[@]}"; do
        IFS='|' read -r name description color <<< "$label_data"
        
        log_info "Creating label: $name"
        gh api repos/:owner/:repo/labels \
            --method POST \
            --field name="$name" \
            --field description="$description" \
            --field color="$color" || log_warning "Label $name might already exist"
    done
    
    log_success "Labels created"
}

# Create Phase 1 issues
create_phase1_issues() {
    log_info "Creating Phase 1 (Foundation) issues..."
    
    # Issue 1: Authentication Service MVP
    gh issue create \
        --title "[TASK] Implement Authentication Service MVP" \
        --body "**Component**: Authentication Service
**Priority**: Critical (MVP blocker)
**Milestone**: Phase 1: Foundation

## Description
Develop core authentication service with JWT tokens, multi-tenant support, and basic user management.

## Acceptance Criteria
- [ ] JWT token generation and validation
- [ ] Multi-tenant user isolation
- [ ] Basic RBAC (admin, analyst, viewer roles)
- [ ] Password hashing and security
- [ ] Health check endpoints
- [ ] Database integration with user/tenant tables
- [ ] API documentation

## Technical Requirements
- Use Go with Gin framework
- PostgreSQL for user/tenant storage
- Redis for session management
- JWT with RS256 signing
- Input validation and error handling

## Estimate
1 week" \
        --label "task,authentication,critical,mvp" \
        --milestone "Phase 1: Foundation"

    # Issue 2: Cross-Platform Network Sensor
    gh issue create \
        --title "[TASK] Build Cross-Platform Network Sensor Foundation" \
        --body "**Component**: Network Sensor
**Priority**: Critical (MVP blocker)
**Milestone**: Phase 1: Foundation

## Description
Develop Go-based network sensor with cross-platform deployment capabilities (Windows, Linux, containers).

## Acceptance Criteria
- [ ] Cross-platform Go binary compilation (Windows, Linux, macOS, ARM)
- [ ] Basic network packet capture
- [ ] TLS handshake detection and analysis
- [ ] Configuration file support
- [ ] Daemon/service mode operation
- [ ] Health status reporting
- [ ] Secure communication with platform

## Technical Requirements
- Single Go binary with no external dependencies
- Support for Windows Service and Linux systemd
- Docker container deployment option
- Configuration via YAML file
- mTLS authentication with platform

## Estimate
1.5 weeks" \
        --label "task,sensors,critical,mvp" \
        --milestone "Phase 1: Foundation"

    # Issue 3: Asset Discovery Service
    gh issue create \
        --title "[TASK] Implement Asset Discovery Service" \
        --body "**Component**: Inventory Service
**Priority**: Critical (MVP blocker)
**Milestone**: Phase 1: Foundation

## Description
Build core asset discovery and crypto implementation tracking system.

## Acceptance Criteria
- [ ] Asset registration and management API
- [ ] Crypto implementation storage and tracking
- [ ] Certificate tracking with expiration monitoring
- [ ] Basic search and filtering capabilities
- [ ] RESTful API endpoints with OpenAPI docs
- [ ] Multi-tenant data isolation

## Technical Requirements
- Go with Gin framework
- PostgreSQL for data storage
- Proper indexing for performance
- Input validation and sanitization
- Comprehensive error handling

## Estimate
1 week" \
        --label "task,inventory,critical,mvp" \
        --milestone "Phase 1: Foundation"

    # Issue 4: Executive Dashboard UI
    gh issue create \
        --title "[TASK] Build Main Dashboard Interface" \
        --body "**Component**: Web Frontend
**Priority**: High (Important for MVP)
**Milestone**: Phase 1: Foundation

## Description
Create impressive executive dashboard with KPIs, risk heat map, and real-time updates.

## Acceptance Criteria
- [ ] Key metrics cards (sensors, assets, compliance score)
- [ ] Risk heat map visualization
- [ ] Real-time activity feed
- [ ] Quick action buttons
- [ ] Responsive design for mobile/tablet
- [ ] Role-based navigation

## Technical Requirements
- React with TypeScript
- Ant Design components
- WebSocket for real-time updates
- Responsive CSS
- Accessibility compliance

## Estimate
1.5 weeks" \
        --label "task,frontend,high,mvp" \
        --milestone "Phase 1: Foundation"

    # Issue 5: Database Schema
    gh issue create \
        --title "[TASK] Create Database Schema and Migrations" \
        --body "**Component**: Database
**Priority**: Critical (MVP blocker)
**Milestone**: Phase 1: Foundation

## Description
Implement comprehensive PostgreSQL schema with proper indexing and constraints.

## Acceptance Criteria
- [ ] All tables created with proper relationships
- [ ] Indexes for performance optimization
- [ ] Migration scripts for schema versioning
- [ ] Seed data for development
- [ ] Backup and restore procedures
- [ ] Multi-tenant isolation at database level

## Technical Requirements
- PostgreSQL 15+
- Proper foreign key relationships
- Performance-optimized indexes
- Data validation constraints
- Audit logging capabilities

## Estimate
4 days" \
        --label "task,database,critical,mvp" \
        --milestone "Phase 1: Foundation"

    log_success "Phase 1 issues created"
}

# Create Phase 2 issues
create_phase2_issues() {
    log_info "Creating Phase 2 (Intelligence) issues..."
    
    # Issue 6: AI Analysis Service
    gh issue create \
        --title "[TASK] Implement AI Analysis Service Infrastructure" \
        --body "**Component**: AI Analysis Service
**Priority**: High (Important for MVP)
**Milestone**: Phase 2: Intelligence

## Description
Set up Python-based AI service with basic anomaly detection capabilities.

## Acceptance Criteria
- [ ] FastAPI service setup with health checks
- [ ] Basic anomaly detection model
- [ ] Risk scoring algorithm
- [ ] API endpoints for analysis
- [ ] Model training pipeline foundation
- [ ] Integration with inventory service

## Technical Requirements
- Python 3.11+ with FastAPI
- TensorFlow/PyTorch for ML models
- PostgreSQL integration
- Docker containerization
- API documentation

## Estimate
1.5 weeks" \
        --label "task,ai,high,mvp" \
        --milestone "Phase 2: Intelligence"

    # Issue 7: PCI-DSS Compliance Framework
    gh issue create \
        --title "[TASK] Implement PCI-DSS Compliance Framework" \
        --body "**Component**: Compliance Engine
**Priority**: High (Important for MVP)
**Milestone**: Phase 2: Intelligence

## Description
Build PCI-DSS compliance assessment engine with gap analysis.

## Acceptance Criteria
- [ ] PCI-DSS rule engine implementation
- [ ] Compliance assessment API
- [ ] Gap analysis reporting
- [ ] Remediation recommendations
- [ ] Score calculation and trending
- [ ] Integration with inventory data

## Technical Requirements
- Go with rule engine pattern
- JSON-based rule definitions
- PostgreSQL for compliance data
- RESTful API endpoints
- Comprehensive testing

## Estimate
1.5 weeks" \
        --label "task,compliance,high,mvp" \
        --milestone "Phase 2: Intelligence"

    log_success "Phase 2 issues created"
}

# Create Phase 3 issues
create_phase3_issues() {
    log_info "Creating Phase 3 (Enterprise) issues..."
    
    # Issue 8: Integration Service Foundation
    gh issue create \
        --title "[TASK] Build Integration Service Foundation" \
        --body "**Component**: Integration Service
**Priority**: High (Important for MVP)
**Milestone**: Phase 3: Enterprise

## Description
Create integration service architecture with plugin framework for ITAM systems.

## Acceptance Criteria
- [ ] Plugin architecture design
- [ ] Connector interface definition
- [ ] Data transformation pipeline
- [ ] Sync scheduling system
- [ ] Error handling and retry logic
- [ ] Configuration management

## Technical Requirements
- Go with plugin architecture
- Interface-based connector design
- NATS for async processing
- PostgreSQL for configuration
- Comprehensive logging

## Estimate
1 week" \
        --label "task,integration,high,mvp" \
        --milestone "Phase 3: Enterprise"

    # Issue 9: ServiceNow Connector
    gh issue create \
        --title "[TASK] Implement ServiceNow Connector" \
        --body "**Component**: Integration Service
**Priority**: High (Important for MVP)
**Milestone**: Phase 3: Enterprise

## Description
Build ServiceNow CMDB integration connector for crypto data synchronization.

## Acceptance Criteria
- [ ] ServiceNow API integration
- [ ] Field mapping configuration
- [ ] Bidirectional sync capability
- [ ] Error handling and logging
- [ ] Configuration UI
- [ ] Real-time and scheduled sync

## Technical Requirements
- ServiceNow REST API integration
- OAuth 2.0 authentication
- JSON field mapping
- Retry logic and error handling
- Rate limiting compliance

## Estimate
1.5 weeks" \
        --label "task,integration,high,mvp" \
        --milestone "Phase 3: Enterprise"

    log_success "Phase 3 issues created"
}

# Create testing and infrastructure issues
create_infrastructure_issues() {
    log_info "Creating infrastructure and testing issues..."
    
    # Issue 10: Comprehensive Testing
    gh issue create \
        --title "[TASK] Implement Comprehensive Test Suite" \
        --body "**Component**: Testing
**Priority**: High (Important for MVP)
**Milestone**: Phase 1: Foundation

## Description
Implement comprehensive testing across all components.

## Acceptance Criteria
- [ ] Unit tests for all services (>80% coverage)
- [ ] Integration tests for API endpoints
- [ ] End-to-end tests for user workflows
- [ ] Load testing for performance validation
- [ ] Security testing for vulnerability assessment
- [ ] CI/CD integration

## Technical Requirements
- Go testing framework for backend
- Jest/React Testing Library for frontend
- Testcontainers for integration tests
- k6 for load testing
- GitHub Actions integration

## Estimate
2 weeks" \
        --label "task,testing,high,mvp" \
        --milestone "Phase 1: Foundation"

    # Issue 11: Security Hardening
    gh issue create \
        --title "[TASK] Implement Security Hardening" \
        --body "**Component**: Security
**Priority**: Critical (MVP blocker)
**Milestone**: Phase 1: Foundation

## Description
Implement comprehensive security measures across the platform.

## Acceptance Criteria
- [ ] Vulnerability scanning integration
- [ ] Security headers implementation
- [ ] Input validation and sanitization
- [ ] Rate limiting on all APIs
- [ ] Comprehensive audit logging
- [ ] Secrets management

## Technical Requirements
- Trivy for vulnerability scanning
- OWASP security headers
- Input validation middleware
- Redis-based rate limiting
- Structured logging with audit trails

## Estimate
1 week" \
        --label "task,security,critical,mvp" \
        --milestone "Phase 1: Foundation"

    log_success "Infrastructure issues created"
}

# Main function
main() {
    echo "🚀 Creating GitHub Issues for Crypto Inventory Platform"
    echo "====================================================="
    echo ""
    
    check_gh_cli
    
    echo "This script will create:"
    echo "  ✓ Project milestones (4 phases)"
    echo "  ✓ Project labels (components, priorities)"
    echo "  ✓ Phase 1 issues (Foundation - 5 issues)"
    echo "  ✓ Phase 2 issues (Intelligence - 2 issues)"
    echo "  ✓ Phase 3 issues (Enterprise - 2 issues)"
    echo "  ✓ Infrastructure issues (Testing, Security - 2 issues)"
    echo ""
    
    read -p "Continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Issue creation cancelled."
        exit 0
    fi
    
    create_milestones
    create_labels
    create_phase1_issues
    create_phase2_issues
    create_phase3_issues
    create_infrastructure_issues
    
    echo ""
    log_success "All GitHub issues created successfully!"
    echo ""
    echo "🎯 Next Steps:"
    echo "  1. Review issues: gh issue list"
    echo "  2. Assign team members to issues"
    echo "  3. Start development with Phase 1 issues"
    echo "  4. Update project board: https://github.com/democorp/crypto-inventory-platform/projects"
    echo ""
    echo "📋 Issue Management:"
    echo "  - View all issues: gh issue list"
    echo "  - Filter by milestone: gh issue list --milestone 'Phase 1: Foundation'"
    echo "  - Filter by label: gh issue list --label mvp"
    echo "  - Assign issues: gh issue edit <number> --assignee @username"
}

# Run main function
main "$@"
