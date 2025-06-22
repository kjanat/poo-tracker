package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// BowelMovementRepository implements bowelmovement.Repository using in-memory storage
type BowelMovementRepository struct {
	mu             sync.RWMutex
	bowelMovements map[string]*bowelmovement.BowelMovement
	details        map[string]*bowelmovement.BowelMovementDetails // keyed by bowelMovementID
}

// NewBowelMovementRepository creates a new in-memory bowel movement repository
func NewBowelMovementRepository() bowelmovement.Repository {
	return &BowelMovementRepository{
		bowelMovements: make(map[string]*bowelmovement.BowelMovement),
		details:        make(map[string]*bowelmovement.BowelMovementDetails),
	}
}

// BowelMovement CRUD operations
func (r *BowelMovementRepository) Create(ctx context.Context, bm *bowelmovement.BowelMovement) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bowelMovements[bm.ID] = bm
	return nil
}

func (r *BowelMovementRepository) GetByID(ctx context.Context, id string) (*bowelmovement.BowelMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bm, exists := r.bowelMovements[id]
	if !exists {
		return nil, shared.ErrNotFound
	}

	return bm, nil
}

func (r *BowelMovementRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*bowelmovement.BowelMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMovements []*bowelmovement.BowelMovement
	for _, bm := range r.bowelMovements {
		if bm.UserID == userID {
			userMovements = append(userMovements, bm)
		}
	}

	// Sort by RecordedAt descending
	sort.Slice(userMovements, func(i, j int) bool {
		return userMovements[i].RecordedAt.After(userMovements[j].RecordedAt)
	})

	// Apply pagination
	start := offset
	if start > len(userMovements) {
		return []*bowelmovement.BowelMovement{}, nil
	}

	end := start + limit
	if end > len(userMovements) {
		end = len(userMovements)
	}

	return userMovements[start:end], nil
}

func (r *BowelMovementRepository) Update(ctx context.Context, id string, update *bowelmovement.BowelMovementUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	bm, exists := r.bowelMovements[id]
	if !exists {
		return shared.ErrNotFound
	}

	// Apply updates
	if update.BristolType != nil {
		bm.BristolType = *update.BristolType
	}
	if update.Volume != nil {
		bm.Volume = update.Volume
	}
	if update.Color != nil {
		bm.Color = update.Color
	}
	if update.Consistency != nil {
		bm.Consistency = update.Consistency
	}
	if update.Floaters != nil {
		bm.Floaters = *update.Floaters
	}
	if update.Pain != nil {
		bm.Pain = *update.Pain
	}
	if update.Strain != nil {
		bm.Strain = *update.Strain
	}
	if update.Satisfaction != nil {
		bm.Satisfaction = *update.Satisfaction
	}
	if update.PhotoURL != nil {
		bm.PhotoURL = *update.PhotoURL
	}
	if update.SmellLevel != nil {
		bm.SmellLevel = update.SmellLevel
	}
	if update.RecordedAt != nil {
		bm.RecordedAt = *update.RecordedAt
	}

	bm.UpdatedAt = time.Now()
	return nil
}

func (r *BowelMovementRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bowelMovements[id]; !exists {
		return shared.ErrNotFound
	}

	delete(r.bowelMovements, id)
	delete(r.details, id) // Also delete any associated details
	return nil
}

// BowelMovementDetails operations
func (r *BowelMovementRepository) CreateDetails(ctx context.Context, details *bowelmovement.BowelMovementDetails) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if bowel movement exists
	if _, exists := r.bowelMovements[details.BowelMovementID]; !exists {
		return shared.ErrNotFound
	}

	r.details[details.BowelMovementID] = details

	// Update HasDetails flag on bowel movement
	if bm, exists := r.bowelMovements[details.BowelMovementID]; exists {
		bm.HasDetails = true
	}

	return nil
}

func (r *BowelMovementRepository) GetDetailsByBowelMovementID(ctx context.Context, bowelMovementID string) (*bowelmovement.BowelMovementDetails, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	details, exists := r.details[bowelMovementID]
	if !exists {
		return nil, shared.ErrNotFound
	}

	return details, nil
}

