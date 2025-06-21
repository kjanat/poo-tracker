// Package bowelmovement provides domain models for bowel movement tracking.
// This file contains the core BowelMovement model and related data structures
// for comprehensive digestive health monitoring.
//
// Key Models:
// - BowelMovement: Core bowel movement tracking with Bristol stool scale
// - BowelMovementDetails: Extended metadata and detailed descriptions
// - BowelMovementUpdate: Update operations with pointer fields for partial updates
//
// Features:
// - Bristol Stool Scale classification (1-7)
// - Pain, strain, and satisfaction metrics (1-10 scale)
// - Duration tracking and urgency levels
// - Photo attachments and detailed notes
// - User-specified timing vs system timestamps
// - Performance-optimized with separate details storage
package bowelmovement

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// BowelMovement represents a bowel movement entry with comprehensive tracking.
// Large text fields and detailed metadata are stored in BowelMovementDetails for performance.
type BowelMovement struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	RecordedAt time.Time `json:"recordedAt"` // User-specified time in their timezone

	// Bristol Stool Chart (1-7)
	BristolType int `json:"bristolType"`

	// Physical characteristics
	Volume      *shared.Volume      `json:"volume,omitempty"`
	Color       *shared.Color       `json:"color,omitempty"`
	Consistency *shared.Consistency `json:"consistency,omitempty"`
	Floaters    bool                `json:"floaters"`

	// Experience (1-10 scales)
	Pain         int `json:"pain"`         // 1-10 scale, default 1
	Strain       int `json:"strain"`       // 1-10 scale, default 1
	Satisfaction int `json:"satisfaction"` // 1-10 scale, default 5

	// Optional fields (light data)
	PhotoURL   string             `json:"photoUrl,omitempty"`
	SmellLevel *shared.SmellLevel `json:"smellLevel,omitempty"`

	// Reference to details (loaded separately)
	HasDetails bool `json:"hasDetails"` // Indicates if details exist for this bowel movement
}

// BowelMovementDetails represents detailed information stored separately for performance.
// This includes large text fields, AI analysis, and other detailed metadata.
type BowelMovementDetails struct {
	ID              string    `json:"id"`
	BowelMovementID string    `json:"bowelMovementId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Large text fields (stored separately for performance)
	Notes          string `json:"notes,omitempty"`
	DetailedNotes  string `json:"detailedNotes,omitempty"`  // More extensive notes
	Environment    string `json:"environment,omitempty"`    // Bathroom conditions, privacy, etc.
	PreConditions  string `json:"preConditions,omitempty"`  // What led to this movement
	PostConditions string `json:"postConditions,omitempty"` // How you felt after

	// AI analysis and metadata
	AIAnalysis        interface{} `json:"aiAnalysis,omitempty"`        // JSON field for AI analysis
	AIConfidence      *float64    `json:"aiConfidence,omitempty"`      // 0.0-1.0 confidence score
	AIRecommendations string      `json:"aiRecommendations,omitempty"` // AI-generated recommendations

	// Extended metadata
	Tags              []string `json:"tags,omitempty"`              // User-defined tags
	WeatherCondition  string   `json:"weatherCondition,omitempty"`  // Weather at time of movement
	StressLevel       *int     `json:"stressLevel,omitempty"`       // 1-10 scale
	SleepQuality      *int     `json:"sleepQuality,omitempty"`      // 1-10 scale (night before)
	ExerciseIntensity *int     `json:"exerciseIntensity,omitempty"` // 1-10 scale (day before)
}

// BowelMovementUpdate represents fields that can be updated on a BowelMovement.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
// Note: Notes and detailed fields are updated via BowelMovementDetails.
type BowelMovementUpdate struct {
	BristolType  *int                `json:"bristolType,omitempty"`
	Volume       *shared.Volume      `json:"volume,omitempty"`
	Color        *shared.Color       `json:"color,omitempty"`
	Consistency  *shared.Consistency `json:"consistency,omitempty"`
	Floaters     *bool               `json:"floaters,omitempty"`
	Pain         *int                `json:"pain,omitempty"`
	Strain       *int                `json:"strain,omitempty"`
	Satisfaction *int                `json:"satisfaction,omitempty"`
	PhotoURL     *string             `json:"photoUrl,omitempty"`
	SmellLevel   *shared.SmellLevel  `json:"smellLevel,omitempty"`
	RecordedAt   *time.Time          `json:"recordedAt,omitempty"`
}

// BowelMovementDetailsUpdate represents fields that can be updated on BowelMovementDetails.
type BowelMovementDetailsUpdate struct {
	Notes             *string     `json:"notes,omitempty"`
	DetailedNotes     *string     `json:"detailedNotes,omitempty"`
	Environment       *string     `json:"environment,omitempty"`
	PreConditions     *string     `json:"preConditions,omitempty"`
	PostConditions    *string     `json:"postConditions,omitempty"`
	AIAnalysis        interface{} `json:"aiAnalysis,omitempty"`
	AIConfidence      *float64    `json:"aiConfidence,omitempty"`
	AIRecommendations *string     `json:"aiRecommendations,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	WeatherCondition  *string     `json:"weatherCondition,omitempty"`
	StressLevel       *int        `json:"stressLevel,omitempty"`
	SleepQuality      *int        `json:"sleepQuality,omitempty"`
	ExerciseIntensity *int        `json:"exerciseIntensity,omitempty"`
}

// NewBowelMovement creates a new BowelMovement with sensible defaults.
func NewBowelMovement(userID string, bristolType int) (BowelMovement, error) {
	if bristolType < 1 || bristolType > 7 {
		return BowelMovement{}, ErrInvalidBristolType
	}

	now := time.Now()
	bm := BowelMovement{
		UserID:       userID,
		BristolType:  bristolType,
		CreatedAt:    now,
		UpdatedAt:    now,
		RecordedAt:   now,
		Pain:         1, // Default: minimal pain
		Strain:       1, // Default: minimal strain
		Satisfaction: 5, // Default: neutral satisfaction
		Floaters:     false,
	}
	return bm, nil
}

// NewBowelMovementDetails creates a new BowelMovementDetails with defaults.
func NewBowelMovementDetails(bowelMovementID string) BowelMovementDetails {
	now := time.Now()
	return BowelMovementDetails{
		BowelMovementID: bowelMovementID,
		CreatedAt:       now,
		UpdatedAt:       now,
		Tags:            []string{},
	}
}
