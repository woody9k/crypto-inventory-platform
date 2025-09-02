# Deployment Guide

## Overview

This guide covers deployment strategies for the crypto inventory platform across different environments: local development, staging, and production. The architecture is designed to be cloud-agnostic while providing specific guidance for major cloud providers.

## Deployment Architecture Principles

1. **Container-First**: All services containerized with Docker
2. **Infrastructure as Code**: Terraform for cloud resources
3. **GitOps**: Automated deployments from Git repositories
4. **Environment Parity**: Consistent environments from dev to production
5. **Security by Default**: Security configurations built into deployment
6. **Scalability**: Horizontal scaling capabilities built-in
7. **Monitoring Ready**: Observability tools deployed with applications

## Local Development Environment

### Prerequisites
```bash
# Required tools
- Docker (20.10+)
- Docker Compose (2.0+)
- Git
- Make (optional, for convenience scripts)

# Optional tools for development
- kubectl (for Kubernetes development)
- terraform (for infrastructure testing)
- go (1.19+ for backend development)
- node.js (18+ for frontend development)
```

### Quick Start
```bash
# Clone the repository
git clone https://github.com/yourorg/crypto-inventory-platform.git
cd crypto-inventory-platform

# Start all services
docker-compose up -d

# Verify services are running
docker-compose ps

# Access the application
# Web UI: http://localhost:3000
# API: http://localhost:8080
# Database: localhost:5432
# InfluxDB: http://localhost:8086
```

### Docker Compose Configuration
```yaml
# docker-compose.yml
version: '3.8'

services:
  # Database Services
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: crypto_inventory
      POSTGRES_USER: crypto_user
      POSTGRES_PASSWORD: crypto_pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U crypto_user -d crypto_inventory"]
      interval: 30s
      timeout: 10s
      retries: 3

  influxdb:
    image: influxdb:2.6-alpine
    environment:
      DOCKER_INFLUXDB_INIT_MODE: setup
      DOCKER_INFLUXDB_INIT_USERNAME: admin
      DOCKER_INFLUXDB_INIT_PASSWORD: adminpass
      DOCKER_INFLUXDB_INIT_ORG: crypto-inventory
      DOCKER_INFLUXDB_INIT_BUCKET: metrics
    volumes:
      - influx_data:/var/lib/influxdb2
    ports:
      - "8086:8086"

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"

  # Message Queue
  nats:
    image: nats:2.9-alpine
    command: ["-js", "-m", "8222"]
    ports:
      - "4222:4222"
      - "8222:8222"

  # Application Services
  auth-service:
    build: 
      context: ./services/auth-service
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://crypto_user:crypto_pass@postgres:5432/crypto_inventory
      REDIS_URL: redis://redis:6379
      JWT_SECRET: dev-secret-key
      LOG_LEVEL: debug
    ports:
      - "8081:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./services/auth-service:/app
      - /app/vendor  # Exclude vendor directory

  inventory-service:
    build:
      context: ./services/inventory-service
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://crypto_user:crypto_pass@postgres:5432/crypto_inventory
      INFLUXDB_URL: http://influxdb:8086
      INFLUXDB_TOKEN: dev-token
      NATS_URL: nats://nats:4222
      LOG_LEVEL: debug
    ports:
      - "8082:8080"
    depends_on:
      postgres:
        condition: service_healthy
      influxdb:
        condition: service_started
      nats:
        condition: service_started
    volumes:
      - ./services/inventory-service:/app
      - /app/vendor

  compliance-service:
    build:
      context: ./services/compliance-service
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://crypto_user:crypto_pass@postgres:5432/crypto_inventory
      NATS_URL: nats://nats:4222
      LOG_LEVEL: debug
    ports:
      - "8083:8080"
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_started
    volumes:
      - ./services/compliance-service:/app
      - /app/vendor

  report-service:
    build:
      context: ./services/report-service
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://crypto_user:crypto_pass@postgres:5432/crypto_inventory
      FILE_STORAGE_PATH: /app/reports
      LOG_LEVEL: debug
    ports:
      - "8084:8080"
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - ./services/report-service:/app
      - /app/vendor
      - report_storage:/app/reports

  sensor-manager:
    build:
      context: ./services/sensor-manager
      dockerfile: Dockerfile.dev
    environment:
      DATABASE_URL: postgres://crypto_user:crypto_pass@postgres:5432/crypto_inventory
      INFLUXDB_URL: http://influxdb:8086
      INFLUXDB_TOKEN: dev-token
      NATS_URL: nats://nats:4222
      LOG_LEVEL: debug
    ports:
      - "8085:8080"
    depends_on:
      postgres:
        condition: service_healthy
      influxdb:
        condition: service_started
      nats:
        condition: service_started
    volumes:
      - ./services/sensor-manager:/app
      - /app/vendor

  # API Gateway
  api-gateway:
    image: nginx:alpine
    volumes:
      - ./config/nginx/dev.conf:/etc/nginx/nginx.conf
    ports:
      - "8080:80"
    depends_on:
      - auth-service
      - inventory-service
      - compliance-service
      - report-service
      - sensor-manager

  # Web Frontend
  web-ui:
    build:
      context: ./web-ui
      dockerfile: Dockerfile.dev
    environment:
      REACT_APP_API_URL: http://localhost:8080
      REACT_APP_WS_URL: ws://localhost:8080
    ports:
      - "3000:3000"
    volumes:
      - ./web-ui:/app
      - /app/node_modules
    depends_on:
      - api-gateway

  # Development Tools
  adminer:
    image: adminer:latest
    environment:
      ADMINER_DEFAULT_SERVER: postgres
    ports:
      - "8090:8080"
    depends_on:
      - postgres

volumes:
  postgres_data:
  influx_data:
  redis_data:
  report_storage:
```

