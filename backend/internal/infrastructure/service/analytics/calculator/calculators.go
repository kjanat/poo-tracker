package calculator

import (
	"errors"
	"log"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// DEPRECATED: This file is kept for backward compatibility only.
// All functionality has moved to health_scores.go using the improved architecture
// with separation of concerns (HealthScoreCalculator and HealthDataSummarizer).
//
// If you're using any methods from this file, please migrate to the implementations
// in health_scores.go. This file will be removed in a future release.

// ErrDeprecated is returned by all methods in this file to indicate they should not be used
var ErrDeprecated = errors.New("deprecated: this method is a stub and has been replaced in health_scores.go")

// LegacyCalculator is a deprecated stub implementation. Use HealthScoreCalculator and
// HealthDataSummarizer from health_scores.go instead.
type LegacyCalculator struct{}

// NewLegacyCalculator creates a new legacy calculator (deprecated)
func NewLegacyCalculator() *LegacyCalculator {
	log.Println("WARNING: Using deprecated LegacyCalculator from calculators.go")
	return &LegacyCalculator{}
}

// CalculateRegularityScore is deprecated - use HealthDataSummarizer.calculateRegularityScore
func (c *LegacyCalculator) CalculateRegularityScore(movements []bowelmovement.BowelMovement) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateRegularityScore")
	return 0.0, ErrDeprecated
}

// CalculateMostCommonBristol is deprecated
func (c *LegacyCalculator) CalculateMostCommonBristol(movements []bowelmovement.BowelMovement) (int, error) {
	log.Println("WARNING: Using deprecated CalculateMostCommonBristol")
	return 0, ErrDeprecated
}

// CalculateAveragePain is deprecated
func (c *LegacyCalculator) CalculateAveragePain(movements []bowelmovement.BowelMovement) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateAveragePain")
	return 0.0, ErrDeprecated
}

// CalculateAverageStrain is deprecated
func (c *LegacyCalculator) CalculateAverageStrain(movements []bowelmovement.BowelMovement) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateAverageStrain")
	return 0.0, ErrDeprecated
}

// CalculateAverageSatisfaction is deprecated
func (c *LegacyCalculator) CalculateAverageSatisfaction(movements []bowelmovement.BowelMovement) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateAverageSatisfaction")
	return 0.0, ErrDeprecated
}

// CalculateTotalCalories is deprecated
func (c *LegacyCalculator) CalculateTotalCalories(meals []meal.Meal) (int, error) {
	log.Println("WARNING: Using deprecated CalculateTotalCalories")
	return 0, ErrDeprecated
}

// CalculateAverageCalories is deprecated
func (c *LegacyCalculator) CalculateAverageCalories(meals []meal.Meal) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateAverageCalories")
	return 0.0, ErrDeprecated
}

// CalculateFiberRichPercent is deprecated
func (c *LegacyCalculator) CalculateFiberRichPercent(meals []meal.Meal) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateFiberRichPercent")
	return 0.0, ErrDeprecated
}

// CalculateMealHealthScore is deprecated
func (c *LegacyCalculator) CalculateMealHealthScore(meals []meal.Meal) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateMealHealthScore")
	return 0.0, ErrDeprecated
}

// CalculateAverageSeverity is deprecated
func (c *LegacyCalculator) CalculateAverageSeverity(symptoms []symptom.Symptom) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateAverageSeverity")
	return 0.0, ErrDeprecated
}

// FindMostCommonCategory is deprecated
func (c *LegacyCalculator) FindMostCommonCategory(symptoms []symptom.Symptom) (string, error) {
	log.Println("WARNING: Using deprecated FindMostCommonCategory")
	return "", ErrDeprecated
}

// FindMostCommonType is deprecated
func (c *LegacyCalculator) FindMostCommonType(symptoms []symptom.Symptom) (string, error) {
	log.Println("WARNING: Using deprecated FindMostCommonType")
	return "", ErrDeprecated
}

// CalculateSymptomTrendDirection is deprecated
func (c *LegacyCalculator) CalculateSymptomTrendDirection(symptoms []symptom.Symptom) (string, error) {
	log.Println("WARNING: Using deprecated CalculateSymptomTrendDirection")
	return "", ErrDeprecated
}

// CountActiveMedications is deprecated
func (c *LegacyCalculator) CountActiveMedications(medications []medication.Medication) (int64, error) {
	log.Println("WARNING: Using deprecated CountActiveMedications")
	return 0, ErrDeprecated
}

// CalculateAdherenceScore is deprecated
func (c *LegacyCalculator) CalculateAdherenceScore(medications []medication.Medication) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateAdherenceScore")
	return 0.0, ErrDeprecated
}

// FindMostCommonMedicationCategory is deprecated
func (c *LegacyCalculator) FindMostCommonMedicationCategory(medications []medication.Medication) (string, error) {
	log.Println("WARNING: Using deprecated FindMostCommonMedicationCategory")
	return "", ErrDeprecated
}

// CalculateMedicationComplexity is deprecated
func (c *LegacyCalculator) CalculateMedicationComplexity(medications []medication.Medication) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateMedicationComplexity")
	return 0.0, ErrDeprecated
}

// CalculateOverallHealthScore is deprecated
func (c *LegacyCalculator) CalculateOverallHealthScore(bowels []bowelmovement.BowelMovement, meals []meal.Meal, symptoms []symptom.Symptom) (float64, error) {
	log.Println("WARNING: Using deprecated CalculateOverallHealthScore")
	return 0.0, ErrDeprecated
}

// CalculateHealthTrend is deprecated
func (c *LegacyCalculator) CalculateHealthTrend(bowels []bowelmovement.BowelMovement, meals []meal.Meal, symptoms []symptom.Symptom) (string, error) {
	log.Println("WARNING: Using deprecated CalculateHealthTrend")
	return "", ErrDeprecated
}
