package meal

import (
	"context"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"gorm.io/gorm"
)

type MealRepository struct {
	db *gorm.DB
}

func NewMealRepository(db *gorm.DB) meal.Repository {
	return &MealRepository{db: db}
}

// Implement all methods from meal.Repository interface
func (r *MealRepository) Create(ctx context.Context, m *meal.Meal) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MealRepository) GetByID(ctx context.Context, id string) (*meal.Meal, error) {
	var m meal.Meal
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MealRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*meal.Meal, error) {
	var ms []*meal.Meal
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&ms).Error
	return ms, err
}

func (r *MealRepository) Update(ctx context.Context, id string, update *meal.MealUpdate) error {
	return r.db.WithContext(ctx).Model(&meal.Meal{}).Where("id = ?", id).Updates(update).Error
}

func (r *MealRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&meal.Meal{}, "id = ?", id).Error
}

// Query
func (r *MealRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*meal.Meal, error) {
	var ms []*meal.Meal
	err := r.db.WithContext(ctx).Where("user_id = ? AND occurred_at BETWEEN ? AND ?", userID, start, end).Find(&ms).Error
	return ms, err
}

func (r *MealRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*meal.Meal, error) {
	var ms []*meal.Meal
	err := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&ms).Error
	return ms, err
}

func (r *MealRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&meal.Meal{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *MealRepository) GetLatestByUserID(ctx context.Context, userID string) (*meal.Meal, error) {
	var m meal.Meal
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("occurred_at desc").First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Analytics
func (r *MealRepository) GetNutritionSummary(ctx context.Context, userID string, start, end time.Time) (*meal.NutritionSummary, error) {
	// Implement analytics logic as needed
	return nil, nil
}

func (r *MealRepository) GetMealFrequency(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	// Implement analytics logic as needed
	return nil, nil
}
