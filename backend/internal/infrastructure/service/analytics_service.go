package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/aggregator"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/analyzer"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/insights"
	analyticsSvc "github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
	"go.uber.org/zap"
)

// AnalyticsError represents an error in the analytics service
type AnalyticsError string

func (e AnalyticsError) Error() string { return string(e) }

const (
	ErrInvalidDateRange AnalyticsError = "invalid date range: start date must be before end date"
)

// AnalyticsServiceConfig holds configuration values
type AnalyticsServiceConfig struct {
	DefaultMedicationLimit int
	DefaultDataWindow      time.Duration
}

// AnalyticsService handles all analytics operations by delegating to specialized submodules
type AnalyticsService struct {
	// Configuration
	config *AnalyticsServiceConfig

	// Core submodules
	HealthScoreCalculator *analytics.HealthScore
	correlationAnalyzer   *analyzer.CorrelationAnalyzer
	trendAnalyzer         *analyzer.TrendAnalyzer
	dataAggregator        *aggregator.DataAggregator
	insightEngine         *insights.InsightEngine

	// Repositories and Services
	bowelMovementRepo bowelmovement.Repository
	mealRepo          meal.Repository
	symptomRepo       symptom.Repository
	medicationService medication.Service
	logger            *zap.Logger
}

// NewAnalyticsService creates a new analytics service with all required submodules
func NewAnalyticsService(
	bowelMovementRepo bowelmovement.Repository,
	mealRepo meal.Repository,
	symptomRepo symptom.Repository,
	medicationService medication.Service,
	logger *zap.Logger,
	config *AnalyticsServiceConfig,
) *AnalyticsService {
	if config == nil {
		config = &AnalyticsServiceConfig{
			DefaultMedicationLimit: 100,
			DefaultDataWindow:      30 * 24 * time.Hour, // 30 days
		}
	}

	return &AnalyticsService{
		// Configuration
		config: config,

		// Initialize submodules
		correlationAnalyzer: analyzer.NewCorrelationAnalyzer(),
		trendAnalyzer:       analyzer.NewTrendAnalyzer(),
		dataAggregator:      aggregator.NewDataAggregator(),
		insightEngine:       insights.NewInsightEngine(),

		// Store dependencies
		bowelMovementRepo: bowelMovementRepo,
		mealRepo:          mealRepo,
		symptomRepo:       symptomRepo,
		medicationService: medicationService,
		logger:            logger,
	}
}

// GetUserHealthOverview generates a comprehensive health overview for a user
func (s *AnalyticsService) GetUserHealthOverview(ctx context.Context, userID string, start, end time.Time) (*analytics.HealthOverview, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get all data for analysis
	bowelMovements, err := s.bowelMovementRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	medications, err := s.medicationService.GetByUserID(ctx, userID, s.config.DefaultMedicationLimit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get medications: %w", err)
	}

	// Convert pointer slices to value slices for analysis
	bowelValues := s.convertBowelMovementPointers(bowelMovements)
	mealValues := s.convertMealPointers(meals)
	symptomValues := s.convertSymptomPointers(symptoms)

	// Calculate statistics using the aggregate components
	bmStats := &analytics.BowelMovementSummary{
		TotalCount:          int64(len(bowelValues)),
		AveragePerDay:       s.calculateAveragePerDay(len(bowelValues), start, end),
		MostCommonBristol:   s.calculateMostCommonBristol(bowelValues),
		AveragePain:         s.calculateAveragePain(bowelValues),
		AverageStrain:       s.calculateAverageStrain(bowelValues),
		AverageSatisfaction: s.calculateAverageSatisfaction(bowelValues),
		RegularityScore:     s.calculateRegularityScore(bowelValues),
	}

	mStats := &analytics.MealSummary{
		TotalMeals:       int64(len(mealValues)),
		AveragePerDay:    s.calculateAveragePerDay(len(mealValues), start, end),
		TotalCalories:    s.calculateTotalCalories(mealValues),
		AverageCalories:  s.calculateAverageCalories(mealValues),
		FiberRichPercent: s.calculateFiberRichPercent(mealValues),
		HealthScore:      s.calculateMealHealthScore(mealValues),
	}

	sStats := &analytics.SymptomSummary{
		TotalSymptoms:      int64(len(symptomValues)),
		AveragePerDay:      s.calculateAveragePerDay(len(symptomValues), start, end),
		AverageSeverity:    s.calculateAverageSeverity(symptomValues),
		MostCommonCategory: s.findMostCommonCategory(symptomValues),
		MostCommonType:     s.findMostCommonType(symptomValues),
		TrendDirection:     s.calculateSymptomTrendDirection(symptomValues),
	}

	medStats := &analytics.MedicationSummary{
		TotalMedications:   int64(len(medications)),
		ActiveMedications:  s.countActiveMedications(medications),
		AdherenceScore:     s.calculateAdherenceScore(medications),
		MostCommonCategory: s.findMostCommonMedicationCategory(medications),
		ComplexityScore:    s.calculateMedicationComplexity(medications),
	}

	overview := &analytics.HealthOverview{
		Period:             fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
		BowelMovementStats: bmStats,
		MealStats:          mStats,
		SymptomStats:       sStats,
		MedicationStats:    medStats,
		TrendDirection:     s.calculateTrendDirection(bowelValues, mealValues, symptomValues),
	}

	return overview, nil
}

