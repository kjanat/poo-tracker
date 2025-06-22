package bowelmovement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BowelMovementDB represents a bowel movement record stored via GORM.
type BowelMovementDB struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID       string         `gorm:"type:varchar(64);not null;index" json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	RecordedAt   time.Time      `gorm:"not null;index" json:"recorded_at"`
	BristolType  int            `gorm:"not null;check:bristol_type >= 1 AND bristol_type <= 7" json:"bristol_type"` // 1-7 Bristol stool scale
	Volume       *float64       `json:"volume,omitempty"`                                                           // Optional volume in grams
	Color        *string        `json:"color,omitempty"`
	Consistency  *string        `json:"consistency,omitempty"`
	Floaters     bool           `json:"floaters"`
	Pain         int            `json:"pain"`
	Strain       int            `json:"strain"`
	Satisfaction int            `json:"satisfaction"`
	PhotoURL     *string        `json:"photo_url,omitempty"`
	SmellLevel   *string        `json:"smell_level,omitempty"`
	HasDetails   bool           `json:"has_details"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
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
