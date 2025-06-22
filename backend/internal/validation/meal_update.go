package validation

import (
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
)

// ValidateMealUpdate validates fields for updating a Meal.
func ValidateMealUpdate(update meal.MealUpdate) ValidationErrors {
	var errors ValidationErrors

	if update.Name != nil {
		if err := ValidateMealName(*update.Name); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if update.Category != nil {
		if err := ValidateEnum(*update.Category, "category"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if update.Calories != nil {
		if err := ValidateCalories(*update.Calories); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if update.SpicyLevel != nil {
		if err := ValidateSpicyLevel(*update.SpicyLevel); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if update.PhotoURL != nil {
		if err := ValidateURL(*update.PhotoURL, "photoUrl"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if update.Description != nil {
		if err := ValidateNotes(*update.Description, "description"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	if update.Notes != nil {
		if err := ValidateNotes(*update.Notes, "notes"); err != nil {
			safeAppendValidationError(&errors, err)
		}
	}
	// Add more field validations as needed
	return errors
}
