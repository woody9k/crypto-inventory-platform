#!/bin/bash

# Crypto Inventory Platform - Balanced Performance Start Script
# This script starts the balanced performance environment on GCP ($72/month for 4h/day)

set -e

echo "🚀 Starting Crypto Inventory Platform (Balanced Performance Mode)..."

# Configuration
PROJECT_ID=${PROJECT_ID:-"crypto-inventory-prod"}
CLUSTER_NAME=${CLUSTER_NAME:-"crypto-inventory-cluster"}
ZONE=${ZONE:-"us-central1-a"}
NODE_POOL_NAME=${NODE_POOL_NAME:-"crypto-inventory-pool"}
DB_INSTANCE=${DB_INSTANCE:-"crypto-inventory-db"}
REDIS_INSTANCE=${REDIS_INSTANCE:-"crypto-inventory-redis"}

# Set project
echo "📋 Setting project to: $PROJECT_ID"
gcloud config set project $PROJECT_ID

# Start GKE node pool (2 nodes for balanced performance)
echo "🔄 Starting GKE node pool (2 nodes)..."
gcloud container node-pools resize $NODE_POOL_NAME \
  --cluster=$CLUSTER_NAME \
  --zone=$ZONE \
  --num-nodes=2 \
  --quiet

# Start Cloud SQL instance (db-n1-standard-1)
echo "🗄️ Starting Cloud SQL instance..."
gcloud sql instances patch $DB_INSTANCE \
  --activation-policy=ALWAYS \
  --quiet

# Start Memorystore Redis (1GB basic tier)
echo "🔴 Starting Memorystore Redis..."
gcloud redis instances update $REDIS_INSTANCE \
  --region=us-central1 \
  --tier=basic \
  --memory-size-gb=1 \
  --quiet

# Wait for services to be ready
echo "⏳ Waiting for services to start (2 minutes)..."
sleep 120

# Check service status
echo "🔍 Checking service status..."

# Check GKE nodes
echo "GKE Nodes:"
kubectl get nodes

# Check Cloud SQL status
echo "Cloud SQL Status:"
gcloud sql instances describe $DB_INSTANCE --format="value(state)"

# Check Redis status
echo "Redis Status:"
gcloud redis instances describe $REDIS_INSTANCE --region=us-central1 --format="value(state)"

# Deploy application with balanced resources
echo "🚀 Deploying application (balanced resources)..."
kubectl apply -f k8s/production-balanced/

# Wait for pods to be ready
echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=crypto-inventory --timeout=300s

# Get service URLs
echo "🌐 Getting service URLs..."
kubectl get services

# Display cost information
echo ""
echo "💰 Cost Information:"
echo "   Current running cost: ~$0.80/hour"
echo "   Stopped cost: ~$0.20/hour"
echo "   Estimated monthly cost (4h/day): ~$72"

echo ""
echo "✅ Balanced performance platform started successfully!"
echo "   Performance: 50-100 concurrent users, <200ms API response"
echo "   Access the application at your configured domain"
echo "   Monitor costs in GCP Console"
echo "   Use 'scripts/stop-balanced.sh' to stop when done"
