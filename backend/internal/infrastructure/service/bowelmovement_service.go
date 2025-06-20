package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// BowelMovementService implements the bowel movement business logic
type BowelMovementService struct {
	repo bowelmovement.Repository
}

// NewBowelMovementService creates a new bowel movement service
func NewBowelMovementService(repo bowelmovement.Repository) bowelmovement.Service {
	return &BowelMovementService{
		repo: repo,
	}
}

// Create creates a new bowel movement with business validation
func (s *BowelMovementService) Create(ctx context.Context, userID string, input *bowelmovement.CreateBowelMovementInput) (*bowelmovement.BowelMovement, error) {
	// Validate user ID
	if userID == "" {
		return nil, bowelmovement.ErrInvalidUserID
	}

	// Validate input
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Set defaults
	if input.RecordedAt.IsZero() {
		input.RecordedAt = time.Now()
	}

	// Convert string pointers to shared types
	var volume *shared.Volume
	if input.Volume != nil {
		v := shared.Volume(*input.Volume)
		if v.IsValid() {
			volume = &v
		}
	}

	var color *shared.Color
	if input.Color != nil {
		c := shared.Color(*input.Color)
		if c.IsValid() {
			color = &c
		}
	}

	var consistency *shared.Consistency
	if input.Consistency != nil {
		cons := shared.Consistency(*input.Consistency)
		if cons.IsValid() {
			consistency = &cons
		}
	}

	var smellLevel *shared.SmellLevel
	if input.SmellLevel != nil {
		smell := shared.SmellLevel(*input.SmellLevel)
		if smell.IsValid() {
			smellLevel = &smell
		}
	}

	// Create bowel movement
	bm := &bowelmovement.BowelMovement{
		ID:           uuid.New().String(),
		UserID:       userID,
		BristolType:  input.BristolType,
		RecordedAt:   input.RecordedAt,
		Volume:       volume,
		Color:        color,
		Consistency:  consistency,
		Floaters:     input.Floaters,
		Pain:         input.Pain,
		Strain:       input.Strain,
		Satisfaction: input.Satisfaction,
		PhotoURL:     input.PhotoURL,
		SmellLevel:   smellLevel,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		HasDetails:   false,
	}

	// Save to repository
	if err := s.repo.Create(ctx, bm); err != nil {
		return nil, fmt.Errorf("failed to create bowel movement: %w", err)
	}

	return bm, nil
}

// GetByID retrieves a bowel movement by ID
func (s *BowelMovementService) GetByID(ctx context.Context, id string) (*bowelmovement.BowelMovement, error) {
	if id == "" {
		return nil, bowelmovement.ErrInvalidID
	}

	bm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, bowelmovement.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get bowel movement: %w", err)
	}

	return bm, nil
}

// GetByUserID retrieves bowel movements for a specific user with pagination
func (s *BowelMovementService) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*bowelmovement.BowelMovement, error) {
	if userID == "" {
		return nil, bowelmovement.ErrInvalidUserID
	}

	// Apply business rules for pagination
	if limit <= 0 || limit > 100 {
		limit = 20 // default
	}
	if offset < 0 {
		offset = 0
	}

	movements, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bowel movements: %w", err)
	}

	return movements, nil
}

// Update updates an existing bowel movement
func (s *BowelMovementService) Update(ctx context.Context, id string, input *bowelmovement.UpdateBowelMovementInput) (*bowelmovement.BowelMovement, error) {
	if id == "" {
		return nil, bowelmovement.ErrInvalidID
	}

	// Get existing bowel movement to verify it exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, bowelmovement.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get bowel movement for update: %w", err)
	}

	// Validate update input
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// Convert input to update struct
	update := s.convertToUpdateStruct(input)

	// Save changes
	if err := s.repo.Update(ctx, id, update); err != nil {
		return nil, fmt.Errorf("failed to update bowel movement: %w", err)
	}

	// Return updated bowel movement
	return s.repo.GetByID(ctx, id)
}

