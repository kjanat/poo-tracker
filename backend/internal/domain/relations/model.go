package relations

import "time"

// CorrelationType represents the type of correlation between entities
// e.g. how strongly a meal influences a symptom
//go:generate stringer -type=CorrelationType

type CorrelationType string

const (
	CorrelationPositive CorrelationType = "POSITIVE"
	CorrelationNegative CorrelationType = "NEGATIVE"
	CorrelationNeutral  CorrelationType = "NEUTRAL"
	CorrelationUnknown  CorrelationType = "UNKNOWN"
)

// AllCorrelationTypes returns all valid CorrelationType values
func AllCorrelationTypes() []CorrelationType {
	return []CorrelationType{
		CorrelationPositive,
		CorrelationNegative,
		CorrelationNeutral,
		CorrelationUnknown,
	}
}

// IsValid checks if the CorrelationType value is valid
func (ct CorrelationType) IsValid() bool {
	for _, valid := range AllCorrelationTypes() {
		if ct == valid {
			return true
		}
	}
	return false
}

// MealBowelMovementRelation represents the relationship between meals and bowel movements
// captured by the user. TimeGapHours indicates the time difference between the
// meal and the movement.
type MealBowelMovementRelation struct {
	ID              string           `json:"id"`
	UserID          string           `json:"userId"`
	MealID          string           `json:"mealId"`
	BowelMovementID string           `json:"bowelMovementId"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	Strength        int              `json:"strength"`
	Notes           string           `json:"notes,omitempty"`
	TimeGapHours    float64          `json:"timeGapHours"`
	UserCorrelation *CorrelationType `json:"userCorrelation,omitempty"`
}

// MealSymptomRelation represents the relationship between meals and symptoms
// The structure mirrors MealBowelMovementRelation.
type MealSymptomRelation struct {
	ID              string           `json:"id"`
	UserID          string           `json:"userId"`
	MealID          string           `json:"mealId"`
	SymptomID       string           `json:"symptomId"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	Strength        int              `json:"strength"`
	Notes           string           `json:"notes,omitempty"`
	TimeGapHours    float64          `json:"timeGapHours"`
	UserCorrelation *CorrelationType `json:"userCorrelation,omitempty"`
}

// NewMealBowelMovementRelation creates a new meal-bowel movement relation
// Default strength is set to 5 (neutral on 1-10 scale)
func NewMealBowelMovementRelation(userID, mealID, bowelMovementID string, timeGapHours float64) MealBowelMovementRelation {
	if userID == "" || mealID == "" || bowelMovementID == "" {
		panic("userID, mealID, and bowelMovementID cannot be empty")
	}
	if timeGapHours < 0 {
		panic("timeGapHours cannot be negative")
	}

	now := time.Now()
	return MealBowelMovementRelation{
		UserID:          userID,
		MealID:          mealID,
		BowelMovementID: bowelMovementID,
		CreatedAt:       now,
		UpdatedAt:       now,
		TimeGapHours:    timeGapHours,
		Strength:        5,
	}
}

// NewMealSymptomRelation creates a new meal-symptom relation
// Default strength is set to 5 (neutral on 1-10 scale)
func NewMealSymptomRelation(userID, mealID, symptomID string, timeGapHours float64) MealSymptomRelation {
	if userID == "" || mealID == "" || symptomID == "" {
		panic("userID, mealID, and symptomID cannot be empty")
	}
	if timeGapHours < 0 {
		panic("timeGapHours cannot be negative")
	}

	now := time.Now()
	return MealSymptomRelation{
		UserID:       userID,
		MealID:       mealID,
		SymptomID:    symptomID,
		CreatedAt:    now,
		UpdatedAt:    now,
		TimeGapHours: timeGapHours,
		Strength:     5,
	}
}
