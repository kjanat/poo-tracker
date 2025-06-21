package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

func TestAnalyzeMealTimings(t *testing.T) {
	ta := &TrendAnalyzer{}
	meals := []meal.Meal{
		{MealTime: time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)},
		{MealTime: time.Date(2024, 6, 1, 12, 15, 0, 0, time.UTC)},
		{MealTime: time.Date(2024, 6, 2, 12, 30, 0, 0, time.UTC)},
		{MealTime: time.Date(2024, 6, 3, 20, 0, 0, 0, time.UTC)},
	}

	got := ta.analyzeMealTimings(meals)
	expected := []shared.MealTiming{
		{TimeOfDay: shared.NewTimeOfDay(8, 0), Frequency: 1},
		{TimeOfDay: shared.NewTimeOfDay(12, 0), Frequency: 2},
		{TimeOfDay: shared.NewTimeOfDay(20, 0), Frequency: 1},
	}
	assert.Equal(t, expected, got)
}

func TestIdentifyCommonIngredients(t *testing.T) {
	ta := &TrendAnalyzer{}
	spicy := 3
	meals := []meal.Meal{
		{Dairy: true},
		{Dairy: true, Gluten: true},
		{Gluten: true},
		{SpicyLevel: &spicy},
	}

	got := ta.identifyCommonIngredients(meals)
	assert.Equal(t, []string{"dairy", "gluten", "spicy"}, got)
}

func TestIdentifyCommonIngredientsFiber(t *testing.T) {
	ta := &TrendAnalyzer{}
	meals := []meal.Meal{
		{FiberRich: true},
	}

	got := ta.identifyCommonIngredients(meals)
	assert.Equal(t, []string{"fiber"}, got)
}

func TestIdentifyProblemIngredients(t *testing.T) {
	ta := &TrendAnalyzer{}
	// Test with insufficient data
	meals := []meal.Meal{{Dairy: true}}
	assert.Empty(t, ta.identifyProblemIngredients(meals))

	// Test with problem correlation data (if available)
	// This would require mock data showing correlation between ingredients and symptoms
}

func TestAnalyzeEatingPatterns(t *testing.T) {
	ta := &TrendAnalyzer{}
	spicy := 5
	meals := []meal.Meal{
		{MealTime: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC), Dairy: true},
		{MealTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), Dairy: true, Gluten: true},
		{MealTime: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC), Gluten: true},
		{MealTime: time.Date(2024, 1, 3, 20, 0, 0, 0, time.UTC), SpicyLevel: &spicy},
	}

	pattern := ta.AnalyzeEatingPatterns(meals)
	expectedTimings := []shared.MealTiming{
		{TimeOfDay: shared.NewTimeOfDay(8, 0), Frequency: 1},
		{TimeOfDay: shared.NewTimeOfDay(12, 0), Frequency: 2},
		{TimeOfDay: shared.NewTimeOfDay(20, 0), Frequency: 1},
	}
	assert.Equal(t, expectedTimings, pattern.MealTimings)
	assert.Equal(t, []string{"dairy", "gluten", "spicy"}, pattern.CommonIngredients)
	assert.Empty(t, pattern.ProblemIngredients)
}
