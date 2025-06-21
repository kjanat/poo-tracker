package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// MedicationRepository defines the interface for medication data operations
type MedicationRepository interface {
	Create(ctx context.Context, medication medication.Medication) (medication.Medication, error)
	GetByID(ctx context.Context, id string) (medication.Medication, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]medication.Medication, error)
	Update(ctx context.Context, id string, updates medication.MedicationUpdate) (medication.Medication, error)
	Delete(ctx context.Context, id string) error
	GetActiveByUserID(ctx context.Context, userID string) ([]medication.Medication, error)
	GetByUserIDAndCategory(ctx context.Context, userID string, category shared.MedicationCategory) ([]medication.Medication, error)
	MarkAsTaken(ctx context.Context, id string, takenAt time.Time) error
}

// memoryMedicationRepository implements MedicationRepository using in-memory storage
type memoryMedicationRepository struct {
	mu          sync.RWMutex
	medications map[string]medication.Medication
}

// NewMemoryMedicationRepository creates a new in-memory medication repository
func NewMemoryMedicationRepository() MedicationRepository {
	return &memoryMedicationRepository{
		medications: make(map[string]medication.Medication),
	}
}

// Create creates a new medication
func (r *memoryMedicationRepository) Create(ctx context.Context, medication medication.Medication) (medication.Medication, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if medication.ID == "" {
		medication.ID = uuid.New().String()
	}

	now := time.Now()
	medication.CreatedAt = now
	medication.UpdatedAt = now

	r.medications[medication.ID] = medication
	return medication, nil
}

// GetByID retrieves a medication by ID
func (r *memoryMedicationRepository) GetByID(ctx context.Context, id string) (medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	med, exists := r.medications[id]
	if !exists {
		return medication.Medication{}, fmt.Errorf("medication not found")
	}

	return med, nil
}

// GetByUserID retrieves medications for a specific user with pagination
func (r *memoryMedicationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMedications []medication.Medication
	for _, medication := range r.medications {
		if medication.UserID == userID {
			userMedications = append(userMedications, medication)
		}
	}

	// Sort by creation date (newest first)
	for i := 0; i < len(userMedications)-1; i++ {
		for j := i + 1; j < len(userMedications); j++ {
			if userMedications[i].CreatedAt.Before(userMedications[j].CreatedAt) {
				userMedications[i], userMedications[j] = userMedications[j], userMedications[i]
			}
		}
	}

	// Apply pagination
	if offset >= len(userMedications) {
		return []medication.Medication{}, nil
	}

	end := offset + limit
	if end > len(userMedications) {
		end = len(userMedications)
	}

	return userMedications[offset:end], nil
}

// Update updates an existing medication
func (r *memoryMedicationRepository) Update(ctx context.Context, id string, updates medication.MedicationUpdate) (medication.Medication, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	med, exists := r.medications[id]
	if !exists {
		return medication.Medication{}, fmt.Errorf("medication not found")
	}

	// Apply updates
	if updates.Name != nil {
		med.Name = *updates.Name
	}
	if updates.GenericName != nil {
		med.GenericName = *updates.GenericName
	}
	if updates.Brand != nil {
		med.Brand = *updates.Brand
	}
	if updates.Category != nil {
		med.Category = updates.Category
	}
	if updates.Dosage != nil {
		med.Dosage = *updates.Dosage
	}
	if updates.Form != nil {
		med.Form = updates.Form
	}
	if updates.Frequency != nil {
		med.Frequency = *updates.Frequency
	}
	if updates.Route != nil {
		med.Route = updates.Route
	}
	if updates.StartDate != nil {
		med.StartDate = updates.StartDate
	}
	if updates.EndDate != nil {
		med.EndDate = updates.EndDate
	}
	if updates.TakenAt != nil {
		med.TakenAt = updates.TakenAt
	}
	if updates.Purpose != nil {
		med.Purpose = *updates.Purpose
	}
	if updates.Notes != nil {
		med.Notes = *updates.Notes
	}
	if updates.IsActive != nil {
		med.IsActive = *updates.IsActive
	}
	if updates.SideEffects != nil {
		med.SideEffects = updates.SideEffects
	}
	med.UpdatedAt = time.Now()
	r.medications[id] = med

	return med, nil
}

// Delete removes a medication
func (r *memoryMedicationRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.medications[id]; !exists {
		return fmt.Errorf("medication not found")
	}

	delete(r.medications, id)
	return nil
}

// GetActiveByUserID retrieves active medications for a user
func (r *memoryMedicationRepository) GetActiveByUserID(ctx context.Context, userID string) ([]medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var medications []medication.Medication
	for _, medication := range r.medications {
		if medication.UserID == userID && medication.IsActive {
			medications = append(medications, medication)
		}
	}

	return medications, nil
}

// GetByUserIDAndCategory retrieves medications for a user by category
func (r *memoryMedicationRepository) GetByUserIDAndCategory(ctx context.Context, userID string, category shared.MedicationCategory) ([]medication.Medication, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var medications []medication.Medication
	for _, medication := range r.medications {
		if medication.UserID == userID && medication.Category != nil && *medication.Category == category {
			medications = append(medications, medication)
		}
	}

	return medications, nil
}

// MarkAsTaken marks a medication as taken at a specific time
func (r *memoryMedicationRepository) MarkAsTaken(ctx context.Context, id string, takenAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	medication, exists := r.medications[id]
	if !exists {
		return fmt.Errorf("medication not found")
	}

	medication.TakenAt = &takenAt
	medication.UpdatedAt = time.Now()
	r.medications[id] = medication

	return nil
}
