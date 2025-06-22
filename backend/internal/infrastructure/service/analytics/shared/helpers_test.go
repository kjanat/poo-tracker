package shared

import (
	"math"
	"testing"
	"time"
)

func TestCalculatePercentile(t *testing.T) {
	values := []float64{1, 2, 3, 4}

	p25 := calculatePercentile(values, 0.25)
	p75 := calculatePercentile(values, 0.75)

	if math.Abs(p25-1.75) > 1e-9 {
		t.Errorf("expected 1.75 for 25th percentile, got %v", p25)
	}
	if math.Abs(p75-3.25) > 1e-9 {
		t.Errorf("expected 3.25 for 75th percentile, got %v", p75)
	}
}

func TestCalculatePercentileEdgeCases(t *testing.T) {
	if val := calculatePercentile([]float64{}, 0.25); val != 0 {
		t.Errorf("expected 0 for empty slice, got %v", val)
	}
	if val := calculatePercentile([]float64{5}, 0.75); val != 5 {
		t.Errorf("expected 5 for single element slice, got %v", val)
	}
}

func TestCalculateStatisticsUsesPercentiles(t *testing.T) {
	stats := CalculateStatistics([]float64{1, 2, 3, 4})

	if math.Abs(stats.Percentile25-1.75) > 1e-9 {
		t.Errorf("expected stats.Percentile25 to be 1.75, got %v", stats.Percentile25)
	}
	if math.Abs(stats.Percentile75-3.25) > 1e-9 {
		t.Errorf("expected stats.Percentile75 to be 3.25, got %v", stats.Percentile75)
	}
}

func TestCalculateCorrelation(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10}
	corr := CalculateCorrelation(x, y)
	if math.Abs(corr-1.0) > 1e-9 {
		t.Errorf("expected correlation 1.0, got %v", corr)
	}

	yInv := []float64{10, 8, 6, 4, 2}
	corrInv := CalculateCorrelation(x, yInv)
	if math.Abs(corrInv+1.0) > 1e-9 {
		t.Errorf("expected correlation -1.0, got %v", corrInv)
	}

	corrZero := CalculateCorrelation([]float64{1, 2}, []float64{1})
	if corrZero != 0.0 {
		t.Errorf("expected 0.0 for unequal lengths, got %v", corrZero)
	}
}

func TestInterpretCorrelationStrength(t *testing.T) {
	cases := []struct {
		val      float64
		expected string
	}{
		{0.85, "Very Strong"},
		{0.7, "Strong"},
		{0.5, "Moderate"},
		{0.3, "Weak"},
		{0.1, "Very Weak"},
		{-0.9, "Very Strong"},
	}
	for _, c := range cases {
		if got := InterpretCorrelationStrength(c.val); got != c.expected {
			t.Errorf("expected %s for %v, got %s", c.expected, c.val, got)
		}
	}
}

func TestCalculateTrendSlope(t *testing.T) {
	points := []TrendPoint{
		{Value: 1}, {Value: 2}, {Value: 3}, {Value: 4}, {Value: 5},
	}
	slope := CalculateTrendSlope(points)
	if math.Abs(slope-1.0) > 1e-9 {
		t.Errorf("expected slope 1.0, got %v", slope)
	}

	flat := []TrendPoint{{Value: 2}, {Value: 2}, {Value: 2}}
	if CalculateTrendSlope(flat) != 0.0 {
		t.Errorf("expected slope 0.0 for flat, got %v", CalculateTrendSlope(flat))
	}

	if CalculateTrendSlope([]TrendPoint{{Value: 1}}) != 0.0 {
		t.Errorf("expected 0.0 for <2 points")
	}
}

func TestInterpretTrendDirection(t *testing.T) {
	if got := InterpretTrendDirection(0.01, 0.1); got != "stable" {
		t.Errorf("expected stable, got %s", got)
	}
	if got := InterpretTrendDirection(0.2, 0.1); got != "improving" {
		t.Errorf("expected improving, got %s", got)
	}
	if got := InterpretTrendDirection(-0.2, 0.1); got != "declining" {
		t.Errorf("expected declining, got %s", got)
	}
}

func TestGroupByDay(t *testing.T) {
	times := []time.Time{
		time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 9, 0, 0, 0, time.UTC),
	}
	groups := GroupByDay(times)
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
	if len(groups["2023-01-01"]) != 2 || len(groups["2023-01-02"]) != 1 {
		t.Errorf("unexpected grouping: %+v", groups)
	}
}

func TestCalculateConfidenceScore(t *testing.T) {
	cases := []struct {
		size     int
		expected float64
	}{
		{0, 0.0}, {5, 0.3}, {10, 0.55}, {20, 0.65}, {30, 0.75}, {50, 0.85}, {100, 0.95}, {150, 0.95},
	}
	for _, c := range cases {
		if got := CalculateConfidenceScore(c.size); math.Abs(got-c.expected) > 1e-9 {
			t.Errorf("expected %v for %d, got %v", c.expected, c.size, got)
		}
	}
}

func TestFindOutliers(t *testing.T) {
	values := []float64{1, 2, 2, 2, 3, 100}
	outliers := FindOutliers(values)
	if len(outliers) != 1 || outliers[0] != 100 {
		t.Errorf("expected [100] as outlier, got %v", outliers)
	}
	if len(FindOutliers([]float64{1, 2, 3})) != 0 {
		t.Errorf("expected no outliers for <4 values")
	}
}

func TestSanitizeFloat64(t *testing.T) {
	if SanitizeFloat64(math.NaN()) != 0.0 {
		t.Errorf("expected 0.0 for NaN")
	}
	if SanitizeFloat64(math.Inf(1)) != 0.0 {
		t.Errorf("expected 0.0 for Inf")
	}
	if SanitizeFloat64(42.0) != 42.0 {
		t.Errorf("expected 42.0 for 42.0")
	}
}

func TestRoundToDecimalPlaces(t *testing.T) {
	if got := RoundToDecimalPlaces(3.14159, 2); math.Abs(got-3.14) > 1e-9 {
		t.Errorf("expected 3.14, got %v", got)
	}
	if got := RoundToDecimalPlaces(2.71828, 3); math.Abs(got-2.718) > 1e-9 {
		t.Errorf("expected 2.718, got %v", got)
	}
}
