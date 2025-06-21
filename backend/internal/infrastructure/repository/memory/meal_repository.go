package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// MealRepository implements meal.Repository using in-memory storage
type MealRepository struct {
	mu    sync.RWMutex
	meals map[string]*meal.Meal
}

// NewMealRepository creates a new in-memory meal repository
func NewMealRepository() meal.Repository {
	return &MealRepository{
		meals: make(map[string]*meal.Meal),
	}
}

// Meal CRUD operations
func (r *MealRepository) Create(ctx context.Context, m *meal.Meal) error {
	if m.ID == "" {
		return shared.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.meals[m.ID] = m
	return nil
}

func (r *MealRepository) GetByID(ctx context.Context, id string) (*meal.Meal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, exists := r.meals[id]
	if !exists {
		return nil, shared.ErrNotFound
	}

	return m, nil
}

func (r *MealRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*meal.Meal, error) {
	// Validate pagination parameters
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10 // default limit
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMeals []*meal.Meal
	for _, m := range r.meals {
		if m.UserID == userID {
			userMeals = append(userMeals, m)
		}
	}

	// Sort by MealTime descending
	sort.Slice(userMeals, func(i, j int) bool {
		return userMeals[i].MealTime.After(userMeals[j].MealTime)
	})

	// Apply pagination
	start := offset
	if start > len(userMeals) {
		return []*meal.Meal{}, nil
	}

	end := start + limit
	if end > len(userMeals) {
		end = len(userMeals)
	}

	return userMeals[start:end], nil
}

func (r *MealRepository) Update(ctx context.Context, id string, update *meal.MealUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, exists := r.meals[id]
	if !exists {
		return shared.ErrNotFound
	}

	// Apply updates
	if update.Name != nil {
		m.Name = *update.Name
	}
	if update.Description != nil {
		m.Description = *update.Description
	}
	if update.MealTime != nil {
		m.MealTime = *update.MealTime
	}
	if update.Category != nil {
		m.Category = update.Category
	}
	if update.Cuisine != nil {
		m.Cuisine = *update.Cuisine
	}
	if update.Calories != nil {
		m.Calories = *update.Calories
	}
	if update.SpicyLevel != nil {
		m.SpicyLevel = update.SpicyLevel
	}
	if update.FiberRich != nil {
		m.FiberRich = *update.FiberRich
	}
	if update.Dairy != nil {
		m.Dairy = *update.Dairy
	}
	if update.Gluten != nil {
		m.Gluten = *update.Gluten
	}
	if update.PhotoURL != nil {
		m.PhotoURL = *update.PhotoURL
	}
	if update.Notes != nil {
		m.Notes = *update.Notes
	}

	m.UpdatedAt = time.Now()
	return nil
}

func (r *MealRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.meals[id]; !exists {
		return shared.ErrNotFound
	}

	delete(r.meals, id)
	return nil
}

// Query operations
func (r *MealRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*meal.Meal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*meal.Meal
	for _, m := range r.meals {
		if m.UserID == userID &&
			!m.MealTime.Before(start) &&
			!m.MealTime.After(end) {
			result = append(result, m)
		}
	}

	// Sort by MealTime ascending
	sort.Slice(result, func(i, j int) bool {
		return result[i].MealTime.Before(result[j].MealTime)
	})

	return result, nil
}

func (r *MealRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*meal.Meal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*meal.Meal
	for _, m := range r.meals {
		if m.UserID == userID && m.Category != nil && string(*m.Category) == category {
			result = append(result, m)
		}
	}

	return result, nil
}

func (r *MealRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := int64(0)
	for _, m := range r.meals {
		if m.UserID == userID {
			count++
		}
	}

	return count, nil
}

func (r *MealRepository) GetLatestByUserID(ctx context.Context, userID string) (*meal.Meal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *meal.Meal
	for _, m := range r.meals {
		if m.UserID == userID {
			if latest == nil || m.MealTime.After(latest.MealTime) {
				latest = m
			}
		}
	}

	if latest == nil {
		return nil, shared.ErrNotFound
	}

	return latest, nil
}

// Analytics operations
func (r *MealRepository) GetNutritionSummary(ctx context.Context, userID string, start, end time.Time) (*meal.NutritionSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var meals []*meal.Meal
	for _, m := range r.meals {
		if m.UserID == userID &&
			!m.MealTime.Before(start) &&
			!m.MealTime.After(end) {
			meals = append(meals, m)
		}
	}

	if len(meals) == 0 {
		return &meal.NutritionSummary{}, nil
	}

	var totalCalories int
	var totalSpiciness float64
	fiberRichCount := int64(0)
	dairyCount := int64(0)
	glutenCount := int64(0)
	spicyCount := 0

	for _, m := range meals {
		totalCalories += m.Calories
		if m.SpicyLevel != nil {
			totalSpiciness += float64(*m.SpicyLevel)
			spicyCount++
		}
		if m.FiberRich {
			fiberRichCount++
		}
		if m.Dairy {
			dairyCount++
		}
		if m.Gluten {
			glutenCount++
		}
	}

	averageSpiciness := 0.0
	if spicyCount > 0 {
		averageSpiciness = totalSpiciness / float64(spicyCount)
	}

	count := int64(len(meals))

	return &meal.NutritionSummary{
		TotalCalories:    totalCalories,
		AverageCalories:  float64(totalCalories) / float64(count),
		FiberRichMeals:   fiberRichCount,
		DairyMeals:       dairyCount,
		GlutenMeals:      glutenCount,
		AverageSpiciness: averageSpiciness,
		MealCount:        count,
	}, nil
}

func (r *MealRepository) GetMealFrequency(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	frequency := make(map[string]int)

	for _, m := range r.meals {
		if m.UserID == userID &&
			!m.MealTime.Before(start) &&
			!m.MealTime.After(end) {
			date := m.MealTime.Format("2006-01-02")
			frequency[date]++
		}
	}

	return frequency, nil
}
