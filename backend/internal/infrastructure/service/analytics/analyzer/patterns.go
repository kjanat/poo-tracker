package analyzer

import (
	"sort"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

// AnalyzeEatingPatterns identifies patterns in eating habits and their potential impacts
func (ta *TrendAnalyzer) AnalyzeEatingPatterns(meals []meal.Meal) *shared.EatingPattern {
	// TODO: Remove this early return once methods are implemented
	if len(meals) == 0 {
		return &shared.EatingPattern{}
	}

	pattern := &shared.EatingPattern{
		MealTimings:        []shared.MealTiming{}, // TODO: Implement analyzeMealTimings
		CommonIngredients:  []string{},            // TODO: Implement identifyCommonIngredients
		ProblemIngredients: []string{},            // TODO: Implement identifyProblemIngredients
	}
	return pattern
}

// AnalyzeBowelPatterns identifies patterns in bowel movements and their correlations with meals
func (ta *TrendAnalyzer) AnalyzeBowelPatterns(movements []bowelmovement.BowelMovement, meals []meal.Meal) *shared.BowelPattern {
	pattern := &shared.BowelPattern{
		Frequency:       ta.analyzeBowelFrequency(movements),
		Consistency:     ta.analyzeBowelConsistency(movements),
		MealCorrelation: ta.analyzeMealCorrelation(movements, meals),
	}
	return pattern
}

// AnalyzeSymptomPatterns identifies patterns in symptoms and their triggers
func (ta *TrendAnalyzer) AnalyzeSymptomPatterns(symptoms []symptom.Symptom) *shared.SymptomPattern {
	pattern := &shared.SymptomPattern{
		CommonSymptoms: ta.identifyCommonSymptomMap(symptoms),
		Frequency:      ta.analyzeSymptomFrequency(symptoms),
		Severity:       ta.analyzeSymptomSeverity(symptoms),
	}
	return pattern
}

// AnalyzeLifestylePatterns identifies patterns in overall lifestyle and health indicators
func (ta *TrendAnalyzer) AnalyzeLifestylePatterns(meals []meal.Meal, movements []bowelmovement.BowelMovement, symptoms []symptom.Symptom) *shared.LifestylePattern {
	pattern := &shared.LifestylePattern{
		DietaryHabits:   ta.analyzeDietaryHabits(meals),
		BowelRegularity: ta.analyzeBowelRegularity(movements),
		SymptomTriggers: ta.analyzeSymptomTriggers(symptoms, meals),
	}
	return pattern
}

// Helper methods for pattern analysis
func (ta *TrendAnalyzer) analyzeMealTimings(meals []meal.Meal) []shared.MealTiming {
	if len(meals) == 0 {
		return []shared.MealTiming{}
	}

	hourFreq := make(map[int]int)
	for _, m := range meals {
		hour := m.MealTime.Hour()
		hourFreq[hour]++
	}

	timings := make([]shared.MealTiming, 0, len(hourFreq))
	for h, freq := range hourFreq {
		timings = append(timings, shared.MealTiming{
			TimeOfDay: shared.NewTimeOfDay(h, 0),
			Frequency: freq,
		})
	}

	sort.Slice(timings, func(i, j int) bool {
		return timings[i].TimeOfDay.Hour < timings[j].TimeOfDay.Hour
	})

	return timings
}

func (ta *TrendAnalyzer) identifyCommonSymptomMap(symptoms []symptom.Symptom) map[string]int {
	// Convert symptom list to frequency map
	freq := make(map[string]int)
	for _, s := range symptoms {
		symptomType := s.Type.String()
		freq[symptomType]++
	}
	return freq
}
