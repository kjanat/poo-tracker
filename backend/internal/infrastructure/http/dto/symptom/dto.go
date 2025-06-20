package symptom

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// CreateSymptomRequest represents the request to create a new symptom
type CreateSymptomRequest struct {
	Name        string    `json:"name" binding:"required,min=1,max=100"`
	Description *string   `json:"description,omitempty" binding:"omitempty,max=500"`
	RecordedAt  time.Time `json:"recorded_at" binding:"required"`
	Category    *string   `json:"category,omitempty" binding:"omitempty,oneof=digestive neurological physical emotional"`
	Severity    int       `json:"severity" binding:"required,min=1,max=10"`
	Duration    *int      `json:"duration,omitempty" binding:"omitempty,min=1"` // minutes
	BodyPart    *string   `json:"body_part,omitempty" binding:"omitempty,max=100"`
	Type        *string   `json:"type,omitempty" binding:"omitempty,oneof=abdominal_pain bloating nausea diarrhea constipation heartburn fatigue headache mood_change other"`
	Triggers    []string  `json:"triggers,omitempty"`
	Notes       *string   `json:"notes,omitempty" binding:"omitempty,max=500"`
	PhotoURL    *string   `json:"photo_url,omitempty" binding:"omitempty,url"`
}

// UpdateSymptomRequest represents the request to update a symptom
type UpdateSymptomRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=500"`
	RecordedAt  *time.Time `json:"recorded_at,omitempty"`
	Category    *string    `json:"category,omitempty" binding:"omitempty,oneof=digestive neurological physical emotional"`
	Severity    *int       `json:"severity,omitempty" binding:"omitempty,min=1,max=10"`
	Duration    *int       `json:"duration,omitempty" binding:"omitempty,min=1"` // minutes
	BodyPart    *string    `json:"body_part,omitempty" binding:"omitempty,max=100"`
	Type        *string    `json:"type,omitempty" binding:"omitempty,oneof=abdominal_pain bloating nausea diarrhea constipation heartburn fatigue headache mood_change other"`
	Triggers    []string   `json:"triggers,omitempty"`
	Notes       *string    `json:"notes,omitempty" binding:"omitempty,max=500"`
	PhotoURL    *string    `json:"photo_url,omitempty" binding:"omitempty,url"`
}

