package medication

import (
	"context"
	"time"
)

// Repository defines the interface for medication data persistence
type Repository interface {
	// Medication CRUD operations
	Create(ctx context.Context, medication *Medication) error
	GetByID(ctx context.Context, id string) (*Medication, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Medication, error)
	Update(ctx context.Context, id string, update *MedicationUpdate) error
	Delete(ctx context.Context, id string) error

	// Query operations
	GetActiveByUserID(ctx context.Context, userID string) ([]*Medication, error)
	GetByCategory(ctx context.Context, userID string, category string) ([]*Medication, error)
	GetAsNeededByUserID(ctx context.Context, userID string) ([]*Medication, error)
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*Medication, error)
	GetCountByUserID(ctx context.Context, userID string) (int64, error)

	// Tracking operations
	RecordDose(ctx context.Context, medicationID string, takenAt time.Time) error
	GetDoseHistory(ctx context.Context, medicationID string, start, end time.Time) ([]*DoseRecord, error)

	// Analytics operations
	GetUsageStats(ctx context.Context, userID string, start, end time.Time) (*UsageStats, error)
	GetCategoryBreakdown(ctx context.Context, userID string) (map[string]int, error)
}

// DoseRecord represents a record of when a medication was taken
type DoseRecord struct {
	ID           string    `json:"id"`
	MedicationID string    `json:"medicationId"`
	TakenAt      time.Time `json:"takenAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

// UsageStats represents medication usage statistics
type UsageStats struct {
	TotalMedications  int64          `json:"totalMedications"`
	ActiveMedications int64          `json:"activeMedications"`
	CategoryBreakdown map[string]int `json:"categoryBreakdown"`
	FormBreakdown     map[string]int `json:"formBreakdown"`
	RouteBreakdown    map[string]int `json:"routeBreakdown"`
	AsNeededCount     int64          `json:"asNeededCount"`
	RegularCount      int64          `json:"regularCount"`
}
