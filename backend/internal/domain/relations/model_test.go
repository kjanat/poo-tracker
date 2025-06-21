package relations

import "testing"

func TestAllCorrelationTypes(t *testing.T) {
	types := AllCorrelationTypes()
	if len(types) != 4 {
		t.Fatalf("expected 4 types, got %d", len(types))
	}
	m := map[CorrelationType]bool{}
	for _, ct := range types {
		m[ct] = true
	}
	if !m[CorrelationPositive] || !m[CorrelationNegative] || !m[CorrelationNeutral] || !m[CorrelationUnknown] {
		t.Error("missing expected correlation types")
	}
}

func TestCorrelationTypeIsValid(t *testing.T) {
	for _, ct := range AllCorrelationTypes() {
		if !ct.IsValid() {
			t.Errorf("expected %s to be valid", ct)
		}
	}
	if CorrelationType("OTHER").IsValid() {
		t.Error("expected OTHER to be invalid")
	}
}

func TestNewMealBowelMovementRelation(t *testing.T) {
	r := NewMealBowelMovementRelation("u1", "meal1", "bm1", 2)
	if r.UserID != "u1" || r.MealID != "meal1" || r.BowelMovementID != "bm1" {
		t.Error("ids not set correctly")
	}
	if r.TimeGapHours != 2 {
		t.Errorf("expected TimeGapHours 2, got %f", r.TimeGapHours)
	}
	if r.Strength != 5 {
		t.Errorf("expected default Strength 5, got %d", r.Strength)
	}
	if r.CreatedAt.IsZero() || r.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}

func TestNewMealSymptomRelation(t *testing.T) {
	r := NewMealSymptomRelation("u1", "meal1", "sym1", 3)
	if r.UserID != "u1" || r.MealID != "meal1" || r.SymptomID != "sym1" {
		t.Error("ids not set correctly")
	}
	if r.TimeGapHours != 3 {
		t.Errorf("expected TimeGapHours 3, got %f", r.TimeGapHours)
	}
	if r.Strength != 5 {
		t.Errorf("expected default Strength 5, got %d", r.Strength)
	}
	if r.CreatedAt.IsZero() || r.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}
