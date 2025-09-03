# Contributing to Crypto Inventory Platform

Thank you for your interest in contributing to the Crypto Inventory Platform! This document provides guidelines and information for contributors.

## 🎯 Project Overview

The Crypto Inventory Platform is an enterprise SaaS solution for discovering, inventorying, and managing cryptographic implementations across networks. Our goal is to help organizations maintain crypto-agility and compliance with security frameworks.

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher
- **Python**: Version 3.11 or higher
- **Docker**: Version 20.10 or higher
- **Git**: Latest version

### Development Environment Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/democorp/crypto-inventory-platform.git
   cd crypto-inventory-platform
   ```

2. **Start the development environment**:
   ```bash
   make start
   ```

3. **Verify services are running**:
   ```bash
   make monitor
   ```

4. **Access the application**:
   - Web UI: http://localhost:3000
   - API Gateway: http://localhost:8080
   - Database Admin: http://localhost:8090

## 📋 Development Workflow

### Branch Strategy

We use **GitFlow** branching strategy:

- **`main`**: Production-ready code
- **`develop`**: Integration branch for features
- **`feature/*`**: Feature development branches
- **`hotfix/*`**: Critical bug fixes
- **`release/*`**: Release preparation branches

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Follow coding standards (see below)
   - Write tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   make test
   make lint
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat(component): add new feature description"
   ```

5. **Push and create a pull request**:
   ```bash
   git push origin feature/your-feature-name
   # Create PR through GitHub interface
   ```

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

**Scopes**:
- `auth`: Authentication service
- `inventory`: Inventory service
- `compliance`: Compliance engine
- `reports`: Report generator
- `sensors`: Sensor manager
- `integration`: Integration service
- `ai`: AI analysis service
- `ui`: Web frontend
- `docs`: Documentation
- `infra`: Infrastructure

**Examples**:
```
feat(auth): add SSO integration with SAML support
fix(sensors): resolve memory leak in network capture
docs(api): update authentication endpoints documentation
test(compliance): add unit tests for PCI-DSS framework
```

## 🏗️ Architecture Guidelines

### Service Design Principles

1. **Single Responsibility**: Each service has a clear, focused purpose
2. **API-First**: Design APIs before implementation
3. **Stateless**: Services should be stateless when possible
4. **Fault Tolerant**: Graceful error handling and recovery
5. **Observable**: Comprehensive logging and metrics

### Code Organization

```
services/<service-name>/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and routing
│   ├── business/         # Business logic
│   ├── config/           # Configuration management
│   ├── database/         # Database access layer
│   └── models/           # Data models
├── pkg/                  # Public library code
├── tests/                # Test files
├── Dockerfile.dev        # Development Docker image
├── Dockerfile.prod       # Production Docker image
├── go.mod               # Go module definition
└── README.md            # Service-specific documentation
```

## 🧪 Testing Guidelines

### Test Categories

1. **Unit Tests**: Test individual functions/methods
   ```bash
   cd services/auth-service
   go test ./...
   ```

2. **Integration Tests**: Test service interactions
   ```bash
   make test-integration
   ```

3. **End-to-End Tests**: Test complete user workflows
   ```bash
   make test-e2e
   ```

### Test Requirements

- **Coverage**: Maintain >80% test coverage
- **Test Data**: Use fixtures and mocks, not production data
- **Isolation**: Tests should be independent and idempotent
- **Performance**: Tests should run quickly (<5 minutes total)

### Writing Tests

**Go Services**:
```go
func TestUserAuthentication(t *testing.T) {
    // Arrange
    user := &models.User{Email: "test@example.com"}
    
    // Act
    result, err := authService.Authenticate(user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result.Token)
}
```

**Frontend (Jest)**:
```javascript
describe('Dashboard Component', () => {
  it('should display KPI cards', () => {
    render(<Dashboard />);
    expect(screen.getByText('Sensors Deployed')).toBeInTheDocument();
  });
});
```

## 🎨 Coding Standards

### Go Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

**Example**:
```go
// Package auth provides authentication and authorization services.
package auth

import (
    "context"
    "fmt"
    "time"
)

// User represents a platform user.
type User struct {
    ID       string    `json:"id" db:"id"`
    Email    string    `json:"email" db:"email"`
    Role     Role      `json:"role" db:"role"`
    Created  time.Time `json:"created_at" db:"created_at"`
}

// Authenticate validates user credentials and returns a JWT token.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*Token, error) {
    if email == "" {
        return nil, fmt.Errorf("email is required")
    }
    
    // Implementation...
    return token, nil
}
```

### React/TypeScript Guidelines

- Use TypeScript for all React components
- Follow [React Best Practices](https://reactjs.org/docs/thinking-in-react.html)
- Use functional components with hooks
- Use Ant Design components consistently

**Example**:
```typescript
interface DashboardProps {
  tenantId: string;
  refreshInterval?: number;
}

export const Dashboard: React.FC<DashboardProps> = ({ 
  tenantId, 
  refreshInterval = 30000 
}) => {
  const [metrics, setMetrics] = useState<KPIMetrics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const data = await api.getKPIMetrics(tenantId);
        setMetrics(data);
      } catch (error) {
        console.error('Failed to fetch metrics:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchMetrics();
  }, [tenantId]);

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <div className="dashboard">
      <KPICards metrics={metrics} />
      <RiskHeatMap data={metrics?.riskData} />
    </div>
  );
};
```

### Python Guidelines

- Follow [PEP 8](https://www.python.org/dev/peps/pep-0008/)
- Use type hints
- Use `black` for formatting
- Use `flake8` for linting

**Example**:
```python
from typing import List, Optional
from pydantic import BaseModel


class CryptoImplementation(BaseModel):
    """Represents a cryptographic implementation."""
    
    id: str
    protocol: str
    cipher_suite: str
    risk_score: float
    confidence: float


class AnalysisService:
    """AI-powered analysis service for crypto implementations."""
    
    def __init__(self, model_path: str) -> None:
        self.model_path = model_path
        self._model: Optional[Any] = None
    
    async def analyze_implementations(
        self, 
        implementations: List[CryptoImplementation]
    ) -> List[AnalysisResult]:
        """Analyze crypto implementations for anomalies and risks."""
        if not implementations:
            return []
        
        # Implementation...
        return results
```

## 📖 Documentation Guidelines

### Code Documentation

- **Go**: Use godoc comments for exported functions
- **TypeScript**: Use JSDoc comments for complex functions
- **Python**: Use docstrings for classes and functions

### API Documentation

- Use OpenAPI 3.0 specifications
- Include request/response examples
- Document error codes and messages
- Provide integration examples

### Architecture Documentation

- Update architecture docs for significant changes
- Include diagrams for complex interactions
- Document design decisions and trade-offs

## 🚀 Performance Guidelines

### Database

- Use proper indexing strategies
- Implement connection pooling
- Use prepared statements
- Monitor query performance

### API Design

- Implement pagination for large datasets
- Use appropriate HTTP status codes
- Implement rate limiting
- Cache frequently accessed data

### Frontend

- Implement lazy loading for components
- Optimize bundle sizes
- Use React.memo for expensive components
- Implement proper error boundaries

## 🔒 Security Guidelines

### Authentication & Authorization

- Never store passwords in plain text
- Use strong JWT signing algorithms
- Implement proper session management
- Follow principle of least privilege

### Input Validation

- Validate all user inputs
- Sanitize data before database storage
- Use parameterized queries
- Implement CSRF protection

### Data Protection

- Encrypt sensitive data at rest
- Use HTTPS for all communications
- Implement proper audit logging
- Follow data retention policies

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment details** (OS, browser, versions)
2. **Steps to reproduce** the issue
3. **Expected vs actual behavior**
4. **Screenshots or logs** (if applicable)
5. **Workarounds** (if any)

Use our [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml) for consistency.

## 💡 Feature Requests

When requesting features, please include:

1. **Problem description** and business value
2. **Proposed solution** with details
3. **Alternative approaches** considered
4. **Acceptance criteria** for completion

Use our [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml) for consistency.

## 📞 Getting Help

- **Documentation**: Check [architecture docs](./architecture_docs/)
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions
- **Code Review**: Request reviews from appropriate team members

## 🏆 Recognition

Contributors will be recognized in our:

- README contributors section
- Release notes
- Project documentation
- Annual contributor highlights

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to the Crypto Inventory Platform! Your efforts help make enterprise cryptography more secure and manageable. 🔐
