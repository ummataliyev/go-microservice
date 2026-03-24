APP_NAME := go-microservice
BUILD_DIR := bin
MAIN_PKG := ./cmd/api

.PHONY: run build test test-unit test-integration lint fmt \
        migrate-up migrate-down migrate-create \
        docker-up-postgres docker-up-mysql docker-up-mongo

## run: Start the application
run:
	go run $(MAIN_PKG)

## build: Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

## test: Run all tests
test:
	go test ./... -v -race -count=1

## test-unit: Run unit tests only (exclude integration)
test-unit:
	go test ./... -v -race -count=1 -short

## test-integration: Run integration tests only
test-integration:
	go test ./... -v -race -count=1 -run Integration

## lint: Run linter
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	go fmt ./...
	goimports -w .

## migrate-up: Run database migrations up
migrate-up:
	go run ./cmd/migrate up

## migrate-down: Roll back last database migration
migrate-down:
	go run ./cmd/migrate down

## migrate-create: Create a new migration (usage: make migrate-create name=create_users)
migrate-create:
	go run ./cmd/migrate create $(name)

## docker-up-postgres: Start PostgreSQL in Docker
docker-up-postgres:
	docker compose up -d postgres

## docker-up-mysql: Start MySQL in Docker
docker-up-mysql:
	docker compose up -d mysql

## docker-up-mongo: Start MongoDB in Docker
docker-up-mongo:
	docker compose up -d mongo