// Delete removes a bowel movement
func (s *BowelMovementService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return bowelmovement.ErrInvalidID
	}

	// Check if exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return bowelmovement.ErrNotFound
		}
		return fmt.Errorf("failed to verify bowel movement exists: %w", err)
	}

	// Delete details first (if any)
	_ = s.repo.DeleteDetails(ctx, id) // Ignore error as details may not exist

	// Delete bowel movement
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete bowel movement: %w", err)
	}

	return nil
}

// CreateDetails creates details for a bowel movement
func (s *BowelMovementService) CreateDetails(ctx context.Context, bowelMovementID string, input *bowelmovement.CreateBowelMovementDetailsInput) (*bowelmovement.BowelMovementDetails, error) {
	if bowelMovementID == "" {
		return nil, bowelmovement.ErrInvalidID
	}

	// Verify bowel movement exists
	_, err := s.repo.GetByID(ctx, bowelMovementID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, bowelmovement.ErrNotFound
		}
		return nil, fmt.Errorf("failed to verify bowel movement exists: %w", err)
	}

	// Validate input
	if err := s.validateCreateDetailsInput(input); err != nil {
		return nil, err
	}

	// Create details
	details := &bowelmovement.BowelMovementDetails{
		ID:                uuid.New().String(),
		BowelMovementID:   bowelMovementID,
		Notes:             input.Notes,
		DetailedNotes:     input.DetailedNotes,
		Environment:       input.Environment,
		PreConditions:     input.PreConditions,
		PostConditions:    input.PostConditions,
		Tags:              input.Tags,
		WeatherCondition:  input.WeatherCondition,
		StressLevel:       input.StressLevel,
		SleepQuality:      input.SleepQuality,
		ExerciseIntensity: input.ExerciseIntensity,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save to repository
	if err := s.repo.CreateDetails(ctx, details); err != nil {
		return nil, fmt.Errorf("failed to create bowel movement details: %w", err)
	}

	return details, nil
}

// GetDetails retrieves details for a bowel movement
func (s *BowelMovementService) GetDetails(ctx context.Context, bowelMovementID string) (*bowelmovement.BowelMovementDetails, error) {
	if bowelMovementID == "" {
		return nil, bowelmovement.ErrInvalidID
	}

	details, err := s.repo.GetDetailsByBowelMovementID(ctx, bowelMovementID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, bowelmovement.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get bowel movement details: %w", err)
	}

	return details, nil
}

// UpdateDetails updates bowel movement details
func (s *BowelMovementService) UpdateDetails(ctx context.Context, bowelMovementID string, input *bowelmovement.UpdateBowelMovementDetailsInput) (*bowelmovement.BowelMovementDetails, error) {
	if bowelMovementID == "" {
		return nil, bowelmovement.ErrInvalidID
	}

	// Get existing details to verify they exist
	_, err := s.repo.GetDetailsByBowelMovementID(ctx, bowelMovementID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, bowelmovement.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get bowel movement details for update: %w", err)
	}

	// Validate update input
	if err := s.validateUpdateDetailsInput(input); err != nil {
		return nil, err
	}

	// Convert to details update struct
	update := s.convertToDetailsUpdateStruct(input)

	// Save changes
	if err := s.repo.UpdateDetails(ctx, bowelMovementID, update); err != nil {
		return nil, fmt.Errorf("failed to update bowel movement details: %w", err)
	}

	// Return updated details
	return s.repo.GetDetailsByBowelMovementID(ctx, bowelMovementID)
}

// GetByDateRange retrieves bowel movements within a date range
func (s *BowelMovementService) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*bowelmovement.BowelMovement, error) {
	if userID == "" {
		return nil, bowelmovement.ErrInvalidUserID
	}

	// Validate date range
	if start.After(end) {
		return nil, bowelmovement.ErrInvalidDateRange
	}

	// Limit date range to reasonable bounds
	maxRange := 365 * 24 * time.Hour // 1 year
	if end.Sub(start) > maxRange {
		return nil, bowelmovement.ErrDateRangeTooLarge
	}

	movements, err := s.repo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements by date range: %w", err)
	}

	return movements, nil
}

