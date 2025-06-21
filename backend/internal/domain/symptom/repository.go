package symptom

import (
	"context"
	"time"
)

// Repository defines the interface for symptom data persistence
type Repository interface {
	// Symptom CRUD operations
	Create(ctx context.Context, symptom *Symptom) error
	GetByID(ctx context.Context, id string) (*Symptom, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Symptom, error)
	Update(ctx context.Context, id string, update *SymptomUpdate) error
	Delete(ctx context.Context, id string) error

	// Query operations
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*Symptom, error)
	GetByCategory(ctx context.Context, userID string, category string) ([]*Symptom, error)
	GetBySeverity(ctx context.Context, userID string, minSeverity, maxSeverity int) ([]*Symptom, error)
	GetCountByUserID(ctx context.Context, userID string) (int64, error)
	GetLatestByUserID(ctx context.Context, userID string) (*Symptom, error)

	// Analytics operations
	GetSeverityStats(ctx context.Context, userID string, start, end time.Time) (*SeverityStats, error)
	GetCategoryFrequency(ctx context.Context, userID string, start, end time.Time) (map[string]int, error)
	GetTriggerAnalysis(ctx context.Context, userID string, start, end time.Time) (map[string]int, error)
}

// SeverityStats represents calculated severity statistics for analytics
type SeverityStats struct {
	AverageSeverity      float64     `json:"averageSeverity"`
	MinSeverity          int         `json:"minSeverity"`
	MaxSeverity          int         `json:"maxSeverity"`
	MostCommonSeverity   int         `json:"mostCommonSeverity"`
	SeverityDistribution map[int]int `json:"severityDistribution"`
	TotalCount           int64       `json:"totalCount"`
}
