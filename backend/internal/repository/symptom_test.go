package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

func TestSymptomRepository_Create(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	sym := symptom.NewSymptom("user1", "Headache", 7, time.Now())

	created, err := repo.Create(ctx, sym)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if created.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if created.Name != "Headache" {
		t.Errorf("Expected name 'Headache', got %s", created.Name)
	}

	if created.Severity != 7 {
		t.Errorf("Expected severity 7, got %d", created.Severity)
	}
}

func TestSymptomRepository_GetByID(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	sym := symptom.NewSymptom("user1", "Nausea", 5, time.Now())
	created, _ := repo.Create(ctx, sym)

	retrieved, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, retrieved.ID)
	}

	if retrieved.Name != "Nausea" {
		t.Errorf("Expected name 'Nausea', got %s", retrieved.Name)
	}
}

func TestSymptomRepository_GetByID_NotFound(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent symptom")
	}
}

func TestSymptomRepository_GetByUserID(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	// Create symptoms for different users
	symptom1 := symptom.NewSymptom("user1", "Headache", 5, time.Now())
	symptom2 := symptom.NewSymptom("user1", "Nausea", 5, time.Now())
	symptom3 := symptom.NewSymptom("user2", "Pain", 5, time.Now())

	// Create test data (ignore errors in test setup)
	_, _ = repo.Create(ctx, symptom1)
	_, _ = repo.Create(ctx, symptom2)
	_, _ = repo.Create(ctx, symptom3)

	symptoms, err := repo.GetByUserID(ctx, "user1", 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(symptoms) != 2 {
		t.Errorf("Expected 2 symptoms for user1, got %d", len(symptoms))
	}

	// Check that symptoms are sorted by recorded time (newest first)
	for i := 0; i < len(symptoms)-1; i++ {
		if symptoms[i].RecordedAt.Before(symptoms[i+1].RecordedAt) {
			t.Error("Expected symptoms to be sorted by recorded time (newest first)")
		}
	}
}

func TestSymptomRepository_Update(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	sym := symptom.NewSymptom("user1", "Original", 5, time.Now())
	created, _ := repo.Create(ctx, sym)

	updates := symptom.SymptomUpdate{
		Name:     stringPtr("Updated"),
		Severity: intPtr(8),
		Notes:    stringPtr("Updated notes"),
	}

	updated, err := repo.Update(ctx, created.ID, updates)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got %s", updated.Name)
	}

	if updated.Severity != 8 {
		t.Errorf("Expected severity 8, got %d", updated.Severity)
	}

	if updated.Notes != "Updated notes" {
		t.Errorf("Expected notes 'Updated notes', got %s", updated.Notes)
	}
}

func TestSymptomRepository_Delete(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	sym := symptom.NewSymptom("user1", "To Delete", 5, time.Now())
	created, _ := repo.Create(ctx, sym)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted symptom")
	}
}

func TestSymptomRepository_GetByUserIDAndDateRange(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	// Create symptoms at different times
	symptom1 := symptom.NewSymptom("user1", "Yesterday", 5, yesterday)
	symptom2 := symptom.NewSymptom("user1", "Today", 5, now)
	symptom3 := symptom.NewSymptom("user1", "Tomorrow", 5, tomorrow)

	_, _ = repo.Create(ctx, symptom1)
	_, _ = repo.Create(ctx, symptom2)
	_, _ = repo.Create(ctx, symptom3)

	// Query for symptoms from yesterday to today
	symptoms, err := repo.GetByUserIDAndDateRange(ctx, "user1", yesterday.Add(-time.Hour), now.Add(time.Hour))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(symptoms) != 2 {
		t.Errorf("Expected 2 symptoms in date range, got %d", len(symptoms))
	}
}

func TestSymptomRepository_GetByUserIDAndCategory(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	category := shared.SymptomCategoryDigestive

	symptom1 := symptom.NewSymptom("user1", "Digestive Issue", 5, time.Now())
	symptom1.Category = &category

	symptom2 := symptom.NewSymptom("user1", "Other Issue", 5, time.Now())
	// No category set

	_, _ = repo.Create(ctx, symptom1)
	_, _ = repo.Create(ctx, symptom2)

	symptoms, err := repo.GetByUserIDAndCategory(ctx, "user1", category)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(symptoms) != 1 {
		t.Errorf("Expected 1 symptom with digestive category, got %d", len(symptoms))
	}

	if symptoms[0].Name != "Digestive Issue" {
		t.Errorf("Expected 'Digestive Issue', got %s", symptoms[0].Name)
	}
}

func TestSymptomRepository_GetByUserIDAndType(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	symptomType := shared.SymptomNausea

	symptom1 := symptom.NewSymptom("user1", "Nausea Symptom", 5, time.Now())
	symptom1.Type = &symptomType

	symptom2 := symptom.NewSymptom("user1", "Other Symptom", 5, time.Now())
	// No type set

	_, _ = repo.Create(ctx, symptom1)
	_, _ = repo.Create(ctx, symptom2)

	symptoms, err := repo.GetByUserIDAndType(ctx, "user1", symptomType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(symptoms) != 1 {
		t.Errorf("Expected 1 symptom with nausea type, got %d", len(symptoms))
	}

	if symptoms[0].Name != "Nausea Symptom" {
		t.Errorf("Expected 'Nausea Symptom', got %s", symptoms[0].Name)
	}
}

// Helper functions for pointer creation
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
