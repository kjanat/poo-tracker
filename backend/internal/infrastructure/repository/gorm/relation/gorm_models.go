package relation

import (
	"time"
)

// MealBowelMovementRelationGORM represents the GORM model for meal-bowel movement relations
type MealBowelMovementRelationGORM struct {
	ID              string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID          string    `gorm:"type:varchar(255);not null;index" json:"user_id"`
	MealID          string    `gorm:"type:varchar(255);not null;index" json:"meal_id"`
	BowelMovementID string    `gorm:"type:varchar(255);not null;index" json:"bowel_movement_id"`
	CreatedAt       time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null" json:"updated_at"`
	Strength        int       `gorm:"type:int;default:5" json:"strength"`
	Notes           string    `gorm:"type:text" json:"notes"`
	TimeGapHours    float64   `gorm:"type:decimal(10,2)" json:"time_gap_hours"`
	UserCorrelation *string   `gorm:"type:varchar(50)" json:"user_correlation"`
}

// TableName returns the table name for GORM
func (MealBowelMovementRelationGORM) TableName() string {
	return "meal_bowel_movement_relations"
}

// MealSymptomRelationGORM represents the GORM model for meal-symptom relations
type MealSymptomRelationGORM struct {
	ID              string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID          string    `gorm:"type:varchar(255);not null;index" json:"user_id"`
	MealID          string    `gorm:"type:varchar(255);not null;index" json:"meal_id"`
	SymptomID       string    `gorm:"type:varchar(255);not null;index" json:"symptom_id"`
	CreatedAt       time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null" json:"updated_at"`
	Strength        int       `gorm:"type:int;default:5" json:"strength"`
	Notes           string    `gorm:"type:text" json:"notes"`
	TimeGapHours    float64   `gorm:"type:decimal(10,2)" json:"time_gap_hours"`
	UserCorrelation *string   `gorm:"type:varchar(50)" json:"user_correlation"`
}

// TableName returns the table name for GORM
func (MealSymptomRelationGORM) TableName() string {
	return "meal_symptom_relations"
}
