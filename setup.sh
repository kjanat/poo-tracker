#!/bin/bash
set -e

echo "🚀 Setting up Poo Tracker development environment..."

    # Install pnpm if not already installed
    if ! command -v pnpm &> /dev/null; then
        echo "📦 Installing pnpm..."
        npm install -g pnpm@latest
    fi

    # Install uv if not already installed
    if ! command -v uv &> /dev/null; then
        echo "🐍 Installing uv..."
        curl -LsSf https://astral.sh/uv/install.sh | sh
        export PATH="$HOME/.cargo/bin:$PATH"
    fi

# Install Node.js dependencies (frontend)
echo "📦 Installing Node.js dependencies..."
pnpm install

# Install Python dependencies for AI service
echo "🐍 Installing Python dependencies..."
cd ai-service && uv sync && cd ..

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    echo "⚙️ Setting up environment variables..."
    cp .env.example .env
    echo "✅ Created .env file - please update with your settings"
fi

# Start Docker services
echo "🐳 Starting Docker services..."
make docker-up

# Wait for PostgreSQL to be ready using pg_isready
echo "⏳ Waiting for PostgreSQL to be ready..."
until docker compose exec -T postgres pg_isready -U poo_user -d poo_tracker >/dev/null 2>&1; do
    echo "🔄 Postgres not ready yet, retrying in 1s..."
    sleep 1
done

# No automatic migrations for Go backend

echo "✅ Setup complete! Run 'make dev' to start all services"
echo ""
echo "🌐 Services will be available at:"
echo "  Frontend:      http://localhost:5173"
echo "  Backend:       http://localhost:3002" 
echo "  AI Service:    http://localhost:8001"
echo "  MinIO Console: http://localhost:9002"
