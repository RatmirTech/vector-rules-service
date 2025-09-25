.PHONY: proto build run test docker-build docker-up docker-down clean

# Proto generation
proto:
	@echo "Generating protobuf code..."
	@mkdir -p internal/transport/grpc/pb
	@protoc \
		--proto_path=proto \
		--go_out=internal/transport/grpc/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/transport/grpc/pb \
		--go-grpc_opt=paths=source_relative \
		proto/*.proto

# Build application
build:
	@echo "Building application..."
	@go build -o bin/vector-rules-service cmd/server/main.go

# Run application
run:
	@echo "Running application..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Docker operations
docker-build:
	@echo "Building Docker image..."
	@docker build -t vector-rules-service:latest .

docker-up:
	@echo "Starting containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping containers..."
	@docker-compose down

db-up:
	@echo "Starting database..."
	@docker-compose up -d postgres

migrate:
	@echo "Running migrations..."
	@docker-compose exec postgres psql -U postgres -d vector_rules -f /docker-entrypoint-initdb.d/001_init.sql || true

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf internal/transport/grpc/pb/

# Setup development environment
setup: proto
	@echo "Setting up development environment..."
	@go mod tidy
	@go mod download

# Help
help:
	@echo "Available commands:"
	@echo "  proto       - Generate protobuf code"
	@echo "  build       - Build application"
	@echo "  run         - Run application"
	@echo "  test        - Run tests"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-up   - Start containers"
	@echo "  docker-down - Stop containers"
	@echo "  db-up       - Start database only"
	@echo "  migrate     - Run database migrations"
	@echo "  clean       - Clean build artifacts"
	@echo "  setup       - Setup development environment"
	@echo "  help        - Show this help"