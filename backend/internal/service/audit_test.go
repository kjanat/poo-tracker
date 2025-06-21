package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/audit"
)

func TestAuditService_LogAction(t *testing.T) {
	svc := NewAuditService()
	ctx := context.Background()

	oldData := map[string]string{"name": "old"}
	newData := map[string]string{"name": "new"}

	svc.LogAction(ctx, "user1", "Bowel", "1", audit.AuditUpdate, oldData, newData)

	if len(svc.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(svc.logs))
	}
	log := svc.logs[0]
	if log.ID == "" {
		t.Error("expected ID to be set")
	}
	if log.UserID != "user1" || log.EntityType != "Bowel" || log.EntityID != "1" {
		t.Error("log fields not set correctly")
	}
	var od, nd map[string]string
	if err := json.Unmarshal([]byte(log.OldData), &od); err != nil {
		t.Fatalf("failed to unmarshal old data: %v", err)
	}
	if err := json.Unmarshal([]byte(log.NewData), &nd); err != nil {
		t.Fatalf("failed to unmarshal new data: %v", err)
	}
	if od["name"] != "old" || nd["name"] != "new" {
		t.Error("old/new data not stored correctly")
	}
	if log.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestAuditService_GetAuditLogs_OrderAndPagination(t *testing.T) {
	svc := NewAuditService()
	ctx := context.Background()

	svc.LogAction(ctx, "user1", "type", "1", audit.AuditCreate, nil, nil)
	svc.LogAction(ctx, "user1", "type", "2", audit.AuditCreate, nil, nil)
	svc.LogAction(ctx, "user1", "type", "3", audit.AuditCreate, nil, nil)

	now := time.Now()
	svc.logs[0].CreatedAt = now.Add(-3 * time.Hour)
	svc.logs[1].CreatedAt = now.Add(-2 * time.Hour)
	svc.logs[2].CreatedAt = now.Add(-1 * time.Hour)

	logs, err := svc.GetAuditLogs(ctx, "user1", 10, 0)
	if err != nil {
		t.Fatalf("GetAuditLogs failed: %v", err)
	}
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(logs))
	}
	if logs[0].EntityID != "3" || logs[1].EntityID != "2" || logs[2].EntityID != "1" {
		t.Error("logs not returned in descending order")
	}

	logs, err = svc.GetAuditLogs(ctx, "user1", 1, 1)
	if err != nil {
		t.Fatalf("GetAuditLogs with pagination failed: %v", err)
	}
	if len(logs) != 1 || logs[0].EntityID != "2" {
		t.Error("pagination did not return expected log")
	}
}

func TestAuditService_GetAuditLogsByEntityType(t *testing.T) {
	svc := NewAuditService()
	ctx := context.Background()

	svc.LogAction(ctx, "user1", "TypeA", "1", audit.AuditCreate, nil, nil)
	svc.LogAction(ctx, "user1", "TypeB", "2", audit.AuditCreate, nil, nil)
	svc.LogAction(ctx, "user1", "TypeA", "3", audit.AuditCreate, nil, nil)

	logs, err := svc.GetAuditLogsByEntityType(ctx, "user1", "TypeA")
	if err != nil {
		t.Fatalf("GetAuditLogsByEntityType failed: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
	ids := map[string]bool{"1": false, "3": false}
	for _, l := range logs {
		ids[l.EntityID] = true
	}
	if !ids["1"] || !ids["3"] {
		t.Error("logs for entity type not returned correctly")
	}
}

func TestAuditService_GetAuditLogsByEntity(t *testing.T) {
	svc := NewAuditService()
	ctx := context.Background()

	svc.LogAction(ctx, "user1", "TypeA", "1", audit.AuditCreate, nil, nil)
	svc.LogAction(ctx, "user1", "TypeA", "2", audit.AuditCreate, nil, nil)

	logs, err := svc.GetAuditLogsByEntity(ctx, "user1", "TypeA", "1")
	if err != nil {
		t.Fatalf("GetAuditLogsByEntity failed: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].EntityID != "1" {
		t.Errorf("expected entity ID 1, got %s", logs[0].EntityID)
	}
}

func TestAuditService_CleanupOldLogs(t *testing.T) {
	svc := NewAuditService()
	ctx := context.Background()

	svc.LogAction(ctx, "user1", "Type", "old", audit.AuditCreate, nil, nil)
	svc.LogAction(ctx, "user1", "Type", "new", audit.AuditCreate, nil, nil)

	now := time.Now()
	svc.logs[0].CreatedAt = now.Add(-3 * time.Hour)
	svc.logs[1].CreatedAt = now.Add(-30 * time.Minute)

	if err := svc.CleanupOldLogs(ctx, 2*time.Hour); err != nil {
		t.Fatalf("CleanupOldLogs failed: %v", err)
	}
	if len(svc.logs) != 1 {
		t.Fatalf("expected 1 log after cleanup, got %d", len(svc.logs))
	}
	if svc.logs[0].EntityID != "new" {
		t.Errorf("expected remaining log to have entityID 'new', got %s", svc.logs[0].EntityID)
	}
}
