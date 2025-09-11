#!/bin/bash

# Crypto Inventory Platform - Production Stop Script
# This script stops the production environment to minimize costs

set -e

echo "🛑 Stopping Crypto Inventory Platform (Production Mode)..."

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

# Scale down GKE node pool to 0
echo "🔄 Scaling down GKE node pool..."
gcloud container node-pools resize $NODE_POOL_NAME \
  --cluster=$CLUSTER_NAME \
  --zone=$ZONE \
  --num-nodes=0 \
  --quiet

# Stop Cloud SQL instance
echo "🗄️ Stopping Cloud SQL instance..."
gcloud sql instances patch $DB_INSTANCE \
  --activation-policy=NEVER \
  --quiet

# Scale down Redis to basic tier
echo "🔴 Scaling down Memorystore Redis..."
gcloud redis instances update $REDIS_INSTANCE \
  --region=us-central1 \
  --tier=basic \
  --memory-size-gb=1 \
  --quiet

# Wait for services to stop
echo "⏳ Waiting for services to stop..."
sleep 60

# Check service status
echo "🔍 Checking service status..."

# Check GKE nodes
echo "GKE Nodes:"
kubectl get nodes || echo "No nodes available (expected when stopped)"

# Check Cloud SQL status
echo "Cloud SQL Status:"
gcloud sql instances describe $DB_INSTANCE --format="value(state)"

# Check Redis status
echo "Redis Status:"
gcloud redis instances describe $REDIS_INSTANCE --region=us-central1 --format="value(state)"

# Display cost information
echo ""
echo "💰 Cost Information:"
echo "   Current stopped cost: ~$0.25/hour"
echo "   Running cost: ~$1.00/hour"
echo "   Estimated monthly cost (4h/day): ~$270"

echo ""
echo "✅ Platform stopped successfully!"
echo "   All compute resources have been scaled down"
echo "   Databases are stopped to minimize costs"
echo "   Use 'scripts/start-production.sh' to restart when needed"
echo ""
echo "💡 Cost Savings:"
echo "   - GKE nodes: $0.40/hour saved"
echo "   - Cloud SQL: $0.17/hour saved"
echo "   - Redis: $0.05/hour saved"
echo "   - Total savings: $0.62/hour when stopped"
