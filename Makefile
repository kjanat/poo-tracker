.PHONY: dev dev-frontend dev-backend dev-ai dev-full build build-frontend build-backend lint lint-frontend lint-backend lint-ai test test-frontend test-backend test-ai docker-up docker-down format format-backend format-ai clean

## Development
dev: dev-full

## Run frontend only
dev-frontend:
	pnpm --filter @poo-tracker/frontend run dev

## Run backend only
dev-backend:
	go run -C backend .

## Run AI service only
dev-ai:
	uv run uvicorn ai_service.main:app --reload

## Run all services
dev-full:
	pnpm --filter @poo-tracker/frontend run dev & \
	go run -C backend . & \
	uv run uvicorn ai_service.main:app --reload & \
	wait

## Build tasks
build: build-frontend build-backend

build-frontend:
	pnpm --filter @poo-tracker/frontend run build

build-backend:
	go build -C backend -o bin/server .

## Linting
lint: lint-frontend lint-backend lint-ai

lint-frontend:
	pnpm --filter @poo-tracker/frontend run lint

lint-backend:
	golangci-lint run ./backend/... || true

lint-ai:
	uv run ruff check ai-service

## Testing
test: test-frontend test-backend test-ai

test-frontend:
	pnpm --filter @poo-tracker/frontend run test

test-backend:
	go test -C backend -tags test \
	./internal/domain/... \
	./internal/repository/... \
	./internal/service/...

test-ai:
	uv run pytest ai-service

## Docker services
docker-up:
	docker compose up -d

docker-down:
	docker compose down

## Formatting
format:
	gofmt -s -w backend/**/*.go
	uv run ruff format ai-service

format-backend:
	gofmt -s -w backend/**/*.go

format-ai:
	uv run ruff format ai-service

clean:
	pnpm --filter @poo-tracker/frontend run clean
