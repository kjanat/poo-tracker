package user

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
	DarkMode            bool      `json:"darkMode"`
	PreferredUnits      *string   `json:"preferredUnits,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// UserUpdate represents fields that can be updated on a User
type UserUpdate struct {
	Email    *string  `json:"email,omitempty"`
	Username *string  `json:"username,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Age      *int     `json:"age,omitempty"`
	Gender   *string  `json:"gender,omitempty"`
	Height   *float64 `json:"height,omitempty"`
	Weight   *float64 `json:"weight,omitempty"`
}

// UserSettingsUpdate represents fields that can be updated on UserSettings
type UserSettingsUpdate struct {
	Timezone            *string `json:"timezone,omitempty"`
	ReminderTime        *string `json:"reminderTime,omitempty"`
	ReminderEnabled     *bool   `json:"reminderEnabled,omitempty"`
	DataRetentionDays   *int    `json:"dataRetentionDays,omitempty"`
	PrivacyLevel        *int    `json:"privacyLevel,omitempty"`
	NotificationEnabled *bool   `json:"notificationEnabled,omitempty"`
	ThemePreference     *string `json:"themePreference,omitempty"`
	DarkMode            *bool   `json:"darkMode,omitempty"`
	PreferredUnits      *string `json:"preferredUnits,omitempty"`
}

// NewUser creates a new User with defaults
func NewUser(email, username, name string) User {
	now := time.Now()
	return User{
		Email:     email,
		Username:  username,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewUserAuth creates a new UserAuth with defaults
func NewUserAuth(userID, passwordHash, provider string) UserAuth {
	now := time.Now()
	return UserAuth{
		UserID:       userID,
		PasswordHash: passwordHash,
		Provider:     provider,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewUserSettings creates a new UserSettings with defaults
func NewUserSettings(userID string) UserSettings {
	now := time.Now()
	return UserSettings{
		UserID:              userID,
		ReminderEnabled:     false,
		DataRetentionDays:   365,
		PrivacyLevel:        3,
		NotificationEnabled: true,
		DarkMode:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
