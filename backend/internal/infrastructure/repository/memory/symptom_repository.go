package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// SymptomRepository implements symptom.Repository using in-memory storage
type SymptomRepository struct {
	mu       sync.RWMutex
	symptoms map[string]*symptom.Symptom
}

// NewSymptomRepository creates a new in-memory symptom repository
func NewSymptomRepository() symptom.Repository {
	return &SymptomRepository{
		symptoms: make(map[string]*symptom.Symptom),
	}
}

// Symptom CRUD operations
func (r *SymptomRepository) Create(ctx context.Context, s *symptom.Symptom) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.symptoms[s.ID] = s
	return nil
}

func (r *SymptomRepository) GetByID(ctx context.Context, id string) (*symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, exists := r.symptoms[id]
	if !exists {
		return nil, shared.ErrNotFound
	}

	return s, nil
}

func (r *SymptomRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userSymptoms []*symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID {
			userSymptoms = append(userSymptoms, s)
		}
	}

	// Sort by RecordedAt descending
	sort.Slice(userSymptoms, func(i, j int) bool {
		return userSymptoms[i].RecordedAt.After(userSymptoms[j].RecordedAt)
	})

	// Apply pagination
	start := offset
	if start > len(userSymptoms) {
		return []*symptom.Symptom{}, nil
	}

	end := start + limit
	if end > len(userSymptoms) {
		end = len(userSymptoms)
	}

	return userSymptoms[start:end], nil
}

func (r *SymptomRepository) Update(ctx context.Context, id string, update *symptom.SymptomUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, exists := r.symptoms[id]
	if !exists {
		return shared.ErrNotFound
	}

	// Apply updates
	if update.Name != nil {
		s.Name = *update.Name
	}
	if update.Description != nil {
		s.Description = *update.Description
	}
	if update.RecordedAt != nil {
		s.RecordedAt = *update.RecordedAt
	}
	if update.Category != nil {
		s.Category = update.Category
	}
	if update.Severity != nil {
		s.Severity = *update.Severity
	}
	if update.Duration != nil {
		s.Duration = update.Duration
	}
	if update.BodyPart != nil {
		s.BodyPart = *update.BodyPart
	}
	if update.Type != nil {
		s.Type = update.Type
	}
	if update.Triggers != nil {
		s.Triggers = update.Triggers
	}
	if update.Notes != nil {
		s.Notes = *update.Notes
	}

	s.UpdatedAt = time.Now()
	return nil
}

func (r *SymptomRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.symptoms[id]; !exists {
		return shared.ErrNotFound
	}

	delete(r.symptoms, id)
	return nil
}

// Query operations
func (r *SymptomRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID &&
			!s.RecordedAt.Before(start) &&
			!s.RecordedAt.After(end) {
			result = append(result, s)
		}
	}

	// Sort by RecordedAt ascending
	sort.Slice(result, func(i, j int) bool {
		return result[i].RecordedAt.Before(result[j].RecordedAt)
	})

	return result, nil
}

func (r *SymptomRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID && s.Category != nil && string(*s.Category) == category {
			result = append(result, s)
		}
	}

	return result, nil
}

func (r *SymptomRepository) GetBySeverity(ctx context.Context, userID string, minSeverity, maxSeverity int) ([]*symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID && s.Severity >= minSeverity && s.Severity <= maxSeverity {
			result = append(result, s)
		}
	}

	return result, nil
}

func (r *SymptomRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := int64(0)
	for _, s := range r.symptoms {
		if s.UserID == userID {
			count++
		}
	}

	return count, nil
}

func (r *SymptomRepository) GetLatestByUserID(ctx context.Context, userID string) (*symptom.Symptom, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID {
			if latest == nil || s.RecordedAt.After(latest.RecordedAt) {
				latest = s
			}
		}
	}

	if latest == nil {
		return nil, shared.ErrNotFound
	}

	return latest, nil
}

// Analytics operations
func (r *SymptomRepository) GetSeverityStats(ctx context.Context, userID string, start, end time.Time) (*symptom.SeverityStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var symptoms []*symptom.Symptom
	for _, s := range r.symptoms {
		if s.UserID == userID &&
			!s.RecordedAt.Before(start) &&
			!s.RecordedAt.After(end) {
			symptoms = append(symptoms, s)
		}
	}

	if len(symptoms) == 0 {
		return &symptom.SeverityStats{}, nil
	}

	var totalSeverity int
	minSeverity := 10
	maxSeverity := 1
	severityDist := make(map[int]int)

	for _, s := range symptoms {
		totalSeverity += s.Severity
		if s.Severity < minSeverity {
			minSeverity = s.Severity
		}
		if s.Severity > maxSeverity {
			maxSeverity = s.Severity
		}
		severityDist[s.Severity]++
	}

	// Find most common severity
	mostCommon := 1
	maxCount := 0
	for severity, count := range severityDist {
		if count > maxCount {
			maxCount = count
			mostCommon = severity
		}
	}

	return &symptom.SeverityStats{
		AverageSeverity:      float64(totalSeverity) / float64(len(symptoms)),
		MinSeverity:          minSeverity,
		MaxSeverity:          maxSeverity,
		MostCommonSeverity:   mostCommon,
		SeverityDistribution: severityDist,
		TotalCount:           int64(len(symptoms)),
	}, nil
}

func (r *SymptomRepository) GetCategoryFrequency(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	frequency := make(map[string]int)

	for _, s := range r.symptoms {
		if s.UserID == userID &&
			!s.RecordedAt.Before(start) &&
			!s.RecordedAt.After(end) {
			category := "Other"
			if s.Category != nil {
				category = string(*s.Category)
			}
			frequency[category]++
		}
	}

	return frequency, nil
}

func (r *SymptomRepository) GetTriggerAnalysis(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	triggerCount := make(map[string]int)

	for _, s := range r.symptoms {
		if s.UserID == userID &&
			!s.RecordedAt.Before(start) &&
			!s.RecordedAt.After(end) {
			for _, trigger := range s.Triggers {
				// Normalize trigger text (lowercase, trim spaces)
				normalizedTrigger := strings.ToLower(strings.TrimSpace(trigger))
				if normalizedTrigger != "" {
					triggerCount[normalizedTrigger]++
				}
			}
		}
	}

	return triggerCount, nil
}
