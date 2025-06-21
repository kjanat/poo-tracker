package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
)

func TestBowelMovementDetailsAPI(t *testing.T) {
	app := setup()
	r := app.Engine

	// First create a bowel movement to attach details to
	bmBody := `{"userId":"u1","bristolType":3}`
	req := httptest.NewRequest(http.MethodPost, "/api/bowel-movements", bytes.NewReader([]byte(bmBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var createdBM bm.BowelMovement
	if err := json.Unmarshal(w.Body.Bytes(), &createdBM); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Test creating details
	detailsBody := `{
		"notes": "Test notes",
		"detailedNotes": "Very detailed notes",
		"environment": "Private bathroom",
		"stressLevel": 3,
		"tags": ["morning", "routine"]
	}`

	req = httptest.NewRequest(http.MethodPost, "/api/bowel-movements/"+createdBM.ID+"/details", bytes.NewReader([]byte(detailsBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 for details creation, got %d", w.Code)
	}

	var createdDetails bm.BowelMovementDetails
	if err := json.Unmarshal(w.Body.Bytes(), &createdDetails); err != nil {
		t.Fatalf("Failed to unmarshal details response: %v", err)
	}

	if createdDetails.Notes != "Test notes" {
		t.Errorf("Expected notes 'Test notes', got %s", createdDetails.Notes)
	}
	if createdDetails.BowelMovementID != createdBM.ID {
		t.Errorf("Expected BowelMovementID %s, got %s", createdBM.ID, createdDetails.BowelMovementID)
	}

	// Test duplicate creation (should fail)
	req = httptest.NewRequest(http.MethodPost, "/api/bowel-movements/"+createdBM.ID+"/details", bytes.NewReader([]byte(detailsBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate details, got %d", w.Code)
	}

	// Test getting details
	req = httptest.NewRequest(http.MethodGet, "/api/bowel-movements/"+createdBM.ID+"/details", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for getting details, got %d", w.Code)
	}

	var retrievedDetails bm.BowelMovementDetails
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedDetails); err != nil {
		t.Fatalf("Failed to unmarshal retrieved details: %v", err)
	}

	if retrievedDetails.Notes != "Test notes" {
		t.Errorf("Expected retrieved notes 'Test notes', got %s", retrievedDetails.Notes)
	}

	// Test updating details
	updateBody := `{
		"notes": "Updated notes",
		"environment": "Public bathroom"
	}`

	req = httptest.NewRequest(http.MethodPut, "/api/bowel-movements/"+createdBM.ID+"/details", bytes.NewReader([]byte(updateBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for updating details, got %d", w.Code)
	}

	var updatedDetails bm.BowelMovementDetails
	if err := json.Unmarshal(w.Body.Bytes(), &updatedDetails); err != nil {
		t.Fatalf("Failed to unmarshal updated details: %v", err)
	}

	if updatedDetails.Notes != "Updated notes" {
		t.Errorf("Expected updated notes 'Updated notes', got %s", updatedDetails.Notes)
	}
	if updatedDetails.Environment != "Public bathroom" {
		t.Errorf("Expected updated environment 'Public bathroom', got %s", updatedDetails.Environment)
	}
	// DetailedNotes should remain unchanged
	if updatedDetails.DetailedNotes != "Very detailed notes" {
		t.Errorf("Expected DetailedNotes to remain unchanged 'Very detailed notes', got %s", updatedDetails.DetailedNotes)
	}

	// Test deleting details
	req = httptest.NewRequest(http.MethodDelete, "/api/bowel-movements/"+createdBM.ID+"/details", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 for deleting details, got %d", w.Code)
	}

	// Verify details are gone
	req = httptest.NewRequest(http.MethodGet, "/api/bowel-movements/"+createdBM.ID+"/details", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 after deleting details, got %d", w.Code)
	}
}

func TestBowelMovementDetailsAPI_NotFound(t *testing.T) {
	app := setup()
	r := app.Engine

	// Test creating details for non-existent bowel movement
	detailsBody := `{"notes": "Test notes"}`
	req := httptest.NewRequest(http.MethodPost, "/api/bowel-movements/nonexistent/details", bytes.NewReader([]byte(detailsBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent bowel movement, got %d", w.Code)
	}

	// Test getting details for non-existent bowel movement
	req = httptest.NewRequest(http.MethodGet, "/api/bowel-movements/nonexistent/details", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for getting non-existent details, got %d", w.Code)
	}

	// Test updating details for non-existent bowel movement
	updateBody := `{"notes": "Updated notes"}`
	req = httptest.NewRequest(http.MethodPut, "/api/bowel-movements/nonexistent/details", bytes.NewReader([]byte(updateBody)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for updating non-existent details, got %d", w.Code)
	}

	// Test deleting details for non-existent bowel movement
	req = httptest.NewRequest(http.MethodDelete, "/api/bowel-movements/nonexistent/details", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for deleting non-existent details, got %d", w.Code)
	}
}
