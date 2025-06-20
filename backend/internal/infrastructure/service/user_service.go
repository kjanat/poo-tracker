package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// UserService implements the user business logic
type UserService struct {
	repo user.Repository
}

// NewUserService creates a new user service
func NewUserService(repo user.Repository) user.Service {
	return &UserService{
		repo: repo,
	}
}

// Create creates a new user with business validation
func (s *UserService) Create(ctx context.Context, input *user.CreateUserInput) (*user.User, error) {
	// Validate input
	if err := s.validateCreateInput(input); err != nil {
		return nil, err
	}

	// Check if email already exists
	if err := s.ValidateEmail(ctx, input.Email); err != nil {
		return nil, err
	}

	// Check if username already exists
	if err := s.ValidateUsername(ctx, input.Username); err != nil {
		return nil, err
	}

	// Create user
	userEntity := &user.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Username:  input.Username,
		Name:      input.Name,
		Age:       input.Age,
		Gender:    input.Gender,
		Height:    input.Height,
		Weight:    input.Weight,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to repository
	if err := s.repo.Create(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userEntity, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) (*user.User, error) {
	if id == "" {
		return nil, user.ErrInvalidID
	}

	userEntity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return userEntity, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if email == "" {
		return nil, user.ErrInvalidEmail
	}

	userEntity, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return userEntity, nil
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	if username == "" {
		return nil, user.ErrInvalidUsername
	}

	userEntity, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return userEntity, nil
}

// Update updates an existing user
func (s *UserService) Update(ctx context.Context, id string, input *user.UpdateUserInput) (*user.User, error) {
	if id == "" {
		return nil, user.ErrInvalidID
	}

	// Get existing user to verify it exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	// Validate update input
	if err := s.validateUpdateInput(input); err != nil {
		return nil, err
	}

	// Check email uniqueness if being updated
	if input.Email != nil {
		existing, err := s.repo.GetByEmail(ctx, *input.Email)
		if err == nil && existing.ID != id {
			return nil, user.ErrEmailAlreadyExists
		}
	}

	// Check username uniqueness if being updated
	if input.Username != nil {
		existing, err := s.repo.GetByUsername(ctx, *input.Username)
		if err == nil && existing.ID != id {
			return nil, user.ErrUsernameAlreadyExists
		}
	}

	// Convert input to repository update struct
	update := s.convertToUserUpdate(input)

	// Save changes
	if err := s.repo.Update(ctx, id, update); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Return updated user
	return s.repo.GetByID(ctx, id)
}

// Delete removes a user and all associated records (including auth and settings)
func (s *UserService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return user.ErrInvalidID
	}

	// Check if user exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to verify user exists: %w", err)
	}

	// Delete user
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves users with pagination
func (s *UserService) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	// Apply business rules for pagination
	if limit <= 0 || limit > 100 {
		limit = 20 // default
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// ListWithCount retrieves users with pagination and total count
func (s *UserService) ListWithCount(ctx context.Context, limit, offset int) ([]*user.User, int64, error) {
	// Apply business rules for pagination
	if limit <= 0 || limit > 100 {
		limit = 20 // default
	}
	if offset < 0 {
		offset = 0
	}

	// Get users and total count concurrently would be better, but keeping it simple for now
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	totalCount, err := s.repo.GetUserCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	return users, totalCount, nil
}

// Register registers a new user with password
func (s *UserService) Register(ctx context.Context, input *user.RegisterInput) (*user.User, error) {
	// Validate input
	if err := s.validateRegisterInput(input); err != nil {
		return nil, err
	}

	// Validate password
	if err := s.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	if err := s.ValidateEmail(ctx, input.Email); err != nil {
		return nil, err
	}

	// Check if username already exists
	if err := s.ValidateUsername(ctx, input.Username); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	userEntity := &user.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Username:  input.Username,
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to repository
	if err := s.repo.Create(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Create user auth
	userAuth := &user.UserAuth{
		UserID:       userEntity.ID,
		PasswordHash: string(hashedPassword),
		Provider:     "local",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateAuth(ctx, userAuth); err != nil {
		// Clean up user if auth creation fails
		_ = s.repo.Delete(ctx, userEntity.ID)
		return nil, fmt.Errorf("failed to create user auth: %w", err)
	}

	// Create default user settings
	userSettings := user.NewUserSettings(userEntity.ID)
	if err := s.repo.CreateSettings(ctx, &userSettings); err != nil {
		// Clean up user and auth if settings creation fails
		_ = s.repo.Delete(ctx, userEntity.ID)
		return nil, fmt.Errorf("failed to create user settings: %w", err)
	}

	return userEntity, nil
}

// Login authenticates a user and returns login result
func (s *UserService) Login(ctx context.Context, input *user.LoginInput) (*user.LoginResult, error) {
	// Validate input
	if err := s.validateLoginInput(input); err != nil {
		return nil, err
	}

	// Determine if input is email or username
	var userEntity *user.User
	var err error

	if strings.Contains(input.EmailOrUsername, "@") {
		userEntity, err = s.repo.GetByEmail(ctx, input.EmailOrUsername)
	} else {
		userEntity, err = s.repo.GetByUsername(ctx, input.EmailOrUsername)
	}

	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, user.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user for login: %w", err)
	}

	// Get user auth
	userAuth, err := s.repo.GetAuthByUserID(ctx, userEntity.ID)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	// Check if user is active
	if !userAuth.IsActive {
		return nil, user.ErrAccountDeactivated
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.PasswordHash), []byte(input.Password)); err != nil {
		return nil, user.ErrInvalidCredentials
	}

	// Update last login time
	_ = s.repo.UpdateLastLogin(ctx, userEntity.ID) // Ignore error for last login update

	return &user.LoginResult{
		User:  userEntity,
		Token: "", // JWT token generation would happen in the handler layer
	}, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID string, input *user.ChangePasswordInput) error {
	if userID == "" {
		return user.ErrInvalidID
	}

	// Validate input
	if err := s.validateChangePasswordInput(input); err != nil {
		return err
	}

	// Get user auth
	userAuth, err := s.repo.GetAuthByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to get user auth for password change: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return user.ErrInvalidCredentials
	}

	// Validate new password
	if err := s.ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := s.repo.UpdateAuth(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeactivateAccount deactivates a user account
func (s *UserService) DeactivateAccount(ctx context.Context, userID string) error {
	if userID == "" {
		return user.ErrInvalidID
	}

	// Check if user exists
	_, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to get user for deactivation: %w", err)
	}

	// Deactivate user auth
	if err := s.repo.DeactivateAuth(ctx, userID); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// GetSettings retrieves user settings
func (s *UserService) GetSettings(ctx context.Context, userID string) (*user.UserSettings, error) {
	if userID == "" {
		return nil, user.ErrInvalidID
	}

	settings, err := s.repo.GetSettingsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}

	return settings, nil
}

// UpdateSettings updates user settings
func (s *UserService) UpdateSettings(ctx context.Context, userID string, input *user.UpdateSettingsInput) (*user.UserSettings, error) {
	if userID == "" {
		return nil, user.ErrInvalidID
	}

	// Validate input
	if err := s.validateUpdateSettingsInput(input); err != nil {
		return nil, err
	}

	// Convert input to settings update struct
	update := s.convertToSettingsUpdateStruct(input)

	// Save changes
	if err := s.repo.UpdateSettings(ctx, userID, update); err != nil {
		return nil, fmt.Errorf("failed to update user settings: %w", err)
	}

	// Return updated settings
	return s.repo.GetSettingsByUserID(ctx, userID)
}

// ValidateEmail validates email uniqueness
func (s *UserService) ValidateEmail(ctx context.Context, email string) error {
	if email == "" {
		return user.ErrInvalidEmail
	}

	// Check if email already exists
	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to validate email: %w", err)
	}

	if exists {
		return user.ErrEmailAlreadyExists
	}

	return nil
}

