package medication

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicationDB struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	GenericName string         `gorm:"type:varchar(100)" json:"generic_name,omitempty"`
	Brand       string         `gorm:"type:varchar(100)" json:"brand,omitempty"`
	Category    *string        `gorm:"type:varchar(50)" json:"category,omitempty"`
	Dosage      string         `gorm:"type:varchar(50)" json:"dosage"`
	Form        *string        `gorm:"type:varchar(50)" json:"form,omitempty"`
	Frequency   string         `gorm:"type:varchar(50)" json:"frequency"`
	Route       *string        `gorm:"type:varchar(50)" json:"route,omitempty"`
	StartDate   *time.Time     `json:"start_date,omitempty"`
	EndDate     *time.Time     `json:"end_date,omitempty"`
	TakenAt     *time.Time     `json:"taken_at,omitempty"`
	Purpose     string         `gorm:"type:text" json:"purpose,omitempty"`
	SideEffects string         `gorm:"type:text" json:"side_effects,omitempty"` // JSON-encoded array
	Notes       string         `gorm:"type:text" json:"notes,omitempty"`
	PhotoURL    string         `gorm:"type:varchar(500)" json:"photo_url,omitempty"`
	IsActive    bool           `json:"is_active"`
	IsAsNeeded  bool           `json:"is_as_needed"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *MedicationDB) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (MedicationDB) TableName() string {
	return "medications"
}
