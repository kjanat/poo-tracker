# Multi-Service Development Makefile
# Manages frontend (Node.js/pnpm), backend (Go), and AI service (Python/uv)

# ==== Project Configuration ====
SHELL := /bin/bash
.DEFAULT_GOAL := help

# Directories
FRONTEND_DIR := frontend
BACKEND_DIR := backend
AI_DIR := ai-service

# ==== Phony Targets ====
.PHONY: help init dev build lint lint-fix test docker format clean \
	init-frontend init-backend init-ai \
	dev-frontend dev-backend dev-ai dev-full \
	build-frontend build-backend build-ai \
	lint-frontend lint-backend lint-ai \
	lint-fix-frontend lint-fix-backend lint-fix-ai \
	test-frontend test-backend test-ai \
	docker-up docker-down \
	format-frontend format-backend format-ai \
	clean-frontend clean-backend clean-ai

# =============================================================================
# INITIALIZATION
# =============================================================================

init: init-frontend init-backend init-ai ## Initialize all services

init-frontend: ## Initialize frontend dependencies
	@echo "🔧 Initializing frontend dependencies..."
	cd $(FRONTEND_DIR) && pnpm install

init-backend: ## Initialize backend dependencies
	@echo "🔧 Initializing backend dependencies..."
	go mod tidy -C $(BACKEND_DIR)

init-ai: ## Initialize AI service environment
	@echo "🔧 Initializing AI service environment..."
	uv venv --directory $(AI_DIR)
	uv sync --directory $(AI_DIR)

# =============================================================================
# DEVELOPMENT
# =============================================================================

dev: dev-full ## Start all services in development mode

dev-frontend: ## Start frontend development server
	@echo "🚀 Starting frontend development server..."
	cd $(FRONTEND_DIR) && pnpm run dev

dev-backend: ## Start backend development server
	@echo "🚀 Starting backend development server..."
	cd $(BACKEND_DIR) && air

dev-ai: ## Start AI service development server
	@echo "🚀 Starting AI service development server..."
	cd $(AI_DIR) && uv run uvicorn ai_service.main:app --reload --host 0.0.0.0

dev-full: ## Start all services simultaneously
	@echo "🚀 Starting all services simultaneously..."
	@echo "📡 Frontend: http://localhost:3000"
	@echo "🔧 Backend: http://localhost:8080"
	@echo "🤖 AI Service: http://localhost:8000"
	@echo "Press Ctrl+C to stop all services"
	@echo ""
	@if command -v parallel >/dev/null 2>&1; then \
		parallel --lb --halt now,fail=1 ::: \
			"echo '🔧 [BACKEND] Starting...' && go run -C $(BACKEND_DIR) ./cmd/server/main.go" \
			"echo '🤖 [AI-SERVICE] Starting...' && cd $(AI_DIR) && uv run uvicorn ai_service.main:app --reload --host 0.0.0.0" \
			"echo '📱 [FRONTEND] Starting...' && cd $(FRONTEND_DIR) && pnpm run dev"; \
	else \
		echo "⚠️  GNU parallel not found, falling back to basic method..."; \
		echo "💡 Install parallel for better output: apt install parallel / brew install parallel"; \
		trap 'echo "🛑 Shutting down all services..."; kill 0' INT TERM; \
		(echo "🔧 [BACKEND] Starting..." && go run -C $(BACKEND_DIR) ./cmd/server/main.go) & \
		(echo "🤖 [AI-SERVICE] Starting..." && cd $(AI_DIR) && uv run uvicorn ai_service.main:app --reload --host 0.0.0.0) & \
		(echo "📱 [FRONTEND] Starting..." && cd $(FRONTEND_DIR) && pnpm run dev) & \
		wait; \
	fi

# =============================================================================
# BUILD
# =============================================================================

build: build-frontend build-backend build-ai ## Build all services

build-frontend: ## Build frontend
	@echo "📦 Building frontend..."
	cd $(FRONTEND_DIR) && pnpm run build

build-backend: ## Build backend
	@echo "📦 Building backend..."
	go build -C $(BACKEND_DIR) -o bin/server ./cmd/server/main.go

build-ai: ## Build AI service
	@echo "📦 Building AI service..."
	uv build --all-packages --directory $(AI_DIR)

# =============================================================================
# LINTING
# =============================================================================

lint: lint-frontend lint-backend lint-ai ## Lint all services

lint-frontend: ## Lint frontend code
	@echo "🔍 Linting frontend..."
	cd $(FRONTEND_DIR) && pnpm lint

lint-backend: ## Lint backend code
	@echo "🔍 Linting backend..."
	cd $(BACKEND_DIR) && golangci-lint run ./...

