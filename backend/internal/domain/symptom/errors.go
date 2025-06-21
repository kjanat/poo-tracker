package symptom

import (
	"errors"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
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
