package symptom

import (
	"testing"
	"time"
)

func TestNewSymptomDefaults(t *testing.T) {
	recorded := time.Now()
	s := NewSymptom("user1", "Nausea", 5, recorded)

	if s.UserID != "user1" {
		t.Errorf("expected userID user1, got %s", s.UserID)
	}
	if s.Name != "Nausea" {
		t.Errorf("expected name Nausea, got %s", s.Name)
	}
	if s.Severity != 5 {
		t.Errorf("expected severity 5, got %d", s.Severity)
	}
	if !s.RecordedAt.Equal(recorded) {
		t.Errorf("expected recordedAt %v, got %v", recorded, s.RecordedAt)
	}
	if len(s.Triggers) != 0 {
		t.Error("expected empty Triggers slice")
	}
	if s.CreatedAt.IsZero() || s.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}
