#!/bin/bash

# Payment Service Environment Variables
export PAYMENT_SERVICE_PORT=8082
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="golunch_payments"

# Order Service URL for HTTP communication
export ORDER_SERVICE_URL="http://localhost:8081"

# Production Service URL for HTTP communication  
export PRODUCTION_SERVICE_URL="http://localhost:8083"

# Mercado Pago Configuration (opcional para testes)
export MERCADO_PAGO_ACCESS_TOKEN="test-token"
export MERCADO_PAGO_SELLER_APP_USER_ID="test-seller"

echo "Payment Service environment variables set:"
echo "PORT: $PAYMENT_SERVICE_PORT"
echo "MongoDB: $MONGODB_URI/$MONGODB_DATABASE" 
echo "Order Service: $ORDER_SERVICE_URL"
echo "Production Service: $PRODUCTION_SERVICE_URL"