// GetUserStats generates analytics for a user's bowel movements
func (s *BowelMovementService) GetUserStats(ctx context.Context, userID string, start, end time.Time) (*bowelmovement.UserBowelMovementStats, error) {
	if userID == "" {
		return nil, bowelmovement.ErrInvalidUserID
	}

	// Get movements in date range
	movements, err := s.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	if len(movements) == 0 {
		return &bowelmovement.UserBowelMovementStats{
			TotalCount:          0,
			BristolDistribution: make(map[int]int),
		}, nil
	}

	// Calculate statistics
	stats := s.calculateStats(movements, start, end)
	return stats, nil
}

// GetLatest retrieves the most recent bowel movement for a user
func (s *BowelMovementService) GetLatest(ctx context.Context, userID string) (*bowelmovement.BowelMovement, error) {
	if userID == "" {
		return nil, bowelmovement.ErrInvalidUserID
	}

	latest, err := s.repo.GetLatestByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, bowelmovement.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest bowel movement: %w", err)
	}

	return latest, nil
}

// validateCreateInput validates create input
func (s *BowelMovementService) validateCreateInput(input *bowelmovement.CreateBowelMovementInput) error {
	if input == nil {
		return bowelmovement.ErrInvalidInput
	}

	if input.BristolType < 1 || input.BristolType > 7 {
		return bowelmovement.ErrInvalidBristolType
	}

	if input.Pain < 1 || input.Pain > 10 {
		return bowelmovement.ErrInvalidPainLevel
	}

	if input.Strain < 1 || input.Strain > 10 {
		return bowelmovement.ErrInvalidStrainLevel
	}

	if input.Satisfaction < 1 || input.Satisfaction > 10 {
		return bowelmovement.ErrInvalidSatisfactionLevel
	}

	// Validate enum values if provided
	if input.Volume != nil {
		v := shared.Volume(*input.Volume)
		if !v.IsValid() {
			return bowelmovement.ErrInvalidVolume
		}
	}

	if input.Color != nil {
		c := shared.Color(*input.Color)
		if !c.IsValid() {
			return bowelmovement.ErrInvalidColor
		}
	}

	if input.Consistency != nil {
		cons := shared.Consistency(*input.Consistency)
		if !cons.IsValid() {
			return bowelmovement.ErrInvalidConsistency
		}
	}

	if input.SmellLevel != nil {
		smell := shared.SmellLevel(*input.SmellLevel)
		if !smell.IsValid() {
			return bowelmovement.ErrInvalidSmellLevel
		}
	}

	return nil
}

// validateUpdateInput validates update input
func (s *BowelMovementService) validateUpdateInput(input *bowelmovement.UpdateBowelMovementInput) error {
	if input == nil {
		return bowelmovement.ErrInvalidInput
	}

	if input.BristolType != nil && (*input.BristolType < 1 || *input.BristolType > 7) {
		return bowelmovement.ErrInvalidBristolType
	}

	if input.Pain != nil && (*input.Pain < 1 || *input.Pain > 10) {
		return bowelmovement.ErrInvalidPainLevel
	}

	if input.Strain != nil && (*input.Strain < 1 || *input.Strain > 10) {
		return bowelmovement.ErrInvalidStrainLevel
	}

	if input.Satisfaction != nil && (*input.Satisfaction < 1 || *input.Satisfaction > 10) {
		return bowelmovement.ErrInvalidSatisfactionLevel
	}

	return nil
}

// validateCreateDetailsInput validates create details input
func (s *BowelMovementService) validateCreateDetailsInput(input *bowelmovement.CreateBowelMovementDetailsInput) error {
	if input == nil {
		return bowelmovement.ErrInvalidInput
	}

	if input.StressLevel != nil && (*input.StressLevel < 1 || *input.StressLevel > 10) {
		return bowelmovement.ErrInvalidStressLevel
	}

	if input.SleepQuality != nil && (*input.SleepQuality < 1 || *input.SleepQuality > 10) {
		return bowelmovement.ErrInvalidSleepQuality
	}

	if input.ExerciseIntensity != nil && (*input.ExerciseIntensity < 1 || *input.ExerciseIntensity > 10) {
		return bowelmovement.ErrInvalidExerciseIntensity
	}

	return nil
}

