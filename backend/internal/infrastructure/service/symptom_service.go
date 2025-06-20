package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// SymptomService implements the symptom business logic
type SymptomService struct {
	repo symptom.Repository
}

// NewSymptomService creates a new SymptomService instance
func NewSymptomService(repo symptom.Repository) *SymptomService {
	return &SymptomService{
		repo: repo,
	}
}

// Create creates a new symptom record
func (s *SymptomService) Create(ctx context.Context, userID string, input *symptom.CreateSymptomInput) (*symptom.Symptom, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if input == nil {
		return nil, shared.ErrValidationFailed
	}

	// Validate required fields
	if input.Name == "" {
		return nil, symptom.ErrSymptomNameRequired
	}

	if input.Severity < 1 || input.Severity > 10 {
		return nil, symptom.ErrInvalidSeverity
	}

	// Create the symptom
	newSymptom := symptom.NewSymptom(userID, input.Name, input.Severity, input.RecordedAt)
	newSymptom.ID = uuid.New().String()

	// Set optional fields
	if input.Description != "" {
		newSymptom.Description = input.Description
	}
	if input.Category != nil {
		category := shared.SymptomCategory(*input.Category)
		if category.IsValid() {
			newSymptom.Category = &category
		} else {
			return nil, symptom.ErrInvalidSymptomCategory
		}
	}
	if input.Duration != nil {
		if *input.Duration < 1 {
			return nil, symptom.ErrInvalidDuration
		}
		newSymptom.Duration = input.Duration
	}
	if input.BodyPart != "" {
		newSymptom.BodyPart = input.BodyPart
	}
	if input.Type != nil {
		symptomType := shared.SymptomType(*input.Type)
		if symptomType.IsValid() {
			newSymptom.Type = &symptomType
		} else {
			return nil, symptom.ErrInvalidSymptomType
		}
	}
	if len(input.Triggers) > 0 {
		newSymptom.Triggers = input.Triggers
	}
	if input.Notes != "" {
		newSymptom.Notes = input.Notes
	}
	if input.PhotoURL != "" {
		newSymptom.PhotoURL = input.PhotoURL
	}

	// Save to repository
	if err := s.repo.Create(ctx, &newSymptom); err != nil {
		return nil, fmt.Errorf("failed to create symptom: %w", err)
	}

	return &newSymptom, nil
}

// GetByID retrieves a symptom by its ID
func (s *SymptomService) GetByID(ctx context.Context, id string) (*symptom.Symptom, error) {
	if id == "" {
		return nil, shared.ErrValidationFailed
	}

	return s.repo.GetByID(ctx, id)
}

// GetByUserID retrieves symptoms for a specific user with pagination
func (s *SymptomService) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*symptom.Symptom, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

// Update updates an existing symptom
func (s *SymptomService) Update(ctx context.Context, id string, input *symptom.UpdateSymptomInput) (*symptom.Symptom, error) {
	if id == "" {
		return nil, shared.ErrValidationFailed
	}

	if input == nil {
		return nil, shared.ErrValidationFailed
	}

	// Check if symptom exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate inputs
	if input.Severity != nil && (*input.Severity < 1 || *input.Severity > 10) {
		return nil, symptom.ErrInvalidSeverity
	}

	if input.Duration != nil && *input.Duration < 1 {
		return nil, symptom.ErrInvalidDuration
	}

	// Create update struct
	update := &symptom.SymptomUpdate{
		Name:        input.Name,
		Description: input.Description,
		RecordedAt:  input.RecordedAt,
		Severity:    input.Severity,
		Duration:    input.Duration,
		BodyPart:    input.BodyPart,
		Triggers:    input.Triggers,
		Notes:       input.Notes,
		PhotoURL:    input.PhotoURL,
	}

	// Handle category conversion
	if input.Category != nil {
		category := shared.SymptomCategory(*input.Category)
		if category.IsValid() {
			update.Category = &category
		} else {
			return nil, symptom.ErrInvalidSymptomCategory
		}
	}

	// Handle type conversion
	if input.Type != nil {
		symptomType := shared.SymptomType(*input.Type)
		if symptomType.IsValid() {
			update.Type = &symptomType
		} else {
			return nil, symptom.ErrInvalidSymptomType
		}
	}

	// Update in repository
	if err := s.repo.Update(ctx, id, update); err != nil {
		return nil, fmt.Errorf("failed to update symptom: %w", err)
	}

	// Return updated symptom
	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// If we can't get the updated record, return the original with basic updates applied
		if input.Name != nil {
			existing.Name = *input.Name
		}
		if input.Severity != nil {
			existing.Severity = *input.Severity
		}
		existing.UpdatedAt = time.Now()
		return existing, nil
	}

	return updated, nil
}

