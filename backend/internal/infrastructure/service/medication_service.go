package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// MedicationService implements the medication business logic
type MedicationService struct {
	repo medication.Repository
}

// NewMedicationService creates a new MedicationService instance
func NewMedicationService(repo medication.Repository) *MedicationService {
	return &MedicationService{
		repo: repo,
	}
}

// Create creates a new medication record
func (s *MedicationService) Create(ctx context.Context, userID string, input *medication.CreateMedicationInput) (*medication.Medication, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if input == nil {
		return nil, shared.ErrValidationFailed
	}

	// Validate required fields
	if input.Name == "" {
		return nil, medication.ErrMedicationNameRequired
	}

	if input.Dosage == "" {
		return nil, medication.ErrDosageRequired
	}

	if input.Frequency == "" {
		return nil, medication.ErrFrequencyRequired
	}

	// Create the medication
	newMedication := medication.NewMedication(userID, input.Name, input.Dosage, input.Frequency)
	newMedication.ID = uuid.New().String()

	// Set optional fields
	if input.GenericName != "" {
		newMedication.GenericName = input.GenericName
	}
	if input.Brand != "" {
		newMedication.Brand = input.Brand
	}
	if input.Category != nil {
		category := shared.MedicationCategory(*input.Category)
		if category.IsValid() {
			newMedication.Category = &category
		} else {
			return nil, medication.ErrInvalidMedicationCategory
		}
	}
	if input.Form != nil {
		form := shared.MedicationForm(*input.Form)
		if form.IsValid() {
			newMedication.Form = &form
		} else {
			return nil, medication.ErrInvalidMedicationForm
		}
	}
	if input.Route != nil {
		route := shared.MedicationRoute(*input.Route)
		if route.IsValid() {
			newMedication.Route = &route
		} else {
			return nil, medication.ErrInvalidMedicationRoute
		}
	}
	if input.StartDate != nil {
		newMedication.StartDate = input.StartDate
	}
	if input.EndDate != nil {
		newMedication.EndDate = input.EndDate
	}
	if input.Purpose != "" {
		newMedication.Purpose = input.Purpose
	}
	if len(input.SideEffects) > 0 {
		newMedication.SideEffects = input.SideEffects
	}
	if input.Notes != "" {
		newMedication.Notes = input.Notes
	}
	if input.PhotoURL != "" {
		newMedication.PhotoURL = input.PhotoURL
	}
	newMedication.IsAsNeeded = input.IsAsNeeded

	// Save to repository
	if err := s.repo.Create(ctx, &newMedication); err != nil {
		return nil, fmt.Errorf("failed to create medication: %w", err)
	}

	return &newMedication, nil
}

// GetByID retrieves a medication by its ID
func (s *MedicationService) GetByID(ctx context.Context, id string) (*medication.Medication, error) {
	if id == "" {
		return nil, shared.ErrValidationFailed
	}

	return s.repo.GetByID(ctx, id)
}

// GetByUserID retrieves medications for a specific user with pagination
func (s *MedicationService) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*medication.Medication, error) {
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

// Update updates an existing medication
func (s *MedicationService) Update(ctx context.Context, id string, input *medication.UpdateMedicationInput) (*medication.Medication, error) {
	if id == "" {
		return nil, shared.ErrValidationFailed
	}

	if input == nil {
		return nil, shared.ErrValidationFailed
	}

	// Check if medication exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create update struct
	update := &medication.MedicationUpdate{
		Name:        input.Name,
		GenericName: input.GenericName,
		Brand:       input.Brand,
		Dosage:      input.Dosage,
		Frequency:   input.Frequency,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Purpose:     input.Purpose,
		SideEffects: input.SideEffects,
		Notes:       input.Notes,
		PhotoURL:    input.PhotoURL,
		IsAsNeeded:  input.IsAsNeeded,
	}

	// Handle category conversion
	if input.Category != nil {
		category := shared.MedicationCategory(*input.Category)
		if category.IsValid() {
			update.Category = &category
		} else {
			return nil, medication.ErrInvalidMedicationCategory
		}
	}

	// Handle form conversion
	if input.Form != nil {
		form := shared.MedicationForm(*input.Form)
		if form.IsValid() {
			update.Form = &form
		} else {
			return nil, medication.ErrInvalidMedicationForm
		}
	}

	// Handle route conversion
	if input.Route != nil {
		route := shared.MedicationRoute(*input.Route)
		if route.IsValid() {
			update.Route = &route
		} else {
			return nil, medication.ErrInvalidMedicationRoute
		}
	}

	// Update in repository
	if err := s.repo.Update(ctx, id, update); err != nil {
		return nil, fmt.Errorf("failed to update medication: %w", err)
	}

	// Return updated medication
	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// If we can't get the updated record, return the original with basic updates applied
		if input.Name != nil {
			existing.Name = *input.Name
		}
		existing.UpdatedAt = time.Now()
		return existing, nil
	}

	return updated, nil
}