lint-ai: ## Lint AI service code
	@echo "🔍 Linting AI service..."
	uvx ruff check $(AI_DIR)

# =============================================================================
# LINT FIX
# =============================================================================

lint-fix: lint-fix-frontend lint-fix-backend lint-fix-ai ## Fix linting issues in all services

lint-fix-frontend: ## Fix frontend linting issues
	@echo "🔧 Fixing frontend lint issues..."
	cd $(FRONTEND_DIR) && pnpm lint:fix

lint-fix-backend: ## Fix backend linting issues
	@echo "🔧 Fixing backend lint issues..."
	cd $(BACKEND_DIR) && golangci-lint run --fix ./...

lint-fix-ai: ## Fix AI service linting issues
	@echo "🔧 Fixing AI service lint issues..."
	uvx ruff check $(AI_DIR) --fix

# =============================================================================
# TESTING
# =============================================================================

test: test-frontend test-backend test-ai ## Run tests for all services

test-frontend: ## Run frontend tests
	@echo "🧪 Running frontend tests..."
	cd $(FRONTEND_DIR) && pnpm test

test-backend: ## Run backend tests
	@echo "🧪 Running backend tests..."
	go test -C $(BACKEND_DIR) ./internal/domain/... ./internal/repository/... ./internal/service/...

test-ai: ## Run AI service tests
	@echo "🧪 Running AI service tests..."
	uv run pytest $(AI_DIR)

# =============================================================================
# DOCKER
# =============================================================================

docker-up: ## Start Docker services
	@echo "🐳 Starting Docker services..."
	docker compose up -d

docker-down: ## Stop Docker services
	@echo "🐳 Stopping Docker services..."
	docker compose down

docker: docker-up ## Alias for docker-up

# =============================================================================
# FORMATTING
# =============================================================================

format: format-frontend format-backend format-ai ## Format code for all services

format-frontend: ## Format frontend code
	@echo "✨ Formatting frontend code..."
	pnpm dlx prettier --write $(FRONTEND_DIR)

format-backend: ## Format backend code
	@echo "✨ Formatting backend code..."
	gofmt -s -w $(BACKEND_DIR)

format-ai: ## Format AI service code
	@echo "✨ Formatting AI service code..."
	uvx ruff format $(AI_DIR)

# =============================================================================
# CLEANUP
# =============================================================================

clean: clean-frontend clean-backend clean-ai ## Clean all build artifacts

clean-frontend: ## Clean frontend artifacts
	@echo "🧹 Cleaning frontend artifacts..."
	pnpm --package=rimraf dlx rimraf \
		**/dist **/build **/coverage **/out **/.next **/node_modules

clean-backend: ## Clean backend artifacts
	@echo "🧹 Cleaning backend artifacts..."
	go clean -testcache
	pnpm --package=rimraf dlx rimraf \
		**/bin/tmp **/tmp $(BACKEND_DIR)/coverage.txt $(BACKEND_DIR)/backend

clean-ai: ## Clean AI service artifacts
	@echo "🧹 Cleaning AI service artifacts..."
	pnpm --package=rimraf dlx rimraf \
		**/.venv **/.ruff_cache **/.mypy_cache **/.pytest_cache \
		**/.cache **/__pycache__/ **/__pycache__/* **/htmlcov \
		**/.coverage **/*.pyc **/*.pyo **/*.pyd
	uv cache clean --directory $(AI_DIR)

# =============================================================================
# HELP
# =============================================================================

help:
	@echo "Available commands:"
	@echo ""
	@echo "🔧 Setup:"
	@echo "  init                 Initialize all services"
	@echo "  init-frontend        Initialize frontend dependencies"
	@echo "  init-backend         Initialize backend dependencies"
	@echo "  init-ai              Initialize AI service environment"
	@echo ""
	@echo "🚀 Development:"
	@echo "  dev                  Start all services in development mode"
	@echo "  dev-frontend         Start frontend development server"
	@echo "  dev-backend          Start backend development server"
	@echo "  dev-ai               Start AI service development server"
	@echo ""
	@echo "📦 Build:"
	@echo "  build                Build all services"
	@echo "  build-frontend       Build frontend"
	@echo "  build-backend        Build backend"
	@echo "  build-ai             Build AI service"
	@echo ""
	@echo "🔍 Quality:"
	@echo "  lint                 Lint all services"
	@echo "  lint-fix             Fix linting issues in all services"
	@echo "  test                 Run tests for all services"
	@echo "  format               Format code for all services"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-up            Start Docker services"
	@echo "  docker-down          Stop Docker services"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  clean                Clean all build artifacts"
