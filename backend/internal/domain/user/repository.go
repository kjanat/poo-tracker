package user

import (
	"context"
)

// Repository defines the interface for user data persistence
type Repository interface {
	// User CRUD operations
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id string, update *UserUpdate) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)

	// UserAuth operations
	CreateAuth(ctx context.Context, auth *UserAuth) error
	GetAuthByUserID(ctx context.Context, userID string) (*UserAuth, error)
	UpdateAuth(ctx context.Context, userID string, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	DeactivateAuth(ctx context.Context, userID string) error

	// UserSettings operations
	CreateSettings(ctx context.Context, settings *UserSettings) error
	GetSettingsByUserID(ctx context.Context, userID string) (*UserSettings, error)
	UpdateSettings(ctx context.Context, userID string, update *UserSettingsUpdate) error

	// Query operations
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	GetUserCount(ctx context.Context) (int64, error)
}
