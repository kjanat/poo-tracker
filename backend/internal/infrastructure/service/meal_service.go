package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// MealService implements the meal business logic
type MealService struct {
	repo meal.Repository
}

// NewMealService creates a new meal service
func NewMealService(repo meal.Repository) meal.Service {
	return &MealService{
		repo: repo,
	}
}

// Create creates a new meal with business validation
func (s *MealService) Create(ctx context.Context, userID string, input *meal.CreateMealInput) (*meal.Meal, error) {
	// Validate user ID
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	// Validate input
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Set defaults
	if input.MealTime.IsZero() {
		input.MealTime = time.Now()
	}

	// Convert string pointers to shared types if needed
	var category *shared.MealCategory
	if input.Category != nil {
		cat := shared.MealCategory(*input.Category)
		if cat.IsValid() {
			category = &cat
		}
	}

	// Create meal
	mealEntity := &meal.Meal{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		MealTime:    input.MealTime,
		Category:    category,
		Cuisine:     input.Cuisine,
		Calories:    input.Calories,
		SpicyLevel:  input.SpicyLevel,
		FiberRich:   input.FiberRich,
		Dairy:       input.Dairy,
		Gluten:      input.Gluten,
		PhotoURL:    input.PhotoURL,
		Notes:       input.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository
	if err := s.repo.Create(ctx, mealEntity); err != nil {
		return nil, fmt.Errorf("failed to create meal: %w", err)
	}

	return mealEntity, nil
}

// GetByID retrieves a meal by ID
func (s *MealService) GetByID(ctx context.Context, id string) (*meal.Meal, error) {
	if id == "" {
		return nil, meal.ErrInvalidID
	}

	mealEntity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, meal.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	return mealEntity, nil
}

// GetByUserID retrieves meals for a specific user with pagination
func (s *MealService) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*meal.Meal, error) {
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	// Apply business rules for pagination
	if limit <= 0 || limit > 100 {
		limit = 20 // default
	}
	if offset < 0 {
		offset = 0
	}

	meals, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meals: %w", err)
	}

	return meals, nil
}

// Update updates an existing meal
func (s *MealService) Update(ctx context.Context, id string, input *meal.UpdateMealInput) (*meal.Meal, error) {
	if id == "" {
		return nil, meal.ErrInvalidID
	}

	// Get existing meal to verify it exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, meal.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get meal for update: %w", err)
	}

	// Validate update input
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// Convert input to update struct
	update := s.convertToUpdateStruct(input)

	// Save changes
	if err := s.repo.Update(ctx, id, update); err != nil {
		return nil, fmt.Errorf("failed to update meal: %w", err)
	}

	// Return updated meal
	return s.repo.GetByID(ctx, id)
}

// Delete removes a meal
func (s *MealService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return meal.ErrInvalidID
	}

	// Check if meal exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return meal.ErrNotFound
		}
		return fmt.Errorf("failed to verify meal exists: %w", err)
	}

	// Delete meal
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete meal: %w", err)
	}

	return nil
}

// GetByDateRange retrieves meals within a date range
func (s *MealService) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*meal.Meal, error) {
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	// Validate date range
	if start.After(end) {
		return nil, meal.ErrInvalidDateRange
	}

	// Limit date range to reasonable bounds
	maxRange := 365 * 24 * time.Hour // 1 year
	if end.Sub(start) > maxRange {
		return nil, meal.ErrDateRangeTooLarge
	}

	meals, err := s.repo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals by date range: %w", err)
	}

	return meals, nil
}

// GetByCategory retrieves meals by category
func (s *MealService) GetByCategory(ctx context.Context, userID string, category string) ([]*meal.Meal, error) {
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	if category == "" {
		return nil, meal.ErrInvalidCategory
	}

	meals, err := s.repo.GetByCategory(ctx, userID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals by category: %w", err)
	}

	return meals, nil
}

