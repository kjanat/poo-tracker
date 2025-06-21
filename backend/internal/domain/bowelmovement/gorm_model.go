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
	BristolType int            `gorm:"not null;check:bristol_type >= 1 AND bristol_type <= 7" json:"bristol_type"` // 1-7 Bristol stool scale
	Volume      *float64       `gorm:"check:volume > 0" json:"volume,omitempty"`                                   // Optional volume in grams
	Color       *string        `gorm:"type:varchar(50)" json:"color,omitempty"`
	Urgency     *int           `gorm:"check:urgency >= 1 AND urgency <= 5" json:"urgency,omitempty"`               // 1-5 scale
	Difficulty  *int           `gorm:"check:difficulty >= 1 AND difficulty <= 5" json:"difficulty,omitempty"`      // 1-5 scale
	Pain        *int           `gorm:"check:pain >= 1 AND pain <= 5" json:"pain,omitempty"`                        // 1-5 scale
	Notes       *string        `gorm:"type:text" json:"notes,omitempty"`
	PhotoURL    *string        `gorm:"type:varchar(500)" json:"photo_url,omitempty"`
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
