package meal

import (
	"testing"
	"time"
)

func TestNewMealDefaults(t *testing.T) {
	mealTime := time.Now()
	m := NewMeal("user1", "Lunch", mealTime)

	if m.UserID != "user1" {
		t.Errorf("expected userID user1, got %s", m.UserID)
	}
	if m.Name != "Lunch" {
		t.Errorf("expected name Lunch, got %s", m.Name)
	}
	if !m.MealTime.Equal(mealTime) {
		t.Errorf("expected mealTime %v, got %v", mealTime, m.MealTime)
	}
	if m.FiberRich || m.Dairy || m.Gluten {
		t.Error("expected default boolean fields to be false")
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}
