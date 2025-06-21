package repository

import (
	"context"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
)

type BowelMovementRepository interface {
	List(ctx context.Context) ([]bm.BowelMovement, error)
	Create(ctx context.Context, bm bm.BowelMovement) (bm.BowelMovement, error)
	Get(ctx context.Context, id string) (bm.BowelMovement, error)
	Update(ctx context.Context, id string, update bm.BowelMovementUpdate) (bm.BowelMovement, error)
	Delete(ctx context.Context, id string) error
}

type MealRepository interface {
	List(ctx context.Context) ([]meal.Meal, error)
	Create(ctx context.Context, m meal.Meal) (meal.Meal, error)
	Get(ctx context.Context, id string) (meal.Meal, error)
	Update(ctx context.Context, id string, update meal.MealUpdate) (meal.Meal, error)
	Delete(ctx context.Context, id string) error
}

type BowelMovementDetailsRepository interface {
	Create(ctx context.Context, details bm.BowelMovementDetails) (bm.BowelMovementDetails, error)
	Get(ctx context.Context, bowelMovementID string) (bm.BowelMovementDetails, error)
	Update(ctx context.Context, bowelMovementID string, update bm.BowelMovementDetailsUpdate) (bm.BowelMovementDetails, error)
	Delete(ctx context.Context, bowelMovementID string) error
	Exists(ctx context.Context, bowelMovementID string) (bool, error)
}
