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

	var diffSum float64
	for _, m := range movements {
		diff := math.Abs(float64(m.BristolType) - 4)
		diffSum += diff
	}

	avgDiff := diffSum / float64(len(movements))
	score := 1 - (avgDiff / 3)
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateRegularityScore(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) < 2 {
		return 0
	}

	sort.Slice(movements, func(i, j int) bool { return movements[i].RecordedAt.Before(movements[j].RecordedAt) })

	intervals := make([]float64, 0, len(movements)-1)
	for i := 1; i < len(movements); i++ {
		hours := movements[i].RecordedAt.Sub(movements[i-1].RecordedAt).Hours()
		intervals = append(intervals, hours)
	}

	stats := shared.CalculateStatistics(intervals)
	if stats.Mean == 0 {
		return 0
	}

	score := 1 - (stats.StdDev / stats.Mean)
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateNutritionScore(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}

	fiberCount := 0
	healthyCalories := 0
	for _, m := range meals {
		if m.FiberRich {
			fiberCount++
		}
		if m.Calories <= HealthyMaxCalories {
			healthyCalories++
		}
	}

	fiberRatio := float64(fiberCount) / float64(len(meals))
	healthyRatio := float64(healthyCalories) / float64(len(meals))

	score := (fiberRatio*0.6 + healthyRatio*0.4) * 100
	return shared.RoundToDecimalPlaces(score, 2)
}

func (da *DataAggregator) calculateComplianceScore(medications []*medication.Medication) float64 {
	if len(medications) == 0 {
		return 0
	}

	active := 0
	for _, m := range medications {
		if m.IsActive {
			active++
		}
	}
	return shared.RoundToDecimalPlaces(float64(active)/float64(len(medications))*100, 2)
}

func (da *DataAggregator) calculateEffectivenessScores(medications []*medication.Medication) map[string]float64 {
	scores := make(map[string]float64)
	if len(medications) == 0 {
		return scores
	}

	active := 0
	for _, m := range medications {
		if m.IsActive {
			active++
		}
	}
	scores["overall"] = shared.RoundToDecimalPlaces(float64(active)/float64(len(medications))*100, 2)
	return scores
}
