package bowelmovement

import (
	"context"
	"time"
)

// Service defines the interface for bowel movement business logic
type Service interface {
	// Core operations
	Create(ctx context.Context, userID string, input *CreateBowelMovementInput) (*BowelMovement, error)
	GetByID(ctx context.Context, id string) (*BowelMovement, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*BowelMovement, error)
	Update(ctx context.Context, id string, input *UpdateBowelMovementInput) (*BowelMovement, error)
	Delete(ctx context.Context, id string) error

	// Details operations
	CreateDetails(ctx context.Context, bowelMovementID string, input *CreateBowelMovementDetailsInput) (*BowelMovementDetails, error)
	GetDetails(ctx context.Context, bowelMovementID string) (*BowelMovementDetails, error)
	UpdateDetails(ctx context.Context, bowelMovementID string, input *UpdateBowelMovementDetailsInput) (*BowelMovementDetails, error)

	// Analytics and queries
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*BowelMovement, error)
	GetUserStats(ctx context.Context, userID string, start, end time.Time) (*UserBowelMovementStats, error)
	GetLatest(ctx context.Context, userID string) (*BowelMovement, error)
}

// CreateBowelMovementInput represents input for creating a bowel movement
type CreateBowelMovementInput struct {
	BristolType  int       `json:"bristolType" binding:"required,min=1,max=7"`
	RecordedAt   time.Time `json:"recordedAt"`
	Volume       *string   `json:"volume,omitempty"`
	Color        *string   `json:"color,omitempty"`
	Consistency  *string   `json:"consistency,omitempty"`
	Floaters     bool      `json:"floaters"`
	Pain         int       `json:"pain" binding:"min=0,max=10"`
	Strain       int       `json:"strain" binding:"min=0,max=10"`
	Satisfaction int       `json:"satisfaction" binding:"min=0,max=10"`
	PhotoURL     string    `json:"photoUrl,omitempty"`
	SmellLevel   *string   `json:"smellLevel,omitempty"`
}

// UpdateBowelMovementInput represents input for updating a bowel movement
type UpdateBowelMovementInput struct {
	BristolType  *int       `json:"bristolType,omitempty" binding:"omitempty,min=1,max=7"`
	RecordedAt   *time.Time `json:"recordedAt,omitempty"`
	Volume       *string    `json:"volume,omitempty"`
	Color        *string    `json:"color,omitempty"`
	Consistency  *string    `json:"consistency,omitempty"`
	Floaters     *bool      `json:"floaters,omitempty"`
	Pain         *int       `json:"pain,omitempty" binding:"omitempty,min=0,max=10"`
	Strain       *int       `json:"strain,omitempty" binding:"omitempty,min=0,max=10"`
	Satisfaction *int       `json:"satisfaction,omitempty" binding:"omitempty,min=0,max=10"`
	PhotoURL     *string    `json:"photoUrl,omitempty"`
	SmellLevel   *string    `json:"smellLevel,omitempty"`
}

// CreateBowelMovementDetailsInput represents input for creating bowel movement details
type CreateBowelMovementDetailsInput struct {
	Notes             string   `json:"notes,omitempty"`
	DetailedNotes     string   `json:"detailedNotes,omitempty"`
	Environment       string   `json:"environment,omitempty"`
	PreConditions     string   `json:"preConditions,omitempty"`
	PostConditions    string   `json:"postConditions,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  string   `json:"weatherCondition,omitempty"`
	StressLevel       *int     `json:"stressLevel,omitempty" binding:"omitempty,min=0,max=10"`
	SleepQuality      *int     `json:"sleepQuality,omitempty" binding:"omitempty,min=0,max=10"`
	ExerciseIntensity *int     `json:"exerciseIntensity,omitempty" binding:"omitempty,min=0,max=10"`
}

// UpdateBowelMovementDetailsInput represents input for updating bowel movement details
type UpdateBowelMovementDetailsInput struct {
	Notes             *string  `json:"notes,omitempty"`
	DetailedNotes     *string  `json:"detailedNotes,omitempty"`
	Environment       *string  `json:"environment,omitempty"`
	PreConditions     *string  `json:"preConditions,omitempty"`
	PostConditions    *string  `json:"postConditions,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  *string  `json:"weatherCondition,omitempty"`
	StressLevel       *int     `json:"stressLevel,omitempty" binding:"omitempty,min=0,max=10"`
	SleepQuality      *int     `json:"sleepQuality,omitempty" binding:"omitempty,min=0,max=10"`
	ExerciseIntensity *int     `json:"exerciseIntensity,omitempty" binding:"omitempty,min=0,max=10"`
}

// UserBowelMovementStats represents analytics data for a user
type UserBowelMovementStats struct {
	TotalCount          int64       `json:"totalCount"`
	AveragePain         float64     `json:"averagePain"`
	AverageStrain       float64     `json:"averageStrain"`
	AverageSatisfaction float64     `json:"averageSatisfaction"`
	MostCommonBristol   int         `json:"mostCommonBristol"`
	FrequencyPerDay     float64     `json:"frequencyPerDay"`
	BristolDistribution map[int]int `json:"bristolDistribution"`
}
