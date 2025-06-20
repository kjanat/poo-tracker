package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/model"
)

// MealBowelMovementRelationRepository defines the interface for meal-bowel movement relation operations
type MealBowelMovementRelationRepository interface {
	Create(ctx context.Context, relation *model.MealBowelMovementRelation) error
	GetByID(ctx context.Context, id, userID string) (*model.MealBowelMovementRelation, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.MealBowelMovementRelation, error)
	GetByMealID(ctx context.Context, mealID, userID string) ([]*model.MealBowelMovementRelation, error)
	GetByBowelMovementID(ctx context.Context, bowelMovementID, userID string) ([]*model.MealBowelMovementRelation, error)
	Update(ctx context.Context, relation *model.MealBowelMovementRelation) error
	Delete(ctx context.Context, id, userID string) error
	Count(ctx context.Context, userID string) (int, error)
}

// MealSymptomRelationRepository defines the interface for meal-symptom relation operations
type MealSymptomRelationRepository interface {
	Create(ctx context.Context, relation *model.MealSymptomRelation) error
	GetByID(ctx context.Context, id, userID string) (*model.MealSymptomRelation, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.MealSymptomRelation, error)
	GetByMealID(ctx context.Context, mealID, userID string) ([]*model.MealSymptomRelation, error)
	GetBySymptomID(ctx context.Context, symptomID, userID string) ([]*model.MealSymptomRelation, error)
	Update(ctx context.Context, relation *model.MealSymptomRelation) error
	Delete(ctx context.Context, id, userID string) error
	Count(ctx context.Context, userID string) (int, error)
}

// Memory implementations
type memoryMealBowelMovementRelationRepository struct {
	mu        sync.RWMutex
	relations map[string]*model.MealBowelMovementRelation
}

func NewMemoryMealBowelMovementRelationRepository() MealBowelMovementRelationRepository {
	return &memoryMealBowelMovementRelationRepository{
		relations: make(map[string]*model.MealBowelMovementRelation),
	}
}

func (r *memoryMealBowelMovementRelationRepository) Create(ctx context.Context, relation *model.MealBowelMovementRelation) error {
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

func (r *memoryMealBowelMovementRelationRepository) GetByID(ctx context.Context, id, userID string) (*model.MealBowelMovementRelation, error) {
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

func (r *memoryMealBowelMovementRelationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userRelations []*model.MealBowelMovementRelation
	for _, relation := range r.relations {
		if relation.UserID == userID {
			userRelations = append(userRelations, relation)
		}
	}

	// Apply pagination
	if offset >= len(userRelations) {
		return []*model.MealBowelMovementRelation{}, nil
	}

	end := offset + limit
	if end > len(userRelations) {
		end = len(userRelations)
	}

	return userRelations[offset:end], nil
}

func (r *memoryMealBowelMovementRelationRepository) GetByMealID(ctx context.Context, mealID, userID string) ([]*model.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*model.MealBowelMovementRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.MealID == mealID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealBowelMovementRelationRepository) GetByBowelMovementID(ctx context.Context, bowelMovementID, userID string) ([]*model.MealBowelMovementRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*model.MealBowelMovementRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.BowelMovementID == bowelMovementID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealBowelMovementRelationRepository) Update(ctx context.Context, relation *model.MealBowelMovementRelation) error {
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
	relations map[string]*model.MealSymptomRelation
}

func NewMemoryMealSymptomRelationRepository() MealSymptomRelationRepository {
	return &memoryMealSymptomRelationRepository{
		relations: make(map[string]*model.MealSymptomRelation),
	}
}

func (r *memoryMealSymptomRelationRepository) Create(ctx context.Context, relation *model.MealSymptomRelation) error {
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

func (r *memoryMealSymptomRelationRepository) GetByID(ctx context.Context, id, userID string) (*model.MealSymptomRelation, error) {
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

func (r *memoryMealSymptomRelationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userRelations []*model.MealSymptomRelation
	for _, relation := range r.relations {
		if relation.UserID == userID {
			userRelations = append(userRelations, relation)
		}
	}

	// Apply pagination
	if offset >= len(userRelations) {
		return []*model.MealSymptomRelation{}, nil
	}

	end := offset + limit
	if end > len(userRelations) {
		end = len(userRelations)
	}

	return userRelations[offset:end], nil
}

func (r *memoryMealSymptomRelationRepository) GetByMealID(ctx context.Context, mealID, userID string) ([]*model.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*model.MealSymptomRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.MealID == mealID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealSymptomRelationRepository) GetBySymptomID(ctx context.Context, symptomID, userID string) ([]*model.MealSymptomRelation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var relations []*model.MealSymptomRelation
	for _, relation := range r.relations {
		if relation.UserID == userID && relation.SymptomID == symptomID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

func (r *memoryMealSymptomRelationRepository) Update(ctx context.Context, relation *model.MealSymptomRelation) error {
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
