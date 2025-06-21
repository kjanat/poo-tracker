package repository

import (
	"context"
	"testing"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
)

func TestMemoryRepo_UpdateWithZeroValues(t *testing.T) {
	repo := NewMemoryBowelRepo()
	ctx := context.Background()

	// Create initial bowel movement
	initial := bm.BowelMovement{
		UserID:      "user1",
		BristolType: 5,
	}
	created, err := repo.Create(ctx, initial)
	if err != nil {
		t.Fatalf("Failed to create initial entry: %v", err)
	}

	// Test updating to zero values - this should work now
	bristolZero := 0
	painZero := 0
	update := bm.BowelMovementUpdate{
		BristolType: &bristolZero,
		Pain:        &painZero,
	}

	updated, err := repo.Update(ctx, created.ID, update)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	// Verify zero values were applied
	if updated.BristolType != 0 {
		t.Errorf("Expected BristolType to be 0, got %d", updated.BristolType)
	}
	if updated.Pain != 0 {
		t.Errorf("Expected Pain to be 0, got %d", updated.Pain)
	}

	// Test partial update - only update BristolType, leave Pain as is
	newBristol := 3
	partialUpdate := bm.BowelMovementUpdate{
		BristolType: &newBristol,
		// Pain is nil, so it should not be updated
	}

	updated2, err := repo.Update(ctx, created.ID, partialUpdate)
	if err != nil {
		t.Fatalf("Failed to partial update entry: %v", err)
	}

	// Verify only BristolType was updated, Pain remained 0
	if updated2.BristolType != 3 {
		t.Errorf("Expected BristolType to be 3, got %d", updated2.BristolType)
	}
	if updated2.Pain != 0 {
		t.Errorf("Expected Pain to remain 0, got %d", updated2.Pain)
	}

	// Test updating Pain back to non-zero while keeping BristolType
	newPain := 5
	painUpdate := bm.BowelMovementUpdate{
		Pain: &newPain,
		// BristolType is nil, so it should not be updated
	}

	updated3, err := repo.Update(ctx, created.ID, painUpdate)
	if err != nil {
		t.Fatalf("Failed to update pain: %v", err)
	}

	// Verify only Pain was updated, BristolType remained 3
	if updated3.BristolType != 3 {
		t.Errorf("Expected BristolType to remain 3, got %d", updated3.BristolType)
	}
	if updated3.Pain != 5 {
		t.Errorf("Expected Pain to be 5, got %d", updated3.Pain)
	}
}
