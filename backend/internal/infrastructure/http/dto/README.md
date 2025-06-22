# DTO Architecture Documentation

## Overview

The Data Transfer Object (DTO) layer provides a clean separation between the HTTP layer and the domain layer, following clean architecture principles. This layer handles request/response transformation, validation, and ensures that domain models are never exposed directly through the API.

## Structure

```graphql
internal/infrastructure/http/dto/
├── bowelmovement/dto.go  # Bowel movement tracking DTOs
├── user/dto.go           # User management and authentication DTOs
├── meal/dto.go           # Meal logging and tracking DTOs
├── symptom/dto.go        # Symptom reporting DTOs
├── medication/dto.go     # Medication tracking DTOs
├── analytics/dto.go      # Analytics and insights DTOs
└── shared/dto.go         # Common DTOs and error handling
```

## Key Features

### 1. Request/Response Separation

Each domain has separate DTOs for:

- **Create Requests**: Required fields for entity creation
- **Update Requests**: Optional pointer fields for partial updates
- **Response DTOs**: Consistent output format with computed fields
- **List Responses**: Paginated collections with metadata

### 2. Comprehensive Validation

All DTOs include extensive validation using Gin binding tags:

- **Required fields**: `binding:"required"`
- **Length constraints**: `binding:"min=1,max=100"`
- **Format validation**: `binding:"email"`, `binding:"url"`
- **Enum validation**: `binding:"oneof=value1 value2"`
- **Cross-field validation**: Custom validation methods

### 3. Domain Mapping

Clean conversion functions between DTOs and domain models:

- `ToDomain*()` methods convert requests to domain models
- `ToResponse()` functions convert domain models to responses
- `ApplyToDomain*()` methods handle partial updates
- Proper handling of optional fields and null values

### 4. Error Handling

Standardized error responses with:

- Consistent error structure across all endpoints
- Field-specific validation errors
- Request ID tracking for debugging
- Actionable error messages and suggestions

## Domain-Specific DTOs

### Bowel Movement DTOs

- **Comprehensive tracking**: Bristol scale, consistency, urgency, timing
- **Photo support**: Image upload for visual tracking
- **Symptom correlation**: Links to related symptoms
- **Validation**: Bristol scale range, positive urgency values

### User DTOs

- **Authentication**: Login, registration, password changes
- **Profile management**: Personal information, preferences
- **Security**: Password hashing, JWT token handling
- **Validation**: Email format, password strength

### Meal DTOs

- **Meal tracking**: Name, timing, category, nutritional info
- **Dietary information**: Spicy level, fiber content, allergens
- **Photo support**: Meal photography for reference
- **Validation**: Spicy level range, positive calorie values

### Symptom DTOs

- **Comprehensive tracking**: Type, severity, duration, triggers
- **Body mapping**: Specific body part localization
- **Time tracking**: Start time, duration, frequency
- **Validation**: Severity scale (1-10), logical time ranges

### Medication DTOs

- **Medication management**: Name, dosage, frequency, timing
- **Clinical details**: Generic name, brand, category, route
- **Side effects**: Tracking of adverse reactions
- **Validation**: Dosage format, date ranges, frequency patterns

### Analytics DTOs

- **Statistics**: Bowel movement frequency, consistency trends
- **Correlations**: Meal-symptom-movement relationships
- **Insights**: Health scores, recommendations, patterns
- **Time ranges**: Flexible date filtering for analysis

### Shared DTOs

- **Error responses**: Standardized error format
- **Pagination**: Common pagination metadata
- **Success responses**: Consistent success format
- **Health checks**: System status responses

## Validation Strategy

### Field-Level Validation

Uses Gin's binding tags for immediate validation:

```go
type CreateRequest struct {
    Name  string `json:"name" binding:"required,min=1,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   *int   `json:"age" binding:"omitempty,min=1,max=150"`
}
```

### Business Logic Validation

Custom validation methods for complex rules:

```go
func (r *CreateRequest) Validate() error {
    if r.EndDate != nil && r.EndDate.Before(r.StartDate) {
        return errors.New("end date must be after start date")
    }
    return nil
}
```

### Error Response Format

Consistent error structure with detailed information:

```go
type ErrorResponse struct {
    Error     ErrorDetail   `json:"error"`
    RequestID string        `json:"request_id"`
    Timestamp time.Time     `json:"timestamp"`
}

type ErrorDetail struct {
    Code        string       `json:"code"`
    Message     string       `json:"message"`
    Fields      []FieldError `json:"fields,omitempty"`
    Suggestions []string     `json:"suggestions,omitempty"`
}
```

## Usage Patterns

### Handler Implementation

```go
func (h *Handler) CreateBowelMovement(c *gin.Context) {
    var req bowelmovement.CreateBowelMovementRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, shared.NewValidationErrorResponse(err))
        return
    }

    if err := req.Validate(); err != nil {
        c.JSON(400, shared.NewErrorResponse("VALIDATION_ERROR", err.Error(), ""))
        return
    }

    domainBM := req.ToDomainBowelMovement(userID)
    result, err := h.service.Create(c.Request.Context(), domainBM)
    if err != nil {
        c.JSON(500, shared.NewErrorResponse("INTERNAL_ERROR", err.Error(), ""))
        return
    }

    response := bowelmovement.ToBowelMovementResponse(result)
    c.JSON(201, shared.NewSuccessResponse("Bowel movement created", response))
}
```

### Pagination Support

```go
func (h *Handler) ListBowelMovements(c *gin.Context) {
    var req shared.BaseListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, shared.NewValidationErrorResponse(err))
        return
    }

    page, pageSize := req.GetDefaultPagination()
    results, totalCount, err := h.service.List(c.Request.Context(), page, pageSize)
    if err != nil {
        c.JSON(500, shared.NewErrorResponse("INTERNAL_ERROR", err.Error(), ""))
        return
    }

    response := bowelmovement.ToBowelMovementListResponse(results, totalCount, page, pageSize)
    pagination := shared.NewPaginationMeta(page, pageSize, totalCount)
    c.JSON(200, shared.NewPaginatedResponse(response.BowelMovements, pagination))
}
```

## Benefits

1. **Clean Architecture**: Complete separation between HTTP and domain layers
2. **Type Safety**: Compile-time guarantees for API contracts
3. **Consistent APIs**: Standardized request/response patterns
4. **Better Testing**: Easy to mock and test individual layers
5. **API Documentation**: Clear DTOs make API documentation generation easier
6. **Validation**: Comprehensive input validation prevents invalid data
7. **Maintainability**: Changes to domain models don't affect API contracts

## Next Steps

With DTOs complete, the next phase involves:

1. **Service Layer Implementation**: Move business logic to domain services
2. **Handler Restructuring**: Update HTTP handlers to use DTOs
3. **Repository Implementation**: Complete GORM repository implementations
4. **Middleware Integration**: Add authentication, logging, and validation middleware
5. **API Documentation**: Generate OpenAPI/Swagger documentation from DTOs

## 🧹 Code Quality & Pre-commit Hooks

All linting, formatting, and type-checking is managed via [pre-commit](https://pre-commit.com) and the `.pre-commit-config.yaml` in the project root. Husky and lint-staged are no longer used.

**Setup:**

```bash
uv tool install pre-commit  # or pip install pre-commit
pre-commit install
```

Hooks run on every commit, or manually with:

```bash
pre-commit run --all-files
```
