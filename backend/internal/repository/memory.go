package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/model"
)

var (
	_ BowelMovementRepository = (*memoryRepo)(nil)
	_ MealRepository          = (*memoryRepo)(nil)
)

type memoryRepo struct {
	mu        sync.RWMutex
	bmStore   map[string]model.BowelMovement
	mealStore map[string]model.Meal
}

func NewMemory() *memoryRepo {
	return &memoryRepo{
		bmStore:   make(map[string]model.BowelMovement),
		mealStore: make(map[string]model.Meal),
	}
}

func (m *memoryRepo) List(ctx context.Context) ([]model.BowelMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]model.BowelMovement, len(m.bmStore))
	i := 0
	for _, v := range m.bmStore {
		res[i] = v
		i++
	}
	return res, nil
}

func (m *memoryRepo) Create(ctx context.Context, bm model.BowelMovement) (model.BowelMovement, error) {
	if bm.ID == "" {
		bm.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	bm.CreatedAt = now
	bm.UpdatedAt = now
	m.mu.Lock()
	m.bmStore[bm.ID] = bm
	m.mu.Unlock()
	return bm, nil
}

func (m *memoryRepo) Get(ctx context.Context, id string) (model.BowelMovement, error) {
	m.mu.RLock()
	bm, ok := m.bmStore[id]
	m.mu.RUnlock()
	if !ok {
		return model.BowelMovement{}, ErrNotFound
	}
	return bm, nil
}

func (m *memoryRepo) Update(ctx context.Context, id string, update model.BowelMovementUpdate) (model.BowelMovement, error) {
	m.mu.Lock()
	existing, ok := m.bmStore[id]
	if !ok {
		m.mu.Unlock()
		return model.BowelMovement{}, ErrNotFound
	}

	// Only update fields that are explicitly provided (non-nil pointers)
	if update.BristolType != nil {
		existing.BristolType = *update.BristolType
	}
	if update.Notes != nil {
		existing.Notes = *update.Notes
	}

	existing.UpdatedAt = time.Now().UTC()
	m.bmStore[id] = existing
	m.mu.Unlock()
	return existing, nil
}

func (m *memoryRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bmStore[id]; !ok {
		return ErrNotFound
	}
	delete(m.bmStore, id)
	return nil
}

func (m *memoryRepo) ListMeals(ctx context.Context) ([]model.Meal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]model.Meal, 0, len(m.mealStore))
	for _, v := range m.mealStore {
		res = append(res, v)
	}
	return res, nil
}

func (m *memoryRepo) CreateMeal(ctx context.Context, meal model.Meal) (model.Meal, error) {
	if meal.ID == "" {
		meal.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	meal.CreatedAt = now
	meal.UpdatedAt = now
	m.mu.Lock()
	m.mealStore[meal.ID] = meal
	m.mu.Unlock()
	return meal, nil
}

func (m *memoryRepo) GetMeal(ctx context.Context, id string) (model.Meal, error) {
	m.mu.RLock()
	meal, ok := m.mealStore[id]
	m.mu.RUnlock()
	if !ok {
		return model.Meal{}, ErrNotFound
	}
	return meal, nil
}

func (m *memoryRepo) UpdateMeal(ctx context.Context, meal model.Meal) (model.Meal, error) {
	m.mu.Lock()
	existing, ok := m.mealStore[meal.ID]
	if !ok {
		m.mu.Unlock()
		return model.Meal{}, ErrNotFound
	}
	if meal.Name != "" {
		existing.Name = meal.Name
	}
	if meal.Calories != 0 {
		existing.Calories = meal.Calories
	}
	existing.UpdatedAt = time.Now().UTC()
	m.mealStore[meal.ID] = existing
	m.mu.Unlock()
	return existing, nil
}

func (m *memoryRepo) DeleteMeal(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.mealStore[id]; !ok {
		return ErrNotFound
	}
	delete(m.mealStore, id)
	return nil
}
