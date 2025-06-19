package repository

import (
	"github.com/kjanat/poo-tracker/backend/internal/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id string) error

	// UserAuth methods
	CreateUserAuth(auth *model.UserAuth) error
	GetUserAuthByUserID(userID string) (*model.UserAuth, error)
	GetUserAuthByEmail(email string) (*model.UserAuth, error)
	UpdateUserAuth(auth *model.UserAuth) error
}

// MemoryUserRepository is an in-memory implementation for testing.
type MemoryUserRepository struct {
	users     map[string]*model.User
	auths     map[string]*model.UserAuth // keyed by userID
	emailToID map[string]string
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:     make(map[string]*model.User),
		auths:     make(map[string]*model.UserAuth),
		emailToID: make(map[string]string),
	}
}

func (r *MemoryUserRepository) CreateUser(user *model.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) GetUserByID(id string) (*model.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *MemoryUserRepository) GetUserByEmail(email string) (*model.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MemoryUserRepository) UpdateUser(user *model.User) error {
	if _, ok := r.users[user.ID]; !ok {
		return ErrNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) DeleteUser(id string) error {
	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *MemoryUserRepository) CreateUserAuth(auth *model.UserAuth) error {
	r.auths[auth.UserID] = auth
	return nil
}

func (r *MemoryUserRepository) GetUserAuthByUserID(userID string) (*model.UserAuth, error) {
	auth, ok := r.auths[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return auth, nil
}

func (r *MemoryUserRepository) GetUserAuthByEmail(email string) (*model.UserAuth, error) {
	for _, user := range r.users {
		if user.Email == email {
			return r.GetUserAuthByUserID(user.ID)
		}
	}
	return nil, ErrNotFound
}

func (r *MemoryUserRepository) UpdateUserAuth(auth *model.UserAuth) error {
	if _, ok := r.auths[auth.UserID]; !ok {
		return ErrNotFound
	}
	r.auths[auth.UserID] = auth
	return nil
}
