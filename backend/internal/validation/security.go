package validation

import "regexp"

// ValidateStrongPassword validates a password against security policies
func ValidateStrongPassword(password string) ValidationErrors {
	var errors ValidationErrors

	if len(password) < 8 {
		errors.Add("password", "must be at least 8 characters long")
	}

	if len(password) > 128 {
		errors.Add("password", "cannot be longer than 128 characters")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors.Add("password", "must contain at least one uppercase letter")
	}

	if !hasLower {
		errors.Add("password", "must contain at least one lowercase letter")
	}

	if !hasDigit {
		errors.Add("password", "must contain at least one digit")
	}

	if !hasSpecial {
		errors.Add("password", "must contain at least one special character")
	}

	// Check for common weak passwords
	weakPasswords := []string{
		"password", "123456", "password123", "admin", "qwerty",
		"letmein", "welcome", "monkey", "dragon", "password1",
	}

	for _, weak := range weakPasswords {
		if password == weak {
			errors.Add("password", "cannot be a commonly used weak password")
			break
		}
	}

	return errors
}

// ValidatePasswordPolicy validates password meets security requirements
func ValidatePasswordPolicy(password string) error {
	errors := ValidateStrongPassword(password)
	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidatePasswordStrength returns a password strength score (1-10)
func ValidatePasswordStrength(password string) int {
	score := 0

	// Length scoring
	if len(password) >= 8 {
		score += 2
	}
	if len(password) >= 12 {
		score += 1
	}
	if len(password) >= 16 {
		score += 1
	}

	// Character type scoring
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score += 1
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score += 1
	}
	if regexp.MustCompile(`[0-9]`).MatchString(password) {
		score += 1
	}
	if regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		score += 2
	}

	// Diversity scoring
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	if len(uniqueChars) > len(password)/2 {
		score += 1
	}

	if score > 10 {
		score = 10
	}

	return score
}
