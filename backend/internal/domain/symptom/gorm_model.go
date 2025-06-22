package symptom

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SymptomDB struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      string         `gorm:"type:varchar(64);not null;index" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	RecordedAt  time.Time      `gorm:"not null;index" json:"recorded_at"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Category    *string        `gorm:"type:varchar(50)" json:"category,omitempty"`
	Severity    int            `json:"severity"`
	Duration    *int           `json:"duration,omitempty"`
	BodyPart    *string        `gorm:"type:varchar(50)" json:"body_part,omitempty"`
	Type        *string        `gorm:"type:varchar(50)" json:"type,omitempty"`
	Triggers    datatypes.JSON `gorm:"type:json" json:"triggers,omitempty"`
	Notes       *string        `gorm:"type:text" json:"notes,omitempty"`
	PhotoURL    *string        `gorm:"type:varchar(500)" json:"photo_url,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *SymptomDB) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (SymptomDB) TableName() string {
	return "symptoms"
}
