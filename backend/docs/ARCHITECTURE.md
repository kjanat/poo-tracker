# Poo Tracker Backend Architecture

## 🏗️ Clean Architecture Overview

This backend follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                         │
│  ┌─────────────────┐ ┌─────────────────┐ ┌───────────────┐  │
│  │    Handlers     │ │   Middleware    │ │     DTOs      │  │
│  └─────────────────┘ └─────────────────┘ └───────────────┘  │
└─────────────────────�────┬─────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Service Layer                             │
│  ┌─────────────────┐ ┌─────────────────┐ ┌───────────────┐  │
│  │  Domain Logic   │ │  Business Rules │ │ Orchestration │  │
│  └─────────────────┘ └─────────────────┘ └───────────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Domain Layer                              │
│  ┌─────────────────┐ ┌─────────────────┐ ┌───────────────┐  │
│  │     Models      │ │   Interfaces    │ │    Errors     │  │
│  └─────────────────┘ └─────────────────┘ └───────────────┘  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│               Infrastructure Layer                          │
│  ┌─────────────────┐ ┌─────────────────┐ ┌───────────────┐  │
│  │  Repositories   │ │    Database     │ │   External    │  │
│  │   (Memory/DB)   │ │   Connections   │ │   Services    │  │
│  └─────────────────┘ └─────────────────┘ └───────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 📂 Project Structure

```
backend/
├── cmd/server/                     # Application entry points
│   └── main.go                     # Main application bootstrap
├── internal/                       # Private application code
│   ├── app/                        # Application configuration & DI
│   │   ├── config.go               # Configuration management
│   │   └── container.go            # Dependency injection container
│   ├── domain/                     # Business domain (Core Layer)
│   │   ├── user/                   # User domain
│   │   │   ├── model.go            # User entities & value objects
│   │   │   ├── repository.go       # Repository interface
│   │   │   ├── service.go          # Service interface
│   │   │   └── errors.go           # Domain-specific errors
│   │   ├── bowelmovement/          # Bowel movement tracking
│   │   ├── meal/                   # Meal tracking
│   │   ├── symptom/                # Symptom tracking
│   │   ├── medication/             # Medication tracking
│   │   ├── analytics/              # Analytics & insights
│   │   └── shared/                 # Shared domain concepts
│   └── infrastructure/             # External concerns (Infrastructure Layer)
│       ├── repository/             # Data access implementations
│       │   ├── memory/             # In-memory implementations (dev/test)
│       │   └── postgres/           # PostgreSQL implementations (prod)
│       ├── service/                # Service implementations
│       ├── http/                   # HTTP transport layer
│       │   ├── handlers/           # HTTP handlers
│       │   ├── middleware/         # HTTP middleware
│       │   └── dto/                # Data transfer objects
│       └── database/               # Database setup & migrations
├── migrations/                     # Database migrations
├── docs/                          # Documentation
└── scripts/                       # Build & deployment scripts
```

## 🔧 Dependency Injection

We use **explicit constructor injection** for clean, testable code:

### Container Setup

```go
// internal/app/container.go
type Container struct {
    // Repositories
    UserRepository user.Repository
    BowelMovementRepository bowelmovement.Repository
    // ... other repos

    // Services
    UserService user.Service
    AnalyticsService analytics.Service
    // ... other services
}

func NewContainer() (*Container, error) {
    container := &Container{}

    // Wire up repositories
    container.UserRepository = memory.NewUserRepository()

    // Wire up services with injected dependencies
    container.UserService = service.NewUserService(container.UserRepository)

    return container, nil
}
```

### Service Construction

```go
// internal/infrastructure/service/user_service.go
type UserService struct {
    repo user.Repository
}

func NewUserService(repo user.Repository) user.Service {
    return &UserService{repo: repo}
}
```

## 🧪 Testing Strategy

### 1. **Domain Unit Tests**

- Test business rules in isolation
- No external dependencies
- Fast and reliable

### 2. **Service Unit Tests**

