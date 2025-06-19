package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
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
		Notes       string `json:"notes,omitempty"`
		// Accept optional enhanced fields
		Volume       *model.Volume      `json:"volume,omitempty"`
		Color        *model.Color       `json:"color,omitempty"`
		Consistency  *model.Consistency `json:"consistency,omitempty"`
		Floaters     bool               `json:"floaters,omitempty"`
		Pain         int                `json:"pain,omitempty"`
		Strain       int                `json:"strain,omitempty"`
		Satisfaction int                `json:"satisfaction,omitempty"`
		PhotoURL     string             `json:"photoUrl,omitempty"`
		SmellLevel   *model.SmellLevel  `json:"smellLevel,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start with a new bowel movement with defaults
	bm := model.NewBowelMovement(req.UserID, req.BristolType)

	// Override with provided values
	if req.Notes != "" {
		bm.Notes = req.Notes
	}
	if req.Volume != nil {
		bm.Volume = req.Volume
	}
	if req.Color != nil {
		bm.Color = req.Color
	}
	if req.Consistency != nil {
		bm.Consistency = req.Consistency
	}
	bm.Floaters = req.Floaters
	if req.Pain > 0 {
		bm.Pain = req.Pain
	}
	if req.Strain > 0 {
		bm.Strain = req.Strain
	}
	if req.Satisfaction > 0 {
		bm.Satisfaction = req.Satisfaction
	}
	if req.PhotoURL != "" {
		bm.PhotoURL = req.PhotoURL
	}
	if req.SmellLevel != nil {
		bm.SmellLevel = req.SmellLevel
	}

	// Validate the complete bowel movement
	if validationErrors := validation.ValidateBowelMovement(bm); validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
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
	var update model.BowelMovementUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors := validateBowelMovementUpdateFields(&update)
	if validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	updated, err := a.repo.Update(c.Request.Context(), id, update)
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

func validateBowelMovementUpdateFields(update *model.BowelMovementUpdate) validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	validationErrors = append(validationErrors, validateBowelMovementUpdateBristolType(update)...) // BristolType
	validationErrors = append(validationErrors, validateBowelMovementUpdateScales(update)...)      // Scales
	validationErrors = append(validationErrors, validateBowelMovementUpdateEnums(update)...)       // Enums
	validationErrors = append(validationErrors, validateBowelMovementUpdateOptionals(update)...)   // Optionals
	return validationErrors
}

func validateBowelMovementUpdateBristolType(update *model.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.BristolType != nil {
		if err := validation.ValidateBristolType(*update.BristolType); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	return errs
}

func validateBowelMovementUpdateScales(update *model.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.Pain != nil {
		if err := validation.ValidateScale(*update.Pain, "pain"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	if update.Strain != nil {
		if err := validation.ValidateScale(*update.Strain, "strain"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	if update.Satisfaction != nil {
		if err := validation.ValidateScale(*update.Satisfaction, "satisfaction"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	return errs
}

func validateBowelMovementUpdateEnums(update *model.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.Volume != nil {
		if err := validation.ValidateEnum(*update.Volume, "volume"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	if update.Color != nil {
		if err := validation.ValidateEnum(*update.Color, "color"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	if update.Consistency != nil {
		if err := validation.ValidateEnum(*update.Consistency, "consistency"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	if update.SmellLevel != nil {
		if err := validation.ValidateEnum(*update.SmellLevel, "smellLevel"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	return errs
}

func validateBowelMovementUpdateOptionals(update *model.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.PhotoURL != nil {
		if err := validation.ValidateURL(*update.PhotoURL, "photoUrl"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	if update.Notes != nil {
		if err := validation.ValidateNotes(*update.Notes, "notes"); err != nil {
			errs = append(errs, err.(validation.ValidationError))
		}
	}
	return errs
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
	list, err := a.meals.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (a *App) createMeal(c *gin.Context) {
	var req model.Meal
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the meal
	if validationErrors := validation.ValidateMeal(req); validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	created, err := a.meals.Create(c.Request.Context(), req)
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

	validationErrors := validateMealUpdateFields(&update)
	if validationErrors.HasErrors() {
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

func validateMealUpdateFields(update *model.MealUpdate) validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	if update.Name != nil {
		if err := validation.ValidateMealName(*update.Name); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	if update.Calories != nil {
		if err := validation.ValidateCalories(*update.Calories); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	if update.SpicyLevel != nil {
		if err := validation.ValidateSpicyLevel(*update.SpicyLevel); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	if update.Category != nil {
		if err := validation.ValidateEnum(*update.Category, "category"); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	if update.PhotoURL != nil {
		if err := validation.ValidateURL(*update.PhotoURL, "photoUrl"); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	if update.Description != nil {
		if err := validation.ValidateNotes(*update.Description, "description"); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	if update.Notes != nil {
		if err := validation.ValidateNotes(*update.Notes, "notes"); err != nil {
			validationErrors = append(validationErrors, err.(validation.ValidationError))
		}
	}
	return validationErrors
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
