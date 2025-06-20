package aggregator

import (
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics"
)

// AggregateBowelMovements calculates statistics for bowel movements
func (da *DataAggregator) AggregateBowelMovements(movements []bowelmovement.BowelMovement) *analytics.BowelMovementStats {
	if len(movements) == 0 {
		return &analytics.BowelMovementStats{}
	}

	// Sort movements by time
	sort.Slice(movements, func(i, j int) bool {
		return movements[i].RecordedAt.Before(movements[j].RecordedAt)
	})

	// Calculate averages and scores
	totalDays := movements[len(movements)-1].RecordedAt.Sub(movements[0].RecordedAt).Hours() / 24
	if totalDays < 1 {
		totalDays = 1
	}

	stats := &analytics.BowelMovementStats{
		TotalCount:    len(movements),
		AveragePerDay: float64(len(movements)) / totalDays,
		LastMovement:  movements[len(movements)-1].RecordedAt,
	}

	stats.ConsistencyScore = da.calculateConsistencyScore(movements)
	stats.RegularityScore = da.calculateRegularityScore(movements)

	return stats
}

// AggregateMeals calculates statistics for meals
func (da *DataAggregator) AggregateMeals(meals []meal.Meal) *analytics.MealStats {
	if len(meals) == 0 {
		return &analytics.MealStats{}
	}

	// Sort meals by time
	sort.Slice(meals, func(i, j int) bool {
		return meals[i].CreatedAt.Before(meals[j].CreatedAt)
	})

	// Calculate averages and collect ingredients
	totalDays := meals[len(meals)-1].CreatedAt.Sub(meals[0].CreatedAt).Hours() / 24
	if totalDays < 1 {
		totalDays = 1
	}

	ingredientFreq := make(map[string]int)
	// Ingredients processing is skipped for now - we'll implement it later
	// when the meal model is updated

	// Get most common ingredients
	type ingredientCount struct {
		name  string
		count int
	}
	var ingredients []ingredientCount
	for ing, count := range ingredientFreq {
		ingredients = append(ingredients, ingredientCount{ing, count})
	}
	sort.Slice(ingredients, func(i, j int) bool {
		return ingredients[i].count > ingredients[j].count
	})

	// Take top 5 common ingredients
	commonIngredients := make([]string, 0, 5)
	for i := 0; i < len(ingredients) && i < 5; i++ {
		commonIngredients = append(commonIngredients, ingredients[i].name)
	}

	stats := &analytics.MealStats{
		TotalCount:        len(meals),
		AveragePerDay:     float64(len(meals)) / totalDays,
		CommonIngredients: commonIngredients,
		LastMeal:          meals[len(meals)-1].CreatedAt,
	}

	stats.NutritionScore = da.calculateNutritionScore(meals)

	return stats
}

// AggregateSymptoms calculates statistics for symptoms
func (da *DataAggregator) AggregateSymptoms(symptoms []symptom.Symptom) *analytics.SymptomStats {
	if len(symptoms) == 0 {
		return &analytics.SymptomStats{
			CommonSymptoms: make(map[string]int),
			SeverityTrends: make(map[string]float64),
		}
	}

	// Sort symptoms by time
	sort.Slice(symptoms, func(i, j int) bool {
		return symptoms[i].CreatedAt.Before(symptoms[j].CreatedAt)
	})

	// Calculate averages and collect symptom frequencies
	totalDays := symptoms[len(symptoms)-1].CreatedAt.Sub(symptoms[0].CreatedAt).Hours() / 24
	if totalDays < 1 {
		totalDays = 1
	}

	symptomFreq := make(map[string]int)
	severitySum := make(map[string]float64)
	severityCount := make(map[string]int)

	for _, s := range symptoms {
		symptomType := s.Type.String() // Assuming Type is an enum with String() method
		symptomFreq[symptomType]++
		severitySum[symptomType] += float64(s.Severity)
		severityCount[symptomType]++
	}

	// Calculate average severity for each symptom
	severityTrends := make(map[string]float64)
	for sType, sum := range severitySum {
		severityTrends[sType] = sum / float64(severityCount[sType])
	}

	return &analytics.SymptomStats{
		TotalCount:      len(symptoms),
		AveragePerDay:   float64(len(symptoms)) / totalDays,
		CommonSymptoms:  symptomFreq,
		SeverityTrends:  severityTrends,
		LastSymptomTime: symptoms[len(symptoms)-1].CreatedAt,
	}
}

// AggregateMedications calculates statistics for medications
func (da *DataAggregator) AggregateMedications(medications []*medication.Medication) *analytics.MedicationStats {
	if len(medications) == 0 {
		return &analytics.MedicationStats{
			ActiveMedications:  make([]string, 0),
			EffectivenessScore: make(map[string]float64),
		}
	}

	// Get active medications
	activeMeds := make([]string, 0)
	now := time.Now()
	for _, med := range medications {
		if med.EndDate == nil || med.EndDate.IsZero() || med.EndDate.After(now) {
			activeMeds = append(activeMeds, med.Name)
		}
	}
	stats := &analytics.MedicationStats{
		TotalCount:        len(medications),
		ActiveMedications: activeMeds,
	}

	stats.ComplianceScore = da.calculateComplianceScore(medications)
	stats.EffectivenessScore = da.calculateEffectivenessScores(medications)

	return stats
}

// Helper methods for score calculations
func (da *DataAggregator) calculateConsistencyScore(movements []bowelmovement.BowelMovement) float64 {
	// Implementation for calculating consistency score
	return 0.0
}

func (da *DataAggregator) calculateRegularityScore(movements []bowelmovement.BowelMovement) float64 {
	// Implementation for calculating regularity score
	return 0.0
}

func (da *DataAggregator) calculateNutritionScore(meals []meal.Meal) float64 {
	// Implementation for calculating nutrition score
	return 0.0
}

func (da *DataAggregator) calculateComplianceScore(medications []*medication.Medication) float64 {
	// Implementation for calculating medication compliance score
	return 0.0
}

func (da *DataAggregator) calculateEffectivenessScores(medications []*medication.Medication) map[string]float64 {
	// Implementation for calculating medication effectiveness scores
	return make(map[string]float64)
}
