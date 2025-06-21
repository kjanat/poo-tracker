package shared

import (
	"math"
	"testing"
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
