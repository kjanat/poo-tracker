package symptom

import (
	"testing"
	"time"
)

func TestNewSymptomDefaults(t *testing.T) {
	now := time.Now()
	s := NewSymptom("u1", "Headache", 4, now)
	if s.UserID != "u1" || s.Name != "Headache" || s.Severity != 4 {
		t.Errorf("unexpected basic fields: %+v", s)
	}
	if !s.RecordedAt.Equal(now) {
		t.Errorf("expected RecordedAt %v, got %v", now, s.RecordedAt)
	}
	if len(s.Triggers) != 0 {
		t.Errorf("expected no triggers, got %d", len(s.Triggers))
	}
	if s.CreatedAt.IsZero() || s.UpdatedAt.IsZero() {
		t.Error("timestamps should be set")
	}
}
