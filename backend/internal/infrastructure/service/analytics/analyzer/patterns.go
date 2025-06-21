package analyzer

import (
	"math"
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

const (
	minSpicyLevel        = 7
	highCalorieThreshold = 700
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
	// Helper methods currently never fail, but this structure allows future
	// implementations to surface errors without changing callers.

	pattern := &shared.LifestylePattern{}

	pattern.DietaryHabits = ta.analyzeDietaryHabits(meals)
	pattern.BowelRegularity = ta.analyzeBowelRegularity(movements)
	pattern.SymptomTriggers = ta.analyzeSymptomTriggers(symptoms, meals)

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
	if len(meals) == 0 {
		return []shared.DietaryHabit{}
	}

	type habitData struct {
		desc  string
		count int
	}

	habits := map[string]*habitData{
		"fiber":   {desc: "Fiber-rich meals"},
		"dairy":   {desc: "Dairy consumption"},
		"gluten":  {desc: "Gluten consumption"},
		"spicy":   {desc: "Spicy meals"},
		"calorie": {desc: "High calorie meals"},
	}

	for _, m := range meals {
		if m.FiberRich {
			habits["fiber"].count++
		}
		if m.Dairy {
			habits["dairy"].count++
		}
		if m.Gluten {
			habits["gluten"].count++
		}
		if m.SpicyLevel != nil && *m.SpicyLevel >= minSpicyLevel {
			habits["spicy"].count++
		}
		if m.Calories > highCalorieThreshold {
			habits["calorie"].count++
		}
	}

	results := make([]shared.DietaryHabit, 0, len(habits))
	for _, h := range habits {
		if h.count == 0 {
			continue
		}
		results = append(results, shared.DietaryHabit{
			Description: h.desc,
			Frequency:   h.count,
			Impact:      float64(h.count) / float64(len(meals)),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Frequency > results[j].Frequency
	})

	return results
}

func (ta *TrendAnalyzer) analyzeBowelRegularity(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	freq := ta.analyzeBowelFrequency(movements)

	times := make([]time.Time, 0, len(movements))
	for _, m := range movements {
		times = append(times, m.RecordedAt)
	}
	groups := shared.GroupByDay(times)
	if len(groups) <= 1 {
		// With data from a single day we can't measure variability.
		// Assume perfect regularity.
		return 1
	}

	counts := make([]float64, 0, len(groups))
	for _, ts := range groups {
		counts = append(counts, float64(len(ts)))
	}

	var sum float64
	for _, c := range counts {
		diff := c - freq
		sum += diff * diff
	}
	stdDev := math.Sqrt(sum / float64(len(counts)))

	if stdDev == 0 {
		return 1
	}

	return 1 / (1 + stdDev)
}

func (ta *TrendAnalyzer) analyzeSymptomTriggers(symptoms []symptom.Symptom, meals []meal.Meal) []shared.SymptomTrigger {
	if len(symptoms) == 0 || len(meals) == 0 {
		return []shared.SymptomTrigger{}
	}

	type counter struct {
		total     int
		triggered int
	}

	trig := map[string]*counter{}
	window := 6 * time.Hour

	for _, m := range meals {
		keys := []string{}
		if m.Dairy {
			keys = append(keys, "Dairy")
		}
		if m.Gluten {
			keys = append(keys, "Gluten")
		}
		if m.SpicyLevel != nil && *m.SpicyLevel >= minSpicyLevel {
			keys = append(keys, "Spicy")
		}

		if len(keys) == 0 {
			continue
		}

		for _, k := range keys {
			c, ok := trig[k]
			if !ok {
				c = &counter{}
				trig[k] = c
			}
			c.total++

			for _, s := range symptoms {
				if s.RecordedAt.After(m.MealTime) && s.RecordedAt.Sub(m.MealTime) <= window {
					c.triggered++
					break
				}
			}
		}
	}

	result := make([]shared.SymptomTrigger, 0, len(trig))
	for k, c := range trig {
		if c.total == 0 {
			continue
		}
		confidence := shared.CalculateConfidenceScore(c.total)
		rate := float64(c.triggered) / float64(c.total)
		result = append(result, shared.SymptomTrigger{
			TriggerType: k,
			Confidence:  rate * confidence,
		})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Confidence > result[j].Confidence })
	return result
}
