package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

// This file will be further refactored to split meal handlers
// Currently keeping minimal content to avoid conflicts

func (a *App) listMeals(c *gin.Context) {
	list, err := a.meals.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// createMealRequest defines the request structure for creating meals
type createMealRequest struct {
	UserID      string               `json:"userId"`
	Name        string               `json:"name"`
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

// applyMealFields applies fields from the request to the meal
func applyMealFields(meal *meal.Meal, req *createMealRequest) {
	meal.Name = req.Name
	if req.Description != nil {
		meal.Description = *req.Description
	}
	if req.MealTime != nil {
		meal.MealTime = *req.MealTime
	}
	if req.Category != nil {
		meal.Category = req.Category
	}
	if req.Cuisine != nil {
		meal.Cuisine = *req.Cuisine
	}
	if req.Calories != nil {
		meal.Calories = *req.Calories
	}
	if req.SpicyLevel != nil {
		meal.SpicyLevel = req.SpicyLevel
	}
	if req.FiberRich != nil {
		meal.FiberRich = *req.FiberRich
	}
	if req.Dairy != nil {
		meal.Dairy = *req.Dairy
	}
	if req.Gluten != nil {
		meal.Gluten = *req.Gluten
	}
	if req.PhotoURL != nil {
		meal.PhotoURL = *req.PhotoURL
	}
	if req.Notes != nil {
		meal.Notes = *req.Notes
	}
}

func (a *App) createMeal(c *gin.Context) {
	var req createMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start with a new meal with defaults
	mt := time.Now()
	if req.MealTime != nil {
		mt = *req.MealTime
	}
	meal := meal.NewMeal(req.UserID, req.Name, mt)

	// Apply fields from request
	applyMealFields(&meal, &req)

	// Validate the complete meal
	if validationErrors := validation.ValidateMeal(meal); validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	created, err := a.meals.Create(c.Request.Context(), meal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (a *App) getMeal(c *gin.Context) {
	id := c.Param("id")
	meal, err := a.meals.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meal)
}

func (a *App) updateMeal(c *gin.Context) {
	id := c.Param("id")
	var update meal.MealUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validationErrors := validation.ValidateMealUpdate(update); validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	updated, err := a.meals.Update(c.Request.Context(), id, update)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (a *App) deleteMeal(c *gin.Context) {
	id := c.Param("id")
	if err := a.meals.Delete(c.Request.Context(), id); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