// GetCorrelationAnalysis analyzes correlations between different data types
func (s *AnalyticsService) GetCorrelationAnalysis(ctx context.Context, userID string, start, end time.Time) (*analytics.CorrelationAnalysis, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get all data for correlation analysis
	bowelMovements, err := s.bowelMovementRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	medications, err := s.medicationService.GetByUserID(ctx, userID, s.config.DefaultMedicationLimit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get medications: %w", err)
	}

	// Convert pointer slices to value slices for analyzers
	bowelMovementValues := s.convertBowelMovementPointers(bowelMovements)
	mealValues := s.convertMealPointers(meals)
	symptomValues := s.convertSymptomPointers(symptoms)
	medicationValues := s.convertMedicationPointers(medications)

	// Calculate correlations using the correlation analyzer
	mealBowelCorrelations := s.correlationAnalyzer.CalculateMealBowelCorrelations(mealValues, bowelMovementValues)
	mealSymptomCorrelations := s.correlationAnalyzer.CalculateMealSymptomCorrelations(mealValues, symptomValues)
	medicationEffectiveness := s.correlationAnalyzer.CalculateMedicationEffectiveness(medicationValues, symptomValues, bowelMovementValues)

	// TODO: Implement trigger analysis in analyzer component
	triggerAnalysis := []*analytics.TriggerEffect{}

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
	bowelMovements, err := s.bowelMovementRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get symptoms: %w", err)
	}

	// Convert pointer slices to value slices for analyzers
	bowelMovementValues := s.convertBowelMovementPointers(bowelMovements)
	mealValues := s.convertMealPointers(meals)
	symptomValues := s.convertSymptomPointers(symptoms)

	// Calculate trends using the trend analyzer
	bowelTrendLines := s.trendAnalyzer.CalculateBowelMovementTrends(bowelMovementValues, start, end)
	symptomTrendLines := s.trendAnalyzer.CalculateSymptomTrends(symptomValues, start, end)
	mealTrendLines := s.trendAnalyzer.CalculateMealTrends(mealValues, start, end)

	// Convert trend lines to DataTrend format
	bowelMovementTrend := s.convertTrendLinesToDataTrend(bowelTrendLines, "bowel_movement")
	symptomTrend := s.convertTrendLinesToDataTrend(symptomTrendLines, "symptoms")
	mealTrend := s.convertTrendLinesToDataTrend(mealTrendLines, "meals")

	// Determine overall trend direction
	overallTrend := s.determineOverallTrendDirection(bowelMovementTrend, symptomTrend, mealTrend)

	return &analytics.TrendAnalysis{
		BowelMovementTrends: bowelMovementTrend,
		SymptomTrends:       symptomTrend,
		MealTrends:          mealTrend,
		OverallTrend:        overallTrend,
	}, nil
}

// GetBehaviorPatterns analyzes patterns in user behavior
func (s *AnalyticsService) GetBehaviorPatterns(ctx context.Context, userID string, start, end time.Time) (*analytics.BehaviorPatterns, error) {
	if start.After(end) {
		return nil, analytics.ErrInvalidDateRange
	}

	// Get data for pattern analysis
	bowelMovements, err := s.bowelMovementRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get bowel movements: %w", err)
	}

	meals, err := s.mealRepo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	symptoms, err := s.symptomRepo.GetByDateRange(ctx, userID, start, end)
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

