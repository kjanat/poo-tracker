package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

// MedicationHandler handles medication-related HTTP requests
type MedicationHandler struct {
	repo repository.MedicationRepository
}

// NewMedicationHandler creates a new medication handler
func NewMedicationHandler(repo repository.MedicationRepository) *MedicationHandler {
	return &MedicationHandler{repo: repo}
}

// CreateMedication creates a new medication
func (h *MedicationHandler) CreateMedication(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Name        string                    `json:"name" binding:"required"`
		GenericName string                    `json:"genericName"`
		Brand       string                    `json:"brand"`
		Category    *model.MedicationCategory `json:"category"`
		Dosage      string                    `json:"dosage" binding:"required"`
		Form        *model.MedicationForm     `json:"form"`
		Frequency   string                    `json:"frequency" binding:"required"`
		Route       *model.MedicationRoute    `json:"route"`
		StartDate   *time.Time                `json:"startDate"`
		EndDate     *time.Time                `json:"endDate"`
		Purpose     string                    `json:"purpose"`
		Prescriber  string                    `json:"prescriber"`
		Pharmacy    string                    `json:"pharmacy"`
		Notes       string                    `json:"notes"`
		IsActive    *bool                     `json:"isActive"`
		IsPRN       *bool                     `json:"isPRN"`
		SideEffects []string                  `json:"sideEffects"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.FormatValidationError(err)})
		return
	}

	// Create medication with defaults
	medication := model.NewMedication(userID, req.Name, req.Dosage, req.Frequency)
	medication.GenericName = req.GenericName
	medication.Brand = req.Brand
	medication.Category = req.Category
	medication.Form = req.Form
	medication.Route = req.Route
	medication.StartDate = req.StartDate
	medication.EndDate = req.EndDate
	medication.Purpose = req.Purpose
	medication.Prescriber = req.Prescriber
	medication.Pharmacy = req.Pharmacy
	medication.Notes = req.Notes
	medication.SideEffects = req.SideEffects

	if req.IsActive != nil {
		medication.IsActive = *req.IsActive
	}
	if req.IsPRN != nil {
		medication.IsPRN = *req.IsPRN
	}

	// Validate enums
	if medication.Category != nil && !medication.Category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication category"})
		return
	}
	if medication.Form != nil && !medication.Form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication form"})
		return
	}
	if medication.Route != nil && !medication.Route.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication route"})
		return
	}

	created, err := h.repo.Create(c.Request.Context(), medication)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medication"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetMedication retrieves a medication by ID
func (h *MedicationHandler) GetMedication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	medication, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
		return
	}

	// Check if user owns this medication
	userID := c.GetString("userID")
	if medication.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, medication)
}

// GetMedications retrieves medications for the authenticated user
func (h *MedicationHandler) GetMedications(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	medications, err := h.repo.GetByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medications": medications})
}

// GetActiveMedications retrieves active medications for the authenticated user
func (h *MedicationHandler) GetActiveMedications(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	medications, err := h.repo.GetActiveByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active medications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medications": medications})
}

// UpdateMedication updates an existing medication
func (h *MedicationHandler) UpdateMedication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Check if user owns this medication
	userID := c.GetString("userID")
	medication, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
		return
	}

	if medication.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var updates model.MedicationUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.FormatValidationError(err)})
		return
	}

	// Validate enums
	if updates.Category != nil && !updates.Category.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication category"})
		return
	}
	if updates.Form != nil && !updates.Form.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication form"})
		return
	}
	if updates.Route != nil && !updates.Route.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication route"})
		return
	}

	updated, err := h.repo.Update(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medication"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteMedication deletes a medication
func (h *MedicationHandler) DeleteMedication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Check if user owns this medication
	userID := c.GetString("userID")
	medication, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
		return
	}

	if medication.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete medication"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medication deleted successfully"})
}

// MarkMedicationAsTaken marks a medication as taken
func (h *MedicationHandler) MarkMedicationAsTaken(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	// Check if user owns this medication
	userID := c.GetString("userID")
	medication, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
		return
	}

	if medication.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		TakenAt *time.Time `json:"takenAt"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.FormatValidationError(err)})
		return
	}

	takenAt := time.Now()
	if req.TakenAt != nil {
		takenAt = *req.TakenAt
	}

	if err := h.repo.MarkAsTaken(c.Request.Context(), id, takenAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark medication as taken"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medication marked as taken"})
}

// GetMedicationsByCategory retrieves medications by category
func (h *MedicationHandler) GetMedicationsByCategory(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryStr := c.Param("category")
	category, err := model.ParseMedicationCategory(categoryStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication category"})
		return
	}

	medications, err := h.repo.GetByUserIDAndCategory(c.Request.Context(), userID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medications": medications})
}
