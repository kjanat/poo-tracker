package medication

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrMedicationNotFound        = errors.New("medication not found")
	ErrMedicationNameRequired    = errors.New("medication name is required")
	ErrDosageRequired            = errors.New("dosage is required")
	ErrFrequencyRequired         = errors.New("frequency is required")
	ErrInvalidMedicationName     = errors.New("medication name is required and must be 1-200 characters")
	ErrInvalidDosage             = errors.New("dosage is required and must be 1-100 characters")
	ErrInvalidFrequency          = errors.New("frequency is required and must be 1-200 characters")
	ErrInvalidMedicationCategory = errors.New("invalid medication category")
	ErrInvalidMedicationForm     = errors.New("invalid medication form")
	ErrInvalidMedicationRoute    = errors.New("invalid medication route")
	ErrInvalidDateRange          = errors.New("end date must be after start date")
	ErrMedicationAlreadyActive   = errors.New("medication is already active")
	ErrMedicationNotActive       = errors.New("medication is not active")
	ErrMedicationAlreadyInactive = errors.New("medication is already inactive")
	ErrDoseRecordNotFound        = errors.New("dose record not found")
	ErrDoseTakenInFuture         = errors.New("dose taken time cannot be in the future")
	ErrInvalidDoseTime           = errors.New("dose time cannot be in the future")
	ErrUserNotAuthorized         = errors.New("user not authorized to access this medication")
	ErrDuplicateMedication       = errors.New("duplicate medication entry detected")
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
