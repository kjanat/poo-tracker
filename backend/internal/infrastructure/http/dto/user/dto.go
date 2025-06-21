package user

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Name     string   `json:"name" binding:"required,min=1,max=100"`
	Age      *int     `json:"age,omitempty" binding:"omitempty,min=1,max=150"`
	Gender   *string  `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Height   *float64 `json:"height,omitempty" binding:"omitempty,min=50,max=300"` // cm
	Weight   *float64 `json:"weight,omitempty" binding:"omitempty,min=20,max=500"` // kg
}

// UpdateUserRequest represents the request to update user information
type UpdateUserRequest struct {
	Username *string  `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    *string  `json:"email,omitempty" binding:"omitempty,email"`
	Name     *string  `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Age      *int     `json:"age,omitempty" binding:"omitempty,min=1,max=150"`
	Gender   *string  `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Height   *float64 `json:"height,omitempty" binding:"omitempty,min=50,max=300"` // cm
	Weight   *float64 `json:"weight,omitempty" binding:"omitempty,min=20,max=500"` // kg
}

// ChangePasswordRequest represents the request to change user password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user data returned in responses
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Age       *int      `json:"age,omitempty"`
	Gender    *string   `json:"gender,omitempty"`
	Height    *float64  `json:"height,omitempty"` // cm
	Weight    *float64  `json:"weight,omitempty"` // kg
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// UserSettingsResponse represents user settings returned in responses
type UserSettingsResponse struct {
	UserID              string    `json:"user_id"`
	Timezone            *string   `json:"timezone,omitempty"`
	ReminderTime        *string   `json:"reminder_time,omitempty"`
	ReminderEnabled     bool      `json:"reminder_enabled"`
	DataRetentionDays   int       `json:"data_retention_days"`
	PrivacyLevel        int       `json:"privacy_level"`
	NotificationEnabled bool      `json:"notification_enabled"`
	ThemePreference     *string   `json:"theme_preference,omitempty"`
	DarkMode            bool      `json:"dark_mode"`
	PreferredUnits      *string   `json:"preferred_units,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// UpdateSettingsRequest represents the request to update user settings
type UpdateSettingsRequest struct {
	Timezone            *string `json:"timezone,omitempty"`
	ReminderTime        *string `json:"reminder_time,omitempty" binding:"omitempty,regexp=^([01]?[0-9]|2[0-3]):[0-5][0-9]$"`
	ReminderEnabled     *bool   `json:"reminder_enabled,omitempty"`
	DataRetentionDays   *int    `json:"data_retention_days,omitempty" binding:"omitempty,min=1,max=3650"`
	PrivacyLevel        *int    `json:"privacy_level,omitempty" binding:"omitempty,min=1,max=5"`
	NotificationEnabled *bool   `json:"notification_enabled,omitempty"`
	ThemePreference     *string `json:"theme_preference,omitempty" binding:"omitempty,oneof=light dark auto"`
	DarkMode            *bool   `json:"dark_mode,omitempty"`
	PreferredUnits      *string `json:"preferred_units,omitempty" binding:"omitempty,oneof=metric imperial"`
}

// ToUserResponse converts a domain User to UserResponse
func ToUserResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Age:       u.Age,
		Gender:    u.Gender,
		Height:    u.Height,
		Weight:    u.Weight,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserListResponse converts a slice of domain Users to UserListResponse
func ToUserListResponse(users []user.User, totalCount int64, page, pageSize int) UserListResponse {
	userRes := make([]UserResponse, len(users))
	for i, u := range users {
		userRes[i] = ToUserResponse(&u)
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return UserListResponse{
		Users:      userRes,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ToDomainUser converts CreateUserRequest to domain User
func (r *CreateUserRequest) ToDomainUser() *user.User {
	return &user.User{
		Username: r.Username,
		Email:    r.Email,
		Name:     r.Name,
		Age:      r.Age,
		Gender:   r.Gender,
		Height:   r.Height,
		Weight:   r.Weight,
	}
}

// ToUserAuth converts CreateUserRequest to domain UserAuth
func (r *CreateUserRequest) ToUserAuth(userID string) *user.UserAuth {
	return &user.UserAuth{
		UserID:   userID,
		Provider: "local", // Default provider
		IsActive: true,
		// PasswordHash will be set by the service after hashing r.Password
	}
}

// ApplyToDomainUser applies UpdateUserRequest to a domain User
func (r *UpdateUserRequest) ApplyToDomainUser(u *user.User) {
	if r.Username != nil {
		u.Username = *r.Username
	}
	if r.Email != nil {
		u.Email = *r.Email
	}
	if r.Name != nil {
		u.Name = *r.Name
	}
	if r.Age != nil {
		u.Age = r.Age
	}
	if r.Gender != nil {
		u.Gender = r.Gender
	}
	if r.Height != nil {
		u.Height = r.Height
	}
	if r.Weight != nil {
		u.Weight = r.Weight
	}
}

// ToUpdateUserInput converts UpdateUserRequest to domain UpdateUserInput
func (r *UpdateUserRequest) ToUpdateUserInput() *user.UpdateUserInput {
	return &user.UpdateUserInput{
		Email:    r.Email,
		Username: r.Username,
		Name:     r.Name,
		Age:      r.Age,
		Gender:   r.Gender,
		Height:   r.Height,
		Weight:   r.Weight,
	}
}

// ToUserSettingsResponse converts a domain UserSettings to UserSettingsResponse
func ToUserSettingsResponse(s *user.UserSettings) UserSettingsResponse {
	return UserSettingsResponse{
		UserID:              s.UserID,
		Timezone:            s.Timezone,
		ReminderTime:        s.ReminderTime,
		ReminderEnabled:     s.ReminderEnabled,
		DataRetentionDays:   s.DataRetentionDays,
		PrivacyLevel:        s.PrivacyLevel,
		NotificationEnabled: s.NotificationEnabled,
		ThemePreference:     s.ThemePreference,
		DarkMode:            s.DarkMode,
		PreferredUnits:      s.PreferredUnits,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}

// ToUpdateSettingsInput converts UpdateSettingsRequest to domain UpdateSettingsInput
func (r *UpdateSettingsRequest) ToUpdateSettingsInput() *user.UpdateSettingsInput {
	return &user.UpdateSettingsInput{
		Timezone:            r.Timezone,
		ReminderTime:        r.ReminderTime,
		ReminderEnabled:     r.ReminderEnabled,
		DataRetentionDays:   r.DataRetentionDays,
		PrivacyLevel:        r.PrivacyLevel,
		NotificationEnabled: r.NotificationEnabled,
		ThemePreference:     r.ThemePreference,
		DarkMode:            r.DarkMode,
		PreferredUnits:      r.PreferredUnits,
	}
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	if len(r.Username) < 3 || len(r.Username) > 50 {
		return user.ErrInvalidUsername
	}
	if len(r.Password) < 8 {
		return user.ErrWeakPassword
	}
	// Email validation is handled by the binding tag
	return nil
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate() error {
	if r.Username != nil && (len(*r.Username) < 3 || len(*r.Username) > 50) {
		return user.ErrInvalidUsername
	}
	return nil
}

// Validate validates the ChangePasswordRequest
func (r *ChangePasswordRequest) Validate() error {
	if len(r.NewPassword) < 8 {
		return user.ErrWeakPassword
	}
	return nil
}
