#!/bin/bash

set -e

echo "🔨 Building BankMore Backend in Go..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build Account API
echo "📦 Building Account API..."
CGO_ENABLED=1 go build -o bin/account-api ./cmd/account-api

# Build Transfer API
echo "📦 Building Transfer API..."
CGO_ENABLED=1 go build -o bin/transfer-api ./cmd/transfer-api

# Build Fee API
echo "📦 Building Fee API..."
CGO_ENABLED=1 go build -o bin/fee-api ./cmd/fee-api

echo "✅ Build completed successfully!"
echo ""
echo "📋 Available binaries:"
echo "  - bin/account-api  (Account API - Port 8001)"
echo "  - bin/transfer-api (Transfer API - Port 8002)"
echo "  - bin/fee-api      (Fee API - Port 8003)"
echo ""
echo "🚀 To run the services:"
echo "  ./bin/account-api"
echo "  ./bin/transfer-api"
echo "  ./bin/fee-api"
echo ""
echo "🐳 To run with Docker:"
echo "  docker-compose -f deployments/docker-compose.yml up --build"
