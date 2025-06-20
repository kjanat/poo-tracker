package medication

import (
	"errors"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Domain errors
var (
	// Entity errors
	ErrMedicationNotFound = errors.New("medication not found")
	ErrDoseRecordNotFound = errors.New("dose record not found")

	// Validation errors
	ErrInvalidMedicationName     = errors.New("medication name is required and must be 1-200 characters")
	ErrInvalidDosage             = errors.New("dosage is required and must be 1-100 characters")
	ErrInvalidFrequency          = errors.New("frequency is required and must be 1-200 characters")
	ErrInvalidMedicationCategory = errors.New("invalid medication category")
	ErrInvalidMedicationForm     = errors.New("invalid medication form")
	ErrInvalidMedicationRoute    = errors.New("invalid medication route")
	ErrInvalidDateRange          = errors.New("end date must be after start date")
	ErrInvalidDoseTime           = errors.New("dose time cannot be in the future")
	ErrDoseTakenInFuture         = errors.New("dose taken time cannot be in the future")

	// Business rule errors
	ErrMedicationAlreadyActive   = errors.New("medication is already active")
	ErrMedicationNotActive       = errors.New("medication is not active")
	ErrMedicationAlreadyInactive = errors.New("medication is already inactive")
	ErrDuplicateMedication       = errors.New("duplicate medication entry detected")

	// Authorization errors
	ErrUserNotAuthorized = errors.New("user not authorized to access this medication")
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