func (r *BowelMovementRepository) UpdateDetails(ctx context.Context, bowelMovementID string, update *bowelmovement.BowelMovementDetailsUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	details, exists := r.details[bowelMovementID]
	if !exists {
		return shared.ErrNotFound
	}

	// Apply updates
	if update.Notes != nil {
		details.Notes = *update.Notes
	}
	if update.DetailedNotes != nil {
		details.DetailedNotes = *update.DetailedNotes
	}
	if update.Environment != nil {
		details.Environment = *update.Environment
	}
	if update.PreConditions != nil {
		details.PreConditions = *update.PreConditions
	}
	if update.PostConditions != nil {
		details.PostConditions = *update.PostConditions
	}
	if update.Tags != nil {
		details.Tags = update.Tags
	}
	if update.WeatherCondition != nil {
		details.WeatherCondition = *update.WeatherCondition
	}
	if update.StressLevel != nil {
		details.StressLevel = update.StressLevel
	}
	if update.SleepQuality != nil {
		details.SleepQuality = update.SleepQuality
	}
	if update.ExerciseIntensity != nil {
		details.ExerciseIntensity = update.ExerciseIntensity
	}

	details.UpdatedAt = time.Now()
	return nil
}

func (r *BowelMovementRepository) DeleteDetails(ctx context.Context, bowelMovementID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.details[bowelMovementID]; !exists {
		return shared.ErrNotFound
	}

	delete(r.details, bowelMovementID)

	// Update HasDetails flag on bowel movement
	if bm, exists := r.bowelMovements[bowelMovementID]; exists {
		bm.HasDetails = false
	}

	return nil
}

// Query operations
func (r *BowelMovementRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*bowelmovement.BowelMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*bowelmovement.BowelMovement
	for _, bm := range r.bowelMovements {
		if bm.UserID == userID &&
			!bm.RecordedAt.Before(start) &&
			!bm.RecordedAt.After(end) {
			result = append(result, bm)
		}
	}

	// Sort by RecordedAt ascending
	sort.Slice(result, func(i, j int) bool {
		return result[i].RecordedAt.Before(result[j].RecordedAt)
	})

	return result, nil
}

func (r *BowelMovementRepository) GetByBristolType(ctx context.Context, userID string, bristolType int) ([]*bowelmovement.BowelMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*bowelmovement.BowelMovement
	for _, bm := range r.bowelMovements {
		if bm.UserID == userID && bm.BristolType == bristolType {
			result = append(result, bm)
		}
	}

	return result, nil
}

func (r *BowelMovementRepository) GetCountByUserID(ctx context.Context, userID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := int64(0)
	for _, bm := range r.bowelMovements {
		if bm.UserID == userID {
			count++
		}
	}

	return count, nil
}

func (r *BowelMovementRepository) GetLatestByUserID(ctx context.Context, userID string) (*bowelmovement.BowelMovement, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *bowelmovement.BowelMovement
	for _, bm := range r.bowelMovements {
		if bm.UserID == userID {
			if latest == nil || bm.RecordedAt.After(latest.RecordedAt) {
				latest = bm
			}
		}
	}

	if latest == nil {
		return nil, shared.ErrNotFound
	}

	return latest, nil
}

// Analytics operations
func (r *BowelMovementRepository) GetAveragesByUserID(ctx context.Context, userID string, start, end time.Time) (*bowelmovement.BowelMovementAverages, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var movements []*bowelmovement.BowelMovement
	for _, bm := range r.bowelMovements {
		if bm.UserID == userID &&
			!bm.RecordedAt.Before(start) &&
			!bm.RecordedAt.After(end) {
			movements = append(movements, bm)
		}
	}

	if len(movements) == 0 {
		return &bowelmovement.BowelMovementAverages{}, nil
	}

	var totalPain, totalStrain, totalSatisfaction float64
	bristolCounts := make(map[int]int)

	for _, bm := range movements {
		totalPain += float64(bm.Pain)
		totalStrain += float64(bm.Strain)
		totalSatisfaction += float64(bm.Satisfaction)
		bristolCounts[bm.BristolType]++
	}

	count := float64(len(movements))

	// Find most common Bristol type
	mostCommon := 0
	maxCount := 0
	for bristol, count := range bristolCounts {
		if count > maxCount {
			maxCount = count
			mostCommon = bristol
		}
	}

	return &bowelmovement.BowelMovementAverages{
		AveragePain:         totalPain / count,
		AverageStrain:       totalStrain / count,
		AverageSatisfaction: totalSatisfaction / count,
		MostCommonBristol:   mostCommon,
		TotalCount:          int64(len(movements)),
	}, nil
}

func (r *BowelMovementRepository) GetFrequencyByUserID(ctx context.Context, userID string, start, end time.Time) (map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	frequency := make(map[string]int)

	for _, bm := range r.bowelMovements {
		if bm.UserID == userID &&
			!bm.RecordedAt.Before(start) &&
			!bm.RecordedAt.After(end) {
			date := bm.RecordedAt.Format("2006-01-02")
			frequency[date]++
		}
	}

	return frequency, nil
}
