package medication

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// CreateMedicationRequest represents the request to create a new medication
type CreateMedicationRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=100"`
	GenericName *string    `json:"generic_name,omitempty" binding:"omitempty,max=100"`
	Brand       *string    `json:"brand,omitempty" binding:"omitempty,max=100"`
	Category    *string    `json:"category,omitempty" binding:"omitempty,oneof=gastrointestinal pain_relief antibiotic probiotics supplements anti_inflammatory other"`
	Dosage      string     `json:"dosage" binding:"required,min=1,max=100"`
	Form        *string    `json:"form,omitempty" binding:"omitempty,oneof=tablet capsule liquid injection topical other"`
	Frequency   string     `json:"frequency" binding:"required,min=1,max=100"`
	Route       *string    `json:"route,omitempty" binding:"omitempty,oneof=oral topical injection sublingual other"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Purpose     *string    `json:"purpose,omitempty" binding:"omitempty,max=200"`
	SideEffects []string   `json:"side_effects,omitempty"`
	Notes       *string    `json:"notes,omitempty" binding:"omitempty,max=500"`
	PhotoURL    *string    `json:"photo_url,omitempty" binding:"omitempty,url"`
	IsAsNeeded  *bool      `json:"is_as_needed,omitempty"`
}

// UpdateMedicationRequest represents the request to update a medication
type UpdateMedicationRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	GenericName *string    `json:"generic_name,omitempty" binding:"omitempty,max=100"`
	Brand       *string    `json:"brand,omitempty" binding:"omitempty,max=100"`
	Category    *string    `json:"category,omitempty" binding:"omitempty,oneof=gastrointestinal pain_relief antibiotic probiotics supplements anti_inflammatory other"`
	Dosage      *string    `json:"dosage,omitempty" binding:"omitempty,min=1,max=100"`
	Form        *string    `json:"form,omitempty" binding:"omitempty,oneof=tablet capsule liquid injection topical other"`
	Frequency   *string    `json:"frequency,omitempty" binding:"omitempty,min=1,max=100"`
	Route       *string    `json:"route,omitempty" binding:"omitempty,oneof=oral topical injection sublingual other"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Purpose     *string    `json:"purpose,omitempty" binding:"omitempty,max=200"`
	SideEffects []string   `json:"side_effects,omitempty"`
	Notes       *string    `json:"notes,omitempty" binding:"omitempty,max=500"`
	PhotoURL    *string    `json:"photo_url,omitempty" binding:"omitempty,url"`
	IsActive    *bool      `json:"is_active,omitempty"`
	IsAsNeeded  *bool      `json:"is_as_needed,omitempty"`
}

// MedicationResponse represents a medication in API responses
type MedicationResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	GenericName string     `json:"generic_name,omitempty"`
	Brand       string     `json:"brand,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Dosage      string     `json:"dosage"`
	Form        *string    `json:"form,omitempty"`
	Frequency   string     `json:"frequency"`
	Route       *string    `json:"route,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	TakenAt     *time.Time `json:"taken_at,omitempty"`
	Purpose     string     `json:"purpose,omitempty"`
	SideEffects []string   `json:"side_effects,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	PhotoURL    string     `json:"photo_url,omitempty"`
	IsActive    bool       `json:"is_active"`
	IsAsNeeded  bool       `json:"is_as_needed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// MedicationListResponse represents a paginated list of medications
