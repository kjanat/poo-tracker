package meal

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// Meal represents a meal entry with comprehensive tracking.
type Meal struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Basic meal information
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MealTime    time.Time `json:"mealTime"`

	// Categorization
	Category *shared.MealCategory `json:"category,omitempty"`
	Cuisine  string               `json:"cuisine,omitempty"`

	// Nutritional and dietary information
	Calories   int  `json:"calories,omitempty"`
	SpicyLevel *int `json:"spicyLevel,omitempty"` // 1-10 scale
	FiberRich  bool `json:"fiberRich"`
	Dairy      bool `json:"dairy"`
	Gluten     bool `json:"gluten"`

	// Optional fields
	PhotoURL string `json:"photoUrl,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

// MealUpdate represents fields that can be updated on a Meal.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
type MealUpdate struct {
	Name        *string              `json:"name,omitempty"`
	Description *string              `json:"description,omitempty"`
	MealTime    *time.Time           `json:"mealTime,omitempty"`
	Category    *shared.MealCategory `json:"category,omitempty"`
	Cuisine     *string              `json:"cuisine,omitempty"`
	Calories    *int                 `json:"calories,omitempty"`
	SpicyLevel  *int                 `json:"spicyLevel,omitempty"`
	FiberRich   *bool                `json:"fiberRich,omitempty"`
	Dairy       *bool                `json:"dairy,omitempty"`
	Gluten      *bool                `json:"gluten,omitempty"`
	PhotoURL    *string              `json:"photoUrl,omitempty"`
	Notes       *string              `json:"notes,omitempty"`
}

// NewMeal creates a new Meal with sensible defaults.
func NewMeal(userID, name string, mealTime time.Time) Meal {
	now := time.Now()
	return Meal{
		UserID:    userID,
		Name:      name,
		MealTime:  mealTime,
		CreatedAt: now,
		UpdatedAt: now,
		FiberRich: false,
		Dairy:     false,
		Gluten:    false,
	}
}
