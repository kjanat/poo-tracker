package analyzer

import (
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/aggregator"
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
	if symptoms == nil {
		return make(map[string]int)
	}

	// Convert symptom list to frequency map
	freq := make(map[string]int)
	for _, s := range symptoms {
		symptomType := s.Type.String()
		freq[symptomType]++
	}
	return freq
}

func (ta *TrendAnalyzer) analyzeBowelFrequency(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	days := make(map[string]struct{})
	for _, m := range movements {
		day := m.RecordedAt.Format("2006-01-02")
		days[day] = struct{}{}
	}

	if len(days) == 0 {
		return 0
	}

	return float64(len(movements)) / float64(len(days))
}

func (ta *TrendAnalyzer) analyzeBowelConsistency(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	total := 0
	for _, m := range movements {
		total += m.BristolType
	}

	return float64(total) / float64(len(movements))
}

func (ta *TrendAnalyzer) analyzeMealCorrelation(movements []bowelmovement.BowelMovement, meals []meal.Meal) float64 {
	if len(movements) == 0 || len(meals) == 0 {
		return 0
	}

	correlated := 0
	for _, bm := range movements {
		for _, m := range meals {
			diff := bm.RecordedAt.Sub(m.MealTime)
			if diff < 0 {
				diff = -diff
			}
			if diff <= 4*time.Hour {
				correlated++
				break
			}
		}
	}

	return float64(correlated) / float64(len(movements))
}

func (ta *TrendAnalyzer) analyzeSymptomFrequency(symptoms []symptom.Symptom) map[string]int {
	freq := make(map[string]int)
	for _, s := range symptoms {
		freq[s.Type.String()]++
	}
	return freq
}

func (ta *TrendAnalyzer) analyzeSymptomSeverity(symptoms []symptom.Symptom) map[string]float64 {
	totals := make(map[string]int)
	counts := make(map[string]int)

	for _, s := range symptoms {
		st := s.Type.String()
		totals[st] += s.Severity
		counts[st]++
	}

	result := make(map[string]float64)
	for t, total := range totals {
		result[t] = float64(total) / float64(counts[t])
	}

	return result
}

func (ta *TrendAnalyzer) analyzeDietaryHabits(meals []meal.Meal) []shared.DietaryHabit {
	habits := []shared.DietaryHabit{}
	if len(meals) == 0 {
		return habits
	}

	spicy, dairy, gluten, fiber := 0, 0, 0, 0
	for _, m := range meals {
		if m.SpicyLevel != nil && *m.SpicyLevel >= aggregator.SpicyThreshold {
			spicy++
		}
		if m.Dairy {
			dairy++
		}
		if m.Gluten {
			gluten++
		}
		if m.FiberRich {
			fiber++
		}
	}

	total := len(meals)
	if spicy > 0 {
		habits = append(habits, shared.DietaryHabit{
			Description: "Frequent spicy meals",
			Frequency:   spicy,
			Impact:      float64(spicy) / float64(total),
		})
	}
	if dairy > 0 {
		habits = append(habits, shared.DietaryHabit{
			Description: "Regular dairy consumption",
			Frequency:   dairy,
			Impact:      float64(dairy) / float64(total),
		})
	}
	if gluten > 0 {
		habits = append(habits, shared.DietaryHabit{
			Description: "Regular gluten consumption",
			Frequency:   gluten,
			Impact:      float64(gluten) / float64(total),
		})
	}
	if fiber > 0 {
		habits = append(habits, shared.DietaryHabit{
			Description: "Fiber rich meals",
			Frequency:   fiber,
			Impact:      float64(fiber) / float64(total),
		})
	}

	return habits
}

func (ta *TrendAnalyzer) analyzeBowelRegularity(movements []bowelmovement.BowelMovement) float64 {
	return ta.analyzeBowelFrequency(movements)
}

func (ta *TrendAnalyzer) analyzeSymptomTriggers(symptoms []symptom.Symptom, meals []meal.Meal) []shared.SymptomTrigger {
	triggers := []shared.SymptomTrigger{}
	if len(symptoms) == 0 || len(meals) == 0 {
		return triggers
	}

	counts := map[string]int{}
	total := 0
	for _, s := range symptoms {
		for _, m := range meals {
			diff := s.CreatedAt.Sub(m.MealTime)
			if diff < 0 {
				diff = -diff
			}
			if diff <= 4*time.Hour {
				if m.Dairy {
					counts["dairy"]++
				}
				if m.Gluten {
					counts["gluten"]++
				}
				if m.SpicyLevel != nil && *m.SpicyLevel >= aggregator.SpicyThreshold {
					counts["spicy"]++
				}
				total++
			}
		}
	}

	for k, c := range counts {
		triggers = append(triggers, shared.SymptomTrigger{
			TriggerType: "food",
			Ingredient:  k,
			Confidence:  float64(c) / float64(total),
		})
	}

	return triggers
}
