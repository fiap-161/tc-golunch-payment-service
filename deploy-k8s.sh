#!/bin/bash

# Deploy Payment Service to Kubernetes
# Usage: ./deploy-k8s.sh [namespace]

NAMESPACE=${1:-golunch}

echo "💳 Deploying Payment Service to namespace: ${NAMESPACE}"
echo "💰 Cost: $0 (using MongoDB StatefulSet)"

# Create namespace if it doesn't exist
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -

echo "🗄️ Deploying MongoDB..."
kubectl apply -f k8s/mongodb-statefulset.yaml -n ${NAMESPACE}

echo "⏳ Waiting for MongoDB to be ready..."
kubectl wait --for=condition=ready pod -l app=mongodb-payment -n ${NAMESPACE} --timeout=300s

echo "📦 Applying ConfigMap..."
kubectl apply -f k8s/payment-service-configmap.yaml -n ${NAMESPACE}

echo "🔐 Applying Secrets..."
kubectl apply -f k8s/payment-service-secrets.yaml -n ${NAMESPACE}

echo "🚀 Applying Deployment..."
kubectl apply -f k8s/payment-service-deployment.yaml -n ${NAMESPACE}

echo "🌐 Applying Service..."
kubectl apply -f k8s/payment-service-service.yaml -n ${NAMESPACE}

echo "📈 Applying HPA..."
kubectl apply -f k8s/payment-service-hpa.yaml -n ${NAMESPACE}

# Wait for deployment to be ready
echo "⏳ Waiting for Payment Service to be ready..."
kubectl rollout status deployment/payment-service -n ${NAMESPACE} --timeout=300s

# Show deployment status
echo ""
echo "✅ Payment Service Deployment Status:"
kubectl get pods -l app=payment-service -n ${NAMESPACE}
kubectl get pods -l app=mongodb-payment -n ${NAMESPACE}
kubectl get svc -n ${NAMESPACE} | grep payment

echo ""
echo "🎉 Payment Service deployed successfully!"
echo ""
echo "📊 Next Steps:"
echo "  • Test: kubectl port-forward svc/payment-service 8082:8082 -n ${NAMESPACE}"
echo "  • Check: curl http://localhost:8082/ping"
echo "  • Logs: kubectl logs -f deployment/payment-service -n ${NAMESPACE}"
echo "  • DB Access: kubectl port-forward svc/mongodb-payment 27017:27017 -n ${NAMESPACE}"