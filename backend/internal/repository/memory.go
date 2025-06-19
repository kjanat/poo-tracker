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

	// Set RecordedAt to now if not provided
	if bm.RecordedAt.IsZero() {
		bm.RecordedAt = now
	}

	// Set default values for experience scales if not provided
	if bm.Pain == 0 {
		bm.Pain = 1 // Default: minimal pain
	}
	if bm.Strain == 0 {
		bm.Strain = 1 // Default: minimal strain
	}
	if bm.Satisfaction == 0 {
		bm.Satisfaction = 5 // Default: neutral satisfaction
	}

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
	if update.Volume != nil {
		existing.Volume = update.Volume
	}
	if update.Color != nil {
		existing.Color = update.Color
	}
	if update.Consistency != nil {
		existing.Consistency = update.Consistency
	}
	if update.Floaters != nil {
		existing.Floaters = *update.Floaters
	}
	if update.Pain != nil {
		existing.Pain = *update.Pain
	}
	if update.Strain != nil {
		existing.Strain = *update.Strain
	}
	if update.Satisfaction != nil {
		existing.Satisfaction = *update.Satisfaction
	}
	if update.PhotoURL != nil {
		existing.PhotoURL = *update.PhotoURL
	}
	if update.SmellLevel != nil {
		existing.SmellLevel = update.SmellLevel
	}
	if update.Notes != nil {
		existing.Notes = *update.Notes
	}
	if update.RecordedAt != nil {
		existing.RecordedAt = *update.RecordedAt
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

	// Set MealTime to now if not provided
	if meal.MealTime.IsZero() {
		meal.MealTime = now
	}

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

func (m *memoryRepo) UpdateMeal(ctx context.Context, id string, update model.MealUpdate) (model.Meal, error) {
	m.mu.Lock()
	existing, ok := m.mealStore[id]
	if !ok {
		m.mu.Unlock()
		return model.Meal{}, ErrNotFound
	}

	// Only update fields that are explicitly provided (non-nil pointers)
	if update.Name != nil {
		existing.Name = *update.Name
	}
	if update.Description != nil {
		existing.Description = *update.Description
	}
	if update.MealTime != nil {
		existing.MealTime = *update.MealTime
	}
	if update.Category != nil {
		existing.Category = update.Category
	}
	if update.Cuisine != nil {
		existing.Cuisine = *update.Cuisine
	}
	if update.Calories != nil {
		existing.Calories = *update.Calories
	}
	if update.SpicyLevel != nil {
		existing.SpicyLevel = update.SpicyLevel
	}
	if update.FiberRich != nil {
		existing.FiberRich = *update.FiberRich
	}
	if update.Dairy != nil {
		existing.Dairy = *update.Dairy
	}
	if update.Gluten != nil {
		existing.Gluten = *update.Gluten
	}
	if update.PhotoURL != nil {
		existing.PhotoURL = *update.PhotoURL
	}
	if update.Notes != nil {
		existing.Notes = *update.Notes
	}

	existing.UpdatedAt = time.Now().UTC()
	m.mealStore[id] = existing
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
