package analyzer

import (
	"sort"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
)

// TrendChangeThreshold defines the minimum percentage change to consider a trend significant
const TrendChangeThreshold = 0.1 // 10% change

// Analysis thresholds
const (
	FrequencyChangeThreshold   = 0.25 // 25% change in frequency
	ConsistencyChangeThreshold = 1    // Change of 1 point in Bristol scale
	MinSampleSize              = 5    // Minimum sample size for trend analysis
)

// Stats types for trend analysis
type (
	mealStats struct {
		fiberCount    int
		totalMeals    int
		totalCalories int
	}

	symptomStats struct {
		totalSeverity int
		count         int
	}

	mealPeriodStats struct {
		avgFiberRatio float64
		avgCalories   float64
	}
)

// DetermineTrendDirection analyzes the trend direction of health metrics
// Returns "improving", "declining", or "stable" based on the trend analysis
func (ta *TrendAnalyzer) DetermineTrendDirection(
	movements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) string {
	healthScore := ta.calculateMovingAverageHealthScore(movements, meals, symptoms)
	if len(healthScore) < 2 {
		return "stable"
	}

	firstHalf := ta.calculateAverageSlice(healthScore[:len(healthScore)/2])
	secondHalf := ta.calculateAverageSlice(healthScore[len(healthScore)/2:])

	if secondHalf > firstHalf*(1+TrendChangeThreshold) {
		return "improving"
	} else if secondHalf < firstHalf*(1-TrendChangeThreshold) {
		return "declining"
	}
	return "stable"
}

// GetSignificantTrends identifies significant health trends
func (ta *TrendAnalyzer) GetSignificantTrends(
	movements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) []string {
	var trends []string

	// Analyze movement trends
	if movementTrends := ta.analyzeBowelTrends(movements); len(movementTrends) > 0 {
		trends = append(trends, movementTrends...)
	}

	// Analyze meal trends
	if mealTrends := ta.analyzeMealTrends(meals); len(mealTrends) > 0 {
		trends = append(trends, mealTrends...)
	}

	// Analyze symptom trends
	if symptomTrends := ta.analyzeSymptomTrends(symptoms); len(symptomTrends) > 0 {
		trends = append(trends, symptomTrends...)
	}

	// Filter to most significant trends only (top 5)
	if len(trends) > 5 {
		trends = trends[:5]
	}

	return trends
}

// Helper methods
func (ta *TrendAnalyzer) calculateMovingAverageHealthScore(
	movements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) []float64 {
	if len(movements) == 0 {
		return []float64{}
	}

	// Group data by days
	dailyScores := make(map[string]float64)
	totalFactors := make(map[string]int)

	// Process bowel movements
	for _, bm := range movements {
		day := bm.RecordedAt.Format("2006-01-02")
		// Add Bristol type score (normalized to 0-1 range, ideal is 4)
		bristolScore := 1.0 - (float64(abs(bm.BristolType-4)) / 6.0)
		dailyScores[day] += bristolScore
		totalFactors[day]++
	}

	// Process meals
	for _, m := range meals {
		day := m.MealTime.Format("2006-01-02")
		// Calculate nutritional score
		nutritionScore := 0.0
		if m.FiberRich {
			nutritionScore += 0.5
		}
		if m.Calories > 0 && m.Calories < 800 {
			nutritionScore += 0.5
		}
		dailyScores[day] += nutritionScore
		totalFactors[day]++
	}

	// Process symptoms (negative impact)
	for _, s := range symptoms {
		day := s.RecordedAt.Format("2006-01-02")
		// Normalize severity to 0-1 range and invert (more severe = lower score)
		severityScore := 1.0 - (float64(s.Severity) / 10.0)
		dailyScores[day] -= severityScore
		totalFactors[day]++
	}

	// Calculate daily averages
	var dates []string
	for date := range dailyScores {
		dates = append(dates, date)
	}
	sortDates(dates)

	// Calculate final scores (normalize to 0-100 range)
	var scores []float64
	for _, date := range dates {
		if factors := totalFactors[date]; factors > 0 {
			avgScore := (dailyScores[date]/float64(factors) + 1) * 50 // Scale -1 to 1 range to 0-100
			scores = append(scores, avgScore)
		}
	}

	// Calculate 7-day moving average
	windowSize := 7
	if len(scores) < windowSize {
		return scores
	}

	var movingAverages []float64
	for i := 0; i <= len(scores)-windowSize; i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += scores[i+j]
		}
		movingAverages = append(movingAverages, sum/float64(windowSize))
	}

	return movingAverages
}