// GetLatest retrieves the most recent meal for a user
func (s *MealService) GetLatest(ctx context.Context, userID string) (*meal.Meal, error) {
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	latest, err := s.repo.GetLatestByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, meal.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest meal: %w", err)
	}

	return latest, nil
}

// GetNutritionStats generates nutrition analytics for a user's meals
func (s *MealService) GetNutritionStats(ctx context.Context, userID string, start, end time.Time) (*meal.MealNutritionStats, error) {
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	// Get meals in date range
	meals, err := s.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	if len(meals) == 0 {
		return &meal.MealNutritionStats{
			MealCount: 0,
		}, nil
	}

	// Calculate nutrition statistics
	stats := s.calculateNutritionStats(meals, start, end)
	return stats, nil
}

// GetMealInsights generates insights for a user's meals
func (s *MealService) GetMealInsights(ctx context.Context, userID string, start, end time.Time) (*meal.MealInsights, error) {
	if userID == "" {
		return nil, meal.ErrInvalidUserID
	}

	// Get meals in date range
	meals, err := s.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	if len(meals) == 0 {
		return &meal.MealInsights{
			AverageMealsPerDay: 0,
		}, nil
	}

	// Calculate insights
	insights := s.calculateInsights(meals, start, end)
	return insights, nil
}

// validateCreateInput validates create input
func (s *MealService) validateCreateInput(input *meal.CreateMealInput) error {
	if input == nil {
		return meal.ErrInvalidInput
	}

	if input.Name == "" {
		return meal.ErrInvalidName
	}

	if input.Calories < 0 || input.Calories > 10000 {
		return meal.ErrInvalidCalories
	}

	if input.SpicyLevel != nil && (*input.SpicyLevel < 1 || *input.SpicyLevel > 10) {
		return meal.ErrInvalidSpicyLevel
	}

	// Validate category if provided
	if input.Category != nil {
		cat := shared.MealCategory(*input.Category)
		if !cat.IsValid() {
			return meal.ErrInvalidCategory
		}
	}

	return nil
}

// validateUpdateInput validates update input
func (s *MealService) validateUpdateInput(input *meal.UpdateMealInput) error {
	if input == nil {
		return meal.ErrInvalidInput
	}

	if input.Name != nil && *input.Name == "" {
		return meal.ErrInvalidName
	}

	if input.Calories != nil && (*input.Calories < 0 || *input.Calories > 10000) {
		return meal.ErrInvalidCalories
	}

	if input.SpicyLevel != nil && (*input.SpicyLevel < 1 || *input.SpicyLevel > 10) {
		return meal.ErrInvalidSpicyLevel
	}

	return nil
}

// convertToUpdateStruct converts service input to repository update struct
func (s *MealService) convertToUpdateStruct(input *meal.UpdateMealInput) *meal.MealUpdate {
	update := &meal.MealUpdate{
		Name:        input.Name,
		Description: input.Description,
		MealTime:    input.MealTime,
		Cuisine:     input.Cuisine,
		Calories:    input.Calories,
		SpicyLevel:  input.SpicyLevel,
		FiberRich:   input.FiberRich,
		Dairy:       input.Dairy,
		Gluten:      input.Gluten,
		PhotoURL:    input.PhotoURL,
		Notes:       input.Notes,
	}

	// Convert string pointer to shared type pointer
	if input.Category != nil {
		cat := shared.MealCategory(*input.Category)
		update.Category = &cat
	}

	return update
}

