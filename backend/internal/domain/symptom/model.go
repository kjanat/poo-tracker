package symptom

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Symptom represents a symptom entry with comprehensive tracking.
type Symptom struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Basic symptom information
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	RecordedAt  time.Time               `json:"recordedAt"` // User-specified time
	Category    *shared.SymptomCategory `json:"category,omitempty"`
	Severity    int                     `json:"severity"`           // 1-10 scale
	Duration    *int                    `json:"duration,omitempty"` // Duration in minutes

	// Location and characteristics
	BodyPart string              `json:"bodyPart,omitempty"`
	Type     *shared.SymptomType `json:"type,omitempty"`
	Triggers []string            `json:"triggers,omitempty"`

	// Optional fields
	Notes    string `json:"notes,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// SymptomUpdate represents fields that can be updated on a Symptom.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
type SymptomUpdate struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	RecordedAt  *time.Time              `json:"recordedAt,omitempty"`
	Category    *shared.SymptomCategory `json:"category,omitempty"`
	Severity    *int                    `json:"severity,omitempty"`
	Duration    *int                    `json:"duration,omitempty"`
	BodyPart    *string                 `json:"bodyPart,omitempty"`
	Type        *shared.SymptomType     `json:"type,omitempty"`
	Triggers    []string                `json:"triggers,omitempty"`
	Notes       *string                 `json:"notes,omitempty"`
	PhotoURL    *string                 `json:"photoUrl,omitempty"`
}

// NewSymptom creates a new Symptom with sensible defaults.
func NewSymptom(userID, name string, severity int, recordedAt time.Time) Symptom {
	now := time.Now()
	return Symptom{
		UserID:     userID,
		Name:       name,
		Severity:   severity,
		RecordedAt: recordedAt,
		CreatedAt:  now,
		UpdatedAt:  now,
		Triggers:   []string{},
	}
}