// TODO: These are stub implementations - replace with proper analytics components
// Pattern analysis methods delegate to specialized analyzers
func (s *AnalyticsService) analyzeEatingPatterns(meals []*meal.Meal) *analytics.EatingPattern {
	mealValues := s.convertMealPointers(meals)
	sharedPattern := s.trendAnalyzer.AnalyzeEatingPatterns(mealValues)
	return s.convertSharedEatingPatternToDomain(sharedPattern)
}

func (s *AnalyticsService) analyzeBowelPatterns(movements []*bowelmovement.BowelMovement, meals []*meal.Meal) *analytics.BowelPattern {
	bowelValues := s.convertBowelMovementPointers(movements)
	mealValues := s.convertMealPointers(meals)
	sharedPattern := s.trendAnalyzer.AnalyzeBowelPatterns(bowelValues, mealValues)
	return s.convertSharedBowelPatternToDomain(sharedPattern)
}

func (s *AnalyticsService) analyzeSymptomPatterns(symptoms []*symptom.Symptom) *analytics.SymptomPattern {
	symptomValues := s.convertSymptomPointers(symptoms)
	sharedPattern := s.trendAnalyzer.AnalyzeSymptomPatterns(symptomValues)
	return s.convertSharedSymptomPatternToDomain(sharedPattern)
}

func (s *AnalyticsService) analyzeLifestylePatterns(meals []*meal.Meal, movements []*bowelmovement.BowelMovement, symptoms []*symptom.Symptom) *analytics.LifestylePattern {
	bowelValues := s.convertBowelMovementPointers(movements)
	mealValues := s.convertMealPointers(meals)
	symptomValues := s.convertSymptomPointers(symptoms)
	sharedPattern := s.trendAnalyzer.AnalyzeLifestylePatterns(mealValues, bowelValues, symptomValues)
	return s.convertSharedLifestylePatternToDomain(sharedPattern)
}

// Pattern conversion methods

func (s *AnalyticsService) convertSharedEatingPatternToDomain(sharedPattern *analyticsSvc.EatingPattern) *analytics.EatingPattern {
	if sharedPattern == nil {
		return &analytics.EatingPattern{
			MealTiming:           make(map[string]float64),
			MealSizeDistribution: make(map[string]float64),
			PreferredCuisines:    []string{},
			DietaryConsistency:   0,
		}
	}

	// Convert the shared pattern to domain pattern
	// Since structures are different, we need to do a manual conversion
	mealTiming := make(map[string]float64)
	for _, timing := range sharedPattern.MealTimings {
		hourKey := fmt.Sprintf("%d", timing.TimeOfDay.Hour)
		mealTiming[hourKey] = float64(timing.Frequency)
	}

	return &analytics.EatingPattern{
		MealTiming:           mealTiming,
		MealSizeDistribution: make(map[string]float64),        // Not available in shared
		PreferredCuisines:    sharedPattern.CommonIngredients, // Use common ingredients as proxy
		DietaryConsistency:   0.8,                             // Default value
	}
}

func (s *AnalyticsService) convertSharedBowelPatternToDomain(sharedPattern *analyticsSvc.BowelPattern) *analytics.BowelPattern {
	if sharedPattern == nil {
		return &analytics.BowelPattern{
			PreferredTiming:     make(map[string]float64),
			RegularityScore:     0,
			ConsistencyPatterns: make(map[string]float64),
			ResponseToMeals:     0,
		}
	}

	return &analytics.BowelPattern{
		PreferredTiming:     make(map[string]float64), // Not directly available in shared
		RegularityScore:     sharedPattern.Frequency,  // Use frequency as proxy
		ConsistencyPatterns: map[string]float64{"consistency": sharedPattern.Consistency},
		ResponseToMeals:     sharedPattern.MealCorrelation * 24, // Convert correlation to hours
	}
}

