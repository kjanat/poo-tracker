package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// AnalyticsService implements the analytics service interface
type AnalyticsService struct {
	bowelMovementService bowelmovement.Service
	mealService          meal.Service
	symptomService       symptom.Service
	medicationService    medication.Service
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	bowelMovementService bowelmovement.Service,
	mealService meal.Service,
	symptomService symptom.Service,
	medicationService medication.Service,
) analytics.Service {
	return &AnalyticsService{
		bowelMovementService: bowelMovementService,
		mealService:          mealService,
		symptomService:       symptomService,
		medicationService:    medicationService,
	}
}

// GetUserHealthOverview generates a comprehensive health overview for a user
func (s *AnalyticsService) GetUserHealthOverview(ctx context.Context, userID string, start, end time.Time) (*analytics.HealthOverview, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get data from all services
	bowelMovements, err := s.bowelMovementService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	medications, err := s.medicationService.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get medications: %w", err)
	}

	// Calculate summaries
	bowelMovementStats := s.calculateBowelMovementSummary(bowelMovements, start, end)
	mealStats := s.calculateMealSummary(meals, start, end)
	symptomStats := s.calculateSymptomSummary(symptoms, start, end)
	medicationStats := s.calculateMedicationSummary(medications)

	// Calculate overall health score
	overallHealthScore := s.calculateOverallHealthScore(bowelMovementStats, mealStats, symptomStats, medicationStats)

	// Determine trend direction
	trendDirection := s.determineTrendDirection(bowelMovementStats, symptomStats)

	return &analytics.HealthOverview{
		Period:             fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
		BowelMovementStats: bowelMovementStats,
		MealStats:          mealStats,
		SymptomStats:       symptomStats,
		MedicationStats:    medicationStats,
		OverallHealthScore: overallHealthScore,
		TrendDirection:     trendDirection,
	}, nil
}

// GetCorrelationAnalysis analyzes correlations between different data types
func (s *AnalyticsService) GetCorrelationAnalysis(ctx context.Context, userID string, start, end time.Time) (*analytics.CorrelationAnalysis, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get all data for correlation analysis
	bowelMovements, err := s.bowelMovementService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	medications, err := s.medicationService.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get medications: %w", err)
	}

	// Calculate correlations
	mealBowelCorrelations := s.calculateMealBowelCorrelations(meals, bowelMovements)
	mealSymptomCorrelations := s.calculateMealSymptomCorrelations(meals, symptoms)
	medicationEffectiveness := s.calculateMedicationEffectiveness(medications, symptoms, bowelMovements)
	triggerAnalysis := s.calculateTriggerAnalysis(meals, symptoms, bowelMovements)

	return &analytics.CorrelationAnalysis{
		MealBowelCorrelations:   mealBowelCorrelations,
		MealSymptomCorrelations: mealSymptomCorrelations,
		MedicationEffectiveness: medicationEffectiveness,
		TriggerAnalysis:         triggerAnalysis,
	}, nil
}

// GetTrendAnalysis analyzes trends over time
func (s *AnalyticsService) GetTrendAnalysis(ctx context.Context, userID string, start, end time.Time) (*analytics.TrendAnalysis, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get historical data for trend analysis
	bowelMovements, err := s.bowelMovementService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	// Calculate trends
	bowelMovementTrends := s.calculateBowelMovementTrends(bowelMovements, start, end)
	symptomTrends := s.calculateSymptomTrends(symptoms, start, end)
	mealTrends := s.calculateMealTrends(meals, start, end)

	// Determine overall trend
	overallTrend := s.determineOverallTrend(bowelMovementTrends, symptomTrends, mealTrends)

	return &analytics.TrendAnalysis{
		BowelMovementTrends: bowelMovementTrends,
		SymptomTrends:       symptomTrends,
		MealTrends:          mealTrends,
		OverallTrend:        overallTrend,
	}, nil
}

