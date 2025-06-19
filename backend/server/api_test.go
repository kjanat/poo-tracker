package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

func setup() *App {
	repo := repository.NewMemory()
	strategy := service.AvgBristol{}
	return New(repo, strategy)
}

func TestBowelMovementCRUD(t *testing.T) {
	app := setup()
	r := app.Engine

	body := `{"userId":"u1","bristolType":3,"notes":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/bowel-movements", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	id := created["id"].(string)

	req = httptest.NewRequest(http.MethodGet, "/api/bowel-movements/"+id, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	updBody := `{"bristolType":4}`
	req = httptest.NewRequest(http.MethodPut, "/api/bowel-movements/"+id, bytes.NewBufferString(updBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/analytics", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/bowel-movements/"+id, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
