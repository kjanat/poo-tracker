package repository

import (
	"context"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
)

func TestMedicationRepository_Create(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	medication := medication.NewMedication("user1", "Ibuprofen", "200mg", "twice daily")

	created, err := repo.Create(ctx, medication)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if created.ID == "" {
		t.Error("Expected ID to be generated")
	}

	if created.Name != "Ibuprofen" {
		t.Errorf("Expected name 'Ibuprofen', got %s", created.Name)
	}

	if created.Dosage != "200mg" {
		t.Errorf("Expected dosage '200mg', got %s", created.Dosage)
	}

	if !created.IsActive {
		t.Error("Expected medication to be active by default")
	}
}

func TestMedicationRepository_GetByID(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	medication := medication.NewMedication("user1", "Aspirin", "100mg", "daily")
	created, _ := repo.Create(ctx, medication)

	retrieved, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, retrieved.ID)
	}

	if retrieved.Name != "Aspirin" {
		t.Errorf("Expected name 'Aspirin', got %s", retrieved.Name)
	}
}

func TestMedicationRepository_GetByID_NotFound(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent medication")
	}
}

func TestMedicationRepository_GetByUserID(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	// Create medications for different users
	med1 := medication.NewMedication("user1", "Medication 1", "10mg", "daily")
	med2 := medication.NewMedication("user1", "Medication 2", "20mg", "twice daily")
	med3 := medication.NewMedication("user2", "Medication 3", "30mg", "weekly")

	// Create test medications (ignore errors in test setup)
	_, _ = repo.Create(ctx, med1)
	_, _ = repo.Create(ctx, med2)
	_, _ = repo.Create(ctx, med3)

	medications, err := repo.GetByUserID(ctx, "user1", 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(medications) != 2 {
		t.Errorf("Expected 2 medications for user1, got %d", len(medications))
	}

	// Check that medications are sorted by creation time (newest first)
	for i := 0; i < len(medications)-1; i++ {
		if medications[i].CreatedAt.Before(medications[i+1].CreatedAt) {
			t.Error("Expected medications to be sorted by creation time (newest first)")
		}
	}
}

func TestMedicationRepository_Update(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	medication := medication.NewMedication("user1", "Original", "10mg", "daily")
	created, _ := repo.Create(ctx, medication)

	updates := medication.MedicationUpdate{
		Name:     stringPtr("Updated"),
		Dosage:   stringPtr("20mg"),
		Notes:    stringPtr("Updated notes"),
		IsActive: boolPtr(false),
	}

	updated, err := repo.Update(ctx, created.ID, updates)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got %s", updated.Name)
	}

	if updated.Dosage != "20mg" {
		t.Errorf("Expected dosage '20mg', got %s", updated.Dosage)
	}

	if updated.Notes != "Updated notes" {
		t.Errorf("Expected notes 'Updated notes', got %s", updated.Notes)
	}

	if updated.IsActive {
		t.Error("Expected medication to be inactive after update")
	}
}

func TestMedicationRepository_Delete(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	medication := medication.NewMedication("user1", "To Delete", "10mg", "daily")
	created, _ := repo.Create(ctx, medication)

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted medication")
	}
}

func TestMedicationRepository_GetActiveByUserID(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	// Create active and inactive medications
	activeMed := medication.NewMedication("user1", "Active Med", "10mg", "daily")
	activeMed.IsActive = true

	inactiveMed := medication.NewMedication("user1", "Inactive Med", "20mg", "daily")
	inactiveMed.IsActive = false

	_, _ = repo.Create(ctx, activeMed)
	_, _ = repo.Create(ctx, inactiveMed)

	medications, err := repo.GetActiveByUserID(ctx, "user1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(medications) != 1 {
		t.Errorf("Expected 1 active medication, got %d", len(medications))
	}

	if medications[0].Name != "Active Med" {
		t.Errorf("Expected 'Active Med', got %s", medications[0].Name)
	}
}

func TestMedicationRepository_GetByUserIDAndCategory(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	category := medication.MedicationCategoryPainRelief

	med1 := medication.NewMedication("user1", "Pain Med", "10mg", "daily")
	med1.Category = &category

	med2 := medication.NewMedication("user1", "Other Med", "20mg", "daily")
	// No category set

	_, _ = repo.Create(ctx, med1)
	_, _ = repo.Create(ctx, med2)

	medications, err := repo.GetByUserIDAndCategory(ctx, "user1", category)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(medications) != 1 {
		t.Errorf("Expected 1 medication with pain relief category, got %d", len(medications))
	}

	if medications[0].Name != "Pain Med" {
		t.Errorf("Expected 'Pain Med', got %s", medications[0].Name)
	}
}

func TestMedicationRepository_MarkAsTaken(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	medication := medication.NewMedication("user1", "Test Med", "10mg", "daily")
	created, _ := repo.Create(ctx, medication)

	takenTime := time.Now()
	err := repo.MarkAsTaken(ctx, created.ID, takenTime)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updated, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.TakenAt == nil {
		t.Error("Expected TakenAt to be set")
	}

	if !updated.TakenAt.Equal(takenTime) {
		t.Errorf("Expected TakenAt to be %v, got %v", takenTime, *updated.TakenAt)
	}
}

func TestMedicationRepository_MarkAsTaken_NotFound(t *testing.T) {
	repo := NewMemoryMedicationRepository()
	ctx := context.Background()

	err := repo.MarkAsTaken(ctx, "nonexistent", time.Now())
	if err == nil {
		t.Error("Expected error for nonexistent medication")
	}
}

// Helper function for boolean pointer creation
func boolPtr(b bool) *bool {
	return &b
}