### Development Scripts
```bash
# Makefile for convenience commands
.PHONY: start stop restart logs test migrate

start:
	docker-compose up -d

stop:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

test:
	docker-compose exec auth-service go test ./...
	docker-compose exec inventory-service go test ./...
	docker-compose exec compliance-service go test ./...

migrate:
	docker-compose exec postgres psql -U crypto_user -d crypto_inventory -f /docker-entrypoint-initdb.d/migrations.sql

clean:
	docker-compose down -v
	docker system prune -f
```

## Staging Environment

### Infrastructure Setup (AWS Example)
```hcl
# terraform/staging/main.tf
terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
  
  backend "s3" {
    bucket = "crypto-inventory-terraform-state"
    key    = "staging/terraform.tfstate"
    region = "us-west-2"
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC and Networking
module "vpc" {
  source = "../modules/vpc"
  
  environment = "staging"
  cidr_block  = "10.1.0.0/16"
  
  availability_zones = ["us-west-2a", "us-west-2b"]
  private_subnets    = ["10.1.1.0/24", "10.1.2.0/24"]
  public_subnets     = ["10.1.101.0/24", "10.1.102.0/24"]
  
  enable_nat_gateway = true
  enable_vpn_gateway = false
}

# EKS Cluster
module "eks" {
  source = "../modules/eks"
  
  cluster_name    = "crypto-inventory-staging"
  cluster_version = "1.27"
  
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnets
  
  node_groups = {
    main = {
      instance_types = ["t3.medium"]
      min_size      = 2
      max_size      = 5
      desired_size  = 3
    }
  }
}

# RDS PostgreSQL
module "database" {
  source = "../modules/rds"
  
  environment = "staging"
  
  engine_version    = "15.3"
  instance_class    = "db.t3.micro"
  allocated_storage = 20
  
  db_name  = "cryptoinventory"
  username = "cryptouser"
  
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnets
  allowed_cidr_blocks = [module.vpc.vpc_cidr_block]
}

# ElastiCache Redis
module "redis" {
  source = "../modules/elasticache"
  
  environment = "staging"
  
  node_type           = "cache.t3.micro"
  num_cache_nodes     = 1
  parameter_group     = "default.redis7"
  
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnets
  allowed_cidr_blocks = [module.vpc.vpc_cidr_block]
}

# Application Load Balancer
module "alb" {
  source = "../modules/alb"
  
  environment = "staging"
  
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.public_subnets
  
  certificate_arn = aws_acm_certificate.main.arn
}

# S3 Bucket for Reports
module "s3" {
  source = "../modules/s3"
  
  bucket_name = "crypto-inventory-reports-staging"
  environment = "staging"
  
  enable_versioning = true
  enable_encryption = true
}
```

### Kubernetes Deployment
```yaml
# k8s/staging/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: crypto-inventory-staging
  labels:
    environment: staging
---
# k8s/staging/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: database-credentials
  namespace: crypto-inventory-staging
type: Opaque
data:
  username: <base64-encoded-username>
  password: <base64-encoded-password>
  host: <base64-encoded-host>
---
# k8s/staging/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: crypto-inventory-staging
data:
  LOG_LEVEL: "info"
  ENVIRONMENT: "staging"
  DATABASE_MAX_CONNECTIONS: "20"
---
# k8s/staging/auth-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: crypto-inventory-staging
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: crypto-inventory/auth-service:staging-latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-credentials
              key: url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret
        envFrom:
        - configMapRef:
            name: app-config
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: crypto-inventory-staging
spec:
  selector:
    app: auth-service
  ports:
  - port: 80
    targetPort: 8080
```

