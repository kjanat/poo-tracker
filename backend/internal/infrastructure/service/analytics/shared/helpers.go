package shared

import (
	"math"
	"sort"
	"time"
)

// calculatePercentile returns the value at the given percentile from a slice of
// float64 values. The slice is copied and sorted internally. If the slice is
// empty, the function returns 0. If it contains a single element, that value is
// returned. Percentile should be provided as a decimal (e.g., 0.25 for 25%).
func calculatePercentile(values []float64, percentile float64) float64 {
	if len(values) == 0 {
		return 0
	}
	if len(values) == 1 {
		return values[0]
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	rank := percentile * float64(len(sorted)-1)
	lower := int(math.Floor(rank))
	upper := int(math.Ceil(rank))

	if lower == upper {
		return sorted[lower]
	}

	weight := rank - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// CalculateStatistics computes basic statistical metrics for a slice of float64 values
func CalculateStatistics(values []float64) StatisticalSummary {
	if len(values) == 0 {
		return StatisticalSummary{}
	}

	// Sort values for percentile calculations
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	sort.Float64s(sortedValues)

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate median
	median := sortedValues[len(sortedValues)/2]
	if len(sortedValues)%2 == 0 {
		median = (sortedValues[len(sortedValues)/2-1] + sortedValues[len(sortedValues)/2]) / 2
	}

	// Calculate standard deviation
	varianceSum := 0.0
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}
	stdDev := math.Sqrt(varianceSum / float64(len(values)))

	// Calculate percentiles
	p25 := calculatePercentile(sortedValues, 0.25)
	p75 := calculatePercentile(sortedValues, 0.75)

	return StatisticalSummary{
		Count:        len(values),
		Mean:         mean,
		Median:       median,
		StdDev:       stdDev,
		Min:          sortedValues[0],
		Max:          sortedValues[len(sortedValues)-1],
		Percentile25: p25,
		Percentile75: p75,
	}
}

// CalculateCorrelation computes Pearson correlation coefficient between two variables
func CalculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0.0
	}

	n := float64(len(x))

	// Calculate means
	meanX, meanY := 0.0, 0.0
	for i := 0; i < len(x); i++ {
		meanX += x[i]
		meanY += y[i]
	}
	meanX /= n
	meanY /= n

	// Calculate correlation components
	numerator := 0.0
	sumXSquared := 0.0
	sumYSquared := 0.0

	for i := 0; i < len(x); i++ {
		diffX := x[i] - meanX
		diffY := y[i] - meanY
		numerator += diffX * diffY
		sumXSquared += diffX * diffX
		sumYSquared += diffY * diffY
	}

	denominator := math.Sqrt(sumXSquared * sumYSquared)
	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

// InterpretCorrelationStrength returns a string interpretation of correlation strength
func InterpretCorrelationStrength(coefficient float64) string {
	abs := math.Abs(coefficient)

	switch {
	case abs >= 0.8:
		return "Very Strong"
	case abs >= 0.6:
		return "Strong"
	case abs >= 0.4:
		return "Moderate"
	case abs >= 0.2:
		return "Weak"
	default:
		return "Very Weak"
	}
}

// CalculateTrendSlope calculates the slope of a trend line using linear regression
func CalculateTrendSlope(points []TrendPoint) float64 {
	if len(points) < 2 {
		return 0.0
	}

	n := float64(len(points))
	sumX, sumY, sumXY, sumXX := 0.0, 0.0, 0.0, 0.0

	for i, point := range points {
		x := float64(i) // Use index as x value for time series
		y := point.Value

		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	// Calculate denominator for the least squares method formula
	denominator := (n*sumXX - sumX*sumX)

	// Check for division by zero
	if math.Abs(denominator) < 1e-10 { // Use small epsilon for floating point comparison
		// If all x values are identical (or nearly so), return a flat slope (0)
		// This indicates no discernible trend in the data
		return 0.0
	}

	// Calculate slope using least squares method
	slope := (n*sumXY - sumX*sumY) / denominator

	return slope
}

// InterpretTrendDirection returns a string interpretation of trend direction
func InterpretTrendDirection(slope float64, threshold float64) string {
	if math.Abs(slope) < threshold {
		return "stable"
	}

	if slope > 0 {
		return "improving"
	}

	return "declining"
}

// GroupByDay groups timestamps by day, returning a map of day to time slice
func GroupByDay(timestamps []time.Time) map[string][]time.Time {
	groups := make(map[string][]time.Time)

	for _, ts := range timestamps {
		dayKey := ts.Format("2006-01-02")
		groups[dayKey] = append(groups[dayKey], ts)
	}

	return groups
}

// CalculateConfidenceScore calculates a confidence score based on sample size
func CalculateConfidenceScore(sampleSize int) float64 {
	// Simple confidence scoring based on sample size
	// More samples = higher confidence, with diminishing returns
	if sampleSize == 0 {
		return 0.0
	}

	if sampleSize >= 100 {
		return 0.95
	}

	if sampleSize >= 50 {
		return 0.85
	}

	if sampleSize >= 30 {
		return 0.75
	}

	if sampleSize >= 20 {
		return 0.65
	}

	if sampleSize >= 10 {
		return 0.55
	}

	return 0.3 // Low confidence for small samples
}

// FindOutliers identifies outliers in a dataset using the IQR method
func FindOutliers(values []float64) []float64 {
	if len(values) < 4 {
		return []float64{} // Need at least 4 values for IQR
	}

	stats := CalculateStatistics(values)
	iqr := stats.Percentile75 - stats.Percentile25
	lowerBound := stats.Percentile25 - 1.5*iqr
	upperBound := stats.Percentile75 + 1.5*iqr

	var outliers []float64
	for _, v := range values {
		if v < lowerBound || v > upperBound {
			outliers = append(outliers, v)
		}
	}

	return outliers
}

// SanitizeFloat64 ensures a float64 value is not NaN or Inf
func SanitizeFloat64(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0.0
	}
	return value
}

// RoundToDecimalPlaces rounds a float64 to specified decimal places
func RoundToDecimalPlaces(value float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(value*multiplier) / multiplier
}
