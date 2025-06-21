package repository

import (
	"context"
	"testing"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
)

func TestBowelDetailsRepo_CRUD(t *testing.T) {
	bowelRepo := NewMemoryBowelRepo()
	detailsRepo := NewMemoryBowelDetailsRepo(bowelRepo)
	ctx := context.Background()

	// First create a bowel movement
	movement := bm.BowelMovement{
	createdBM, err := bowelRepo.Create(ctx, movement)
        BristolType: 4,
    }
    createdBM, err := bowelRepo.Create(ctx, initialBM)
	createdBM, err := bowelRepo.Create(ctx, bm)
	if err != nil {
		t.Fatalf("Failed to create bowel movement: %v", err)
	}

	// Verify HasDetails is initially false
	if createdBM.HasDetails {
		t.Error("Expected HasDetails to be false initially")
	}

	// Test Create
	details := bm.NewBowelMovementDetails(createdBM.ID)
	details.Notes = "Test notes"
	details.DetailedNotes = "Very detailed notes"
	details.Environment = "Private bathroom"
	stressLevel := 3
	details.StressLevel = &stressLevel

	created, err := detailsRepo.Create(ctx, details)
	if err != nil {
		t.Fatalf("Failed to create details: %v", err)
	}

	if created.BowelMovementID != createdBM.ID {
		t.Errorf("Expected BowelMovementID %s, got %s", createdBM.ID, created.BowelMovementID)
	}
	if created.Notes != "Test notes" {
		t.Errorf("Expected Notes 'Test notes', got %s", created.Notes)
	}

	// Verify HasDetails is now true in the bowel movement
	updatedBM, err := bowelRepo.Get(ctx, createdBM.ID)
	if err != nil {
		t.Fatalf("Failed to get updated bowel movement: %v", err)
	}
	if !updatedBM.HasDetails {
		t.Error("Expected HasDetails to be true after creating details")
	}

	// Test Get
	retrieved, err := detailsRepo.Get(ctx, createdBM.ID)
	if err != nil {
		t.Fatalf("Failed to get details: %v", err)
	}
	if retrieved.Notes != "Test notes" {
		t.Errorf("Expected Notes 'Test notes', got %s", retrieved.Notes)
	}

	// Test Exists
	exists, err := detailsRepo.Exists(ctx, createdBM.ID)
	if err != nil {
		t.Fatalf("Failed to check exists: %v", err)
	}
	if !exists {
		t.Error("Expected details to exist")
	}

	// Test Update
	newNotes := "Updated notes"
	newEnvironment := "Public bathroom"
	update := bm.BowelMovementDetailsUpdate{
		Notes:       &newNotes,
		Environment: &newEnvironment,
	}

	updated, err := detailsRepo.Update(ctx, createdBM.ID, update)
	if err != nil {
		t.Fatalf("Failed to update details: %v", err)
	}
	if updated.Notes != "Updated notes" {
		t.Errorf("Expected Notes 'Updated notes', got %s", updated.Notes)
	}
	if updated.Environment != "Public bathroom" {
		t.Errorf("Expected Environment 'Public bathroom', got %s", updated.Environment)
	}
	// DetailedNotes should remain unchanged
	if updated.DetailedNotes != "Very detailed notes" {
		t.Errorf("Expected DetailedNotes to remain 'Very detailed notes', got %s", updated.DetailedNotes)
	}

	// Test Delete
	err = detailsRepo.Delete(ctx, createdBM.ID)
	if err != nil {
		t.Fatalf("Failed to delete details: %v", err)
	}

	// Verify HasDetails is now false in the bowel movement
	updatedBM2, err := bowelRepo.Get(ctx, createdBM.ID)
	if err != nil {
		t.Fatalf("Failed to get bowel movement after delete: %v", err)
	}
	if updatedBM2.HasDetails {
		t.Error("Expected HasDetails to be false after deleting details")
	}

	// Verify details no longer exist
	exists, err = detailsRepo.Exists(ctx, createdBM.ID)
	if err != nil {
		t.Fatalf("Failed to check exists after delete: %v", err)
	}
	if exists {
		t.Error("Expected details to not exist after delete")
	}

	// Verify Get returns ErrNotFound
	_, err = detailsRepo.Get(ctx, createdBM.ID)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestBowelDetailsRepo_NotFound(t *testing.T) {
	bowelRepo := NewMemoryBowelRepo()
	detailsRepo := NewMemoryBowelDetailsRepo(bowelRepo)
	ctx := context.Background()

	// Test Get with non-existent ID
	_, err := detailsRepo.Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	// Test Update with non-existent ID
	update := bm.BowelMovementDetailsUpdate{
		Notes: &[]string{"test"}[0],
	}
	_, err = detailsRepo.Update(ctx, "nonexistent", update)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	// Test Delete with non-existent ID
	err = detailsRepo.Delete(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	// Test Exists with non-existent ID
	exists, err := detailsRepo.Exists(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if exists {
		t.Error("Expected exists to be false for non-existent ID")
	}
}
