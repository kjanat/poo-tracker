package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// MedicationRepository implements medication.Repository using in-memory storage
type MedicationRepository struct {
	mu          sync.RWMutex
	medications map[string]*medication.Medication
	doseRecords map[string]*medication.DoseRecord
}

// NewMedicationRepository creates a new in-memory medication repository
func NewMedicationRepository() medication.Repository {
	return &MedicationRepository{
		medications: make(map[string]*medication.Medication),
		doseRecords: make(map[string]*medication.DoseRecord),
	}
}

// Medication CRUD operations
func (r *MedicationRepository) Create(ctx context.Context, m *medication.Medication) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.medications[m.ID] = m
	return nil
}

func (r *MedicationRepository) GetByID(ctx context.Context, id string) (*medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, exists := r.medications[id]
	if !exists {
		return nil, shared.ErrNotFound
	}

	return m, nil
}

func (r *MedicationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMedications []*medication.Medication
	for _, m := range r.medications {
		if m.UserID == userID {
			userMedications = append(userMedications, m)
		}
	}

	// Sort by Name ascending
	sort.Slice(userMedications, func(i, j int) bool {
		return userMedications[i].Name < userMedications[j].Name
	})

	// Apply pagination
	start := offset
	if start > len(userMedications) {
		return []*medication.Medication{}, nil
	}

	end := start + limit
	if end > len(userMedications) {
		end = len(userMedications)
	}

	return userMedications[start:end], nil
}

func (r *MedicationRepository) Update(ctx context.Context, id string, update *medication.MedicationUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, exists := r.medications[id]
	if !exists {
		return shared.ErrNotFound
	}

	// Apply updates
	if update.Name != nil {
		m.Name = *update.Name
	}
	if update.GenericName != nil {
		m.GenericName = *update.GenericName
	}
	if update.Brand != nil {
		m.Brand = *update.Brand
	}
	if update.Category != nil {
		m.Category = update.Category
	}
	if update.Dosage != nil {
		m.Dosage = *update.Dosage
	}
	if update.Form != nil {
		m.Form = update.Form
	}
	if update.Frequency != nil {
		m.Frequency = *update.Frequency
	}
	if update.Route != nil {
		m.Route = update.Route
	}
	if update.StartDate != nil {
		m.StartDate = update.StartDate
	}
	if update.EndDate != nil {
		m.EndDate = update.EndDate
	}
	if update.TakenAt != nil {
		m.TakenAt = update.TakenAt
	}
	if update.Purpose != nil {
		m.Purpose = *update.Purpose
	}
	if update.IsActive != nil {
		m.IsActive = *update.IsActive
	}
	if update.IsAsNeeded != nil {
		m.IsAsNeeded = *update.IsAsNeeded
	}
	if update.PhotoURL != nil {
		m.PhotoURL = *update.PhotoURL
	}
	if update.SideEffects != nil {
		m.SideEffects = update.SideEffects
	}
	if update.Notes != nil {
		m.Notes = *update.Notes
	}

	m.UpdatedAt = time.Now()
	return nil
}

func (r *MedicationRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.medications[id]; !exists {
		return shared.ErrNotFound
	}

	delete(r.medications, id)

	// Also delete associated dose records
	for recordID, record := range r.doseRecords {
		if record.MedicationID == id {
			delete(r.doseRecords, recordID)
		}
	}

	return nil
}

// Query operations
func (r *MedicationRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var activeMedications []*medication.Medication
	now := time.Now()

	for _, m := range r.medications {
		if m.UserID == userID {
			isActive := true

			// Check if medication has ended
			if m.EndDate != nil && m.EndDate.Before(now) {
				isActive = false
			}

			// Check if medication hasn't started yet
			if m.StartDate != nil && m.StartDate.After(now) {
				isActive = false
			}

			if isActive {
				activeMedications = append(activeMedications, m)
			}
		}
	}

	return activeMedications, nil
}

