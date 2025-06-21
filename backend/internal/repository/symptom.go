package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// SymptomRepository defines the interface for symptom data operations
type SymptomRepository interface {
	Create(ctx context.Context, symptom symptom.Symptom) (symptom.Symptom, error)
	GetByID(ctx context.Context, id string) (symptom.Symptom, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]symptom.Symptom, error)
	Update(ctx context.Context, id string, updates symptom.SymptomUpdate) (symptom.Symptom, error)
	Delete(ctx context.Context, id string) error
	GetByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]symptom.Symptom, error)
	GetByUserIDAndCategory(ctx context.Context, userID string, category symptom.SymptomCategory) ([]symptom.Symptom, error)
	GetByUserIDAndType(ctx context.Context, userID string, symptomType symptom.SymptomType) ([]symptom.Symptom, error)
}

// memorySymptomRepository implements SymptomRepository using in-memory storage
type memorySymptomRepository struct {
	mu       sync.RWMutex
	symptoms map[string]symptom.Symptom
}

// NewMemorySymptomRepository creates a new in-memory symptom repository
func NewMemorySymptomRepository() SymptomRepository {
	return &memorySymptomRepository{
		symptoms: make(map[string]symptom.Symptom),
	}
}

// Create creates a new symptom
func (r *memorySymptomRepository) Create(ctx context.Context, symptom symptom.Symptom) (symptom.Symptom, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if symptom.ID == "" {
		symptom.ID = uuid.New().String()
	}

	now := time.Now()
	symptom.CreatedAt = now
	symptom.UpdatedAt = now

	r.symptoms[symptom.ID] = symptom
	return symptom, nil
}

// GetByID retrieves a symptom by ID
func (r *memorySymptomRepository) GetByID(ctx context.Context, id string) (symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	symptom, exists := r.symptoms[id]
	if !exists {
		return symptom.Symptom{}, fmt.Errorf("symptom not found")
	}

	return symptom, nil
}

// GetByUserID retrieves symptoms for a specific user with pagination
func (r *memorySymptomRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userSymptoms []symptom.Symptom
	for _, symptom := range r.symptoms {
		if symptom.UserID == userID {
			userSymptoms = append(userSymptoms, symptom)
		}
	}

	// Sort by recorded time (newest first)
	sort.Slice(userSymptoms, func(i, j int) bool {
		return userSymptoms[i].RecordedAt.After(userSymptoms[j].RecordedAt)
	})

	// Apply pagination
	if offset >= len(userSymptoms) {
		return []symptom.Symptom{}, nil
	}

	end := offset + limit
	if end > len(userSymptoms) {
		end = len(userSymptoms)
	}

	return userSymptoms[offset:end], nil
}

// Update updates an existing symptom
func (r *memorySymptomRepository) Update(ctx context.Context, id string, updates symptom.SymptomUpdate) (symptom.Symptom, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	symptom, exists := r.symptoms[id]
	if !exists {
		return symptom.Symptom{}, fmt.Errorf("symptom not found")
	}

	// Apply updates
	if updates.Name != nil {
		symptom.Name = *updates.Name
	}
	if updates.Description != nil {
		symptom.Description = *updates.Description
	}
	if updates.RecordedAt != nil {
		symptom.RecordedAt = *updates.RecordedAt
	}
	if updates.Category != nil {
		symptom.Category = updates.Category
	}
	if updates.Severity != nil {
		symptom.Severity = *updates.Severity
	}
	if updates.Duration != nil {
		symptom.Duration = updates.Duration
	}
	if updates.BodyPart != nil {
		symptom.BodyPart = *updates.BodyPart
	}
	if updates.Type != nil {
		symptom.Type = updates.Type
	}
	if updates.Triggers != nil {
		symptom.Triggers = updates.Triggers
	}
	if updates.Notes != nil {
		symptom.Notes = *updates.Notes
	}
	if updates.PhotoURL != nil {
		symptom.PhotoURL = *updates.PhotoURL
	}

	symptom.UpdatedAt = time.Now()
	r.symptoms[id] = symptom

	return symptom, nil
}

// Delete removes a symptom
func (r *memorySymptomRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.symptoms[id]; !exists {
		return fmt.Errorf("symptom not found")
	}

	delete(r.symptoms, id)
	return nil
}

// GetByUserIDAndDateRange retrieves symptoms for a user within a date range
func (r *memorySymptomRepository) GetByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var symptoms []symptom.Symptom
	for _, symptom := range r.symptoms {
		if symptom.UserID == userID &&
			!symptom.RecordedAt.Before(startDate) &&
			!symptom.RecordedAt.After(endDate) {
			symptoms = append(symptoms, symptom)
		}
	}

	return symptoms, nil
}

// GetByUserIDAndCategory retrieves symptoms for a user by category
func (r *memorySymptomRepository) GetByUserIDAndCategory(ctx context.Context, userID string, category symptom.SymptomCategory) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var symptoms []symptom.Symptom
	for _, symptom := range r.symptoms {
		if symptom.UserID == userID && symptom.Category != nil && *symptom.Category == category {
			symptoms = append(symptoms, symptom)
		}
	}

	return symptoms, nil
}

// GetByUserIDAndType retrieves symptoms for a user by type
func (r *memorySymptomRepository) GetByUserIDAndType(ctx context.Context, userID string, symptomType symptom.SymptomType) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var symptoms []symptom.Symptom
	for _, symptom := range r.symptoms {
		if symptom.UserID == userID && symptom.Type != nil && *symptom.Type == symptomType {
			symptoms = append(symptoms, symptom)
		}
	}

	return symptoms, nil
}
