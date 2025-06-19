package model

import "time"

// BowelMovement represents a bowel movement entry with comprehensive tracking.
type BowelMovement struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	RecordedAt time.Time `json:"recordedAt"` // User-specified time in their timezone

	// Bristol Stool Chart (1-7)
	BristolType int `json:"bristolType"`

	// Physical characteristics
	Volume      *Volume      `json:"volume,omitempty"`
	Color       *Color       `json:"color,omitempty"`
	Consistency *Consistency `json:"consistency,omitempty"`
	Floaters    bool         `json:"floaters"`

	// Experience (1-10 scales)
	Pain         int `json:"pain"`         // 1-10 scale, default 1
	Strain       int `json:"strain"`       // 1-10 scale, default 1
	Satisfaction int `json:"satisfaction"` // 1-10 scale, default 5

	// Optional fields
	PhotoURL   string      `json:"photoUrl,omitempty"`
	SmellLevel *SmellLevel `json:"smellLevel,omitempty"`
	Notes      string      `json:"notes,omitempty"`
}

// BowelMovementDetails represents detailed information stored separately for performance.
type BowelMovementDetails struct {
	ID              string      `json:"id"`
	BowelMovementID string      `json:"bowelMovementId"`
	Notes           string      `json:"notes,omitempty"`
	AIAnalysis      interface{} `json:"aiAnalysis,omitempty"` // JSON field for AI analysis
}

// BowelMovementUpdate represents fields that can be updated on a BowelMovement.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
type BowelMovementUpdate struct {
	BristolType  *int         `json:"bristolType,omitempty"`
	Volume       *Volume      `json:"volume,omitempty"`
	Color        *Color       `json:"color,omitempty"`
	Consistency  *Consistency `json:"consistency,omitempty"`
	Floaters     *bool        `json:"floaters,omitempty"`
	Pain         *int         `json:"pain,omitempty"`
	Strain       *int         `json:"strain,omitempty"`
	Satisfaction *int         `json:"satisfaction,omitempty"`
	PhotoURL     *string      `json:"photoUrl,omitempty"`
	SmellLevel   *SmellLevel  `json:"smellLevel,omitempty"`
	Notes        *string      `json:"notes,omitempty"`
	RecordedAt   *time.Time   `json:"recordedAt,omitempty"`
}

// NewBowelMovement creates a new BowelMovement with sensible defaults.
func NewBowelMovement(userID string, bristolType int) BowelMovement {
	now := time.Now()
	return BowelMovement{
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
}
