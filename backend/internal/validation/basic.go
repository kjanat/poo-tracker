package validation

import (
	"fmt"

	rel "github.com/kjanat/poo-tracker/backend/internal/domain/relations"
)

// BristolTypeRange defines valid Bristol Stool Chart range (1-7)
const (
	BristolTypeMin = 1
	BristolTypeMax = 7
)

// ScaleRange defines valid scale ranges (1-10) for pain, strain, satisfaction, etc.
const (
	ScaleMin = 1
	ScaleMax = 10
)

// ValidateBristolType validates Bristol Stool Chart type (1-7)
func ValidateBristolType(bristolType int) error {
	if bristolType < BristolTypeMin || bristolType > BristolTypeMax {
		return ValidationError{
			Field:   "bristolType",
			Message: fmt.Sprintf("must be between %d and %d", BristolTypeMin, BristolTypeMax),
		}
	}
	return nil
}

// ValidateScale validates a 1-10 scale value (pain, strain, satisfaction, etc.)
func ValidateScale(value int, fieldName string) error {
	if value < ScaleMin || value > ScaleMax {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must be between %d and %d", ScaleMin, ScaleMax),
		}
	}
	return nil
}

// ValidateEnum validates enum values using the IsValid method
func ValidateEnum[T interface{ IsValid() bool }](value T, fieldName string) error {
	if !value.IsValid() {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("invalid value: %v", value),
		}
	}
	return nil
}

// IsValidCorrelationType validates a correlation type
func IsValidCorrelationType(correlationType string) bool {
	ct := rel.CorrelationType(correlationType)
	return ct.IsValid()
}
