# GitHub Actions Workflows

This directory contains the CI/CD workflows for the Crypto Inventory Platform.

## Workflows Overview

### 🔄 CI/CD Pipeline (`ci.yml`)

**Trigger**: Push to `main`/`develop` branches, Pull Requests to `main`

**Jobs**:

1. **Backend Tests** (`test-backend`)
   - Sets up PostgreSQL and Redis test databases
   - Tests all Go services with coverage reporting
   - Runs Go linting with golangci-lint
   - Caches Go modules for faster builds

2. **AI Service Tests** (`test-ai-service`)
   - Sets up Python 3.11 environment
   - Installs AI service dependencies
   - Runs Python tests and linting
   - Caches pip dependencies

3. **Frontend Tests** (`test-frontend`)
   - Sets up Node.js 18 environment
   - Runs frontend tests with coverage
   - Performs linting and builds
   - Caches npm dependencies

4. **Security Scanning** (`security-scan`)
   - Runs Trivy vulnerability scanner
   - Uploads results to GitHub Security tab
   - Scans both code and dependencies

5. **Build Images** (`build-images`)
   - Builds Docker images for services
   - Validates image builds without pushing
   - Runs after successful tests

6. **Integration Tests** (`integration-tests`)
   - Starts full service stack with Docker Compose
   - Runs end-to-end integration tests
   - Validates complete system functionality

## Workflow Features

### 🚀 Performance Optimizations
- **Parallel Jobs**: Tests run in parallel for faster feedback
- **Caching**: Go modules, npm packages, and pip dependencies cached
- **Matrix Builds**: Test across multiple versions (future)

### 🔒 Security Integration
- **Vulnerability Scanning**: Trivy scans for security issues
- **Secret Scanning**: GitHub native secret detection
- **Dependency Checking**: Automated dependency vulnerability checks

### 📊 Quality Gates
- **Test Coverage**: Minimum coverage requirements
- **Linting**: Code quality enforcement
- **Build Validation**: Ensure all components build successfully

### 🐛 Debugging Workflows

If workflows fail, check:

1. **Logs**: Click on failed job to view detailed logs
2. **Services**: Ensure database services start correctly
3. **Dependencies**: Check for version conflicts
4. **Environment**: Verify environment variables are set

### 📝 Workflow Status

Current status badges for README:

```markdown
![CI/CD Pipeline](https://github.com/democorp/crypto-inventory-platform/workflows/CI/CD%20Pipeline/badge.svg)
![Security Scan](https://github.com/democorp/crypto-inventory-platform/workflows/Security%20Scan/badge.svg)
```

## Future Enhancements

### 🚀 Planned Additions

1. **Release Automation**
   - Automated versioning and tagging
   - Changelog generation
   - Docker image publishing to registry

2. **Deployment Workflows**
   - Staging environment deployment
   - Production deployment with approvals
   - Rollback capabilities

3. **Performance Testing**
   - Load testing with k6
   - Performance regression detection
   - Benchmark tracking

4. **Multi-Platform Builds**
   - ARM64 support
   - Windows container builds
   - Cross-platform sensor binaries

### 🔧 Workflow Customization

To customize workflows:

1. **Environment Variables**: Add to repository secrets
2. **Job Conditions**: Modify `if` conditions for job execution
3. **Matrix Strategies**: Add multiple versions/platforms
4. **Custom Actions**: Create reusable actions for common tasks

### 📋 Workflow Maintenance

Regular maintenance tasks:

- **Update Dependencies**: Keep action versions current
- **Review Permissions**: Ensure minimal required permissions
- **Monitor Performance**: Track workflow execution times
- **Update Documentation**: Keep this README current

## Local Testing

Test workflows locally using:

```bash
# Install act for local GitHub Actions testing
# https://github.com/nektos/act

# Run specific workflow
act push

# Run specific job
act -j test-backend

# Use custom environment
act --env-file .env.test
```

## Workflow Configuration

### Service Dependencies

The workflows start these services for testing:

- **PostgreSQL 15**: Main database
- **Redis 7**: Caching and sessions
- **InfluxDB 2**: Time-series metrics (future)
- **NATS 2**: Message queue (future)

### Environment Variables

Required environment variables for workflows:

- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing secret (test value)

### Secrets Management

Repository secrets used:

- `DOCKER_REGISTRY_URL`: Container registry URL
- `DOCKER_REGISTRY_USERNAME`: Registry username
- `DOCKER_REGISTRY_PASSWORD`: Registry password
- `SECURITY_SCAN_TOKEN`: Security scanning API token

## Contributing to Workflows

When modifying workflows:

1. **Test Locally**: Use `act` to test changes
2. **Small Changes**: Make incremental improvements
3. **Documentation**: Update this README for significant changes
4. **Review**: Have workflow changes reviewed by team
5. **Monitor**: Watch first few runs of modified workflows

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   ```yaml
   # Solution: Increase health check timeout
   options: >-
     --health-cmd pg_isready
     --health-interval 10s
     --health-timeout 10s
     --health-retries 10
   ```

2. **Cache Misses**
   ```yaml
   # Solution: Update cache key pattern
   key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
   ```

3. **Timeout Issues**
   ```yaml
   # Solution: Increase timeout
   timeout-minutes: 30
   ```

4. **Permission Errors**
   ```yaml
   # Solution: Add required permissions
   permissions:
     contents: read
     security-events: write
   ```

For additional help:
- Check GitHub Actions documentation
- Review workflow run logs
- Search GitHub Community forums
- Contact the development team

---

**Keep workflows simple, fast, and reliable!** 🚀
