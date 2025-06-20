package symptom

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrSymptomNotFound        = errors.New("symptom not found")
	ErrSymptomNameRequired    = errors.New("symptom name is required")
	ErrInvalidSymptomName     = errors.New("symptom name is required and must be 1-200 characters")
	ErrInvalidSeverity        = errors.New("severity must be between 1 and 10")
	ErrInvalidDuration        = errors.New("duration must be positive")
	ErrInvalidSymptomCategory = errors.New("invalid symptom category")
	ErrInvalidSymptomType     = errors.New("invalid symptom type")
	ErrRecordedTimeRequired   = errors.New("recorded time is required")
	ErrRecordedTimeInFuture   = errors.New("recorded time cannot be in the future")
	ErrUserNotAuthorized      = errors.New("user not authorized to access this symptom")
	ErrDuplicateSymptom       = errors.New("duplicate symptom entry detected")
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
