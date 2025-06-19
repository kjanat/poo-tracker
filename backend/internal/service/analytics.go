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
