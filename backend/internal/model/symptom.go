package model

import "time"

// Symptom represents a symptom entry with comprehensive tracking.
type Symptom struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Basic symptom information
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	RecordedAt  time.Time        `json:"recordedAt"` // User-specified time
	Category    *SymptomCategory `json:"category,omitempty"`
	Severity    int              `json:"severity"`           // 1-10 scale
	Duration    *int             `json:"duration,omitempty"` // Duration in minutes

	// Location and characteristics
	BodyPart string       `json:"bodyPart,omitempty"`
	Type     *SymptomType `json:"type,omitempty"`
	Triggers []string     `json:"triggers,omitempty"`

	// Optional fields
	Notes    string `json:"notes,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// SymptomUpdate represents fields that can be updated on a Symptom.
type SymptomUpdate struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	RecordedAt  *time.Time       `json:"recordedAt,omitempty"`
	Category    *SymptomCategory `json:"category,omitempty"`
	Severity    *int             `json:"severity,omitempty"`
	Duration    *int             `json:"duration,omitempty"`
	BodyPart    *string          `json:"bodyPart,omitempty"`
	Type        *SymptomType     `json:"type,omitempty"`
	Triggers    []string         `json:"triggers,omitempty"`
	Notes       *string          `json:"notes,omitempty"`
	PhotoURL    *string          `json:"photoUrl,omitempty"`
}

// NewSymptom creates a new Symptom with sensible defaults.
func NewSymptom(userID, name string) Symptom {
	now := time.Now()
	return Symptom{
		UserID:     userID,
		Name:       name,
		CreatedAt:  now,
		UpdatedAt:  now,
		RecordedAt: now,
		Severity:   1,
	}
}
