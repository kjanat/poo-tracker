package bowelmovement

import (
	"context"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(&bowelmovement.BowelMovementDB{})
	return db
}

func TestCreateAndGetBowelMovement(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBowelMovementRepository(db)
	ctx := context.Background()
	bm := &bowelmovement.BowelMovement{
		ID:           "test-id",
		UserID:       "user-1",
		BristolType:  4, // Use a valid Bristol type (1-7)
		Pain:         1,
		Strain:       1,
		Satisfaction: 5,
		Floaters:     false,
	}
	assert.NoError(t, repo.Create(ctx, bm))
	got, err := repo.GetByID(ctx, "test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", got.ID)
}
