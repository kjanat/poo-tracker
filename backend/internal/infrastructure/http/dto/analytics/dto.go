package analytics

import (
	"errors"
	"time"
)

// BowelMovementStatsRequest represents the request for bowel movement statistics
type BowelMovementStatsRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Period    *string    `json:"period,omitempty" binding:"omitempty,oneof=day week month year"`
}

// MealCorrelationRequest represents the request for meal correlation analysis
type MealCorrelationRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// SymptomPatternRequest represents the request for symptom pattern analysis
type SymptomPatternRequest struct {
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	SymptomTypes []string   `json:"symptom_types,omitempty"`
}

// TrendAnalysisRequest represents the request for trend analysis
type TrendAnalysisRequest struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Metrics   []string   `json:"metrics,omitempty" binding:"omitempty,dive,oneof=frequency consistency type urgency symptoms"`
}

// BowelMovementStatsResponse represents bowel movement statistics
type BowelMovementStatsResponse struct {
	Period           string           `json:"period"`
	TotalMovements   int              `json:"total_movements"`
	AverageDaily     float64          `json:"average_daily"`
	ConsistencyStats ConsistencyStats `json:"consistency_stats"`
	TypeStats        TypeStats        `json:"type_stats"`
	UrgencyStats     UrgencyStats     `json:"urgency_stats"`
	TimePatterns     TimePatterns     `json:"time_patterns"`
	WeeklyTrends     []WeeklyTrend    `json:"weekly_trends"`
}

// ConsistencyStats represents consistency statistics
type ConsistencyStats struct {
	Average      float64            `json:"average"`
	Distribution map[string]int     `json:"distribution"`
	Trends       []ConsistencyTrend `json:"trends"`
}

// TypeStats represents Bristol stool type statistics
type TypeStats struct {
	Distribution map[string]int `json:"distribution"`
	MostCommon   string         `json:"most_common"`
	Trends       []TypeTrend    `json:"trends"`
}

// UrgencyStats represents urgency statistics
type UrgencyStats struct {
	Average      float64        `json:"average"`
	Distribution map[string]int `json:"distribution"`
	HighUrgency  int            `json:"high_urgency_count"`
	Trends       []UrgencyTrend `json:"trends"`
}

// TimePatterns represents patterns in timing
type TimePatterns struct {
	HourlyDistribution  map[string]int `json:"hourly_distribution"`
	MostCommonHour      int            `json:"most_common_hour"`
	AverageTimeOfDay    string         `json:"average_time_of_day"`
	WeekdayDistribution map[string]int `json:"weekday_distribution"`
}

// WeeklyTrend represents weekly trends
type WeeklyTrend struct {
	Week               string  `json:"week"`
	Count              int     `json:"count"`
	AverageConsistency float64 `json:"average_consistency"`
	MostCommonType     string  `json:"most_common_type"`
}

// ConsistencyTrend represents consistency trend data
type ConsistencyTrend struct {
	Date    string  `json:"date"`
	Average float64 `json:"average"`
}

// TypeTrend represents type trend data
type TypeTrend struct {
	Date         string         `json:"date"`
	Distribution map[string]int `json:"distribution"`
}

// UrgencyTrend represents urgency trend data
type UrgencyTrend struct {
	Date    string  `json:"date"`
	Average float64 `json:"average"`
}

// MealCorrelationResponse represents meal correlation analysis results
type MealCorrelationResponse struct {
	Correlations    []MealCorrelation `json:"correlations"`
	TriggerFoods    []TriggerFood     `json:"trigger_foods"`
	DigestiveTiming DigestiveTiming   `json:"digestive_timing"`
}

// MealCorrelation represents correlation between meals and bowel movements
type MealCorrelation struct {
	MealName     string    `json:"meal_name"`
	MealTime     time.Time `json:"meal_time"`
	MovementTime time.Time `json:"movement_time"`
	TimeDiff     int       `json:"time_diff_hours"`
	Consistency  int       `json:"consistency"`
	Type         int       `json:"type"`
	Urgency      int       `json:"urgency"`
}

