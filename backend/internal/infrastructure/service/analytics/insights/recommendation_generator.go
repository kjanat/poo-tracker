//go:build ignore

package insights

import (
	"fmt"
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

// GenerateRecommendations generates recommendations based on health data
func (ie *InsightEngine) GenerateRecommendations(
	bowelValues []bowelmovement.BowelMovement,
	mealValues []meal.Meal,
	symptomValues []symptom.Symptom,
) []*shared.InsightRecommendation {
	var recommendations []*shared.InsightRecommendation

	// Add recommendations based on bowel health
	bowelRecs := ie.analyzeBowelHealth(bowelValues)
	recommendations = append(recommendations, bowelRecs...)

	// Add recommendations based on diet
	dietRecs := ie.analyzeDiet(mealValues)
	recommendations = append(recommendations, dietRecs...)

	// Add recommendations based on symptoms
	symptomRecs := ie.analyzeSymptoms(symptomValues)
	recommendations = append(recommendations, symptomRecs...)

	// Create correlation-based recommendations
	correlationRecs := ie.analyzeCorrelations(bowelValues, mealValues, symptomValues)
	recommendations = append(recommendations, correlationRecs...)

	// Prioritize recommendations
	ie.prioritizeRecommendations(recommendations)

	// Return top recommendations
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations
}

// Helper methods
func (ie *InsightEngine) analyzeBowelHealth(movements []bowelmovement.BowelMovement) []*shared.InsightRecommendation {
	var recommendations []*shared.InsightRecommendation

	if len(movements) == 0 {
		now := time.Now()
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "bowel_tracking_start",
			Type:        "TRACKING",
			Priority:    "HIGH",
			Title:       "Track Bowel Movements",
			Description: "Start tracking your bowel movements to get personalized recommendations",
			Evidence:    []string{},
			Actions:     []string{"Open tracking form", "Add first entry"},
			Context:     make(map[string]any),
			CreatedAt:   now,
		})
		return recommendations
	}

	if len(movements) < 3 {
		return recommendations
	}

	painSum := 0.0
	satisfactionSum := 0.0

	for _, bm := range movements {
		painSum += float64(bm.Pain)
		satisfactionSum += float64(bm.Satisfaction)
	}

	avgPain := painSum / float64(len(movements))
	avgSatisfaction := satisfactionSum / float64(len(movements))

	now := time.Now()

	if avgPain > 6 {
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "bowel_pain_high",
			Type:        "LIFESTYLE",
			Priority:    "MEDIUM",
			Title:       "High Bowel Movement Pain",
			Description: fmt.Sprintf("Average pain level %.1f detected in recent bowel movements", avgPain),
			Evidence:    []string{fmt.Sprintf("Average pain: %.1f", avgPain)},
			Actions: []string{
				"Review diet for potential triggers",
				"Discuss pain management with a professional",
			},
			Context:   make(map[string]any),
			CreatedAt: now,
		})
	}

	if avgSatisfaction < 5 {
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "bowel_satisfaction_low",
			Type:        "LIFESTYLE",
			Priority:    "MEDIUM",
			Title:       "Low Bowel Movement Satisfaction",
			Description: fmt.Sprintf("Average satisfaction %.1f suggests discomfort", avgSatisfaction),
			Evidence:    []string{fmt.Sprintf("Average satisfaction: %.1f", avgSatisfaction)},
			Actions: []string{
				"Increase fiber and hydration",
				"Track meals leading to low satisfaction",
			},
			Context:   make(map[string]any),
			CreatedAt: now,
		})
	}

	return recommendations
}

func (ie *InsightEngine) analyzeDiet(meals []meal.Meal) []*shared.InsightRecommendation {
	var recommendations []*shared.InsightRecommendation

	if len(meals) == 0 {
		now := time.Now()
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "meal_tracking_start",
			Type:        "TRACKING",
			Priority:    "HIGH",
			Title:       "Track Meals",
			Description: "Start tracking your meals to get personalized dietary recommendations",
			Evidence:    []string{},
			Actions:     []string{"Open meal form", "Add first meal"},
			Context:     make(map[string]any),
			CreatedAt:   now,
		})
		return recommendations
	}

	if len(meals) < 3 {
		return recommendations
	}

	fiberCount := 0
	calorieSum := 0
	for _, m := range meals {
		if m.FiberRich {
			fiberCount++
		}
		calorieSum += m.Calories
	}

	fiberRatio := float64(fiberCount) / float64(len(meals))
	avgCalories := float64(calorieSum) / float64(len(meals))
	now := time.Now()

	if fiberRatio < 0.3 {
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "low_fiber_intake",
			Type:        "DIET",
			Priority:    "MEDIUM",
			Title:       "Increase Fiber Intake",
			Description: "Your recent meals appear low in fiber which can impact digestion",
			Evidence:    []string{fmt.Sprintf("Only %.0f%% of meals were fiber rich", fiberRatio*100)},
			Actions: []string{
				"Add more whole grains, fruits and vegetables",
			},
			Context:   make(map[string]any),
			CreatedAt: now,
		})
	}

	if avgCalories > 800 {
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "high_calorie_meals",
			Type:        "DIET",
			Priority:    "LOW",
			Title:       "High Average Meal Calories",
			Description: fmt.Sprintf("Average calories per meal is %.0f which may be high", avgCalories),
			Evidence:    []string{fmt.Sprintf("Average calories: %.0f", avgCalories)},
			Actions: []string{
				"Consider reducing portion sizes",
				"Balance meals with vegetables",
			},
			Context:   make(map[string]any),
			CreatedAt: now,
		})
	}

	return recommendations
}

