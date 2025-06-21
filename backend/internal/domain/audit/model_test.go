package audit

import "testing"

func TestAllAuditActions(t *testing.T) {
	actions := AllAuditActions()
	if len(actions) != 3 {
		t.Fatalf("expected 3 actions, got %d", len(actions))
	}
	m := map[AuditAction]bool{}
	for _, a := range actions {
		m[a] = true
	}
	if !m[AuditCreate] || !m[AuditUpdate] || !m[AuditDelete] {
		t.Error("missing expected actions")
	}
}

func TestAuditActionIsValid(t *testing.T) {
	for _, a := range AllAuditActions() {
		if !a.IsValid() {
			t.Errorf("expected %s to be valid", a)
		}
	}
	if AuditAction("OTHER").IsValid() {
		t.Error("expected OTHER to be invalid")
	}
}

func TestNewAuditLog(t *testing.T) {
	log := NewAuditLog("u1", "User", "123", AuditCreate)
	if log.UserID != "u1" || log.EntityType != "User" || log.EntityID != "123" {
		t.Error("fields not set correctly")
	}
	if log.Action != AuditCreate {
		t.Errorf("expected action %s, got %s", AuditCreate, log.Action)
	}
	if log.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}
