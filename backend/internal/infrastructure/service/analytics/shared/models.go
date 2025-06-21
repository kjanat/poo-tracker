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
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// TrendLine represents a calculated data trend
type TrendLine struct {
	Name         string       `json:"name"`
	Points       []TrendPoint `json:"points"`
	Direction    string       `json:"direction"`
	Slope        float64      `json:"slope"`
	Confidence   float64      `json:"confidence"`
	Significance string       `json:"significance"`
}

// Use domain analytics types for shared structures
type CorrelationPair = analytics.Correlation
type PatternMatch = analytics.Insight
type HealthMetric = analytics.ScoreFactor

// InsightRecommendation represents a recommendation produced by the analytics
// engine before it is converted to a domain Recommendation.
type InsightRecommendation struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category,omitempty"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Priority    string                 `json:"priority"`
	Confidence  float64                `json:"confidence,omitempty"`
	Evidence    []string               `json:"evidence"`
	Actions     []string               `json:"actions"`
	Context     map[string]interface{} `json:"context,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}

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
