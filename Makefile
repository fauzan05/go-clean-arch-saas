.PHONY: help setup migrate-up migrate-down migrate-create seed test test-coverage dev build clean docker-up docker-down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Initial setup (copy .env, install deps, create db, run migrations)
	@echo "🚀 Setting up project..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "✅ Created .env file"; fi
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

migrate-up: ## Run database migrations
	@echo "📦 Running migrations..."
	@migrate -database "mysql://$(shell grep DB_USERNAME .env | cut -d '=' -f2):$(shell grep DB_PASSWORD .env | cut -d '=' -f2)@tcp($(shell grep DB_HOST .env | cut -d '=' -f2):$(shell grep DB_PORT .env | cut -d '=' -f2))/$(shell grep DB_NAME .env | cut -d '=' -f2)?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up

migrate-down: ## Rollback last migration
	@echo "⏮️  Rolling back migration..."
	@migrate -database "mysql://$(shell grep DB_USERNAME .env | cut -d '=' -f2):$(shell grep DB_PASSWORD .env | cut -d '=' -f2)@tcp($(shell grep DB_HOST .env | cut -d '=' -f2):$(shell grep DB_PORT .env | cut -d '=' -f2))/$(shell grep DB_NAME .env | cut -d '=' -f2)?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations down 1

migrate-create: ## Create new migration (usage: make migrate-create name=create_table_xyz)
	@migrate create -ext sql -dir db/migrations -seq $(name)

seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@mysql -u$(shell grep DB_USERNAME .env | cut -d '=' -f2) -p$(shell grep DB_PASSWORD .env | cut -d '=' -f2) $(shell grep DB_NAME .env | cut -d '=' -f2) < scripts/seed.sql
	@echo "✅ Database seeded"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./test/

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./test/
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

dev: ## Run with hot reload (requires air)
	@if ! command -v air &> /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air

build: ## Build binary
	@echo "🔨 Building..."
	@go build -o bin/app cmd/web/main.go
	@echo "✅ Binary created: bin/app"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/ coverage.out coverage.html
	@echo "✅ Cleaned"

docker-up: ## Start docker containers
	@echo "🐳 Starting containers..."
	@docker-compose up -d
	@echo "✅ Containers started"

docker-down: ## Stop docker containers
	@echo "🐳 Stopping containers..."
	@docker-compose down
	@echo "✅ Containers stopped"

docker-build: ## Build docker image
	@echo "🐳 Building docker image..."
	@docker build -t go-saas-starter:latest .
	@echo "✅ Image built"

.DEFAULT_GOAL := help