type MedicationListResponse struct {
	Medications []MedicationResponse `json:"medications"`
	TotalCount  int64                `json:"total_count"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	TotalPages  int                  `json:"total_pages"`
}

// MedicationSummaryResponse represents a simplified medication summary
type MedicationSummaryResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Dosage    string     `json:"dosage"`
	Frequency string     `json:"frequency"`
	IsActive  bool       `json:"is_active"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// CreateMedicationLogRequest represents the request to log a medication dose
type CreateMedicationLogRequest struct {
	MedicationID string    `json:"medication_id" binding:"required"`
	TakenAt      time.Time `json:"taken_at" binding:"required"`
	Dosage       *string   `json:"dosage,omitempty" binding:"omitempty,min=1,max=100"`
	Notes        *string   `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// MedicationLogResponse represents a medication log entry in API responses
type MedicationLogResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	MedicationID string    `json:"medication_id"`
	TakenAt      time.Time `json:"taken_at"`
	Dosage       string    `json:"dosage"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToMedicationResponse converts a domain Medication to MedicationResponse
func ToMedicationResponse(m *medication.Medication) MedicationResponse {
	response := MedicationResponse{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		GenericName: m.GenericName,
		Brand:       m.Brand,
		Dosage:      m.Dosage,
		Frequency:   m.Frequency,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		TakenAt:     m.TakenAt,
		Purpose:     m.Purpose,
		SideEffects: m.SideEffects,
		Notes:       m.Notes,
		PhotoURL:    m.PhotoURL,
		IsActive:    m.IsActive,
		IsAsNeeded:  m.IsAsNeeded,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	// Convert category if present
	if m.Category != nil {
		categoryStr := string(*m.Category)
		response.Category = &categoryStr
	}

	// Convert form if present
	if m.Form != nil {
		formStr := string(*m.Form)
		response.Form = &formStr
	}

	// Convert route if present
	if m.Route != nil {
		routeStr := string(*m.Route)
		response.Route = &routeStr
	}

	return response
}

// ToMedicationListResponse converts a slice of domain Medications to MedicationListResponse
func ToMedicationListResponse(medications []medication.Medication, totalCount int64, page, pageSize int) MedicationListResponse {
	medicationRes := make([]MedicationResponse, len(medications))
	for i, m := range medications {
		medicationRes[i] = ToMedicationResponse(&m)
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return MedicationListResponse{
		Medications: medicationRes,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}
}

// ToMedicationSummaryResponse converts a domain Medication to MedicationSummaryResponse
func ToMedicationSummaryResponse(m *medication.Medication) MedicationSummaryResponse {
	return MedicationSummaryResponse{
		ID:        m.ID,
		Name:      m.Name,
		Dosage:    m.Dosage,
		Frequency: m.Frequency,
		IsActive:  m.IsActive,
		EndDate:   m.EndDate,
	}
}

// ToDomainMedication converts CreateMedicationRequest to domain Medication
func (r *CreateMedicationRequest) ToDomainMedication(userID string) *medication.Medication {
	m := medication.NewMedication(userID, r.Name, r.Dosage, r.Frequency)

	// Set optional string fields
	if r.GenericName != nil {
		m.GenericName = *r.GenericName
	}
	if r.Brand != nil {
		m.Brand = *r.Brand
	}
	if r.Purpose != nil {
		m.Purpose = *r.Purpose
	}
	if r.Notes != nil {
		m.Notes = *r.Notes
	}
	if r.PhotoURL != nil {
		m.PhotoURL = *r.PhotoURL
	}

	// Set optional time fields
	if r.StartDate != nil {
		m.StartDate = r.StartDate
	}
	if r.EndDate != nil {
		m.EndDate = r.EndDate
	}

	// Set optional slice fields
	if r.SideEffects != nil {
		m.SideEffects = r.SideEffects
	}

	// Set optional bool fields
	if r.IsAsNeeded != nil {
		m.IsAsNeeded = *r.IsAsNeeded
	}

	// Set optional enum fields
	if r.Category != nil {
		category := shared.MedicationCategory(*r.Category)
		m.Category = &category
	}
	if r.Form != nil {
		form := shared.MedicationForm(*r.Form)
		m.Form = &form
	}
	if r.Route != nil {
		route := shared.MedicationRoute(*r.Route)
		m.Route = &route
	}

	return &m
}

// ApplyToDomainMedication applies UpdateMedicationRequest to a domain Medication
func (r *UpdateMedicationRequest) ApplyToDomainMedication(m *medication.Medication) {
	if r.Name != nil {
		m.Name = *r.Name
	}
	if r.GenericName != nil {
		m.GenericName = *r.GenericName
	}
	if r.Brand != nil {
		m.Brand = *r.Brand
	}
	if r.Dosage != nil {
		m.Dosage = *r.Dosage
	}
	if r.Frequency != nil {
		m.Frequency = *r.Frequency
	}
	if r.StartDate != nil {
		m.StartDate = r.StartDate
	}
	if r.EndDate != nil {
		m.EndDate = r.EndDate
	}
	if r.Purpose != nil {
		m.Purpose = *r.Purpose
	}
	if r.SideEffects != nil {
		m.SideEffects = r.SideEffects
	}
	if r.Notes != nil {
		m.Notes = *r.Notes
	}
	if r.PhotoURL != nil {
		m.PhotoURL = *r.PhotoURL
	}
	if r.IsActive != nil {
		m.IsActive = *r.IsActive
	}
	if r.IsAsNeeded != nil {
		m.IsAsNeeded = *r.IsAsNeeded
	}

	// Handle enum fields
	if r.Category != nil {
		category := shared.MedicationCategory(*r.Category)
		m.Category = &category
	}
	if r.Form != nil {
		form := shared.MedicationForm(*r.Form)
		m.Form = &form
	}
	if r.Route != nil {
		route := shared.MedicationRoute(*r.Route)
		m.Route = &route
	}
}

// Validate validates the CreateMedicationRequest
func (r *CreateMedicationRequest) Validate() error {
	if len(r.Name) == 0 || len(r.Name) > 100 {
		return medication.ErrInvalidMedicationName
	}
	if len(r.Dosage) == 0 || len(r.Dosage) > 100 {
		return medication.ErrInvalidDosage
	}
	if len(r.Frequency) == 0 || len(r.Frequency) > 100 {
		return medication.ErrInvalidFrequency
	}
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return medication.ErrInvalidDateRange
	}
	return nil
}

// Validate validates the UpdateMedicationRequest
func (r *UpdateMedicationRequest) Validate() error {
	if r.Name != nil && (len(*r.Name) == 0 || len(*r.Name) > 100) {
		return medication.ErrInvalidMedicationName
	}
	if r.Dosage != nil && (len(*r.Dosage) == 0 || len(*r.Dosage) > 100) {
		return medication.ErrInvalidDosage
	}
	if r.Frequency != nil && (len(*r.Frequency) == 0 || len(*r.Frequency) > 100) {
		return medication.ErrInvalidFrequency
	}
	if r.StartDate != nil && r.EndDate != nil && r.EndDate.Before(*r.StartDate) {
		return medication.ErrInvalidDateRange
	}
	return nil
}

// Validate validates the CreateMedicationLogRequest
func (r *CreateMedicationLogRequest) Validate() error {
	if len(r.MedicationID) == 0 {
		return medication.ErrInvalidDoseTime
	}
	return nil
}
