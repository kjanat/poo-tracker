package repository

import (
	"context"

	"github.com/kjanat/poo-tracker/backend/internal/model"
)

type BowelMovementRepository interface {
	List(ctx context.Context) ([]model.BowelMovement, error)
	Create(ctx context.Context, bm model.BowelMovement) (model.BowelMovement, error)
	Get(ctx context.Context, id string) (model.BowelMovement, error)
	Update(ctx context.Context, id string, update model.BowelMovementUpdate) (model.BowelMovement, error)
	Delete(ctx context.Context, id string) error
}

type MealRepository interface {
	List(ctx context.Context) ([]model.Meal, error)
	Create(ctx context.Context, m model.Meal) (model.Meal, error)
	Get(ctx context.Context, id string) (model.Meal, error)
	Update(ctx context.Context, id string, update model.MealUpdate) (model.Meal, error)
	Delete(ctx context.Context, id string) error
}

type BowelMovementDetailsRepository interface {
	Create(ctx context.Context, details model.BowelMovementDetails) (model.BowelMovementDetails, error)
	Get(ctx context.Context, bowelMovementID string) (model.BowelMovementDetails, error)
	Update(ctx context.Context, bowelMovementID string, update model.BowelMovementDetailsUpdate) (model.BowelMovementDetails, error)
	Delete(ctx context.Context, bowelMovementID string) error
	Exists(ctx context.Context, bowelMovementID string) (bool, error)
}
