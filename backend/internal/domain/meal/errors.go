package meal

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrMealNotFound        = errors.New("meal not found")
	ErrInvalidMealName     = errors.New("meal name is required and must be 1-200 characters")
	ErrInvalidCalories     = errors.New("calories must be between 0 and 10000")
	ErrInvalidSpicyLevel   = errors.New("spicy level must be between 1 and 10")
	ErrInvalidMealCategory = errors.New("invalid meal category")
	ErrMealTimeRequired    = errors.New("meal time is required")
	ErrMealTimeInFuture    = errors.New("meal time cannot be in the future")
	ErrUserNotAuthorized   = errors.New("user not authorized to access this meal")
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// BusinessRuleError represents a business rule violation
type BusinessRuleError struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func (e BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation '%s': %s", e.Rule, e.Message)
}

// NewBusinessRuleError creates a new business rule error
func NewBusinessRuleError(rule, message string) BusinessRuleError {
	return BusinessRuleError{
		Rule:    rule,
		Message: message,
	}
}
