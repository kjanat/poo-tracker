package medication

import (
	"context"
	"encoding/json"
	"time"

	domainmedication "github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"gorm.io/gorm"
)

type MedicationRepository struct {
	db *gorm.DB
}

func NewMedicationRepository(db *gorm.DB) domainmedication.Repository {
	return &MedicationRepository{db: db}
}

// Implement all methods from medication.Repository interface
func (r *MedicationRepository) Create(ctx context.Context, m *domainmedication.Medication) error {
	dbModel, err := ToMedicationDB(m)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(dbModel).Error
}

func (r *MedicationRepository) GetByID(ctx context.Context, id string) (*domainmedication.Medication, error) {
	var dbModel domainmedication.MedicationDB
	err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return ToMedication(&dbModel)
}

func (r *MedicationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domainmedication.Medication, error) {
	var dbModels []domainmedication.MedicationDB
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainmedication.Medication, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToMedication(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *MedicationRepository) Update(ctx context.Context, id string, update *domainmedication.MedicationUpdate) error {
	updates := map[string]interface{}{}
	if update.Name != nil {
		updates["name"] = *update.Name
	}
	if update.GenericName != nil {
		updates["generic_name"] = *update.GenericName
	}
	if update.Brand != nil {
		updates["brand"] = *update.Brand
	}
	if update.Category != nil {
		cat := string(*update.Category)
		updates["category"] = cat
	}
	if update.Dosage != nil {
		updates["dosage"] = *update.Dosage
	}
	if update.Form != nil {
		form := string(*update.Form)
		updates["form"] = form
	}
	if update.Frequency != nil {
		updates["frequency"] = *update.Frequency
	}
	if update.Route != nil {
		route := string(*update.Route)
		updates["route"] = route
	}
	if update.StartDate != nil {
		updates["start_date"] = *update.StartDate
	}
	if update.EndDate != nil {
		updates["end_date"] = *update.EndDate
	}
	if update.TakenAt != nil {
		updates["taken_at"] = *update.TakenAt
	}
	if update.Purpose != nil {
		updates["purpose"] = *update.Purpose
	}
	if update.SideEffects != nil {
		sideEffectsJSON, err := json.Marshal(update.SideEffects)
		if err != nil {
			return err
		}
		updates["side_effects"] = sideEffectsJSON
	}
	if update.Notes != nil {
		updates["notes"] = *update.Notes
	}
	if update.PhotoURL != nil {
		updates["photo_url"] = *update.PhotoURL
	}
	if update.IsActive != nil {
		updates["is_active"] = *update.IsActive
	}
	if update.IsAsNeeded != nil {
		updates["is_as_needed"] = *update.IsAsNeeded
	}
	return r.db.WithContext(ctx).Model(&domainmedication.MedicationDB{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MedicationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domainmedication.MedicationDB{}, "id = ?", id).Error
}

// Query
func (r *MedicationRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*domainmedication.Medication, error) {
	var dbModels []domainmedication.MedicationDB
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainmedication.Medication, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToMedication(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *MedicationRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*domainmedication.Medication, error) {
	var dbModels []domainmedication.MedicationDB
	err := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainmedication.Medication, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToMedication(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *MedicationRepository) GetAsNeededByUserID(ctx context.Context, userID string) ([]*domainmedication.Medication, error) {
	var dbModels []domainmedication.MedicationDB
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_as_needed = true", userID).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainmedication.Medication, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToMedication(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *MedicationRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*domainmedication.Medication, error) {
	var ms []*domainmedication.Medication
	err := r.db.WithContext(ctx).Where("user_id = ? AND occurred_at BETWEEN ? AND ?", userID, start, end).Find(&ms).Error
	return ms, err
}

func (r *MedicationRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainmedication.Medication{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Tracking
func (r *MedicationRepository) RecordDose(ctx context.Context, medicationID string, takenAt time.Time) error {
	// Implement dose tracking logic as needed
	return nil
}

func (r *MedicationRepository) GetDoseHistory(ctx context.Context, medicationID string, start, end time.Time) ([]*domainmedication.DoseRecord, error) {
	// Implement dose history logic as needed
	return nil, nil
}

// Analytics
func (r *MedicationRepository) GetUsageStats(ctx context.Context, userID string, start, end time.Time) (*domainmedication.UsageStats, error) {
	// Implement analytics logic as needed
	return nil, nil
}

func (r *MedicationRepository) GetCategoryBreakdown(ctx context.Context, userID string) (map[string]int, error) {
	// Implement analytics logic as needed
	return nil, nil
}
