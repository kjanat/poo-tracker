package bowelmovement

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	// Core errors
	ErrNotFound          = errors.New("bowel movement not found")
	ErrInvalidID         = errors.New("invalid bowel movement ID")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDateRangeTooLarge = errors.New("date range too large")
	ErrInvalidDateRange  = errors.New("invalid date range")

	// Existing errors
	ErrBowelMovementNotFound        = errors.New("bowel movement not found")
	ErrBowelMovementDetailsNotFound = errors.New("bowel movement details not found")
	ErrInvalidBristolType           = errors.New("bristol type must be between 1 and 7")
	ErrInvalidPainLevel             = errors.New("pain level must be between 1 and 10")
	ErrInvalidStrainLevel           = errors.New("strain level must be between 1 and 10")
	ErrInvalidSatisfactionLevel     = errors.New("satisfaction level must be between 1 and 10")
	ErrInvalidStressLevel           = errors.New("stress level must be between 1 and 10")
	ErrInvalidSleepQuality          = errors.New("sleep quality must be between 1 and 10")
	ErrInvalidExerciseIntensity     = errors.New("exercise intensity must be between 1 and 10")
	ErrInvalidVolume                = errors.New("invalid volume value")
	ErrInvalidColor                 = errors.New("invalid color value")
	ErrInvalidConsistency           = errors.New("invalid consistency value")
	ErrInvalidSmellLevel            = errors.New("invalid smell level value")
	ErrUserNotAuthorized            = errors.New("user not authorized to access this bowel movement")
	ErrDetailsAlreadyExist          = errors.New("details already exist for this bowel movement")
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