// TriggerFood represents foods that may trigger symptoms
type TriggerFood struct {
	Food             string   `json:"food"`
	SymptomCount     int      `json:"symptom_count"`
	CorrelationScore float64  `json:"correlation_score"`
	CommonSymptoms   []string `json:"common_symptoms"`
}

// DigestiveTiming represents digestive timing patterns
type DigestiveTiming struct {
	AverageTransitTime  int    `json:"average_transit_time_hours"`
	FastestTransitTime  int    `json:"fastest_transit_time_hours"`
	SlowestTransitTime  int    `json:"slowest_transit_time_hours"`
	MostCommonTimeRange string `json:"most_common_time_range"`
}

// SymptomPatternResponse represents symptom pattern analysis results
type SymptomPatternResponse struct {
	Patterns          []SymptomPattern  `json:"patterns"`
	FrequencyAnalysis FrequencyAnalysis `json:"frequency_analysis"`
	SeverityTrends    []SeverityTrend   `json:"severity_trends"`
	Triggers          []SymptomTrigger  `json:"triggers"`
}

// SymptomPattern represents a pattern in symptoms
type SymptomPattern struct {
	Type            string   `json:"type"`
	Frequency       int      `json:"frequency"`
	AverageSeverity float64  `json:"average_severity"`
	CommonTimes     []string `json:"common_times"`
	Duration        int      `json:"average_duration_minutes"`
}

// FrequencyAnalysis represents frequency analysis of symptoms
type FrequencyAnalysis struct {
	TotalSymptoms    int            `json:"total_symptoms"`
	DailyAverage     float64        `json:"daily_average"`
	MostCommonType   string         `json:"most_common_type"`
	TypeDistribution map[string]int `json:"type_distribution"`
}

// SeverityTrend represents severity trends over time
type SeverityTrend struct {
	Date            string  `json:"date"`
	AverageSeverity float64 `json:"average_severity"`
	SymptomCount    int     `json:"symptom_count"`
}

// SymptomTrigger represents potential symptom triggers
type SymptomTrigger struct {
	Trigger          string   `json:"trigger"`
	SymptomTypes     []string `json:"symptom_types"`
	Frequency        int      `json:"frequency"`
	CorrelationScore float64  `json:"correlation_score"`
}

// TrendAnalysisResponse represents trend analysis results
type TrendAnalysisResponse struct {
	FrequencyTrend   []TrendPoint `json:"frequency_trend"`
	ConsistencyTrend []TrendPoint `json:"consistency_trend"`
	TypeTrend        []TrendPoint `json:"type_trend"`
	UrgencyTrend     []TrendPoint `json:"urgency_trend"`
	SymptomTrend     []TrendPoint `json:"symptom_trend"`
	Insights         []Insight    `json:"insights"`
}

// TrendPoint represents a point in a trend analysis
type TrendPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Label string  `json:"label,omitempty"`
}

// Insight represents an analytical insight
type Insight struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"` // low, medium, high
	ActionItems []string `json:"action_items,omitempty"`
}

// HealthMetricsResponse represents overall health metrics
type HealthMetricsResponse struct {
	OverallScore     int              `json:"overall_score"` // 0-100
	RegularityScore  int              `json:"regularity_score"`
	ConsistencyScore int              `json:"consistency_score"`
	SymptomScore     int              `json:"symptom_score"`
	TrendDirection   string           `json:"trend_direction"` // improving, stable, declining
	LastUpdated      time.Time        `json:"last_updated"`
	Recommendations  []Recommendation `json:"recommendations"`
}

// Recommendation represents a health recommendation
type Recommendation struct {
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"` // low, medium, high
	Actions     []string `json:"actions"`
}

// Validate validates the BowelMovementStatsRequest
func (r *BowelMovementStatsRequest) Validate() error {
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

// Validate validates the MealCorrelationRequest
func (r *MealCorrelationRequest) Validate() error {
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

// Validate validates the SymptomPatternRequest
func (r *SymptomPatternRequest) Validate() error {
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

// Validate validates the TrendAnalysisRequest
func (r *TrendAnalysisRequest) Validate() error {
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}
