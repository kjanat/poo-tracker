package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

func TestSymptomAPI_CreateSymptom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.POST("/symptoms", handler.CreateSymptom)

	requestBody := map[string]interface{}{
		"name":     "Test Headache",
		"severity": 7,
		"notes":    "Severe headache after meal",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/symptoms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.Symptom
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != "Test Headache" {
		t.Errorf("Expected name 'Test Headache', got %s", response.Name)
	}

	if response.Severity != 7 {
		t.Errorf("Expected severity 7, got %d", response.Severity)
	}

	if response.UserID != "test-user" {
		t.Errorf("Expected userID 'test-user', got %s", response.UserID)
	}
}

func TestSymptomAPI_CreateSymptom_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.POST("/symptoms", handler.CreateSymptom)

	// Missing required name field
	requestBody := map[string]interface{}{
		"severity": 7,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/symptoms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSymptomAPI_GetSymptom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	// Create a symptom first
	symptom := model.NewSymptom("test-user", "Test Symptom")
	created, _ := repo.Create(nil, symptom)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/symptoms/:id", handler.GetSymptom)

	req, _ := http.NewRequest("GET", "/symptoms/"+created.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.Symptom
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, response.ID)
	}

	if response.Name != "Test Symptom" {
		t.Errorf("Expected name 'Test Symptom', got %s", response.Name)
	}
}

func TestSymptomAPI_GetSymptom_AccessDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	// Create a symptom for different user
	symptom := model.NewSymptom("other-user", "Other User Symptom")
	created, _ := repo.Create(nil, symptom)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/symptoms/:id", handler.GetSymptom)

	req, _ := http.NewRequest("GET", "/symptoms/"+created.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestSymptomAPI_GetSymptoms(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	// Create symptoms for the user
	symptom1 := model.NewSymptom("test-user", "Symptom 1")
	symptom2 := model.NewSymptom("test-user", "Symptom 2")
	symptom3 := model.NewSymptom("other-user", "Other User Symptom")

	repo.Create(nil, symptom1)
	repo.Create(nil, symptom2)
	repo.Create(nil, symptom3)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/symptoms", handler.GetSymptoms)

	req, _ := http.NewRequest("GET", "/symptoms", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]model.Symptom
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	symptoms := response["symptoms"]
	if len(symptoms) != 2 {
		t.Errorf("Expected 2 symptoms, got %d", len(symptoms))
	}
}

func TestSymptomAPI_UpdateSymptom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	// Create a symptom first
	symptom := model.NewSymptom("test-user", "Original")
	created, _ := repo.Create(nil, symptom)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.PUT("/symptoms/:id", handler.UpdateSymptom)

	updateBody := map[string]interface{}{
		"name":     "Updated",
		"severity": 9,
		"notes":    "Updated notes",
	}

	jsonBody, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest("PUT", "/symptoms/"+created.ID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.Symptom
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got %s", response.Name)
	}

	if response.Severity != 9 {
		t.Errorf("Expected severity 9, got %d", response.Severity)
	}

	if response.Notes != "Updated notes" {
		t.Errorf("Expected notes 'Updated notes', got %s", response.Notes)
	}
}

func TestSymptomAPI_DeleteSymptom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	// Create a symptom first
	symptom := model.NewSymptom("test-user", "To Delete")
	created, _ := repo.Create(nil, symptom)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.DELETE("/symptoms/:id", handler.DeleteSymptom)

	req, _ := http.NewRequest("DELETE", "/symptoms/"+created.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify symptom is deleted
	_, err := repo.GetByID(nil, created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted symptom")
	}
}

func TestSymptomAPI_GetSymptomsByDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// Create symptoms at different times
	symptom1 := model.NewSymptom("test-user", "Yesterday")
	symptom1.RecordedAt = yesterday

	symptom2 := model.NewSymptom("test-user", "Today")
	symptom2.RecordedAt = now

	repo.Create(nil, symptom1)
	repo.Create(nil, symptom2)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/symptoms/date-range", handler.GetSymptomsByDateRange)

	startDate := yesterday.Add(-time.Hour).Format(time.RFC3339)
	endDate := now.Add(time.Hour).Format(time.RFC3339)

	req, _ := http.NewRequest("GET", "/symptoms/date-range?startDate="+startDate+"&endDate="+endDate, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]model.Symptom
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	symptoms := response["symptoms"]
	if len(symptoms) != 2 {
		t.Errorf("Expected 2 symptoms in date range, got %d", len(symptoms))
	}
}

func TestSymptomAPI_GetSymptomsByCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewMemorySymptomRepository()
	handler := NewSymptomHandler(repo)

	category := model.SymptomCategoryDigestive

	symptom1 := model.NewSymptom("test-user", "Digestive")
	symptom1.Category = &category

	symptom2 := model.NewSymptom("test-user", "Other")
	// No category

	repo.Create(nil, symptom1)
	repo.Create(nil, symptom2)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/symptoms/category/:category", handler.GetSymptomsByCategory)

	req, _ := http.NewRequest("GET", "/symptoms/category/DIGESTIVE", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]model.Symptom
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	symptoms := response["symptoms"]
	if len(symptoms) != 1 {
		t.Errorf("Expected 1 symptom with digestive category, got %d", len(symptoms))
	}

	if symptoms[0].Name != "Digestive" {
		t.Errorf("Expected 'Digestive', got %s", symptoms[0].Name)
	}
}
