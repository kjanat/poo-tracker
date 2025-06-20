package meal

import (
	"errors"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Domain errors
var (
	// Input validation errors
	ErrInvalidID         = errors.New("invalid meal ID")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidCalories   = errors.New("calories must be between 0 and 10000")
	ErrInvalidSpicyLevel = errors.New("spicy level must be between 1 and 10")
	ErrDateRangeTooLarge = errors.New("date range too large")
	ErrInvalidDateRange  = errors.New("invalid date range")

	// Entity and field validation errors
	ErrMealNotFound        = errors.New("meal not found")
	ErrInvalidMealName     = errors.New("meal name is required and must be 1-200 characters")
	ErrInvalidMealCategory = errors.New("invalid meal category")
	ErrMealTimeRequired    = errors.New("meal time is required")
	ErrMealTimeInFuture    = errors.New("meal time cannot be in the future")

	// Authorization errors
	ErrUserNotAuthorized = errors.New("user not authorized to access this meal")
)

// ValidationError represents a validation error with field information.
// Deprecated: use shared.ValidationError instead.
type ValidationError = shared.ValidationError

// NewValidationError creates a new validation error.
// Deprecated: use shared.NewValidationError instead.
var NewValidationError = shared.NewValidationError

// BusinessRuleError represents a business rule violation.
// Deprecated: use shared.BusinessRuleError instead.
type BusinessRuleError = shared.BusinessRuleError

// NewBusinessRuleError creates a new business rule error.
// Deprecated: use shared.NewBusinessRuleError instead.
var NewBusinessRuleError = shared.NewBusinessRuleError
