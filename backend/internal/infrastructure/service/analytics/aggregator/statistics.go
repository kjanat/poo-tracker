package aggregator

import (
	"math"
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
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

	// Collect statistics
	bristolFreq := make(map[int]int)
	var painSum, strainSum, satisfactionSum float64
	for _, m := range movements {
		bristolFreq[m.BristolType]++
		painSum += float64(m.Pain)
		strainSum += float64(m.Strain)
		satisfactionSum += float64(m.Satisfaction)
	}

	mostBristol := 0
	maxCount := 0
	for b, count := range bristolFreq {
		if count > maxCount {
			maxCount = count
			mostBristol = b
		}
	}

	return &analytics.BowelMovementStats{
		TotalCount:          int64(len(movements)),
		AveragePerDay:       float64(len(movements)) / totalDays,
		MostCommonBristol:   mostBristol,
		AveragePain:         painSum / float64(len(movements)),
		AverageStrain:       strainSum / float64(len(movements)),
		AverageSatisfaction: satisfactionSum / float64(len(movements)),
		RegularityScore:     da.calculateRegularityScore(movements),
	}
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

	var calorieSum int
	fiberRich := 0
	for _, m := range meals {
		calorieSum += m.Calories
		if m.FiberRich {
			fiberRich++
		}
	}

	return &analytics.MealStats{
		TotalMeals:       int64(len(meals)),
		AveragePerDay:    float64(len(meals)) / totalDays,
		TotalCalories:    calorieSum,
		AverageCalories:  float64(calorieSum) / float64(len(meals)),
		FiberRichPercent: float64(fiberRich) / float64(len(meals)) * 100,
		HealthScore:      da.calculateNutritionScore(meals),
	}
}

// AggregateSymptoms calculates statistics for symptoms
func (da *DataAggregator) AggregateSymptoms(symptoms []symptom.Symptom) *analytics.SymptomStats {
	if len(symptoms) == 0 {
		return &analytics.SymptomStats{}
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

	var severityTotal float64
	catFreq := make(map[string]int)
	typeFreq := make(map[string]int)
	for _, s := range symptoms {
		severityTotal += float64(s.Severity)
		if s.Category != nil {
			catFreq[s.Category.String()]++
		}
		if s.Type != nil {
			typeFreq[s.Type.String()]++
		}
	}

	mostCat := ""
	mostCatCount := 0
	for c, count := range catFreq {
		if count > mostCatCount {
			mostCat = c
			mostCatCount = count
		}
	}

	mostType := ""
	mostTypeCount := 0
	for t, count := range typeFreq {
		if count > mostTypeCount {
			mostType = t
			mostTypeCount = count
		}
	}

	return &analytics.SymptomStats{
		TotalSymptoms:      int64(len(symptoms)),
		AveragePerDay:      float64(len(symptoms)) / totalDays,
		AverageSeverity:    severityTotal / float64(len(symptoms)),
		MostCommonCategory: mostCat,
		MostCommonType:     mostType,
		TrendDirection:     "stable",
	}
}

// AggregateMedications calculates statistics for medications
func (da *DataAggregator) AggregateMedications(medications []*medication.Medication) *analytics.MedicationStats {
	if len(medications) == 0 {
		return &analytics.MedicationStats{}
	}

	// Get active medications
	activeMeds := make([]string, 0)
	now := time.Now()
	for _, med := range medications {
		if med.EndDate == nil || med.EndDate.IsZero() || med.EndDate.After(now) {
			activeMeds = append(activeMeds, med.Name)
		}
	}
	categoryFreq := make(map[string]int)
	for _, med := range medications {
		if med.Category != nil {
			categoryFreq[med.Category.String()]++
		}
	}
	mostCat := ""
	mostCount := 0
	for c, count := range categoryFreq {
		if count > mostCount {
			mostCat = c
			mostCount = count
		}
	}

	return &analytics.MedicationStats{
		TotalMedications:   int64(len(medications)),
		ActiveMedications:  int64(len(activeMeds)),
		AdherenceScore:     da.calculateComplianceScore(medications),
		MostCommonCategory: mostCat,
		ComplexityScore:    da.calculateEffectivenessScores(medications)["overall"],
	}
}

// Helper methods for score calculations
func (da *DataAggregator) calculateConsistencyScore(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	sort.Slice(movements, func(i, j int) bool { return movements[i].RecordedAt.Before(movements[j].RecordedAt) })

	freq := make(map[int]int)
	changes := 0
	for i, m := range movements {
		freq[m.BristolType]++
		if i > 0 && m.BristolType != movements[i-1].BristolType {
			changes++
		}
	}

	healthyCount := 0
	for _, t := range []int{3, 4, 5} {
		healthyCount += freq[t]
	}

	healthRatio := float64(healthyCount) / float64(len(movements))
	changeRatio := 0.0
	if len(movements) > 1 {
		changeRatio = float64(changes) / float64(len(movements)-1)
	}

	score := healthRatio * (1 - changeRatio)
	if score < 0 {
		score = 0
	} else if score > 1 {
		score = 1
	}

	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateRegularityScore(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) < 2 {
		return 0
	}

	sort.Slice(movements, func(i, j int) bool { return movements[i].RecordedAt.Before(movements[j].RecordedAt) })

	var totalVariance float64
	last := movements[0].RecordedAt
	for i := 1; i < len(movements); i++ {
		interval := movements[i].RecordedAt.Sub(last).Hours()
		diff := interval - 24
		totalVariance += diff * diff
		last = movements[i].RecordedAt
	}

	avgVariance := totalVariance / float64(len(movements)-1)
	maxVariance := 24.0
	score := 1 - (avgVariance / (maxVariance * maxVariance))
	if score < 0 {
		score = 0
	} else if score > 1 {
		score = 1
	}
	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateNutritionScore(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}

	totalCalories := 0
	fiberRich := 0
	for _, m := range meals {
		totalCalories += m.Calories
		if m.FiberRich {
			fiberRich++
		}
	}

	avgCalories := float64(totalCalories) / float64(len(meals))
	ideal := 600.0
	calorieScore := 1 - math.Abs(avgCalories-ideal)/ideal
	if calorieScore < 0 {
		calorieScore = 0
	}
	fiberScore := float64(fiberRich) / float64(len(meals))
	score := 0.7*calorieScore + 0.3*fiberScore
	if score < 0 {
		score = 0
	} else if score > 1 {
		score = 1
	}
	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateComplianceScore(medications []*medication.Medication) float64 {
	if len(medications) == 0 {
		return 0
	}

	now := time.Now()
	adhered := 0
	for _, med := range medications {
		if med == nil || !med.IsActive {
			continue
		}
		if med.TakenAt == nil || now.Sub(*med.TakenAt).Hours() <= 36 {
			adhered++
		}
	}

	score := float64(adhered) / float64(len(medications))
	if score < 0 {
		score = 0
	} else if score > 1 {
		score = 1
	}

	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateEffectivenessScores(medications []*medication.Medication) map[string]float64 {
	scores := make(map[string]float64)
	if len(medications) == 0 {
		return scores
	}

	now := time.Now()
	total := 0.0
	count := 0
	for _, med := range medications {
		if med == nil {
			continue
		}
		score := 0.0
		if med.IsActive {
			score = 0.7
			if med.TakenAt != nil && now.Sub(*med.TakenAt).Hours() <= 36 {
				score = 1.0
			}
		}
		scores[med.Name] = shared.RoundToDecimalPlaces(score, 2)
		total += score
		count++
	}

	if count > 0 {
		scores["overall"] = shared.RoundToDecimalPlaces(total/float64(count), 2)
	}

	return scores
}