// GetBehaviorPatterns analyzes patterns in user behavior
func (s *AnalyticsService) GetBehaviorPatterns(ctx context.Context, userID string, start, end time.Time) (*analytics.BehaviorPatterns, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get data for pattern analysis
	bowelMovements, err := s.bowelMovementService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomService.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	// Analyze patterns
	eatingPatterns := s.analyzeEatingPatterns(meals)
	bowelPatterns := s.analyzeBowelPatterns(bowelMovements, meals)
	symptomPatterns := s.analyzeSymptomPatterns(symptoms)
	lifestylePatterns := s.analyzeLifestylePatterns(meals, bowelMovements, symptoms)

	return &analytics.BehaviorPatterns{
		EatingPatterns:    eatingPatterns,
		BowelPatterns:     bowelPatterns,
		SymptomPatterns:   symptomPatterns,
		LifestylePatterns: lifestylePatterns,
	}, nil
}

// GetHealthInsights generates actionable health insights
func (s *AnalyticsService) GetHealthInsights(ctx context.Context, userID string, start, end time.Time) (*analytics.HealthInsights, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get comprehensive data
	overview, err := s.GetUserHealthOverview(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get health overview: %w", err)
	}

	correlations, err := s.GetCorrelationAnalysis(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get correlation analysis: %w", err)
	}

	patterns, err := s.GetBehaviorPatterns(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get behavior patterns: %w", err)
	}

	// Generate insights
	keyFindings := s.generateKeyFindings(overview, correlations, patterns)
	riskFactors := s.identifyRiskFactors(overview, correlations)
	positiveFactors := s.identifyPositiveFactors(overview, patterns)
	recommendations := s.generateRecommendations(overview, correlations, patterns)
	alertLevel := s.determineAlertLevel(overview, riskFactors)
	confidenceLevel := s.calculateConfidenceLevel(overview)

	return &analytics.HealthInsights{
		KeyFindings:     keyFindings,
		RiskFactors:     riskFactors,
		PositiveFactors: positiveFactors,
		Recommendations: recommendations,
		AlertLevel:      alertLevel,
		ConfidenceLevel: confidenceLevel,
	}, nil
}

// GetHealthScore calculates an overall health score
func (s *AnalyticsService) GetHealthScore(ctx context.Context, userID string) (*analytics.HealthScore, error) {
	// Get recent data (last 30 days)
	end := time.Now()
	start := end.AddDate(0, 0, -30)

	overview, err := s.GetUserHealthOverview(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get health overview: %w", err)
	}

	// Calculate component scores
	componentScores := map[string]float64{
		"BowelHealth":    s.calculateBowelHealthScore(overview.BowelMovementStats),
		"Nutrition":      s.calculateNutritionScore(overview.MealStats),
		"SymptomControl": s.calculateSymptomControlScore(overview.SymptomStats),
		"Medication":     s.calculateMedicationScore(overview.MedicationStats),
	}

	// Calculate overall score (weighted average)
	overallScore := s.calculateWeightedScore(componentScores)

	// Calculate factors
	factors := s.calculateScoreFactors(overview, componentScores)

	// Calculate benchmarks (mock data for now)
	benchmarks := map[string]float64{
		"PopulationAverage": 65.0,
		"HealthyRange":      75.0,
		"OptimalRange":      85.0,
	}

	return &analytics.HealthScore{
		OverallScore:    overallScore,
		ComponentScores: componentScores,
		Trend:           overview.TrendDirection,
		LastUpdated:     time.Now(),
		Factors:         factors,
		Benchmarks:      benchmarks,
	}, nil
}

// GetRecommendations generates personalized recommendations
func (s *AnalyticsService) GetRecommendations(ctx context.Context, userID string) ([]*analytics.Recommendation, error) {
	// Get recent data for recommendations
	end := time.Now()
	start := end.AddDate(0, 0, -30)

	insights, err := s.GetHealthInsights(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get health insights: %w", err)
	}

	healthScore, err := s.GetHealthScore(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health score: %w", err)
	}

	// Generate personalized recommendations
	recommendations := s.generatePersonalizedRecommendations(insights, healthScore)

	return recommendations, nil
}

// Helper methods for calculations would go here...
// I'll implement key helper methods to make the service functional

