package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	medicationDomain "github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

func TestMedicationAPI_CreateMedication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.POST("/medications", handler.CreateMedication)

	requestBody := map[string]interface{}{
		"name":      "Ibuprofen",
		"dosage":    "200mg",
		"frequency": "twice daily",
		"purpose":   "Pain relief",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/medications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response medicationDomain.Medication
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != "Ibuprofen" {
		t.Errorf("Expected name 'Ibuprofen', got %s", response.Name)
	}

	if response.Dosage != "200mg" {
		t.Errorf("Expected dosage '200mg', got %s", response.Dosage)
	}

	if response.UserID != "test-user" {
		t.Errorf("Expected userID 'test-user', got %s", response.UserID)
	}
}

func TestMedicationAPI_CreateMedication_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.POST("/medications", handler.CreateMedication)

	// Missing required fields
	requestBody := map[string]interface{}{
		"name": "Test Med",
		// Missing dosage and frequency
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/medications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMedicationAPI_GetMedication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create a medication first
	medication := medicationDomain.NewMedication("test-user", "Test Med", "10mg", "daily")
	created, _ := repo.Create(context.TODO(), medication)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/medications/:id", handler.GetMedication)

	req, _ := http.NewRequest("GET", "/medications/"+created.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response medicationDomain.Medication
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, response.ID)
	}

	if response.Name != "Test Med" {
		t.Errorf("Expected name 'Test Med', got %s", response.Name)
	}
}

func TestMedicationAPI_GetMedication_AccessDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create a medication for different user
	medication := medicationDomain.NewMedication("other-user", "Other Med", "10mg", "daily")
	created, _ := repo.Create(context.TODO(), medication)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/medications/:id", handler.GetMedication)

	req, _ := http.NewRequest("GET", "/medications/"+created.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestMedicationAPI_GetMedications(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create medications for the user
	med1 := medicationDomain.NewMedication("test-user", "Med 1", "10mg", "daily")
	med2 := medicationDomain.NewMedication("test-user", "Med 2", "20mg", "twice daily")
	med3 := medicationDomain.NewMedication("other-user", "Other Med", "30mg", "weekly")

	// Create test data (ignore errors in test setup)
	_, _ = repo.Create(context.TODO(), med1)
	_, _ = repo.Create(context.TODO(), med2)
	_, _ = repo.Create(context.TODO(), med3)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/medications", handler.GetMedications)

	req, _ := http.NewRequest("GET", "/medications", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]medicationDomain.Medication
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	medications := response["medications"]
	if len(medications) != 2 {
		t.Errorf("Expected 2 medications, got %d", len(medications))
	}
}

func TestMedicationAPI_GetActiveMedications(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create active and inactive medications
	activeMed := medicationDomain.NewMedication("test-user", "Active Med", "10mg", "daily")
	activeMed.IsActive = true

	inactiveMed := medicationDomain.NewMedication("test-user", "Inactive Med", "20mg", "daily")
	inactiveMed.IsActive = false

	// Create test data (ignore errors in test setup)
	_, _ = repo.Create(context.TODO(), activeMed)
	_, _ = repo.Create(context.TODO(), inactiveMed)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/medications/active", handler.GetActiveMedications)

	req, _ := http.NewRequest("GET", "/medications/active", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]medicationDomain.Medication
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	medications := response["medications"]
	if len(medications) != 1 {
		t.Errorf("Expected 1 active medication, got %d", len(medications))
	}

	if medications[0].Name != "Active Med" {
		t.Errorf("Expected 'Active Med', got %s", medications[0].Name)
	}
}

func TestMedicationAPI_UpdateMedication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create a medication first
	medication := medicationDomain.NewMedication("test-user", "Original", "10mg", "daily")
	created, _ := repo.Create(context.TODO(), medication)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.PUT("/medications/:id", handler.UpdateMedication)

	updateBody := map[string]interface{}{
		"name":     "Updated",
		"dosage":   "20mg",
		"notes":    "Updated notes",
		"isActive": false,
	}

	jsonBody, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest("PUT", "/medications/"+created.ID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response medicationDomain.Medication
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got %s", response.Name)
	}

	if response.Dosage != "20mg" {
		t.Errorf("Expected dosage '20mg', got %s", response.Dosage)
	}

	if response.IsActive {
		t.Error("Expected medication to be inactive after update")
	}
}

func TestMedicationAPI_DeleteMedication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create a medication first
	medication := medicationDomain.NewMedication("test-user", "To Delete", "10mg", "daily")
	created, _ := repo.Create(context.TODO(), medication)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.DELETE("/medications/:id", handler.DeleteMedication)

	req, _ := http.NewRequest("DELETE", "/medications/"+created.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify medication is deleted
	_, err := repo.GetByID(context.TODO(), created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted medication")
	}
}

func TestMedicationAPI_MarkAsTaken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemoryMedicationRepository()
	handler := NewMedicationHandler(repo)

	// Create a medication first
	medication := medicationDomain.NewMedication("test-user", "Test Med", "10mg", "daily")
	created, _ := repo.Create(context.TODO(), medication)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.POST("/medications/:id/taken", handler.MarkMedicationAsTaken)

	requestBody := map[string]interface{}{
		// Empty body should use current time
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/medications/"+created.ID+"/taken", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify medication is marked as taken
	updated, err := repo.GetByID(context.TODO(), created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.TakenAt == nil {
		t.Error("Expected TakenAt to be set")
	}
}