// calculateNutritionStats calculates nutrition statistics from meals
func (s *MealService) calculateNutritionStats(meals []*meal.Meal, start, end time.Time) *meal.MealNutritionStats {
	mealCount := int64(len(meals))
	var totalCalories int64
	fiberRichCount := int64(0)
	dairyCount := int64(0)
	glutenCount := int64(0)
	var totalSpiciness int64
	spicyMealsCount := int64(0)

	categoryBreakdown := make(map[string]int)
	cuisineBreakdown := make(map[string]int)

	for _, m := range meals {
		totalCalories += int64(m.Calories)
		if m.FiberRich {
			fiberRichCount++
		}
		if m.Dairy {
			dairyCount++
		}
		if m.Gluten {
			glutenCount++
		}
		if m.SpicyLevel != nil {
			totalSpiciness += int64(*m.SpicyLevel)
			spicyMealsCount++
		}

		// Count categories
		if m.Category != nil {
			categoryBreakdown[string(*m.Category)]++
		}

		// Count cuisines
		if m.Cuisine != "" {
			cuisineBreakdown[m.Cuisine]++
		}
	}

	// Calculate averages
	var avgCalories, avgSpiciness float64
	if mealCount > 0 {
		avgCalories = float64(totalCalories) / float64(mealCount)
	}
	if spicyMealsCount > 0 {
		avgSpiciness = float64(totalSpiciness) / float64(spicyMealsCount)
	}

	return &meal.MealNutritionStats{
		TotalCalories:     int(totalCalories),
		AverageCalories:   avgCalories,
		FiberRichMeals:    fiberRichCount,
		DairyMeals:        dairyCount,
		GlutenMeals:       glutenCount,
		AverageSpiciness:  avgSpiciness,
		MealCount:         mealCount,
		CategoryBreakdown: categoryBreakdown,
		CuisineBreakdown:  cuisineBreakdown,
	}
}

// calculateInsights calculates insights from meals
func (s *MealService) calculateInsights(meals []*meal.Meal, start, end time.Time) *meal.MealInsights {
	totalMeals := int64(len(meals))
	categoryDistribution := make(map[string]int)
	cuisineDistribution := make(map[string]int)

	var mostCommonCategory, mostCommonCuisine string
	maxCategoryCount, maxCuisineCount := 0, 0

	for _, m := range meals {
		// Count categories
		if m.Category != nil {
			categoryStr := string(*m.Category)
			categoryDistribution[categoryStr]++
			if categoryDistribution[categoryStr] > maxCategoryCount {
				maxCategoryCount = categoryDistribution[categoryStr]
				mostCommonCategory = categoryStr
			}
		}

		// Count cuisines
		if m.Cuisine != "" {
			cuisineDistribution[m.Cuisine]++
			if cuisineDistribution[m.Cuisine] > maxCuisineCount {
				maxCuisineCount = cuisineDistribution[m.Cuisine]
				mostCommonCuisine = m.Cuisine
			}
		}
	}

	// Calculate frequency per day
	days := end.Sub(start).Hours() / 24
	if days <= 0 {
		days = 1
	}
	averageMealsPerDay := float64(totalMeals) / days

	// Calculate meal time patterns (hour -> frequency)
	mealTimePatterns := make(map[string]float64)
	for _, m := range meals {
		hour := m.MealTime.Hour()
		hourStr := fmt.Sprintf("%d", hour)
		mealTimePatterns[hourStr]++
	}
	// Normalize to percentages
	for hour := range mealTimePatterns {
		mealTimePatterns[hour] = mealTimePatterns[hour] / float64(totalMeals)
	}

	// Calculate simple health score based on fiber content
	healthScore := 5.0 // default middle score
	if totalMeals > 0 {
		fiberRichCount := 0
		for _, m := range meals {
			if m.FiberRich {
				fiberRichCount++
			}
		}
		// Health score 1-10 based on fiber content
		fiberRatio := float64(fiberRichCount) / float64(totalMeals)
		healthScore = 1.0 + (fiberRatio * 9.0) // 1-10 scale
	}

	return &meal.MealInsights{
		MostCommonCategory: mostCommonCategory,
		MostCommonCuisine:  mostCommonCuisine,
		AverageMealsPerDay: averageMealsPerDay,
		MealTimePatterns:   mealTimePatterns,
		HealthScore:        healthScore,
	}
}
