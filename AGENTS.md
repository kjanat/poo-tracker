# AGENTS.md

## Overview

This file specifies instructions and standards for all AI agents, bots, GitHub Actions, and automated tools operating within the **poo-tracker** repository. All agents must follow these guidelines when reading, analyzing, changing, or generating code.

**Target Audience**: AI coding assistants, GitHub Actions, automated deployment tools, linting bots, and any other autonomous systems interacting with this codebase.

---

## Project Structure & Architecture

### 🏗️ Monorepo Layout

```tree
poo-tracker/
├── frontend/          # React + Vite + TypeScript + TailwindCSS
├── backend/           # Node.js + Express v5 + Prisma ORM
├── ai-service/        # Python FastAPI + Redis + ML/AI features
├── .github/           # GitHub Actions workflows and custom actions
│   ├── workflows/     # CI/CD pipelines
│   └── actions/       # Custom reusable actions (e.g., svg-converter)
├── branding/          # SVG assets and generated favicons/components
├── docker-compose.yml # Development environment setup
└── *.md               # Documentation files
```

### 🎯 Key Technologies

- **Frontend**: React 18, Vite, TypeScript, TailwindCSS v4
- **Backend**: Node.js, Express v5, Prisma ORM, PostgreSQL
- **AI Service**: Python, FastAPI, Redis, NumPy, Pandas
- **Package Management**: pnpm (Node.js), uv (Python)
- **Infrastructure**: Docker, PostgreSQL, Redis, MinIO (S3-compatible)

---

## Coding Conventions & Style Guide

### 🎨 Frontend (React/TypeScript)

- **TypeScript is mandatory** - No plain JavaScript allowed
- Use **functional components** with hooks
- **TailwindCSS utility classes** - No CSS modules or external stylesheets
- Component naming: **PascalCase** (e.g., `UserDashboard.tsx`)
- File structure: Components in `src/components/`, pages in `src/pages/`
- Import order: React imports → Third-party → Local imports

### 🚀 Backend (Node.js/Express)

- **TypeScript is mandatory**
- Express v5 syntax - NO wildcard routes (`app.use('*', ...)`)
- Use Prisma ORM for all database operations
- RESTful API design - No GraphQL
- Environment variables in `.env` files (never commit secrets)
- Error handling with proper HTTP status codes

### 🤖 AI Service (Python)

- **Python 3.11+** required
- Use `uv` and `uvx` for package management (NOT `pip`)
- Follow **PEP 8** style guidelines
- **Black** for code formatting
- **Ruff** for linting and additional formatting
- **isort** for import sorting
- Type hints required for all function signatures
- FastAPI for API endpoints
- Comprehensive error handling and logging

### 📦 Package Management

- **Frontend/Backend**: Use `pnpm` exclusively
- **Python**: Use `uv pip install` and `uvx` for tool execution
- **Workspaces**: Leverage pnpm workspaces for monorepo management

---

## Testing Protocols

### 🧪 Required Tests

- **Frontend**: Vitest for unit/component tests
- **Backend**: Jest + Supertest for API testing
- **AI Service**: pytest + pytest-asyncio for FastAPI testing

### 🔬 Test Execution

Before making changes:

```bash
# Frontend tests
pnpm run --filter frontend test

# Backend tests
pnpm run --filter backend test

# AI service tests
cd ai-service && python -m pytest

# Full test suite
pnpm test # Runs all workspace tests
```

### ✅ Test Requirements

- All new features MUST include tests
- Minimum 80% code coverage expected
- Tests must pass in CI before merge
- Mock external dependencies (Redis, PostgreSQL) appropriately

---

## Automated Checks & CI/CD

### 🔍 Pre-Commit Requirements

All code changes must pass these checks:

**Frontend/Backend**:

```bash
pnpm run lint  # ESLint
pnpm run tsc   # TypeScript compilation
pnpm run build # Production build test
```

**AI Service**:

```bash
uvx black --check .       # Code formatting
uvx ruff check .          # Linting
uvx ruff format --check . # Format verification
```

### 🚦 GitHub Actions Workflows

