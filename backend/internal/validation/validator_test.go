package validation

import (
	"testing"
	"time"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

func TestValidateBristolType(t *testing.T) {
	tests := []struct {
		name        string
		bristolType int
		wantErr     bool
	}{
		{"Valid 1", 1, false},
		{"Valid 4", 4, false},
		{"Valid 7", 7, false},
		{"Invalid 0", 0, true},
		{"Invalid 8", 8, true},
		{"Invalid negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBristolType(tt.bristolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBristolType(%d) error = %v, wantErr %v", tt.bristolType, err, tt.wantErr)
			}
		})
	}
}

func TestValidateScale(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		field   string
		wantErr bool
	}{
		{"Valid 1", 1, "pain", false},
		{"Valid 5", 5, "satisfaction", false},
		{"Valid 10", 10, "strain", false},
		{"Invalid 0", 0, "pain", true},
		{"Invalid 11", 11, "satisfaction", true},
		{"Invalid negative", -1, "strain", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScale(tt.value, tt.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScale(%d, %s) error = %v, wantErr %v", tt.value, tt.field, err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{"Valid short ID", "user123", false},
		{"Valid long ID", "user_with_long_id_but_still_valid_12345", false},
		{"Empty ID", "", true},
		{"Too long ID", string(make([]rune, 101)), true}, // 101 characters
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserID(%s) error = %v, wantErr %v", tt.userID, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid basic", "test@example.com", false},
		{"Valid with plus", "test+tag@example.com", false},
		{"Valid with subdomain", "user@mail.example.com", false},
		{"Empty email", "", true},
		{"Invalid format", "invalid-email", true},
		{"Missing @", "userexample.com", true},
		{"Missing domain", "user@", true},
		{"Missing TLD", "user@example", true},
		{"Too long", "test@" + string(make([]rune, 250)) + ".com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%s) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "Password123", false},
		{"Valid complex", "MySecure123!", false},
		{"Too short", "Pass1", true},
		{"No uppercase", "password123", true},
		{"No lowercase", "PASSWORD123", true},
		{"No digit", "Password", true},
		{"Too long", string(make([]rune, 129)), true}, // 129 characters
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%s) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTimezone(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{"Valid UTC", "UTC", false},
		{"Valid New York", "America/New_York", false},
		{"Valid London", "Europe/London", false},
		{"Empty (optional)", "", false},
		{"Invalid timezone", "Invalid/Timezone", true},
		// Note: "EST" is actually a valid timezone in Go, so removing this test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimezone(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimezone(%s) error = %v, wantErr %v", tt.timezone, err, tt.wantErr)
			}
		})
	}
}

