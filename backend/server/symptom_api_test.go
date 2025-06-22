package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	symptomDomain "github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
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

	var response symptomDomain.Symptom
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
	symptom := symptomDomain.NewSymptom("test-user", "Test Symptom", 5, time.Now())
	created, _ := repo.Create(context.TODO(), symptom)

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

	var response symptomDomain.Symptom
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
	symptom := symptomDomain.NewSymptom("other-user", "Other User Symptom", 5, time.Now())
	created, _ := repo.Create(context.TODO(), symptom)

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
	symptom1 := symptomDomain.NewSymptom("test-user", "Symptom 1", 5, time.Now())
	symptom2 := symptomDomain.NewSymptom("test-user", "Symptom 2", 5, time.Now())
	symptom3 := symptomDomain.NewSymptom("other-user", "Other User Symptom", 5, time.Now())

	// Create test data (ignore errors in test setup)
	_, _ = repo.Create(context.TODO(), symptom1)
	_, _ = repo.Create(context.TODO(), symptom2)
	_, _ = repo.Create(context.TODO(), symptom3)

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

	var response map[string][]symptomDomain.Symptom
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
	symptom := symptomDomain.NewSymptom("test-user", "Original", 5, time.Now())
	created, _ := repo.Create(context.TODO(), symptom)

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

	var response symptomDomain.Symptom
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
	symptom := symptomDomain.NewSymptom("test-user", "To Delete", 5, time.Now())
	created, _ := repo.Create(context.TODO(), symptom)

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
	_, err := repo.GetByID(context.TODO(), created.ID)
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
	symptom1 := symptomDomain.NewSymptom("test-user", "Yesterday", 5, yesterday)
	symptom2 := symptomDomain.NewSymptom("test-user", "Today", 5, now)

	_, _ = repo.Create(context.TODO(), symptom1)
	_, _ = repo.Create(context.TODO(), symptom2)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "test-user")
		c.Next()
	})
	router.GET("/symptoms/date-range", handler.GetSymptomsByDateRange)

	startDate := yesterday.Add(-time.Hour).Format(time.RFC3339)
	endDate := now.Add(time.Hour).Format(time.RFC3339)

	req, _ := http.NewRequest("GET", "/symptoms/date-range?startDate="+url.QueryEscape(startDate)+"&endDate="+url.QueryEscape(endDate), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response body: %s", http.StatusOK, w.Code, w.Body.String())
		return
	}

	var response map[string][]symptomDomain.Symptom
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

	category := shared.SymptomCategoryDigestive

	symptom1 := symptomDomain.NewSymptom("test-user", "Digestive", 5, time.Now())
	symptom1.Category = &category

	symptom2 := symptomDomain.NewSymptom("test-user", "Other", 5, time.Now())
	// No category

	_, _ = repo.Create(context.TODO(), symptom1)
	_, _ = repo.Create(context.TODO(), symptom2)

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

	var response map[string][]symptomDomain.Symptom
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
