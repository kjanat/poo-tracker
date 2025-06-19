package repository

import (
	"context"

	"github.com/kjanat/poo-tracker/backend/internal/model"
)

type BowelMovementRepository interface {
	List(ctx context.Context) ([]model.BowelMovement, error)
	Create(ctx context.Context, bm model.BowelMovement) (model.BowelMovement, error)
	Get(ctx context.Context, id string) (model.BowelMovement, error)
	Update(ctx context.Context, bm model.BowelMovement) (model.BowelMovement, error)
	Delete(ctx context.Context, id string) error
}
