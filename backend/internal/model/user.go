package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Age       *int      `json:"age,omitempty"`
	Gender    *string   `json:"gender,omitempty"`
	Height    *float64  `json:"height,omitempty"` // cm
	Weight    *float64  `json:"weight,omitempty"` // kg
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserAuth represents authentication data for a user
type UserAuth struct {
	UserID       string     `json:"userId"`
	PasswordHash string     `json:"-"` // Never serialize password hash
	Provider     string     `json:"provider"`
	IsActive     bool       `json:"isActive"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// UserSettings represents user preferences and settings
type UserSettings struct {
	UserID              string    `json:"userId"`
	Timezone            *string   `json:"timezone,omitempty"`
	ReminderTime        *string   `json:"reminderTime,omitempty"` // HH:MM format
	ReminderEnabled     bool      `json:"reminderEnabled"`
	DataRetentionDays   int       `json:"dataRetentionDays"`
	PrivacyLevel        int       `json:"privacyLevel"` // 1-5 scale
	NotificationEnabled bool      `json:"notificationEnabled"`
	ThemePreference     *string   `json:"themePreference,omitempty"` // light, dark, auto
	DarkMode            bool      `json:"dark_mode"`
	PreferredUnits      *string   `json:"preferred_units,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Name     string   `json:"name" binding:"required"`
	Age      *int     `json:"age,omitempty"`
	Gender   *string  `json:"gender,omitempty"`
	Height   *float64 `json:"height,omitempty"`
	Weight   *float64 `json:"weight,omitempty"`
}

// UpdateUserRequest represents the request to update user information
type UpdateUserRequest struct {
	Name   *string  `json:"name,omitempty"`
	Age    *int     `json:"age,omitempty"`
	Gender *string  `json:"gender,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Weight *float64 `json:"weight,omitempty"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User      User   `json:"user"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

// UpdateUserSettingsRequest represents the request to update user settings
type UpdateUserSettingsRequest struct {
	Timezone            *string `json:"timezone,omitempty"`
	ReminderTime        *string `json:"reminderTime,omitempty"`
	ReminderEnabled     *bool   `json:"reminderEnabled,omitempty"`
	DataRetentionDays   *int    `json:"dataRetentionDays,omitempty"`
	PrivacyLevel        *int    `json:"privacyLevel,omitempty"`
	NotificationEnabled *bool   `json:"notificationEnabled,omitempty"`
	ThemePreference     *string `json:"themePreference,omitempty"`
	DarkMode            *bool   `json:"dark_mode,omitempty"`
	PreferredUnits      *string `json:"preferred_units,omitempty"`
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

// RefreshTokenRequest represents the request to refresh access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// UserProfile represents a user's public profile
type UserProfile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Age       *int      `json:"age,omitempty"`
	Gender    *string   `json:"gender,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
