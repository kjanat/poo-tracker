package service

import bm "github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"

type AvgBristol struct{}

func (AvgBristol) Summary(list []bm.BowelMovement) map[string]any {
	if len(list) == 0 {
		return map[string]any{
			"total":      0,
			"avgBristol": 0,
		}
	}
	sum := 0
	for _, bm := range list {
		sum += bm.BristolType
	}
	return map[string]any{
		"total":      len(list),
		"avgBristol": float64(sum) / float64(len(list)),
	}
}
