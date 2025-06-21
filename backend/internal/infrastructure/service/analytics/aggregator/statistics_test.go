package aggregator

import (
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
)

func TestCalculateConsistencyScore(t *testing.T) {
	da := &DataAggregator{}
	now := time.Now()
	movements := []bowelmovement.BowelMovement{
		{BristolType: 4, RecordedAt: now},
		{BristolType: 4, RecordedAt: now.Add(24 * time.Hour)},
		{BristolType: 4, RecordedAt: now.Add(48 * time.Hour)},
	}
	score := da.calculateConsistencyScore(movements)
	if score != 1 {
		t.Errorf("expected 1, got %f", score)
	}

	varied := []bowelmovement.BowelMovement{
		{BristolType: 2, RecordedAt: now},
		{BristolType: 5, RecordedAt: now.Add(24 * time.Hour)},
		{BristolType: 7, RecordedAt: now.Add(48 * time.Hour)},
	}
	low := da.calculateConsistencyScore(varied)
	if low != 0 {
		t.Errorf("expected 0, got %f", low)
	}
}

func TestCalculateRegularityScore(t *testing.T) {
	da := &DataAggregator{}
	now := time.Now()
	if score := da.calculateRegularityScore(nil); score != 0 {
		t.Errorf("expected 0 for nil slice, got %f", score)
	}
	if score := da.calculateRegularityScore([]bowelmovement.BowelMovement{{RecordedAt: now}}); score != 0 {
		t.Errorf("expected 0 for single movement, got %f", score)
	}
	regular := []bowelmovement.BowelMovement{
		{RecordedAt: now},
		{RecordedAt: now.Add(24 * time.Hour)},
		{RecordedAt: now.Add(48 * time.Hour)},
	}
	high := da.calculateRegularityScore(regular)
	if high != 1 {
		t.Errorf("expected 1, got %f", high)
	}

	irregular := []bowelmovement.BowelMovement{
		{RecordedAt: now},
		{RecordedAt: now.Add(10 * time.Hour)},
		{RecordedAt: now.Add(50 * time.Hour)},
	}
	low := da.calculateRegularityScore(irregular)
	if low >= high {
		t.Errorf("expected lower score for irregular intervals")
	}
}

func TestCalculateNutritionScore(t *testing.T) {
	da := &DataAggregator{}
	meals := []meal.Meal{{Calories: 600, FiberRich: true}, {Calories: 600, FiberRich: true}}
	score := da.calculateNutritionScore(meals)
	if score != 1 {
		t.Errorf("expected 1, got %f", score)
	}

	bad := []meal.Meal{{Calories: 1200}, {Calories: 1200}}
	low := da.calculateNutritionScore(bad)
	if low != 0 {
		t.Errorf("expected 0, got %f", low)
	}
}

func TestCalculateComplianceScore(t *testing.T) {
	da := &DataAggregator{}
	now := time.Now()
	meds := []*medication.Medication{
		{Name: "A", IsActive: true, TakenAt: ptrTime(now)},
		{Name: "B", IsActive: true, TakenAt: ptrTime(now)},
		{Name: "C", IsActive: false},
	}
	score := da.calculateComplianceScore(meds)
	if score <= 0.6 || score >= 0.8 {
		t.Errorf("unexpected score %f", score)
	}

	empty := da.calculateComplianceScore(nil)
	if empty != 0 {
		t.Errorf("expected 0, got %f", empty)
	}
}

func TestCalculateEffectivenessScores(t *testing.T) {
	da := &DataAggregator{}
	now := time.Now()
	meds := []*medication.Medication{
		{Name: "A", IsActive: true, TakenAt: ptrTime(now)},
		{Name: "B", IsActive: false},
	}
	scores := da.calculateEffectivenessScores(meds)
	if scores["A"] != 1 {
		t.Errorf("expected 1 for active med, got %f", scores["A"])
	}
	if scores["overall"] != 0.5 {
		t.Errorf("expected overall 0.5, got %f", scores["overall"])
	}

	empty := da.calculateEffectivenessScores(nil)
	if len(empty) != 0 {
		t.Errorf("expected empty map, got %v", empty)
	}
}

func ptrTime(t time.Time) *time.Time { return &t }
