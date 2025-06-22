package relation

import (
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/relations"
	"gorm.io/gorm"
)

// MealBowelMovementRelationRepository implements repository for meal-bowel movement relations
type MealBowelMovementRelationRepository struct {
	db *gorm.DB
}

// NewMealBowelMovementRelationRepository creates a new repository instance
func NewMealBowelMovementRelationRepository(db *gorm.DB) *MealBowelMovementRelationRepository {
	return &MealBowelMovementRelationRepository{db: db}
}

// Create creates a new meal-bowel movement relation
func (r *MealBowelMovementRelationRepository) Create(relation *relations.MealBowelMovementRelation) error {
	if relation.ID == "" {
		relation.ID = uuid.New().String()
	}

	gormModel := r.toGORM(relation)
	return r.db.Create(&gormModel).Error
}

// GetByID retrieves a relation by ID
func (r *MealBowelMovementRelationRepository) GetByID(id string) (*relations.MealBowelMovementRelation, error) {
	var gormModel MealBowelMovementRelationGORM
	err := r.db.First(&gormModel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return r.toDomain(&gormModel), nil
}

// GetByUserID retrieves all relations for a user
func (r *MealBowelMovementRelationRepository) GetByUserID(userID string) ([]*relations.MealBowelMovementRelation, error) {
	var gormModels []MealBowelMovementRelationGORM
	err := r.db.Where("user_id = ?", userID).Find(&gormModels).Error
	if err != nil {
		return nil, err
	}

	result := make([]*relations.MealBowelMovementRelation, len(gormModels))
	for i, model := range gormModels {
		result[i] = r.toDomain(&model)
	}
	return result, nil
}

// GetByMealID retrieves all relations for a meal
func (r *MealBowelMovementRelationRepository) GetByMealID(mealID string) ([]*relations.MealBowelMovementRelation, error) {
	var gormModels []MealBowelMovementRelationGORM
	err := r.db.Where("meal_id = ?", mealID).Find(&gormModels).Error
	if err != nil {
		return nil, err
	}

	result := make([]*relations.MealBowelMovementRelation, len(gormModels))
	for i, model := range gormModels {
		result[i] = r.toDomain(&model)
	}
	return result, nil
}

// Update updates an existing relation
func (r *MealBowelMovementRelationRepository) Update(relation *relations.MealBowelMovementRelation) error {
	relation.UpdatedAt = time.Now()
	gormModel := r.toGORM(relation)
	return r.db.Save(&gormModel).Error
}

// Delete deletes a relation by ID
func (r *MealBowelMovementRelationRepository) Delete(id string) error {
	return r.db.Delete(&MealBowelMovementRelationGORM{}, "id = ?", id).Error
}

// Helper methods for conversion
func (r *MealBowelMovementRelationRepository) toGORM(domain *relations.MealBowelMovementRelation) MealBowelMovementRelationGORM {
	var userCorrelation *string
	if domain.UserCorrelation != nil {
		str := string(*domain.UserCorrelation)
		userCorrelation = &str
	}

	return MealBowelMovementRelationGORM{
		ID:              domain.ID,
		UserID:          domain.UserID,
		MealID:          domain.MealID,
		BowelMovementID: domain.BowelMovementID,
		CreatedAt:       domain.CreatedAt,
		UpdatedAt:       domain.UpdatedAt,
		Strength:        domain.Strength,
		Notes:           domain.Notes,
		TimeGapHours:    domain.TimeGapHours,
		UserCorrelation: userCorrelation,
	}
}

func (r *MealBowelMovementRelationRepository) toDomain(gorm *MealBowelMovementRelationGORM) *relations.MealBowelMovementRelation {
	var userCorrelation *relations.CorrelationType
	if gorm.UserCorrelation != nil {
		ct := relations.CorrelationType(*gorm.UserCorrelation)
		userCorrelation = &ct
	}

	return &relations.MealBowelMovementRelation{
		ID:              gorm.ID,
		UserID:          gorm.UserID,
		MealID:          gorm.MealID,
		BowelMovementID: gorm.BowelMovementID,
		CreatedAt:       gorm.CreatedAt,
		UpdatedAt:       gorm.UpdatedAt,
		Strength:        gorm.Strength,
		Notes:           gorm.Notes,
		TimeGapHours:    gorm.TimeGapHours,
		UserCorrelation: userCorrelation,
	}
}

// MealSymptomRelationRepository implements repository for meal-symptom relations
type MealSymptomRelationRepository struct {
	db *gorm.DB
}

// NewMealSymptomRelationRepository creates a new repository instance
func NewMealSymptomRelationRepository(db *gorm.DB) *MealSymptomRelationRepository {
	return &MealSymptomRelationRepository{db: db}
}

// Create creates a new meal-symptom relation
func (r *MealSymptomRelationRepository) Create(relation *relations.MealSymptomRelation) error {
	if relation.ID == "" {
		relation.ID = uuid.New().String()
	}

	gormModel := r.toGORM(relation)
	return r.db.Create(&gormModel).Error
}

// GetByID retrieves a relation by ID
func (r *MealSymptomRelationRepository) GetByID(id string) (*relations.MealSymptomRelation, error) {
	var gormModel MealSymptomRelationGORM
	err := r.db.First(&gormModel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return r.toDomain(&gormModel), nil
}

// GetByUserID retrieves all relations for a user
func (r *MealSymptomRelationRepository) GetByUserID(userID string) ([]*relations.MealSymptomRelation, error) {
	var gormModels []MealSymptomRelationGORM
	err := r.db.Where("user_id = ?", userID).Find(&gormModels).Error
	if err != nil {
		return nil, err
	}

	result := make([]*relations.MealSymptomRelation, len(gormModels))
	for i, model := range gormModels {
		result[i] = r.toDomain(&model)
	}
	return result, nil
}

// GetByMealID retrieves all relations for a meal
func (r *MealSymptomRelationRepository) GetByMealID(mealID string) ([]*relations.MealSymptomRelation, error) {
	var gormModels []MealSymptomRelationGORM
	err := r.db.Where("meal_id = ?", mealID).Find(&gormModels).Error
	if err != nil {
		return nil, err
	}

	result := make([]*relations.MealSymptomRelation, len(gormModels))
	for i, model := range gormModels {
		result[i] = r.toDomain(&model)
	}
	return result, nil
}

// Update updates an existing relation
func (r *MealSymptomRelationRepository) Update(relation *relations.MealSymptomRelation) error {
	relation.UpdatedAt = time.Now()
	gormModel := r.toGORM(relation)
	return r.db.Save(&gormModel).Error
}

// Delete deletes a relation by ID
func (r *MealSymptomRelationRepository) Delete(id string) error {
	return r.db.Delete(&MealSymptomRelationGORM{}, "id = ?", id).Error
}

// Helper methods for conversion
func (r *MealSymptomRelationRepository) toGORM(domain *relations.MealSymptomRelation) MealSymptomRelationGORM {
	var userCorrelation *string
	if domain.UserCorrelation != nil {
		str := string(*domain.UserCorrelation)
		userCorrelation = &str
	}

	return MealSymptomRelationGORM{
		ID:              domain.ID,
		UserID:          domain.UserID,
		MealID:          domain.MealID,
		SymptomID:       domain.SymptomID,
		CreatedAt:       domain.CreatedAt,
		UpdatedAt:       domain.UpdatedAt,
		Strength:        domain.Strength,
		Notes:           domain.Notes,
		TimeGapHours:    domain.TimeGapHours,
		UserCorrelation: userCorrelation,
	}
}

func (r *MealSymptomRelationRepository) toDomain(gorm *MealSymptomRelationGORM) *relations.MealSymptomRelation {
	var userCorrelation *relations.CorrelationType
	if gorm.UserCorrelation != nil {
		ct := relations.CorrelationType(*gorm.UserCorrelation)
		userCorrelation = &ct
	}

	return &relations.MealSymptomRelation{
		ID:              gorm.ID,
		UserID:          gorm.UserID,
		MealID:          gorm.MealID,
		SymptomID:       gorm.SymptomID,
		CreatedAt:       gorm.CreatedAt,
		UpdatedAt:       gorm.UpdatedAt,
		Strength:        gorm.Strength,
		Notes:           gorm.Notes,
		TimeGapHours:    gorm.TimeGapHours,
		UserCorrelation: userCorrelation,
	}
}
