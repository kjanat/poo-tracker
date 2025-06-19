package repository

import (
	"errors"

	"github.com/kjanat/poo-tracker/backend/internal/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id string) error
}

// MemoryUserRepository is an in-memory implementation for testing.
type MemoryUserRepository struct {
	users map[string]*model.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{users: make(map[string]*model.User)}
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

// Only declare ErrNotFound once, at the bottom of the file.
var ErrNotFound = errors.New("user not found")
