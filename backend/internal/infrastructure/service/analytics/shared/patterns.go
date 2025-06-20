package shared

import (
	"fmt"
	"time"
)

// EatingPattern represents identified patterns in eating habits
type EatingPattern struct {
	MealTimings        []MealTiming
	CommonIngredients  []string
	ProblemIngredients []string
}

// BowelPattern represents identified patterns in bowel movements
type BowelPattern struct {
	Frequency       float64
	Consistency     float64
	MealCorrelation float64
}

// SymptomPattern represents identified patterns in symptoms
type SymptomPattern struct {
	CommonSymptoms map[string]int
	Frequency      map[string]int
	Severity       map[string]float64
}

// LifestylePattern represents overall lifestyle patterns
type LifestylePattern struct {
	DietaryHabits   []DietaryHabit
	BowelRegularity float64
	SymptomTriggers []SymptomTrigger
}

// TimeOfDay represents a time of day without date information
type TimeOfDay struct {
	Hour   int `json:"hour"`   // 0-23
	Minute int `json:"minute"` // 0-59
}

// NewTimeOfDay creates a new TimeOfDay from hour and minute
func NewTimeOfDay(hour, minute int) TimeOfDay {
	return TimeOfDay{
		Hour:   hour % 24,
		Minute: minute % 60,
	}
}

// FromTime creates a TimeOfDay from a time.Time
func (tod *TimeOfDay) FromTime(t time.Time) {
	tod.Hour = t.Hour()
	tod.Minute = t.Minute()
}

// String returns the time in HH:MM format
func (tod TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", tod.Hour, tod.Minute)
}

// MealTiming represents timing patterns for meals
type MealTiming struct {
	TimeOfDay TimeOfDay
	Frequency int
}

// DietaryHabit represents identified dietary habits
type DietaryHabit struct {
	Description string
	Frequency   int
	Impact      float64
}

// SymptomTrigger represents identified triggers for symptoms
type SymptomTrigger struct {
	TriggerType string
	Ingredient  string
	Confidence  float64
}
