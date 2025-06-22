package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	symptomDomain "github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

const maxSymptomLimit = 100

// SymptomHandler handles symptom-related HTTP requests
type SymptomHandler struct {
	repo repository.SymptomRepository
}

// NewSymptomHandler creates a new symptom handler
func NewSymptomHandler(repo repository.SymptomRepository) *SymptomHandler {
	return &SymptomHandler{repo: repo}
}

// CreateSymptom creates a new symptom
func (h *SymptomHandler) CreateSymptom(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Name        string                  `json:"name" binding:"required"`
		Description string                  `json:"description"`
		RecordedAt  *time.Time              `json:"recordedAt"`
		Category    *shared.SymptomCategory `json:"category"`
		Severity    int                     `json:"severity" binding:"min=1,max=10"`
		Duration    *int                    `json:"duration"`
		BodyPart    string                  `json:"bodyPart"`
		Type        *shared.SymptomType     `json:"type"`
		Triggers    []string                `json:"triggers"`
		Notes       string                  `json:"notes"`
		PhotoURL    string                  `json:"photoUrl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.FormatValidationError(err)})
		return
	}

	// Create symptom with defaults
	recordedAt := time.Now()
	if req.RecordedAt != nil {
		recordedAt = *req.RecordedAt
	}
	symptom := symptomDomain.NewSymptom(userID, req.Name, req.Severity, recordedAt)
	symptom.Description = req.Description

	if req.RecordedAt != nil {
		symptom.RecordedAt = *req.RecordedAt
	}

	symptom.Category = req.Category
	if req.Severity > 0 {
		symptom.Severity = req.Severity
	}
	symptom.Duration = req.Duration
	symptom.BodyPart = req.BodyPart
	symptom.Type = req.Type
	symptom.Triggers = req.Triggers
	symptom.Notes = req.Notes
	symptom.PhotoURL = req.PhotoURL

	// Validate enums
	if symptom.Category != nil && !symptom.Category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid symptom category"})
		return
	}
	if symptom.Type != nil && !symptom.Type.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid symptom type"})
		return
	}

	created, err := h.repo.Create(c.Request.Context(), symptom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create symptom"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetSymptom retrieves a symptom by ID
func (h *SymptomHandler) GetSymptom(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	symptom, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Symptom not found"})
		return
	}

	// Check if user owns this symptom
	userID := c.GetString("userID")
	if symptom.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, symptom)
}

// GetSymptoms retrieves symptoms for the authenticated user
func (h *SymptomHandler) GetSymptoms(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			if parsedLimit > 0 {
				limit = parsedLimit
			}
		}
	}

	if limit > maxSymptomLimit {
		limit = maxSymptomLimit
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			if parsedOffset >= 0 {
				offset = parsedOffset
			}
		}
	}

	symptoms, err := h.repo.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve symptoms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"symptoms": symptoms})
}

// UpdateSymptom updates an existing symptom
func (h *SymptomHandler) UpdateSymptom(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Check if user owns this symptom
	userID := c.GetString("userID")
	symptom, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Symptom not found"})
		return
	}

	if symptom.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updates symptomDomain.SymptomUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.FormatValidationError(err)})
		return
	}

	// Validate enums
	if updates.Category != nil && !updates.Category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid symptom category"})
		return
	}
	if updates.Type != nil && !updates.Type.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid symptom type"})
		return
	}
	if updates.Severity != nil && (*updates.Severity < 1 || *updates.Severity > 10) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Severity must be between 1 and 10"})
		return
	}

	updated, err := h.repo.Update(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update symptom"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSymptom deletes a symptom
func (h *SymptomHandler) DeleteSymptom(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Check if user owns this symptom
	userID := c.GetString("userID")
	symptom, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Symptom not found"})
		return
	}

	if symptom.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete symptom"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Symptom deleted successfully"})
}

// GetSymptomsByDateRange retrieves symptoms within a date range
func (h *SymptomHandler) GetSymptomsByDateRange(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDate and endDate parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use RFC3339 format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use RFC3339 format"})
		return
	}

	symptoms, err := h.repo.GetByUserIDAndDateRange(c.Request.Context(), userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve symptoms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"symptoms": symptoms})
}

// GetSymptomsByCategory retrieves symptoms by category
func (h *SymptomHandler) GetSymptomsByCategory(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryStr := c.Param("category")
	category, err := shared.ParseSymptomCategory(categoryStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid symptom category"})
		return
	}

	symptoms, err := h.repo.GetByUserIDAndCategory(c.Request.Context(), userID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve symptoms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"symptoms": symptoms})
}

// GetSymptomsByType retrieves symptoms by type
func (h *SymptomHandler) GetSymptomsByType(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	typeStr := c.Param("type")
	symptomType, err := shared.ParseSymptomType(typeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid symptom type"})
		return
	}

	symptoms, err := h.repo.GetByUserIDAndType(c.Request.Context(), userID, symptomType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve symptoms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"symptoms": symptoms})
}
