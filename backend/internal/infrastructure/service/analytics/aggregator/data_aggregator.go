package aggregator

import (
	"fmt"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

const (
	// SpicyThreshold is the minimum spicy level to consider a meal as spicy
	SpicyThreshold = 2

	// HealthyMaxCalories is the maximum calories for a meal to be considered healthy
	HealthyMaxCalories = 600

	// LargeMealCalories is the minimum calories to consider a meal as large
	LargeMealCalories = 800
)

// DataAggregator handles data aggregation operations
type DataAggregator struct{}

// NewDataAggregator creates a new data aggregator
func NewDataAggregator() *DataAggregator {
	return &DataAggregator{}
}

// AggregateDailyData aggregates all data types by day for comprehensive analysis
func (da *DataAggregator) AggregateDailyData(
	movements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
	medications []medication.Medication,
	start, end time.Time,
) []shared.DailyAggregation {
	dailyMap := make(map[string]*shared.DailyAggregation)

	// Initialize days in range
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		dayKey := d.Format("2006-01-02")
		dailyMap[dayKey] = &shared.DailyAggregation{
			Date:            d,
			BowelMovements:  []bowelmovement.BowelMovement{},
			Meals:           []meal.Meal{},
			Symptoms:        []symptom.Symptom{},
			Medications:     []medication.Medication{},
			BristolAverage:  0,
			PainAverage:     0,
			StrainAverage:   0,
			SatisfactionAvg: 0,
			MealCount:       0,
			SymptomCount:    0,
			SpicyMealCount:  0,
			DairyMealCount:  0,
			GlutenMealCount: 0,
		}
	}

	// Aggregate bowel movements
	for _, movement := range movements {
		dayKey := movement.RecordedAt.Format("2006-01-02")
		if data, exists := dailyMap[dayKey]; exists {
			data.BowelMovements = append(data.BowelMovements, movement)
		}
	}

	// Aggregate meals
	for _, meal := range meals {
		dayKey := meal.MealTime.Format("2006-01-02")
		if data, exists := dailyMap[dayKey]; exists {
			data.Meals = append(data.Meals, meal)
			data.MealCount++

			// Count special meal types
			if meal.SpicyLevel != nil && *meal.SpicyLevel > 0 {
				data.SpicyMealCount++
			}

			// Check for dairy and gluten using boolean fields
			if meal.Dairy {
				data.DairyMealCount++
			}
			if meal.Gluten {
				data.GlutenMealCount++
			}
		}
	}

	// Aggregate symptoms
	for _, symptom := range symptoms {
		dayKey := symptom.RecordedAt.Format("2006-01-02")
		if data, exists := dailyMap[dayKey]; exists {
			data.Symptoms = append(data.Symptoms, symptom)
			data.SymptomCount++
		}
	}

	// Aggregate medications
	for _, medication := range medications {
		if medication.StartDate != nil {
			dayKey := medication.StartDate.Format("2006-01-02")
			if data, exists := dailyMap[dayKey]; exists {
				data.Medications = append(data.Medications, medication)
			}
		}
	}

	// Calculate daily averages
	for _, data := range dailyMap {
		if len(data.BowelMovements) > 0 {
			bristolSum, painSum, strainSum, satisfactionSum := 0.0, 0.0, 0.0, 0.0
			for _, bm := range data.BowelMovements {
				bristolSum += float64(bm.BristolType)
				painSum += float64(bm.Pain)
				strainSum += float64(bm.Strain)
				satisfactionSum += float64(bm.Satisfaction)
			}
			count := float64(len(data.BowelMovements))
			data.BristolAverage = bristolSum / count
			data.PainAverage = painSum / count
			data.StrainAverage = strainSum / count
			data.SatisfactionAvg = satisfactionSum / count
		}
	}

	// Convert map to sorted slice
	var result []shared.DailyAggregation
	for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
		dayKey := d.Format("2006-01-02")
		if data, exists := dailyMap[dayKey]; exists {
			result = append(result, *data)
		}
	}

	return result
}

// CalculateAverage calculates the average of a slice of float64 values
func (da *DataAggregator) CalculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// CalculateAverageSymptomSeverity calculates average symptom severity
func (da *DataAggregator) CalculateAverageSymptomSeverity(symptoms []symptom.Symptom) float64 {
	if len(symptoms) == 0 {
		return 0
	}

	total := 0.0
	for _, s := range symptoms {
		total += float64(s.Severity)
	}
	return total / float64(len(symptoms))
}

// CalculateAveragePain calculates average pain from bowel movements
func (da *DataAggregator) CalculateAveragePain(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	total := 0.0
	for _, m := range movements {
		total += float64(m.Pain)
	}
	return total / float64(len(movements))
}

// CalculateAverageSatisfaction calculates average satisfaction from bowel movements
func (da *DataAggregator) CalculateAverageSatisfaction(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	total := 0.0
	for _, m := range movements {
		total += float64(m.Satisfaction)
	}
	return total / float64(len(movements))
}

// GroupBowelMovementsByWeek groups bowel movements by week for trend analysis
func (da *DataAggregator) GroupBowelMovementsByWeek(movements []bowelmovement.BowelMovement) map[string][]bowelmovement.BowelMovement {
	weekGroups := make(map[string][]bowelmovement.BowelMovement)

	for _, movement := range movements {
		// Get Monday of the week as the key
		year, week := movement.RecordedAt.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		weekGroups[weekKey] = append(weekGroups[weekKey], movement)
	}

	return weekGroups
}

// GroupMealsByType groups meals by their characteristics for analysis
func (da *DataAggregator) GroupMealsByType(meals []meal.Meal) map[string][]meal.Meal {
	groups := map[string][]meal.Meal{
		"spicy":   {},
		"dairy":   {},
		"gluten":  {},
		"healthy": {},
		"large":   {},
	}

	for _, meal := range meals {
		if meal.SpicyLevel != nil && *meal.SpicyLevel > SpicyThreshold {
			groups["spicy"] = append(groups["spicy"], meal)
		}

		// Check boolean fields for dairy/gluten
		if meal.Dairy {
			groups["dairy"] = append(groups["dairy"], meal)
		}
		if meal.Gluten {
			groups["gluten"] = append(groups["gluten"], meal)
		}

		// Classify as healthy based on simple criteria
		if meal.Calories > 0 && meal.Calories < HealthyMaxCalories && meal.FiberRich {
			groups["healthy"] = append(groups["healthy"], meal)
		}

		// Classify as large meal
		if meal.Calories > LargeMealCalories {
			groups["large"] = append(groups["large"], meal)
		}
	}

	return groups
}

// GetActiveMedicationPercentage returns the percentage of medications that are currently active.
//
// This function is a simplified metric and does not track individual dose adherence.
// It merely reflects how many medications are marked active versus the total count.
func (da *DataAggregator) GetActiveMedicationPercentage(medications []medication.Medication, days int) float64 {
	if len(medications) == 0 || days == 0 {
		return 0
	}

	// Since the medication model doesn't have dose tracking fields,
	// this calculation is a rough estimate based solely on the number of
	// medications marked as active.
	totalMedications := len(medications)
	activeMedications := 0

	for _, med := range medications {
		if med.IsActive {
			activeMedications++
		}
	}

	if totalMedications == 0 {
		return 0
	}

	return float64(activeMedications) / float64(totalMedications) * 100
}
