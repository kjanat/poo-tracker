package repository

import (
	"context"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/model"
)

func TestMemoryRepo_UpdateWithZeroValues(t *testing.T) {
	repo := NewMemory()
	ctx := context.Background()

	// Create initial bowel movement
	initial := model.BowelMovement{
		UserID:      "user1",
		BristolType: 5,
		Notes:       "Initial notes",
	}
	created, err := repo.Create(ctx, initial)
	if err != nil {
		t.Fatalf("Failed to create initial entry: %v", err)
	}

	// Test updating to zero values - this should work now
	bristolZero := 0
	emptyNotes := ""
	update := model.BowelMovementUpdate{
		BristolType: &bristolZero,
		Notes:       &emptyNotes,
	}

	updated, err := repo.Update(ctx, created.ID, update)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	// Verify zero values were applied
	if updated.BristolType != 0 {
		t.Errorf("Expected BristolType to be 0, got %d", updated.BristolType)
	}
	if updated.Notes != "" {
		t.Errorf("Expected Notes to be empty, got %q", updated.Notes)
	}

	// Test partial update - only update BristolType, leave Notes as is
	newBristol := 3
	partialUpdate := model.BowelMovementUpdate{
		BristolType: &newBristol,
		// Notes is nil, so it should not be updated
	}

	updated2, err := repo.Update(ctx, created.ID, partialUpdate)
	if err != nil {
		t.Fatalf("Failed to partial update entry: %v", err)
	}

	// Verify only BristolType was updated, Notes remained empty
	if updated2.BristolType != 3 {
		t.Errorf("Expected BristolType to be 3, got %d", updated2.BristolType)
	}
	if updated2.Notes != "" {
		t.Errorf("Expected Notes to remain empty, got %q", updated2.Notes)
	}

	// Test updating Notes back to non-empty while keeping BristolType
	newNotes := "Updated notes"
	notesUpdate := model.BowelMovementUpdate{
		Notes: &newNotes,
		// BristolType is nil, so it should not be updated
	}

	updated3, err := repo.Update(ctx, created.ID, notesUpdate)
	if err != nil {
		t.Fatalf("Failed to update notes: %v", err)
	}

	// Verify only Notes was updated, BristolType remained 3
	if updated3.BristolType != 3 {
		t.Errorf("Expected BristolType to remain 3, got %d", updated3.BristolType)
	}
	if updated3.Notes != "Updated notes" {
		t.Errorf("Expected Notes to be 'Updated notes', got %q", updated3.Notes)
	}
}
