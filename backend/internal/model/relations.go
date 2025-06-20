package model

import "time"

// MealBowelMovementRelation represents the relationship between meals and bowel movements
type MealBowelMovementRelation struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userId"`
	MealID          string    `json:"mealId"`
	BowelMovementID string    `json:"bowelMovementId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relationship strength (1-10) - how related the user thinks they are
	Strength int    `json:"strength"` // 1-10 scale
	Notes    string `json:"notes,omitempty"`

	// Time gap between meal and bowel movement (in hours)
	TimeGapHours float64 `json:"timeGapHours"`

	// User-reported correlation
	UserCorrelation *CorrelationType `json:"userCorrelation,omitempty"`
}

// MealSymptomRelation represents the relationship between meals and symptoms
type MealSymptomRelation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	MealID    string    `json:"mealId"`
	SymptomID string    `json:"symptomId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relationship strength (1-10) - how related the user thinks they are
	Strength int    `json:"strength"` // 1-10 scale
	Notes    string `json:"notes,omitempty"`

	// Time gap between meal and symptom (in hours)
	TimeGapHours float64 `json:"timeGapHours"`

	// User-reported correlation
	UserCorrelation *CorrelationType `json:"userCorrelation,omitempty"`
}

// NewMealBowelMovementRelation creates a new meal-bowel movement relation
func NewMealBowelMovementRelation(userID, mealID, bowelMovementID string, timeGapHours float64) MealBowelMovementRelation {
	now := time.Now()
	return MealBowelMovementRelation{
		UserID:          userID,
		MealID:          mealID,
		BowelMovementID: bowelMovementID,
		CreatedAt:       now,
		UpdatedAt:       now,
		TimeGapHours:    timeGapHours,
		Strength:        5, // Default neutral strength
	}
}

// NewMealSymptomRelation creates a new meal-symptom relation
func NewMealSymptomRelation(userID, mealID, symptomID string, timeGapHours float64) MealSymptomRelation {
	now := time.Now()
	return MealSymptomRelation{
		UserID:       userID,
		MealID:       mealID,
		SymptomID:    symptomID,
		CreatedAt:    now,
		UpdatedAt:    now,
		TimeGapHours: timeGapHours,
		Strength:     5, // Default neutral strength
	}
}

// AuditLog represents an audit log entry for tracking changes
type AuditLog struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	EntityType string      `json:"entityType"` // e.g., "bowel_movement", "meal", "symptom"
	EntityID   string      `json:"entityId"`
	Action     AuditAction `json:"action"`            // CREATE, UPDATE, DELETE
	OldData    string      `json:"oldData,omitempty"` // JSON representation of old data
	NewData    string      `json:"newData,omitempty"` // JSON representation of new data
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
