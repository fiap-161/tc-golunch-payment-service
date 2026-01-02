# Payment Service Makefile

.PHONY: build run test test-unit test-integration test-bdd clean coverage lint ci db-setup help

# Variáveis
BINARY_NAME=payment-service
GO_MODULE=github.com/fiap-161/tc-golunch-payment-service

# Build
build:
	@echo "🔨 Building Payment Service..."
	go build -o bin/$(BINARY_NAME) cmd/api/main.go

# Executar aplicação
run:
	@echo "🚀 Starting Payment Service on port 8082..."
	go run cmd/api/main.go

# Testes
test: test-unit test-integration

test-unit:
	@echo "🧪 Running Unit Tests..."
	go test -v ./internal/... -coverprofile=coverage-unit.out
	go tool cover -html=coverage-unit.out -o coverage-unit.html

test-integration:
	@echo "🔗 Running Integration Tests (with mocked dependencies)..."
	go test -v ./tests/... -coverprofile=coverage-integration.out
	go tool cover -html=coverage-integration.out -o coverage-integration.html

# BDD Tests
test-bdd:
	@echo "🥒 Running BDD Tests..."
	@if command -v ginkgo > /dev/null; then \
		ginkgo -r --cover --coverprofile=coverage-bdd.out; \
	else \
		echo "⚠️  Ginkgo not installed. Running standard BDD-style tests..."; \
		go test -v ./tests/... -tags=bdd; \
	fi

# Coverage total (80%+ obrigatório)
coverage:
	@echo "📊 Generating Total Coverage Report..."
	go test -v ./... -coverprofile=coverage-total.out
	go tool cover -html=coverage-total.out -o coverage-total.html
	@echo "📈 Coverage Summary:"
	go tool cover -func=coverage-total.out | grep total
	@echo "🎯 Target: 80% minimum coverage"

# Linting
lint:
	@echo "🔍 Running Linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Verificar dependências
mod-tidy:
	@echo "📦 Tidying modules..."
	go mod tidy

# Verificar vulnerabilidades
security-check:
	@echo "🔒 Running Security Check..."
	@if command -v govulncheck > /dev/null; then \
		govulncheck ./...; \
	else \
		echo "⚠️  govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

# Pipeline de CI/CD completa
ci: mod-tidy lint test coverage security-check
	@echo "✅ Payment Service CI Pipeline completed successfully!"
	@echo "📊 Verifying 80% coverage requirement..."
	@go tool cover -func=coverage-total.out | grep total | awk '{if ($$3+0 >= 80.0) print "✅ Coverage OK:", $$3; else print "❌ Coverage LOW:", $$3, "- Need 80%+"}'

# Limpar arquivos gerados
clean:
	@echo "🧹 Cleaning up..."
	rm -f bin/$(BINARY_NAME)
	rm -f coverage-*.out coverage-*.html
	go clean -testcache

# Docker
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t payment-service:latest .

docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8082:8082 --name payment-service payment-service:latest

# Database setup (MongoDB)
db-setup:
	@echo "🗄️ Setting up Payment Service database (MongoDB)..."
	docker run -d \
		--name golunch_payments_db \
		-e MONGO_INITDB_ROOT_USERNAME=golunch_user \
		-e MONGO_INITDB_ROOT_PASSWORD=golunch_password \
		-e MONGO_INITDB_DATABASE=golunch_payments \
		-p 27017:27017 \
		mongo:4.4

db-stop:
	@echo "🛑 Stopping Payment Service database..."
	docker stop golunch_payments_db || true
	docker rm golunch_payments_db || true

# Test com dependências mockadas
test-mock-deps:
	@echo "🎭 Running tests with mocked external dependencies..."
	@echo "   - Core Service: Mocked"
	@echo "   - Operation Service: Mocked"  
	@echo "   - MercadoPago API: Mocked"
	go test -v ./tests/... -tags=mock

# Verificar saúde do serviço
health-check:
	@echo "🏥 Checking Payment Service health..."
	@if curl -s http://localhost:8082/ping > /dev/null; then \
		echo "✅ Payment Service is healthy"; \
	else \
		echo "❌ Payment Service is not responding"; \
	fi

# Mostrar ajuda
help:
	@echo "📚 Payment Service - Available commands:"
	@echo ""
	@echo "🚀 Development:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application (port 8082)"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests (mocked deps)"
	@echo "  test-bdd           - Run BDD tests"
	@echo "  test-mock-deps     - Run tests with all dependencies mocked"
	@echo "  coverage           - Generate coverage report (80%+ required)"
	@echo ""
	@echo "🔍 Quality:"
	@echo "  lint               - Run linter"
	@echo "  security-check     - Run security vulnerability check"
	@echo "  ci                 - Run full CI pipeline"
	@echo ""
	@echo "🗄️ Database:"
	@echo "  db-setup           - Setup MongoDB database"
	@echo "  db-stop            - Stop database"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo ""
	@echo "🏥 Monitoring:"
	@echo "  health-check       - Check service health"
	@echo "  help               - Show this help"
	@echo ""
	@echo "📋 Note: This service communicates with Core and Operation services via HTTP"

# Default target
.DEFAULT_GOAL := help