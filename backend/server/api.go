package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

func (a *App) listBowelMovements(c *gin.Context) {
	list, err := a.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (a *App) createBowelMovement(c *gin.Context) {
	var req struct {
		UserID      string `json:"userId"`
		BristolType int    `json:"bristolType"`
		Notes       string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.BristolType < 1 || req.BristolType > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bristolType must be between 1 and 7"})
		return
	}
	bm := model.BowelMovement{
		UserID:      req.UserID,
		BristolType: req.BristolType,
		Notes:       req.Notes,
	}
	created, err := a.repo.Create(c.Request.Context(), bm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (a *App) getBowelMovement(c *gin.Context) {
	id := c.Param("id")
	bm, err := a.repo.Get(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bm)
}

func (a *App) updateBowelMovement(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		BristolType *int    `json:"bristolType"`
		Notes       *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bm := model.BowelMovement{ID: id}
	if req.BristolType != nil {
		bm.BristolType = *req.BristolType
	}
	if req.Notes != nil {
		bm.Notes = *req.Notes
	}
	updated, err := a.repo.Update(c.Request.Context(), bm)
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

func (a *App) deleteBowelMovement(c *gin.Context) {
	id := c.Param("id")
	if err := a.repo.Delete(c.Request.Context(), id); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *App) getAnalytics(c *gin.Context) {
	stats, err := a.analytics.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (a *App) listMeals(c *gin.Context) {
	list, err := a.meals.ListMeals(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (a *App) createMeal(c *gin.Context) {
	var req struct {
		UserID   string `json:"userId"`
		Name     string `json:"name"`
		Calories int    `json:"calories"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := model.Meal{UserID: req.UserID, Name: req.Name, Calories: req.Calories}
	created, err := a.meals.CreateMeal(c.Request.Context(), meal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (a *App) getMeal(c *gin.Context) {
	id := c.Param("id")
	meal, err := a.meals.GetMeal(c.Request.Context(), id)
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
	var req struct {
		Name     *string `json:"name"`
		Calories *int    `json:"calories"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := model.Meal{ID: id}
	if req.Name != nil {
		meal.Name = *req.Name
	}
	if req.Calories != nil {
		meal.Calories = *req.Calories
	}
	updated, err := a.meals.UpdateMeal(c.Request.Context(), meal)
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
	if err := a.meals.DeleteMeal(c.Request.Context(), id); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