func (ta *TrendAnalyzer) calculateAverageSlice(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (ta *TrendAnalyzer) analyzeBowelTrends(movements []bowelmovement.BowelMovement) []string {
	if len(movements) < MinSampleSize {
		return nil
	}

	var trends []string

	// Sort movements by date
	sort.Slice(movements, func(i, j int) bool {
		return movements[i].RecordedAt.Before(movements[j].RecordedAt)
	})

	// Calculate daily frequencies and average Bristol scores
	type dailyStats struct {
		count        int
		bristolSum   int
		bristolCount int
	}
	dailyData := make(map[string]*dailyStats)

	for _, bm := range movements {
		day := bm.RecordedAt.Format("2006-01-02")
		if _, exists := dailyData[day]; !exists {
			dailyData[day] = &dailyStats{}
		}
		stats := dailyData[day]
		stats.count++
		stats.bristolSum += bm.BristolType
		stats.bristolCount++
	}

	// Calculate averages for first and second half of the period
	var dates []string
	for date := range dailyData {
		dates = append(dates, date)
	}
	sortDates(dates)

	midpoint := len(dates) / 2
	firstHalfFreq := 0.0
	firstHalfBristol := 0.0
	firstHalfDays := 0

	secondHalfFreq := 0.0
	secondHalfBristol := 0.0
	secondHalfDays := 0

	for i, date := range dates {
		stats := dailyData[date]
		if i < midpoint {
			firstHalfFreq += float64(stats.count)
			if stats.bristolCount > 0 {
				firstHalfBristol += float64(stats.bristolSum) / float64(stats.bristolCount)
			}
			firstHalfDays++
		} else {
			secondHalfFreq += float64(stats.count)
			if stats.bristolCount > 0 {
				secondHalfBristol += float64(stats.bristolSum) / float64(stats.bristolCount)
			}
			secondHalfDays++
		}
	}

	// Calculate averages
	firstHalfFreq /= float64(firstHalfDays)
	firstHalfBristol /= float64(firstHalfDays)
	secondHalfFreq /= float64(secondHalfDays)
	secondHalfBristol /= float64(secondHalfDays)

	// Analyze frequency changes
	freqChange := (secondHalfFreq - firstHalfFreq) / firstHalfFreq
	if abs(int(freqChange*100)) >= int(FrequencyChangeThreshold*100) {
		if freqChange > 0 {
			trends = append(trends, "Increased bowel movement frequency")
		} else {
			trends = append(trends, "Decreased bowel movement frequency")
		}
	}

	// Analyze consistency changes
	bristolChange := secondHalfBristol - firstHalfBristol
	if abs(int(bristolChange)) >= ConsistencyChangeThreshold {
		if bristolChange > 0 {
			trends = append(trends, "Stool consistency becoming looser")
		} else {
			trends = append(trends, "Stool consistency becoming firmer")
		}
	}

	return trends
}

func (ta *TrendAnalyzer) analyzeMealTrends(meals []meal.Meal) []string {
	if len(meals) < MinSampleSize {
		return nil
	}

	var trends []string

	// Sort meals by date
	sort.Slice(meals, func(i, j int) bool {
		return meals[i].MealTime.Before(meals[j].MealTime)
	})

	// Calculate daily nutritional stats
	dailyData := make(map[string]*mealStats)

	for _, m := range meals {
		day := m.MealTime.Format("2006-01-02")
		if _, exists := dailyData[day]; !exists {
			dailyData[day] = &mealStats{}
		}
		stats := dailyData[day]
		stats.totalMeals++
		if m.FiberRich {
			stats.fiberCount++
		}
		if m.Calories > 0 {
			stats.totalCalories += m.Calories
		}
	}

	// Analyze trends in first vs second half of the period
	var dates []string
	for date := range dailyData {
		dates = append(dates, date)
	}
	sortDates(dates)

	midpoint := len(dates) / 2
	firstHalf := analyzeHalfPeriodMeals(dailyData, dates[:midpoint])
	secondHalf := analyzeHalfPeriodMeals(dailyData, dates[midpoint:])

	// Compare and identify trends
	if firstHalf.avgFiberRatio > 0 && secondHalf.avgFiberRatio > 0 {
		change := (secondHalf.avgFiberRatio - firstHalf.avgFiberRatio) / firstHalf.avgFiberRatio
		if change >= 0.2 {
			trends = append(trends, "Increasing fiber intake")
		} else if change <= -0.2 {
			trends = append(trends, "Decreasing fiber intake")
		}
	}

	if firstHalf.avgCalories > 0 && secondHalf.avgCalories > 0 {
		change := (secondHalf.avgCalories - firstHalf.avgCalories) / firstHalf.avgCalories
		if change >= 0.15 {
			trends = append(trends, "Increasing caloric intake")
		} else if change <= -0.15 {
			trends = append(trends, "Decreasing caloric intake")
		}
	}

	return trends
}

func analyzeHalfPeriodMeals(dailyData map[string]*mealStats, dates []string) mealPeriodStats {
	var stats mealPeriodStats
	if len(dates) == 0 {
		return stats
	}

	totalDays := float64(len(dates))
	for _, date := range dates {
		data := dailyData[date]
		if data.totalMeals > 0 {
			stats.avgFiberRatio += float64(data.fiberCount) / float64(data.totalMeals)
			stats.avgCalories += float64(data.totalCalories) / float64(data.totalMeals)
		}
	}

	stats.avgFiberRatio /= totalDays
	stats.avgCalories /= totalDays

	return stats
}

func (ta *TrendAnalyzer) analyzeSymptomTrends(symptoms []symptom.Symptom) []string {
	if len(symptoms) < MinSampleSize {
		return nil
	}

	var trends []string

	// Sort symptoms by date
	sort.Slice(symptoms, func(i, j int) bool {
		return symptoms[i].RecordedAt.Before(symptoms[j].RecordedAt)
	})

	// Calculate daily symptom severity
	dailyData := make(map[string]*symptomStats)

	for _, s := range symptoms {
		day := s.RecordedAt.Format("2006-01-02")
		if _, exists := dailyData[day]; !exists {
			dailyData[day] = &symptomStats{}
		}
		stats := dailyData[day]
		stats.totalSeverity += s.Severity
		stats.count++
	}

	// Get dates and sort them
	var dates []string
	for date := range dailyData {
		dates = append(dates, date)
	}
	sortDates(dates)

	if len(dates) < MinSampleSize {
		return nil
	}

	// Calculate average severity for first and second half
	midpoint := len(dates) / 2
	firstHalfSeverity := calculateAverageSeverity(dailyData, dates[:midpoint])
	secondHalfSeverity := calculateAverageSeverity(dailyData, dates[midpoint:])

	// Compare severity trends
	if firstHalfSeverity > 0 && secondHalfSeverity > 0 {
		change := (secondHalfSeverity - firstHalfSeverity) / firstHalfSeverity
		if change >= 0.25 {
			trends = append(trends, "Increasing symptom severity")
		} else if change <= -0.25 {
			trends = append(trends, "Decreasing symptom severity")
		}
	}

	return trends
}

func calculateAverageSeverity(dailyData map[string]*symptomStats, dates []string) float64 {
	if len(dates) == 0 {
		return 0
	}

	var totalSeverity float64
	var totalDays int

	for _, date := range dates {
		data := dailyData[date]
		if data.count > 0 {
			totalSeverity += float64(data.totalSeverity) / float64(data.count)
			totalDays++
		}
	}

	if totalDays == 0 {
		return 0
	}

	return totalSeverity / float64(totalDays)
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sortDates(dates []string) {
	sort.Slice(dates, func(i, j int) bool {
		return dates[i] < dates[j]
	})
}
