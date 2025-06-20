package analytics

import "time"

// BowelMovementSummary represents bowel movement analytics summary
type BowelMovementSummary struct {
	TotalCount          int64   `json:"totalCount"`
	AveragePerDay       float64 `json:"averagePerDay"`
	MostCommonBristol   int     `json:"mostCommonBristol"`
	AveragePain         float64 `json:"averagePain"`
	AverageStrain       float64 `json:"averageStrain"`
	AverageSatisfaction float64 `json:"averageSatisfaction"`
	RegularityScore     float64 `json:"regularityScore"`
}

// MealSummary represents meal analytics summary
type MealSummary struct {
	TotalMeals       int64   `json:"totalMeals"`
	AveragePerDay    float64 `json:"averagePerDay"`
	TotalCalories    int     `json:"totalCalories"`
	AverageCalories  float64 `json:"averageCalories"`
	FiberRichPercent float64 `json:"fiberRichPercent"`
	HealthScore      float64 `json:"healthScore"`
}

// SymptomSummary represents symptom analytics summary
type SymptomSummary struct {
	TotalSymptoms      int64   `json:"totalSymptoms"`
	AveragePerDay      float64 `json:"averagePerDay"`
	AverageSeverity    float64 `json:"averageSeverity"`
	MostCommonCategory string  `json:"mostCommonCategory"`
	MostCommonType     string  `json:"mostCommonType"`
	TrendDirection     string  `json:"trendDirection"`
}

// MedicationSummary represents medication analytics summary
type MedicationSummary struct {
	TotalMedications   int64   `json:"totalMedications"`
	ActiveMedications  int64   `json:"activeMedications"`
	AdherenceScore     float64 `json:"adherenceScore"`
	MostCommonCategory string  `json:"mostCommonCategory"`
	ComplexityScore    float64 `json:"complexityScore"`
}

// ScoreFactor represents a factor contributing to the health score
type ScoreFactor struct {
	Name        string  `json:"name"`
	Impact      float64 `json:"impact"` // -100 to 100
	Weight      float64 `json:"weight"` // 0-1
	Description string  `json:"description"`
	Trend       string  `json:"trend"`
}

// DataTrend represents a trend in data over time
type DataTrend struct {
	Direction   string             `json:"direction"`  // IMPROVING, STABLE, DECLINING
	Slope       float64            `json:"slope"`      // Rate of change
	Confidence  float64            `json:"confidence"` // 0 to 1
	TimePoints  []time.Time        `json:"timePoints"`
	Values      []float64          `json:"values"`
	Seasonality map[string]float64 `json:"seasonality"` // Seasonal patterns
}

// Insight represents an actionable insight
type Insight struct {
	Type        string   `json:"type"` // DIETARY, LIFESTYLE, MEDICAL, BEHAVIORAL
	Category    string   `json:"category"`
	Message     string   `json:"message"`
	Evidence    string   `json:"evidence"`
	Priority    string   `json:"priority"`   // LOW, MEDIUM, HIGH
	Confidence  float64  `json:"confidence"` // 0-1
	ActionItems []string `json:"actionItems"`
}