func (s *AnalyticsService) convertSharedSymptomPatternToDomain(sharedPattern *analyticsSvc.SymptomPattern) *analytics.SymptomPattern {
	if sharedPattern == nil {
		return &analytics.SymptomPattern{
			SymptomTiming:     make(map[string]float64),
			TriggerPatterns:   make(map[string]float64),
			SeasonalVariation: make(map[string]float64),
			CyclicalPatterns:  false,
		}
	}

	// Convert common symptoms map to timing pattern
	timingMap := make(map[string]float64)
	for symptom, count := range sharedPattern.CommonSymptoms {
		timingMap[symptom] = float64(count)
	}

	// Convert frequency to trigger patterns
	triggerMap := make(map[string]float64)
	for trigger, freq := range sharedPattern.Frequency {
		triggerMap[trigger] = float64(freq)
	}

	return &analytics.SymptomPattern{
		SymptomTiming:     timingMap,
		TriggerPatterns:   triggerMap,
		SeasonalVariation: make(map[string]float64),              // Not available in shared
		CyclicalPatterns:  len(sharedPattern.CommonSymptoms) > 2, // Simple heuristic
	}
}

func (s *AnalyticsService) convertSharedLifestylePatternToDomain(sharedPattern *analyticsSvc.LifestylePattern) *analytics.LifestylePattern {
	if sharedPattern == nil {
		return &analytics.LifestylePattern{
			StressLevels:       make(map[string]float64),
			SleepQuality:       make(map[string]float64),
			ExerciseFrequency:  make(map[string]float64),
			WeatherSensitivity: make(map[string]float64),
		}
	}

	return &analytics.LifestylePattern{
		StressLevels:       make(map[string]float64), // Not directly available in shared
		SleepQuality:       make(map[string]float64), // Not directly available in shared
		ExerciseFrequency:  make(map[string]float64), // Not directly available in shared
		WeatherSensitivity: make(map[string]float64), // Not directly available in shared
	}
}

func (s *AnalyticsService) generateKeyFindings(overview *analytics.HealthOverview, correlations *analytics.CorrelationAnalysis, patterns *analytics.BehaviorPatterns) []string {
	return []string{"Analysis completed successfully"}
}

func (s *AnalyticsService) identifyRiskFactors(overview *analytics.HealthOverview, correlations *analytics.CorrelationAnalysis) []string {
	return []string{}
}

func (s *AnalyticsService) identifyPositiveFactors(overview *analytics.HealthOverview, patterns *analytics.BehaviorPatterns) []string {
	return []string{}
}

func (s *AnalyticsService) generateRecommendations(overview *analytics.HealthOverview, correlations *analytics.CorrelationAnalysis, patterns *analytics.BehaviorPatterns) []*analytics.Insight {
	return []*analytics.Insight{
		{
			Type:       "LIFESTYLE",
			Category:   "General",
			Message:    "Continue monitoring your health data",
			Evidence:   "Based on current analysis",
			Priority:   "LOW",
			Confidence: 0.8,
		},
	}
}

func (s *AnalyticsService) determineAlertLevel(overview *analytics.HealthOverview, riskFactors []string) string {
	return "NORMAL"
}

func (s *AnalyticsService) calculateConfidenceLevel(overview *analytics.HealthOverview) float64 {
	return 0.8
}

func (s *AnalyticsService) calculateBowelHealthScore(stats *analytics.BowelMovementSummary) float64 {
	return 75.0
}

func (s *AnalyticsService) calculateNutritionScore(stats *analytics.MealSummary) float64 {
	return 70.0
}

func (s *AnalyticsService) calculateSymptomControlScore(stats *analytics.SymptomSummary) float64 {
	return 80.0
}

func (s *AnalyticsService) calculateMedicationScore(stats *analytics.MedicationSummary) float64 {
	return 85.0
}

func (s *AnalyticsService) calculateWeightedScore(componentScores map[string]float64) float64 {
	total := 0.0
	count := 0
	for _, score := range componentScores {
		total += score
		count++
	}
	if count > 0 {
		return total / float64(count)
	}
	return 0.0
}

func (s *AnalyticsService) calculateScoreFactors(overview *analytics.HealthOverview, scores map[string]float64) []*analytics.ScoreFactor {
	return []*analytics.ScoreFactor{}
}

