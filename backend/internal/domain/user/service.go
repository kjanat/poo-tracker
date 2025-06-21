package user

import (
	"context"
)

// Service defines the interface for user business logic
type Service interface {
	// User operations
	Create(ctx context.Context, input *CreateUserInput) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, id string, input *UpdateUserInput) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	ListWithCount(ctx context.Context, limit, offset int) ([]*User, int64, error)

	// Authentication operations
	Register(ctx context.Context, input *RegisterInput) (*User, error)
	Login(ctx context.Context, input *LoginInput) (*LoginResult, error)
	ChangePassword(ctx context.Context, userID string, input *ChangePasswordInput) error
	DeactivateAccount(ctx context.Context, userID string) error

	// Settings operations
	GetSettings(ctx context.Context, userID string) (*UserSettings, error)
	UpdateSettings(ctx context.Context, userID string, input *UpdateSettingsInput) (*UserSettings, error)

	// Validation operations
	ValidateEmail(ctx context.Context, email string) error
	ValidateUsername(ctx context.Context, username string) error
	ValidatePassword(password string) error
}

// CreateUserInput represents input for creating a user
type CreateUserInput struct {
	Email    string   `json:"email" binding:"required,email"`
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Name     string   `json:"name" binding:"required,min=1,max=100"`
	Age      *int     `json:"age,omitempty" binding:"omitempty,min=1,max=150"`
	Gender   *string  `json:"gender,omitempty"`
	Height   *float64 `json:"height,omitempty" binding:"omitempty,min=50,max=300"`
	Weight   *float64 `json:"weight,omitempty" binding:"omitempty,min=20,max=500"`
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	Email    *string  `json:"email,omitempty" binding:"omitempty,email"`
	Username *string  `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Name     *string  `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Age      *int     `json:"age,omitempty" binding:"omitempty,min=1,max=150"`
	Gender   *string  `json:"gender,omitempty"`
	Height   *float64 `json:"height,omitempty" binding:"omitempty,min=50,max=300"`
	Weight   *float64 `json:"weight,omitempty" binding:"omitempty,min=20,max=500"`
}

// RegisterInput represents input for user registration
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginInput represents input for user login
type LoginInput struct {
	EmailOrUsername string `json:"emailOrUsername" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// LoginResult represents the result of a successful login
type LoginResult struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// ChangePasswordInput represents input for changing password
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// UpdateSettingsInput represents input for updating user settings
type UpdateSettingsInput struct {
	Timezone            *string `json:"timezone,omitempty"`
	ReminderTime        *string `json:"reminderTime,omitempty"`
	ReminderEnabled     *bool   `json:"reminderEnabled,omitempty"`
	DataRetentionDays   *int    `json:"dataRetentionDays,omitempty" binding:"omitempty,min=1,max=3650"`
	PrivacyLevel        *int    `json:"privacyLevel,omitempty" binding:"omitempty,min=1,max=5"`
	NotificationEnabled *bool   `json:"notificationEnabled,omitempty"`
	ThemePreference     *string `json:"themePreference,omitempty"`
	DarkMode            *bool   `json:"darkMode,omitempty"`
	PreferredUnits      *string `json:"preferredUnits,omitempty"`
}
