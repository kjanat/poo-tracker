#!/bin/bash

# This script runs every time the container starts
echo "🚽 Starting Poo Tracker development environment..."

# Wait for Docker daemon to be ready
echo "🐳 Waiting for Docker daemon..."
while ! docker info > /dev/null 2>&1; do
  echo "Docker daemon not ready, waiting..."
  sleep 2
done

echo "✅ Docker daemon is ready!"

# Check if infrastructure services are running
if ! docker compose ps | grep -q "Up"; then
  echo "🏗️ Infrastructure services not running. Starting them..."
  echo "💡 Tip: You can also run 'pnpm docker:up' manually"

  # Start infrastructure services with development configuration
  docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml up -d

  # Wait for services to be ready
  echo "⏳ Waiting for services to be ready..."

  # Wait for PostgreSQL
  echo "🐘 Waiting for PostgreSQL..."
  while ! docker exec poo-tracker-postgres pg_isready -U poo_user -d poo_tracker > /dev/null 2>&1; do
    sleep 2
  done

  # Wait for Redis
  echo "🔴 Waiting for Redis..."
  while ! docker exec poo-tracker-redis redis-cli ping > /dev/null 2>&1; do
    sleep 2
  done

  # Wait for MinIO
  echo "🪣 Waiting for MinIO..."
  while ! curl -f http://localhost:9000/minio/health/live > /dev/null 2>&1; do
    sleep 2
  done

  echo "✅ All infrastructure services are ready!"
else
  echo "✅ Infrastructure services are already running!"
fi

# Check if database migrations need to be run
echo "🔍 Checking database status..."
if ! docker exec poo-tracker-postgres psql -U poo_user -d poo_tracker -c "SELECT 1 FROM _prisma_migrations LIMIT 1;" > /dev/null 2>&1; then
  echo "🗃️ Running database migrations..."
  cd backend && pnpm run db:migrate && cd ..
else
  echo "✅ Database is up to date!"
fi

# Show service status
echo ""
echo "🎯 Service Status:"
echo "  Frontend:  http://localhost:3000"
echo "  Backend:   http://localhost:3001"
echo "  AI Service: http://localhost:8001"
echo "  PostgreSQL: localhost:5432"
echo "  Redis:     localhost:6379"
echo "  MinIO API: http://localhost:9000"
echo "  MinIO Console: http://localhost:9002"
echo ""
echo "🚀 Ready to start development!"
echo "   Run 'pnpm dev' to start the development servers"
