package service

import (
	"context"

	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

type AnalyticsStrategy interface {
	Summary([]model.BowelMovement) map[string]any
}

type Service struct {
	repo     repository.BowelMovementRepository
	strategy AnalyticsStrategy
}

func New(repo repository.BowelMovementRepository, s AnalyticsStrategy) *Service {
	return &Service{repo: repo, strategy: s}
}

func (s *Service) Stats(ctx context.Context) (map[string]any, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return s.strategy.Summary(list), nil
}

// ToTypedResponse converts the map response to a typed response for better type safety
func ToTypedResponse(data map[string]any) AnalyticsResponse {
	response := AnalyticsResponse{}

	if total, ok := data["total"].(int); ok {
		response.Total = total
	}
	if avgBristol, ok := data["avgBristol"].(float64); ok {
		response.AverageBristol = avgBristol
	}
	if avgPerDay, ok := data["averagePerDay"].(float64); ok {
		response.AveragePerDay = avgPerDay
	}
	if mostCommon, ok := data["mostCommonBristolType"].(int); ok {
		response.MostCommonType = mostCommon
	}
	if dist, ok := data["bristolTypeDistribution"].(map[int]int); ok {
		response.BristolDistribution = dist
	}

	return response
}
