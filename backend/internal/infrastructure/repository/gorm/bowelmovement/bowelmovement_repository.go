package bowelmovement

import (
	"context"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"gorm.io/gorm"
)

type BowelMovementRepository struct {
	db *gorm.DB
}

func NewBowelMovementRepository(db *gorm.DB) bowelmovement.Repository {
	return &BowelMovementRepository{db: db}
}

// Implement all methods from bowelmovement.Repository interface
func (r *BowelMovementRepository) Create(ctx context.Context, bm *bowelmovement.BowelMovement) error {
	return r.db.WithContext(ctx).Create(bm).Error
}

func (r *BowelMovementRepository) GetByID(ctx context.Context, id string) (*bowelmovement.BowelMovement, error) {
	var bm bowelmovement.BowelMovement
	err := r.db.WithContext(ctx).First(&bm, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bm, nil
}

func (r *BowelMovementRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*bowelmovement.BowelMovement, error) {
	var bms []*bowelmovement.BowelMovement
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&bms).Error
	return bms, err
}

func (r *BowelMovementRepository) Update(ctx context.Context, id string, update *bowelmovement.BowelMovementUpdate) error {
	return r.db.WithContext(ctx).Model(&bowelmovement.BowelMovement{}).Where("id = ?", id).Updates(update).Error
}

func (r *BowelMovementRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&bowelmovement.BowelMovement{}, "id = ?", id).Error
}

// Details
func (r *BowelMovementRepository) CreateDetails(ctx context.Context, details *bowelmovement.BowelMovementDetails) error {
	return r.db.WithContext(ctx).Create(details).Error
}

func (r *BowelMovementRepository) GetDetailsByBowelMovementID(ctx context.Context, bowelMovementID string) (*bowelmovement.BowelMovementDetails, error) {
	var details bowelmovement.BowelMovementDetails
	err := r.db.WithContext(ctx).First(&details, "bowel_movement_id = ?", bowelMovementID).Error
	if err != nil {
		return nil, err
	}
	return &details, nil
}

func (r *BowelMovementRepository) UpdateDetails(ctx context.Context, bowelMovementID string, update *bowelmovement.BowelMovementDetailsUpdate) error {
	return r.db.WithContext(ctx).Model(&bowelmovement.BowelMovementDetails{}).Where("bowel_movement_id = ?", bowelMovementID).Updates(update).Error
}

func (r *BowelMovementRepository) DeleteDetails(ctx context.Context, bowelMovementID string) error {
	return r.db.WithContext(ctx).Delete(&bowelmovement.BowelMovementDetails{}, "bowel_movement_id = ?", bowelMovementID).Error
}

// Query
func (r *BowelMovementRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*bowelmovement.BowelMovement, error) {
	var bms []*bowelmovement.BowelMovement
	err := r.db.WithContext(ctx).Where("user_id = ? AND occurred_at BETWEEN ? AND ?", userID, start, end).Find(&bms).Error
	return bms, err
}

func (r *BowelMovementRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&bowelmovement.BowelMovement{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *BowelMovementRepository) GetLatestByUserID(ctx context.Context, userID string) (*bowelmovement.BowelMovement, error) {
	var bm bowelmovement.BowelMovement
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("occurred_at desc").First(&bm).Error
	if err != nil {
		return nil, err
	}
	return &bm, nil
}

// Analytics
func (r *BowelMovementRepository) GetAveragesByUserID(ctx context.Context, userID string, start, end time.Time) (*bowelmovement.BowelMovementAverages, error) {
	// Implement analytics logic as needed
	return nil, nil
}

func (r *BowelMovementRepository) GetFrequencyByUserID(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	// Implement analytics logic as needed
	return nil, nil
}
