# Poo Tracker Backend (Go)

This service implements the API for Poo Tracker using Go and the Gin framework. The previous Node.js backend has been fully rewritten in Go with dependency injection and in-memory repositories for easy testing.

## Development

```bash
# Run the server
go run ./backend

# Run tests
go test ./...
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
- JWT authentication with user registration and login
- Comprehensive validation for all endpoints
- Full CRUD operations for bowel movements and meals
- Analytics service with pluggable strategies
- Comprehensive test coverage
- RESTful API design

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
