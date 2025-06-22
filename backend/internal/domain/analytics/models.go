// Package analytics provides domain models for health analytics functionality.
// This file contains all the data structures used for health analytics, including
// summaries, trends, insights, and scoring systems.
//
// Key Models:
// - BowelMovementSummary: Aggregated statistics for bowel movements
// - MealSummary: Nutritional and dietary analytics
// - SymptomSummary: Symptom tracking and trend analysis
// - MedicationSummary: Medication adherence and effectiveness metrics
// - DataTrend: Time-series trend analysis with direction and confidence
// - Insight: Actionable health insights with evidence and priority
// - ScoreFactor: Individual factors contributing to overall health scores
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

// HealthOverview represents a comprehensive health overview
type HealthOverview struct {
	Period             string                `json:"period"`
	BowelMovementStats *BowelMovementSummary `json:"bowelMovementStats"`
	MealStats          *MealSummary          `json:"mealStats"`
	SymptomStats       *SymptomSummary       `json:"symptomStats"`
	MedicationStats    *MedicationSummary    `json:"medicationStats"`
	OverallHealthScore float64               `json:"overallHealthScore"`
	TrendDirection     string                `json:"trendDirection"` // IMPROVING, STABLE, DECLINING
}

// CorrelationAnalysis represents analysis of correlations between different data types
type CorrelationAnalysis struct {
	MealBowelCorrelations   []*Correlation      `json:"mealBowelCorrelations"`
	MealSymptomCorrelations []*Correlation      `json:"mealSymptomCorrelations"`
	MedicationEffectiveness []*MedicationEffect `json:"medicationEffectiveness"`
	TriggerAnalysis         []*TriggerEffect    `json:"triggerAnalysis"`
}

// Correlation represents a correlation between two data points
type Correlation struct {
	Factor      string  `json:"factor"`
	Outcome     string  `json:"outcome"`
	Strength    float64 `json:"strength"`   // -1 to 1
	Confidence  float64 `json:"confidence"` // 0 to 1
	Description string  `json:"description"`
	SampleSize  int     `json:"sampleSize"`
}

// MedicationEffect represents the effect of a medication
type MedicationEffect struct {
	MedicationName     string             `json:"medicationName"`
	SymptomImprovement map[string]float64 `json:"symptomImprovement"`
	BowelImprovement   float64            `json:"bowelImprovement"`
	EffectivenessScore float64            `json:"effectivenessScore"`
}

// TriggerEffect represents the effect of various triggers
type TriggerEffect struct {
	Trigger         string  `json:"trigger"`
	SymptomIncrease float64 `json:"symptomIncrease"`
	BowelImpact     float64 `json:"bowelImpact"`
	Frequency       int     `json:"frequency"`
	Severity        float64 `json:"severity"`
}

// TrendAnalysis represents trend analysis over time
type TrendAnalysis struct {
	BowelMovementTrends *DataTrend `json:"bowelMovementTrends"`
	SymptomTrends       *DataTrend `json:"symptomTrends"`
	MealTrends          *DataTrend `json:"mealTrends"`
	OverallTrend        string     `json:"overallTrend"`
}

// BehaviorPatterns represents patterns in user behavior
type BehaviorPatterns struct {
	EatingPatterns    *EatingPattern    `json:"eatingPatterns"`
	BowelPatterns     *BowelPattern     `json:"bowelPatterns"`
	SymptomPatterns   *SymptomPattern   `json:"symptomPatterns"`
	LifestylePatterns *LifestylePattern `json:"lifestylePatterns"`
}

// EatingPattern represents patterns in eating behavior
type EatingPattern struct {
	MealTiming           map[string]float64 `json:"mealTiming"`           // Hour -> frequency
	MealSizeDistribution map[string]float64 `json:"mealSizeDistribution"` // Size -> frequency
	PreferredCuisines    []string           `json:"preferredCuisines"`
	DietaryConsistency   float64            `json:"dietaryConsistency"` // 0-1
}

// BowelPattern represents patterns in bowel movements
type BowelPattern struct {
	PreferredTiming     map[string]float64 `json:"preferredTiming"`     // Hour -> frequency
	RegularityScore     float64            `json:"regularityScore"`     // 0-1
	ConsistencyPatterns map[string]float64 `json:"consistencyPatterns"` // Bristol type -> frequency
	ResponseToMeals     float64            `json:"responseToMeals"`     // Hours after meals
}

// SymptomPattern represents patterns in symptom occurrence
type SymptomPattern struct {
	SymptomTiming     map[string]float64 `json:"symptomTiming"`     // Hour -> frequency
	TriggerPatterns   map[string]float64 `json:"triggerPatterns"`   // Trigger -> frequency
	SeasonalVariation map[string]float64 `json:"seasonalVariation"` // Season -> frequency
	CyclicalPatterns  bool               `json:"cyclicalPatterns"`  // Whether symptoms follow cycles
}

// LifestylePattern represents patterns in lifestyle factors
type LifestylePattern struct {
	StressLevels       map[string]float64 `json:"stressLevels"`       // Level -> frequency
	SleepQuality       map[string]float64 `json:"sleepQuality"`       // Quality -> frequency
	ExerciseFrequency  map[string]float64 `json:"exerciseFrequency"`  // Intensity -> frequency
	WeatherSensitivity map[string]float64 `json:"weatherSensitivity"` // Weather -> impact
}

// HealthInsights represents actionable health insights
type HealthInsights struct {
	KeyFindings     []string   `json:"keyFindings"`
	RiskFactors     []string   `json:"riskFactors"`
	PositiveFactors []string   `json:"positiveFactors"`
	Recommendations []*Insight `json:"recommendations"`
	AlertLevel      string     `json:"alertLevel"`      // LOW, MEDIUM, HIGH
	ConfidenceLevel float64    `json:"confidenceLevel"` // 0-1
}

// HealthScore represents an overall health score
type HealthScore struct {
	OverallScore    float64            `json:"overallScore"`    // 0-100
	ComponentScores map[string]float64 `json:"componentScores"` // Component -> score
	Trend           string             `json:"trend"`           // IMPROVING, STABLE, DECLINING
	LastUpdated     time.Time          `json:"lastUpdated"`
	Factors         []*ScoreFactor     `json:"factors"`
	Benchmarks      map[string]float64 `json:"benchmarks"` // Comparison to population averages
}

// Recommendation represents a personalized recommendation
type Recommendation struct {
	ID             string     `json:"id"`
	Type           string     `json:"type"` // DIETARY, LIFESTYLE, MEDICAL, TRACKING
	Category       string     `json:"category"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`   // LOW, MEDIUM, HIGH
	Confidence     float64    `json:"confidence"` // 0-1
	Evidence       []string   `json:"evidence"`
	ActionSteps    []string   `json:"actionSteps"`
	ExpectedImpact string     `json:"expectedImpact"`
	Timeline       string     `json:"timeline"` // When to expect results
	CreatedAt      time.Time  `json:"createdAt"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
}
