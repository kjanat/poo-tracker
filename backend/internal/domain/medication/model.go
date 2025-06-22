package medication

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Medication represents a medication entry with comprehensive tracking.
type Medication struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Basic medication information
	Name        string                     `json:"name"`
	GenericName string                     `json:"genericName,omitempty"`
	Brand       string                     `json:"brand,omitempty"`
	Category    *shared.MedicationCategory `json:"category,omitempty"`

	// Dosage information
	Dosage    string                  `json:"dosage"` // e.g., "10mg", "2 tablets"
	Form      *shared.MedicationForm  `json:"form,omitempty"`
	Frequency string                  `json:"frequency"` // e.g., "twice daily", "as needed"
	Route     *shared.MedicationRoute `json:"route,omitempty"`

	// Timing
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	TakenAt   *time.Time `json:"takenAt,omitempty"` // When last taken (for PRN medications)

	// Additional information
	Purpose     string   `json:"purpose,omitempty"`
	SideEffects []string `json:"sideEffects,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	PhotoURL    string   `json:"photoUrl,omitempty"`

	// Status tracking
	IsActive   bool `json:"isActive"`
	IsAsNeeded bool `json:"isAsNeeded"` // PRN (Pro Re Nata) medications
}

// MedicationUpdate represents fields that can be updated on a Medication.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
type MedicationUpdate struct {
	Name        *string                    `json:"name,omitempty"`
	GenericName *string                    `json:"genericName,omitempty"`
	Brand       *string                    `json:"brand,omitempty"`
	Category    *shared.MedicationCategory `json:"category,omitempty"`
	Dosage      *string                    `json:"dosage,omitempty"`
	Form        *shared.MedicationForm     `json:"form,omitempty"`
	Frequency   *string                    `json:"frequency,omitempty"`
	Route       *shared.MedicationRoute    `json:"route,omitempty"`
	StartDate   *time.Time                 `json:"startDate,omitempty"`
	EndDate     *time.Time                 `json:"endDate,omitempty"`
	TakenAt     *time.Time                 `json:"takenAt,omitempty"`
	Purpose     *string                    `json:"purpose,omitempty"`
	SideEffects []string                   `json:"sideEffects,omitempty"`
	Notes       *string                    `json:"notes,omitempty"`
	PhotoURL    *string                    `json:"photoUrl,omitempty"`
	IsActive    *bool                      `json:"isActive,omitempty"`
	IsAsNeeded  *bool                      `json:"isAsNeeded,omitempty"`
}

// NewMedication creates a new Medication with sensible defaults.
func NewMedication(userID, name, dosage, frequency string) Medication {
	now := time.Now()
	return Medication{
		UserID:      userID,
		Name:        name,
		Dosage:      dosage,
		Frequency:   frequency,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
		IsAsNeeded:  false,
		SideEffects: []string{},
	}
}
