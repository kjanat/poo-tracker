package analytics

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
)

// Use domain analytics types for shared structures
type BowelMovementStats = analytics.BowelMovementSummary
type MealStats = analytics.MealSummary
type SymptomStats = analytics.SymptomSummary
type MedicationStats = analytics.MedicationSummary
type HealthOverview = analytics.HealthOverview
type CorrelationAnalysis = analytics.CorrelationAnalysis
type BowelPattern = analytics.BowelPattern
type SymptomPattern = analytics.SymptomPattern
type EatingPattern = analytics.EatingPattern
type BehaviorPatterns = analytics.BehaviorPatterns
type Correlation = analytics.Correlation
type HealthInsights = analytics.HealthInsights
type Recommendation = analytics.Recommendation

// RiskFactor represents an identified health risk
type RiskFactor struct {
	Factor     string   `json:"factor"`
	Risk       string   `json:"risk"`       // LOW, MEDIUM, HIGH
	Confidence float64  `json:"confidence"` // 0 to 1
	Impact     string   `json:"impact"`
	Mitigation []string `json:"mitigation"`
}

// TrendAnalysis represents trend analysis over time
type TrendAnalysis struct {
	Period         string               `json:"period"`
	Start          time.Time            `json:"start"`
	End            time.Time            `json:"end"`
	HealthTrends   map[string][]float64 `json:"healthTrends"`
	BehaviorTrends map[string][]string  `json:"behaviorTrends"`
	RiskFactors    []RiskFactor         `json:"riskFactors"`
	LastAnalyzed   time.Time            `json:"lastAnalyzed"`
}

// MetricsProvider defines the interface for obtaining analytics metrics
type MetricsProvider interface {
	// GetMetrics retrieves health metrics for a given time range
	GetMetrics(userID string, start, end time.Time) (*HealthMetrics, error)
}

// CorrelationAnalyzer defines the interface for analyzing correlations between health data
type CorrelationAnalyzer interface {
	// AnalyzeCorrelations performs correlation analysis on health data
	AnalyzeCorrelations(userID string, start, end time.Time) (*CorrelationAnalysis, error)
}

// TrendAnalyzer defines the interface for analyzing health trends
type TrendAnalyzer interface {
	// AnalyzeTrends performs trend analysis on health data
	AnalyzeTrends(userID string, start, end time.Time) (*TrendAnalysis, error)
}

// HealthMetrics represents calculated health metrics over a period
type HealthMetrics struct {
	Period         string    `json:"period"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	OverallScore   float64   `json:"overallScore"`
	BowelScore     float64   `json:"bowelScore"`
	DietScore      float64   `json:"dietScore"`
	SymptomScore   float64   `json:"symptomScore"`
	LifestyleScore float64   `json:"lifestyleScore"`
	TrendDirection string    `json:"trendDirection"` // IMPROVING, STABLE, DECLINING
	LastCalculated time.Time `json:"lastCalculated"`
}

// Correlation represents a correlation between two factors

// TriggerAnalysis represents analysis of triggering factors
type TriggerAnalysis struct {
	Trigger      string    `json:"trigger"`
	Effect       string    `json:"effect"`
	Probability  float64   `json:"probability"` // 0 to 1
	TimeToOnset  string    `json:"timeToOnset"` // Duration string
	LastObserved time.Time `json:"lastObserved"`
}

// LifestyleImpact represents the impact of lifestyle factors
type LifestyleImpact struct {
	Factor       string    `json:"factor"`
	Impact       float64   `json:"impact"`     // -1 to 1
	Confidence   float64   `json:"confidence"` // 0 to 1
	Description  string    `json:"description"`
	LastAssessed time.Time `json:"lastAssessed"`
}

// HealthScore represents overall health scoring
type HealthScore struct {
	Overall        float64            `json:"overall"`
	Components     map[string]float64 `json:"components"`
	TrendDirection string             `json:"trendDirection"`
	LastCalculated time.Time          `json:"lastCalculated"`
}

// HealthInsights represents health insights and recommendations

// Recommendation represents a health recommendation
