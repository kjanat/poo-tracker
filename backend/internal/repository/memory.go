package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	rel "github.com/kjanat/poo-tracker/backend/internal/domain/relations"
)

var (
	_ BowelMovementRepository        = (*memoryBowelRepo)(nil)
	_ BowelMovementDetailsRepository = (*memoryBowelDetailsRepo)(nil)
	_ MealRepository                 = (*memoryMealRepo)(nil)
)

// Separate memory repositories for interface compliance

type memoryBowelRepo struct {
	mu      sync.RWMutex
	bmStore map[string]bm.BowelMovement
}

type memoryBowelDetailsRepo struct {
	mu           sync.RWMutex
	detailsStore map[string]bm.BowelMovementDetails // keyed by BowelMovementID
	bowelRepo    *memoryBowelRepo                   // reference to update HasDetails flag
}

type memoryMealRepo struct {
	mu        sync.RWMutex
	mealStore map[string]meal.Meal
}

func NewMemoryBowelRepo() *memoryBowelRepo {
	return &memoryBowelRepo{
		bmStore: make(map[string]bm.BowelMovement),
	}
}

func NewMemoryBowelDetailsRepo(bowelRepo *memoryBowelRepo) *memoryBowelDetailsRepo {
	return &memoryBowelDetailsRepo{
		detailsStore: make(map[string]bm.BowelMovementDetails),
		bowelRepo:    bowelRepo,
	}
}

func NewMemoryMealRepo() *memoryMealRepo {
	return &memoryMealRepo{
		mealStore: make(map[string]meal.Meal),
	}
}

// BowelMovementRepository methods for memoryBowelRepo
func (m *memoryBowelRepo) List(ctx context.Context) ([]bm.BowelMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]bm.BowelMovement, 0, len(m.bmStore))
	for _, v := range m.bmStore {
		res = append(res, v)
	}
	return res, nil
}

func (m *memoryBowelRepo) Create(ctx context.Context, bm bm.BowelMovement) (bm.BowelMovement, error) {
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

func (m *memoryBowelRepo) Get(ctx context.Context, id string) (bm.BowelMovement, error) {
	m.mu.RLock()
	bm, ok := m.bmStore[id]
	m.mu.RUnlock()
	if !ok {
		return bm.BowelMovement{}, ErrNotFound
	}
	return bm, nil
}

func (m *memoryBowelRepo) Update(ctx context.Context, id string, update bm.BowelMovementUpdate) (bm.BowelMovement, error) {
	m.mu.Lock()
	existing, ok := m.bmStore[id]
	if !ok {
		m.mu.Unlock()
		return bm.BowelMovement{}, ErrNotFound
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
	if update.RecordedAt != nil {
		existing.RecordedAt = *update.RecordedAt
	}

	existing.UpdatedAt = time.Now().UTC()
	m.bmStore[id] = existing
	m.mu.Unlock()
	return existing, nil
}

// updateHasDetails updates the HasDetails flag for a bowel movement
// This is an internal method used by the details repository
func (m *memoryBowelRepo) updateHasDetails(id string, hasDetails bool) error {
	m.mu.Lock()
	existing, ok := m.bmStore[id]
	if !ok {
		m.mu.Unlock()
		return ErrNotFound
	}
	existing.HasDetails = hasDetails
	existing.UpdatedAt = time.Now().UTC()
	m.bmStore[id] = existing
	m.mu.Unlock()
	return nil
}

func (m *memoryBowelRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bmStore[id]; !ok {
		return ErrNotFound
	}
	delete(m.bmStore, id)
	return nil
}

// MealRepository methods for memoryMealRepo
func (m *memoryMealRepo) List(ctx context.Context) ([]meal.Meal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]meal.Meal, 0, len(m.mealStore))
	for _, v := range m.mealStore {
		res = append(res, v)
	}
	return res, nil
}