- Mock repository interfaces
- Test service orchestration
- Verify error handling

```go
func TestUserService_Register(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

    // Test the service
    result, err := service.Register(ctx, input)

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### 3. **Repository Integration Tests**

- Test against real in-memory implementations
- Verify data persistence behavior
- Test repository interface compliance

### 4. **HTTP Handler Tests** (Future)

- Test HTTP layer in isolation
- Mock service dependencies
- Verify request/response handling

## 🔌 Interface-Driven Development

Every layer depends on **interfaces**, not concrete implementations:

```go
// Domain defines the contract
type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

// Infrastructure implements the contract
type memoryRepository struct {
    users map[string]*User
}

func (r *memoryRepository) Create(ctx context.Context, user *User) error {
    // Implementation details
}
```

## 🏃‍♂️ Development Workflow

### 1. **Start Development Server**

```bash
go run cmd/server/main.go
```

### 2. **Run Tests**

```bash
# All tests
go test ./...

# Specific package
go test ./internal/infrastructure/service -v

# With coverage
go test -cover ./...
```

### 3. **Linting & Formatting**

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run
```

### 4. **Database Migrations** (Future)

```bash
# Run migrations
go run cmd/migrate/main.go up

# Rollback migrations
go run cmd/migrate/main.go down
```

## 📊 Package Dependencies

### Domain Layer (Core)

- ✅ **Zero external dependencies**
- ✅ **Pure business logic**
- ✅ **Framework agnostic**

### Service Layer

- ➡️ Depends on: Domain interfaces only
- ✅ **No database dependencies**
- ✅ **No HTTP dependencies**

### Infrastructure Layer

- ➡️ Depends on: Domain interfaces + external libraries
- 📦 **GORM** (database ORM)
- 📦 **Gin** (HTTP framework)
- 📦 **bcrypt** (password hashing)

## 🚀 Deployment

### Environment Configuration

```bash
# Development
export APP_ENV=development
export DB_HOST=localhost
export DB_PORT=5432

# Production
export APP_ENV=production
export DB_HOST=prod-db-host
export JWT_SECRET=your-production-secret
```

### Build for Production

```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

## 🔧 Configuration Management

Configuration is centralized in `internal/app/config.go`:

```go
type Config struct {
    Environment string
    Server ServerConfig
    Database DatabaseConfig
    JWT JWTConfig
}

func LoadConfig() *Config {
    return &Config{
        Environment: getEnvOrDefault("APP_ENV", "development"),
        Server: ServerConfig{
            Port: getEnvOrDefault("SERVER_PORT", "8080"),
        },
        // ... other config
    }
}
```

## 📈 Performance Considerations

- **Connection Pooling**: Database connections are pooled and reused
- **Context Propagation**: All operations support context cancellation
- **Memory Efficiency**: In-memory repositories use maps for O(1) lookups
- **Lazy Loading**: Services are initialized only when needed

## 🛡️ Security Features

- **Password Hashing**: bcrypt with configurable cost
- **JWT Authentication**: Stateless token-based auth
- **Input Validation**: All inputs validated at service layer
- **SQL Injection Prevention**: GORM handles parameterized queries
- **CORS Protection**: Configurable CORS middleware

## 🧩 Extending the Architecture

### Adding a New Domain

1. **Create domain package**: `internal/domain/newdomain/`
2. **Define models**: `model.go`, `repository.go`, `service.go`
3. **Implement repository**: `internal/infrastructure/repository/memory/`
4. **Implement service**: `internal/infrastructure/service/`
5. **Wire in container**: Add to `internal/app/container.go`
6. **Add tests**: Repository + service test files

### Adding New Endpoints

1. **Create DTOs**: `internal/infrastructure/http/dto/`
2. **Create handlers**: `internal/infrastructure/http/handlers/`
3. **Wire routes**: Update route registration
4. **Add middleware**: Authentication, validation, etc.

This architecture provides a solid foundation that's **maintainable**, **testable**, and **scalable**! 🚀
