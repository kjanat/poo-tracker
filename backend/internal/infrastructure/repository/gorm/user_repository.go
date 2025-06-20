package gorm

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// UserRepository implements user.Repository using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM user repository
func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) GetByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of error for not found
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&user.User{}, id).Error
}

func (r *UserRepository) List(limit, offset int) ([]*user.User, error) {
	var users []*user.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
