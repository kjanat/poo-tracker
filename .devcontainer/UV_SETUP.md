# Python UV Virtual Environment Management

This guide explains how to manage Python virtual environments using UV in both host and DevContainer environments.

## 🐍 How UV Integration Works

The DevContainer uses Docker volumes to manage Python virtual environments:

- **Virtual Environment**: `/workspaces/poo-tracker/.venv` (inside container)
- **UV Cache**: `/home/kjanat/.cache/uv` (persistent volume)
- **Host Isolation**: Host and container use separate environments

## 📋 Key Benefits

1. **Consistency**: Same Python version and packages everywhere
2. **Isolation**: Container environment doesn't affect host
3. **Performance**: UV cache shared between container rebuilds
4. **Portability**: Works the same for all developers

## 🚀 Automatic Setup

The DevContainer automatically:

1. **Creates Virtual Environment**: Using `uv venv .venv --python 3.13`
2. **Installs Dependencies**: All AI service requirements
3. **Configures Environment**: Proper paths and activation
4. **Manages Cache**: Persistent UV cache across rebuilds

## 🛠️ Manual Management

Use the UV manager script for manual control:

```bash
# Create virtual environment
.devcontainer/uv-manager.sh create

# Recreate from scratch
.devcontainer/uv-manager.sh recreate

# Install dependencies
.devcontainer/uv-manager.sh install

# Check status
.devcontainer/uv-manager.sh status

# Run tests
.devcontainer/uv-manager.sh test

# Clean everything
.devcontainer/uv-manager.sh clean
```

## 🔧 Direct UV Commands

Inside the DevContainer:

```bash
# Activate virtual environment
source .venv/bin/activate

# Install packages
uv pip install package-name

# Install from requirements
uv pip install -r ai-service/requirements.txt

# Sync exact versions
uv pip sync ai-service/requirements.txt

# Add new dependency
cd ai-service
uv add package-name

# Remove dependency
uv remove package-name
```

## 📊 Environment Status

Check your environment status:

```bash
# Python version
python --version

# Virtual environment path
echo $VIRTUAL_ENV

# UV cache location
echo $UV_CACHE_DIR

# Installed packages
pip list

# Environment info
.devcontainer/uv-manager.sh status
```

## 🎯 Workflow Examples

### Adding a New Python Package

```bash
# In DevContainer
cd ai-service

# Add to requirements.txt
echo "numpy>=1.24.0" >> requirements.txt

# Install it
uv pip install -r requirements.txt

# Or use uv add (if using pyproject.toml)
uv add numpy
```

### Updating Dependencies

```bash
# Update requirements.txt with new versions
# Then sync to exact versions
uv pip sync requirements.txt

# Or update specific package
uv pip install --upgrade package-name
```

### Running AI Service

```bash
# Make sure virtual environment is active
source .venv/bin/activate

# Run the service
cd ai-service
python main.py

# Or run with uvicorn
uvicorn main:app --reload --host 0.0.0.0 --port 8000
```

## 🔍 Troubleshooting

### Virtual Environment Not Found

```bash
# Check if .venv exists
ls -la .venv/

# Recreate if missing
.devcontainer/uv-manager.sh recreate
```

### Wrong Python Version

```bash
# Check Python version
python --version

# Should be Python 3.13.x
# If not, recreate virtual environment
.devcontainer/uv-manager.sh recreate
```

### Package Installation Fails

```bash
# Check if virtual environment is active
echo $VIRTUAL_ENV

# Activate if needed
source .venv/bin/activate

# Try installing again
uv pip install -r ai-service/requirements.txt
```

### Cache Issues

```bash
# Clear UV cache
uv cache clean

# Or use manager script
.devcontainer/uv-manager.sh clean
```

## 💡 Best Practices

1. **Use Requirements Files**: Keep `requirements.txt` updated
2. **Pin Versions**: Use specific versions for reproducibility
3. **Test Locally**: Always test in the DevContainer environment
4. **Document Changes**: Update requirements when adding packages
5. **Use UV Commands**: Prefer `uv` over `pip` for better performance

## 🔄 Host vs Container Differences

| Aspect           | Host System        | DevContainer             |
| ---------------- | ------------------ | ------------------------ |
| Python Version   | Varies             | Python 3.13              |
| Virtual Env Path | `.venv/`           | `.venv/` (Docker volume) |
| UV Cache         | `~/.cache/uv`      | Volume-backed            |
| Package Manager  | uv/pip             | uv (preferred)           |
| Isolation        | Shared with system | Fully isolated           |

## 🚀 Performance Tips

1. **Use UV Cache**: Shared cache speeds up installations
2. **Volume Mounts**: Faster than bind mounts for virtual environments
3. **Batch Operations**: Install multiple packages at once
4. **Sync Instead of Install**: Use `uv pip sync` for exact versions

The DevContainer handles all the complexity of managing Python environments consistently across different development setups! 🐍✨
