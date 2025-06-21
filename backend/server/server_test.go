package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

func TestHealth(t *testing.T) {
	repo := repository.NewMemoryBowelRepo()
	meals := repository.NewMemoryMealRepo()

	// Create a mock auth service for testing
	userRepo := repository.NewMemoryUserRepository()
	authService := &service.JWTAuthService{
		UserRepo: userRepo,
		Secret:   "test-secret",
		Expiry:   24 * time.Hour,
	}

	detailsRepo := repository.NewMemoryBowelDetailsRepo(repo)
	symptomRepo := repository.NewMemorySymptomRepository()
	medicationRepo := repository.NewMemoryMedicationRepository()
	mealBowelRelationRepo := repository.NewMemoryMealBowelMovementRelationRepository()
	mealSymptomRelationRepo := repository.NewMemoryMealSymptomRelationRepository()
	twoFactorRepo := repository.NewMemoryUserTwoFactorRepository()
	app := New(repo, detailsRepo, meals, symptomRepo, medicationRepo, mealBowelRelationRepo, mealSymptomRelationRepo, service.AvgBristol{}, authService, twoFactorRepo, userRepo)
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	app.Engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
