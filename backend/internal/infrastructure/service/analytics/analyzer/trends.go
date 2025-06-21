
package analyzer

import (
	"fmt"
	"sort"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

// TrendAnalyzer analyzes trends over time
type TrendAnalyzer struct{}

// NewTrendAnalyzer creates a new trend analyzer
func NewTrendAnalyzer() *TrendAnalyzer {
	return &TrendAnalyzer{}
}

// CalculateBowelMovementTrends analyzes bowel movement trends over time
func (ta *TrendAnalyzer) CalculateBowelMovementTrends(
	movements []bowelmovement.BowelMovement,
	start, end time.Time,
) []*shared.TrendLine {
	trends := make([]*shared.TrendLine, 0)

	// Bristol type trend
	bristolTrend := ta.calculateBristolTypeTrend(movements, start, end)
	if bristolTrend != nil {
		trends = append(trends, bristolTrend)
	}

	// Pain trend
	painTrend := ta.calculatePainTrend(movements, start, end)
	if painTrend != nil {
		trends = append(trends, painTrend)
	}

	// Satisfaction trend
	satisfactionTrend := ta.calculateSatisfactionTrend(movements, start, end)
	if satisfactionTrend != nil {
		trends = append(trends, satisfactionTrend)
	}

	// Frequency trend
	frequencyTrend := ta.calculateFrequencyTrend(movements, start, end)
	if frequencyTrend != nil {
		trends = append(trends, frequencyTrend)
	}

	return trends
}

// CalculateSymptomTrends analyzes symptom trends over time
func (ta *TrendAnalyzer) CalculateSymptomTrends(
	symptoms []symptom.Symptom,
	start, end time.Time,
) []*shared.TrendLine {
	trends := make([]*shared.TrendLine, 0)

	// Symptom severity trend
	severityTrend := ta.calculateSymptomSeverityTrend(symptoms, start, end)
	if severityTrend != nil {
		trends = append(trends, severityTrend)
	}

	// Symptom frequency trend
	symptomFrequencyTrend := ta.calculateSymptomFrequencyTrend(symptoms, start, end)
	if symptomFrequencyTrend != nil {
		trends = append(trends, symptomFrequencyTrend)
	}

	return trends
}

// CalculateMealTrends analyzes meal trends over time
func (ta *TrendAnalyzer) CalculateMealTrends(
	meals []meal.Meal,
	start, end time.Time,
) []*shared.TrendLine {
	trends := make([]*shared.TrendLine, 0)

	// Calorie trend
	calorieTrend := ta.calculateCalorieTrend(meals, start, end)
	if calorieTrend != nil {
		trends = append(trends, calorieTrend)
	}

	// Meal frequency trend
	mealFrequencyTrend := ta.calculateMealFrequencyTrend(meals, start, end)
	if mealFrequencyTrend != nil {
		trends = append(trends, mealFrequencyTrend)
	}

	return trends
}

// DetermineOverallTrend determines the overall health trend direction
func (ta *TrendAnalyzer) DetermineOverallTrend(
	bowelTrends, symptomTrends, mealTrends []*shared.TrendLine,
) string {
	trendScores := make([]float64, 0)

	// Collect all trend slopes
	for _, trend := range bowelTrends {
		// For bowel movements, positive slope for satisfaction is good, negative for pain is good
		switch trend.Name {
		case "Satisfaction":
			trendScores = append(trendScores, trend.Slope)
		case "Pain", "Strain":
			trendScores = append(trendScores, -trend.Slope) // Invert negative trends
		default:
			trendScores = append(trendScores, trend.Slope)
		}
	}

	for _, trend := range symptomTrends {
		// For symptoms, negative slope (decreasing) is good
		trendScores = append(trendScores, -trend.Slope)
	}

	for _, trend := range mealTrends {
		// For meals, neutral to slightly positive trends are good
		trendScores = append(trendScores, trend.Slope)
	}

	if len(trendScores) == 0 {
		return "STABLE"
	}

	// Calculate average trend score
	totalScore := 0.0
	for _, score := range trendScores {
		totalScore += score
	}
	averageScore := totalScore / float64(len(trendScores))

	// Interpret trend direction
	if averageScore > 0.1 {
		return "IMPROVING"
	} else if averageScore < -0.1 {
		return "DECLINING"
	} else {
		return "STABLE"
	}
}

// Helper functions for bowel movement trends

func (ta *TrendAnalyzer) calculateBristolTypeTrend(
	movements []bowelmovement.BowelMovement,
	start, end time.Time,
) *shared.TrendLine {
	if len(movements) < 3 {
		return nil
	}

	// Group by week and calculate averages
	weeklyData := ta.groupBowelMovementsByWeek(movements, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: data.AverageBristol,
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.1)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Bristol Type",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

func (ta *TrendAnalyzer) calculatePainTrend(
	movements []bowelmovement.BowelMovement,
	start, end time.Time,
) *shared.TrendLine {
	if len(movements) < 3 {
		return nil
	}

	weeklyData := ta.groupBowelMovementsByWeek(movements, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: data.AveragePain,
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.1)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Pain",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

func (ta *TrendAnalyzer) calculateSatisfactionTrend(
	movements []bowelmovement.BowelMovement,
	start, end time.Time,
) *shared.TrendLine {
	if len(movements) < 3 {
		return nil
	}

	weeklyData := ta.groupBowelMovementsByWeek(movements, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: data.AverageSatisfaction,
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.1)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Satisfaction",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

func (ta *TrendAnalyzer) calculateFrequencyTrend(
	movements []bowelmovement.BowelMovement,
	start, end time.Time,
) *shared.TrendLine {
	if len(movements) < 3 {
		return nil
	}

	weeklyData := ta.groupBowelMovementsByWeek(movements, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: float64(data.Count),
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.5)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Frequency",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

// Helper functions for symptom trends

func (ta *TrendAnalyzer) calculateSymptomSeverityTrend(
	symptoms []symptom.Symptom,
	start, end time.Time,
) *shared.TrendLine {
	if len(symptoms) < 3 {
		return nil
	}

	weeklyData := ta.groupSymptomsByWeek(symptoms, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: data.AverageSeverity,
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.1)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Symptom Severity",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

func (ta *TrendAnalyzer) calculateSymptomFrequencyTrend(
	symptoms []symptom.Symptom,
	start, end time.Time,
) *shared.TrendLine {
	if len(symptoms) < 3 {
		return nil
	}

	weeklyData := ta.groupSymptomsByWeek(symptoms, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: float64(data.Count),
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.5)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Symptom Frequency",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

// Helper functions for meal trends

func (ta *TrendAnalyzer) calculateCalorieTrend(
	meals []meal.Meal,
	start, end time.Time,
) *shared.TrendLine {
	if len(meals) < 3 {
		return nil
	}

	weeklyData := ta.groupMealsByWeek(meals, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: data.AverageCalories,
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 50.0) // 50 calories threshold
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Average Calories",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

func (ta *TrendAnalyzer) calculateMealFrequencyTrend(
	meals []meal.Meal,
	start, end time.Time,
) *shared.TrendLine {
	if len(meals) < 3 {
		return nil
	}

	weeklyData := ta.groupMealsByWeek(meals, start, end)
	if len(weeklyData) < 2 {
		return nil
	}

	points := make([]shared.TrendPoint, 0, len(weeklyData))
	for _, data := range weeklyData {
		points = append(points, shared.TrendPoint{
			Date:  data.Date,
			Value: float64(data.Count),
		})
	}

	slope := shared.CalculateTrendSlope(points)
	direction := shared.InterpretTrendDirection(slope, 0.5)
	confidence := shared.CalculateConfidenceScore(len(points))

	return &shared.TrendLine{
		Name:         "Meal Frequency",
		Points:       points,
		Direction:    direction,
		Slope:        shared.RoundToDecimalPlaces(slope, 4),
		Confidence:   confidence,
		Significance: ta.interpretTrendSignificance(slope, confidence),
	}
}

// Data grouping helper types and functions

type WeeklyBowelMovementData struct {
	Date                time.Time
	Count               int
	AverageBristol      float64
	AveragePain         float64
	AverageSatisfaction float64
}

type WeeklySymptomData struct {
	Date            time.Time
	Count           int
	AverageSeverity float64
}

type WeeklyMealData struct {
	Date            time.Time
	Count           int
	AverageCalories float64
}

func (ta *TrendAnalyzer) groupBowelMovementsByWeek(
	movements []bowelmovement.BowelMovement,
	start, end time.Time,
) []WeeklyBowelMovementData {
	weekGroups := make(map[string][]bowelmovement.BowelMovement)

	for _, movement := range movements {
		if movement.RecordedAt.Before(start) || movement.RecordedAt.After(end) {
			continue
		}

		year, week := movement.RecordedAt.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		weekGroups[weekKey] = append(weekGroups[weekKey], movement)
	}

	result := make([]WeeklyBowelMovementData, 0, len(weekGroups))
	for _, movements := range weekGroups {
		if len(movements) == 0 {
			continue
		}

		var totalBristol, totalPain, totalSatisfaction float64
		for _, bm := range movements {
			totalBristol += float64(bm.BristolType)
			totalPain += float64(bm.Pain)
			totalSatisfaction += float64(bm.Satisfaction)
		}

		count := float64(len(movements))
		result = append(result, WeeklyBowelMovementData{
			Date:                movements[0].RecordedAt,
			Count:               len(movements),
			AverageBristol:      totalBristol / count,
			AveragePain:         totalPain / count,
			AverageSatisfaction: totalSatisfaction / count,
		})
	}

	// Sort results by date for consistent trend calculation
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result
}

func (ta *TrendAnalyzer) groupSymptomsByWeek(
	symptoms []symptom.Symptom,
	start, end time.Time,
) []WeeklySymptomData {
	weekGroups := make(map[string][]symptom.Symptom)

	for _, symptom := range symptoms {
		if symptom.RecordedAt.Before(start) || symptom.RecordedAt.After(end) {
			continue
		}

		year, week := symptom.RecordedAt.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		weekGroups[weekKey] = append(weekGroups[weekKey], symptom)
	}

	result := make([]WeeklySymptomData, 0, len(weekGroups))
	for _, symptoms := range weekGroups {
		if len(symptoms) == 0 {
			continue
		}

		var totalSeverity float64
		for _, s := range symptoms {
			totalSeverity += float64(s.Severity)
		}

		result = append(result, WeeklySymptomData{
			Date:            symptoms[0].RecordedAt,
			Count:           len(symptoms),
			AverageSeverity: totalSeverity / float64(len(symptoms)),
		})
	}

	// Sort results by date for consistent trend calculation
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result
}

func (ta *TrendAnalyzer) groupMealsByWeek(
	meals []meal.Meal,
	start, end time.Time,
) []WeeklyMealData {
	weekGroups := make(map[string][]meal.Meal)

	for _, meal := range meals {
		if meal.MealTime.Before(start) || meal.MealTime.After(end) {
			continue
		}

		year, week := meal.MealTime.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		weekGroups[weekKey] = append(weekGroups[weekKey], meal)
	}

	result := make([]WeeklyMealData, 0, len(weekGroups))
	for _, meals := range weekGroups {
		if len(meals) == 0 {
			continue
		}

		var totalCalories float64
		for _, m := range meals {
			totalCalories += float64(m.Calories)
		}

		result = append(result, WeeklyMealData{
			Date:            meals[0].MealTime,
			Count:           len(meals),
			AverageCalories: totalCalories / float64(len(meals)),
		})
	}

	// Sort results by date for consistent trend calculation
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result
}

func (ta *TrendAnalyzer) interpretTrendSignificance(slope, confidence float64) string {
	sanitizedSlope := shared.SanitizeFloat64(slope)

	if sanitizedSlope < 0 {
		return "Low"
	} else if confidence < 0.7 {
		return "Moderate"
	} else {
		return "High"
	}
}
