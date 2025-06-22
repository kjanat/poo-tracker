package symptom

import (
	"context"
	"encoding/json"
	"time"

	domainsymptom "github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"gorm.io/gorm"
)

type SymptomRepository struct {
	db *gorm.DB
}

func NewSymptomRepository(db *gorm.DB) domainsymptom.Repository {
	return &SymptomRepository{db: db}
}

// Implement all methods from symptom.Repository interface
func (r *SymptomRepository) Create(ctx context.Context, s *domainsymptom.Symptom) error {
	dbModel, err := ToSymptomDB(s)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(dbModel).Error
}

func (r *SymptomRepository) GetByID(ctx context.Context, id string) (*domainsymptom.Symptom, error) {
	var dbModel domainsymptom.SymptomDB
	err := r.db.WithContext(ctx).First(&dbModel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return ToSymptom(&dbModel)
}

func (r *SymptomRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domainsymptom.Symptom, error) {
	var dbModels []domainsymptom.SymptomDB
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainsymptom.Symptom, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToSymptom(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *SymptomRepository) Update(ctx context.Context, id string, update *domainsymptom.SymptomUpdate) error {
	// For Triggers, marshal []string to JSON if present
	updates := map[string]interface{}{}
	if update.Name != nil {
		updates["name"] = *update.Name
	}
	if update.Description != nil {
		updates["description"] = *update.Description
	}
	if update.RecordedAt != nil {
		updates["recorded_at"] = *update.RecordedAt
	}
	if update.Category != nil {
		cat := string(*update.Category)
		updates["category"] = cat
	}
	if update.Severity != nil {
		updates["severity"] = *update.Severity
	}
	if update.Duration != nil {
		updates["duration"] = *update.Duration
	}
	if update.BodyPart != nil {
		updates["body_part"] = *update.BodyPart
	}
	if update.Type != nil {
		typeStr := string(*update.Type)
		updates["type"] = typeStr
	}
	if update.Triggers != nil {
		triggersJSON, err := json.Marshal(update.Triggers)
		if err != nil {
			return err
		}
		updates["triggers"] = triggersJSON
	}
	if update.Notes != nil {
		updates["notes"] = *update.Notes
	}
	if update.PhotoURL != nil {
		updates["photo_url"] = *update.PhotoURL
	}
	return r.db.WithContext(ctx).Model(&domainsymptom.SymptomDB{}).Where("id = ?", id).Updates(updates).Error
}

func (r *SymptomRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domainsymptom.SymptomDB{}, "id = ?", id).Error
}

// Query
func (r *SymptomRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*domainsymptom.Symptom, error) {
	var dbModels []domainsymptom.SymptomDB
	err := r.db.WithContext(ctx).Where("user_id = ? AND recorded_at BETWEEN ? AND ?", userID, start, end).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainsymptom.Symptom, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToSymptom(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *SymptomRepository) GetByCategory(ctx context.Context, userID string, category string) ([]*domainsymptom.Symptom, error) {
	var dbModels []domainsymptom.SymptomDB
	err := r.db.WithContext(ctx).Where("user_id = ? AND category = ?", userID, category).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainsymptom.Symptom, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToSymptom(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *SymptomRepository) GetBySeverity(ctx context.Context, userID string, minSeverity, maxSeverity int) ([]*domainsymptom.Symptom, error) {
	var dbModels []domainsymptom.SymptomDB
	err := r.db.WithContext(ctx).Where("user_id = ? AND severity BETWEEN ? AND ?", userID, minSeverity, maxSeverity).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domainsymptom.Symptom, 0, len(dbModels))
	for _, dbModel := range dbModels {
		domainModel, err := ToSymptom(&dbModel)
		if err != nil {
			return nil, err
		}
		result = append(result, domainModel)
	}
	return result, nil
}

func (r *SymptomRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domainsymptom.SymptomDB{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *SymptomRepository) GetLatestByUserID(ctx context.Context, userID string) (*domainsymptom.Symptom, error) {
	var dbModel domainsymptom.SymptomDB
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("recorded_at desc").First(&dbModel).Error
	if err != nil {
		return nil, err
	}
	return ToSymptom(&dbModel)
}

// Analytics
func (r *SymptomRepository) GetSeverityStats(ctx context.Context, userID string, start, end time.Time) (*domainsymptom.SeverityStats, error) {
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
