package repository

import (
	"sync"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

type UserRepository interface {
	CreateUser(user *user.User) error
	GetUserByID(id string) (*user.User, error)
	GetUserByEmail(email string) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id string) error

	// UserAuth methods
	CreateUserAuth(auth *user.UserAuth) error
	GetUserAuthByUserID(userID string) (*user.UserAuth, error)
	GetUserAuthByEmail(email string) (*user.UserAuth, error)
	UpdateUserAuth(auth *user.UserAuth) error
}

// MemoryUserRepository is an in-memory implementation for testing.
type MemoryUserRepository struct {
	mu        sync.RWMutex
	users     map[string]*user.User
	auths     map[string]*user.UserAuth // keyed by userID
	emailToID map[string]string
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:     make(map[string]*user.User),
		auths:     make(map[string]*user.UserAuth),
		emailToID: make(map[string]string),
	}
}

func (r *MemoryUserRepository) CreateUser(user *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	r.emailToID[user.Email] = user.ID
	return nil
}

func (r *MemoryUserRepository) GetUserByID(id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *MemoryUserRepository) GetUserByEmail(email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, ok := r.emailToID[email]
	if !ok {
		return nil, ErrNotFound
	}
	user, ok := r.users[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *MemoryUserRepository) UpdateUser(user *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return ErrNotFound
	}
	r.users[user.ID] = user
	r.emailToID[user.Email] = user.ID
	return nil
}

func (r *MemoryUserRepository) DeleteUser(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	delete(r.emailToID, user.Email)
	delete(r.auths, id) // Also remove auth data
	return nil
}

func (r *MemoryUserRepository) CreateUserAuth(auth *user.UserAuth) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.auths[auth.UserID] = auth
	return nil
}

func (r *MemoryUserRepository) GetUserAuthByUserID(userID string) (*user.UserAuth, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	auth, ok := r.auths[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return auth, nil
}

func (r *MemoryUserRepository) GetUserAuthByEmail(email string) (*user.UserAuth, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, ok := r.emailToID[email]
	if !ok {
		return nil, ErrNotFound
	}
	auth, ok := r.auths[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return auth, nil
}

func (r *MemoryUserRepository) UpdateUserAuth(auth *user.UserAuth) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.auths[auth.UserID]; !ok {
		return ErrNotFound
	}
	r.auths[auth.UserID] = auth
	return nil
}
