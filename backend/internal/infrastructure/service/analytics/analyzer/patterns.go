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
	pattern := &shared.EatingPattern{}

	if len(meals) == 0 {
		return pattern
	}

	pattern.MealTimings = ta.analyzeMealTimings(meals)
	pattern.CommonIngredients = ta.identifyCommonIngredients(meals)
	pattern.ProblemIngredients = ta.identifyProblemIngredients(meals)

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

// identifyCommonIngredients returns a sorted list of frequently occurring meal
// attributes. The current model doesn't store detailed ingredient lists, so we
// approximate ingredients using boolean flags and characteristics such as
// Dairy, Gluten, FiberRich and SpicyLevel.
func (ta *TrendAnalyzer) identifyCommonIngredients(meals []meal.Meal) []string {
	if len(meals) == 0 {
		return []string{}
	}

	counts := make(map[string]int)
	for _, m := range meals {
		if m.Dairy {
			counts["dairy"]++
		}
		if m.Gluten {
			counts["gluten"]++
		}
		if m.FiberRich {
			counts["fiber"]++
		}
		if m.SpicyLevel != nil && *m.SpicyLevel > aggregator.SpicyThreshold {
			counts["spicy"]++
		}
	}

	type pair struct {
		name  string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, pair{name: k, count: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].name < pairs[j].name
		}
		return pairs[i].count > pairs[j].count
	})

	result := make([]string, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, p.name)
	}

	return result
}

// identifyProblemIngredients identifies ingredients that may cause issues for the user
// This analyzes meal properties (dairy, gluten, spicy level) and correlates with symptoms
func (ta *TrendAnalyzer) identifyProblemIngredients(meals []meal.Meal) []string {
	problemIngredients := make([]string, 0)

	// Create a map to track ingredient frequency and associated issues
	ingredientIssues := make(map[string]int)
	ingredientCount := make(map[string]int)

	// For each meal, extract triggers and check for issues in following days
	for i, meal := range meals {
		triggers := ta.extractMealTriggers(meal)

		// Look for symptoms/issues in next few meals
		hasIssues := ta.hasSubsequentIssues(meals, i)

		for _, trigger := range triggers {
			ingredientCount[trigger]++
			if hasIssues {
				ingredientIssues[trigger]++
			}
		}
	}

	// Identify ingredients with high issue rate (>50% correlation)
	for ingredient, issueCount := range ingredientIssues {
		totalCount := ingredientCount[ingredient]
		if totalCount >= 3 && float64(issueCount)/float64(totalCount) > 0.5 {
			problemIngredients = append(problemIngredients, ingredient)
		}
	}

	return problemIngredients
}

// extractMealTriggers extracts potential trigger ingredients from meal properties
func (ta *TrendAnalyzer) extractMealTriggers(meal meal.Meal) []string {
	triggers := make([]string, 0)

	if meal.Dairy {
		triggers = append(triggers, "dairy")
	}
	if meal.Gluten {
		triggers = append(triggers, "gluten")
	}
	if meal.SpicyLevel != nil && *meal.SpicyLevel > 6 {
		triggers = append(triggers, "spicy_food")
	}
	if meal.Calories > 800 {
		triggers = append(triggers, "high_calorie_meal")
	}
	if meal.Cuisine != "" {
		triggers = append(triggers, meal.Cuisine)
	}

	return triggers
}

// hasSubsequentIssues checks if there are digestive issues after this meal
func (ta *TrendAnalyzer) hasSubsequentIssues(meals []meal.Meal, currentIndex int) bool {
	if currentIndex >= len(meals)-1 {
		return false
	}

	currentMeal := meals[currentIndex]

	// Look at meals in the next 24-48 hours for digestive distress indicators
	// This is a simplified heuristic - in reality we'd correlate with actual symptom data
	for i := currentIndex + 1; i < len(meals) && i < currentIndex+3; i++ {
		nextMeal := meals[i]

		// Check if next meal is within reasonable timeframe (48 hours)
		if nextMeal.MealTime.Sub(currentMeal.MealTime) > 48*time.Hour {
			break
		}

		// Heuristic: if next meal is notably smaller or different, it might indicate issues
		// This is a simplified approach - ideally we'd have actual symptom correlation
		if ta.indicatesDigestiveIssues(currentMeal, nextMeal) {
			return true
		}
	}

	return false
}