func (s *AnalyticsService) generatePersonalizedRecommendations(insights *analytics.HealthInsights, healthScore *analytics.HealthScore) []*analytics.Recommendation {
	// For now, return a simple set of recommendations based on the insights
	recommendations := make([]*analytics.Recommendation, 0)

	// Add recommendations based on alert level
	if insights.AlertLevel == "HIGH" {
		recommendations = append(recommendations, &analytics.Recommendation{
			ID:             "high-alert-1",
			Type:           "MEDICAL",
			Category:       "Healthcare",
			Title:          "Consult Healthcare Provider",
			Description:    "Your health indicators suggest you should consult with a healthcare provider",
			Priority:       "HIGH",
			Confidence:     0.9,
			Evidence:       insights.RiskFactors,
			ActionSteps:    []string{"Schedule appointment", "Prepare health summary"},
			ExpectedImpact: "Improved health management and early intervention",
			Timeline:       "Within 1-2 weeks",
			CreatedAt:      time.Now(),
		})
	}

	// Add general recommendations based on score
	if healthScore.OverallScore < 50 {
		recommendations = append(recommendations, &analytics.Recommendation{
			ID:             "low-score-1",
			Type:           "LIFESTYLE",
			Category:       "General Health",
			Title:          "Improve Health Habits",
			Description:    "Your overall health score could be improved with lifestyle changes",
			Priority:       "MEDIUM",
			Confidence:     0.7,
			Evidence:       []string{fmt.Sprintf("Overall score: %.1f", healthScore.OverallScore)},
			ActionSteps:    []string{"Review dietary patterns", "Monitor symptoms closely"},
			ExpectedImpact: "Gradual improvement in overall health score",
			Timeline:       "2-4 weeks",
			CreatedAt:      time.Now(),
		})
	}

	return recommendations
}

// Type conversion helpers between pointer and value types

func (s *AnalyticsService) convertBowelMovementPointers(bms []*bowelmovement.BowelMovement) []bowelmovement.BowelMovement {
	result := make([]bowelmovement.BowelMovement, len(bms))
	for i, bm := range bms {
		if bm != nil {
			result[i] = *bm
		}
	}
	return result
}

func (s *AnalyticsService) convertMealPointers(meals []*meal.Meal) []meal.Meal {
	result := make([]meal.Meal, len(meals))
	for i, m := range meals {
		if m != nil {
			result[i] = *m
		}
	}
	return result
}

func (s *AnalyticsService) convertSymptomPointers(symptoms []*symptom.Symptom) []symptom.Symptom {
	result := make([]symptom.Symptom, len(symptoms))
	for i, s := range symptoms {
		if s != nil {
			result[i] = *s
		}
	}
	return result
}

func (s *AnalyticsService) convertMedicationPointers(medications []*medication.Medication) []medication.Medication {
	result := make([]medication.Medication, len(medications))
	for i, m := range medications {
		if m != nil {
			result[i] = *m
		}
	}
	return result
}

// Calculation helper methods

func (s *AnalyticsService) calculateAveragePerDay(count int, start, end time.Time) float64 {
	days := end.Sub(start).Hours() / 24
	if days <= 0 {
		return 0
	}
	return float64(count) / days
}

func (s *AnalyticsService) calculateMostCommonBristol(bms []bowelmovement.BowelMovement) int {
	if len(bms) == 0 {
		return 0
	}
	bristolCounts := make(map[int]int)
	maxCount, maxBristol := 0, 0
	for _, bm := range bms {
		bristolCounts[bm.BristolType]++
		if bristolCounts[bm.BristolType] > maxCount {
			maxCount = bristolCounts[bm.BristolType]
			maxBristol = bm.BristolType
		}
	}
	return maxBristol
}

func (s *AnalyticsService) calculateAveragePain(bms []bowelmovement.BowelMovement) float64 {
	if len(bms) == 0 {
		return 0
	}
	sum := 0.0
	for _, bm := range bms {
		sum += float64(bm.Pain)
	}
	return sum / float64(len(bms))
}

func (s *AnalyticsService) calculateAverageStrain(bms []bowelmovement.BowelMovement) float64 {
	if len(bms) == 0 {
		return 0
	}
	sum := 0.0
	for _, bm := range bms {
		sum += float64(bm.Strain)
	}
	return sum / float64(len(bms))
}

func (s *AnalyticsService) calculateAverageSatisfaction(bms []bowelmovement.BowelMovement) float64 {
	if len(bms) == 0 {
		return 0
	}
	sum := 0.0
	for _, bm := range bms {
		sum += float64(bm.Satisfaction)
	}
	return sum / float64(len(bms))
}

