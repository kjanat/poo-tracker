package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	rel "github.com/kjanat/poo-tracker/backend/internal/domain/relations"
)

// MealBowelMovementRelationRepository defines the interface for meal-bowel movement relation operations
type MealBowelMovementRelationRepository interface {
	Create(ctx context.Context, relation *rel.MealBowelMovementRelation) error
	GetByID(ctx context.Context, id, userID string) (*rel.MealBowelMovementRelation, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*rel.MealBowelMovementRelation, error)
	GetByMealID(ctx context.Context, mealID, userID string) ([]*rel.MealBowelMovementRelation, error)
	GetByBowelMovementID(ctx context.Context, bowelMovementID, userID string) ([]*rel.MealBowelMovementRelation, error)
	Update(ctx context.Context, relation *rel.MealBowelMovementRelation) error
	Delete(ctx context.Context, id, userID string) error
	Count(ctx context.Context, userID string) (int, error)
}

// MealSymptomRelationRepository defines the interface for meal-symptom relation operations
type MealSymptomRelationRepository interface {
	Create(ctx context.Context, relation *rel.MealSymptomRelation) error
	GetByID(ctx context.Context, id, userID string) (*rel.MealSymptomRelation, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*rel.MealSymptomRelation, error)
	GetByMealID(ctx context.Context, mealID, userID string) ([]*rel.MealSymptomRelation, error)
	GetBySymptomID(ctx context.Context, symptomID, userID string) ([]*rel.MealSymptomRelation, error)
	Update(ctx context.Context, relation *rel.MealSymptomRelation) error
	Delete(ctx context.Context, id, userID string) error
	Count(ctx context.Context, userID string) (int, error)
}

// Memory implementations
type memoryMealBowelMovementRelationRepository struct {
	mu        sync.RWMutex
	relations map[string]*rel.MealBowelMovementRelation
}

func NewMemoryMealBowelMovementRelationRepository() MealBowelMovementRelationRepository {
	return &memoryMealBowelMovementRelationRepository{
		relations: make(map[string]*rel.MealBowelMovementRelation),
	}
}

func (r *memoryMealBowelMovementRelationRepository) Create(ctx context.Context, relation *rel.MealBowelMovementRelation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if relation.ID == "" {
		relation.ID = uuid.New().String()
	}

	now := time.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	// Check for existing relation
	for _, existing := range r.relations {
		if existing.UserID == relation.UserID &&
			existing.MealID == relation.MealID &&
			existing.BowelMovementID == relation.BowelMovementID {
			return ErrRelationAlreadyExists
		}
	}

	r.relations[relation.ID] = relation
	return nil
}

func (r *memoryMealBowelMovementRelationRepository) GetByID(ctx context.Context, id, userID string) (*rel.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	relation, exists := r.relations[id]
	if !exists {
		return nil, ErrNotFound
	}

	if relation.UserID != userID {
		return nil, ErrNotFound
	}

	return relation, nil
}

func (r *memoryMealBowelMovementRelationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*rel.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userRelations []*rel.MealBowelMovementRelation
	for _, relation := range r.relations {
		if relation.UserID == userID {
			userRelations = append(userRelations, relation)
		}
	}

	// Apply pagination
	if offset >= len(userRelations) {
		return []*rel.MealBowelMovementRelation{}, nil
	}

	end := offset + limit
	if end > len(userRelations) {
		end = len(userRelations)
	}

	return userRelations[offset:end], nil
}

func (r *memoryMealBowelMovementRelationRepository) GetByMealID(ctx context.Context, mealID, userID string) ([]*rel.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*rel.MealBowelMovementRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.MealID == mealID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealBowelMovementRelationRepository) GetByBowelMovementID(ctx context.Context, bowelMovementID, userID string) ([]*rel.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*rel.MealBowelMovementRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.BowelMovementID == bowelMovementID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealBowelMovementRelationRepository) Update(ctx context.Context, relation *rel.MealBowelMovementRelation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.relations[relation.ID]
	if !exists {
		return ErrNotFound
	}

	if existing.UserID != relation.UserID {
		return ErrNotFound
	}

	relation.UpdatedAt = time.Now()
	r.relations[relation.ID] = relation
	return nil
}

func (r *memoryMealBowelMovementRelationRepository) Delete(ctx context.Context, id, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	relation, exists := r.relations[id]
	if !exists {
		return ErrNotFound
	}

	if relation.UserID != userID {
		return ErrNotFound
	}

	delete(r.relations, id)
	return nil
}

func (r *memoryMealBowelMovementRelationRepository) Count(ctx context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, relation := range r.relations {
		if relation.UserID == userID {
			count++
		}
	}

	return count, nil
}

// Memory implementation for MealSymptomRelation
type memoryMealSymptomRelationRepository struct {
	mu        sync.RWMutex
	relations map[string]*rel.MealSymptomRelation
}

func NewMemoryMealSymptomRelationRepository() MealSymptomRelationRepository {
	return &memoryMealSymptomRelationRepository{
		relations: make(map[string]*rel.MealSymptomRelation),
	}
}

func (r *memoryMealSymptomRelationRepository) Create(ctx context.Context, relation *rel.MealSymptomRelation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if relation.ID == "" {
		relation.ID = uuid.New().String()
	}

	now := time.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	// Check for existing relation
	for _, existing := range r.relations {
		if existing.UserID == relation.UserID &&
			existing.MealID == relation.MealID &&
			existing.SymptomID == relation.SymptomID {
			return ErrRelationAlreadyExists
		}
	}

	r.relations[relation.ID] = relation
	return nil
}

func (r *memoryMealSymptomRelationRepository) GetByID(ctx context.Context, id, userID string) (*rel.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	relation, exists := r.relations[id]
	if !exists {
		return nil, ErrNotFound
	}

	if relation.UserID != userID {
		return nil, ErrNotFound
	}

	return relation, nil
}

func (r *memoryMealSymptomRelationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*rel.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userRelations []*rel.MealSymptomRelation
	for _, relation := range r.relations {
		if relation.UserID == userID {
			userRelations = append(userRelations, relation)
		}
	}

	// Apply pagination
	if offset >= len(userRelations) {
		return []*rel.MealSymptomRelation{}, nil
	}

	end := offset + limit
	if end > len(userRelations) {
		end = len(userRelations)
	}

	return userRelations[offset:end], nil
}

func (r *memoryMealSymptomRelationRepository) GetByMealID(ctx context.Context, mealID, userID string) ([]*rel.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*rel.MealSymptomRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.MealID == mealID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealSymptomRelationRepository) GetBySymptomID(ctx context.Context, symptomID, userID string) ([]*rel.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*rel.MealSymptomRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.SymptomID == symptomID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealSymptomRelationRepository) Update(ctx context.Context, relation *rel.MealSymptomRelation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.relations[relation.ID]
	if !exists {
		return ErrNotFound
	}

	if existing.UserID != relation.UserID {
		return ErrNotFound
	}

	relation.UpdatedAt = time.Now()
	r.relations[relation.ID] = relation
	return nil
}

func (r *memoryMealSymptomRelationRepository) Delete(ctx context.Context, id, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	relation, exists := r.relations[id]
	if !exists {
		return ErrNotFound
	}

	if relation.UserID != userID {
		return ErrNotFound
	}

	delete(r.relations, id)
	return nil
}

func (r *memoryMealSymptomRelationRepository) Count(ctx context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, relation := range r.relations {
		if relation.UserID == userID {
			count++
		}
	}

	return count, nil
}
