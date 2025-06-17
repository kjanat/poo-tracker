#!/bin/bash

# This script runs after the container is created
echo "🚽 Setting up Poo Tracker development environment..."

# Set proper permissions (ignore failures due to mounted volumes)
chmod +x backend/src/utils/seed.ts 2>/dev/null || true
chmod +x .devcontainer/devcontainer.sh 2>/dev/null || true
chmod +x .devcontainer/post-start.sh 2>/dev/null || true
chmod +x .devcontainer/uv-manager.sh 2>/dev/null || true

# Configure GPG
echo "🔐 Configuring GPG for commit signing..."

if [ -d ~/.gnupg ] && [ "$(ls -A ~/.gnupg 2>/dev/null)" ]; then
  echo "GPG keys found, configuring signing..."
  # Fix permissions for GPG (ignore failures on mounted volumes)
  chmod 700 ~/.gnupg 2>/dev/null || true
  chmod 600 ~/.gnupg/* 2>/dev/null || true

  # Configure git to use GPG signing if keys are available
  if gpg --list-secret-keys --keyid-format LONG | grep -q 'sec'; then
    GPG_KEY_ID=$(gpg --list-secret-keys --keyid-format LONG | grep 'sec' | head -1 | sed 's/.*\/\([A-F0-9]*\).*/\1/')
    git config --global user.signingkey "$GPG_KEY_ID"
    git config --global commit.gpgsign true
    git config --global tag.gpgsign true
    echo "✅ GPG signing configured with key: $GPG_KEY_ID"
  else
    echo "⚠️ No GPG private keys found, skipping signing configuration"
  fi
else
  echo "⚠️ No GPG directory found, skipping GPG configuration"
fi

# Set up Python UV virtual environment
echo "🐍 Setting up Python virtual environment with UV..."
if [ ! -d ".venv" ]; then
  echo "Creating new virtual environment..."
  uv venv .venv --python 3.13
else
  echo "Virtual environment already exists"
fi

# Activate virtual environment for this session
if [ -f ".venv/bin/activate" ]; then
  source .venv/bin/activate
  echo "✅ Virtual environment activated: $(which python)"
else
  echo "⚠️ Virtual environment activation script not found, using system Python: $(which python)"
fi

# Set up PNPM first
echo "📦 Setting up PNPM..."
pnpm setup || echo "⚠️ PNPM setup failed, continuing..."

# Install global dependencies
echo "📦 Installing global Node.js dependencies..."
pnpm install -g concurrently husky 2>/dev/null || echo "⚠️ Global package installation skipped (no global bin dir)"

# Install project dependencies
echo "📦 Installing project dependencies..."
pnpm install

# Install frontend dependencies
echo "🎨 Installing frontend dependencies..."
cd frontend && pnpm install && cd .. || exit 1

# Install backend dependencies
echo "🖥️ Installing backend dependencies..."
cd backend && pnpm install && cd .. || exit 1

# Install AI service dependencies with UV
echo "🤖 Installing AI service dependencies with UV..."
cd ai-service || exit 1
uv pip install -r requirements.txt
cd .. || exit 1

# Set up Git hooks
echo "🪝 Setting up Git hooks..."
pnpm prepare

# Set up environment files
echo "⚙️ Setting up environment files..."

# Backend .env
if [ ! -f "backend/.env" ]; then
  cat > backend/.env << EOF
# Database
DATABASE_URL="postgresql://poo_user:secure_password_123@localhost:5432/poo_tracker"

# Redis
REDIS_URL="redis://localhost:6379"

# S3 Storage
S3_ENDPOINT="http://localhost:9000"
S3_ACCESS_KEY="minioadmin"
S3_SECRET_KEY="minioadmin123"
S3_BUCKET="poo-photos"
S3_REGION="us-east-1"

# AI Service
AI_SERVICE_URL="http://localhost:8001"

# Authentication
JWT_SECRET="your-super-secret-jwt-key-change-in-production"
JWT_EXPIRES_IN="7d"

# Server
PORT=3001
NODE_ENV=development
CORS_ORIGIN="http://localhost:3000"
EOF
fi

# Frontend .env
if [ ! -f "frontend/.env" ]; then
  cat > frontend/.env << EOF
VITE_API_URL=http://localhost:3001
VITE_APP_NAME="Poo Tracker"
VITE_APP_VERSION="1.0.0"
EOF
fi

# AI Service .env
if [ ! -f "ai-service/.env" ]; then
  cat > ai-service/.env << EOF
REDIS_URL=redis://localhost:6379
ENVIRONMENT=development
LOG_LEVEL=info
EOF
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p backend/uploads
mkdir -p backend/logs
mkdir -p ai-service/logs

# Set proper permissions (ignore failures due to mounted volumes)
chmod +x backend/src/utils/seed.ts 2>/dev/null || true
chmod +x .devcontainer/devcontainer.sh 2>/dev/null || true
chmod +x .devcontainer/post-start.sh 2>/dev/null || true

echo "✅ Post-create setup completed!"
echo ""
echo "🎯 Next steps:"
echo "1. Run 'pnpm docker:up' to start the infrastructure services"
echo "2. Run 'pnpm db:migrate' to set up the database"
echo "3. Run 'pnpm db:seed' to seed the database (optional)"
echo "4. Run 'pnpm dev' to start the development servers"
echo ""
echo "Happy coding, you magnificent bastard! 💩"