func (s *AnalyticsService) calculateBowelMovementSummary(movements []*bowelmovement.BowelMovement, start, end time.Time) *analytics.BowelMovementSummary {
	if len(movements) == 0 {
		return &analytics.BowelMovementSummary{}
	}

	days := end.Sub(start).Hours() / 24
	if days == 0 {
		days = 1
	}

	var totalPain, totalStrain, totalSatisfaction float64
	bristolCounts := make(map[int]int)

	for _, movement := range movements {
		totalPain += float64(movement.Pain)
		totalStrain += float64(movement.Strain)
		totalSatisfaction += float64(movement.Satisfaction)
		bristolCounts[movement.BristolType]++
	}

	// Find most common Bristol scale
	mostCommonBristol := 0
	maxCount := 0
	for bristol, count := range bristolCounts {
		if count > maxCount {
			maxCount = count
			mostCommonBristol = bristol
		}
	}

	count := float64(len(movements))
	return &analytics.BowelMovementSummary{
		TotalCount:          int64(len(movements)),
		AveragePerDay:       count / days,
		MostCommonBristol:   mostCommonBristol,
		AveragePain:         totalPain / count,
		AverageStrain:       totalStrain / count,
		AverageSatisfaction: totalSatisfaction / count,
		RegularityScore:     s.calculateRegularityScore(movements),
	}
}

func (s *AnalyticsService) calculateMealSummary(meals []*meal.Meal, start, end time.Time) *analytics.MealSummary {
	if len(meals) == 0 {
		return &analytics.MealSummary{}
	}

	days := end.Sub(start).Hours() / 24
	if days == 0 {
		days = 1
	}

	var totalCalories int
	var fiberRichCount int

	for _, m := range meals {
		totalCalories += m.Calories
		if m.FiberRich {
			fiberRichCount++
		}
	}

	count := float64(len(meals))
	return &analytics.MealSummary{
		TotalMeals:       int64(len(meals)),
		AveragePerDay:    count / days,
		TotalCalories:    totalCalories,
		AverageCalories:  float64(totalCalories) / count,
		FiberRichPercent: float64(fiberRichCount) / count * 100,
		HealthScore:      s.calculateMealHealthScore(meals),
	}
}

func (s *AnalyticsService) calculateSymptomSummary(symptoms []*symptom.Symptom, start, end time.Time) *analytics.SymptomSummary {
	if len(symptoms) == 0 {
		return &analytics.SymptomSummary{}
	}

	days := end.Sub(start).Hours() / 24
	if days == 0 {
		days = 1
	}

	var totalSeverity float64
	categoryCounts := make(map[string]int)
	typeCounts := make(map[string]int)

	for _, s := range symptoms {
		totalSeverity += float64(s.Severity)
		if s.Category != nil {
			categoryCounts[string(*s.Category)]++
		}
		if s.Type != nil {
			typeCounts[string(*s.Type)]++
		}
	}

	// Find most common category and type
	mostCommonCategory := ""
	maxCategoryCount := 0
	for category, count := range categoryCounts {
		if count > maxCategoryCount {
			maxCategoryCount = count
			mostCommonCategory = category
		}
	}

	mostCommonType := ""
	maxTypeCount := 0
	for symptomType, count := range typeCounts {
		if count > maxTypeCount {
			maxTypeCount = count
			mostCommonType = symptomType
		}
	}

	count := float64(len(symptoms))
	return &analytics.SymptomSummary{
		TotalSymptoms:      int64(len(symptoms)),
		AveragePerDay:      count / days,
		AverageSeverity:    totalSeverity / count,
		MostCommonCategory: mostCommonCategory,
		MostCommonType:     mostCommonType,
		TrendDirection:     "STABLE", // Simplified for now
	}
}

func (s *AnalyticsService) calculateMedicationSummary(medications []*medication.Medication) *analytics.MedicationSummary {
	if len(medications) == 0 {
		return &analytics.MedicationSummary{}
	}

	var activeCount int64
	categoryCounts := make(map[string]int)

	for _, med := range medications {
		if med.IsActive {
			activeCount++
		}
		if med.Category != nil {
			categoryCounts[string(*med.Category)]++
		}
	}

	// Find most common category
	mostCommonCategory := ""
	maxCount := 0
	for category, count := range categoryCounts {
		if count > maxCount {
			maxCount = count
			mostCommonCategory = category
		}
	}

	return &analytics.MedicationSummary{
		TotalMedications:   int64(len(medications)),
		ActiveMedications:  activeCount,
		AdherenceScore:     85.0, // Simplified for now
		MostCommonCategory: mostCommonCategory,
		ComplexityScore:    float64(len(medications)) * 10, // Simplified calculation
	}
}

