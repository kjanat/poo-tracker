# Poo Tracker Backend

A comprehensive health tracking application backend built in Go, following Domain-Driven Design (DDD) principles with clean architecture patterns.

## Project Overview

The Poo Tracker backend provides a REST API for comprehensive digestive health monitoring, including bowel movement tracking, meal logging, symptom monitoring, medication management, and advanced health analytics.

## Architecture

The project follows a clean architecture with clear separation of concerns:

```
backend/
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── domain/           # Business logic and domain models
│   ├── infrastructure/   # External dependencies and implementations
│   └── app/             # Application services and orchestration
├── server/               # HTTP server and API handlers
├── docs/                # API documentation
└── data/                # Database files (development)
```

## Domain Structure

### Core Domains

#### 1. Bowel Movement Domain (`internal/domain/bowelmovement/`)

**Purpose**: Tracks bowel movements with comprehensive health metrics

**Key Models**:

- `BowelMovement`: Core tracking model with Bristol Stool Scale (1-7)
- `BowelMovementDetails`: Extended metadata and descriptions
- Pain, strain, satisfaction metrics (1-10 scale)
- Duration tracking and urgency levels
- Photo attachments and detailed notes

**Service Operations**:

- Create, read, update, delete bowel movement entries
- Date range queries for analytics
- User-specific data with privacy controls
- Statistical aggregations

#### 2. Meal Domain (`internal/domain/meal/`)

**Purpose**: Comprehensive meal and nutrition tracking

**Key Models**:

- `Meal`: Nutritional tracking with calorie counting
- Ingredient lists and dietary restrictions
- Meal timing and portion sizes
- Fiber content and dietary trigger tracking
- Photo documentation

**Service Operations**:

- Meal logging with nutritional analysis
- Dietary pattern analysis
- Trigger food identification
- Nutrition goal tracking

#### 3. Symptom Domain (`internal/domain/symptom/`)

**Purpose**: Symptom tracking and pattern recognition

**Key Models**:

- `Symptom`: Comprehensive symptom tracking
- Severity scales (1-10) and duration tracking
- Body part localization and symptom categorization
- Trigger identification and correlation analysis

**Service Operations**:

- Symptom logging with severity tracking
- Pattern recognition and trigger analysis
- Correlation with meals and medications
- Trend analysis over time

#### 4. Medication Domain (`internal/domain/medication/`)

**Purpose**: Medication management and adherence tracking

**Key Models**:

- `Medication`: Comprehensive medication tracking
- Dosage, frequency, and administration tracking
- Side effect monitoring and effectiveness assessment
- Start/end dates and PRN (as-needed) medications

**Service Operations**:

- Medication schedule management
- Adherence tracking and reporting
- Side effect correlation analysis
- Effectiveness assessment

#### 5. Analytics Domain (`internal/domain/analytics/`)

**Purpose**: Advanced health analytics and insights

**Key Features**:

- Cross-domain health overview generation
- Correlation analysis between health factors
- Trend analysis and pattern recognition
- Personalized health scoring (0-100)
- Evidence-based recommendations
- Risk factor identification

**Service Operations**:

- `GetUserHealthOverview()`: Comprehensive health summary
- `GetCorrelationAnalysis()`: Factor relationship analysis
- `GetTrendAnalysis()`: Time-series health trends
- `GetBehaviorPatterns()`: Behavioral pattern recognition
- `GetHealthInsights()`: Actionable health insights
- `GetHealthScore()`: Overall health scoring
- `GetRecommendations()`: Personalized recommendations

## Development

```bash
# Run the server
go run ./cmd/server/main.go

# Run tests
go test ./...

# Build for production
go build -o bin/server ./cmd/server/main.go
```

### Quick Start

```bash
# Install dependencies
go mod download

# Set up environment
cp .env.example .env
# Edit .env with your configuration

# Start development server
go run cmd/server/main.go
```

### Architecture

- `internal/model` – domain models
- `internal/repository` – repository interfaces and implementations
- `internal/service` – business logic with pluggable analytics strategies
- `server` – HTTP handlers and routing

The `main.go` file wires dependencies using constructor functions. A memory repository is used by default but can be swapped out for a real database implementation.

### Endpoints

- `GET /health` – basic health check
- `GET /api/bowel-movements` – list entries
- `POST /api/bowel-movements` – create entry
  - Request body: `{"userId": "string", "bristolType": 1-7, "notes": "optional"}`
- `GET /api/bowel-movements/:id` – get entry
- `PUT /api/bowel-movements/:id` – update entry
  - Request body: `{"bristolType": 1-7, "notes": "optional"}` (partial updates supported)
- `DELETE /api/bowel-movements/:id` – delete entry
- `GET /api/meals` – list meals
- `POST /api/meals` – create meal
  - Request body: `{"userId": "string", "name": "string", "calories": number}`
- `GET /api/meals/:id` – get meal
- `PUT /api/meals/:id` – update meal
  - Request body: `{"name": "string", "calories": number}` (partial updates supported)
- `DELETE /api/meals/:id` – delete meal
- `GET /api/analytics` – summary statistics
  - Response: Analytics data based on configured strategy

#### BowelMovement Details (Enhanced tracking)

- `POST /api/bowel-movements/:id/details` – create detailed information for bowel movement
  - Request body: `{"notes": "string", "detailedNotes": "string", "environment": "string", "preConditions": "string", "postConditions": "string", "aiRecommendations": "string", "tags": ["string"], "weatherCondition": "string", "stressLevel": 1-10, "sleepQuality": 1-10, "exerciseIntensity": 1-10}`
- `GET /api/bowel-movements/:id/details` – get detailed information
- `PUT /api/bowel-movements/:id/details` – update detailed information
- `DELETE /api/bowel-movements/:id/details` – delete detailed information

#### User Management

- `POST /api/register` – create user account
  - Request body: `{"email": "string", "password": "string", "name": "string"}`
- `POST /api/login` – authenticate user
  - Request body: `{"email": "string", "password": "string"}`
  - Response: User data with JWT token
- `GET /api/profile` – get authenticated user profile (requires auth header)

## Current Implementation Status

### ✅ Completed Features

- Clean architecture with dependency injection
- In-memory repositories for bowel movements, meals, and users
- **Enhanced BowelMovement model with separated details for performance**
- **BowelMovementDetails with comprehensive tracking fields and AI analysis**
- JWT authentication with user registration and login
- Comprehensive validation for all endpoints
- Full CRUD operations for bowel movements, meals, and details
- Analytics service with pluggable strategies
- Comprehensive test coverage
- RESTful API design
- **Automatic HasDetails flag synchronization between models**

### 🔄 In Progress / Planned

- **Database**: Migration from in-memory to PostgreSQL
- **File Storage**: Photo upload integration with MinIO/S3
- **Advanced Models**: Symptoms, medications, and their relationships
- **Enhanced Security**: Rate limiting, 2FA, password reset
- **Data Export**: PDF reports and data export functionality
- **Advanced Analytics**: Pattern detection and health insights

### 🏗️ Architecture Notes

- Uses Go's built-in dependency injection via constructor functions
- Strategy pattern for analytics (easily extensible)
- Middleware-based authentication using JWT
- Memory repositories can be swapped for PostgreSQL implementations
- Clean separation: handlers → services → repositories
