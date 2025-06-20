package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/model"
)

// AuditService handles audit logging for the application
type AuditService struct {
	logs []model.AuditLog
}

// NewAuditService creates a new audit service
func NewAuditService() *AuditService {
	return &AuditService{
		logs: make([]model.AuditLog, 0),
	}
}

// LogAction logs an action to the audit trail
func (s *AuditService) LogAction(ctx context.Context, userID, entityType, entityID string, action model.AuditAction, oldData, newData interface{}) {
	auditLog := model.NewAuditLog(userID, entityType, entityID, action)
	auditLog.ID = uuid.New().String()

	if oldData != nil {
		if oldDataJSON, err := json.Marshal(oldData); err == nil {
			auditLog.OldData = string(oldDataJSON)
		}
	}

	if newData != nil {
		if newDataJSON, err := json.Marshal(newData); err == nil {
			auditLog.NewData = string(newDataJSON)
		}
	}

	s.logs = append(s.logs, auditLog)

	// Log to console for debugging (in production, this would go to a proper log aggregator)
	log.Printf("AUDIT: User %s performed %s on %s %s", userID, action, entityType, entityID)
}

// GetAuditLogs retrieves audit logs for a user
func (s *AuditService) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]model.AuditLog, error) {
	var userLogs []model.AuditLog

	for _, log := range s.logs {
		if log.UserID == userID {
			userLogs = append(userLogs, log)
		}
	}

	// Sort by creation time (newest first)
	for i := 0; i < len(userLogs)-1; i++ {
		for j := i + 1; j < len(userLogs); j++ {
			if userLogs[i].CreatedAt.Before(userLogs[j].CreatedAt) {
				userLogs[i], userLogs[j] = userLogs[j], userLogs[i]
			}
		}
	}

	// Apply pagination
	if offset >= len(userLogs) {
		return []model.AuditLog{}, nil
	}

	end := offset + limit
	if end > len(userLogs) {
		end = len(userLogs)
	}

	return userLogs[offset:end], nil
}

// GetAuditLogsByEntityType retrieves audit logs for a specific entity type
func (s *AuditService) GetAuditLogsByEntityType(ctx context.Context, userID, entityType string) ([]model.AuditLog, error) {
	var logs []model.AuditLog

	for _, log := range s.logs {
		if log.UserID == userID && log.EntityType == entityType {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

// GetAuditLogsByEntity retrieves audit logs for a specific entity
func (s *AuditService) GetAuditLogsByEntity(ctx context.Context, userID, entityType, entityID string) ([]model.AuditLog, error) {
	var logs []model.AuditLog

	for _, log := range s.logs {
		if log.UserID == userID && log.EntityType == entityType && log.EntityID == entityID {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

// CleanupOldLogs removes audit logs older than the specified duration
func (s *AuditService) CleanupOldLogs(ctx context.Context, maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)

	var filteredLogs []model.AuditLog
	for _, log := range s.logs {
		if log.CreatedAt.After(cutoff) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	s.logs = filteredLogs
	log.Printf("AUDIT: Cleaned up old audit logs, %d logs remain", len(s.logs))

	return nil
}
