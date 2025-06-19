package validation

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/model"
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

// Error types for validation
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors holds multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	msg := fmt.Sprintf("%d validation errors: ", len(e))
	for i, err := range e {
		if i > 0 {
			msg += "; "
		}
		msg += err.Error()
	}
	return msg
}

// safeAppendValidationError safely appends an error to ValidationErrors if it's a ValidationError
func safeAppendValidationError(errors *ValidationErrors, err error) {
	if validationErr, ok := err.(ValidationError); ok {
		*errors = append(*errors, validationErr)
	}
}

// HasErrors returns true if there are any validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Add adds a new validation error
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, ValidationError{Field: field, Message: message})
}

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

// ValidateBowelMovement validates a complete BowelMovement model
func ValidateBowelMovement(bm model.BowelMovement) ValidationErrors {
	var errors ValidationErrors

	errors = append(errors, validateBowelMovementRequired(bm)...)  // Required fields
	errors = append(errors, validateBowelMovementEnums(bm)...)     // Enum fields
	errors = append(errors, validateBowelMovementScales(bm)...)    // Scale fields
	errors = append(errors, validateBowelMovementOptionals(bm)...) // Optional fields

	return errors
}

func validateBowelMovementRequired(bm model.BowelMovement) ValidationErrors {
	var errors ValidationErrors
	if err := ValidateUserID(bm.UserID); err != nil {
		safeAppendValidationError(&errors, err)
	}
	if err := ValidateBristolType(bm.BristolType); err != nil {
		safeAppendValidationError(&errors, err)
	}
	return errors
}

func validateBowelMovementEnums(bm model.BowelMovement) ValidationErrors {
	var errors ValidationErrors
	if bm.Volume != nil {
		if err := ValidateEnum(*bm.Volume, "volume"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if bm.Color != nil {
		if err := ValidateEnum(*bm.Color, "color"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if bm.Consistency != nil {
		if err := ValidateEnum(*bm.Consistency, "consistency"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if bm.SmellLevel != nil {
		if err := ValidateEnum(*bm.SmellLevel, "smellLevel"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	return errors
}

func validateBowelMovementScales(bm model.BowelMovement) ValidationErrors {
	var errors ValidationErrors
	if err := ValidateScale(bm.Pain, "pain"); err != nil {
		errors = append(errors, err.(ValidationError))
	}
	if err := ValidateScale(bm.Strain, "strain"); err != nil {
		errors = append(errors, err.(ValidationError))
	}
	if err := ValidateScale(bm.Satisfaction, "satisfaction"); err != nil {
		errors = append(errors, err.(ValidationError))
	}
	return errors
}

func validateBowelMovementOptionals(bm model.BowelMovement) ValidationErrors {
	var errors ValidationErrors
	if err := ValidateURL(bm.PhotoURL, "photoUrl"); err != nil {
		errors = append(errors, err.(ValidationError))
	}
	if err := ValidateNotes(bm.Notes, "notes"); err != nil {
		errors = append(errors, err.(ValidationError))
	}
	return errors
}

// ValidateMeal validates a complete Meal model
func ValidateMeal(meal model.Meal) ValidationErrors {
	var errors ValidationErrors

	// Required fields
	if err := ValidateUserID(meal.UserID); err != nil {
		errors = append(errors, err.(ValidationError))
	}

	if err := ValidateMealName(meal.Name); err != nil {
		errors = append(errors, err.(ValidationError))
	}

	// Optional enum fields
	if meal.Category != nil {
		if err := ValidateEnum(*meal.Category, "category"); err != nil {
			errors = append(errors, err.(ValidationError))
		}
	}

	// Optional scale fields
	if meal.SpicyLevel != nil {
		if err := ValidateSpicyLevel(*meal.SpicyLevel); err != nil {
			errors = append(errors, err.(ValidationError))
		}
	}

	// Optional fields
	if err := ValidateURL(meal.PhotoURL, "photoUrl"); err != nil {
		errors = append(errors, err.(ValidationError))
	}

	if err := ValidateNotes(meal.Description, "description"); err != nil {
		errors = append(errors, err.(ValidationError))
	}

	if err := ValidateNotes(meal.Notes, "notes"); err != nil {
		errors = append(errors, err.(ValidationError))
	}

	return errors
}
