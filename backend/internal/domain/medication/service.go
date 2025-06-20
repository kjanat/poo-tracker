package medication

import (
	"context"
	"time"
)

// Service defines the interface for medication business logic
type Service interface {
	// Core operations
	Create(ctx context.Context, userID string, input *CreateMedicationInput) (*Medication, error)
	GetByID(ctx context.Context, id string) (*Medication, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Medication, error)
	Update(ctx context.Context, id string, input *UpdateMedicationInput) (*Medication, error)
	Delete(ctx context.Context, id string) error

	// Medication management
	GetActiveMedications(ctx context.Context, userID string) ([]*Medication, error)
	GetAsNeededMedications(ctx context.Context, userID string) ([]*Medication, error)
	DeactivateMedication(ctx context.Context, id string) error
	ReactivateMedication(ctx context.Context, id string) error

	// Dose tracking
	RecordDose(ctx context.Context, medicationID string, input *RecordDoseInput) error
	GetDoseHistory(ctx context.Context, medicationID string, start, end time.Time) ([]*DoseRecord, error)

	// Analytics operations
	GetMedicationStats(ctx context.Context, userID string, start, end time.Time) (*MedicationStats, error)
	GetMedicationInsights(ctx context.Context, userID string, start, end time.Time) (*MedicationInsights, error)
	GetSideEffectAnalysis(ctx context.Context, userID string) (*SideEffectAnalysis, error)
}

// CreateMedicationInput represents input for creating a medication
type CreateMedicationInput struct {
	Name        string     `json:"name" binding:"required,min=1,max=200"`
	GenericName string     `json:"genericName,omitempty"`
	Brand       string     `json:"brand,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Dosage      string     `json:"dosage" binding:"required,min=1,max=100"`
	Form        *string    `json:"form,omitempty"`
	Frequency   string     `json:"frequency" binding:"required,min=1,max=200"`
	Route       *string    `json:"route,omitempty"`
	StartDate   *time.Time `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Purpose     string     `json:"purpose,omitempty"`
	SideEffects []string   `json:"sideEffects,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	PhotoURL    string     `json:"photoUrl,omitempty"`
	IsAsNeeded  bool       `json:"isAsNeeded"`
}

// UpdateMedicationInput represents input for updating a medication
type UpdateMedicationInput struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	GenericName *string    `json:"genericName,omitempty"`
	Brand       *string    `json:"brand,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Dosage      *string    `json:"dosage,omitempty" binding:"omitempty,min=1,max=100"`
	Form        *string    `json:"form,omitempty"`
	Frequency   *string    `json:"frequency,omitempty" binding:"omitempty,min=1,max=200"`
	Route       *string    `json:"route,omitempty"`
	StartDate   *time.Time `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Purpose     *string    `json:"purpose,omitempty"`
	SideEffects []string   `json:"sideEffects,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
	PhotoURL    *string    `json:"photoUrl,omitempty"`
	IsAsNeeded  *bool      `json:"isAsNeeded,omitempty"`
}

// RecordDoseInput represents input for recording a dose
type RecordDoseInput struct {
	TakenAt time.Time `json:"takenAt" binding:"required"`
	Notes   string    `json:"notes,omitempty"`
}

// MedicationStats represents medication usage statistics
type MedicationStats struct {
	TotalMedications    int64          `json:"totalMedications"`
	ActiveMedications   int64          `json:"activeMedications"`
	AsNeededMedications int64          `json:"asNeededMedications"`
	RegularMedications  int64          `json:"regularMedications"`
	CategoryBreakdown   map[string]int `json:"categoryBreakdown"`
	FormBreakdown       map[string]int `json:"formBreakdown"`
	RouteBreakdown      map[string]int `json:"routeBreakdown"`
	DosesThisPeriod     int64          `json:"dosesThisPeriod"`
	AverageDoesPerDay   float64        `json:"averageDosesPerDay"`
}

// MedicationInsights represents insights from medication data
type MedicationInsights struct {
	MostUsedCategory string  `json:"mostUsedCategory"`
	MostCommonForm   string  `json:"mostCommonForm"`
	MostCommonRoute  string  `json:"mostCommonRoute"`
	AdhereanceScore  float64 `json:"adherenceScore"`   // 0-1 based on regular medication consistency
	ComplexityScore  float64 `json:"complexityScore"`  // 0-1 based on number of medications and frequency
	MedicationBurden int     `json:"medicationBurden"` // Total number of active medications
	InteractionRisk  string  `json:"interactionRisk"`  // LOW, MEDIUM, HIGH
}

// SideEffectAnalysis represents analysis of reported side effects
type SideEffectAnalysis struct {
	MostCommonSideEffects   []string            `json:"mostCommonSideEffects"`
	SideEffectFrequency     map[string]int      `json:"sideEffectFrequency"`
	MedicationSideEffectMap map[string][]string `json:"medicationSideEffectMap"`
	TotalUniqueSideEffects  int                 `json:"totalUniqueSideEffects"`
}
