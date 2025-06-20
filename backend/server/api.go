package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
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
	UserID      string              `json:"userId"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	MealTime    *time.Time          `json:"mealTime,omitempty"`
	Category    *model.MealCategory `json:"category,omitempty"`
	Cuisine     string              `json:"cuisine,omitempty"`
	Calories    int                 `json:"calories,omitempty"`
	SpicyLevel  *int                `json:"spicyLevel,omitempty"`
	FiberRich   bool                `json:"fiberRich,omitempty"`
	Dairy       bool                `json:"dairy,omitempty"`
	Gluten      bool                `json:"gluten,omitempty"`
	PhotoURL    string              `json:"photoUrl,omitempty"`
	Notes       string              `json:"notes,omitempty"`
}

// applyMealFields applies fields from the request to the meal
func applyMealFields(meal *model.Meal, req *createMealRequest) {
	meal.Name = req.Name
	if req.Description != "" {
		meal.Description = req.Description
	}
	if req.MealTime != nil {
		meal.MealTime = *req.MealTime
	}
	if req.Category != nil {
		meal.Category = req.Category
	}
	if req.Cuisine != "" {
		meal.Cuisine = req.Cuisine
	}
	if req.Calories > 0 {
		meal.Calories = req.Calories
	}
	if req.SpicyLevel != nil {
		meal.SpicyLevel = req.SpicyLevel
	}
	meal.FiberRich = req.FiberRich
	meal.Dairy = req.Dairy
	meal.Gluten = req.Gluten
	if req.PhotoURL != "" {
		meal.PhotoURL = req.PhotoURL
	}
	if req.Notes != "" {
		meal.Notes = req.Notes
	}
}

func (a *App) createMeal(c *gin.Context) {
	var req createMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start with a new meal with defaults
	meal := model.NewMeal(req.UserID, req.Name)

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
	var update model.MealUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, just do basic validation - TODO: implement ValidateMealUpdate
	// if validationErrors := validation.ValidateMealUpdate(update); validationErrors.HasErrors() {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
	//     return
	// }

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