func (s *AnalyticsService) calculateRegularityScore(bms []bowelmovement.BowelMovement) float64 {
	if len(bms) < 2 {
		return 0
	}
	// Calculate variance in time between movements
	var totalVariance float64
	var lastTime time.Time
	for i, bm := range bms {
		if i == 0 {
			lastTime = bm.RecordedAt
			continue
		}
		interval := bm.RecordedAt.Sub(lastTime).Hours()
		// Ideal interval is 24 hours, calculate variance from that
		variance := (interval - 24) * (interval - 24)
		totalVariance += variance
		lastTime = bm.RecordedAt
	}
	avgVariance := totalVariance / float64(len(bms)-1)
	// Convert to a 0-1 score where lower variance = higher score
	maxVarianceHours := 24.0 // Max reasonable variance is 24 hours
	score := 1 - (avgVariance / (maxVarianceHours * maxVarianceHours))
	if score < 0 {
		score = 0
	}
	return score
}

func (s *AnalyticsService) calculateTotalCalories(meals []meal.Meal) int {
	total := 0
	for _, m := range meals {
		total += m.Calories
	}
	return total
}

func (s *AnalyticsService) calculateAverageCalories(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}
	total := s.calculateTotalCalories(meals)
	return float64(total) / float64(len(meals))
}

func (s *AnalyticsService) calculateFiberRichPercent(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}
	fiberRich := 0
	for _, m := range meals {
		if m.FiberRich {
			fiberRich++
		}
	}
	return float64(fiberRich) / float64(len(meals)) * 100
}

// Updated meal methods to use actual model fields

func (s *AnalyticsService) calculateMealHealthScore(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}
	var totalScore float64
	for _, m := range meals {
		score := 0.0
		// Add points for desired attributes
		if m.FiberRich {
			score += 0.2
		}
		if !m.Dairy && !m.Gluten { // Lower risk of triggers
			score += 0.3
		}
		if m.SpicyLevel != nil && *m.SpicyLevel <= 5 { // Moderate spice level
			score += 0.2
		}
		// Add calories score (1500-2500 per meal is reasonable range)
		if m.Calories > 0 {
			calScore := s.calculateCalorieHealthScore(m.Calories)
			score += calScore * 0.3
		}
		totalScore += score
	}
	return totalScore / float64(len(meals))
}

func (s *AnalyticsService) calculateCalorieHealthScore(calories int) float64 {
	// Assume 1500-2500 is ideal range for a meal
	if calories <= 0 {
		return 0
	}
	if calories < 1500 {
		// Score linearly from 0 to 1 as calories approach 1500
		return float64(calories) / 1500
	}
	if calories <= 2500 {
		// Perfect score in ideal range
		return 1.0
	}
	// Score decreases linearly above 2500
	excess := float64(calories - 2500)
	score := 1.0 - (excess / 1000) // Lose 0.1 points per 100 calories over
	if score < 0 {
		score = 0
	}
	return score
}

func (s *AnalyticsService) hasCommonTriggers(m meal.Meal) bool {
	if m.Dairy || m.Gluten {
		return true
	}
	if m.SpicyLevel != nil && *m.SpicyLevel > 7 {
		return true
	}
	return false
}

func (s *AnalyticsService) hasGoodHydration(m meal.Meal) bool {
	// In the current model, we don't have detailed ingredient info,
	// so we'll base this primarily on the meal category
	if m.Category != nil {
		switch *m.Category {
		case "SOUP", "SALAD", "FRUIT", "SMOOTHIE", "BEVERAGE":
			return true
		}
	}
	return false
}

// Updated trend analysis methods

func (s *AnalyticsService) calculateAverageFiber(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}
	fiberCount := 0
	for _, m := range meals {
		if m.FiberRich {
			fiberCount++
		}
	}
	return float64(fiberCount) / float64(len(meals))
}

func (s *AnalyticsService) calculateMacroBalanceScore(meals []meal.Meal) float64 {
	// With the current model, we can't calculate true macro balance
	// Instead, use a simpler scoring based on available fields
	if len(meals) == 0 {
		return 0
	}

	var totalScore float64
	for _, m := range meals {
		score := 0.0
		if m.FiberRich {
			score += 0.3
		}
		if m.Category != nil {
			switch *m.Category {
			case "BALANCED", "PROTEIN", "WHOLE_GRAIN":
				score += 0.7
			case "VEGETABLE", "FRUIT":
				score += 0.5
			}
		}
		totalScore += score
	}
	return totalScore / float64(len(meals))
}