// Delete deletes a symptom
func (s *SymptomService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return shared.ErrValidationFailed
	}

	// Check if symptom exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// GetByDateRange retrieves symptoms within a date range
func (s *SymptomService) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*symptom.Symptom, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	return s.repo.GetByDateRange(ctx, userID, start, end)
}

// GetByCategory retrieves symptoms by category
func (s *SymptomService) GetByCategory(ctx context.Context, userID string, category string) ([]*symptom.Symptom, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if category == "" {
		return nil, shared.ErrValidationFailed
	}

	return s.repo.GetByCategory(ctx, userID, category)
}

// GetBySeverityRange retrieves symptoms within a severity range
func (s *SymptomService) GetBySeverityRange(ctx context.Context, userID string, minSeverity, maxSeverity int) ([]*symptom.Symptom, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if minSeverity < 1 || maxSeverity > 10 || minSeverity > maxSeverity {
		return nil, symptom.ErrInvalidSeverity
	}

	return s.repo.GetBySeverity(ctx, userID, minSeverity, maxSeverity)
}

// GetLatest retrieves the most recent symptom for a user
func (s *SymptomService) GetLatest(ctx context.Context, userID string) (*symptom.Symptom, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	return s.repo.GetLatestByUserID(ctx, userID)
}

// GetSymptomStats retrieves symptom statistics for analytics
func (s *SymptomService) GetSymptomStats(ctx context.Context, userID string, start, end time.Time) (*symptom.SymptomStats, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	// Get symptoms for the period
	symptoms, err := s.repo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve symptoms for stats: %w", err)
	}

	if len(symptoms) == 0 {
		return &symptom.SymptomStats{
			TotalCount:           0,
			AverageSeverity:      0,
			CategoryBreakdown:    make(map[string]int),
			TypeBreakdown:        make(map[string]int),
			SeverityDistribution: make(map[int]int),
			AverageDuration:      0,
		}, nil
	}

	// Calculate statistics
	stats := &symptom.SymptomStats{
		TotalCount:           int64(len(symptoms)),
		CategoryBreakdown:    make(map[string]int),
		TypeBreakdown:        make(map[string]int),
		SeverityDistribution: make(map[int]int),
	}

	var totalSeverity int
	var totalDuration int
	var durationCount int
	categoryCount := make(map[string]int)
	typeCount := make(map[string]int)

	for _, sym := range symptoms {
		totalSeverity += sym.Severity
		stats.SeverityDistribution[sym.Severity]++

		if sym.Category != nil {
			categoryName := string(*sym.Category)
			categoryCount[categoryName]++
		}

		if sym.Type != nil {
			typeName := string(*sym.Type)
			typeCount[typeName]++
		}

		if sym.Duration != nil {
			totalDuration += *sym.Duration
			durationCount++
		}
	}

	stats.AverageSeverity = float64(totalSeverity) / float64(len(symptoms))
	if durationCount > 0 {
		stats.AverageDuration = float64(totalDuration) / float64(durationCount)
	}

	// Find most common category and type
	maxCatCount := 0
	maxTypeCount := 0
	for cat, count := range categoryCount {
		stats.CategoryBreakdown[cat] = count
		if count > maxCatCount {
			maxCatCount = count
			stats.MostCommonCategory = cat
		}
	}

	for typ, count := range typeCount {
		stats.TypeBreakdown[typ] = count
		if count > maxTypeCount {
			maxTypeCount = count
			stats.MostCommonType = typ
		}
	}

	return stats, nil
}