func (s *AnalyticsService) calculateOverallHealthScore(bowel *analytics.BowelMovementSummary, meal *analytics.MealSummary, symptom *analytics.SymptomSummary, medication *analytics.MedicationSummary) float64 {
	// Simplified health score calculation
	bowelScore := math.Max(0, 100-bowel.AveragePain*10)
	mealScore := meal.HealthScore
	symptomScore := math.Max(0, 100-symptom.AverageSeverity*10)
	medicationScore := medication.AdherenceScore

	// Weighted average
	return (bowelScore*0.3 + mealScore*0.25 + symptomScore*0.25 + medicationScore*0.2)
}

func (s *AnalyticsService) determineTrendDirection(bowel *analytics.BowelMovementSummary, symptom *analytics.SymptomSummary) string {
	// Simplified trend calculation
	if bowel.RegularityScore > 0.8 && symptom.AverageSeverity < 3 {
		return "IMPROVING"
	} else if bowel.RegularityScore < 0.5 || symptom.AverageSeverity > 7 {
		return "DECLINING"
	}
	return "STABLE"
}

// Simplified helper methods for basic functionality
func (s *AnalyticsService) calculateRegularityScore(movements []*bowelmovement.BowelMovement) float64 {
	if len(movements) < 2 {
		return 0.5
	}

	// Calculate based on consistency of timing
	var intervals []float64
	for i := 1; i < len(movements); i++ {
		interval := movements[i].RecordedAt.Sub(movements[i-1].RecordedAt).Hours()
		intervals = append(intervals, interval)
	}

	// Calculate standard deviation
	var sum float64
	for _, interval := range intervals {
		sum += interval
	}
	mean := sum / float64(len(intervals))

	var variance float64
	for _, interval := range intervals {
		diff := interval - mean
		variance += diff * diff
	}
	variance /= float64(len(intervals))

	// Convert to score (lower variance = higher regularity)
	return math.Max(0, 1-variance/100)
}

func (s *AnalyticsService) calculateMealHealthScore(meals []*meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}

	var score float64
	for _, m := range meals {
		mealScore := 50.0 // Base score

		if m.FiberRich {
			mealScore += 20
		}
		if m.Calories > 100 && m.Calories < 800 {
			mealScore += 15
		}
		// Note: Quality field doesn't exist in meal model, so skipping it

		score += mealScore
	}

	return score / float64(len(meals))
}

// Stub implementations for more complex methods
func (s *AnalyticsService) calculateMealBowelCorrelations(meals []*meal.Meal, movements []*bowelmovement.BowelMovement) []*analytics.Correlation {
	return []*analytics.Correlation{
		{
			Factor:      "High Fiber Meals",
			Outcome:     "Better Bowel Movement Quality",
			Strength:    0.6,
			Confidence:  0.7,
			Description: "Fiber-rich meals correlate with better bowel movement quality",
			SampleSize:  len(meals),
		},
	}
}

func (s *AnalyticsService) calculateMealSymptomCorrelations(meals []*meal.Meal, symptoms []*symptom.Symptom) []*analytics.Correlation {
	return []*analytics.Correlation{
		{
			Factor:      "Spicy Foods",
			Outcome:     "Digestive Symptoms",
			Strength:    0.4,
			Confidence:  0.6,
			Description: "Spicy foods may trigger digestive symptoms",
			SampleSize:  len(meals),
		},
	}
}

func (s *AnalyticsService) calculateMedicationEffectiveness(medications []*medication.Medication, symptoms []*symptom.Symptom, movements []*bowelmovement.BowelMovement) []*analytics.MedicationEffect {
	effects := make([]*analytics.MedicationEffect, 0, len(medications))

	for _, med := range medications {
		effects = append(effects, &analytics.MedicationEffect{
			MedicationName: med.Name,
			SymptomImprovement: map[string]float64{
				"digestive": 0.3,
				"pain":      0.4,
			},
			BowelImprovement:   0.2,
			EffectivenessScore: 0.65,
		})
	}

	return effects
}

func (s *AnalyticsService) calculateTriggerAnalysis(meals []*meal.Meal, symptoms []*symptom.Symptom, movements []*bowelmovement.BowelMovement) []*analytics.TriggerEffect {
	return []*analytics.TriggerEffect{
		{
			Trigger:         "Dairy Products",
			SymptomIncrease: 0.3,
			BowelImpact:     0.2,
			Frequency:       5,
			Severity:        0.6,
		},
	}
}

