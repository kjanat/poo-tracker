package aggregator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

func TestAggregateDailyData(t *testing.T) {
	da := NewDataAggregator()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	movements := []bowelmovement.BowelMovement{
		{RecordedAt: start, BristolType: 4, Pain: 2, Strain: 1, Satisfaction: 7},
		{RecordedAt: end, BristolType: 5, Pain: 3, Strain: 1, Satisfaction: 6},
	}
	spicy := 3
	meals := []meal.Meal{
		{MealTime: start.Add(12 * time.Hour), Calories: 500, SpicyLevel: &spicy, Dairy: true},
		{MealTime: end.Add(18 * time.Hour), Calories: 700, Gluten: true},
	}
	symptoms := []symptom.Symptom{
		{RecordedAt: start.Add(10 * time.Hour), Severity: 4},
	}
	meds := []medication.Medication{
		{Name: "med1", IsActive: true, StartDate: &start},
	}

	got := da.AggregateDailyData(movements, meals, symptoms, meds, start, end)
	assert.Len(t, got, 2)

	d1 := got[0]
	assert.Equal(t, start, d1.Date)
	assert.Equal(t, 1, d1.MealCount)
	assert.Equal(t, 1, d1.SpicyMealCount)
	assert.Equal(t, 1, d1.DairyMealCount)
	assert.Equal(t, 1, d1.SymptomCount)
}

func TestCalculateAverage(t *testing.T) {
	da := NewDataAggregator()
	avg := da.CalculateAverage([]float64{1, 2, 3, 4})
	assert.Equal(t, 2.5, avg)
}

func TestCalculateAverageSymptomSeverity(t *testing.T) {
	da := NewDataAggregator()
	symptoms := []symptom.Symptom{
		{Severity: 3},
		{Severity: 5},
	}
	avg := da.CalculateAverageSymptomSeverity(symptoms)
	assert.Equal(t, 4.0, avg)
}

func TestGroupMealsByType(t *testing.T) {
	da := NewDataAggregator()
	spicy := 4
	meals := []meal.Meal{
		{Calories: 900},
		{SpicyLevel: &spicy},
		{Dairy: true},
		{Gluten: true},
		{Calories: 500, FiberRich: true},
	}
	groups := da.GroupMealsByType(meals)
	assert.Len(t, groups["large"], 1)
	assert.Len(t, groups["spicy"], 1)
	assert.Len(t, groups["dairy"], 1)
	assert.Len(t, groups["gluten"], 1)
	assert.Len(t, groups["healthy"], 1)
}

func TestGetActiveMedicationPercentage(t *testing.T) {
	da := NewDataAggregator()
	meds := []medication.Medication{
		{IsActive: true},
		{IsActive: false},
	}
	pct := da.GetActiveMedicationPercentage(meds, 7)
	assert.Equal(t, 50.0, pct)
}