// GetTriggerInsights retrieves insights about symptom triggers
func (s *SymptomService) GetTriggerInsights(ctx context.Context, userID string, start, end time.Time) (*symptom.TriggerInsights, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	// Get symptoms for the period
	symptoms, err := s.repo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve symptoms for trigger insights: %w", err)
	}

	insights := &symptom.TriggerInsights{
		TriggerFrequency:   make(map[string]int),
		TriggerSeverityMap: make(map[string]float64),
	}

	if len(symptoms) == 0 {
		return insights, nil
	}

	// Analyze triggers
	triggerSeveritySum := make(map[string]int)
	triggerOccurrences := make(map[string]int)

	for _, sym := range symptoms {
		for _, trigger := range sym.Triggers {
			if trigger != "" {
				insights.TriggerFrequency[trigger]++
				triggerSeveritySum[trigger] += sym.Severity
				triggerOccurrences[trigger]++
			}
		}
	}

	// Calculate average severity per trigger
	for trigger, count := range triggerOccurrences {
		if count > 0 {
			insights.TriggerSeverityMap[trigger] = float64(triggerSeveritySum[trigger]) / float64(count)
		}
	}

	insights.UniqueTriggerCount = len(insights.TriggerFrequency)

	// Find most common triggers (top 10)
	type triggerInfo struct {
		trigger string
		count   int
	}

	var triggers []triggerInfo
	for trigger, count := range insights.TriggerFrequency {
		triggers = append(triggers, triggerInfo{trigger, count})
	}

	// Simple sort by count (descending)
	for i := 0; i < len(triggers)-1; i++ {
		for j := i + 1; j < len(triggers); j++ {
			if triggers[j].count > triggers[i].count {
				triggers[i], triggers[j] = triggers[j], triggers[i]
			}
		}
	}

	// Take top 10
	limit := 10
	if len(triggers) < limit {
		limit = len(triggers)
	}

	for i := 0; i < limit; i++ {
		insights.MostCommonTriggers = append(insights.MostCommonTriggers, triggers[i].trigger)
	}

	return insights, nil
}

// GetSymptomPatterns retrieves patterns in symptom occurrence
func (s *SymptomService) GetSymptomPatterns(ctx context.Context, userID string, start, end time.Time) (*symptom.SymptomPatterns, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	// Get symptoms for the period
	symptoms, err := s.repo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve symptoms for pattern analysis: %w", err)
	}

	patterns := &symptom.SymptomPatterns{
		TimeOfDayPatterns: make(map[string]float64),
		DayOfWeekPatterns: make(map[string]float64),
		SeasonalPatterns:  make(map[string]float64),
	}

	if len(symptoms) == 0 {
		return patterns, nil
	}

	// Analyze patterns
	totalDays := end.Sub(start).Hours() / 24
	if totalDays <= 0 {
		totalDays = 1
	}

	patterns.FrequencyPerDay = float64(len(symptoms)) / totalDays

	hourCounts := make(map[int]int)
	weekdayCounts := make(map[string]int)
	seasonCounts := make(map[string]int)

	var totalTimeBetween float64
	var timeBetweenCount int

	for i, sym := range symptoms {
		// Time of day patterns
		hour := sym.RecordedAt.Hour()
		hourCounts[hour]++

		// Day of week patterns
		weekday := sym.RecordedAt.Weekday().String()
		weekdayCounts[weekday]++

		// Seasonal patterns (simplified)
		month := sym.RecordedAt.Month()
		var season string
		switch {
		case month >= 3 && month <= 5:
			season = "Spring"
		case month >= 6 && month <= 8:
			season = "Summer"
		case month >= 9 && month <= 11:
			season = "Fall"
		default:
			season = "Winter"
		}
		seasonCounts[season]++

		// Time between symptoms
		if i > 0 {
			timeDiff := sym.RecordedAt.Sub(symptoms[i-1].RecordedAt).Hours()
			totalTimeBetween += timeDiff
			timeBetweenCount++
		}
	}

	// Convert counts to frequencies (normalize)
	totalSymptoms := float64(len(symptoms))

	for hour, count := range hourCounts {
		patterns.TimeOfDayPatterns[fmt.Sprintf("%02d:00", hour)] = float64(count) / totalSymptoms
	}

	for weekday, count := range weekdayCounts {
		patterns.DayOfWeekPatterns[weekday] = float64(count) / totalSymptoms
	}

	for season, count := range seasonCounts {
		patterns.SeasonalPatterns[season] = float64(count) / totalSymptoms
	}

	if timeBetweenCount > 0 {
		patterns.AverageTimeBetween = totalTimeBetween / float64(timeBetweenCount)
	}

	return patterns, nil
}