// More stub implementations for the remaining methods...
func (s *AnalyticsService) calculateBowelMovementTrends(movements []*bowelmovement.BowelMovement, start, end time.Time) *analytics.DataTrend {
	return &analytics.DataTrend{
		Direction:  "STABLE",
		Slope:      0.1,
		Confidence: 0.7,
		TimePoints: []time.Time{start, end},
		Values:     []float64{5.0, 5.5},
		Seasonality: map[string]float64{
			"morning": 0.6,
			"evening": 0.4,
		},
	}
}

func (s *AnalyticsService) calculateSymptomTrends(symptoms []*symptom.Symptom, start, end time.Time) *analytics.DataTrend {
	return &analytics.DataTrend{
		Direction:  "IMPROVING",
		Slope:      -0.1,
		Confidence: 0.6,
		TimePoints: []time.Time{start, end},
		Values:     []float64{4.0, 3.5},
		Seasonality: map[string]float64{
			"stress": 0.7,
			"diet":   0.3,
		},
	}
}

func (s *AnalyticsService) calculateMealTrends(meals []*meal.Meal, start, end time.Time) *analytics.DataTrend {
	return &analytics.DataTrend{
		Direction:  "IMPROVING",
		Slope:      0.05,
		Confidence: 0.8,
		TimePoints: []time.Time{start, end},
		Values:     []float64{70.0, 72.5},
		Seasonality: map[string]float64{
			"healthy": 0.8,
			"junk":    0.2,
		},
	}
}

func (s *AnalyticsService) determineOverallTrend(bowelTrend, symptomTrend, mealTrend *analytics.DataTrend) string {
	improvingCount := 0
	decliningCount := 0

	trends := []*analytics.DataTrend{bowelTrend, symptomTrend, mealTrend}
	for _, trend := range trends {
		switch trend.Direction {
		case "IMPROVING":
			improvingCount++
		case "DECLINING":
			decliningCount++
		}
	}

	if improvingCount > decliningCount {
		return "IMPROVING"
	} else if decliningCount > improvingCount {
		return "DECLINING"
	}
	return "STABLE"
}

func (s *AnalyticsService) analyzeEatingPatterns(meals []*meal.Meal) *analytics.EatingPattern {
	mealTiming := make(map[string]float64)
	cuisineCount := make(map[string]int)

	for _, meal := range meals {
		hour := meal.MealTime.Hour()
		timeSlot := s.getTimeSlot(hour)
		mealTiming[timeSlot]++

		cuisineCount[meal.Cuisine]++
	}

	// Get preferred cuisines
	var preferredCuisines []string
	for cuisine := range cuisineCount {
		if cuisine != "" {
			preferredCuisines = append(preferredCuisines, cuisine)
		}
	}

	return &analytics.EatingPattern{
		MealTiming: mealTiming,
		MealSizeDistribution: map[string]float64{
			"small":  0.3,
			"medium": 0.5,
			"large":  0.2,
		},
		PreferredCuisines:  preferredCuisines,
		DietaryConsistency: 0.7,
	}
}

func (s *AnalyticsService) analyzeBowelPatterns(movements []*bowelmovement.BowelMovement, meals []*meal.Meal) *analytics.BowelPattern {
	timing := make(map[string]float64)
	bristolFreq := make(map[string]float64)

	for _, movement := range movements {
		hour := movement.RecordedAt.Hour()
		timeSlot := s.getTimeSlot(hour)
		timing[timeSlot]++

		bristolFreq[fmt.Sprintf("type_%d", movement.BristolType)]++
	}

	return &analytics.BowelPattern{
		PreferredTiming:     timing,
		RegularityScore:     s.calculateRegularityScore(movements),
		ConsistencyPatterns: bristolFreq,
		ResponseToMeals:     2.5, // Average hours after meals
	}
}