// Delete deletes a medication
func (s *MedicationService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return shared.ErrValidationFailed
	}

	// Check if medication exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// GetActiveMedications retrieves all active medications for a user
func (s *MedicationService) GetActiveMedications(ctx context.Context, userID string) ([]*medication.Medication, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	return s.repo.GetActiveByUserID(ctx, userID)
}

// GetAsNeededMedications retrieves all as-needed medications for a user
func (s *MedicationService) GetAsNeededMedications(ctx context.Context, userID string) ([]*medication.Medication, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	return s.repo.GetAsNeededByUserID(ctx, userID)
}

// DeactivateMedication deactivates a medication
func (s *MedicationService) DeactivateMedication(ctx context.Context, id string) error {
	if id == "" {
		return shared.ErrValidationFailed
	}

	// Check if medication exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !existing.IsActive {
		return medication.ErrMedicationNotActive
	}

	// Update to deactivate
	isActive := false
	update := &medication.MedicationUpdate{
		IsActive: &isActive,
	}

	return s.repo.Update(ctx, id, update)
}

// ReactivateMedication reactivates a medication
func (s *MedicationService) ReactivateMedication(ctx context.Context, id string) error {
	if id == "" {
		return shared.ErrValidationFailed
	}

	// Check if medication exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existing.IsActive {
		return medication.ErrMedicationAlreadyActive
	}

	// Update to reactivate
	isActive := true
	update := &medication.MedicationUpdate{
		IsActive: &isActive,
	}

	return s.repo.Update(ctx, id, update)
}

// RecordDose records a dose taken
func (s *MedicationService) RecordDose(ctx context.Context, medicationID string, input *medication.RecordDoseInput) error {
	if medicationID == "" {
		return shared.ErrValidationFailed
	}

	if input == nil {
		return shared.ErrValidationFailed
	}

	// Check if medication exists
	_, err := s.repo.GetByID(ctx, medicationID)
	if err != nil {
		return err
	}

	// Validate taken time is not in the future
	if input.TakenAt.After(time.Now()) {
		return medication.ErrDoseTakenInFuture
	}

	return s.repo.RecordDose(ctx, medicationID, input.TakenAt)
}

// GetDoseHistory retrieves dose history for a medication
func (s *MedicationService) GetDoseHistory(ctx context.Context, medicationID string, start, end time.Time) ([]*medication.DoseRecord, error) {
	if medicationID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	// Check if medication exists
	_, err := s.repo.GetByID(ctx, medicationID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetDoseHistory(ctx, medicationID, start, end)
}

// GetMedicationStats retrieves medication statistics for analytics
func (s *MedicationService) GetMedicationStats(ctx context.Context, userID string, start, end time.Time) (*medication.MedicationStats, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	// Get medications for the period
	medications, err := s.repo.GetByUserID(ctx, userID, 1000, 0) // Get all medications
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve medications for stats: %w", err)
	}

	stats := &medication.MedicationStats{
		CategoryBreakdown: make(map[string]int),
		FormBreakdown:     make(map[string]int),
		RouteBreakdown:    make(map[string]int),
	}

	if len(medications) == 0 {
		return stats, nil
	}

	var activeMeds, asNeededMeds, regularMeds int64

	for _, med := range medications {
		stats.TotalMedications++

		if med.IsActive {
			activeMeds++
		}

		if med.IsAsNeeded {
			asNeededMeds++
		} else {
			regularMeds++
		}

		// Category breakdown
		if med.Category != nil {
			categoryName := string(*med.Category)
			stats.CategoryBreakdown[categoryName]++
		}

		// Form breakdown
		if med.Form != nil {
			formName := string(*med.Form)
			stats.FormBreakdown[formName]++
		}

		// Route breakdown
		if med.Route != nil {
			routeName := string(*med.Route)
			stats.RouteBreakdown[routeName]++
		}
	}

	stats.ActiveMedications = activeMeds
	stats.AsNeededMedications = asNeededMeds
	stats.RegularMedications = regularMeds

	// Calculate dose statistics for the period
	totalDoses := int64(0)
	for _, med := range medications {
		doseHistory, err := s.repo.GetDoseHistory(ctx, med.ID, start, end)
		if err == nil {
			totalDoses += int64(len(doseHistory))
		}
	}

	stats.DosesThisPeriod = totalDoses
	days := end.Sub(start).Hours() / 24
	if days > 0 {
		stats.AverageDoesPerDay = float64(totalDoses) / days
	}

	return stats, nil
}

