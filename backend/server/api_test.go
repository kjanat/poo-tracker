package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBowelMovementCRUD(t *testing.T) {
	r := New()

	// create
	body := `{"userId":"u1","bristolType":3,"notes":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/bowel-movements", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created BowelMovement
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	// get
	req = httptest.NewRequest(http.MethodGet, "/api/bowel-movements/"+created.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// update
	updBody := `{"bristolType":4}`
	req = httptest.NewRequest(http.MethodPut, "/api/bowel-movements/"+created.ID, bytes.NewBufferString(updBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// analytics
	req = httptest.NewRequest(http.MethodGet, "/api/analytics", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/api/bowel-movements/"+created.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
