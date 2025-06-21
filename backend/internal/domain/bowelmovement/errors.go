package bowelmovement

import (
	"errors"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Domain errors
var (
	// Entity errors
	ErrBowelMovementNotFound        = errors.New("bowel movement not found")
	ErrBowelMovementDetailsNotFound = errors.New("bowel movement details not found")

	// Input validation errors
	ErrInvalidID         = errors.New("invalid bowel movement ID")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDateRangeTooLarge = errors.New("date range too large")
	ErrInvalidDateRange  = errors.New("invalid date range")

	// Field validation errors
	ErrInvalidBristolType       = errors.New("bristol type must be between 1 and 7")
	ErrInvalidPainLevel         = errors.New("pain level must be between 1 and 10")
	ErrInvalidStrainLevel       = errors.New("strain level must be between 1 and 10")
	ErrInvalidSatisfactionLevel = errors.New("satisfaction level must be between 1 and 10")
	ErrInvalidStressLevel       = errors.New("stress level must be between 1 and 10")
	ErrInvalidSleepQuality      = errors.New("sleep quality must be between 1 and 10")
	ErrInvalidExerciseIntensity = errors.New("exercise intensity must be between 1 and 10")
	ErrInvalidVolume            = errors.New("invalid volume value")
	ErrInvalidColor             = errors.New("invalid color value")
	ErrInvalidConsistency       = errors.New("invalid consistency value")
	ErrInvalidSmellLevel        = errors.New("invalid smell level value")

	// Authorization and business rule errors
	ErrUserNotAuthorized   = errors.New("user not authorized to access this bowel movement")
	ErrDetailsAlreadyExist = errors.New("details already exist for this bowel movement")
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
