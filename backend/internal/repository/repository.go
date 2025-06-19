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
	ListMeals(ctx context.Context) ([]model.Meal, error)
	CreateMeal(ctx context.Context, m model.Meal) (model.Meal, error)
	GetMeal(ctx context.Context, id string) (model.Meal, error)
	UpdateMeal(ctx context.Context, id string, update model.MealUpdate) (model.Meal, error)
	DeleteMeal(ctx context.Context, id string) error
}