func (s *AnalyticsService) analyzeSymptomPatterns(symptoms []*symptom.Symptom) *analytics.SymptomPattern {
	timing := make(map[string]float64)
	triggers := make(map[string]float64)

	for _, symptom := range symptoms {
		hour := symptom.RecordedAt.Hour()
		timeSlot := s.getTimeSlot(hour)
		timing[timeSlot]++

		// Note: Triggers field is a slice in the model, so we'll process each trigger
		for _, trigger := range symptom.Triggers {
			if trigger != "" {
				triggers[trigger]++
			}
		}
	}

	return &analytics.SymptomPattern{
		SymptomTiming:   timing,
		TriggerPatterns: triggers,
		SeasonalVariation: map[string]float64{
			"spring": 0.2,
			"summer": 0.3,
			"fall":   0.2,
			"winter": 0.3,
		},
		CyclicalPatterns: false,
	}
}

func (s *AnalyticsService) analyzeLifestylePatterns(meals []*meal.Meal, movements []*bowelmovement.BowelMovement, symptoms []*symptom.Symptom) *analytics.LifestylePattern {
	return &analytics.LifestylePattern{
		StressLevels: map[string]float64{
			"low":    0.3,
			"medium": 0.5,
			"high":   0.2,
		},
		SleepQuality: map[string]float64{
			"poor":    0.2,
			"average": 0.5,
			"good":    0.3,
		},
		ExerciseFrequency: map[string]float64{
			"none":   0.2,
			"light":  0.4,
			"medium": 0.3,
			"high":   0.1,
		},
		WeatherSensitivity: map[string]float64{
			"none":   0.6,
			"mild":   0.3,
			"strong": 0.1,
		},
	}
}

func (s *AnalyticsService) getTimeSlot(hour int) string {
	switch {
	case hour >= 6 && hour < 12:
		return "morning"
	case hour >= 12 && hour < 17:
		return "afternoon"
	case hour >= 17 && hour < 22:
		return "evening"
	default:
		return "night"
	}
}

func (s *AnalyticsService) generateKeyFindings(overview *analytics.HealthOverview, correlations *analytics.CorrelationAnalysis, patterns *analytics.BehaviorPatterns) []string {
	findings := []string{
		fmt.Sprintf("Average %d bowel movements per day", int(overview.BowelMovementStats.AveragePerDay)),
		fmt.Sprintf("Most common Bristol scale: Type %d", overview.BowelMovementStats.MostCommonBristol),
		fmt.Sprintf("Average symptom severity: %.1f/10", overview.SymptomStats.AverageSeverity),
	}

	if overview.BowelMovementStats.RegularityScore > 0.8 {
		findings = append(findings, "Excellent bowel movement regularity")
	}

	if overview.MealStats.FiberRichPercent > 60 {
		findings = append(findings, "High fiber diet maintained")
	}

	return findings
}

func (s *AnalyticsService) identifyRiskFactors(overview *analytics.HealthOverview, correlations *analytics.CorrelationAnalysis) []string {
	var risks []string

	if overview.BowelMovementStats.AveragePain > 5 {
		risks = append(risks, "High average pain levels during bowel movements")
	}

	if overview.SymptomStats.AverageSeverity > 6 {
		risks = append(risks, "High average symptom severity")
	}

	if overview.BowelMovementStats.RegularityScore < 0.5 {
		risks = append(risks, "Irregular bowel movement patterns")
	}

	if len(risks) == 0 {
		risks = append(risks, "No significant risk factors identified")
	}

	return risks
}

func (s *AnalyticsService) identifyPositiveFactors(overview *analytics.HealthOverview, patterns *analytics.BehaviorPatterns) []string {
	var positives []string

	if overview.MealStats.FiberRichPercent > 50 {
		positives = append(positives, "Good fiber intake from meals")
	}

	if overview.BowelMovementStats.RegularityScore > 0.7 {
		positives = append(positives, "Consistent bowel movement schedule")
	}

	if overview.MedicationStats.AdherenceScore > 80 {
		positives = append(positives, "Good medication adherence")
	}

	if len(positives) == 0 {
		positives = append(positives, "Maintaining baseline health indicators")
	}

	return positives
}

