package medication

import (
	"testing"
)

func TestNewMedicationDefaults(t *testing.T) {
	m := NewMedication("user1", "Ibuprofen", "200mg", "daily")

	if m.UserID != "user1" {
		t.Errorf("expected userID user1, got %s", m.UserID)
	}
	if m.Name != "Ibuprofen" {
		t.Errorf("expected name Ibuprofen, got %s", m.Name)
	}
	if m.Dosage != "200mg" {
		t.Errorf("expected dosage 200mg, got %s", m.Dosage)
	}
	if m.Frequency != "daily" {
		t.Errorf("expected frequency daily, got %s", m.Frequency)
	}
	if !m.IsActive {
		t.Error("expected IsActive true by default")
	}
	if m.IsAsNeeded {
		t.Error("expected IsAsNeeded false by default")
	}
	if len(m.SideEffects) != 0 {
		t.Error("expected empty SideEffects slice")
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}
