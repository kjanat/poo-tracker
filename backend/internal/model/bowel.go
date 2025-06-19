package model

import "time"

// BowelMovement represents a bowel movement entry.
type BowelMovement struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	BristolType int       `json:"bristolType"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
