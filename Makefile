APP_NAME := algoplatform
BIN_DIR := bin
APP_PORT := 8080
DOCKER_IMAGE := shohratd15/$(APP_NAME)
DOCKER_TAG := latest

MIGRATIONS_DIR := migrations
DB_URL := "postgres://user:password@algoplatform_db:5432/algoplatform_db?sslmode=disable"

.PHONY: build run test clean \
        docker-up docker-down docker-build docker-push \
        migrate-up migrate-down migrate-create

## ================================
## BUILD & RUN
## ================================
build:
	@echo "🚀 Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/main.go

run: build
	@echo "🏃 Running $(APP_NAME) on port $(PORT)..."
	@./$(BIN_DIR)/$(APP_NAME) -port $(PORT)

test:
	@echo "🧪 Testing $(APP_NAME)..."
	@go test -v ./...

clean:
	@echo "🧹 Cleaning $(APP_NAME)..."
	@rm -rf $(BIN_DIR)
	@go clean

## ================================
## DOCKER
## ================================
docker-up:
	@echo "🐳 Starting main docker-compose environment..."
	docker compose up -d --build

docker-down:
	@echo "🛑 Stopping main docker-compose environment..."
	docker compose down -v

docker-build:
	@echo "📦 Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push:
	@echo "⬆️ Pushing Docker image $(DOCKER_IMAGE):$(DOCKER_TAG) to Docker Hub..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

## ================================
## MIGRATIONS
## ================================
migrate-up:
	@echo "📜 Running migrations UP..."
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	@echo "📜 Running migrations DOWN..."
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: provide a migration name, e.g. make migrate-create name=create_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