### CI/CD Pipeline (GitHub Actions)
```yaml
# .github/workflows/deploy-staging.yml
name: Deploy to Staging

on:
  push:
    branches:
      - develop
  workflow_dispatch:

env:
  AWS_REGION: us-west-2
  EKS_CLUSTER_NAME: crypto-inventory-staging

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.19
    
    - name: Run tests
      run: |
        go test ./services/auth-service/...
        go test ./services/inventory-service/...
        go test ./services/compliance-service/...
    
    - name: Run frontend tests
      working-directory: ./web-ui
      run: |
        npm ci
        npm test -- --coverage --watchAll=false

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v2
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: ${{ env.AWS_REGION }}
    
    - name: Login to Amazon ECR
      id: login-ecr
      uses: aws-actions/amazon-ecr-login@v1
    
    - name: Build and push images
      env:
        ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
        IMAGE_TAG: staging-${{ github.sha }}
      run: |
        # Build and push each service
        services=("auth-service" "inventory-service" "compliance-service" "report-service" "sensor-manager")
        for service in "${services[@]}"; do
          docker build -t $ECR_REGISTRY/$service:$IMAGE_TAG ./services/$service
          docker push $ECR_REGISTRY/$service:$IMAGE_TAG
        done
        
        # Build and push web UI
        docker build -t $ECR_REGISTRY/web-ui:$IMAGE_TAG ./web-ui
        docker push $ECR_REGISTRY/web-ui:$IMAGE_TAG

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v2
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: ${{ env.AWS_REGION }}
    
    - name: Update kubeconfig
      run: |
        aws eks update-kubeconfig --region ${{ env.AWS_REGION }} --name ${{ env.EKS_CLUSTER_NAME }}
    
    - name: Deploy to Kubernetes
      env:
        IMAGE_TAG: staging-${{ github.sha }}
      run: |
        # Update image tags in Kubernetes manifests
        sed -i "s|:staging-latest|:$IMAGE_TAG|g" k8s/staging/*.yaml
        
        # Apply Kubernetes manifests
        kubectl apply -f k8s/staging/
        
        # Wait for deployment to complete
        kubectl rollout status deployment/auth-service -n crypto-inventory-staging
        kubectl rollout status deployment/inventory-service -n crypto-inventory-staging
    
    - name: Run integration tests
      run: |
        # Wait for services to be ready
        kubectl wait --for=condition=ready pod -l app=auth-service -n crypto-inventory-staging --timeout=300s
        
        # Run integration tests against staging environment
        go test ./tests/integration/... -tags=integration
```

## Production Environment

### High Availability Setup
```hcl
# terraform/production/main.tf
module "vpc" {
  source = "../modules/vpc"
  
  environment = "production"
  cidr_block  = "10.0.0.0/16"
  
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets    = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  
  enable_nat_gateway = true
  single_nat_gateway = false  # Multi-AZ NAT for HA
}

module "eks" {
  source = "../modules/eks"
  
  cluster_name    = "crypto-inventory-production"
  cluster_version = "1.27"
  
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets
  
  node_groups = {
    main = {
      instance_types = ["m5.large"]
      min_size      = 3
      max_size      = 20
      desired_size  = 6
    }
    spot = {
      instance_types = ["m5.large", "m5a.large", "m4.large"]
      capacity_type  = "SPOT"
      min_size      = 0
      max_size      = 10
      desired_size  = 3
    }
  }
}

module "database" {
  source = "../modules/rds"
  
  environment = "production"
  
  engine_version    = "15.3"
  instance_class    = "db.r5.large"
  allocated_storage = 100
  
  multi_az               = true
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  monitoring_interval = 60
  performance_insights_enabled = true
  
  replica_count = 1  # Read replica for reporting
}
```

### Multi-Tenant Isolation
```yaml
# k8s/production/tenant-isolation.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: tenant-enterprise-customer
  labels:
    tenant: enterprise-customer
    isolation: dedicated
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: tenant-isolation
  namespace: tenant-enterprise-customer
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: crypto-inventory-shared
    - podSelector: {}
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: crypto-inventory-shared
    - podSelector: {}
  - to: []  # Allow external traffic
    ports:
    - protocol: TCP
      port: 443
    - protocol: TCP
      port: 80
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: tenant-quota
  namespace: tenant-enterprise-customer
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.cpu: "8"
    limits.memory: 16Gi
    pods: "50"
    services: "10"
```

