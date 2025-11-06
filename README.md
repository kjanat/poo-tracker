# 💩 Poo Tracker

<img align="right" width="150" height="150" src="branding/logo.svg">

Ever wondered if your gut's on a winning streak, or if your last kebab is about to take its revenge?
Welcome to **Poo Tracker**, the brutally honest app that lets you log every majestic turd, rabbit
pellet, and volcanic diarrhea eruption without shame or censorship. Because your bowel movements
_do_ matter, and no, your GP isn't getting this level of data.

## 🚀 Quick Start

### Prerequisites

- [Node.js](https://nodejs.org/) 22+
- [Docker](https://docs.docker.com/engine/) & [Docker Compose](https://docs.docker.com/compose/)
- [pnpm](https://pnpm.io/) 9+ (because npm is for amateurs)
- [uv](https://docs.astral.sh/uv/) for Python (because pip is for the weak)

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/kjanat/poo-tracker.git
   cd poo-tracker
   ```

1. **Install dependencies:**

   ```bash
   # Install Node.js dependencies (frontend + backend)
   pnpm install

   # Install Python dependencies for AI service
   cd ai-service && uv sync && cd ..
   ```

1. **Set up environment variables:**

   ```bash
   cp .env.example .env
   cp frontend/.env.example frontend/.env.local
   cp backend/.env.example backend/.env
   cp ai-service/.env.example ai-service/.env
   # Edit each file with your local values
   ```

1. **Start the services:**

   ```bash
   # Start database and supporting services
   pnpm docker:up

   # Start all development servers (frontend + backend + AI)
   pnpm dev:full

   # OR start just frontend + backend
   pnpm dev
   ```

1. **Set up the database:**

   ```bash
   # Run database migrations
   pnpm db:migrate

   # (Optional) Seed with test data
   pnpm db:seed
   ```

   A sample set of credentials for API testing is provided in
   [`login.example.json`](login.example.json). Copy this file to `login.json` if you need a local
   login file that's ignored by Git.

1. **Open your browser:**

- Frontend: <http://localhost:5173>
- Backend API: <http://localhost:3002>
- AI Service: <http://localhost:8001>
- MinIO Console: <http://localhost:9002> (minioadmin/minioadmin123)

## 🏗️ Tech Stack

- **Frontend**: React 19 + Vite 6 + TypeScript 5.9 + TailwindCSS v4
- **Backend**: Node.js 22 + Express v5 + TypeScript 5.9 + Prisma v6.10
- **Database**: PostgreSQL 17
- **Storage**: MinIO (S3-compatible for photos)
- **AI Service**: Python 3.9+ + FastAPI + scikit-learn
- **Infrastructure**: Docker + Docker Compose
- **Package Management**: pnpm 9+ (Node.js) + uv (Python)
- **Container Registry**: GitHub Container Registry (GHCR)

## 🗝️ Environment Variables

Each package ships with an `.env.example` file. Copy them and tweak the values before running
anything:

```bash
cp .env.example .env
cp frontend/.env.example frontend/.env.local
cp backend/.env.example backend/.env
cp ai-service/.env.example ai-service/.env
```

These example files contain **sample credentials only**. Replace them with your real secrets and
**never commit private keys or passwords** to the repository.

## 📖 What's the fucking point?

- **Track your shits**: Log every glorious bowel movement in a timeline worthy of its own Netflix
  docu-series. Date, time, volume, Bristol score, color, floaters—you name it.
- **Photo evidence**: Snap a pic for science (or just to traumatize your friends and family). Upload
  and catalog the full visual glory of your rectal output.
- **Meal log**: Record what you've stuffed down your throat so you can finally prove that Taco
  Tuesday was, in fact, a war crime against your colon.
- **AI-powered analysis**: Our cold, heartless machine learning model doesn't judge—just coldly
  correlates your food, stress, and fiber intake with your latest gut grenade.
- **Export & share**: PDF reports for your doctor, gastroenterologist, or weird fetish community.
  It's your data, so overshare as you wish.

## 🔧 Development

### Available Scripts

```bash
# Development - All Services
pnpm dev:full         # Start frontend + backend + AI service
pnpm dev              # Start frontend + backend only
pnpm dev:frontend     # Start frontend only
pnpm dev:backend      # Start backend only
pnpm dev:ai           # Start AI service only

# Building
pnpm build            # Build all projects
pnpm build:frontend   # Build frontend only
pnpm build:backend    # Build backend only

# Database Operations
pnpm db:migrate       # Run Prisma migrations
pnpm db:seed          # Seed database with test data
pnpm db:studio        # Open Prisma Studio
pnpm db:reset         # Reset database (⚠️ destructive!)

# Docker Services
pnpm docker:up        # Start PostgreSQL, Redis, MinIO
pnpm docker:down      # Stop all Docker services

# Testing & Quality
pnpm test             # Run all tests (frontend + backend)
pnpm test:watch       # Run tests in watch mode
pnpm typecheck        # Run TypeScript type checking
pnpm lint             # Run linters on all projects
pnpm lint:fix         # Auto-fix linting issues
pnpm clean            # Clean all build artifacts

# Code Formatting
pnpm format           # Format all files (Node.js + Python)
pnpm format:all       # Format all JavaScript/TypeScript files
```

### Project Structure

```graphql
poo-tracker/
├── frontend/           # React frontend (Vite + TypeScript + TailwindCSS)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── backend/            # Express.js API (TypeScript + Prisma)
│   ├── src/
│   ├── prisma/
│   └── package.json
├── ai-service/         # Python FastAPI AI service
│   ├── main.py
│   ├── pyproject.toml
│   └── uv.lock
├── docker-compose.yml  # Infrastructure setup
├── pnpm-workspace.yaml # Workspace configuration
└── package.json        # Root workspace config
```

### Workspace Management

This project uses **pnpm workspaces** for efficient monorepo management:

```bash
# Install dependencies for specific workspace
pnpm --filter @poo-tracker/frontend add react-router-dom
pnpm --filter @poo-tracker/backend add express-rate-limit

# Run commands on specific workspaces
pnpm --filter @poo-tracker/frontend build
pnpm --filter @poo-tracker/backend test

# Run commands on all workspaces
pnpm --recursive run build
pnpm --parallel run dev
```

### API Endpoints

| Endpoint                 | Method   | Description                           |
| ------------------------ | -------- | ------------------------------------- |
| `/api/auth/register`     | `POST`   | Register a new user                   |
| `/api/auth/login`        | `POST`   | Authenticate user                     |
| `/api/entries`           | `GET`    | Get bowel movement entries            |
| `/api/entries`           | `POST`   | Create a new bowel movement entry     |
| `/api/entries/:id`       | `PUT`    | Update a bowel movement entry         |
| `/api/entries/:id`       | `DELETE` | Delete a bowel movement entry         |
| `/api/uploads/photo`     | `POST`   | Upload a photo of your bowel movement |
| `/api/analytics/summary` | `GET`    | Get AI analysis summary               |

## 📝 How it works

1. **Eat** something questionable.
1. **Shit**. (Preferably in a toilet, but we're not here to kink-shame.)
1. **Log** your experience in Poo Tracker—optionally add a photo, rate your suffering, describe the
   smell if you're brave enough.
1. **Repeat** until you finally realize you're lactose intolerant or have IBS.

## 🔒 Privacy, baby

We encrypt your brown notes and hide them away. Nobody's reading your logs except you—and whatever
godforsaken AI wants to learn about the day-after effects of your sushi buffet.

## 👥 Who is this for?

- People with IBS, Crohn, colitis, or "nobody shits like me" syndrome.
- Control freaks and data lovers.
- Curious bastards who just want to know if spicy food really _is_ their nemesis.

## 🎯 Features

### Bristol Stool Chart Integration

- Complete 7-type Bristol Stool Chart support
- Visual indicators and descriptions
- Trend analysis over time

### AI-Powered Insights

- Pattern recognition in bowel movements
- Meal correlation analysis
- Health recommendations
- Risk factor identification

### Photo Documentation

- Secure photo upload and storage
- Image compression and optimization
- Privacy-focused MinIO storage

### Export & Analytics

- Detailed charts and visualizations
- Data export capabilities
- Comprehensive health insights

## 🤝 Contributing

1. Fork the repository
1. Create a feature branch (`git checkout -b feat/amazing-feature`)
1. Follow the coding standards (see [AGENTS.md](AGENTS.md))
1. Write tests for new features
1. Commit using conventional commits (`feat: add Bristol chart selector`)
1. Push and create a Pull Request

### Coding Standards

- TypeScript everywhere (no JS allowed)
- Component-based architecture
- TailwindCSS for styling (no CSS modules)
- RESTful API design
- Comprehensive test coverage
- ESLint with @typescript-eslint and Prettier (follow the config, don't "fix" it)
- Use pnpm workspace commands for consistent development

### Pre-Commit Hooks

This project uses both `husky` with `lint-staged` and supports `pre-commit` framework:

#### Option 1: Husky (Automatically Installed)

Husky hooks are automatically set up when you run `pnpm install`. They will:

- Type check changed TypeScript files
- Format changed files with Prettier
- Lint and format Python files with Ruff

#### Option 2: Pre-Commit Framework (Optional)

For more comprehensive checks, install the `pre-commit` framework:

```bash
# Install pre-commit (choose one)
pip install pre-commit
# or
pipx install pre-commit
# or
uvx --from pre-commit pre-commit install

# Install the git hooks
pre-commit install
pre-commit install --hook-type commit-msg

# Run on all files (optional)
pre-commit run --all-files
```

The pre-commit hooks include:

- Trailing whitespace removal
- End-of-file fixing
- YAML/JSON/TOML validation
- Large file detection
- Merge conflict detection
- Secret detection
- Markdown formatting
- Shell script linting (shellcheck)
- Docker file linting (hadolint)
- Conventional commit validation

## 🚀 Deployment

### Docker Images

Pre-built Docker images are automatically published to GitHub Container Registry (GHCR) on every
push to `master`:

```bash
# Pull the latest AI service image
docker pull ghcr.io/kjanat/poo-tracker/ai-service:latest

# Run with Docker Compose (uses published images)
docker compose up -d
```

**Image Features:**

- ✅ Multi-architecture support (amd64, arm64)
- ✅ SBOM (Software Bill of Materials) attestation
- ✅ Provenance attestation for supply chain security
- ✅ Signed with Cosign (keyless signing via Sigstore)
- ✅ Automated security scanning
- ✅ OCI-compliant metadata and labels

**Verify Image Signature:**

```bash
# Install cosign
brew install cosign  # macOS
# or: https://docs.sigstore.dev/cosign/installation/

# Verify the image signature
cosign verify ghcr.io/kjanat/poo-tracker/ai-service:latest \
  --certificate-identity-regexp="https://github.com/kjanat/poo-tracker" \
  --certificate-oidc-issuer=https://token.actions.githubusercontent.com
```

**View SBOM:**

```bash
# Download and view the SBOM
docker buildx imagetools inspect ghcr.io/kjanat/poo-tracker/ai-service:latest \
  --format "{{ json .SBOM }}"
```

### Environment Variables

```env
# Database
DATABASE_URL="postgresql://poo_user:secure_password_123@localhost:5432/poo_tracker"

# Storage (MinIO)
MINIO_ENDPOINT="localhost:9000"
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY="minioadmin123"

# API Configuration
API_PORT=3002
CORS_ORIGIN="http://localhost:5173"

# AI Service
AI_SERVICE_URL="http://localhost:8001"

# JWT
JWT_SECRET="your-super-secret-jwt-key-change-in-production"
```

### AI Service Environment Variables

The AI service has its own `.env` file located at `ai-service/.env.example`. The table below lists
the most relevant options. See [ai-service/README.md](ai-service/README.md) for full configuration
details.

| Variable                 | Purpose                                      |
| ------------------------ | -------------------------------------------- |
| `DEBUG`                  | Enable debug mode for verbose logging        |
| `ENVIRONMENT`            | Name of the runtime environment              |
| `HOST`                   | Interface to bind the FastAPI server         |
| `PORT`                   | Port for the HTTP service                    |
| `WORKERS`                | Number of worker processes                   |
| `REDIS_URL`              | Connection string for Redis caching          |
| `REDIS_TIMEOUT`          | How long to wait for Redis responses         |
| `REDIS_RETRY_ON_TIMEOUT` | Whether to retry when Redis times out        |
| `CACHE_TTL`              | Default cache time-to-live in seconds        |
| `CACHE_PREFIX`           | Prefix for all cache keys                    |
| `ML_MODEL_PATH`          | Directory containing trained models          |
| `ENABLE_ML_FEATURES`     | Toggle machine learning analysis             |
| `MAX_ANALYSIS_DAYS`      | How far back to analyze entries              |
| `MIN_ENTRIES_FOR_ML`     | Minimum records required before ML runs      |
| `RATE_LIMIT_REQUESTS`    | Number of requests allowed per window        |
| `RATE_LIMIT_WINDOW`      | Duration of the rate limit window in seconds |
| `LOG_LEVEL`              | Logging verbosity level                      |
| `LOG_FORMAT`             | Format of log output                         |
| `BACKEND_URL`            | Base URL of the backend API                  |
| `BACKEND_TIMEOUT`        | Timeout for backend API calls                |

### Production Build

```bash
# Build all projects
pnpm build

# Or build individually
pnpm build:frontend
pnpm build:backend

# AI service
uv run uvicorn ai_service.main:app
```

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0** (AGPLv3).

See the [LICENSE](LICENSE) file for full license details.

### What this means?

- ✅ Free to use, modify, and distribute
- ✅ Commercial use allowed
- ⚠️ Must share source code if you distribute
- ⚠️ Network use triggers copyleft (AGPLv3 special sauce)
- 📋 Must include license and copyright notices

**TL;DR**: Keep it open source, keep it free. Don't be a corporate asshole.

______________________________________________________________________

**Disclaimer:** Not responsible for phone screen damage caused by ill-advised photo documentation.
Use with pride, shame, or scientific detachment. Up to you.

______________________________________________________________________

**Golden Rule:** If your code stinks, it won't get merged. And yes, we'll know. 💩

Ready to track some legendary logs? Get started above!

<!--

  markdownlint-configure-file {
    "no-inline-html": false,
    "no-alt-text": false,
  }

-->
