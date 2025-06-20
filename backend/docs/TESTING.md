# Testing Guide

## 🧪 Testing Philosophy

Our testing strategy follows the **Testing Pyramid**:

```
           ∆
          / \
         /   \
        /     \
       /   UI  \     ← Few, Slow, Expensive
      /_________\
     /           \
    / Integration \  ← Some, Medium Speed
   /_______________\
  /                 \
 /    Unit Tests     \  ← Many, Fast, Cheap
/_____________________\
```

## 📊 Test Categories

### 1. Unit Tests

**What**: Test individual functions/methods in isolation  
**Speed**: Fast (< 1ms per test)  
**Coverage**: 80%+ of business logic

```go
func TestUserService_Register(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    input := &user.RegisterInput{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "password123",
    }

    mockRepo.On("EmailExists", mock.Anything, input.Email).Return(false, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)

    // Act
    result, err := service.Register(context.Background(), input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, input.Email, result.Email)
    mockRepo.AssertExpectations(t)
}
```

### 2. Integration Tests

**What**: Test component interactions (service + repository)  
**Speed**: Medium (1-10ms per test)  
**Coverage**: Critical user flows

```go
func TestUserRepository_CRUD(t *testing.T) {
    // Use real repository implementation
    repo := memory.NewUserRepository()
    ctx := context.Background()

    // Create user
    user := &user.User{
        ID:       "test-id",
        Email:    "test@example.com",
        Username: "testuser",
        Name:     "Test User",
    }

    err := repo.Create(ctx, user)
    assert.NoError(t, err)

    // Retrieve user
    retrieved, err := repo.GetByID(ctx, user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, retrieved.Email)
}
```

### 3. End-to-End Tests (Future)

**What**: Test complete user journeys via HTTP  
**Speed**: Slow (100ms+ per test)  
**Coverage**: Happy path scenarios

## 🛠️ Testing Tools

### Core Testing

- **Go Standard Library**: `testing` package
- **Testify**: Assertions and mocking (`github.com/stretchr/testify`)

### Mocking Strategy

```go
// mockgen generates mocks from interfaces
//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}
```

## 📁 Test Organization

```
internal/
├── domain/
│   └── user/
│       ├── model.go
│       ├── repository.go
│       └── service.go
└── infrastructure/
    ├── repository/
    │   └── memory/
    │       ├── user_repository.go
    │       └── user_repository_test.go  ← Integration tests
    └── service/
        ├── user_service.go
        └── user_service_test.go         ← Unit tests with mocks
```

## 🎯 Test Naming Conventions

### Function Naming

```go
func TestServiceName_MethodName_Scenario(t *testing.T)

// Examples:
func TestUserService_Register_Success(t *testing.T)
func TestUserService_Register_EmailExists(t *testing.T)
func TestUserRepository_GetByID_NotFound(t *testing.T)
```

### File Naming

- `*_test.go` - Tests in same package
- `*_integration_test.go` - Integration tests
- `testdata/` - Test fixtures and data

## 🧩 Testing Patterns

### Table-Driven Tests

```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   *RegisterInput
        wantErr bool
        errType error
    }{
        {
            name: "valid input",
            input: &RegisterInput{
                Email:    "test@example.com",
                Username: "testuser",
                Password: "password123",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            input: &RegisterInput{
                Email:    "invalid-email",
                Username: "testuser",
                Password: "password123",
            },
            wantErr: true,
            errType: user.ErrInvalidEmail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateRegisterInput(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.True(t, errors.Is(err, tt.errType))
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Test Fixtures

```go
// testdata/users.go
package testdata

func ValidUser() *user.User {
    return &user.User{
        ID:       "test-user-id",
        Email:    "test@example.com",
        Username: "testuser",
        Name:     "Test User",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func ValidRegisterInput() *user.RegisterInput {
    return &user.RegisterInput{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "password123",
        Name:     "Test User",
    }
}
```

## 🏃‍♂️ Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/infrastructure/service

# Run specific test
go test -run TestUserService_Register ./internal/infrastructure/service

# Verbose output
go test -v ./...

# Run tests in parallel
go test -parallel 4 ./...
```

### Test Coverage Goals

- **Domain Logic**: 90%+ coverage
- **Service Layer**: 85%+ coverage
- **Repository Layer**: 80%+ coverage
- **Overall**: 80%+ coverage

### Continuous Integration

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - run: go test -race -coverprofile=coverage.out ./...
      - run: go tool cover -func=coverage.out
```

## 🐛 Testing Best Practices

### DO ✅

- **Test behavior, not implementation**
- **Use descriptive test names**
- **Follow AAA pattern** (Arrange, Act, Assert)
- **Test error cases**
- **Use table-driven tests for multiple scenarios**
- **Mock external dependencies**
- **Test at the interface boundary**

### DON'T ❌

- **Test private functions directly**
- **Mock everything (test some real integrations)**
- **Write tests after code (TDD preferred)**
- **Ignore flaky tests**
- **Test framework code (e.g., Gin handlers)**
- **Share state between tests**

### Error Testing

```go
func TestUserService_Register_EmailExists(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    input := &user.RegisterInput{
        Email: "existing@example.com",
    }

    // Mock repository to return "email exists" error
    mockRepo.On("EmailExists", mock.Anything, input.Email).Return(true, nil)

    result, err := service.Register(context.Background(), input)

    assert.Error(t, err)
    assert.Nil(t, result)
    assert.True(t, errors.Is(err, user.ErrEmailAlreadyExists))
}
```

## 📈 Test Metrics

### Performance Benchmarks

```go
func BenchmarkUserService_Register(b *testing.B) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    input := &user.RegisterInput{
        Email:    "test@example.com",
        Username: "testuser",
        Password: "password123",
    }

    mockRepo.On("EmailExists", mock.Anything, mock.Anything).Return(false, nil)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.Register(context.Background(), input)
    }
}
```

### Memory Profiling

```bash
# Run with memory profiling
go test -memprofile=mem.prof ./...

# Analyze memory usage
go tool pprof mem.prof
```

## 🔄 Test-Driven Development

### Red-Green-Refactor Cycle

1. **Red**: Write failing test
2. **Green**: Write minimal code to pass
3. **Refactor**: Improve code while keeping tests green

```go
// 1. RED: Write failing test first
func TestUserService_ChangePassword(t *testing.T) {
    service := NewUserService(nil)
    err := service.ChangePassword(context.Background(), "user-id", &ChangePasswordInput{
        OldPassword: "old",
        NewPassword: "new",
    })
    assert.NoError(t, err)
}

// 2. GREEN: Implement minimal functionality
func (s *UserService) ChangePassword(ctx context.Context, userID string, input *ChangePasswordInput) error {
    return nil // Minimal implementation
}

// 3. REFACTOR: Add real implementation
func (s *UserService) ChangePassword(ctx context.Context, userID string, input *ChangePasswordInput) error {
    // Real implementation with validation, hashing, etc.
}
```

This testing strategy ensures our codebase remains **reliable**, **maintainable**, and **regression-free**! 🚀
