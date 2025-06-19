package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kjanat/poo-tracker/backend/internal/model"
)

type memoryRepo struct {
	mu    sync.RWMutex
	store map[string]model.BowelMovement
}

func NewMemory() BowelMovementRepository {
	return &memoryRepo{store: make(map[string]model.BowelMovement)}
}

func (m *memoryRepo) List(ctx context.Context) ([]model.BowelMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]model.BowelMovement, 0, len(m.store))
	for _, v := range m.store {
		res = append(res, v)
	}
	return res, nil
}

func (m *memoryRepo) Create(ctx context.Context, bm model.BowelMovement) (model.BowelMovement, error) {
	if bm.ID == "" {
		bm.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	bm.CreatedAt = now
	bm.UpdatedAt = now
	m.mu.Lock()
	m.store[bm.ID] = bm
	m.mu.Unlock()
	return bm, nil
}

func (m *memoryRepo) Get(ctx context.Context, id string) (model.BowelMovement, error) {
	m.mu.RLock()
	bm, ok := m.store[id]
	m.mu.RUnlock()
	if !ok {
		return model.BowelMovement{}, ErrNotFound
	}
	return bm, nil
}

func (m *memoryRepo) Update(ctx context.Context, bm model.BowelMovement) (model.BowelMovement, error) {
	m.mu.Lock()
	existing, ok := m.store[bm.ID]
	if !ok {
		m.mu.Unlock()
		return model.BowelMovement{}, ErrNotFound
	}
	if bm.BristolType != 0 {
		existing.BristolType = bm.BristolType
	}
	if bm.Notes != "" {
		existing.Notes = bm.Notes
	}
	existing.UpdatedAt = time.Now().UTC()
	m.store[bm.ID] = existing
	m.mu.Unlock()
	return existing, nil
}

func (m *memoryRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.store[id]; !ok {
		return ErrNotFound
	}
	delete(m.store, id)
	return nil
}