func (ie *InsightEngine) analyzeSymptoms(symptoms []symptom.Symptom) []*shared.InsightRecommendation {
	var recommendations []*shared.InsightRecommendation

	if len(symptoms) == 0 {
		now := time.Now()
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "symptom_tracking_start",
			Type:        "TRACKING",
			Priority:    "HIGH",
			Title:       "Track Symptoms",
			Description: "Start tracking your symptoms to get personalized recommendations",
			Evidence:    []string{},
			Actions:     []string{"Open symptom form", "Add first symptom"},
			Context:     make(map[string]any),
			CreatedAt:   now,
		})
		return recommendations
	}

	severitySum := 0.0
	freq := make(map[string]int)
	for _, s := range symptoms {
		severitySum += float64(s.Severity)
		if s.Type != nil {
			freq[string(*s.Type)]++
		}
	}

	avgSeverity := severitySum / float64(len(symptoms))
	common := ""
	max := 0
	for t, c := range freq {
		if c > max {
			common = t
			max = c
		}
	}

	now := time.Now()

	if avgSeverity > 6 {
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "symptom_severity_high",
			Type:        "SYMPTOM",
			Priority:    "MEDIUM",
			Title:       "High Symptom Severity",
			Description: fmt.Sprintf("Average symptom severity is %.1f", avgSeverity),
			Evidence:    []string{fmt.Sprintf("Common symptom: %s", common)},
			Actions: []string{
				"Review potential triggers",
				"Consult a healthcare professional",
			},
			Context:   make(map[string]any),
			CreatedAt: now,
		})
	}

	return recommendations
}

func (ie *InsightEngine) analyzeCorrelations(
	movements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) []*shared.InsightRecommendation {
	var recommendations []*shared.InsightRecommendation

	if len(movements) == 0 || len(meals) == 0 {
		return recommendations
	}

	if len(movements) < 3 || len(meals) < 3 {
		return recommendations
	}

	// Simple correlation between daily spicy meal count and pain level
	dailySpice := make(map[string]float64)
	dailyPain := make(map[string]float64)
	dailyCount := make(map[string]int)
	dailyMovementCount := make(map[string]int)

	for _, m := range meals {
		if m.SpicyLevel != nil {
			key := m.MealTime.Format("2006-01-02")
			dailySpice[key] += float64(*m.SpicyLevel)
			dailyCount[key]++
		}
	}

	for _, bm := range movements {
		key := bm.RecordedAt.Format("2006-01-02")
		dailyPain[key] += float64(bm.Pain)
		dailyMovementCount[key]++
	}

	var spiceVals, painVals []float64
	for day, spice := range dailySpice {
		if c, ok := dailyCount[day]; ok && c > 0 {
			avgSpice := spice / float64(c)
			if pain, exists := dailyPain[day]; exists {
				if mc, ok := dailyMovementCount[day]; ok && mc > 0 {
					avgPain := pain / float64(mc)
					spiceVals = append(spiceVals, avgSpice)
					painVals = append(painVals, avgPain)
				}
			}
		}
	}

	if len(spiceVals) < 3 {
		return recommendations
	}

	coeff := shared.CalculateCorrelation(spiceVals, painVals)
	if coeff > 0.5 {
		now := time.Now()
		recommendations = append(recommendations, &shared.InsightRecommendation{
			ID:          "spice_pain_correlation",
			Type:        "CORRELATION",
			Priority:    "MEDIUM",
			Title:       "Spicy Food Linked to Pain",
			Description: "Higher spicy food levels correlate with increased bowel movement pain",
			Evidence:    []string{fmt.Sprintf("Correlation coefficient %.2f", coeff)},
			Actions: []string{
				"Reduce spicy food intake",
				"Monitor pain levels when avoiding spicy meals",
			},
			Context:   make(map[string]any),
			CreatedAt: now,
		})
	}

	return recommendations
}

func (ie *InsightEngine) prioritizeRecommendations(recommendations []*shared.InsightRecommendation) {
	priorityValue := func(p string) int {
		switch p {
		case "HIGH":
			return 1
		case "MEDIUM":
			return 2
		default:
			return 3
		}
	}

	sort.SliceStable(recommendations, func(i, j int) bool {
		return priorityValue(recommendations[i].Priority) < priorityValue(recommendations[j].Priority)
	})
}
