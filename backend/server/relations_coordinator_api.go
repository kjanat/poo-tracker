package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

// RelationCoordinatorHandler coordinates between meal-bowel and meal-symptom relations
type RelationCoordinatorHandler struct {
	mealBowelRepo   repository.MealBowelMovementRelationRepository
	mealSymptomRepo repository.MealSymptomRelationRepository
}

func NewRelationCoordinatorHandler(mealBowelRepo repository.MealBowelMovementRelationRepository, mealSymptomRepo repository.MealSymptomRelationRepository) *RelationCoordinatorHandler {
	return &RelationCoordinatorHandler{
		mealBowelRepo:   mealBowelRepo,
		mealSymptomRepo: mealSymptomRepo,
	}
}

// GetRelationsByMeal returns both bowel movement and symptom relations for a meal
func (h *RelationCoordinatorHandler) GetRelationsByMeal(c *gin.Context) {
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

	bowelRelations, err := h.mealBowelRepo.GetByMealID(ctx, mealID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bowel movement relations"})
		return
	}

	symptomRelations, err := h.mealSymptomRepo.GetByMealID(ctx, mealID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch symptom relations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bowelMovementRelations": bowelRelations,
		"symptomRelations":       symptomRelations,
	})
}
