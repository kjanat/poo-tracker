#!/bin/bash

# UV Virtual Environment Management Script for DevContainer
# This script helps manage Python virtual environments consistently
# between host and DevContainer environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
VENV_PATH="$PROJECT_ROOT/.venv"

case "${1:-help}" in
  "create")
    echo "🐍 Creating UV virtual environment..."
    cd "$PROJECT_ROOT"
    if [ -d "$VENV_PATH" ]; then
      echo "⚠️ Virtual environment already exists at $VENV_PATH"
      echo "Use 'recreate' to remove and recreate it"
      exit 1
    fi
    uv venv "$VENV_PATH" --python 3.13
    echo "✅ Virtual environment created at $VENV_PATH"
    ;;

  "recreate")
    echo "🔄 Recreating UV virtual environment..."
    cd "$PROJECT_ROOT"
    if [ -d "$VENV_PATH" ]; then
      rm -rf "$VENV_PATH"
      echo "🗑️ Removed existing virtual environment"
    fi
    uv venv "$VENV_PATH" --python 3.13
    echo "✅ Virtual environment recreated at $VENV_PATH"
    ;;

  "activate")
    echo "🔌 Activating virtual environment..."
    if [ ! -d "$VENV_PATH" ]; then
      echo "❌ Virtual environment not found at $VENV_PATH"
      echo "Run '$0 create' first"
      exit 1
    fi
    echo "Run: source $VENV_PATH/bin/activate"
    ;;

  "install")
    echo "📦 Installing AI service dependencies..."
    cd "$PROJECT_ROOT/ai-service"
    if [ ! -d "$VENV_PATH" ]; then
      echo "❌ Virtual environment not found"
      echo "Run '$0 create' first"
      exit 1
    fi
    # Use UV to install in the specific virtual environment
    VIRTUAL_ENV="$VENV_PATH" uv pip install -r requirements.txt
    echo "✅ Dependencies installed"
    ;;

  "sync")
    echo "🔄 Syncing AI service dependencies..."
    cd "$PROJECT_ROOT/ai-service"
    if [ ! -d "$VENV_PATH" ]; then
      echo "❌ Virtual environment not found"
      echo "Run '$0 create' first"
      exit 1
    fi
    # Use UV to sync dependencies (installs exact versions)
    VIRTUAL_ENV="$VENV_PATH" uv pip sync requirements.txt
    echo "✅ Dependencies synced"
    ;;

  "status")
    echo "📊 Virtual Environment Status:"
    echo "  Path: $VENV_PATH"
    if [ -d "$VENV_PATH" ]; then
      echo "  Exists: ✅ Yes"
      echo "  Python: $("$VENV_PATH"/bin/python --version 2> /dev/null || echo 'Not accessible')"
      echo "  UV Cache: ${UV_CACHE_DIR:-$HOME/.cache/uv}"
      if [ -f "$VENV_PATH/pyvenv.cfg" ]; then
        echo "  Config:"
        grep -E "home|version" "$VENV_PATH/pyvenv.cfg" | sed 's/^/    /'
      fi
      echo "  Packages:"
      if [ -x "$VENV_PATH/bin/pip" ]; then
        "$VENV_PATH/bin/pip" list --format=columns | head -10 | sed 's/^/    /'
        pkg_count=$("$VENV_PATH/bin/pip" list --format=freeze | wc -l)
        echo "    ... and $((pkg_count - 5)) more packages"
      fi
    else
      echo "  Exists: ❌ No"
    fi
    ;;

  "clean")
    echo "🧹 Cleaning UV cache and virtual environment..."
    cd "$PROJECT_ROOT"
    uv cache clean
    if [ -d "$VENV_PATH" ]; then
      rm -rf "$VENV_PATH"
      echo "🗑️ Removed virtual environment"
    fi
    echo "✅ Cleanup completed"
    ;;

  "test")
    echo "🧪 Testing AI service in virtual environment..."
    cd "$PROJECT_ROOT/ai-service"
    if [ ! -d "$VENV_PATH" ]; then
      echo "❌ Virtual environment not found"
      echo "Run '$0 create' first"
      exit 1
    fi
    echo "Running tests with virtual environment Python..."
    "$VENV_PATH/bin/python" -m pytest test_main.py -v
    ;;

  "help" | *)
    echo "🐍 UV Virtual Environment Manager for Poo Tracker"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  create     Create new virtual environment"
    echo "  recreate   Remove and recreate virtual environment"
    echo "  activate   Show activation command"
    echo "  install    Install AI service dependencies"
    echo "  sync       Sync dependencies (exact versions)"
    echo "  status     Show virtual environment status"
    echo "  clean      Clean cache and remove virtual environment"
    echo "  test       Run AI service tests"
    echo "  help       Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  UV_CACHE_DIR   UV cache directory (default: ~/.cache/uv)"
    echo "  VIRTUAL_ENV    Virtual environment path"
    echo ""
    echo "The virtual environment is stored in a Docker volume when"
    echo "running in DevContainer, ensuring consistency between"
    echo "container rebuilds while avoiding host path conflicts."
    ;;
esac
