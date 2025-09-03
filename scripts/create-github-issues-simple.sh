#!/bin/bash

# Simplified GitHub Issues Creation Script
# This script creates essential labels and development issues

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

# Create essential labels only
create_labels() {
    log_info "Creating essential project labels..."
    
    essential_labels=(
        "task|Development Task|0e8a16"
        "mvp|Critical for MVP|e11d21"
        "authentication|Authentication Service|0052cc"
        "sensors|Sensor Manager|c2e0c6"
        "frontend|Web Frontend|fef2c0"
        "database|Database Related|bfd4f2"
        "critical|Critical Priority|b60205"
        "high|High Priority|d93f0b"
    )
    
    for label_data in "${essential_labels[@]}"; do
        IFS='|' read -r name description color <<< "$label_data"
        
        log_info "Creating label: $name"
        gh api repos/:owner/:repo/labels \
            --method POST \
            --field name="$name" \
            --field description="$description" \
            --field color="$color" 2>/dev/null || log_warning "Label $name might already exist"
    done
    
    log_success "Essential labels processed"
}

# Create core issues only
create_core_issues() {
    log_info "Creating core development issues..."
    
    # Issue 1: Authentication Service MVP
    log_info "Creating Authentication Service MVP issue..."
    gh issue create \
        --title "[TASK] Implement Authentication Service MVP" \
        --body "**Component**: Authentication Service
**Priority**: Critical (MVP blocker)

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
        --label "task,authentication,critical,mvp" || log_error "Failed to create authentication issue"

    # Issue 2: Cross-Platform Network Sensor
    log_info "Creating Network Sensor issue..."
    gh issue create \
        --title "[TASK] Build Cross-Platform Network Sensor Foundation" \
        --body "**Component**: Network Sensor
**Priority**: Critical (MVP blocker)

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
        --label "task,sensors,critical,mvp" || log_error "Failed to create sensor issue"

    # Issue 3: Database Schema
    log_info "Creating Database Schema issue..."
    gh issue create \
        --title "[TASK] Create Database Schema and Migrations" \
        --body "**Component**: Database
**Priority**: Critical (MVP blocker)

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
        --label "task,database,critical,mvp" || log_error "Failed to create database issue"

    # Issue 4: Executive Dashboard UI
    log_info "Creating Dashboard UI issue..."
    gh issue create \
        --title "[TASK] Build Executive Dashboard Interface" \
        --body "**Component**: Web Frontend
**Priority**: High (Important for MVP)

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
        --label "task,frontend,high,mvp" || log_error "Failed to create dashboard issue"

    log_success "Core issues created"
}

# Main function
main() {
    echo "🚀 Creating Essential GitHub Issues for Crypto Inventory Platform"
    echo "=============================================================="
    echo ""
    
    check_gh_cli
    
    echo "This simplified script will create:"
    echo "  ✓ Essential project labels"
    echo "  ✓ 4 core development issues (Foundation)"
    echo ""
    
    read -p "Continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Issue creation cancelled."
        exit 0
    fi
    
    create_labels
    create_core_issues
    
    echo ""
    log_success "Essential GitHub issues created successfully!"
    echo ""
    echo "🎯 Next Steps:"
    echo "  1. Review issues: gh issue list"
    echo "  2. Assign team members: gh issue edit <number> --assignee @username"
    echo "  3. Start development with core issues"
    echo "  4. Add more issues later as needed"
    echo ""
    echo "📋 Issue Management:"
    echo "  - View all issues: gh issue list"
    echo "  - Filter by label: gh issue list --label mvp"
    echo "  - Create additional issues using GitHub web interface"
}

# Run main function
main "$@"