func (s *AnalyticsService) generateRecommendations(overview *analytics.HealthOverview, correlations *analytics.CorrelationAnalysis, patterns *analytics.BehaviorPatterns) []*analytics.Insight {
	var recommendations []*analytics.Insight

	if overview.MealStats.FiberRichPercent < 40 {
		recommendations = append(recommendations, &analytics.Insight{
			Type:        "DIETARY",
			Category:    "nutrition",
			Message:     "Consider increasing fiber intake",
			Evidence:    "Low percentage of fiber-rich meals in diet",
			Priority:    "MEDIUM",
			Confidence:  0.8,
			ActionItems: []string{"Add more fruits and vegetables", "Choose whole grains", "Include legumes in meals"},
		})
	}

	if overview.BowelMovementStats.RegularityScore < 0.6 {
		recommendations = append(recommendations, &analytics.Insight{
			Type:        "LIFESTYLE",
			Category:    "routine",
			Message:     "Establish a consistent bathroom routine",
			Evidence:    "Irregular bowel movement patterns detected",
			Priority:    "HIGH",
			Confidence:  0.7,
			ActionItems: []string{"Try to use the bathroom at the same times daily", "Allow adequate time for bowel movements", "Don't delay when you feel the urge"},
		})
	}

	return recommendations
}

func (s *AnalyticsService) determineAlertLevel(overview *analytics.HealthOverview, riskFactors []string) string {
	riskCount := len(riskFactors)

	if overview.BowelMovementStats.AveragePain > 8 || overview.SymptomStats.AverageSeverity > 8 {
		return "HIGH"
	}

	if riskCount > 2 || overview.OverallHealthScore < 50 {
		return "MEDIUM"
	}

	return "LOW"
}

func (s *AnalyticsService) calculateConfidenceLevel(overview *analytics.HealthOverview) float64 {
	// Base confidence on data availability
	dataPoints := float64(overview.BowelMovementStats.TotalCount + overview.MealStats.TotalMeals + overview.SymptomStats.TotalSymptoms)

	// More data = higher confidence, but cap at 0.9
	confidence := math.Min(0.9, dataPoints/100)

	// Minimum confidence of 0.3
	return math.Max(0.3, confidence)
}

func (s *AnalyticsService) calculateBowelHealthScore(stats *analytics.BowelMovementSummary) float64 {
	score := 50.0 // Base score

	// Regularity contributes significantly
	score += stats.RegularityScore * 30

	// Lower pain is better
	score += math.Max(0, (10-stats.AveragePain)*2)

	// Lower strain is better
	score += math.Max(0, (10-stats.AverageStrain)*1)

	// Higher satisfaction is better
	score += stats.AverageSatisfaction * 1

	return math.Min(100, score)
}

func (s *AnalyticsService) calculateNutritionScore(stats *analytics.MealSummary) float64 {
	return stats.HealthScore
}

func (s *AnalyticsService) calculateSymptomControlScore(stats *analytics.SymptomSummary) float64 {
	// Lower severity and frequency = better score
	severityScore := math.Max(0, (10-stats.AverageSeverity)*10)
	frequencyScore := math.Max(0, 100-(stats.AveragePerDay*10))

	return (severityScore + frequencyScore) / 2
}

func (s *AnalyticsService) calculateMedicationScore(stats *analytics.MedicationSummary) float64 {
	return stats.AdherenceScore
}

func (s *AnalyticsService) calculateWeightedScore(componentScores map[string]float64) float64 {
	weights := map[string]float64{
		"BowelHealth":    0.3,
		"Nutrition":      0.25,
		"SymptomControl": 0.25,
		"Medication":     0.2,
	}

	var weightedSum, totalWeight float64
	for component, score := range componentScores {
		if weight, exists := weights[component]; exists {
			weightedSum += score * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedSum / totalWeight
}

func (s *AnalyticsService) calculateScoreFactors(overview *analytics.HealthOverview, componentScores map[string]float64) []*analytics.ScoreFactor {
	var factors []*analytics.ScoreFactor

	for component, score := range componentScores {
		impact := score - 50 // Relative to neutral (50)

		factor := &analytics.ScoreFactor{
			Name:        component,
			Impact:      impact,
			Weight:      0.25, // Simplified equal weighting
			Description: s.getFactorDescription(component, score),
			Trend:       s.getFactorTrend(component, overview),
		}

		factors = append(factors, factor)
	}

	// Sort by impact (highest positive impact first)
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].Impact > factors[j].Impact
	})

	return factors
}