func (r *MedicationRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*medication.Medication
	for _, m := range r.medications {
		if m.UserID == userID && m.Category != nil && string(*m.Category) == category {
			result = append(result, m)
		}
	}

	return result, nil
}

func (r *MedicationRepository) GetAsNeededByUserID(ctx context.Context, userID string) ([]*medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*medication.Medication
	for _, m := range r.medications {
		if m.UserID == userID && m.IsAsNeeded {
			result = append(result, m)
		}
	}

	return result, nil
}

func (r *MedicationRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*medication.Medication
	for _, m := range r.medications {
		if m.UserID == userID {
			// Check if medication was active during the date range
			medStart := m.CreatedAt
			if m.StartDate != nil {
				medStart = *m.StartDate
			}

			medEnd := end
			if m.EndDate != nil {
				medEnd = *m.EndDate
			}

			// Check for overlap
			if !medStart.After(end) && !medEnd.Before(start) {
				result = append(result, m)
			}
		}
	}

	return result, nil
}

func (r *MedicationRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := int64(0)
	for _, m := range r.medications {
		if m.UserID == userID {
			count++
		}
	}

	return count, nil
}

// Tracking operations
func (r *MedicationRepository) RecordDose(ctx context.Context, medicationID string, takenAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if medication exists
	if _, exists := r.medications[medicationID]; !exists {
		return shared.ErrNotFound
	}

	record := &medication.DoseRecord{
		ID:           uuid.New().String(),
		MedicationID: medicationID,
		TakenAt:      takenAt,
		CreatedAt:    time.Now(),
	}

	r.doseRecords[record.ID] = record
	return nil
}

func (r *MedicationRepository) GetDoseHistory(ctx context.Context, medicationID string, start, end time.Time) ([]*medication.DoseRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var records []*medication.DoseRecord
	for _, record := range r.doseRecords {
		if record.MedicationID == medicationID &&
			!record.TakenAt.Before(start) &&
			!record.TakenAt.After(end) {
			records = append(records, record)
		}
	}

	// Sort by TakenAt ascending
	sort.Slice(records, func(i, j int) bool {
		return records[i].TakenAt.Before(records[j].TakenAt)
	})

	return records, nil
}

// Analytics operations
func (r *MedicationRepository) GetUsageStats(ctx context.Context, userID string, start, end time.Time) (*medication.UsageStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	medications, err := r.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	activeMedications, err := r.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	categoryBreakdown := make(map[string]int)
	formBreakdown := make(map[string]int)
	routeBreakdown := make(map[string]int)
	asNeededCount := int64(0)
	regularCount := int64(0)

	for _, m := range medications {
		// Category breakdown
		category := "Other"
		if m.Category != nil {
			category = string(*m.Category)
		}
		categoryBreakdown[category]++

		// Form breakdown
		form := "Other"
		if m.Form != nil {
			form = string(*m.Form)
		}
		formBreakdown[form]++

		// Route breakdown
		route := "Other"
		if m.Route != nil {
			route = string(*m.Route)
		}
		routeBreakdown[route]++

		// As needed vs regular
		if m.IsAsNeeded {
			asNeededCount++
		} else {
			regularCount++
		}
	}

	return &medication.UsageStats{
		TotalMedications:  int64(len(medications)),
		ActiveMedications: int64(len(activeMedications)),
		CategoryBreakdown: categoryBreakdown,
		FormBreakdown:     formBreakdown,
		RouteBreakdown:    routeBreakdown,
		AsNeededCount:     asNeededCount,
		RegularCount:      regularCount,
	}, nil
}

func (r *MedicationRepository) GetCategoryBreakdown(ctx context.Context, userID string) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	breakdown := make(map[string]int)

	for _, m := range r.medications {
		if m.UserID == userID {
			category := "Other"
			if m.Category != nil {
				category = string(*m.Category)
			}
			breakdown[category]++
		}
	}

	return breakdown, nil
}
