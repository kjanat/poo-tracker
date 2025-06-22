package meal

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// CreateMealRequest represents the request to create a new meal
type CreateMealRequest struct {
	Name        string    `json:"name" binding:"required,min=1,max=100"`
	Description *string   `json:"description,omitempty" binding:"omitempty,max=500"`
	MealTime    time.Time `json:"meal_time" binding:"required"`
	Category    *string   `json:"category,omitempty" binding:"omitempty,oneof=breakfast lunch dinner snack"`
	Cuisine     *string   `json:"cuisine,omitempty" binding:"omitempty,max=50"`
	Calories    *int      `json:"calories,omitempty" binding:"omitempty,min=0"`
	SpicyLevel  *int      `json:"spicy_level,omitempty" binding:"omitempty,min=1,max=10"`
	FiberRich   *bool     `json:"fiber_rich,omitempty"`
	Dairy       *bool     `json:"dairy,omitempty"`
	Gluten      *bool     `json:"gluten,omitempty"`
	PhotoURL    *string   `json:"photo_url,omitempty" binding:"omitempty,url"`
	Notes       *string   `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// UpdateMealRequest represents the request to update a meal
type UpdateMealRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=500"`
	MealTime    *time.Time `json:"meal_time,omitempty"`
	Category    *string    `json:"category,omitempty" binding:"omitempty,oneof=breakfast lunch dinner snack"`
	Cuisine     *string    `json:"cuisine,omitempty" binding:"omitempty,max=50"`
	Calories    *int       `json:"calories,omitempty" binding:"omitempty,min=0"`
	SpicyLevel  *int       `json:"spicy_level,omitempty" binding:"omitempty,min=1,max=10"`
	FiberRich   *bool      `json:"fiber_rich,omitempty"`
	Dairy       *bool      `json:"dairy,omitempty"`
	Gluten      *bool      `json:"gluten,omitempty"`
	PhotoURL    *string    `json:"photo_url,omitempty" binding:"omitempty,url"`
	Notes       *string    `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// MealResponse represents a meal in API responses
type MealResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MealTime    time.Time `json:"meal_time"`
	Category    *string   `json:"category,omitempty"`
	Cuisine     *string   `json:"cuisine,omitempty"`
	Calories    int       `json:"calories,omitempty"`
	SpicyLevel  *int      `json:"spicy_level,omitempty"`
	FiberRich   bool      `json:"fiber_rich"`
	Dairy       bool      `json:"dairy"`
	Gluten      bool      `json:"gluten"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MealListResponse represents a paginated list of meals
type MealListResponse struct {
	Meals      []MealResponse `json:"meals"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// MealSummaryResponse represents a simplified meal summary
type MealSummaryResponse struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	MealTime time.Time `json:"meal_time"`
	Category *string   `json:"category,omitempty"`
	Calories int       `json:"calories,omitempty"`
}

// ToMealResponse converts a domain Meal to MealResponse
func ToMealResponse(m *meal.Meal) MealResponse {
	response := MealResponse{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		Description: m.Description,
		MealTime:    m.MealTime,
		Calories:    m.Calories,
		SpicyLevel:  m.SpicyLevel,
		FiberRich:   m.FiberRich,
		Dairy:       m.Dairy,
		Gluten:      m.Gluten,
		PhotoURL:    m.PhotoURL,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	// Convert category if present
	if m.Category != nil {
		categoryStr := string(*m.Category)
		response.Category = &categoryStr
	}
	// Convert cuisine if present
	if m.Cuisine != "" {
		response.Cuisine = &m.Cuisine
	}

	return response
}

// ToMealListResponse converts a slice of domain Meals to MealListResponse
func ToMealListResponse(meals []meal.Meal, totalCount int64, page, pageSize int) MealListResponse {
	mealRes := make([]MealResponse, len(meals))
	for i, m := range meals {
		mealRes[i] = ToMealResponse(&m)
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return MealListResponse{
		Meals:      mealRes,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToMealSummaryResponse converts a domain Meal to MealSummaryResponse
func ToMealSummaryResponse(m *meal.Meal) MealSummaryResponse {
	summary := MealSummaryResponse{
		ID:       m.ID,
		Name:     m.Name,
		MealTime: m.MealTime,
		Calories: m.Calories,
	}

	// Convert category if present
	if m.Category != nil {
		categoryStr := string(*m.Category)
		summary.Category = &categoryStr
	}

	return summary
}

// ToDomainMeal converts CreateMealRequest to domain Meal
func (r *CreateMealRequest) ToDomainMeal(userID string) *meal.Meal {
	m := meal.NewMeal(userID, r.Name, r.MealTime)

	// Set optional fields
	if r.Description != nil {
		m.Description = *r.Description
	}
	if r.Category != nil {
		category := shared.MealCategory(*r.Category)
		m.Category = &category
	}
	if r.Cuisine != nil {
		m.Cuisine = *r.Cuisine
	}
	if r.Calories != nil {
		m.Calories = *r.Calories
	}
	if r.SpicyLevel != nil {
		m.SpicyLevel = r.SpicyLevel
	}
	if r.FiberRich != nil {
		m.FiberRich = *r.FiberRich
	}
	if r.Dairy != nil {
		m.Dairy = *r.Dairy
	}
	if r.Gluten != nil {
		m.Gluten = *r.Gluten
	}
	if r.PhotoURL != nil {
		m.PhotoURL = *r.PhotoURL
	}
	if r.Notes != nil {
		m.Notes = *r.Notes
	}

	return &m
}

// ApplyToDomainMeal applies UpdateMealRequest to a domain Meal
func (r *UpdateMealRequest) ApplyToDomainMeal(m *meal.Meal) {
	if r.Name != nil {
		m.Name = *r.Name
	}
	if r.Description != nil {
		m.Description = *r.Description
	}
	if r.MealTime != nil {
		m.MealTime = *r.MealTime
	}
	if r.Category != nil {
		category := shared.MealCategory(*r.Category)
		m.Category = &category
	}
	if r.Cuisine != nil {
		m.Cuisine = *r.Cuisine
	}
	if r.Calories != nil {
		m.Calories = *r.Calories
	}
	if r.SpicyLevel != nil {
		m.SpicyLevel = r.SpicyLevel
	}
	if r.FiberRich != nil {
		m.FiberRich = *r.FiberRich
	}
	if r.Dairy != nil {
		m.Dairy = *r.Dairy
	}
	if r.Gluten != nil {
		m.Gluten = *r.Gluten
	}
	if r.PhotoURL != nil {
		m.PhotoURL = *r.PhotoURL
	}
	if r.Notes != nil {
		m.Notes = *r.Notes
	}
}

// Validate validates the CreateMealRequest
func (r *CreateMealRequest) Validate() error {
	if len(r.Name) == 0 || len(r.Name) > 100 {
		return meal.ErrInvalidMealName
	}
	if r.SpicyLevel != nil && (*r.SpicyLevel < 1 || *r.SpicyLevel > 10) {
		return meal.ErrInvalidSpicyLevel
	}
	if r.Calories != nil && *r.Calories < 0 {
		return meal.ErrInvalidCalories
	}
	return nil
}

// Validate validates the UpdateMealRequest
func (r *UpdateMealRequest) Validate() error {
	if r.Name != nil && (len(*r.Name) == 0 || len(*r.Name) > 100) {
		return meal.ErrInvalidMealName
	}
	if r.SpicyLevel != nil && (*r.SpicyLevel < 1 || *r.SpicyLevel > 10) {
		return meal.ErrInvalidSpicyLevel
	}
	if r.Calories != nil && *r.Calories < 0 {
		return meal.ErrInvalidCalories
	}
	return nil
}
