package meal

import (
	"context"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(&meal.MealDB{})
	return db
}

func TestCreateAndGetMeal(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMealRepository(db)
	ctx := context.Background()
	m := &meal.Meal{ID: "test-id", UserID: "user-1", Name: "Lunch"}
	assert.NoError(t, repo.Create(ctx, m))
	got, err := repo.GetByID(ctx, "test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", got.ID)
}