- **CI Pipeline** (`.github/workflows/ci.yml`): Runs tests, linting, builds
- **Format Pipeline** (`.github/workflows/format.yml`): Code formatting checks
- **SVG Converter** (`.github/workflows/svg-convert.yml`): Asset generation
- **Release Pipeline** (`.github/workflows/release.yml`): Deployment automation

### 🐳 Docker & Infrastructure

- Use `docker compose` (V2 syntax) - NOT `docker-compose`
- PostgreSQL for primary database
- Redis for caching and session storage
- MinIO for S3-compatible file storage

---

## Pull Request Guidelines

### 📝 Commit Message Format

Follow conventional commits with emojis:

```txt
feat: add Bristol stool chart selector
fix: correct off-by-one in shit counter
chore: update dependencies
🎨 feat: add SVG Converter Pro GitHub Action
```

### 🔥 PR Requirements

- **Title**: Use emoji + conventional commit format
- **Description**: Include what, why, and how
- **Tests**: All tests must pass
- **Reviews**: At least one approval required
- **Branch**: Use feature branches (`feat/`, `fix/`, `chore/`)
- **Target**: PRs to `master` branch

### 📋 PR Template Checklist

- [ ] Tests added/updated and passing
- [ ] Linting and formatting checks pass
- [ ] Documentation updated if needed
- [ ] No breaking changes (or clearly documented)
- [ ] Environment variables documented in `.env.example`

---

## Task & Workflow Definitions

### 🎨 Asset Management Agent

- **Trigger**: SVG file changes in `branding/` directory
- **Action**: Convert SVGs to ICO, PNG, React, React Native formats
- **Tool**: `.github/actions/svg-converter`
- **Output**: Auto-commit generated assets

### 🔧 Dependency Update Agent

- **Trigger**: Weekly schedule or manual dispatch
- **Action**: Update npm/Python dependencies
- **Requirements**: Must run full test suite before committing

### 🚀 Deployment Agent

- **Trigger**: Merge to `master` branch
- **Action**: Deploy to staging/production environments
- **Requirements**: All CI checks must pass

---

## Directory/Scope Rules

### 📁 File Organization

- **Root AGENTS.md**: Governs entire repository
- **Service-specific rules**: Follow this file unless overridden by local configs
- **Generated files**: Never manually edit files in `branding/` (auto-generated)

### 🚨 Protected Areas

- **Never modify**: `pnpm-lock.yaml`, generated migration files
- **Careful with**: Database schema changes, Docker configurations
- **Always validate**: Environment variable changes

---

## Error Handling & Escalation

### ❌ When Things Go Wrong

1. **Log detailed error messages** with context
2. **Include relevant file paths** and line numbers
3. **Capture command output** for debugging
4. **Reference this AGENTS.md** for guidance

### 🆘 Escalation Process

- **Syntax Errors**: Fix automatically if possible, document in PR
- **Test Failures**: Stop and report - never bypass failing tests
- **Dependency Issues**: Check for version conflicts, suggest updates
- **Unclear Requirements**: Flag for human review

### 📞 Contact Information

- **Repository Owner**: [@kjanat](https://github.com/kjanat)
- **Issues**: Use GitHub Issues for bugs/feature requests
- **Discussions**: Use GitHub Discussions for questions

---

## Golden Rules for Agents

### 💎 The Sacred Laws

1. **🚫 No broken builds** - If tests fail, fix or abort
2. **🎯 TypeScript everywhere** - No JavaScript allowed
3. **🧪 Test everything** - No untested code
4. **📚 Document changes** - Update docs when needed
5. **🔒 Never commit secrets** - Use environment variables
6. **🎨 Follow the style** - Use project formatters and linters
7. **🔄 Use proper tools** - pnpm for Node.js, uv for Python
8. **🐳 Docker for dev** - Use docker-compose for local development

### 🎉 Success Metrics

- All CI checks pass ✅
- Code coverage maintained or improved 📈
- No linting errors 🧹
- Documentation updated 📖
- PR approved by human reviewer 👤

---

## Version & Updates

**Last Updated**: June 14, 2025  
**Version**: 1.0.0  
**Next Review**: When major architecture changes occur

---

_Remember: If your code stinks, it won't get merged. Keep it clean, tested, and documented!_ 💩✨
