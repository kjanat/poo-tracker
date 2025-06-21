package meal

import (
	"context"
	"time"
)

// Repository defines the interface for meal data persistence
type Repository interface {
	// Meal CRUD operations
	Create(ctx context.Context, meal *Meal) error
	GetByID(ctx context.Context, id string) (*Meal, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Meal, error)
	Update(ctx context.Context, id string, update *MealUpdate) error
	Delete(ctx context.Context, id string) error

	// Query operations
	GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*Meal, error)
	GetByCategory(ctx context.Context, userID string, category string) ([]*Meal, error)
	GetCountByUserID(ctx context.Context, userID string) (int64, error)
	GetLatestByUserID(ctx context.Context, userID string) (*Meal, error)

	// Analytics operations
	GetNutritionSummary(ctx context.Context, userID string, start, end time.Time) (*NutritionSummary, error)
	GetMealFrequency(ctx context.Context, userID string, start, end time.Time) (map[string]int, error)
}

// NutritionSummary represents calculated nutrition data for analytics
type NutritionSummary struct {
	TotalCalories    int     `json:"totalCalories"`
	AverageCalories  float64 `json:"averageCalories"`
	FiberRichMeals   int64   `json:"fiberRichMeals"`
	DairyMeals       int64   `json:"dairyMeals"`
	GlutenMeals      int64   `json:"glutenMeals"`
	AverageSpiciness float64 `json:"averageSpiciness"`
	MealCount        int64   `json:"mealCount"`
}