func (s *AnalyticsService) calculateTriggerScore(meals []meal.Meal) float64 {
	if len(meals) == 0 {
		return 0
	}
	triggerCount := 0
	for _, m := range meals {
		if s.hasCommonTriggers(m) {
			triggerCount++
		}
	}
	return float64(triggerCount) / float64(len(meals))
}

// Missing calculation methods for symptoms

func (s *AnalyticsService) calculateAverageSeverity(symptoms []symptom.Symptom) float64 {
	if len(symptoms) == 0 {
		return 0
	}
	total := 0.0
	for _, sym := range symptoms {
		total += float64(sym.Severity)
	}
	return total / float64(len(symptoms))
}

func (s *AnalyticsService) findMostCommonCategory(symptoms []symptom.Symptom) string {
	if len(symptoms) == 0 {
		return ""
	}
	categoryCounts := make(map[string]int)
	for _, sym := range symptoms {
		if sym.Category != nil {
			categoryCounts[string(*sym.Category)]++
		}
	}

	maxCount := 0
	mostCommonCategory := ""
	for category, count := range categoryCounts {
		if count > maxCount {
			maxCount = count
			mostCommonCategory = category
		}
	}
	return mostCommonCategory
}

func (s *AnalyticsService) findMostCommonType(symptoms []symptom.Symptom) string {
	if len(symptoms) == 0 {
		return ""
	}
	typeCounts := make(map[string]int)
	for _, sym := range symptoms {
		if sym.Type != nil {
			typeCounts[string(*sym.Type)]++
		}
	}

	maxCount := 0
	mostCommonType := ""
	for symType, count := range typeCounts {
		if count > maxCount {
			maxCount = count
			mostCommonType = symType
		}
	}
	return mostCommonType
}

func (s *AnalyticsService) calculateSymptomTrendDirection(symptoms []symptom.Symptom) string {
	if len(symptoms) < 2 {
		return "STABLE"
	}

	// Sort symptoms by date
	sortedSymptoms := make([]symptom.Symptom, len(symptoms))
	copy(sortedSymptoms, symptoms)
	sort.Slice(sortedSymptoms, func(i, j int) bool {
		return sortedSymptoms[i].RecordedAt.Before(sortedSymptoms[j].RecordedAt)
	})

	// Compare first half with second half
	midpoint := len(sortedSymptoms) / 2
	firstHalf := sortedSymptoms[:midpoint]
	secondHalf := sortedSymptoms[midpoint:]

	firstHalfAvg := s.calculateAverageSeverity(firstHalf)
	secondHalfAvg := s.calculateAverageSeverity(secondHalf)

	if secondHalfAvg > firstHalfAvg+0.5 {
		return "DECLINING" // Getting worse
	} else if firstHalfAvg > secondHalfAvg+0.5 {
		return "IMPROVING" // Getting better
	}
	return "STABLE"
}

// Missing calculation methods for medications

func (s *AnalyticsService) countActiveMedications(medications []*medication.Medication) int64 {
	count := int64(0)
	for _, med := range medications {
		if med != nil && med.IsActive {
			count++
		}
	}
	return count
}

func (s *AnalyticsService) calculateAdherenceScore(medications []*medication.Medication) float64 {
	if len(medications) == 0 {
		return 0
	}

	totalScore := 0.0
	for _, med := range medications {
		if med != nil && med.IsActive {
			// Simple adherence calculation - in reality this would be more complex
			// For now, assume 90% adherence if medication is active
			totalScore += 0.9
		}
	}
	return totalScore / float64(len(medications))
}

func (s *AnalyticsService) findMostCommonMedicationCategory(medications []*medication.Medication) string {
	if len(medications) == 0 {
		return ""
	}

	categoryCounts := make(map[string]int)
	for _, med := range medications {
		if med != nil && med.Category != nil {
			categoryCounts[string(*med.Category)]++
		}
	}

	maxCount := 0
	mostCommonCategory := ""
	for category, count := range categoryCounts {
		if count > maxCount {
			maxCount = count
			mostCommonCategory = category
		}
	}
	return mostCommonCategory
}

