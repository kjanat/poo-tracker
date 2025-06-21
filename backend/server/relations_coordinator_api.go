package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	rel "github.com/kjanat/poo-tracker/backend/internal/domain/relations"
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
	if _, err := uuid.Parse(mealID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var (
		bowelRelations   []*rel.MealBowelMovementRelation
		symptomRelations []*rel.MealSymptomRelation
		bowelErr         error
		symptomErr       error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bowelRelations, bowelErr = h.mealBowelRepo.GetByMealID(ctx, mealID, userID)
	}()

	go func() {
		defer wg.Done()
		symptomRelations, symptomErr = h.mealSymptomRepo.GetByMealID(ctx, mealID, userID)
	}()

	wg.Wait()

	if bowelErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bowel movement relations"})
		return
	}

	if symptomErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch symptom relations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bowelMovementRelations": bowelRelations,
		"symptomRelations":       symptomRelations,
	})
}