// validateUpdateDetailsInput validates update details input
func (s *BowelMovementService) validateUpdateDetailsInput(input *bowelmovement.UpdateBowelMovementDetailsInput) error {
	if input == nil {
		return bowelmovement.ErrInvalidInput
	}

	if input.StressLevel != nil && (*input.StressLevel < 1 || *input.StressLevel > 10) {
		return bowelmovement.ErrInvalidStressLevel
	}

	if input.SleepQuality != nil && (*input.SleepQuality < 1 || *input.SleepQuality > 10) {
		return bowelmovement.ErrInvalidSleepQuality
	}

	if input.ExerciseIntensity != nil && (*input.ExerciseIntensity < 1 || *input.ExerciseIntensity > 10) {
		return bowelmovement.ErrInvalidExerciseIntensity
	}

	return nil
}

// convertToUpdateStruct converts service input to repository update struct
func (s *BowelMovementService) convertToUpdateStruct(input *bowelmovement.UpdateBowelMovementInput) *bowelmovement.BowelMovementUpdate {
	update := &bowelmovement.BowelMovementUpdate{
		BristolType:  input.BristolType,
		Floaters:     input.Floaters,
		Pain:         input.Pain,
		Strain:       input.Strain,
		Satisfaction: input.Satisfaction,
		PhotoURL:     input.PhotoURL,
		RecordedAt:   input.RecordedAt,
	}

	// Convert string pointers to shared type pointers
	if input.Volume != nil {
		v := shared.Volume(*input.Volume)
		update.Volume = &v
	}

	if input.Color != nil {
		c := shared.Color(*input.Color)
		update.Color = &c
	}

	if input.Consistency != nil {
		cons := shared.Consistency(*input.Consistency)
		update.Consistency = &cons
	}

	if input.SmellLevel != nil {
		smell := shared.SmellLevel(*input.SmellLevel)
		update.SmellLevel = &smell
	}

	return update
}

// convertToDetailsUpdateStruct converts service input to repository details update struct
func (s *BowelMovementService) convertToDetailsUpdateStruct(input *bowelmovement.UpdateBowelMovementDetailsInput) *bowelmovement.BowelMovementDetailsUpdate {
	return &bowelmovement.BowelMovementDetailsUpdate{
		Notes:             input.Notes,
		DetailedNotes:     input.DetailedNotes,
		Environment:       input.Environment,
		PreConditions:     input.PreConditions,
		PostConditions:    input.PostConditions,
		Tags:              input.Tags,
		WeatherCondition:  input.WeatherCondition,
		StressLevel:       input.StressLevel,
		SleepQuality:      input.SleepQuality,
		ExerciseIntensity: input.ExerciseIntensity,
	}
}

// calculateStats calculates user statistics from bowel movements
func (s *BowelMovementService) calculateStats(movements []*bowelmovement.BowelMovement, start, end time.Time) *bowelmovement.UserBowelMovementStats {
	totalCount := int64(len(movements))

	var totalPain, totalStrain, totalSatisfaction int64
	bristolDistribution := make(map[int]int)
	bristolCounts := make(map[int]int)

	for _, bm := range movements {
		totalPain += int64(bm.Pain)
		totalStrain += int64(bm.Strain)
		totalSatisfaction += int64(bm.Satisfaction)

		bristolDistribution[bm.BristolType]++
		bristolCounts[bm.BristolType]++
	}

	// Calculate averages
	avgPain := float64(totalPain) / float64(totalCount)
	avgStrain := float64(totalStrain) / float64(totalCount)
	avgSatisfaction := float64(totalSatisfaction) / float64(totalCount)

	// Find most common Bristol type
	mostCommonBristol := 1
	maxCount := 0
	for bristol, count := range bristolCounts {
		if count > maxCount {
			maxCount = count
			mostCommonBristol = bristol
		}
	}

	// Calculate frequency per day
	days := end.Sub(start).Hours() / 24
	if days <= 0 {
		days = 1
	}
	frequencyPerDay := float64(totalCount) / days

	return &bowelmovement.UserBowelMovementStats{
		TotalCount:          totalCount,
		AveragePain:         avgPain,
		AverageStrain:       avgStrain,
		AverageSatisfaction: avgSatisfaction,
		MostCommonBristol:   mostCommonBristol,
		FrequencyPerDay:     frequencyPerDay,
		BristolDistribution: bristolDistribution,
	}
}
