package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/aggregator"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/analyzer"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/calculator"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/insights"
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
	healthScoreCalculator *calculator.HealthScoreCalculator
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
		healthScoreCalculator: calculator.NewHealthScoreCalculator(),
		correlationAnalyzer:   analyzer.NewCorrelationAnalyzer(),
		trendAnalyzer:         analyzer.NewTrendAnalyzer(),
		dataAggregator:        aggregator.NewDataAggregator(),
		insightEngine:         insights.NewInsightEngine(),

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
		OverallHealthScore: s.healthScoreCalculator.CalculateOverallHealthScore(bmStats, mStats, sStats, medStats),
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
	bowelMovementTrends := s.convertTrendLinesToDataTrend(bowelTrendLines, "bowel_movement")
	symptomTrends := s.convertTrendLinesToDataTrend(symptomTrendLines, "symptoms")
	mealTrends := s.convertTrendLinesToDataTrend(mealTrendLines, "meals")

	// Determine overall trend direction
	overallTrend := s.determineOverallTrendDirection(bowelMovementTrends, symptomTrends, mealTrends)

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
func (s *AnalyticsService) analyzeEatingPatterns(meals []*meal.Meal) *shared.EatingPattern {
	mealValues := s.convertMealPointers(meals)
	return s.trendAnalyzer.AnalyzeEatingPatterns(mealValues)
}

func (s *AnalyticsService) analyzeBowelPatterns(movements []*bowelmovement.BowelMovement, meals []*meal.Meal) *shared.BowelPattern {
	bowelValues := s.convertBowelMovementPointers(movements)
	mealValues := s.convertMealPointers(meals)
	return s.trendAnalyzer.AnalyzeBowelPatterns(bowelValues, mealValues)
}

func (s *AnalyticsService) analyzeSymptomPatterns(symptoms []*symptom.Symptom) *shared.SymptomPattern {
	symptomValues := s.convertSymptomPointers(symptoms)
	return s.trendAnalyzer.AnalyzeSymptomPatterns(symptomValues)
}

func (s *AnalyticsService) analyzeLifestylePatterns(meals []*meal.Meal, movements []*bowelmovement.BowelMovement, symptoms []*symptom.Symptom) *shared.LifestylePattern {
	bowelValues := s.convertBowelMovementPointers(movements)
	mealValues := s.convertMealPointers(meals)
	symptomValues := s.convertSymptomPointers(symptoms)
	return s.trendAnalyzer.AnalyzeLifestylePatterns(mealValues, bowelValues, symptomValues)
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
	// Use the insight engine to generate recommendations - simplified call for now
	recommendations := s.insightEngine.GenerateRecommendations([]*shared.InsightRecommendation{})
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
