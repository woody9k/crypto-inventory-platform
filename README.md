# Crypto Inventory SaaS Platform

A comprehensive SaaS platform for discovering, inventorying, and managing cryptographic implementations across enterprise networks.

## 🎯 Overview

This platform helps organizations:
- **Discover** all cryptographic implementations across their networks
- **Analyze compliance** against frameworks like PCI-DSS, NIST, FIPS
- **Generate reports** for executive and compliance teams
- **Monitor in real-time** for changes in crypto configurations
- **Prepare for post-quantum** cryptography migration

## 🏗️ Architecture

### Core Services
- **Auth Service**: Multi-tenant authentication and authorization
- **Inventory Service**: Asset discovery and crypto implementation tracking
- **Compliance Engine**: Framework-specific compliance analysis
- **Report Generator**: PDF/Excel report generation and scheduling
- **Sensor Manager**: Network sensor coordination and data ingestion
- **Integration Service**: ITAM system connectors and data synchronization
- **AI Analysis Service**: Machine learning for anomaly detection and insights

### Network Sensor
- **Cross-Platform**: Windows, Linux, macOS, ARM support
- **Flexible Deployment**: Native binaries, containers, Windows Services, systemd
- **Passive Analysis**: TLS, SSH, IPSec, VPN discovery
- **Edge AI**: Local anomaly detection capabilities

## 🚀 Quick Start

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.19+ (for development)
- Node.js 18+ (for frontend development)

### Local Development
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Access the application
# Web UI: http://localhost:3000
# API Gateway: http://localhost:8080
# Database Admin: http://localhost:8090
```

### Building Services
```bash
# Build all Go services
make build-services

# Build network sensor
make build-sensor

# Build frontend
make build-frontend

# Run tests
make test
```

## 📁 Project Structure

```
crypto-inventory-platform/
├── services/                   # Backend microservices
│   ├── auth-service/          # Authentication & authorization
│   ├── inventory-service/     # Asset & crypto inventory
│   ├── compliance-engine/     # Framework compliance analysis
│   ├── report-generator/      # PDF/Excel report generation
│   ├── sensor-manager/        # Sensor coordination
│   ├── integration-service/   # ITAM system connectors
│   └── ai-analysis-service/   # AI/ML analysis and inference
├── sensor/                    # Cross-platform network sensor
├── web-ui/                    # React frontend application
├── infrastructure/            # Infrastructure as Code
│   ├── terraform/            # Cloud infrastructure
│   └── helm/                 # Kubernetes deployments
├── k8s/                      # Kubernetes manifests
├── tests/                    # Test suites
├── scripts/                  # Deployment and utility scripts
└── docs/                     # Additional documentation
```

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend Services** | Go + Gin | High performance, cross-platform |
| **AI/ML Service** | Python + TensorFlow | Machine learning and analysis |
| **Network Sensor** | Go (multi-platform) | Cross-platform agent deployment |
| **Frontend** | React + TypeScript | Enterprise web interface |
| **Primary Database** | PostgreSQL | Relational data storage |
| **Time-Series DB** | InfluxDB | Metrics and time-stamped data |
| **Cache** | Redis | Session storage and caching |
| **Message Queue** | NATS | Inter-service communication |
| **Orchestration** | Kubernetes | Container orchestration |
| **Infrastructure** | Terraform | Infrastructure as Code |

## 🔒 Security

- **Zero Trust Architecture**: Never trust, always verify
- **Encryption Everywhere**: TLS 1.3, encryption at rest
- **Multi-tenant Isolation**: Secure tenant separation
- **RBAC**: Granular role-based access control
- **Compliance Ready**: SOC 2, GDPR, NIST framework support

## 📊 Monitoring

- **Metrics**: Prometheus + Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Tracing**: Jaeger for distributed tracing
- **Health Checks**: Built into all services

## 🚢 Deployment

### Development
```bash
docker-compose up -d
```

### Staging/Production
```bash
# Deploy with Terraform
cd infrastructure/terraform
terraform apply

# Deploy applications with Helm
helm install crypto-inventory ./infrastructure/helm
```

## 🧪 Testing

```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# End-to-end tests
make test-e2e

# Load tests
make test-load
```

## 📖 Documentation

- [Architecture Overview](./architecture_docs/02_system_architecture.md)
- [API Documentation](./architecture_docs/05_api_specifications.md)
- [Deployment Guide](./architecture_docs/06_deployment_guide.md)
- [Development Setup](./docs/development.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is proprietary software. All rights reserved.

## 🆘 Support

For support and questions:
- Check the [documentation](./docs/)
- Review [architecture decisions](./architecture_docs/)
- Open an issue for bugs or feature requests

---

*Building the future of cryptographic visibility and compliance.*
