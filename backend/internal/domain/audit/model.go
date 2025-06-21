package audit

import (
	"time"
)

// AuditAction represents the type of action performed for audit logging
//go:generate stringer -type=AuditAction
// Note: we keep simple string enums for now

type AuditAction string

const (
	AuditCreate AuditAction = "CREATE"
	AuditUpdate AuditAction = "UPDATE"
	AuditDelete AuditAction = "DELETE"
)

// AllAuditActions returns all valid AuditAction values
func AllAuditActions() []AuditAction {
	return []AuditAction{AuditCreate, AuditUpdate, AuditDelete}
}

// IsValid checks if the AuditAction value is valid
func (a AuditAction) IsValid() bool {
	for _, valid := range AllAuditActions() {
		if a == valid {
			return true
		}
	}
	return false
}

// AuditLog represents an audit log entry for tracking changes
// used by the in-memory audit service.
type AuditLog struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	EntityType string      `json:"entityType"`
	EntityID   string      `json:"entityId"`
	Action     AuditAction `json:"action"`
	OldData    string      `json:"oldData,omitempty"`
	NewData    string      `json:"newData,omitempty"`
	IPAddress  string      `json:"ipAddress,omitempty"`
	UserAgent  string      `json:"userAgent,omitempty"`
	CreatedAt  time.Time   `json:"createdAt"`
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(userID, entityType, entityID string, action AuditAction) AuditLog {
	return AuditLog{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		CreatedAt:  time.Now(),
	}
}
