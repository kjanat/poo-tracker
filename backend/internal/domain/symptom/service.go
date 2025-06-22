package symptom

import (
	"context"
	"time"
)

// Service defines the interface for symptom business logic
type Service interface {
	// Core operations
	Create(ctx context.Context, userID string, input *CreateSymptomInput) (*Symptom, error)
	GetByID(ctx context.Context, id string) (*Symptom, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Symptom, error)
	Update(ctx context.Context, id string, input *UpdateSymptomInput) (*Symptom, error)
	Delete(ctx context.Context, id string) error

	// Query operations
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*Symptom, error)
	GetByCategory(ctx context.Context, userID string, category string) ([]*Symptom, error)
	GetBySeverityRange(ctx context.Context, userID string, minSeverity, maxSeverity int) ([]*Symptom, error)
	GetLatest(ctx context.Context, userID string) (*Symptom, error)

	// Analytics operations
	GetSymptomStats(ctx context.Context, userID string, start, end time.Time) (*SymptomStats, error)
	GetTriggerInsights(ctx context.Context, userID string, start, end time.Time) (*TriggerInsights, error)
	GetSymptomPatterns(ctx context.Context, userID string, start, end time.Time) (*SymptomPatterns, error)
}

// CreateSymptomInput represents input for creating a symptom
type CreateSymptomInput struct {
	Name        string    `json:"name" binding:"required,min=1,max=200"`
	Description string    `json:"description,omitempty"`
	RecordedAt  time.Time `json:"recordedAt" binding:"required"`
	Category    *string   `json:"category,omitempty"`
	Severity    int       `json:"severity" binding:"required,min=1,max=10"`
	Duration    *int      `json:"duration,omitempty" binding:"omitempty,min=1"`
	BodyPart    string    `json:"bodyPart,omitempty"`
	Type        *string   `json:"type,omitempty"`
	Triggers    []string  `json:"triggers,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	PhotoURL    string    `json:"photoUrl,omitempty"`
}

// UpdateSymptomInput represents input for updating a symptom
type UpdateSymptomInput struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description,omitempty"`
	RecordedAt  *time.Time `json:"recordedAt,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Severity    *int       `json:"severity,omitempty" binding:"omitempty,min=1,max=10"`
	Duration    *int       `json:"duration,omitempty" binding:"omitempty,min=1"`
	BodyPart    *string    `json:"bodyPart,omitempty"`
	Type        *string    `json:"type,omitempty"`
	Triggers    []string   `json:"triggers,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
	PhotoURL    *string    `json:"photoUrl,omitempty"`
}

// SymptomStats represents analytics data for symptoms
type SymptomStats struct {
	TotalCount           int64          `json:"totalCount"`
	AverageSeverity      float64        `json:"averageSeverity"`
	MostCommonCategory   string         `json:"mostCommonCategory"`
	MostCommonType       string         `json:"mostCommonType"`
	CategoryBreakdown    map[string]int `json:"categoryBreakdown"`
	TypeBreakdown        map[string]int `json:"typeBreakdown"`
	SeverityDistribution map[int]int    `json:"severityDistribution"`
	AverageDuration      float64        `json:"averageDuration"`
}

// TriggerInsights represents analysis of symptom triggers
type TriggerInsights struct {
	MostCommonTriggers []string           `json:"mostCommonTriggers"`
	TriggerFrequency   map[string]int     `json:"triggerFrequency"`
	TriggerSeverityMap map[string]float64 `json:"triggerSeverityMap"`
	UniqueTriggerCount int                `json:"uniqueTriggerCount"`
}

// SymptomPatterns represents patterns in symptom occurrence
type SymptomPatterns struct {
	TimeOfDayPatterns  map[string]float64 `json:"timeOfDayPatterns"` // Hour -> frequency
	DayOfWeekPatterns  map[string]float64 `json:"dayOfWeekPatterns"` // Weekday -> frequency
	SeasonalPatterns   map[string]float64 `json:"seasonalPatterns"`  // Season -> frequency
	FrequencyPerDay    float64            `json:"frequencyPerDay"`
	AverageTimeBetween float64            `json:"averageTimeBetween"` // Hours between symptoms
}
