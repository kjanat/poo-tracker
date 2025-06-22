package symptom

import (
	"context"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(&symptom.SymptomDB{})
	return db
}

func TestCreateAndGetSymptom(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSymptomRepository(db)
	ctx := context.Background()
	validUUID := "123e4567-e89b-12d3-a456-426614174001"
	s := &symptom.Symptom{ID: validUUID, UserID: "user-1", Name: "Headache"}
	assert.NoError(t, repo.Create(ctx, s))
	got, err := repo.GetByID(ctx, validUUID)
	assert.NoError(t, err)
	assert.Equal(t, validUUID, got.ID)
}
