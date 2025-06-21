package calculator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

func TestHealthScoreCalculationSimple(t *testing.T) {
	hds := NewHealthDataSummarizer()
	hsc := NewHealthScoreCalculator()

	day := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	movements := []bowelmovement.BowelMovement{{RecordedAt: day, BristolType: 4, Pain: 1, Strain: 1, Satisfaction: 7}}
	meals := []meal.Meal{{MealTime: day.Add(8 * time.Hour), Calories: 500, FiberRich: true}}
	symptoms := []symptom.Symptom{{RecordedAt: day.Add(2 * time.Hour), Severity: 2}}

	bmSummary := hds.CalculateBowelMovementSummary(movements, 1)
	mealSummary := hds.CalculateMealSummary(meals, 1)
	symptomSummary := hds.CalculateSymptomSummary(symptoms, 1)
	medSummary := hds.CalculateMedicationSummary(nil)

	overall := hsc.CalculateOverallHealthScore(bmSummary, mealSummary, symptomSummary, medSummary)
	assert.InEpsilon(t, 89.0, overall, 0.01)
}