// GetMedicationInsights retrieves insights from medication data
func (s *MedicationService) GetMedicationInsights(ctx context.Context, userID string, start, end time.Time) (*medication.MedicationInsights, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	if start.After(end) {
		return nil, shared.ErrInvalidDateRange
	}

	// Get basic stats first
	stats, err := s.GetMedicationStats(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication stats for insights: %w", err)
	}

	insights := &medication.MedicationInsights{
		MedicationBurden: int(stats.ActiveMedications),
	}

	// Find most common category, form, and route
	maxCatCount := 0
	maxFormCount := 0
	maxRouteCount := 0

	for category, count := range stats.CategoryBreakdown {
		if count > maxCatCount {
			maxCatCount = count
			insights.MostUsedCategory = category
		}
	}

	for form, count := range stats.FormBreakdown {
		if count > maxFormCount {
			maxFormCount = count
			insights.MostCommonForm = form
		}
	}

	for route, count := range stats.RouteBreakdown {
		if count > maxRouteCount {
			maxRouteCount = count
			insights.MostCommonRoute = route
		}
	}

	// Calculate adherence score (simplified - based on dose consistency for regular meds)
	if stats.RegularMedications > 0 {
		// This is a simplified calculation - in a real system you'd track scheduled vs actual doses
		expectedDoses := stats.RegularMedications * int64(end.Sub(start).Hours()/24) // Assuming once daily
		if expectedDoses > 0 {
			insights.AdhereanceScore = float64(stats.DosesThisPeriod) / float64(expectedDoses)
			if insights.AdhereanceScore > 1.0 {
				insights.AdhereanceScore = 1.0 // Cap at 100%
			}
		}
	}

	// Calculate complexity score (based on number of medications and frequency)
	insights.ComplexityScore = float64(stats.TotalMedications) / 10.0 // Simple scale where 10+ meds = max complexity
	if insights.ComplexityScore > 1.0 {
		insights.ComplexityScore = 1.0
	}

	// Determine interaction risk (simplified)
	switch {
	case stats.TotalMedications >= 10:
		insights.InteractionRisk = "HIGH"
	case stats.TotalMedications >= 5:
		insights.InteractionRisk = "MEDIUM"
	default:
		insights.InteractionRisk = "LOW"
	}

	return insights, nil
}

// GetSideEffectAnalysis retrieves analysis of reported side effects
func (s *MedicationService) GetSideEffectAnalysis(ctx context.Context, userID string) (*medication.SideEffectAnalysis, error) {
	if userID == "" {
		return nil, shared.ErrValidationFailed
	}

	// Get all medications for the user
	medications, err := s.repo.GetByUserID(ctx, userID, 1000, 0) // Get all medications
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve medications for side effect analysis: %w", err)
	}

	analysis := &medication.SideEffectAnalysis{
		SideEffectFrequency:     make(map[string]int),
		MedicationSideEffectMap: make(map[string][]string),
	}

	if len(medications) == 0 {
		return analysis, nil
	}

	allSideEffects := make(map[string]bool)

	for _, med := range medications {
		if len(med.SideEffects) > 0 {
			analysis.MedicationSideEffectMap[med.Name] = med.SideEffects

			for _, sideEffect := range med.SideEffects {
				if sideEffect != "" {
					analysis.SideEffectFrequency[sideEffect]++
					allSideEffects[sideEffect] = true
				}
			}
		}
	}

	analysis.TotalUniqueSideEffects = len(allSideEffects)

	// Find most common side effects (top 10)
	type sideEffectInfo struct {
		effect string
		count  int
	}

	var sideEffects []sideEffectInfo
	for effect, count := range analysis.SideEffectFrequency {
		sideEffects = append(sideEffects, sideEffectInfo{effect, count})
	}

	// Simple sort by count (descending)
	for i := 0; i < len(sideEffects)-1; i++ {
		for j := i + 1; j < len(sideEffects); j++ {
			if sideEffects[j].count > sideEffects[i].count {
				sideEffects[i], sideEffects[j] = sideEffects[j], sideEffects[i]
			}
		}
	}

	// Take top 10
	limit := 10
	if len(sideEffects) < limit {
		limit = len(sideEffects)
	}

	for i := 0; i < limit; i++ {
		analysis.MostCommonSideEffects = append(analysis.MostCommonSideEffects, sideEffects[i].effect)
	}

	return analysis, nil
}
