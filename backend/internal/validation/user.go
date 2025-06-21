package validation

import (
	"regexp"
	"time"
)

// ValidateUserID validates a user ID (non-empty string)
func ValidateUserID(userID string) error {
	if userID == "" {
		return ValidationError{Field: "userId", Message: "cannot be empty"}
	}
	if len(userID) > 100 {
		return ValidationError{Field: "userId", Message: "cannot be longer than 100 characters"}
	}
	return nil
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	if email == "" {
		return ValidationError{Field: "email", Message: "cannot be empty"}
	}
	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "invalid email format"}
	}
	if len(email) > 254 {
		return ValidationError{Field: "email", Message: "cannot be longer than 254 characters"}
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ValidationError{Field: "password", Message: "must be at least 8 characters long"}
	}
	if len(password) > 128 {
		return ValidationError{Field: "password", Message: "cannot be longer than 128 characters"}
	}

	// Check for at least one uppercase, one lowercase, one digit
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit {
		return ValidationError{
			Field:   "password",
			Message: "must contain at least one uppercase letter, one lowercase letter, and one digit",
		}
	}
	return nil
}

// ValidateName validates a person's name
func ValidateName(name string) error {
	if name != "" && len(name) > 100 {
		return ValidationError{Field: "name", Message: "cannot be longer than 100 characters"}
	}
	return nil
}

// ValidateTimezone validates IANA timezone format
func ValidateTimezone(timezone string) error {
	if timezone == "" {
		return nil // Optional field
	}
	_, err := time.LoadLocation(timezone)
	if err != nil {
		return ValidationError{Field: "timezone", Message: "invalid IANA timezone format"}
	}
	return nil
}

// ValidateReminderTime validates HH:MM format (24-hour)
func ValidateReminderTime(reminderTime string) error {
	if reminderTime == "" {
		return nil // Optional field
	}
	timeRegex := regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !timeRegex.MatchString(reminderTime) {
		return ValidationError{Field: "reminderTime", Message: "must be in HH:MM format (24-hour)"}
	}
	return nil
}

// ValidateDataRetentionDays validates data retention period
func ValidateDataRetentionDays(days int) error {
	if days < 30 {
		return ValidationError{Field: "dataRetentionDays", Message: "must be at least 30 days"}
	}
	if days > 3650 { // 10 years max
		return ValidationError{Field: "dataRetentionDays", Message: "cannot exceed 3650 days (10 years)"}
	}
	return nil
}
