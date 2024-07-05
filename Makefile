# List all your services or packages here
SERVICES := api

.PHONY: all lint test build tools clean docker-up docker-down docker-test $(SERVICES)

# Default target
all: sqlc lint docker-test build

# generate sql
sqlc:
	@echo "Running sqlc..."
	@for service in $(SERVICES); do \
		echo "sqlc $$service"; \
		cd $$service && sqlc generate && cd ..; \
	done

# Lint all services
lint:
	@echo "Running linter..."
	@for service in $(SERVICES); do \
		echo "Linting $$service"; \
		cd $$service && golangci-lint run && cd ..; \
	done

# Test all services
test:
	@echo "Running tests..."
	@for service in $(SERVICES); do \
		echo "Testing $$service"; \
		cd $$service; \
		DB_CONNECTION_STRING=postgresql://postgres:templatepass@localhost:5432/postgres?sslmode=disable \
		go test -v ./...; \
		cd ..; \
	done

# Build all services
build:
	@echo "Building services..."
	mkdir -p ./bin
	@for service in $(SERVICES); do \
		echo "Building $$service"; \
		cd $$service && go build -o ../bin/$$service ./cmd && cd ..; \
	done

# Install tools
tools:
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go get github.com/sqlc-dev/pqtype@latest

# Clean up
clean:
	@echo "Cleaning..."
	go clean
	rm -rf bin/


# Start PostgreSQL container
docker-up:
	@echo "Starting PostgreSQL container..."
	@for service in $(SERVICES); do \
		echo "Copy schema $$service"; \
		cp $$service/db/migrations/*.sql ./test-db-init-scripts; \
	done
	docker-compose up -d

# Stop and remove PostgreSQL container
docker-down:
	@echo "Stopping PostgreSQL container..."
	docker-compose down

# Run tests with Docker PostgreSQL
docker-test: docker-up
	@echo "Running tests with Docker PostgreSQL..."
	@sleep 5  # Wait for PostgreSQL to be ready
	make test
	@make docker-down

# Individual service targets
$(SERVICES):
	@echo "Processing $@..."
	@cd $@ && golangci-lint run && go test -v ./... && go build -o ../bin/$@ && cd ..
