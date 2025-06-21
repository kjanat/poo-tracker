package validation

import (
	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// ValidateBowelMovement validates a complete BowelMovement model
func ValidateBowelMovement(bm bm.BowelMovement) ValidationErrors {
	var errors ValidationErrors

	errors = append(errors, validateBowelMovementRequired(bm)...)  // Required fields
	errors = append(errors, validateBowelMovementEnums(bm)...)     // Enum fields
	errors = append(errors, validateBowelMovementScales(bm)...)    // Scale fields
	errors = append(errors, validateBowelMovementOptionals(bm)...) // Optional fields

	return errors
}

func validateBowelMovementRequired(bm bm.BowelMovement) ValidationErrors {
	var errors ValidationErrors
	if err := ValidateUserID(bm.UserID); err != nil {
		safeAppendValidationError(&errors, err)
	}
	if err := ValidateBristolType(bm.BristolType); err != nil {
		safeAppendValidationError(&errors, err)
	}
	return errors
}

func validateBowelMovementEnums(bm bm.BowelMovement) ValidationErrors {
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

func validateBowelMovementScales(bm bm.BowelMovement) ValidationErrors {
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

func validateBowelMovementOptionals(bm bm.BowelMovement) ValidationErrors {
	var errors ValidationErrors
	if err := ValidateURL(bm.PhotoURL, "photoUrl"); err != nil {
		errors = append(errors, err.(ValidationError))
	}
	// Notes validation moved to BowelMovementDetails
	return errors
}

// ValidateMeal validates a complete Meal model
func ValidateMeal(meal meal.Meal) ValidationErrors {
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
