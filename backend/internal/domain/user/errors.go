package user

import (
	"errors"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Domain errors
var (
	// Input validation errors
	ErrInvalidID        = errors.New("invalid user ID")
	ErrInvalidInput     = errors.New("invalid input")
	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordTooLong  = errors.New("password too long")
	ErrSamePassword     = errors.New("new password must be different from current password")

	// Entity errors
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAuthNotFound     = errors.New("user authentication not found")
	ErrUserSettingsNotFound = errors.New("user settings not found")

	// Field validation errors
	ErrInvalidEmail         = errors.New("invalid email format")
	ErrInvalidUsername      = errors.New("invalid username format")
	ErrWeakPassword         = errors.New("password is too weak")
	ErrPasswordMismatch     = errors.New("password does not match")
	ErrInvalidAge           = errors.New("age must be between 1 and 150")
	ErrInvalidHeight        = errors.New("height must be between 50 and 300 cm")
	ErrInvalidWeight        = errors.New("weight must be between 20 and 500 kg")
	ErrInvalidPrivacyLevel  = errors.New("privacy level must be between 1 and 5")
	ErrInvalidDataRetention = errors.New("data retention days must be between 1 and 3650")
	ErrInvalidReminderTime  = errors.New("reminder time must be in HH:MM format")
	ErrInvalidTimezone      = errors.New("invalid timezone")

	// Business rule errors
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrAuthAlreadyExists     = errors.New("user authentication already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrAccountDeactivated    = errors.New("account is deactivated")

	// Authorization errors
	ErrUserNotAuthorized = errors.New("user not authorized to perform this action")
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