// ValidateUsername validates username uniqueness
func (s *UserService) ValidateUsername(ctx context.Context, username string) error {
	if username == "" {
		return user.ErrInvalidUsername
	}

	// Check if username already exists
	exists, err := s.repo.UsernameExists(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to validate username: %w", err)
	}

	if exists {
		return user.ErrUsernameAlreadyExists
	}

	return nil
}

// ValidatePassword validates password strength
func (s *UserService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return user.ErrPasswordTooShort
	}
	if len(password) > 128 {
		return user.ErrPasswordTooLong
	}
	// Add more password validation rules as needed
	return nil
}

// Helper validation methods
func (s *UserService) validateCreateInput(input *user.CreateUserInput) error {
	if input == nil {
		return user.ErrInvalidInput
	}

	if input.Email == "" {
		return user.ErrInvalidEmail
	}

	if input.Username == "" {
		return user.ErrInvalidUsername
	}

	if input.Name == "" {
		return user.ErrInvalidInput
	}

	if input.Age != nil && (*input.Age < 1 || *input.Age > 150) {
		return user.ErrInvalidAge
	}

	if input.Height != nil && (*input.Height < 50 || *input.Height > 300) {
		return user.ErrInvalidHeight
	}

	if input.Weight != nil && (*input.Weight < 20 || *input.Weight > 500) {
		return user.ErrInvalidWeight
	}

	return nil
}

