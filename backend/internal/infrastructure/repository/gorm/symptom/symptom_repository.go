package symptom

import (
	"context"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"gorm.io/gorm"
)

type SymptomRepository struct {
	db *gorm.DB
}

func NewSymptomRepository(db *gorm.DB) symptom.Repository {
	return &SymptomRepository{db: db}
}

// Implement all methods from symptom.Repository interface
func (r *SymptomRepository) Create(ctx context.Context, s *symptom.Symptom) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *SymptomRepository) GetByID(ctx context.Context, id string) (*symptom.Symptom, error) {
	var s symptom.Symptom
	err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SymptomRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*symptom.Symptom, error) {
	var ss []*symptom.Symptom
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&ss).Error
	return ss, err
}

func (r *SymptomRepository) Update(ctx context.Context, id string, update *symptom.SymptomUpdate) error {
	return r.db.WithContext(ctx).Model(&symptom.Symptom{}).Where("id = ?", id).Updates(update).Error
}

func (r *SymptomRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&symptom.Symptom{}, "id = ?", id).Error
}

// Query
func (r *SymptomRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*symptom.Symptom, error) {
	var ss []*symptom.Symptom
	err := r.db.WithContext(ctx).Where("user_id = ? AND occurred_at BETWEEN ? AND ?", userID, start, end).Find(&ss).Error
	return ss, err
}

func (r *SymptomRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*symptom.Symptom, error) {
	var ss []*symptom.Symptom
	err := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&ss).Error
	return ss, err
}

func (r *SymptomRepository) GetBySeverity(ctx context.Context, userID string, minSeverity, maxSeverity int) ([]*symptom.Symptom, error) {
	var ss []*symptom.Symptom
	err := r.db.WithContext(ctx).Where("user_id = ? AND severity BETWEEN ? AND ?", userID, minSeverity, maxSeverity).Find(&ss).Error
	return ss, err
}

func (r *SymptomRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&symptom.Symptom{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *SymptomRepository) GetLatestByUserID(ctx context.Context, userID string) (*symptom.Symptom, error) {
	var s symptom.Symptom
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("occurred_at desc").First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Analytics
func (r *SymptomRepository) GetSeverityStats(ctx context.Context, userID string, start, end time.Time) (*symptom.SeverityStats, error) {
	// Implement analytics logic as needed
	return nil, nil
}

func (r *SymptomRepository) GetCategoryFrequency(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	// Implement analytics logic as needed
	return nil, nil
}

func (r *SymptomRepository) GetTriggerAnalysis(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	// Implement analytics logic as needed
	return nil, nil
}
