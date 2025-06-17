#!/bin/bash

# DevContainer Management Script for Poo Tracker
# Usage: ./devcontainer.sh [command]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

case "${1:-help}" in
  "build")
    echo "🔨 Building DevContainer..."
    cd "$PROJECT_ROOT"
    docker build -t poo-tracker-devcontainer -f .devcontainer/Dockerfile .
    ;;

  "up")
    echo "🚀 Starting DevContainer services..."
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml up -d
    ;;

  "down")
    echo "🛑 Stopping DevContainer services..."
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml down
    ;;

  "logs")
    echo "📝 Showing service logs..."
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml logs -f
    ;;

  "status")
    echo "📊 Service Status:"
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml ps
    ;;

  "clean")
    echo "🧹 Cleaning up DevContainer..."
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml down -v
    docker volume prune -f
    echo "✅ Cleanup completed!"
    ;;

  "reset")
    echo "🔄 Resetting DevContainer environment..."
    cd "$PROJECT_ROOT"
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml down -v
    docker volume prune -f
    docker compose -f docker-compose.yml -f .devcontainer/docker-compose.dev.yml up -d
    echo "✅ Environment reset completed!"
    ;;

  "help" | *)
    echo "🚽 Poo Tracker DevContainer Management"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build    Build the DevContainer image"
    echo "  up       Start all services"
    echo "  down     Stop all services"
    echo "  logs     Show service logs"
    echo "  status   Show service status"
    echo "  clean    Stop services and remove volumes"
    echo "  reset    Clean and restart everything"
    echo "  help     Show this help message"
    echo ""
    echo "For VS Code DevContainer usage:"
    echo "1. Open project in VS Code"
    echo "2. Install 'Remote - Containers' extension"
    echo "3. Cmd/Ctrl+Shift+P → 'Reopen in Container'"
    echo ""
    echo "Happy coding, you magnificent bastard! 💩"
    ;;
esac
