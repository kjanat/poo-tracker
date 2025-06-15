# AGENTS.md

## Overview

This file specifies instructions and standards for all AI agents, bots, GitHub Actions, and automated tools operating within the **poo-tracker** repository. All agents must follow these guidelines when reading, analyzing, changing, or generating code.

**Target Audience**: AI coding assistants, GitHub Actions, automated deployment tools, linting bots, and any other autonomous systems interacting with this codebase.

---

## Project Structure & Architecture

### 🏗️ pnpm Workspace Monorepo Layout

```tree
poo-tracker/
├── frontend/           # React + Vite + TypeScript + TailwindCSS v4
│   ├── src/
│   ├── package.json    # @poo-tracker/frontend
│   └── vite.config.ts
├── backend/            # Node.js + Express v5 + Prisma ORM
│   ├── src/
│   ├── prisma/
│   └── package.json    # @poo-tracker/backend
├── ai-service/         # Python FastAPI + Redis + ML/AI features
│   ├── main.py
│   ├── pyproject.toml  # Modern uv-compatible configuration
│   └── uv.lock
├── .github/            # GitHub Actions workflows and custom actions
│   ├── workflows/      # CI/CD pipelines
│   └── actions/        # Custom reusable actions (e.g., svg-converter)
├── branding/           # SVG assets and generated favicons/components
├── docker-compose.yml  # Development environment setup
├── pnpm-workspace.yaml # Workspace configuration
└── package.json        # Root workspace with shared configs
```

### 🎯 Key Technologies

- **Frontend**: React 19, Vite, TypeScript, TailwindCSS v4
- **Backend**: Node.js, Express v5, Prisma ORM, PostgreSQL
- **AI Service**: Python 3.9+, FastAPI, Redis, NumPy, Pandas, scikit-learn
- **Package Management**: pnpm 9+ (Node.js), uv (Python)
- **Infrastructure**: Docker, PostgreSQL, Redis, MinIO (S3-compatible)
- **Workspace**: pnpm workspaces for monorepo management

---

## Package Management & Workspace Commands

### 📦 pnpm Workspace Usage

**ALWAYS use pnpm workspace commands** - never `cd` into directories:

```bash
# ✅ Correct workspace commands
pnpm --filter @poo-tracker/frontend add react-router-dom
pnpm --filter @poo-tracker/backend run dev
pnpm --parallel run build
pnpm --recursive run test

# ❌ Incorrect - avoid cd commands in scripts
cd frontend && pnpm add react-router-dom
cd backend && pnpm run dev
```

### 🐍 Python Package Management

**ALWAYS use uv** - never pip:

```bash
# ✅ Correct uv commands
cd ai-service && uv sync
cd ai-service && uv add fastapi
uvx ruff format .
uvx black .

# ❌ Incorrect - avoid pip
pip install fastapi
python -m pip install -r requirements.txt
```

---

## Coding Conventions & Style Guide

### 🎨 Frontend (React/TypeScript)

- **TypeScript is mandatory** - No plain JavaScript allowed
- Use **functional components** with hooks
- **TailwindCSS v4 utility classes** - No CSS modules or external stylesheets
- Component naming: **PascalCase** (e.g., `UserDashboard.tsx`)
- File structure: Components in `src/components/`, pages in `src/pages/`
- Import order: React imports → Third-party → Local imports
- **Workspace-aware imports**: Use relative paths within workspace

### 🚀 Backend (Node.js/Express)

- **TypeScript is mandatory**
- Express v5 syntax - NO wildcard routes (`app.use('*', ...)`)
- Use Prisma ORM for all database operations
- RESTful API design - No GraphQL
- Environment variables in `.env` files (never commit secrets)
- Error handling with proper HTTP status codes
- **Workspace commands**: Use `pnpm --filter @poo-tracker/backend` for operations

### 🤖 AI Service (Python)

- **Python 3.9+** required
- Use `uv` and `uvx` for ALL package management (NOT `pip`)
- **[`pyproject.toml`](ai-service/pyproject.toml)** configuration (no requirements.txt)
- Follow **PEP 8** style guidelines
- **ruff** for linting and formatting (preferred over black + isort)
- Type hints required for all function signatures
- FastAPI for API endpoints
- Comprehensive error handling and logging