### Monitoring and Observability
```yaml
# k8s/monitoring/prometheus.yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
  namespace: monitoring
spec:
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      team: crypto-inventory
  ruleSelector:
    matchLabels:
      team: crypto-inventory
  resources:
    requests:
      memory: 400Mi
  retention: 30d
  storage:
    volumeClaimTemplate:
      spec:
        storageClassName: gp2
        resources:
          requests:
            storage: 100Gi
---
# k8s/monitoring/grafana.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - name: grafana
        image: grafana/grafana:9.5.0
        env:
        - name: GF_SECURITY_ADMIN_PASSWORD
          valueFrom:
            secretKeyRef:
              name: grafana-admin
              key: password
        - name: GF_AUTH_ANONYMOUS_ENABLED
          value: "false"
        - name: GF_AUTH_OAUTH_AUTO_LOGIN
          value: "true"
        volumeMounts:
        - name: grafana-storage
          mountPath: /var/lib/grafana
        - name: grafana-config
          mountPath: /etc/grafana/provisioning
      volumes:
      - name: grafana-storage
        persistentVolumeClaim:
          claimName: grafana-pvc
      - name: grafana-config
        configMap:
          name: grafana-config
```

## Security Considerations

### SSL/TLS Configuration
```yaml
# k8s/security/tls-config.yaml
apiVersion: v1
kind: Secret
metadata:
  name: tls-certificate
  namespace: crypto-inventory-production
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-certificate>
  tls.key: <base64-encoded-private-key>
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: crypto-inventory-ingress
  namespace: crypto-inventory-production
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - api.cryptoinventory.com
    secretName: api-tls-secret
  - hosts:
    - app.cryptoinventory.com
    secretName: app-tls-secret
  rules:
  - host: api.cryptoinventory.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 80
  - host: app.cryptoinventory.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web-ui
            port:
              number: 80
```

### Pod Security Standards
```yaml
# k8s/security/pod-security.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: crypto-inventory-production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
---
apiVersion: v1
kind: SecurityContext
metadata:
  name: restricted-security-context
spec:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
  seccompProfile:
    type: RuntimeDefault
  capabilities:
    drop:
    - ALL
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
```

## Backup and Disaster Recovery

### Database Backup Strategy
```bash
#!/bin/bash
# scripts/backup-database.sh

# RDS Automated Backups
aws rds create-db-snapshot \
  --db-instance-identifier crypto-inventory-prod \
  --db-snapshot-identifier crypto-inventory-$(date +%Y%m%d-%H%M%S)

# Cross-region backup
aws rds copy-db-snapshot \
  --source-db-snapshot-identifier crypto-inventory-$(date +%Y%m%d-%H%M%S) \
  --target-db-snapshot-identifier crypto-inventory-$(date +%Y%m%d-%H%M%S)-dr \
  --source-region us-west-2 \
  --target-region us-east-1
```

### Application Data Backup
```yaml
# k8s/backup/velero-backup.yaml
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: crypto-inventory-backup
  namespace: velero
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  template:
    includedNamespaces:
    - crypto-inventory-production
    - tenant-*
    excludedResources:
    - events
    - logs
    ttl: 720h  # 30 days
    storageLocation: default
    volumeSnapshotLocations:
    - default
```

## Scaling Strategies

### Horizontal Pod Autoscaler
```yaml
# k8s/scaling/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: inventory-service-hpa
  namespace: crypto-inventory-production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: inventory-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: custom_queue_depth
      target:
        type: AverageValue
        averageValue: "30"
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

### Cluster Autoscaler
```yaml
# k8s/scaling/cluster-autoscaler.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cluster-autoscaler
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: cluster-autoscaler
  template:
    metadata:
      labels:
        app: cluster-autoscaler
    spec:
      containers:
      - image: k8s.gcr.io/autoscaling/cluster-autoscaler:v1.27.0
        name: cluster-autoscaler
        command:
        - ./cluster-autoscaler
        - --v=4
        - --stderrthreshold=info
        - --cloud-provider=aws
        - --skip-nodes-with-local-storage=false
        - --expander=least-waste
        - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/crypto-inventory-production
        - --balance-similar-node-groups
        - --skip-nodes-with-system-pods=false
```

---

*This deployment guide provides comprehensive instructions for deploying the crypto inventory platform across all environments while maintaining security, scalability, and reliability.*