// indicatesDigestiveIssues uses meal patterns to infer digestive issues
func (ta *TrendAnalyzer) indicatesDigestiveIssues(currentMeal, nextMeal meal.Meal) bool {
	// Simplified heuristic indicators:

	// 1. Significantly smaller next meal (possible appetite loss)
	if currentMeal.Calories > 0 && nextMeal.Calories > 0 {
		if float64(nextMeal.Calories) < float64(currentMeal.Calories)*0.5 {
			return true
		}
	}

	// 2. Large gap between meals (possible nausea/avoidance)
	timeGap := nextMeal.MealTime.Sub(currentMeal.MealTime)
	if timeGap > 8*time.Hour && currentMeal.MealTime.Hour() < 20 { // Not dinner
		return true
	}

	// 3. Switch to bland foods (avoiding triggers)
	if ta.isBlandMeal(nextMeal) && !ta.isBlandMeal(currentMeal) {
		return true
	}

	return false
}

// isBlandMeal determines if a meal is "bland" (likely chosen to avoid triggers)
func (ta *TrendAnalyzer) isBlandMeal(meal meal.Meal) bool {
	// Characteristics of bland meals:
	return !meal.Dairy &&
		!meal.Gluten &&
		(meal.SpicyLevel == nil || *meal.SpicyLevel <= 2) &&
		meal.Calories < 400
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
	if len(meals) == 0 {
		return nil
	}

	counts := map[string]int{
		"spicy meals":  0,
		"large meals":  0,
		"fiber rich":   0,
		"dairy meals":  0,
		"gluten meals": 0,
	}

	for _, m := range meals {
		if m.SpicyLevel != nil && *m.SpicyLevel > 0 {
			counts["spicy meals"]++
		}
		if m.Calories > 800 {
			counts["large meals"]++
		}
		if m.FiberRich {
			counts["fiber rich"]++
		}
		if m.Dairy {
			counts["dairy meals"]++
		}
		if m.Gluten {
			counts["gluten meals"]++
		}
	}

	habits := make([]shared.DietaryHabit, 0, len(counts))
	for desc, freq := range counts {
		if freq == 0 {
			continue
		}
		habits = append(habits, shared.DietaryHabit{
			Description: desc,
			Frequency:   freq,
			Impact:      0,
		})
	}

	sort.Slice(habits, func(i, j int) bool { return habits[i].Frequency > habits[j].Frequency })
	return habits
}

func (ta *TrendAnalyzer) analyzeBowelRegularity(movements []bowelmovement.BowelMovement) float64 {
	return ta.analyzeBowelFrequency(movements)
}

func (ta *TrendAnalyzer) analyzeSymptomTriggers(symptoms []symptom.Symptom, meals []meal.Meal) []shared.SymptomTrigger {
	if len(symptoms) == 0 || len(meals) == 0 {
		return nil
	}

	triggerCounts := make(map[string]int)
	for _, s := range symptoms {
		for _, m := range meals {
			if s.RecordedAt.After(m.MealTime) {
				if s.RecordedAt.Sub(m.MealTime) <= 6*time.Hour {
					triggerCounts[m.Name]++
				}
			}
		}
	}

	triggers := make([]shared.SymptomTrigger, 0, len(triggerCounts))
	for ing, count := range triggerCounts {
		confidence := float64(count) / float64(len(symptoms))
		triggers = append(triggers, shared.SymptomTrigger{
			TriggerType: "meal",
			Ingredient:  ing,
			Confidence:  confidence,
		})
	}

	sort.Slice(triggers, func(i, j int) bool { return triggers[i].Confidence > triggers[j].Confidence })
	return triggers
}
