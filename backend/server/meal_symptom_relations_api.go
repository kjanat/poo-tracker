package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

// MealSymptomRelationHandler handles HTTP requests for meal-symptom relations
type MealSymptomRelationHandler struct {
	mealSymptomRepo repository.MealSymptomRelationRepository
}

func NewMealSymptomRelationHandler(mealSymptomRepo repository.MealSymptomRelationRepository) *MealSymptomRelationHandler {
	return &MealSymptomRelationHandler{
		mealSymptomRepo: mealSymptomRepo,
	}
}

// Request/Response Types

type CreateMealSymptomRelationRequest struct {
	MealID          string  `json:"mealId" binding:"required"`
	SymptomID       string  `json:"symptomId" binding:"required"`
	Strength        int     `json:"strength" binding:"min=1,max=10"`
	Notes           string  `json:"notes"`
	TimeGapHours    float64 `json:"timeGapHours" binding:"min=0"`
	UserCorrelation *string `json:"userCorrelation"`
}

type UpdateMealSymptomRelationRequest struct {
	Strength        *int     `json:"strength,omitempty" binding:"omitempty,min=1,max=10"`
	Notes           *string  `json:"notes,omitempty"`
	TimeGapHours    *float64 `json:"timeGapHours,omitempty" binding:"omitempty,min=0"`
	UserCorrelation *string  `json:"userCorrelation,omitempty"`
}

// Handler Methods

func (h *MealSymptomRelationHandler) CreateMealSymptomRelation(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateMealSymptomRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Validate correlation type if provided
	var correlationType *model.CorrelationType
	if req.UserCorrelation != nil {
		if !validation.IsValidCorrelationType(*req.UserCorrelation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid correlation type"})
			return
		}
		ct := model.CorrelationType(*req.UserCorrelation)
		correlationType = &ct
	}

	relation := model.NewMealSymptomRelation(userID, req.MealID, req.SymptomID, req.TimeGapHours)
	relation.Strength = req.Strength
	relation.Notes = req.Notes
	relation.UserCorrelation = correlationType

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.mealSymptomRepo.Create(ctx, &relation); err != nil {
		if err == repository.ErrRelationAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Relation already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relation"})
		return
	}

	c.JSON(http.StatusCreated, relation)
}

func (h *MealSymptomRelationHandler) GetMealSymptomRelations(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	relations, err := h.mealSymptomRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relations"})
		return
	}

	count, _ := h.mealSymptomRepo.Count(ctx, userID)

	c.JSON(http.StatusOK, gin.H{
		"relations": relations,
		"total":     count,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *MealSymptomRelationHandler) GetMealSymptomRelation(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	relationID := c.Param("id")
	if relationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Relation ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	relation, err := h.mealSymptomRepo.GetByID(ctx, relationID, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Relation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relation"})
		return
	}

	c.JSON(http.StatusOK, relation)
}

func (h *MealSymptomRelationHandler) UpdateMealSymptomRelation(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	relationID := c.Param("id")
	if relationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Relation ID is required"})
		return
	}

	var req UpdateMealSymptomRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	relation, err := h.mealSymptomRepo.GetByID(ctx, relationID, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Relation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch relation"})
		return
	}

	// Update fields
	if req.Strength != nil {
		relation.Strength = *req.Strength
	}
	if req.Notes != nil {
		relation.Notes = *req.Notes
	}
	if req.TimeGapHours != nil {
		relation.TimeGapHours = *req.TimeGapHours
	}
	if req.UserCorrelation != nil {
		if !validation.IsValidCorrelationType(*req.UserCorrelation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid correlation type"})
			return
		}
		ct := model.CorrelationType(*req.UserCorrelation)
		relation.UserCorrelation = &ct
	}

	if err := h.mealSymptomRepo.Update(ctx, relation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update relation"})
		return
	}

	c.JSON(http.StatusOK, relation)
}

func (h *MealSymptomRelationHandler) DeleteMealSymptomRelation(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	relationID := c.Param("id")
	if relationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Relation ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.mealSymptomRepo.Delete(ctx, relationID, userID); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Relation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relation"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *MealSymptomRelationHandler) GetSymptomRelationsByMeal(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	mealID := c.Param("mealId")
	if mealID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meal ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	symptomRelations, err := h.mealSymptomRepo.GetByMealID(ctx, mealID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch symptom relations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symptomRelations": symptomRelations,
	})
}
