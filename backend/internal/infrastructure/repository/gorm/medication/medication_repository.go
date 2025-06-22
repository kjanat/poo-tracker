package medication

import (
	"context"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"gorm.io/gorm"
)

type MedicationRepository struct {
	db *gorm.DB
}

func NewMedicationRepository(db *gorm.DB) medication.Repository {
	return &MedicationRepository{db: db}
}

// Implement all methods from medication.Repository interface
func (r *MedicationRepository) Create(ctx context.Context, m *medication.Medication) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MedicationRepository) GetByID(ctx context.Context, id string) (*medication.Medication, error) {
	var m medication.Medication
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MedicationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*medication.Medication, error) {
	var ms []*medication.Medication
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&ms).Error
	return ms, err
}

func (r *MedicationRepository) Update(ctx context.Context, id string, update *medication.MedicationUpdate) error {
	return r.db.WithContext(ctx).Model(&medication.Medication{}).Where("id = ?", id).Updates(update).Error
}

func (r *MedicationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&medication.Medication{}, "id = ?", id).Error
}

// Query
func (r *MedicationRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*medication.Medication, error) {
	var ms []*medication.Medication
	err := r.db.WithContext(ctx).Where("user_id = ? AND active = true", userID).Find(&ms).Error
	return ms, err
}

func (r *MedicationRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*medication.Medication, error) {
	var ms []*medication.Medication
	err := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&ms).Error
	return ms, err
}

func (r *MedicationRepository) GetAsNeededByUserID(ctx context.Context, userID string) ([]*medication.Medication, error) {
	var ms []*medication.Medication
	err := r.db.WithContext(ctx).Where("user_id = ? AND as_needed = true", userID).Find(&ms).Error
	return ms, err
}

func (r *MedicationRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*medication.Medication, error) {
	var ms []*medication.Medication
	err := r.db.WithContext(ctx).Where("user_id = ? AND occurred_at BETWEEN ? AND ?", userID, start, end).Find(&ms).Error
	return ms, err
}

func (r *MedicationRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&medication.Medication{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Tracking
func (r *MedicationRepository) RecordDose(ctx context.Context, medicationID string, takenAt time.Time) error {
	// Implement dose tracking logic as needed
	return nil
}

func (r *MedicationRepository) GetDoseHistory(ctx context.Context, medicationID string, start, end time.Time) ([]*medication.DoseRecord, error) {
	// Implement dose history logic as needed
	return nil, nil
}

// Analytics
func (r *MedicationRepository) GetUsageStats(ctx context.Context, userID string, start, end time.Time) (*medication.UsageStats, error) {
	// Implement analytics logic as needed
	return nil, nil
}

func (r *MedicationRepository) GetCategoryBreakdown(ctx context.Context, userID string) (map[string]int, error) {
	// Implement analytics logic as needed
	return nil, nil
}
