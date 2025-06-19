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

// BowelMovementUpdate represents fields that can be updated on a BowelMovement.
// Pointer fields allow distinguishing between "not provided" and "set to zero value".
type BowelMovementUpdate struct {
	BristolType *int    `json:"bristolType,omitempty"`
	Notes       *string `json:"notes,omitempty"`
}
