package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

func TestHealth(t *testing.T) {
	repo := repository.NewMemory()
	app := New(repo, repo, service.AvgBristol{})
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
