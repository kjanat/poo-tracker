package shared

import "errors"

// Common domain errors
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden access")
	ErrInvalidID         = errors.New("invalid ID")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrDateRangeTooLarge = errors.New("date range too large")
	ErrInvalidDateRange  = errors.New("invalid date range")
)
