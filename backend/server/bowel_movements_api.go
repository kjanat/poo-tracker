package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	rel "github.com/kjanat/poo-tracker/backend/internal/domain/relations"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

// This file contains bowel movement related API handlers
// Moved from api.go to reduce complexity

// createBowelMovementRequest defines the request structure for creating bowel movements
type createBowelMovementRequest struct {
	UserID       string              `json:"userId"`
	BristolType  int                 `json:"bristolType"`
	Volume       *shared.Volume      `json:"volume,omitempty"`
	Color        *shared.Color       `json:"color,omitempty"`
	Consistency  *shared.Consistency `json:"consistency,omitempty"`
	Floaters     *bool               `json:"floaters,omitempty"`
	Pain         *int                `json:"pain,omitempty"`
	Strain       *int                `json:"strain,omitempty"`
	Satisfaction *int                `json:"satisfaction,omitempty"`
	PhotoURL     *string             `json:"photoUrl,omitempty"`
	SmellLevel   *shared.SmellLevel  `json:"smellLevel,omitempty"`
}

// createBowelMovementDetailsRequest defines the request structure for creating bowel movement details
type createBowelMovementDetailsRequest struct {
	Notes             *string  `json:"notes,omitempty"`
	DetailedNotes     *string  `json:"detailedNotes,omitempty"`
	Environment       *string  `json:"environment,omitempty"`
	PreConditions     *string  `json:"preConditions,omitempty"`
	PostConditions    *string  `json:"postConditions,omitempty"`
	AIRecommendations *string  `json:"aiRecommendations,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  *string  `json:"weatherCondition,omitempty"`
	StressLevel       *int     `json:"stressLevel,omitempty"`
	SleepQuality      *int     `json:"sleepQuality,omitempty"`
	ExerciseIntensity *int     `json:"exerciseIntensity,omitempty"`
}

// Bowel Movement Handlers

func (a *App) listBowelMovements(c *gin.Context) {
	list, err := a.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (a *App) createBowelMovement(c *gin.Context) {
	var req createBowelMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start with a new bowel movement with defaults
	bmEntry, err := bm.NewBowelMovement(req.UserID, req.BristolType)
	if err != nil {
		if err == bm.ErrInvalidBristolType {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	bm := bmEntry

	// Apply optional fields
	applyOptionalFields(&bm, &req)

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
	var update bm.BowelMovementUpdate
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

// Bowel Movement Details Handlers

func (a *App) createBowelMovementDetails(c *gin.Context) {
	bowelMovementID := c.Param("id")

	// Check if the bowel movement exists
	_, err := a.repo.Get(c.Request.Context(), bowelMovementID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "bowel movement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if details already exist
	exists, err := a.details.Exists(c.Request.Context(), bowelMovementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "details already exist for this bowel movement"})
		return
	}

	var req createBowelMovementDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new details
	details := bm.NewBowelMovementDetails(bowelMovementID)

	// Apply optional fields
	applyBowelMovementDetailsFields(&details, &req)

	// Validate the details
	if validationErrors := validateBowelMovementDetailsFields(details); validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	created, err := a.details.Create(c.Request.Context(), details)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (a *App) getBowelMovementDetails(c *gin.Context) {
	bowelMovementID := c.Param("id")
	details, err := a.details.Get(c.Request.Context(), bowelMovementID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "details not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, details)
}

func (a *App) updateBowelMovementDetails(c *gin.Context) {
	bowelMovementID := c.Param("id")

	var update bm.BowelMovementDetailsUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the update fields
	if validationErrors := validateBowelMovementDetailsUpdate(update); validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErrors.Error()})
		return
	}

	updated, err := a.details.Update(c.Request.Context(), bowelMovementID, update)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "details not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (a *App) deleteBowelMovementDetails(c *gin.Context) {
	bowelMovementID := c.Param("id")
	if err := a.details.Delete(c.Request.Context(), bowelMovementID); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "details not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Helper functions

// applyOptionalFields applies optional fields from the request to the bowel movement
func applyOptionalFields(bm *bm.BowelMovement, req *createBowelMovementRequest) {
	if req.Volume != nil {
		bm.Volume = req.Volume
	}
	if req.Color != nil {
		bm.Color = req.Color
	}
	if req.Consistency != nil {
		bm.Consistency = req.Consistency
	}
	if req.Floaters != nil {
		bm.Floaters = *req.Floaters
	}
	if req.Pain != nil {
		bm.Pain = *req.Pain
	}
	if req.Strain != nil {
		bm.Strain = *req.Strain
	}
	if req.Satisfaction != nil {
		bm.Satisfaction = *req.Satisfaction
	}
	if req.PhotoURL != nil {
		bm.PhotoURL = *req.PhotoURL
	}
	if req.SmellLevel != nil {
		bm.SmellLevel = req.SmellLevel
	}
}

func applyBowelMovementDetailsFields(details *bm.BowelMovementDetails, req *createBowelMovementDetailsRequest) {
	if req.Notes != nil {
		details.Notes = *req.Notes
	}
	if req.DetailedNotes != nil {
		details.DetailedNotes = *req.DetailedNotes
	}
	if req.Environment != nil {
		details.Environment = *req.Environment
	}
	if req.PreConditions != nil {
		details.PreConditions = *req.PreConditions
	}
	if req.PostConditions != nil {
		details.PostConditions = *req.PostConditions
	}
	if req.AIRecommendations != nil {
		details.AIRecommendations = *req.AIRecommendations
	}
	if req.Tags != nil {
		details.Tags = req.Tags
	}
	if req.WeatherCondition != nil {
		details.WeatherCondition = *req.WeatherCondition
	}
	if req.StressLevel != nil {
		details.StressLevel = req.StressLevel
	}
	if req.SleepQuality != nil {
		details.SleepQuality = req.SleepQuality
	}
	if req.ExerciseIntensity != nil {
		details.ExerciseIntensity = req.ExerciseIntensity
	}
}

// Validation functions

func validateBowelMovementUpdateFields(update *bm.BowelMovementUpdate) validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	validationErrors = append(validationErrors, validateBowelMovementUpdateBristolType(update)...) // BristolType
	validationErrors = append(validationErrors, validateBowelMovementUpdateScales(update)...)      // Scales
	validationErrors = append(validationErrors, validateBowelMovementUpdateEnums(update)...)       // Enums
	validationErrors = append(validationErrors, validateBowelMovementUpdateOptionals(update)...)   // Optionals
	return validationErrors
}

func validateBowelMovementUpdateBristolType(update *bm.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.BristolType != nil {
		if err := validation.ValidateBristolType(*update.BristolType); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "bristolType",
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}

func validateBowelMovementUpdateScales(update *bm.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.Pain != nil {
		if err := validation.ValidateScale(*update.Pain, "pain"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "pain",
					Message: err.Error(),
				})
			}
		}
	}
	if update.Strain != nil {
		if err := validation.ValidateScale(*update.Strain, "strain"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "strain",
					Message: err.Error(),
				})
			}
		}
	}
	if update.Satisfaction != nil {
		if err := validation.ValidateScale(*update.Satisfaction, "satisfaction"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "satisfaction",
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}

func validateBowelMovementUpdateEnums(update *bm.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.Volume != nil {
		if err := validation.ValidateEnum(*update.Volume, "volume"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "volume",
					Message: err.Error(),
				})
			}
		}
	}
	if update.Color != nil {
		if err := validation.ValidateEnum(*update.Color, "color"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "color",
					Message: err.Error(),
				})
			}
		}
	}
	if update.Consistency != nil {
		if err := validation.ValidateEnum(*update.Consistency, "consistency"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "consistency",
					Message: err.Error(),
				})
			}
		}
	}
	if update.SmellLevel != nil {
		if err := validation.ValidateEnum(*update.SmellLevel, "smellLevel"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "smellLevel",
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}

func validateBowelMovementUpdateOptionals(update *bm.BowelMovementUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.PhotoURL != nil {
		if err := validation.ValidateURL(*update.PhotoURL, "photoUrl"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "photoUrl",
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}

func validateBowelMovementDetailsFields(details bm.BowelMovementDetails) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if details.Notes != "" {
		if err := validation.ValidateNotes(details.Notes, "notes"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "notes",
					Message: err.Error(),
				})
			}
		}
	}
	if details.DetailedNotes != "" {
		if err := validation.ValidateNotes(details.DetailedNotes, "detailedNotes"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "detailedNotes",
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}

func validateBowelMovementDetailsUpdate(update bm.BowelMovementDetailsUpdate) validation.ValidationErrors {
	var errs validation.ValidationErrors
	if update.Notes != nil {
		if err := validation.ValidateNotes(*update.Notes, "notes"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "notes",
					Message: err.Error(),
				})
			}
		}
	}
	if update.DetailedNotes != nil {
		if err := validation.ValidateNotes(*update.DetailedNotes, "detailedNotes"); err != nil {
			if verr, ok := err.(validation.ValidationError); ok {
				errs = append(errs, verr)
			} else {
				// Handle non-ValidationError types
				errs = append(errs, validation.ValidationError{
					Field:   "detailedNotes",
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}
