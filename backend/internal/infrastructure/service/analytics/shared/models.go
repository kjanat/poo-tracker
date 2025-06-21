package shared

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// DailyAggregation represents aggregated data for a single day
type DailyAggregation struct {
	Date            time.Time
	BowelMovements  []bowelmovement.BowelMovement
	Meals           []meal.Meal
	Symptoms        []symptom.Symptom
	Medications     []medication.Medication
	BristolAverage  float64
	PainAverage     float64
	StrainAverage   float64
	SatisfactionAvg float64
	MealCount       int
	SymptomCount    int
	SpicyMealCount  int
	DairyMealCount  int
	GlutenMealCount int
}

// TrendPoint represents a single data point in a trend analysis
type TrendPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

// Use domain analytics types for shared structures
type CorrelationPair = analytics.Correlation
type TrendLine = analytics.DataTrend
type PatternMatch = analytics.Insight
type HealthMetric = analytics.ScoreFactor

// InsightRecommendation references the domain Recommendation type for
// actionable advice returned by the insight engine.
type InsightRecommendation = analytics.Recommendation

// StatisticalSummary represents basic summary statistics for a data set.
type StatisticalSummary struct {
	Count        int
	Mean         float64
	Median       float64
	StdDev       float64
	Min          float64
	Max          float64
	Percentile25 float64
	Percentile75 float64
}
type BowelMovementStats = analytics.BowelMovementSummary
type MealStats = analytics.MealSummary
type SymptomStats = analytics.SymptomSummary
type MedicationStats = analytics.MedicationSummary
