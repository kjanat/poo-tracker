# 💩 Poo Tracker

<img align="right" width="150" height="150" src="branding/logo.svg">

Ever wondered if your gut's on a winning streak, or if your last kebab is about to take its revenge? Welcome to **Poo Tracker**, the brutally honest app that lets you log every majestic turd, rabbit pellet, and volcanic diarrhea eruption without shame or censorship. Because your bowel movements _do_ matter, and no, your GP isn't getting this level of data.

## 🚀 Quick Start

### Prerequisites

- [Go](https://golang.org/) 1.22+ (backend, uses Gin + GORM)
- [Node.js](https://nodejs.org/) 22+ (frontend only)
- [Docker](https://docs.docker.com/engine/) & [Docker Compose](https://docs.docker.com/compose/)
- [pnpm](https://pnpm.io/) 9+ (because npm is for amateurs)
- [uv](https://docs.astral.sh/uv/) for Python (because pip is for the weak)

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/kjanat/poo-tracker.git
   cd poo-tracker
   ```

2. **Install dependencies:**

   ```bash
   # Install Node.js dependencies (frontend only)
   cd frontend && pnpm install

   # Install Python dependencies for AI service
   cd ai-service && uv sync && cd ..
   ```

3. **Set up environment variables:**

   ```bash
   cp .env.example .env
   cp frontend/.env.example frontend/.env.local
   cp ai-service/.env.example ai-service/.env
   # Edit each file with your local values
   ```

4. **Start the services:**

   ```bash
   # Start database and supporting services
   make docker-up

   # Start all development servers (frontend + Go backend + AI)
   make dev
   ```

5. **Set up the database:**

   Start the PostgreSQL database using Docker Compose:

   ```bash
   make docker-up
   ```

   A sample set of credentials for API testing is provided in
   [`login.example.json`](login.example.json). Copy this file to
   `login.json` if you need a local login file that's ignored by Git.

6. **Open your browser:**

- Frontend: <http://localhost:5173>
- Backend API: <http://localhost:3002>
- AI Service: <http://localhost:8001>
- MinIO Console: <http://localhost:9002> (minioadmin/minioadmin123)

## 🏗️ Tech Stack

- **Frontend**: React + Vite + TypeScript + TailwindCSS v4
- **Backend**: Go + Gin + GORM
- **Database**: PostgreSQL (or SQLite for dev)
- **Storage**: MinIO (S3-compatible for photos)
- **AI Service**: Python + FastAPI + scikit-learn
- **Infrastructure**: Docker + Docker Compose
- **Package Management**: pnpm (Node.js) + uv (Python)

## 📝 Environment Variables

Each package ships with an `.env.example` file. Copy them and tweak the values before running anything:

```bash
cp .env.example .env
cp frontend/.env.example frontend/.env.local
cp ai-service/.env.example ai-service/.env
```

These example files contain **sample credentials only**. Replace them with your real secrets and **never commit private keys or passwords** to the repository.

## 📖 What's the fucking point?

- **Track your shits**: Log every glorious bowel movement in a timeline worthy of its own Netflix docu-series. Date, time, volume, Bristol score, color, floaters—you name it.
- **Photo evidence**: Snap a pic for science (or just to traumatize your friends and family). Upload and catalog the full visual glory of your rectal output.
- **Meal log**: Record what you've stuffed down your throat so you can finally prove that Taco Tuesday was, in fact, a war crime against your colon.
- **AI-powered analysis**: Our cold, heartless machine learning model doesn't judge—just coldly correlates your food, stress, and fiber intake with your latest gut grenade.
- **Export & share**: PDF reports for your doctor, gastroenterologist, or weird fetish community. It's your data, so overshare as you wish.

## 🔧 Development

### Available Scripts

```bash
# Development - All Services
make dev              # Start frontend, Go backend and AI service
make dev-frontend     # Start frontend only
make dev-backend      # Start Go backend only
make dev-ai           # Start AI service only

# Building
make build            # Build all projects
make build-frontend   # Build frontend only
make build-backend    # Build Go backend only

# Database Operations
# (Database migrations handled automatically by GORM in Go backend)

# Testing & Quality
make test             # Run all tests (frontend + backend)
make test-frontend    # Run frontend tests only
make test-backend     # Run backend tests only
make lint             # Run linters on all projects
make lint-fix         # Auto-fix linting issues
make clean            # Clean all build artifacts

# Code Formatting
make format           # Format backend and AI service

# For watch mode in frontend tests, use:
#   make test-frontend WATCH=1
# (or add a Makefile target if needed)
```

> **Note:** Always use Makefile targets for all workflows (test, lint, build, dev). Do not use pnpm workspace commands directly.

### Project Structure

```text
poo-tracker/
├── frontend/           # React frontend (Vite + TypeScript + TailwindCSS)
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── backend/            # Go API (Gin + GORM)
│   ├── *.go
│   ├── go.mod
│   └── README.md
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
cd frontend && pnpm add react-router-dom

# Run commands on specific workspaces
cd frontend && pnpm run build

# Run commands on all workspaces
pnpm --recursive run build
pnpm --parallel run dev
```

### API Endpoints

| Endpoint                   | Method   | Description                       |
| -------------------------- | -------- | --------------------------------- |
| `/api/bowel-movements`     | `GET`    | List bowel movement entries       |
| `/api/bowel-movements`     | `POST`   | Create a new bowel movement entry |
| `/api/bowel-movements/:id` | `GET`    | Get a specific entry              |
| `/api/bowel-movements/:id` | `PUT`    | Update a bowel movement entry     |
| `/api/bowel-movements/:id` | `DELETE` | Delete a bowel movement entry     |
| `/api/meals`               | `GET`    | List meal entries                 |
| `/api/meals`               | `POST`   | Create a new meal entry           |
| `/api/meals/:id`           | `GET`    | Get a specific meal               |
| `/api/meals/:id`           | `PUT`    | Update a meal entry               |
| `/api/meals/:id`           | `DELETE` | Delete a meal entry               |
| `/api/analytics`           | `GET`    | Get simple analytics              |

## 📝 How it works

1. **Eat** something questionable.
2. **Shit**. (Preferably in a toilet, but we're not here to kink-shame.)
3. **Log** your experience in Poo Tracker—optionally add a photo, rate your suffering, describe the smell if you're brave enough.
4. **Repeat** until you finally realize you're lactose intolerant or have IBS.

## 🔒 Privacy, baby

We encrypt your brown notes and hide them away. Nobody's reading your logs except you—and whatever godforsaken AI wants to learn about the day-after effects of your sushi buffet.

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
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Follow the coding standards (see [AGENTS.md](AGENTS.md))
4. Write tests for new features
5. Commit using conventional commits (`feat: add Bristol chart selector`)
6. Push and create a Pull Request

### Coding Standards

- TypeScript everywhere (frontend)
- Go (Gin + GORM) for backend
- Component-based architecture
- TailwindCSS for styling (no CSS modules)
- RESTful API design
- Comprehensive test coverage
- ESLint with @typescript-eslint and Prettier (frontend)
- Use pnpm workspace commands or Makefile targets for consistent development (never use `cd` in scripts)
- All code quality checks are managed by pre-commit hooks

## 🚀 Deployment

### Environment Variables

```env
# Database
DATABASE_URL="postgresql://poo_user:secure_password_123@localhost:5432/poo_tracker"
# For SQLite (development only):
# SQLITE_DSN="backend/data/poo-tracker.db"
# DB_TYPE="sqlite" # or "postgres"

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

The AI service has its own `.env` file located at `ai-service/.env.example`. The
table below lists the most relevant options. See
[ai-service/README.md](ai-service/README.md) for full configuration details.

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
make build

# Or build individually
make build-frontend
make build-backend    # Build Go backend only

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

---

**Disclaimer:** Not responsible for phone screen damage caused by ill-advised photo documentation. Use with pride, shame, or scientific detachment. Up to you.

---

**Golden Rule:** If your code stinks, it won't get merged. And yes, we'll know. 💩

Ready to track some legendary logs? Get started above!

<!--

  markdownlint-configure-file {
    "no-inline-html": false,
    "no-alt-text": false,
  }

-->
