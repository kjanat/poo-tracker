package validation

import "regexp"

// ValidateNotes validates notes/description fields
func ValidateNotes(notes string, fieldName string) error {
	if len(notes) > 2000 {
		return ValidationError{
			Field:   fieldName,
			Message: "cannot be longer than 2000 characters",
		}
	}
	return nil
}

// ValidateMealName validates meal name
func ValidateMealName(name string) error {
	if name == "" {
		return ValidationError{Field: "name", Message: "cannot be empty"}
	}
	if len(name) > 200 {
		return ValidationError{Field: "name", Message: "cannot be longer than 200 characters"}
	}
	return nil
}

// ValidateCalories validates calorie count
func ValidateCalories(calories int) error {
	if calories < 0 {
		return ValidationError{Field: "calories", Message: "cannot be negative"}
	}
	if calories > 50000 { // Reasonable upper limit
		return ValidationError{Field: "calories", Message: "cannot exceed 50,000 calories"}
	}
	return nil
}

// ValidateSpicyLevel validates spicy level (1-10 scale)
func ValidateSpicyLevel(level int) error {
	return ValidateScale(level, "spicyLevel")
}

// ValidateURL validates a URL format (for photo URLs)
func ValidateURL(url string, fieldName string) error {
	if url == "" {
		return nil // Optional field
	}
	// Strict regex for HTTP/HTTPS URLs only with proper structure
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*(:[0-9]{1,5})?(/[^\s]*)?$`)
	if !urlRegex.MatchString(url) {
		return ValidationError{Field: fieldName, Message: "invalid URL format - only HTTP/HTTPS URLs are allowed"}
	}
	if len(url) > 2048 {
		return ValidationError{Field: fieldName, Message: "URL cannot be longer than 2048 characters"}
	}
	return nil
}
