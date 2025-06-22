package relation

import (
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/relations"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(&MealBowelMovementRelationGORM{}, &MealSymptomRelationGORM{})
	if err != nil {
		panic("failed to migrate test database")
	}

	return db
}

func TestMealBowelMovementRelationRepository_CRUD(t *testing.T) {
	db := setupTestDB()
	repo := NewMealBowelMovementRelationRepository(db)

	// Create
	relation := &relations.MealBowelMovementRelation{
		UserID:          "user-123",
		MealID:          "meal-123",
		BowelMovementID: "bowel-123",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Strength:        7,
		Notes:           "Test relation",
		TimeGapHours:    2.5,
	}

	err := repo.Create(relation)
	assert.NoError(t, err)
	assert.NotEmpty(t, relation.ID)

	// Read
	retrieved, err := repo.GetByID(relation.ID)
	assert.NoError(t, err)
	assert.Equal(t, relation.UserID, retrieved.UserID)
	assert.Equal(t, relation.MealID, retrieved.MealID)
	assert.Equal(t, relation.BowelMovementID, retrieved.BowelMovementID)
	assert.Equal(t, relation.Strength, retrieved.Strength)
	assert.Equal(t, relation.Notes, retrieved.Notes)
	assert.Equal(t, relation.TimeGapHours, retrieved.TimeGapHours)

	// Update
	retrieved.Strength = 9
	retrieved.Notes = "Updated notes"
	err = repo.Update(retrieved)
	assert.NoError(t, err)

	updated, err := repo.GetByID(relation.ID)
	assert.NoError(t, err)
	assert.Equal(t, 9, updated.Strength)
	assert.Equal(t, "Updated notes", updated.Notes)

	// Delete
	err = repo.Delete(relation.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(relation.ID)
	assert.Error(t, err)
}

func TestMealBowelMovementRelationRepository_GetByUserID(t *testing.T) {
	db := setupTestDB()
	repo := NewMealBowelMovementRelationRepository(db)

	userID := "user-123"

	// Create multiple relations for the same user
	relation1 := &relations.MealBowelMovementRelation{
		UserID:          userID,
		MealID:          "meal-1",
		BowelMovementID: "bowel-1",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Strength:        5,
		TimeGapHours:    1.0,
	}

	relation2 := &relations.MealBowelMovementRelation{
		UserID:          userID,
		MealID:          "meal-2",
		BowelMovementID: "bowel-2",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Strength:        8,
		TimeGapHours:    3.0,
	}

	err := repo.Create(relation1)
	assert.NoError(t, err)

	err = repo.Create(relation2)
	assert.NoError(t, err)

	// Get all relations for the user
	relations, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Len(t, relations, 2)
}

func TestMealSymptomRelationRepository_CRUD(t *testing.T) {
	db := setupTestDB()
	repo := NewMealSymptomRelationRepository(db)

	// Create
	relation := &relations.MealSymptomRelation{
		UserID:       "user-123",
		MealID:       "meal-123",
		SymptomID:    "symptom-123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Strength:     6,
		Notes:        "Test symptom relation",
		TimeGapHours: 1.5,
	}

	err := repo.Create(relation)
	assert.NoError(t, err)
	assert.NotEmpty(t, relation.ID)

	// Read
	retrieved, err := repo.GetByID(relation.ID)
	assert.NoError(t, err)
	assert.Equal(t, relation.UserID, retrieved.UserID)
	assert.Equal(t, relation.MealID, retrieved.MealID)
	assert.Equal(t, relation.SymptomID, retrieved.SymptomID)
	assert.Equal(t, relation.Strength, retrieved.Strength)
	assert.Equal(t, relation.Notes, retrieved.Notes)
	assert.Equal(t, relation.TimeGapHours, retrieved.TimeGapHours)

	// Update
	retrieved.Strength = 4
	retrieved.Notes = "Updated symptom notes"
	err = repo.Update(retrieved)
	assert.NoError(t, err)

	updated, err := repo.GetByID(relation.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4, updated.Strength)
	assert.Equal(t, "Updated symptom notes", updated.Notes)

	// Delete
	err = repo.Delete(relation.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(relation.ID)
	assert.Error(t, err)
}

func TestMealSymptomRelationRepository_GetByMealID(t *testing.T) {
	db := setupTestDB()
	repo := NewMealSymptomRelationRepository(db)

	mealID := "meal-123"

	// Create multiple relations for the same meal
	relation1 := &relations.MealSymptomRelation{
		UserID:       "user-1",
		MealID:       mealID,
		SymptomID:    "symptom-1",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Strength:     5,
		TimeGapHours: 1.0,
	}

	relation2 := &relations.MealSymptomRelation{
		UserID:       "user-1",
		MealID:       mealID,
		SymptomID:    "symptom-2",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Strength:     7,
		TimeGapHours: 2.0,
	}

	err := repo.Create(relation1)
	assert.NoError(t, err)

	err = repo.Create(relation2)
	assert.NoError(t, err)

	// Get all relations for the meal
	relations, err := repo.GetByMealID(mealID)
	assert.NoError(t, err)
	assert.Len(t, relations, 2)
}
