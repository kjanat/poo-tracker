package model

import "time"

// Medication represents a medication entry with comprehensive tracking.
type Medication struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Basic medication information
	Name        string              `json:"name"`
	GenericName string              `json:"genericName,omitempty"`
	Brand       string              `json:"brand,omitempty"`
	Category    *MedicationCategory `json:"category,omitempty"`

	// Dosage information
	Dosage    string           `json:"dosage"` // e.g., "10mg", "2 tablets"
	Form      *MedicationForm  `json:"form,omitempty"`
	Frequency string           `json:"frequency"` // e.g., "twice daily", "as needed"
	Route     *MedicationRoute `json:"route,omitempty"`

	// Timing
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	TakenAt   *time.Time `json:"takenAt,omitempty"` // When last taken (for PRN medications)

	// Additional information
	Purpose    string `json:"purpose,omitempty"`
	Prescriber string `json:"prescriber,omitempty"`
	Pharmacy   string `json:"pharmacy,omitempty"`
	Notes      string `json:"notes,omitempty"`

	// Status
	IsActive bool `json:"isActive"`
	IsPRN    bool `json:"isPRN"` // As needed (pro re nata)

	// Side effects tracking
	SideEffects []string `json:"sideEffects,omitempty"`
}

// MedicationUpdate represents fields that can be updated on a Medication.
type MedicationUpdate struct {
	Name        *string             `json:"name,omitempty"`
	GenericName *string             `json:"genericName,omitempty"`
	Brand       *string             `json:"brand,omitempty"`
	Category    *MedicationCategory `json:"category,omitempty"`
	Dosage      *string             `json:"dosage,omitempty"`
	Form        *MedicationForm     `json:"form,omitempty"`
	Frequency   *string             `json:"frequency,omitempty"`
	Route       *MedicationRoute    `json:"route,omitempty"`
	StartDate   *time.Time          `json:"startDate,omitempty"`
	EndDate     *time.Time          `json:"endDate,omitempty"`
	TakenAt     *time.Time          `json:"takenAt,omitempty"`
	Purpose     *string             `json:"purpose,omitempty"`
	Prescriber  *string             `json:"prescriber,omitempty"`
	Pharmacy    *string             `json:"pharmacy,omitempty"`
	Notes       *string             `json:"notes,omitempty"`
	IsActive    *bool               `json:"isActive,omitempty"`
	IsPRN       *bool               `json:"isPRN,omitempty"`
	SideEffects []string            `json:"sideEffects,omitempty"`
}

// NewMedication creates a new Medication with sensible defaults.
func NewMedication(userID, name, dosage, frequency string) Medication {
	now := time.Now()
	return Medication{
		UserID:    userID,
		Name:      name,
		Dosage:    dosage,
		Frequency: frequency,
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
		IsPRN:     false,
	}
}
