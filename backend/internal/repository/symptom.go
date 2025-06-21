package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
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
	GetByUserIDAndCategory(ctx context.Context, userID string, category shared.SymptomCategory) ([]symptom.Symptom, error)
	GetByUserIDAndType(ctx context.Context, userID string, symptomType shared.SymptomType) ([]symptom.Symptom, error)
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
func (r *memorySymptomRepository) Create(ctx context.Context, s symptom.Symptom) (symptom.Symptom, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	r.symptoms[s.ID] = s
	return s, nil
}

// GetByID retrieves a symptom by ID
func (r *memorySymptomRepository) GetByID(ctx context.Context, id string) (symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, exists := r.symptoms[id]
	if !exists {
		return symptom.Symptom{}, fmt.Errorf("symptom not found")
	}

	return s, nil
}

// GetByUserID retrieves symptoms for a specific user with pagination
func (r *memorySymptomRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userSymptoms []symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID {
			userSymptoms = append(userSymptoms, s)
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

	s, exists := r.symptoms[id]
	if !exists {
		return symptom.Symptom{}, fmt.Errorf("symptom not found")
	}

	// Apply updates
	if updates.Name != nil {
		s.Name = *updates.Name
	}
	if updates.Description != nil {
		s.Description = *updates.Description
	}
	if updates.RecordedAt != nil {
		s.RecordedAt = *updates.RecordedAt
	}
	if updates.Category != nil {
		s.Category = updates.Category
	}
	if updates.Severity != nil {
		s.Severity = *updates.Severity
	}
	if updates.Duration != nil {
		s.Duration = updates.Duration
	}
	if updates.BodyPart != nil {
		s.BodyPart = *updates.BodyPart
	}
	if updates.Type != nil {
		s.Type = updates.Type
	}
	if updates.Triggers != nil {
		s.Triggers = updates.Triggers
	}
	if updates.Notes != nil {
		s.Notes = *updates.Notes
	}
	if updates.PhotoURL != nil {
		s.PhotoURL = *updates.PhotoURL
	}

	s.UpdatedAt = time.Now()
	r.symptoms[id] = s

	return s, nil
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
	for _, s := range r.symptoms {
		if s.UserID == userID &&
			!s.RecordedAt.Before(startDate) &&
			!s.RecordedAt.After(endDate) {
			symptoms = append(symptoms, s)
		}
	}

	return symptoms, nil
}

// GetByUserIDAndCategory retrieves symptoms for a user by category
func (r *memorySymptomRepository) GetByUserIDAndCategory(ctx context.Context, userID string, category shared.SymptomCategory) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var symptoms []symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID && s.Category != nil && *s.Category == category {
			symptoms = append(symptoms, s)
		}
	}

	return symptoms, nil
}

// GetByUserIDAndType retrieves symptoms for a user by type
func (r *memorySymptomRepository) GetByUserIDAndType(ctx context.Context, userID string, symptomType shared.SymptomType) ([]symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var symptoms []symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID && s.Type != nil && *s.Type == symptomType {
			symptoms = append(symptoms, s)
		}
	}

	return symptoms, nil
}