### 📦 Workspace Package Management

**Frontend/Backend** (pnpm workspace):

```bash
# Add dependencies to specific workspace
pnpm --filter @poo-tracker/frontend add axios
pnpm --filter @poo-tracker/backend add express-rate-limit

# Add dev dependencies to root workspace
pnpm add -Dw prettier eslint

# Run scripts on specific workspace
pnpm --filter @poo-tracker/frontend run build
pnpm --filter @poo-tracker/backend run test

# Run scripts on all workspaces
pnpm --parallel run dev
pnpm --recursive run build
```

**Python** (uv):

```bash
cd ai-service
uv add fastapi          # Add runtime dependency
uv add --dev pytest     # Add dev dependency
uv sync                 # Install all dependencies
uvx ruff format .       # Format code
uvx ruff check .        # Lint code
```

---

## Testing Protocols

### 🧪 Required Tests

- **Frontend**: Vitest for unit/component tests
- **Backend**: Vitest + Supertest for API testing
- **AI Service**: pytest + pytest-asyncio for FastAPI testing

### 🔬 Test Execution

**Use workspace commands**:

```bash
# Individual workspace testing
pnpm --filter @poo-tracker/frontend test
pnpm --filter @poo-tracker/backend test

# All workspaces
pnpm test                   # Runs all workspace tests
pnpm test:watch             # Watch mode for all workspaces

# AI service testing
cd ai-service && uv run pytest
```

### ✅ Test Requirements

- All new features MUST include tests
- Minimum 80% code coverage expected
- Tests must pass in CI before merge
- Mock external dependencies (Redis, PostgreSQL) appropriately
- Use workspace-aware test commands

---

## Automated Checks & CI/CD

### 🔍 Pre-Commit Requirements

All code changes must pass these checks:

**Frontend/Backend** (workspace commands):

```bash
pnpm lint              # ESLint across all workspaces
pnpm lint:fix          # Auto-fix linting issues
pnpm build             # Production build test
pnpm test              # All workspace tests
```

**AI Service**:

```bash
cd ai-service
uvx ruff format --check .  # Code formatting
uvx ruff check .           # Linting
uv run pytest              # Run tests
```

### 🚦 GitHub Actions Workflows

- **CI Pipeline** ([`.github/workflows/ci.yml`](.github/workflows/ci.yml)): Code formatting, linting, and testing
- **Format Pipeline** ([`.github/workflows/format.yml`](.github/workflows/format.yml)): Asset generation
- **SVG Converter** ([`.github/workflows/svg-convert.yml`](.github/workflows/svg-convert.yml)): Asset generation
- **Release Pipeline** ([`.github/workflows/release.yml`](.github/workflows/release.yml)): Deployment automation

### 🐳 Docker & Infrastructure

- Use `docker compose` (V2 syntax) - NOT `docker-compose`
- PostgreSQL for primary database
- Redis for caching and session storage
- MinIO for S3-compatible file storage
- Workspace commands: `pnpm docker:up`, `pnpm docker:down`

---

## Pull Request Guidelines

### 📝 Commit Message Format

Follow conventional commits with emojis:

```txt
feat: add Bristol stool chart selector
fix: correct off-by-one in shit counter
chore: update dependencies
```

### 🔥 PR Requirements

- **Title**: Use emoji + conventional commit format
- **Description**: Include what, why, and how
- **Tests**: All tests must pass
- **Reviews**: At least one approval required
- **Branch**: Use feature branches (`feat/`, `fix/`, `chore/`)
  - If only working on the frontend, use `frontend/` branch
  - If only working on the backend, use `backend/` branch
  - If only working on the AI service, use `ai/` branch
- **Target**: PRs to `master` branch
- **Workspace**: Use proper workspace commands in scripts

### 📋 PR Template Checklist

- [ ] Tests added/updated and passing
- [ ] Linting and formatting checks pass
- [ ] Used proper workspace commands (`pnpm --filter`, `uv` commands)
- [ ] Documentation updated if needed
- [ ] No breaking changes (or clearly documented)
- [ ] Environment variables documented in [`.env.example`](.env.example)

---

## Task & Workflow Definitions

### 🎨 Asset Management Agent

