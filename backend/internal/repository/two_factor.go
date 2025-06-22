package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// UserTwoFactorRepository defines the interface for two-factor authentication operations
type UserTwoFactorRepository interface {
	Create(ctx context.Context, twoFactor *user.UserTwoFactor) error
	GetByUserID(ctx context.Context, userID string) (*user.UserTwoFactor, error)
	Update(ctx context.Context, twoFactor *user.UserTwoFactor) error
	Delete(ctx context.Context, userID string) error
	UpdateLastUsed(ctx context.Context, userID string) error
	RemoveBackupCode(ctx context.Context, userID string, code string) error
}

// memoryUserTwoFactorRepository is the in-memory implementation
type memoryUserTwoFactorRepository struct {
	mu         sync.RWMutex
	twoFactors map[string]*user.UserTwoFactor // userID -> UserTwoFactor
}

// NewMemoryUserTwoFactorRepository creates a new in-memory repository
func NewMemoryUserTwoFactorRepository() UserTwoFactorRepository {
	return &memoryUserTwoFactorRepository{
		twoFactors: make(map[string]*user.UserTwoFactor),
	}
}

func (r *memoryUserTwoFactorRepository) Create(ctx context.Context, twoFactor *user.UserTwoFactor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if twoFactor.ID == "" {
		twoFactor.ID = uuid.New().String()
	}

	now := time.Now()
	twoFactor.CreatedAt = now
	twoFactor.UpdatedAt = now

	r.twoFactors[twoFactor.UserID] = twoFactor
	return nil
}

func (r *memoryUserTwoFactorRepository) GetByUserID(ctx context.Context, userID string) (*user.UserTwoFactor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	twoFactor, exists := r.twoFactors[userID]
	if !exists {
		return nil, ErrNotFound
	}

	return twoFactor, nil
}

func (r *memoryUserTwoFactorRepository) Update(ctx context.Context, twoFactor *user.UserTwoFactor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.twoFactors[twoFactor.UserID]
	if !exists {
		return ErrNotFound
	}

	// Update fields
	existing.Secret = twoFactor.Secret
	existing.IsEnabled = twoFactor.IsEnabled
	existing.BackupCodes = twoFactor.BackupCodes
	existing.LastUsedAt = twoFactor.LastUsedAt
	existing.UpdatedAt = time.Now()

	return nil
}

func (r *memoryUserTwoFactorRepository) Delete(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.twoFactors[userID]; !exists {
		return ErrNotFound
	}

	delete(r.twoFactors, userID)
	return nil
}

func (r *memoryUserTwoFactorRepository) UpdateLastUsed(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	twoFactor, exists := r.twoFactors[userID]
	if !exists {
		return ErrNotFound
	}

	now := time.Now()
	twoFactor.LastUsedAt = &now
	twoFactor.UpdatedAt = now

	return nil
}

func (r *memoryUserTwoFactorRepository) RemoveBackupCode(ctx context.Context, userID string, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	twoFactor, exists := r.twoFactors[userID]
	if !exists {
		return ErrNotFound
	}

	// Find and remove the backup code
	for i, backupCode := range twoFactor.BackupCodes {
		if backupCode == code {
			// Remove the code from slice
			twoFactor.BackupCodes = append(twoFactor.BackupCodes[:i], twoFactor.BackupCodes[i+1:]...)
			twoFactor.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrNotFound // Backup code not found
}
