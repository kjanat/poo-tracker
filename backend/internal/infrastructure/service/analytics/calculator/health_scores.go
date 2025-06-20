package calculator

import (
	"math"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

// HealthScoreCalculator calculates various health scores
type HealthScoreCalculator struct{}

// NewHealthScoreCalculator creates a new health score calculator
func NewHealthScoreCalculator() *HealthScoreCalculator {
	return &HealthScoreCalculator{}
}

// HealthDataSummarizer calculates summary statistics for health data
type HealthDataSummarizer struct{}

// NewHealthDataSummarizer creates a new health data summarizer
func NewHealthDataSummarizer() *HealthDataSummarizer {
	return &HealthDataSummarizer{}
}

// CalculateOverallHealthScore calculates an overall health score (0-100)
func (hsc *HealthScoreCalculator) CalculateOverallHealthScore(
	bowelMovementStats *analytics.BowelMovementSummary,
	mealStats *analytics.MealSummary,
	symptomStats *analytics.SymptomSummary,
	medicationStats *analytics.MedicationSummary,
) float64 {
	// Weight factors for different components
	const (
		bowelWeight      = 0.4 // 40% - most important for digestive health
		symptomWeight    = 0.3 // 30% - symptom management
		mealWeight       = 0.2 // 20% - nutrition
		medicationWeight = 0.1 // 10% - medication adherence
	)

	// Calculate individual scores (0-100)
	bowelScore := hsc.calculateBowelMovementScore(bowelMovementStats)
	symptomScore := hsc.calculateSymptomScore(symptomStats)
	mealScore := hsc.calculateMealScore(mealStats)
	medicationScore := hsc.calculateMedicationScore(medicationStats)

	// Calculate weighted average
	overallScore := (bowelScore*bowelWeight +
		symptomScore*symptomWeight +
		mealScore*mealWeight +
		medicationScore*medicationWeight)

	return shared.RoundToDecimalPlaces(shared.SanitizeFloat64(overallScore), 1)
}

// calculateBowelMovementScore calculates health score based on bowel movements (0-100)
func (hsc *HealthScoreCalculator) calculateBowelMovementScore(stats *analytics.BowelMovementSummary) float64 {
	const (
		idealBristolValue           = 3.5
		bristolPenaltyMultiplier    = 10.0
		painPenaltyMultiplier       = 5.0
		strainPenaltyMultiplier     = 5.0
		satisfactionBonusMultiplier = 2.0
		lowFrequencyThreshold       = 0.5
		highFrequencyThreshold      = 4.0
		lowFrequencyPenalty         = 20.0
		highFrequencyPenalty        = 15.0
	)

	if stats == nil || stats.TotalCount == 0 {
		return 50.0 // Neutral score if no data
	}

	score := 100.0

	// Bristol type score (ideal range 3-4)
	bristolDeviation := math.Abs(float64(stats.MostCommonBristol) - idealBristolValue)
	bristolPenalty := bristolDeviation * bristolPenaltyMultiplier // Penalty increases with deviation
	score -= bristolPenalty

	// Pain score (lower is better, 1-10 scale)
	painPenalty := (stats.AveragePain - 1) * painPenaltyMultiplier // Scale 1-10 to 0-45 penalty
	score -= painPenalty

	// Strain score (lower is better, 1-10 scale)
	strainPenalty := (stats.AverageStrain - 1) * strainPenaltyMultiplier // Scale 1-10 to 0-45 penalty
	score -= strainPenalty

	// Satisfaction bonus (higher is better, 1-10 scale)
	satisfactionBonus := (stats.AverageSatisfaction - 5) * satisfactionBonusMultiplier // Scale around middle
	score += satisfactionBonus

	// Frequency score (use average per day)
	if stats.AveragePerDay < lowFrequencyThreshold {
		score -= lowFrequencyPenalty // Too infrequent
	} else if stats.AveragePerDay > highFrequencyThreshold {
		score -= highFrequencyPenalty // Too frequent
	}

	return math.Max(0, math.Min(100, score))
}

// calculateSymptomScore calculates health score based on symptoms (0-100)
func (hsc *HealthScoreCalculator) calculateSymptomScore(stats *analytics.SymptomSummary) float64 {
	if stats == nil {
		return 100.0 // No symptoms = perfect score
	}

	if stats.TotalSymptoms == 0 {
		return 100.0 // No symptoms = perfect score
	}

	score := 100.0

	// Penalty based on symptom count
	countPenalty := float64(stats.TotalSymptoms) * 2 // 2 points per symptom
	score -= countPenalty

	// Penalty based on severity (1-10 scale)
	severityPenalty := stats.AverageSeverity * 5 // Scale severity to penalty
	score -= severityPenalty

	// Additional penalty for high symptom frequency
	if stats.TotalSymptoms > 10 {
		score -= 20 // High symptom burden
	}

	return math.Max(0, math.Min(100, score))
}

// calculateMealScore calculates health score based on meals (0-100)
func (hsc *HealthScoreCalculator) calculateMealScore(stats *analytics.MealSummary) float64 {
	if stats == nil || stats.TotalMeals == 0 {
		return 50.0 // Neutral score if no data
	}

	score := 100.0

	// Penalty for excessive calories
	if stats.AverageCalories > 2500 {
		score -= 15
	} else if stats.AverageCalories > 2000 {
		score -= 5
	}

	// Penalty for too few calories
	if stats.AverageCalories < 1200 {
		score -= 20
	}

	// Bonus for balanced meal frequency (3-4 meals per day ideal)
	if stats.AveragePerDay >= 2.5 && stats.AveragePerDay <= 4.5 {
		score += 10 // Good meal frequency
	} else if stats.AveragePerDay < 2 {
		score -= 15 // Too few meals
	} else if stats.AveragePerDay > 6 {
		score -= 10 // Too many meals
	}

	// Additional scoring based on meal diversity would go here
	// (requires more detailed meal analysis)

	return math.Max(0, math.Min(100, score))
}

// calculateMedicationScore calculates health score based on medication adherence (0-100)
func (hsc *HealthScoreCalculator) calculateMedicationScore(stats *analytics.MedicationSummary) float64 {
	if stats == nil || stats.TotalMedications == 0 {
		return 100.0 // No medications needed = perfect score
	}

	score := 100.0

	// Score based on active medications and adherence
	if stats.AdherenceScore < 0.8 {
		score -= 30 // Poor medication management
	} else if stats.AdherenceScore < 0.9 {
		score -= 15 // Moderate medication management
	}

	return math.Max(0, math.Min(100, score))
}

// CalculateBowelMovementSummary calculates summary statistics for bowel movements
func (hds *HealthDataSummarizer) CalculateBowelMovementSummary(movements []bowelmovement.BowelMovement, days int) *analytics.BowelMovementSummary {
	if len(movements) == 0 {
		return &analytics.BowelMovementSummary{}
	}

	var totalPain, totalStrain, totalSatisfaction float64
	bristolCount := make(map[int]int)

	for _, movement := range movements {
		totalPain += float64(movement.Pain)
		totalStrain += float64(movement.Strain)
		totalSatisfaction += float64(movement.Satisfaction)
		bristolCount[movement.BristolType]++
	}

	// Find most common Bristol type
	mostCommonBristol := 0
	maxCount := 0
	for bristolType, count := range bristolCount {
		if count > maxCount {
			mostCommonBristol = bristolType
			maxCount = count
		}
	}

	count := float64(len(movements))
	// Avoid division by zero if days is 0
	if days == 0 {
		days = 1
	}
	averagePerDay := count / float64(days)

	return &analytics.BowelMovementSummary{
		TotalCount:          int64(len(movements)),
		AveragePerDay:       shared.RoundToDecimalPlaces(averagePerDay, 1),
		MostCommonBristol:   mostCommonBristol,
		AveragePain:         shared.RoundToDecimalPlaces(totalPain/count, 1),
		AverageStrain:       shared.RoundToDecimalPlaces(totalStrain/count, 1),
		AverageSatisfaction: shared.RoundToDecimalPlaces(totalSatisfaction/count, 1),
		RegularityScore:     hds.calculateRegularityScore(movements, days),
	}
}

// CalculateMealSummary calculates summary statistics for meals
func (hds *HealthDataSummarizer) CalculateMealSummary(meals []meal.Meal, days int) *analytics.MealSummary {
	if len(meals) == 0 {
		return &analytics.MealSummary{}
	}

	var totalCalories int
	var fiberRichCount int
	for _, meal := range meals {
		totalCalories += meal.Calories
		if meal.FiberRich {
			fiberRichCount++
		}
	}

	count := float64(len(meals))
	// Avoid division by zero if days is 0
	if days == 0 {
		days = 1
	}
	averagePerDay := count / float64(days)

	return &analytics.MealSummary{
		TotalMeals:       int64(len(meals)),
		AveragePerDay:    shared.RoundToDecimalPlaces(averagePerDay, 1),
		AverageCalories:  float64(totalCalories) / count,
		FiberRichPercent: (float64(fiberRichCount) / count) * 100,
		HealthScore:      hds.calculateMealHealthScore(meals),
		TotalCalories:    totalCalories,
	}
}

// CalculateSymptomSummary calculates summary statistics for symptoms
func (hds *HealthDataSummarizer) CalculateSymptomSummary(symptoms []symptom.Symptom, days int) *analytics.SymptomSummary {
	if len(symptoms) == 0 {
		return &analytics.SymptomSummary{}
	}

	var totalSeverity float64
	categoryCount := make(map[string]int)
	typeCount := make(map[string]int)

	for _, symptom := range symptoms {
		totalSeverity += float64(symptom.Severity)
		if symptom.Category != nil {
			categoryCount[symptom.Category.String()]++
		}
		if symptom.Type != nil {
			typeCount[symptom.Type.String()]++
		}
	}

	// Find most common category and type
	mostCommonCategory := hds.findMostCommon(categoryCount)
	mostCommonType := hds.findMostCommon(typeCount)

	return &analytics.SymptomSummary{
		TotalSymptoms:      int64(len(symptoms)),
		AverageSeverity:    shared.RoundToDecimalPlaces(totalSeverity/float64(len(symptoms)), 1),
		MostCommonCategory: mostCommonCategory,
		MostCommonType:     mostCommonType,
		TrendDirection:     "Stable", // Placeholder
	}
}

// CalculateMedicationSummary calculates summary statistics for medications
func (hds *HealthDataSummarizer) CalculateMedicationSummary(medications []medication.Medication) *analytics.MedicationSummary {
	if len(medications) == 0 {
		return &analytics.MedicationSummary{}
	}

	var activeCount int
	categoryCount := make(map[string]int)

	for _, med := range medications {
		if med.EndDate == nil || med.EndDate.After(time.Now()) {
			activeCount++
		}
		if med.Category != nil {
			categoryCount[med.Category.String()]++
		}
	}

	return &analytics.MedicationSummary{
		TotalMedications:   int64(len(medications)),
		ActiveMedications:  int64(activeCount),
		MostCommonCategory: hds.findMostCommon(categoryCount),
		AdherenceScore:     0.85, // Placeholder
		ComplexityScore:    hds.calculateMedicationComplexity(medications),
	}
}

// Helper functions

func (hds *HealthDataSummarizer) calculateRegularityScore(movements []bowelmovement.BowelMovement, days int) float64 {
	if len(movements) < 2 || days == 0 {
		return 50.0 // Neutral score for insufficient data
	}

	// Ideal frequency is 1-3 times per day
	averagePerDay := float64(len(movements)) / float64(days)
	if averagePerDay >= 1 && averagePerDay <= 3 {
		return 90.0
	} else if averagePerDay > 0.5 && averagePerDay < 1 {
		return 70.0
	} else if averagePerDay < 0.5 {
		return 40.0
	} else {
		return 50.0 // for > 3 per day
	}
}

func (hds *HealthDataSummarizer) calculateMealHealthScore(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 75.0 // Neutral score
	}

	// Simplified health scoring based on available fields
	var score float64
	var fiberRichCount int

	for _, meal := range meals {
		if meal.FiberRich {
			fiberRichCount++
		}
	}

	fiberRatio := float64(fiberRichCount) / float64(len(meals))
	score = fiberRatio * 100

	return math.Max(0, math.Min(100, score))
}

