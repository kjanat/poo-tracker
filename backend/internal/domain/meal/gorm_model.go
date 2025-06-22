package meal

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MealDB struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	MealTime    time.Time      `gorm:"not null;index" json:"meal_time"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Category    *string        `gorm:"type:varchar(50)" json:"category,omitempty"`
	Cuisine     string         `gorm:"type:varchar(50)" json:"cuisine,omitempty"`
	Calories    int            `json:"calories,omitempty"`
	SpicyLevel  *int           `json:"spicy_level,omitempty"`
	FiberRich   bool           `json:"fiber_rich"`
	Dairy       bool           `json:"dairy"`
	Gluten      bool           `json:"gluten"`
	PhotoURL    string         `gorm:"type:varchar(500)" json:"photo_url,omitempty"`
	Notes       string         `gorm:"type:text" json:"notes,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *MealDB) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (MealDB) TableName() string {
	return "meals"
}