func (m *memoryMealRepo) Create(ctx context.Context, meal meal.Meal) (meal.Meal, error) {
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

func (m *memoryMealRepo) Get(ctx context.Context, id string) (meal.Meal, error) {
	m.mu.RLock()
	meal, ok := m.mealStore[id]
	m.mu.RUnlock()
	if !ok {
		return meal.Meal{}, ErrNotFound
	}
	return meal, nil
}

func (m *memoryMealRepo) Update(ctx context.Context, id string, update meal.MealUpdate) (meal.Meal, error) {
	m.mu.Lock()
	existing, ok := m.mealStore[id]
	if !ok {
		m.mu.Unlock()
		return meal.Meal{}, ErrNotFound
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

func (m *memoryMealRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.mealStore[id]; !ok {
		return ErrNotFound
	}
	delete(m.mealStore, id)
	return nil
}

// BowelMovementDetailsRepository methods for memoryBowelDetailsRepo
func (r *memoryBowelDetailsRepo) Create(ctx context.Context, details bm.BowelMovementDetails) (bm.BowelMovementDetails, error) {
	if details.ID == "" {
		details.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	details.CreatedAt = now
	details.UpdatedAt = now

	r.mu.Lock()
	r.detailsStore[details.BowelMovementID] = details
	r.mu.Unlock()

	// Update the corresponding BowelMovement to indicate it has details
	// Intentionally ignore error: details creation should not fail if HasDetails sync fails
	_ = r.bowelRepo.updateHasDetails(details.BowelMovementID, true)

	return details, nil
}

func (r *memoryBowelDetailsRepo) Get(ctx context.Context, bowelMovementID string) (bm.BowelMovementDetails, error) {
	r.mu.RLock()
	details, ok := r.detailsStore[bowelMovementID]
	r.mu.RUnlock()
	if !ok {
		return bm.BowelMovementDetails{}, ErrNotFound
	}
	return details, nil
}

func (r *memoryBowelDetailsRepo) Update(ctx context.Context, bowelMovementID string, update bm.BowelMovementDetailsUpdate) (bm.BowelMovementDetails, error) {
	r.mu.Lock()
	existing, ok := r.detailsStore[bowelMovementID]
	if !ok {
		r.mu.Unlock()
		return bm.BowelMovementDetails{}, ErrNotFound
	}

	// Apply updates
	if update.Notes != nil {
		existing.Notes = *update.Notes
	}
	if update.DetailedNotes != nil {
		existing.DetailedNotes = *update.DetailedNotes
	}
	if update.Environment != nil {
		existing.Environment = *update.Environment
	}
	if update.PreConditions != nil {
		existing.PreConditions = *update.PreConditions
	}
	if update.PostConditions != nil {
		existing.PostConditions = *update.PostConditions
	}
	if update.AIAnalysis != nil {
		existing.AIAnalysis = update.AIAnalysis
	}
	if update.AIConfidence != nil {
		existing.AIConfidence = update.AIConfidence
	}
	if update.AIRecommendations != nil {
		existing.AIRecommendations = *update.AIRecommendations
	}
	if update.Tags != nil {
		existing.Tags = update.Tags
	}
	if update.WeatherCondition != nil {
		existing.WeatherCondition = *update.WeatherCondition
	}
	if update.StressLevel != nil {
		existing.StressLevel = update.StressLevel
	}
	if update.SleepQuality != nil {
		existing.SleepQuality = update.SleepQuality
	}
	if update.ExerciseIntensity != nil {
		existing.ExerciseIntensity = update.ExerciseIntensity
	}

	existing.UpdatedAt = time.Now().UTC()
	r.detailsStore[bowelMovementID] = existing
	r.mu.Unlock()
	return existing, nil
}

func (r *memoryBowelDetailsRepo) Delete(ctx context.Context, bowelMovementID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.detailsStore[bowelMovementID]; !ok {
		return ErrNotFound
	}
	delete(r.detailsStore, bowelMovementID)

	// Update the corresponding BowelMovement to indicate it no longer has details
	// Intentionally ignore error: details deletion should not fail if HasDetails sync fails
	_ = r.bowelRepo.updateHasDetails(bowelMovementID, false)
	return nil
}

func (r *memoryBowelDetailsRepo) Exists(ctx context.Context, bowelMovementID string) (bool, error) {
	r.mu.RLock()
	_, exists := r.detailsStore[bowelMovementID]
	r.mu.RUnlock()
	return exists, nil
}
