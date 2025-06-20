# Backend Architecture Refactoring Progress

## Overview

Comprehensive restructuring of the poo-tracker backend from functional but messy architecture to production-grade clean architecture following Go best practices and dependency injection patterns.

## Current State Analysis

### ❌ Issues Identified (2025-06-20)

1. **Mixed Concerns**: Server package contains routing, handlers, and server setup
2. **Monolithic Dependencies**: All wiring crammed into main.go without clear separation
3. **Inconsistent Handler Patterns**: Different handler types with different interfaces
4. **Missing Service Layer**: Handlers calling repositories directly bypassing business logic
5. **No DTOs**: Using domain models directly in HTTP layer violating clean architecture
6. **Scattered Configuration**: JWT secrets, database config spread across files
7. **Inconsistent Error Handling**: Each handler has different error response patterns
8. **No Dependency Container**: Manual dependency wiring makes testing difficult
9. **Mixed Repository Concerns**: Business logic bleeding into repository layer
10. **Poor Testability**: Hard to mock dependencies due to tight coupling

### 📊 Current Structure Audit

```text
backend/
├── main.go (47 lines) - Monolithic dependency wiring
├── server/
│   ├── server.go (169 lines) - Mixed concerns: routing + handlers + setup
│   ├── api.go (remaining meal handlers)
│   ├── bowel_movements_api.go (451 lines) - Too large, mixed concerns
│   ├── user_api.go + user_api_test.go
│   ├── symptom_api.go + symptom_api_test.go
│   ├── medication_api.go + medication_api_test.go
│   ├── two_factor_api.go (153 lines)
│   ├── meal_bowel_relations_api.go
│   ├── meal_symptom_relations_api.go
│   └── relations_coordinator_api.go
├── internal/
│   ├── model/ (domain models mixed with HTTP concerns)
│   ├── repository/ (memory implementations, mixed business logic)
│   ├── service/ (incomplete service layer)
│   ├── validation/ (properly separated ✅)
│   └── middleware/ (properly separated ✅)
```

## Target Architecture

### 🎯 Clean Architecture Structure

```text
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Minimal main, delegates to app
├── internal/
│   ├── app/                        # Application layer
│   │   ├── app.go                  # App constructor & dependency wiring
│   │   ├── config.go               # Configuration management
│   │   └── container.go            # Dependency injection container
│   ├── domain/                     # Core business logic (no external deps)
│   │   ├── bowelmovement/
│   │   │   ├── model.go            # Domain models
│   │   │   ├── repository.go       # Repository interfaces
│   │   │   ├── service.go          # Business logic
│   │   │   └── errors.go           # Domain-specific errors
│   │   ├── user/
│   │   ├── meal/
│   │   ├── symptom/
│   │   ├── medication/
│   │   ├── analytics/
│   │   └── shared/                 # Shared domain concepts
│   │       ├── errors.go
│   │       ├── interfaces.go
│   │       └── events.go
│   ├── infrastructure/             # External concerns
│   │   ├── repository/
│   │   │   ├── memory/             # In-memory implementations
│   │   │   └── postgres/           # Future postgres implementations
│   │   ├── http/                   # HTTP transport layer
│   │   │   ├── handlers/           # HTTP handlers by domain
│   │   │   │   ├── bowelmovement/
│   │   │   │   ├── user/
│   │   │   │   ├── meal/
│   │   │   │   ├── analytics/
│   │   │   │   └── health/
│   │   │   ├── middleware/
│   │   │   ├── dto/                # Request/Response DTOs
│   │   │   ├── router.go           # Route registration
│   │   │   └── server.go           # HTTP server setup
│   │   └── security/               # Security implementations
│   └── shared/                     # Shared utilities
│       ├── validation/
│       ├── logger/
│       └── testing/
```

## Implementation Phases

### Phase 1 - Infrastructure Setup ⏳ IN PROGRESS

**Priority: High** | **Started:** 2025-06-20 | **Estimated:** 2-3 hours

#### Tasks:

- [ ] Create new directory structure
- [ ] Move existing files to appropriate locations
- [ ] Create placeholder files for new architecture
- [ ] Update import paths for moved files
- [ ] Ensure all tests still pass after file moves

#### Files to Create:

- [ ] `cmd/server/main.go` (minimal main)
- [ ] `internal/app/app.go` (application setup)
- [ ] `internal/app/config.go` (configuration management)
- [ ] `internal/app/container.go` (dependency injection)
- [ ] `internal/domain/*/` (domain package structure)
- [ ] `internal/infrastructure/http/` (HTTP layer structure)

#### Migration Strategy:

1. Create new structure alongside existing
2. Copy files to new locations with updated imports
3. Gradually migrate functionality
4. Remove old files once everything works

**Success Criteria:**

- [ ] New directory structure exists
- [ ] All existing functionality preserved
- [ ] All tests pass
- [ ] Build succeeds

---

### Phase 2 - Domain Layer Extraction ⏳ PENDING

**Priority: High** | **Dependencies:** Phase 1 ✅ | **Estimated:** 4-5 hours

#### Tasks:

- [ ] Extract domain models from current model package
- [ ] Define repository interfaces in domain layer
- [ ] Create domain service interfaces
- [ ] Define domain-specific errors
- [ ] Implement domain validation rules

#### Domains to Extract:

- [ ] `domain/bowelmovement/` - Core bowel movement business logic
- [ ] `domain/user/` - User management and authentication
- [ ] `domain/meal/` - Meal tracking and management
- [ ] `domain/symptom/` - Symptom tracking
- [ ] `domain/medication/` - Medication management
- [ ] `domain/analytics/` - Analytics and reporting
- [ ] `domain/shared/` - Shared domain concepts

**Success Criteria:**

- [ ] Domain layer has no external dependencies
- [ ] Clean separation of business logic
- [ ] Repository interfaces defined
- [ ] Domain services designed

---

### Phase 3 - DTO and HTTP Layer ⏳ PENDING

**Priority: High** | **Dependencies:** Phase 2 ✅ | **Estimated:** 3-4 hours

#### Tasks:

- [ ] Create request/response DTOs for each endpoint
- [ ] Implement DTO validation
- [ ] Create domain ↔ DTO mapping functions
- [ ] Restructure HTTP handlers by domain
- [ ] Implement consistent error response format

#### DTOs to Create:

- [ ] `http/dto/bowelmovement/` - BM request/response DTOs
- [ ] `http/dto/user/` - User management DTOs
- [ ] `http/dto/meal/` - Meal tracking DTOs
- [ ] `http/dto/symptom/` - Symptom DTOs
- [ ] `http/dto/medication/` - Medication DTOs
- [ ] `http/dto/analytics/` - Analytics DTOs
- [ ] `http/dto/shared/` - Common DTOs and responses

**Success Criteria:**

- [ ] No domain models in HTTP layer
- [ ] Consistent request/response patterns
- [ ] Proper validation on all inputs
- [ ] Clean error handling

---

### Phase 4 - Service Layer Implementation ⏳ PENDING

**Priority: Medium** | **Dependencies:** Phase 2, 3 ✅ | **Estimated:** 4-5 hours

#### Tasks:

- [ ] Implement domain service interfaces
- [ ] Move business logic from handlers to services
- [ ] Implement service-to-service communication
- [ ] Add proper transaction handling
- [ ] Implement domain events

#### Services to Implement:

- [ ] `bowelmovement.Service` - BM business logic
- [ ] `user.Service` - User management
- [ ] `meal.Service` - Meal management
- [ ] `symptom.Service` - Symptom tracking
- [ ] `medication.Service` - Medication management
- [ ] `analytics.Service` - Analytics calculations
- [ ] `auth.Service` - Authentication logic

**Success Criteria:**

- [ ] All business logic in services
- [ ] Handlers only handle HTTP concerns
- [ ] Services communicate through interfaces
- [ ] Proper error handling and logging

---

### Phase 5 - Dependency Injection Container ⏳ PENDING

**Priority: Medium** | **Dependencies:** Phase 4 ✅ | **Estimated:** 2-3 hours

#### Tasks:

- [ ] Implement dependency injection container
- [ ] Create container builder pattern
- [ ] Wire up all dependencies
- [ ] Replace manual dependency creation
- [ ] Add container validation

#### Container Components:

- [ ] `Container` struct with all dependencies
- [ ] `NewContainer()` builder function
- [ ] Repository container
- [ ] Service container
- [ ] Handler container
- [ ] Middleware container

**Success Criteria:**

- [ ] All dependencies injected through container
- [ ] Easy to swap implementations
- [ ] Clean main.go with minimal wiring
- [ ] Enhanced testability

---

### Phase 6 - Configuration Management ⏳ PENDING

**Priority: Low** | **Dependencies:** Phase 5 ✅ | **Estimated:** 2 hours

#### Tasks:

- [ ] Centralize configuration management
- [ ] Environment variable handling
- [ ] Configuration validation
- [ ] Default value management
- [ ] Configuration documentation

#### Configuration Areas:

- [ ] Server configuration (port, timeouts)
- [ ] Database configuration
- [ ] JWT configuration (secret, expiry)
- [ ] Security configuration (rate limits, etc.)
- [ ] Logging configuration
- [ ] Feature flags

**Success Criteria:**

- [ ] Single source of configuration truth
- [ ] Environment-based configuration
- [ ] Validation of required settings
- [ ] Clear documentation

---

### Phase 7 - Enhanced Testing ⏳ PENDING

**Priority: Medium** | **Dependencies:** Phase 5 ✅ | **Estimated:** 3-4 hours

#### Tasks:

- [ ] Update tests for new architecture
- [ ] Implement mock repositories for unit tests
- [ ] Create integration test helpers
- [ ] Add service layer tests
- [ ] Enhance test coverage

#### Testing Improvements:

- [ ] Domain service unit tests with mocks
- [ ] HTTP handler integration tests
- [ ] Repository integration tests
- [ ] End-to-end API tests
- [ ] Performance tests

**Success Criteria:**

- [ ] All tests pass with new architecture
- [ ] High test coverage maintained
- [ ] Mocking enables fast unit tests
- [ ] Integration tests verify behavior

---

### Phase 8 - Documentation and Cleanup ⏳ PENDING

**Priority: Low** | **Dependencies:** Phase 7 ✅ | **Estimated:** 1-2 hours

#### Tasks:

- [ ] Update documentation for new architecture
- [ ] Create architecture decision records (ADRs)
- [ ] Clean up old code and comments
- [ ] Add code generation scripts if needed
- [ ] Performance optimization

#### Documentation:

- [ ] Architecture overview
- [ ] Package structure explanation
- [ ] Dependency injection guide
- [ ] Testing strategy
- [ ] Development workflow

**Success Criteria:**

- [ ] Clear documentation for maintainers
- [ ] No dead code or old patterns
- [ ] Consistent code style
- [ ] Performance baselines established

---

## Progress Tracking

### 📊 Overall Statistics

- **Total Phases:** 8
- **Completed:** 0 (0%)
- **In Progress:** 1 (12.5%)
- **Pending:** 7 (87.5%)
- **Estimated Total Time:** 20-25 hours

### 🎯 Current Focus

**Active Phase:** Phase 1 - Infrastructure Setup
**Next Actions:**

1. Create new directory structure
2. Move existing files to appropriate locations
3. Update import paths
4. Verify tests pass

### 📈 Success Metrics

- [ ] **Maintainability**: Reduced cyclomatic complexity
- [ ] **Testability**: All services mockable, fast unit tests
- [ ] **Modularity**: Clear separation of concerns
- [ ] **Extensibility**: Easy to add new features
- [ ] **Performance**: No regression in API response times
- [ ] **Documentation**: Clear architecture documentation

---

## Risk Assessment

### 🔴 High Risk

- Breaking existing functionality during migration
- Import path conflicts during transition
- Test failures due to structural changes

### 🟡 Medium Risk

- Increased complexity during transition period
- Temporary reduction in development velocity
- Integration issues between layers

### 🟢 Low Risk

- Configuration changes
- Documentation updates
- Performance optimizations

---

## Rollback Plan

1. **Git branching strategy**: Work in feature branch `refactor/clean-architecture`
2. **Incremental commits**: Each phase has its own commits
3. **Testing checkpoints**: All tests must pass before proceeding
4. **Backup strategy**: Keep old structure until refactor complete
5. **Gradual migration**: Run both old and new side-by-side when possible

---

## Last Updated: 2025-06-20
