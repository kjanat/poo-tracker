package medication

import (
	"context"
	"testing"

	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(&medication.MedicationDB{})
	return db
}

func TestCreateAndGetMedication(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMedicationRepository(db)
	ctx := context.Background()
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	m := &medication.Medication{ID: validUUID, UserID: "user-1", Name: "Aspirin"}
	assert.NoError(t, repo.Create(ctx, m))
	got, err := repo.GetByID(ctx, validUUID)
	assert.NoError(t, err)
	assert.Equal(t, validUUID, got.ID)
}
