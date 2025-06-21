package medication

import "testing"

func TestNewMedicationDefaults(t *testing.T) {
	m := NewMedication("user1", "Name", "10mg", "daily")
	if m.UserID != "user1" {
		t.Errorf("expected user ID user1, got %s", m.UserID)
	}
	if m.Name != "Name" || m.Dosage != "10mg" || m.Frequency != "daily" {
		t.Errorf("unexpected basic fields: %+v", m)
	}
	if !m.IsActive {
		t.Error("expected IsActive true")
	}
	if m.IsAsNeeded {
		t.Error("expected IsAsNeeded false")
	}
	if len(m.SideEffects) != 0 {
		t.Errorf("expected no side effects, got %d", len(m.SideEffects))
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		t.Error("timestamps should be set")
	}
}
