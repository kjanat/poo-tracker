package meal

import (
	"context"
	"time"
)

// Service defines the interface for meal business logic
type Service interface {
	// Core operations
	Create(ctx context.Context, userID string, input *CreateMealInput) (*Meal, error)
	GetByID(ctx context.Context, id string) (*Meal, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Meal, error)
	Update(ctx context.Context, id string, input *UpdateMealInput) (*Meal, error)
	Delete(ctx context.Context, id string) error

	// Query operations
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*Meal, error)
	GetByCategory(ctx context.Context, userID string, category string) ([]*Meal, error)
	GetLatest(ctx context.Context, userID string) (*Meal, error)

	// Analytics operations
	GetNutritionStats(ctx context.Context, userID string, start, end time.Time) (*MealNutritionStats, error)
	GetMealInsights(ctx context.Context, userID string, start, end time.Time) (*MealInsights, error)
}

// CreateMealInput represents input for creating a meal
type CreateMealInput struct {
	Name        string    `json:"name" binding:"required,min=1,max=200"`
	Description string    `json:"description,omitempty"`
	MealTime    time.Time `json:"mealTime" binding:"required"`
	Category    *string   `json:"category,omitempty"`
	Cuisine     string    `json:"cuisine,omitempty"`
	Calories    int       `json:"calories,omitempty" binding:"omitempty,min=0,max=10000"`
	SpicyLevel  *int      `json:"spicyLevel,omitempty" binding:"omitempty,min=1,max=10"`
	FiberRich   bool      `json:"fiberRich"`
	Dairy       bool      `json:"dairy"`
	Gluten      bool      `json:"gluten"`
	PhotoURL    string    `json:"photoUrl,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}

// UpdateMealInput represents input for updating a meal
type UpdateMealInput struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description,omitempty"`
	MealTime    *time.Time `json:"mealTime,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Cuisine     *string    `json:"cuisine,omitempty"`
	Calories    *int       `json:"calories,omitempty" binding:"omitempty,min=0,max=10000"`
	SpicyLevel  *int       `json:"spicyLevel,omitempty" binding:"omitempty,min=1,max=10"`
	FiberRich   *bool      `json:"fiberRich,omitempty"`
	Dairy       *bool      `json:"dairy,omitempty"`
	Gluten      *bool      `json:"gluten,omitempty"`
	PhotoURL    *string    `json:"photoUrl,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
}

// MealNutritionStats represents nutrition analytics for a user
type MealNutritionStats struct {
	TotalCalories     int            `json:"totalCalories"`
	AverageCalories   float64        `json:"averageCalories"`
	FiberRichMeals    int64          `json:"fiberRichMeals"`
	DairyMeals        int64          `json:"dairyMeals"`
	GlutenMeals       int64          `json:"glutenMeals"`
	AverageSpiciness  float64        `json:"averageSpiciness"`
	MealCount         int64          `json:"mealCount"`
	CategoryBreakdown map[string]int `json:"categoryBreakdown"`
	CuisineBreakdown  map[string]int `json:"cuisineBreakdown"`
}

// MealInsights represents behavioral insights from meal data
type MealInsights struct {
	MostCommonCategory string             `json:"mostCommonCategory"`
	MostCommonCuisine  string             `json:"mostCommonCuisine"`
	AverageMealsPerDay float64            `json:"averageMealsPerDay"`
	MealTimePatterns   map[string]float64 `json:"mealTimePatterns"` // Hour -> frequency
	HealthScore        float64            `json:"healthScore"`      // 1-10 based on fiber, calories, etc.
}