func TestValidateReminderTime(t *testing.T) {
	tests := []struct {
		name         string
		reminderTime string
		wantErr      bool
	}{
		{"Valid 09:00", "09:00", false},
		{"Valid 23:59", "23:59", false},
		{"Valid 00:00", "00:00", false},
		{"Valid 9:00", "9:00", false}, // Single digit hour is actually valid with our regex
		{"Empty (optional)", "", false},
		{"Invalid hour", "25:00", true},
		{"Invalid minute", "09:60", true},
		{"Invalid format", "09-00", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReminderTime(tt.reminderTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReminderTime(%s) error = %v, wantErr %v", tt.reminderTime, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDataRetentionDays(t *testing.T) {
	tests := []struct {
		name    string
		days    int
		wantErr bool
	}{
		{"Valid 30", 30, false},
		{"Valid 365", 365, false},
		{"Valid 3650", 3650, false},
		{"Too few", 29, true},
		{"Too many", 3651, true},
		{"Negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDataRetentionDays(tt.days)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDataRetentionDays(%d) error = %v, wantErr %v", tt.days, err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		field   string
		wantErr bool
	}{
		{"Valid HTTP", "http://example.com", "photoUrl", false},
		{"Valid HTTPS", "https://example.com/path", "photoUrl", false},
		{"Valid with query", "https://example.com/path?query=1", "photoUrl", false},
		{"Empty (optional)", "", "photoUrl", false},
		{"Invalid format", "not-a-url", "photoUrl", true},
		{"Invalid protocol", "ftp://example.com", "photoUrl", true},
		{"Too long", "https://" + string(make([]rune, 2050)), "photoUrl", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url, tt.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%s, %s) error = %v, wantErr %v", tt.url, tt.field, err, tt.wantErr)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	var errors ValidationErrors

	// Test empty errors
	if errors.HasErrors() {
		t.Error("Empty ValidationErrors should not have errors")
	}

	// Test adding errors
	errors.Add("field1", "error message 1")
	errors.Add("field2", "error message 2")

	if !errors.HasErrors() {
		t.Error("ValidationErrors should have errors after adding")
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	// Test error message formatting
	errMsg := errors.Error()
	if errMsg == "" {
		t.Error("ValidationErrors.Error() should not be empty")
	}

	// Test single error message
	var singleError ValidationErrors
	singleError.Add("field", "message")
	singleErrMsg := singleError.Error()
	expected := "validation error on field 'field': message"
	if singleErrMsg != expected {
		t.Errorf("Single error message = %q, want %q", singleErrMsg, expected)
	}
}

func TestValidateBowelMovement(t *testing.T) {
	// Valid bowel movement
	validBM := bm.BowelMovement{
		UserID:       "user123",
		BristolType:  4,
		Pain:         3,
		Strain:       2,
		Satisfaction: 7,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		RecordedAt:   time.Now(),
	}

	errors := ValidateBowelMovement(validBM)
	if errors.HasErrors() {
		t.Errorf("Valid bowel movement should not have errors, got: %v", errors)
	}

	// Invalid bowel movement
	invalidBM := bm.BowelMovement{
		UserID:       "", // Invalid: empty
		BristolType:  8,  // Invalid: out of range
		Pain:         11, // Invalid: out of range
		Strain:       0,  // Invalid: out of range
		Satisfaction: -1, // Invalid: out of range
	}

	errors = ValidateBowelMovement(invalidBM)
	if !errors.HasErrors() {
		t.Error("Invalid bowel movement should have errors")
	}

	// Check that we have the expected number of errors
	expectedErrors := 5 // userID, bristolType, pain, strain, satisfaction
	if len(errors) != expectedErrors {
		t.Errorf("Expected %d errors, got %d: %v", expectedErrors, len(errors), errors)
	}
}

func TestValidateMeal(t *testing.T) {
	// Valid meal
	validMeal := meal.Meal{
		UserID:    "user123",
		Name:      "Test Meal",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		MealTime:  time.Now(),
	}

	errors := ValidateMeal(validMeal)
	if errors.HasErrors() {
		t.Errorf("Valid meal should not have errors, got: %v", errors)
	}

	// Invalid meal
	invalidMeal := meal.Meal{
		UserID: "", // Invalid: empty
		Name:   "", // Invalid: empty
	}

	errors = ValidateMeal(invalidMeal)
	if !errors.HasErrors() {
		t.Error("Invalid meal should have errors")
	}

	// Check that we have the expected number of errors
	expectedErrors := 2 // userID, name
	if len(errors) != expectedErrors {
		t.Errorf("Expected %d errors, got %d: %v", expectedErrors, len(errors), errors)
	}
}

func TestValidateEnum(t *testing.T) {
	// Valid enum
	volume := shared.VolumeSmall
	err := ValidateEnum(volume, "volume")
	if err != nil {
		t.Errorf("Valid enum should not have error, got: %v", err)
	}

	// Test with different valid enums
	color := shared.ColorBrown
	err = ValidateEnum(color, "color")
	if err != nil {
		t.Errorf("Valid color enum should not have error, got: %v", err)
	}

	consistency := shared.ConsistencySoft
	err = ValidateEnum(consistency, "consistency")
	if err != nil {
		t.Errorf("Valid consistency enum should not have error, got: %v", err)
	}
}
