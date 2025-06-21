package service

import (
	"testing"

	bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
)

func TestAvgBristol_Summary(t *testing.T) {
	strat := AvgBristol{}

	t.Run("Empty Input", func(t *testing.T) {
		got := strat.Summary([]bm.BowelMovement{})
		if got["total"] != 0 {
			t.Errorf("expected total 0, got %v", got["total"])
		}
		if got["avgBristol"] != 0 {
			t.Errorf("expected avgBristol 0, got %v", got["avgBristol"])
		}
	})

	t.Run("Multiple Entries", func(t *testing.T) {
		list := []bm.BowelMovement{
			{BristolType: 3},
			{BristolType: 4},
			{BristolType: 5},
		}
		got := strat.Summary(list)

		total, ok := got["total"].(int)
		if !ok {
			t.Fatalf("total should be int, got %T", got["total"])
		}
		if total != 3 {
			t.Errorf("expected total 3, got %d", total)
		}

		avg, ok := got["avgBristol"].(float64)
		if !ok {
			t.Fatalf("avgBristol should be float64, got %T", got["avgBristol"])
		}
		if avg != 4 {
			t.Errorf("expected avgBristol 4, got %f", avg)
		}
	})
}