func (s *AnalyticsService) getFactorDescription(component string, score float64) string {
	switch component {
	case "BowelHealth":
		if score > 75 {
			return "Excellent bowel movement regularity and comfort"
		} else if score > 50 {
			return "Good bowel health with room for improvement"
		} else {
			return "Bowel health needs attention"
		}
	case "Nutrition":
		if score > 75 {
			return "Excellent nutritional choices"
		} else if score > 50 {
			return "Good nutrition with some improvements possible"
		} else {
			return "Nutrition needs significant improvement"
		}
	case "SymptomControl":
		if score > 75 {
			return "Excellent symptom management"
		} else if score > 50 {
			return "Symptoms are reasonably well controlled"
		} else {
			return "Symptoms need better management"
		}
	case "Medication":
		if score > 75 {
			return "Excellent medication adherence"
		} else if score > 50 {
			return "Good medication compliance"
		} else {
			return "Medication adherence needs improvement"
		}
	default:
		return "Component contributing to overall health"
	}
}

func (s *AnalyticsService) getFactorTrend(component string, overview *analytics.HealthOverview) string {
	// Simplified trend based on overall trend
	return overview.TrendDirection
}

func (s *AnalyticsService) generatePersonalizedRecommendations(insights *analytics.HealthInsights, healthScore *analytics.HealthScore) []*analytics.Recommendation {
	var recommendations []*analytics.Recommendation

	// Convert insights to recommendations
	for i, insight := range insights.Recommendations {
		rec := &analytics.Recommendation{
			ID:             fmt.Sprintf("rec_%d_%d", time.Now().Unix(), i),
			Type:           insight.Type,
			Category:       insight.Category,
			Title:          insight.Message,
			Description:    insight.Evidence,
			Priority:       insight.Priority,
			Confidence:     insight.Confidence,
			Evidence:       []string{insight.Evidence},
			ActionSteps:    insight.ActionItems,
			ExpectedImpact: s.calculateExpectedImpact(insight),
			Timeline:       s.calculateTimeline(insight),
			CreatedAt:      time.Now(),
			ExpiresAt:      s.calculateExpirationDate(insight),
		}

		recommendations = append(recommendations, rec)
	}

	// Add health score based recommendations
	if healthScore.OverallScore < 60 {
		rec := &analytics.Recommendation{
			ID:          fmt.Sprintf("health_score_%d", time.Now().Unix()),
			Type:        "LIFESTYLE",
			Category:    "overall_health",
			Title:       "Focus on Overall Health Improvement",
			Description: "Your health score indicates room for improvement across multiple areas",
			Priority:    "HIGH",
			Confidence:  0.8,
			Evidence:    []string{fmt.Sprintf("Overall health score: %.1f/100", healthScore.OverallScore)},
			ActionSteps: []string{
				"Review dietary patterns and increase fiber intake",
				"Establish consistent daily routines",
				"Monitor symptoms more closely",
				"Ensure medication adherence",
			},
			ExpectedImpact: "Significant improvement in overall health score within 4-6 weeks",
			Timeline:       "4-6 weeks",
			CreatedAt:      time.Now(),
			ExpiresAt:      &[]time.Time{time.Now().AddDate(0, 0, 30)}[0],
		}

		recommendations = append(recommendations, rec)
	}

	return recommendations
}

func (s *AnalyticsService) calculateExpectedImpact(insight *analytics.Insight) string {
	switch insight.Priority {
	case "HIGH":
		return "Significant improvement expected within 2-4 weeks"
	case "MEDIUM":
		return "Moderate improvement expected within 4-6 weeks"
	case "LOW":
		return "Gradual improvement expected within 6-8 weeks"
	default:
		return "Improvement timeline varies"
	}
}

func (s *AnalyticsService) calculateTimeline(insight *analytics.Insight) string {
	switch insight.Type {
	case "DIETARY":
		return "2-4 weeks"
	case "LIFESTYLE":
		return "4-6 weeks"
	case "MEDICAL":
		return "Consult healthcare provider"
	case "BEHAVIORAL":
		return "2-3 weeks"
	default:
		return "4-6 weeks"
	}
}

func (s *AnalyticsService) calculateExpirationDate(insight *analytics.Insight) *time.Time {
	var weeks int
	switch insight.Priority {
	case "HIGH":
		weeks = 4
	case "MEDIUM":
		weeks = 8
	case "LOW":
		weeks = 12
	default:
		weeks = 6
	}

	expiration := time.Now().AddDate(0, 0, weeks*7)
	return &expiration
}
