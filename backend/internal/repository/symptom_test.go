package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

func TestSymptomRepository_Create(t *testing.T) {
	repo := NewMemorySymptomRepository()
	ctx := context.Background()

	symptom := symptom.NewSymptom("user1", "Headache")
	symptom.Severity = 7

	created, err := repo.Create(ctx, symptom)
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

	symptom := symptom.NewSymptom("user1", "Nausea")
	created, _ := repo.Create(ctx, symptom)

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
	symptom1 := symptom.NewSymptom("user1", "Headache")
	symptom2 := symptom.NewSymptom("user1", "Nausea")
	symptom3 := symptom.NewSymptom("user2", "Pain")

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

	symptom := symptom.NewSymptom("user1", "Original")
	created, _ := repo.Create(ctx, symptom)

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

	symptom := symptom.NewSymptom("user1", "To Delete")
	created, _ := repo.Create(ctx, symptom)

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
	symptom1 := symptom.NewSymptom("user1", "Yesterday")
	symptom1.RecordedAt = yesterday

	symptom2 := symptom.NewSymptom("user1", "Today")
	symptom2.RecordedAt = now

	symptom3 := symptom.NewSymptom("user1", "Tomorrow")
	symptom3.RecordedAt = tomorrow

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

	category := symptom.SymptomCategoryDigestive

	symptom1 := symptom.NewSymptom("user1", "Digestive Issue")
	symptom1.Category = &category

	symptom2 := symptom.NewSymptom("user1", "Other Issue")
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

	symptomType := symptom.SymptomNausea

	symptom1 := symptom.NewSymptom("user1", "Nausea Symptom")
	symptom1.Type = &symptomType

	symptom2 := symptom.NewSymptom("user1", "Other Symptom")
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