// SymptomResponse represents a symptom in API responses
type SymptomResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	RecordedAt  time.Time `json:"recorded_at"`
	Category    *string   `json:"category,omitempty"`
	Severity    int       `json:"severity"`
	Duration    *int      `json:"duration,omitempty"` // minutes
	BodyPart    string    `json:"body_part,omitempty"`
	Type        *string   `json:"type,omitempty"`
	Triggers    []string  `json:"triggers,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SymptomListResponse represents a paginated list of symptoms
type SymptomListResponse struct {
	Symptoms   []SymptomResponse `json:"symptoms"`
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// SymptomSummaryResponse represents a simplified symptom summary
type SymptomSummaryResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Severity   int       `json:"severity"`
	RecordedAt time.Time `json:"recorded_at"`
	Duration   *int      `json:"duration,omitempty"`
}

// SymptomFrequencyResponse represents symptom frequency data for analytics
type SymptomFrequencyResponse struct {
	Type         string    `json:"type"`
	Count        int       `json:"count"`
	AvgSeverity  float64   `json:"avg_severity"`
	LastOccurred time.Time `json:"last_occurred"`
}

// ToSymptomResponse converts a domain Symptom to SymptomResponse
func ToSymptomResponse(s *symptom.Symptom) SymptomResponse {
	response := SymptomResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		Name:        s.Name,
		Description: s.Description,
		RecordedAt:  s.RecordedAt,
		Severity:    s.Severity,
		Duration:    s.Duration,
		BodyPart:    s.BodyPart,
		Triggers:    s.Triggers,
		Notes:       s.Notes,
		PhotoURL:    s.PhotoURL,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	// Convert category if present
	if s.Category != nil {
		categoryStr := string(*s.Category)
		response.Category = &categoryStr
	}

	// Convert type if present
	if s.Type != nil {
		typeStr := string(*s.Type)
		response.Type = &typeStr
	}

	return response
}

// ToSymptomListResponse converts a slice of domain Symptoms to SymptomListResponse
func ToSymptomListResponse(symptoms []symptom.Symptom, totalCount int64, page, pageSize int) SymptomListResponse {
	symptomRes := make([]SymptomResponse, len(symptoms))
	for i, s := range symptoms {
		symptomRes[i] = ToSymptomResponse(&s)
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return SymptomListResponse{
		Symptoms:   symptomRes,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToSymptomSummaryResponse converts a domain Symptom to SymptomSummaryResponse
func ToSymptomSummaryResponse(s *symptom.Symptom) SymptomSummaryResponse {
	summary := SymptomSummaryResponse{
		ID:         s.ID,
		Name:       s.Name,
		Severity:   s.Severity,
		RecordedAt: s.RecordedAt,
		Duration:   s.Duration,
	}

	return summary
}

// ToDomainSymptom converts CreateSymptomRequest to domain Symptom
func (r *CreateSymptomRequest) ToDomainSymptom(userID string) *symptom.Symptom {
	s := &symptom.Symptom{
		UserID:     userID,
		Name:       r.Name,
		RecordedAt: r.RecordedAt,
		Severity:   r.Severity,
		Duration:   r.Duration,
		Triggers:   r.Triggers,
	}

	// Set optional string fields
	if r.Description != nil {
		s.Description = *r.Description
	}
	if r.BodyPart != nil {
		s.BodyPart = *r.BodyPart
	}
	if r.Notes != nil {
		s.Notes = *r.Notes
	}
	if r.PhotoURL != nil {
		s.PhotoURL = *r.PhotoURL
	}

	// Set optional enum fields
	if r.Category != nil {
		category := shared.SymptomCategory(*r.Category)
		s.Category = &category
	}
	if r.Type != nil {
		symptomType := shared.SymptomType(*r.Type)
		s.Type = &symptomType
	}

	return s
}

// ApplyToDomainSymptom applies UpdateSymptomRequest to a domain Symptom
func (r *UpdateSymptomRequest) ApplyToDomainSymptom(s *symptom.Symptom) {
	if r.Name != nil {
		s.Name = *r.Name
	}
	if r.Description != nil {
		s.Description = *r.Description
	}
	if r.RecordedAt != nil {
		s.RecordedAt = *r.RecordedAt
	}
	if r.Severity != nil {
		s.Severity = *r.Severity
	}
	if r.Duration != nil {
		s.Duration = r.Duration
	}
	if r.BodyPart != nil {
		s.BodyPart = *r.BodyPart
	}
	if r.Triggers != nil {
		s.Triggers = r.Triggers
	}
	if r.Notes != nil {
		s.Notes = *r.Notes
	}
	if r.PhotoURL != nil {
		s.PhotoURL = *r.PhotoURL
	}

	// Handle enum fields
	if r.Category != nil {
		category := shared.SymptomCategory(*r.Category)
		s.Category = &category
	}
	if r.Type != nil {
		symptomType := shared.SymptomType(*r.Type)
		s.Type = &symptomType
	}
}

// Validate validates the CreateSymptomRequest
func (r *CreateSymptomRequest) Validate() error {
	if r.Severity < 1 || r.Severity > 10 {
		return symptom.ErrInvalidSeverity
	}
	return nil
}

// Validate validates the UpdateSymptomRequest
func (r *UpdateSymptomRequest) Validate() error {
	if r.Severity != nil && (*r.Severity < 1 || *r.Severity > 10) {
		return symptom.ErrInvalidSeverity
	}
	return nil
}
