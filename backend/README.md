# Poo Tracker Backend (Go)

This Go service exposes a REST API for the Poo Tracker app. It now replicates the
core functionality of the old Node.js backend using the Gin framework.

## Development

```bash
# Run the server
go run ./backend

# Run tests
go test ./...
```

### Endpoints

- `GET /health` – basic health check
- `GET /api/bowel-movements` – list entries
- `POST /api/bowel-movements` – create entry
- `GET /api/bowel-movements/:id` – get entry
- `PUT /api/bowel-movements/:id` – update entry
- `DELETE /api/bowel-movements/:id` – delete entry
- `GET /api/analytics` – summary statistics