func (s *UserService) validateUpdateInput(input *user.UpdateUserInput) error {
	if input == nil {
		return user.ErrInvalidInput
	}

	if input.Email != nil && *input.Email == "" {
		return user.ErrInvalidEmail
	}

	if input.Username != nil && *input.Username == "" {
		return user.ErrInvalidUsername
	}

	if input.Name != nil && *input.Name == "" {
		return user.ErrInvalidInput
	}

	if input.Age != nil && (*input.Age < 1 || *input.Age > 150) {
		return user.ErrInvalidAge
	}

	if input.Height != nil && (*input.Height < 50 || *input.Height > 300) {
		return user.ErrInvalidHeight
	}

	if input.Weight != nil && (*input.Weight < 20 || *input.Weight > 500) {
		return user.ErrInvalidWeight
	}

	return nil
}

func (s *UserService) validateRegisterInput(input *user.RegisterInput) error {
	if input == nil {
		return user.ErrInvalidInput
	}

	if input.Email == "" {
		return user.ErrInvalidEmail
	}

	if input.Username == "" {
		return user.ErrInvalidUsername
	}

	if input.Name == "" {
		return user.ErrInvalidInput
	}

	if input.Password == "" {
		return user.ErrWeakPassword
	}

	return nil
}

func (s *UserService) validateLoginInput(input *user.LoginInput) error {
	if input == nil {
		return user.ErrInvalidInput
	}

	if input.EmailOrUsername == "" {
		return user.ErrInvalidCredentials
	}

	if input.Password == "" {
		return user.ErrWeakPassword
	}

	return nil
}

func (s *UserService) validateChangePasswordInput(input *user.ChangePasswordInput) error {
	if input == nil {
		return user.ErrInvalidInput
	}

	if input.CurrentPassword == "" {
		return user.ErrWeakPassword
	}

	if input.NewPassword == "" {
		return user.ErrWeakPassword
	}

	if input.CurrentPassword == input.NewPassword {
		return user.ErrSamePassword
	}

	return nil
}

func (s *UserService) validateUpdateSettingsInput(input *user.UpdateSettingsInput) error {
	if input == nil {
		return user.ErrInvalidInput
	}

	if input.PrivacyLevel != nil && (*input.PrivacyLevel < 1 || *input.PrivacyLevel > 5) {
		return user.ErrInvalidPrivacyLevel
	}

	if input.DataRetentionDays != nil && (*input.DataRetentionDays < 1 || *input.DataRetentionDays > 3650) {
		return user.ErrInvalidDataRetention
	}

	return nil
}

// Helper conversion methods
func (s *UserService) convertToUserUpdate(input *user.UpdateUserInput) *user.UserUpdate {
	return &user.UserUpdate{
		Email:    input.Email,
		Username: input.Username,
		Name:     input.Name,
		Age:      input.Age,
		Gender:   input.Gender,
		Height:   input.Height,
		Weight:   input.Weight,
	}
}

func (s *UserService) convertToSettingsUpdateStruct(input *user.UpdateSettingsInput) *user.UserSettingsUpdate {
	return &user.UserSettingsUpdate{
		Timezone:            input.Timezone,
		ReminderTime:        input.ReminderTime,
		ReminderEnabled:     input.ReminderEnabled,
		DataRetentionDays:   input.DataRetentionDays,
		PrivacyLevel:        input.PrivacyLevel,
		NotificationEnabled: input.NotificationEnabled,
		ThemePreference:     input.ThemePreference,
		DarkMode:            input.DarkMode,
		PreferredUnits:      input.PreferredUnits,
	}
}
