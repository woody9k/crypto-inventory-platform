# Changelog

All notable changes to the Crypto Inventory Management System will be documented in this file.

## [2025-09-12] - Authentication & Assets Page Fixes

### Fixed
- **Authentication Issues**: Fixed malformed password hashes in seed data that prevented user login
  - Updated all seed scripts to use proper Argon2id password hashes
  - Standardized all passwords to `Password123!` format
  - Fixed bcrypt vs Argon2id hash format mismatches

- **Assets Page Empty Issue**: Resolved JSONB field scanning errors in inventory service
  - Fixed PostgreSQL JSONB column scanning in Go
  - Added proper JSON parsing for tags and metadata fields
  - Implemented fallback handling for malformed JSON data

### Changed
- **Password Policy**: All default passwords now use strong password requirements
  - Uppercase letters, lowercase letters, numbers, and special characters
  - Consistent password format across all demo users

- **Database Schema**: Updated seed scripts to use correct column names
  - Fixed `subscription_tier` vs `subscription_tier_id` references
  - Updated tenant creation scripts for current schema

### Added
- **Test Data**: Added sample network assets for demo tenant
  - 3 sample assets: web server, API gateway, and database server
  - Proper JSON tags and metadata for testing

- **Documentation**: Enhanced deployment and troubleshooting guides
  - Added default credentials section
  - Added troubleshooting for common issues
  - Updated database schema documentation

### Technical Details
- **JSONB Handling**: Implemented proper PostgreSQL JSONB scanning in Go
  - Cast JSONB fields to text in SQL queries
  - Manual JSON parsing with error handling
  - Fallback to empty maps for malformed JSON

- **Password Hashing**: Standardized on Argon2id for all password hashes
  - Updated seed scripts with correct hash format
  - Fixed parameter mismatches (p=2 vs p=4)
  - Ensured compatibility with auth service expectations

### Files Modified
- `scripts/database/seed.sql`
- `scripts/database/001_auth_schema.sql`
- `scripts/database/seed_brian_debban.sql`
- `scripts/database/06-rbac-seed.sql`
- `scripts/database/rbac_seed.sql`
- `services/inventory-service/internal/services/asset_service.go`
- `services/inventory-service/internal/models/asset.go`
- `docs/DEPLOYMENT_GUIDE.md`
- `docs/DATABASE_SCHEMA.md`

### Breaking Changes
- None

### Migration Notes
- Existing deployments will need to update password hashes
- Restart required after code changes to pick up JSONB fixes
- All changes are backward compatible with existing data
