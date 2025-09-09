#!/bin/bash

# =================================================================
# RBAC Migration Script
# =================================================================
# This script runs the RBAC migration to add role-based access control
# to the crypto inventory platform

set -e

echo "🚀 Starting RBAC Migration..."

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running"
    exit 1
fi

# Check if PostgreSQL container is running
if ! docker-compose ps postgres | grep -q "Up"; then
    echo "❌ Error: PostgreSQL container is not running"
    echo "Please start the database first: docker-compose up -d postgres"
    exit 1
fi

echo "📊 Running RBAC migration..."

# Run the RBAC migration
docker-compose exec -T postgres psql -U crypto_user -d crypto_inventory -f /docker-entrypoint-initdb.d/rbac_migration.sql

if [ $? -eq 0 ]; then
    echo "✅ RBAC migration completed successfully"
else
    echo "❌ RBAC migration failed"
    exit 1
fi

echo "🌱 Seeding RBAC data..."

# Run the RBAC seed data
docker-compose exec -T postgres psql -U crypto_user -d crypto_inventory -f /docker-entrypoint-initdb.d/rbac_seed.sql

if [ $? -eq 0 ]; then
    echo "✅ RBAC seed data completed successfully"
else
    echo "❌ RBAC seed data failed"
    exit 1
fi

echo "🔄 Restarting auth service to load RBAC changes..."

# Restart the auth service to load the new RBAC code
docker-compose restart auth-service

if [ $? -eq 0 ]; then
    echo "✅ Auth service restarted successfully"
else
    echo "❌ Failed to restart auth service"
    exit 1
fi

echo ""
echo "🎉 RBAC Migration Complete!"
echo ""
echo "What was added:"
echo "  ✅ Platform-level roles (super_admin, platform_admin, support_admin)"
echo "  ✅ Tenant-level roles (tenant_owner, tenant_admin, security_admin, analyst, viewer, api_user)"
echo "  ✅ Granular permissions system"
echo "  ✅ User role assignments"
echo "  ✅ Permission checking functions"
echo "  ✅ Audit logging"
echo "  ✅ Enhanced middleware for permission enforcement"
echo "  ✅ API endpoints for role and permission management"
echo ""
echo "Next steps:"
echo "  1. Test the new RBAC system"
echo "  2. Build frontend components for role management"
echo "  3. Update existing services to use RBAC permissions"
echo ""
echo "🔗 API Endpoints available:"
echo "  GET    /api/v1/permissions                    # List all permissions"
echo "  GET    /api/v1/tenant/:tenantId/roles         # List tenant roles"
echo "  GET    /api/v1/tenant/:tenantId/users/:userId/roles  # List user roles"
echo "  POST   /api/v1/tenant/:tenantId/users/:userId/roles  # Assign role to user"
echo "  POST   /api/v1/permissions/check              # Check user permission"
echo "  GET    /api/v1/platform/users                 # List platform users (admin only)"
echo "  GET    /api/v1/audit/logs                     # View audit logs"
echo ""
