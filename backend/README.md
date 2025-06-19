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
- `GET /api/bowel-movements/:id` – get entry
- `PUT /api/bowel-movements/:id` – update entry
- `DELETE /api/bowel-movements/:id` – delete entry
- `GET /api/meals` – list meals
- `POST /api/meals` – create meal
- `GET /api/meals/:id` – get meal
- `PUT /api/meals/:id` – update meal
- `DELETE /api/meals/:id` – delete meal
- `GET /api/analytics` – summary statistics
