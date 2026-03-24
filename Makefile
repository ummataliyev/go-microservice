APP_NAME := go-microservice
BUILD_DIR := bin
MAIN_PKG := ./cmd/api
COMPOSE_FILE := infra/docker-compose.local.yml
DOCKER_IMAGE := $(APP_NAME)

.PHONY: run build test test-unit test-integration lint fmt swagger \
        migrate-up migrate-down migrate-create \
        docker-build docker-up docker-down docker-restart docker-logs docker-ps docker-clean

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

## swagger: Generate Swagger docs
swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

## migrate-up: Run database migrations up
migrate-up:
	go run ./cmd/migrate up

## migrate-down: Roll back last database migration
migrate-down:
	go run ./cmd/migrate down

## migrate-create: Create a new migration (usage: make migrate-create name=create_users)
migrate-create:
	go run ./cmd/migrate create $(name)

## docker-build: Build the Docker image
docker-build:
	docker build -f infra/Dockerfile -t $(DOCKER_IMAGE) .

## docker-up: Start all services
docker-up:
	docker compose -f $(COMPOSE_FILE) up -d --build

## docker-down: Stop all services
docker-down:
	docker compose -f $(COMPOSE_FILE) down

## docker-restart: Restart all services
docker-restart: docker-down docker-up

## docker-logs: Tail logs from all services
docker-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

## docker-ps: Show running containers
docker-ps:
	docker compose -f $(COMPOSE_FILE) ps

## docker-clean: Stop services and remove volumes
docker-clean:
	docker compose -f $(COMPOSE_FILE) down -v --remove-orphans