// MedicationComplexityThresholds defines the boundaries for complexity scoring.
// These thresholds are based on clinical guidelines where polypharmacy (the use of multiple drugs)
// is often categorized by the number of medications a patient is taking.
//
// Rationale for medication count thresholds:
// - 0 medications: No complexity (optimal for digestive health)
// - 1-2 medications: Low complexity (minimal drug interactions, simple regimen)
// - 3-5 medications: Medium complexity (moderate potential for interactions, requires tracking)
// - 6-10 medications: High complexity (significant potential for interactions, difficult adherence)
// - >10 medications: Very high complexity (high risk of adverse events, challenging management)
//
// These thresholds align with common clinical definitions where polypharmacy is often
// defined as taking 5+ medications, and excessive polypharmacy as taking 10+ medications.
const (
	// Medication count thresholds
	lowComplexityMax    = 2  // 1-2 medications: low complexity
	mediumComplexityMax = 5  // 3-5 medications: medium complexity (clinical definition of polypharmacy often starts at 5)
	highComplexityMax   = 10 // 6-10 medications: high complexity

	// Complexity score values (0-100 scale)
	noMedicationsScore      = 0   // No medications = optimal score
	lowComplexityScore      = 25  // Low complexity burden
	mediumComplexityScore   = 50  // Medium complexity burden
	highComplexityScore     = 75  // High complexity burden
	veryHighComplexityScore = 100 // Very high complexity burden
)

func (hds *HealthDataSummarizer) calculateMedicationComplexity(medications []medication.Medication) float64 {
	// Calculate complexity score based on the number of medications
	// Higher medication counts indicate greater regimen complexity, increased
	// risk of drug interactions, and potential adherence challenges
	count := len(medications)

	if count == 0 {
		// No medications = no complexity
		return noMedicationsScore
	} else if count <= lowComplexityMax {
		// 1-2 medications = simple regimen with minimal interactions
		return lowComplexityScore
	} else if count <= mediumComplexityMax {
		// 3-5 medications = moderate complexity with potential interactions
		return mediumComplexityScore
	} else if count <= highComplexityMax {
		// 6-10 medications = complex regimen with significant interaction risk
		return highComplexityScore
	} else {
		// >10 medications = very complex regimen with high interaction risk
		return veryHighComplexityScore
	}
}

func (hds *HealthDataSummarizer) findMostCommon(counts map[string]int) string {
	if len(counts) == 0 {
		return ""
	}

	maxCount := 0
	mostCommon := ""
	for item, count := range counts {
		if count > maxCount {
			maxCount = count
			mostCommon = item
		}
	}
	return mostCommon
}
