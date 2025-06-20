package bowelmovement

import (
	"context"
	"time"
)

// Repository defines the interface for bowel movement data persistence
type Repository interface {
	// BowelMovement CRUD operations
	Create(ctx context.Context, bm *BowelMovement) error
	GetByID(ctx context.Context, id string) (*BowelMovement, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*BowelMovement, error)
	Update(ctx context.Context, id string, update *BowelMovementUpdate) error
	Delete(ctx context.Context, id string) error

	// BowelMovementDetails operations
	CreateDetails(ctx context.Context, details *BowelMovementDetails) error
	GetDetailsByBowelMovementID(ctx context.Context, bowelMovementID string) (*BowelMovementDetails, error)
	UpdateDetails(ctx context.Context, bowelMovementID string, update *BowelMovementDetailsUpdate) error
	DeleteDetails(ctx context.Context, bowelMovementID string) error

	// Query operations
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*BowelMovement, error)
	GetCountByUserID(ctx context.Context, userID string) (int64, error)
	GetLatestByUserID(ctx context.Context, userID string) (*BowelMovement, error)

	// Analytics operations
	GetAveragesByUserID(ctx context.Context, userID string, start, end time.Time) (*BowelMovementAverages, error)
	GetFrequencyByUserID(ctx context.Context, userID string, start, end time.Time) (map[string]int, error)
}

// BowelMovementAverages represents calculated averages for analytics
type BowelMovementAverages struct {
	AveragePain         float64 `json:"averagePain"`
	AverageStrain       float64 `json:"averageStrain"`
	AverageSatisfaction float64 `json:"averageSatisfaction"`
	MostCommonBristol   int     `json:"mostCommonBristol"`
	TotalCount          int64   `json:"totalCount"`
}