func (s *AnalyticsService) calculateMedicationComplexity(medications []*medication.Medication) float64 {
	if len(medications) == 0 {
		return 0
	}

	// Simple complexity score based on number of active medications
	activeCount := float64(s.countActiveMedications(medications))

	// Complexity increases with number of medications
	// Scale from 0-1 where 1 is low complexity, 0 is high complexity
	if activeCount <= 1 {
		return 1.0 // Low complexity
	} else if activeCount <= 3 {
		return 0.7 // Medium complexity
	} else if activeCount <= 5 {
		return 0.4 // High complexity
	}
	return 0.1 // Very high complexity
}

// Missing trend calculation methods

func (s *AnalyticsService) calculateTrendDirection(bowelMovements []bowelmovement.BowelMovement, meals []meal.Meal, symptoms []symptom.Symptom) string {
	// Simple trend calculation based on recent data
	bowelTrend := s.calculateBowelTrendDirection(bowelMovements)
	symptomTrend := s.calculateSymptomTrendDirection(symptoms)

	// Combine trends - if both are improving, overall is improving
	if bowelTrend == "IMPROVING" && symptomTrend == "IMPROVING" {
		return "IMPROVING"
	} else if bowelTrend == "DECLINING" || symptomTrend == "DECLINING" {
		return "DECLINING"
	}
	return "STABLE"
}

func (s *AnalyticsService) calculateBowelTrendDirection(bowelMovements []bowelmovement.BowelMovement) string {
	if len(bowelMovements) < 2 {
		return "STABLE"
	}

	// Sort movements by date
	sortedMovements := make([]bowelmovement.BowelMovement, len(bowelMovements))
	copy(sortedMovements, bowelMovements)
	sort.Slice(sortedMovements, func(i, j int) bool {
		return sortedMovements[i].RecordedAt.Before(sortedMovements[j].RecordedAt)
	})

	// Compare first half with second half regularity
	midpoint := len(sortedMovements) / 2
	firstHalf := sortedMovements[:midpoint]
	secondHalf := sortedMovements[midpoint:]

	firstHalfRegularity := s.calculateRegularityScore(firstHalf)
	secondHalfRegularity := s.calculateRegularityScore(secondHalf)

	if secondHalfRegularity > firstHalfRegularity+0.1 {
		return "IMPROVING"
	} else if firstHalfRegularity > secondHalfRegularity+0.1 {
		return "DECLINING"
	}
	return "STABLE"
}

// Missing trend conversion methods

func (s *AnalyticsService) convertTrendLinesToDataTrend(trendLines []*analyticsSvc.TrendLine, dataType string) *analytics.DataTrend {
	if len(trendLines) == 0 {
		return &analytics.DataTrend{
			Direction:   "STABLE",
			Slope:       0,
			Confidence:  0,
			TimePoints:  []time.Time{},
			Values:      []float64{},
			Seasonality: make(map[string]float64),
		}
	}

	// Aggregate multiple trend lines into a single DataTrend
	// For simplicity, take the first trend line as the primary trend
	firstTrend := trendLines[0]
	if firstTrend == nil {
		return &analytics.DataTrend{
			Direction:   "STABLE",
			Slope:       0,
			Confidence:  0,
			TimePoints:  []time.Time{},
			Values:      []float64{},
			Seasonality: make(map[string]float64),
		}
	}

	return &analytics.DataTrend{
		Direction:   firstTrend.Direction,
		Slope:       firstTrend.Slope,
		Confidence:  firstTrend.Confidence,
		TimePoints:  []time.Time{}, // TrendLine doesn't have TimePoints, so empty for now
		Values:      []float64{},   // TrendLine doesn't have Values, so empty for now
		Seasonality: make(map[string]float64),
	}
}

func (s *AnalyticsService) determineOverallTrendDirection(bowelTrend, symptomTrend, mealTrend *analytics.DataTrend) string {
	trends := []*analytics.DataTrend{bowelTrend, symptomTrend, mealTrend}

	improvingCount := 0
	decliningCount := 0
	totalTrends := 0

	for _, trend := range trends {
		if trend != nil {
			totalTrends++
			switch trend.Direction {
			case "IMPROVING":
				improvingCount++
			case "DECLINING":
				decliningCount++
			}
		}
	}

	if totalTrends == 0 {
		return "STABLE"
	}

	improvingRatio := float64(improvingCount) / float64(totalTrends)
	decliningRatio := float64(decliningCount) / float64(totalTrends)

	if improvingRatio > 0.6 {
		return "IMPROVING"
	} else if decliningRatio > 0.6 {
		return "DECLINING"
	}
	return "STABLE"
}
