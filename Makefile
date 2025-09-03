.PHONY: build clean test run-account run-transfer run-fee docker-up docker-down help

# Build all services
build:
	@echo "🔨 Building all services..."
	@chmod +x scripts/build.sh
	@./scripts/build.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -cover ./...

# Run Account API
run-account:
	@echo "🏦 Starting Account API..."
	@./bin/account-api

# Run Transfer API
run-transfer:
	@echo "💸 Starting Transfer API..."
	@./bin/transfer-api

# Run Fee API
run-fee:
	@echo "💰 Starting Fee API..."
	@./bin/fee-api

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Generate swagger docs
swagger:
	@echo "📚 Generating Swagger documentation..."
	@swag init -g cmd/account-api/main.go -o docs/account
	@swag init -g cmd/transfer-api/main.go -o docs/transfer
	@swag init -g cmd/fee-api/main.go -o docs/fee

# Docker commands
docker-up:
	@echo "🐳 Starting services with Docker..."
	@docker-compose -f deployments/docker-compose.yml up --build

docker-down:
	@echo "🐳 Stopping Docker services..."
	@docker-compose -f deployments/docker-compose.yml down

docker-logs:
	@echo "📋 Showing Docker logs..."
	@docker-compose -f deployments/docker-compose.yml logs -f

# Development setup
dev-setup:
	@echo "🛠️ Setting up development environment..."
	@go mod tidy
	@chmod +x scripts/build.sh

# Health check all services
health-check:
	@echo "🏥 Checking service health..."
	@curl -f http://localhost:8001/health || echo "Account API not responding"
	@curl -f http://localhost:8002/health || echo "Transfer API not responding"
	@curl -f http://localhost:8003/health || echo "Fee API not responding"

# Help
help:
	@echo "📖 Available commands:"
	@echo "  build         - Build all services"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  run-account   - Run Account API"
	@echo "  run-transfer  - Run Transfer API"
	@echo "  run-fee       - Run Fee API"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  swagger       - Generate Swagger docs"
	@echo "  docker-up     - Start with Docker"
	@echo "  docker-down   - Stop Docker services"
	@echo "  docker-logs   - Show Docker logs"
	@echo "  dev-setup     - Setup development environment"
	@echo "  health-check  - Check service health"
	@echo "  help          - Show this help"
