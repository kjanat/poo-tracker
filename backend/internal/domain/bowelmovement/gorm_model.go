package bowelmovement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BowelMovementDB represents a bowel movement record stored via GORM.
type BowelMovementDB struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Timestamp   time.Time      `gorm:"not null;index" json:"timestamp"`
	BristolType int            `gorm:"not null" json:"bristol_type"` // 1-7 Bristol stool scale
	Volume      *float64       `json:"volume,omitempty"`             // Optional volume in grams
	Color       string         `json:"color,omitempty"`
	Urgency     *int           `json:"urgency,omitempty"`    // 1-5 scale
	Difficulty  *int           `json:"difficulty,omitempty"` // 1-5 scale
	Pain        *int           `json:"pain,omitempty"`       // 1-5 scale
	Notes       string         `json:"notes,omitempty"`
	PhotoURL    string         `json:"photo_url,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate GORM hook to set UUID if not provided
func (b *BowelMovementDB) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for BowelMovement
func (BowelMovementDB) TableName() string {
	return "bowel_movements"
}