- **Trigger**: SVG file changes in `branding/` directory
- **Action**: Convert SVGs to ICO, PNG, React, React Native formats
- **Tool**: [`kjanat/svg-converter-action`](https://github.com/kjanat/svg-converter-action "GitHub: kjanat/svg-converter-action"), a custom GitHub Action
- **Output**: Auto-commit generated assets

### 🔧 Dependency Update Agent

- **Trigger**: Weekly schedule or manual dispatch
- **Action**: Update (p)npm/Python dependencies using workspace commands
- **Requirements**: Must run full test suite before committing
- **Commands**: Use `pnpm update` and `uv sync --upgrade`

### 🚀 Deployment Agent

- **Trigger**: Merge to `master` branch
- **Action**: Deploy to staging/production environments
- **Requirements**: All CI checks must pass
- **Build**: Use `pnpm build` for workspace builds

---

## Directory/Scope Rules

### 📁 File Organization

- **Root [`AGENTS.md`](AGENTS.md)**: Governs entire repository
- **[`pnpm-workspace.yaml`](pnpm-workspace.yaml)**: Defines workspace packages
- **Root [`package.json`](package.json)**: Shared scripts and workspace configuration
- **Service-specific rules**: Follow this file unless overridden by local configs
- **Generated files**: Never manually edit files in `branding/` (auto-generated)

### 🚨 Protected Areas

- **Never modify**: [`pnpm-lock.yaml`](pnpm-lock.yaml), [`uv.lock`](ai-service/uv.lock), generated migration files
- **Careful with**: Database schema changes, Docker configurations, workspace structure
- **Always validate**: Environment variable changes, workspace dependencies

---

## Error Handling & Escalation

### ❌ When Things Go Wrong

1. **Log detailed error messages** with context
2. **Include relevant file paths** and line numbers
3. **Capture command output** for debugging
4. **Check workspace configuration** for dependency issues
5. **Reference this AGENTS.md** for guidance

### 🆘 Escalation Process

- **Syntax Errors**: Fix automatically if possible, document in PR
- **Test Failures**: Stop and report - never bypass failing tests
- **Dependency Issues**: Check workspace conflicts, suggest updates
- **Workspace Issues**: Verify pnpm-workspace.yaml and package scoping
- **Unclear Requirements**: Flag for human review

### 📞 Contact Information

- **Repository Owner**: [@kjanat](https://github.com/kjanat)
- **Issues**: Use GitHub Issues for bugs/feature requests
- **Discussions**: Use GitHub Discussions for questions

---

## Golden Rules for Agents

### 💎 The Sacred Laws

1. **🚫 No broken builds** - If tests fail, fix or abort
2. **🎯 TypeScript everywhere** - No JavaScript allowed in frontend/backend
3. **🧪 Test everything** - No untested code
4. **📚 Document changes** - Update docs when needed
5. **🔒 Never commit secrets** - Use environment variables
6. **🎨 Follow the style** - Use project formatters and linters
7. **🔄 Use proper tools** - pnpm for Node.js, uv for Python
8. **🐳 Docker for dev** - Use docker-compose for local development
9. **🏗️ Respect workspace** - Use pnpm workspace commands exclusively
10. **📦 Modern tooling** - uv for Python, pnpm 9+ for Node.js

### 🔧 Workspace-Specific Rules

- **NEVER use `cd` in scripts** - Use `pnpm --filter` instead
- **Respect package scoping** - `@poo-tracker/frontend`, `@poo-tracker/backend`
- **Use workspace scripts** - `pnpm dev`, `pnpm build`, `pnpm test`
- **Share common configs** - TypeScript, Prettier, ESLint in root when possible
- **Python isolation** - Keep Python dependencies in [`ai-service/pyproject.toml`](ai-service/pyproject.toml) and use `uv` for management

### 🎉 Success Metrics

- All CI checks pass ✅
- Code coverage maintained or improved 📈
- No linting errors 🧹
- Workspace commands used correctly 🏗️
- Documentation updated 📖
- PR approved by human reviewer 👤

---

## Version & Updates

**Last Updated**: June 15, 2025  
**Version**: 2.0.0  
**Major Changes**: Added pnpm workspace configuration, uv Python tooling  
**Next Review**: When major architecture changes occur

---

_Remember: If your code stinks, it won't get merged. Use the workspace properly, keep it clean, tested, and documented!_ 💩✨
