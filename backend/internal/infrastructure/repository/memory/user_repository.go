package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// UserRepository implements user.Repository interface using in-memory storage
type UserRepository struct {
	mu       sync.RWMutex
	users    map[string]*user.User
	auths    map[string]*user.UserAuth     // userID -> UserAuth
	settings map[string]*user.UserSettings // userID -> UserSettings
}

// NewUserRepository creates a new in-memory user repository
func NewUserRepository() user.Repository {
	return &UserRepository{
		users:    make(map[string]*user.User),
		settings: make(map[string]*user.UserSettings),
		auths:    make(map[string]*user.UserAuth),
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email or username already exists in a single loop
	for _, existingUser := range r.users {
		if existingUser.Email == u.Email {
			return user.ErrEmailAlreadyExists
		}
		if existingUser.Username == u.Username {
			return user.ErrUsernameAlreadyExists
		}
	}

	r.users[u.ID] = u
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[id]
	if !exists {
		return nil, user.ErrUserNotFound
	}

	return u, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

// Update updates an existing user
// Update updates an existing user with partial data
func (r *UserRepository) Update(ctx context.Context, id string, update *user.UserUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id == "" {
		return user.ErrInvalidID
	}

	existingUser, exists := r.users[id]
	if !exists {
		return user.ErrUserNotFound
	}

	// Check email uniqueness if updating email
	if update.Email != nil {
		for userID, u := range r.users {
			if userID != id && u.Email == *update.Email {
				return user.ErrEmailAlreadyExists
			}
		}
		existingUser.Email = *update.Email
	}

	// Check username uniqueness if updating username
	if update.Username != nil {
		for userID, u := range r.users {
			if userID != id && u.Username == *update.Username {
				return user.ErrUsernameAlreadyExists
			}
		}
		existingUser.Username = *update.Username
	}

	// Update other fields if provided
	if update.Name != nil {
		existingUser.Name = *update.Name
	}
	if update.Age != nil {
		existingUser.Age = update.Age
	}
	if update.Gender != nil {
		existingUser.Gender = update.Gender
	}
	if update.Height != nil {
		existingUser.Height = update.Height
	}
	if update.Weight != nil {
		existingUser.Weight = update.Weight
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return user.ErrUserNotFound
	}

	delete(r.users, id)
	delete(r.settings, id) // Also delete associated settings
	delete(r.auths, id)    // Also delete associated auth record
	return nil
}

// List lists users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	// Sort by created date (newest first)
	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.After(users[j].CreatedAt)
	})

	// Apply pagination
	start := offset
	if start > len(users) {
		return []*user.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

// Count returns the total number of users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.users)), nil
}

// CreateSettings creates user settings
func (r *UserRepository) CreateSettings(ctx context.Context, settings *user.UserSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user exists
	if _, exists := r.users[settings.UserID]; !exists {
		return user.ErrUserNotFound
	}

	r.settings[settings.UserID] = settings
	return nil
}

// GetSettingsByUserID retrieves user settings by user ID
func (r *UserRepository) GetSettingsByUserID(ctx context.Context, userID string) (*user.UserSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	settings, exists := r.settings[userID]
	if !exists {
		return nil, user.ErrUserSettingsNotFound
	}

	return settings, nil
}

// UpdateSettings updates user settings with partial data
func (r *UserRepository) UpdateSettings(ctx context.Context, userID string, update *user.UserSettingsUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if userID == "" {
		return user.ErrInvalidID
	}

	// Check if user exists
	if _, exists := r.users[userID]; !exists {
		return user.ErrUserNotFound
	}

	// Get existing settings or create new ones
	existingSettings, exists := r.settings[userID]
	if !exists {
		return user.ErrUserSettingsNotFound
	}

	// Update fields if provided
	if update.Timezone != nil {
		existingSettings.Timezone = update.Timezone
	}
	if update.ReminderTime != nil {
		existingSettings.ReminderTime = update.ReminderTime
	}
	if update.ReminderEnabled != nil {
		existingSettings.ReminderEnabled = *update.ReminderEnabled
	}
	if update.DataRetentionDays != nil {
		existingSettings.DataRetentionDays = *update.DataRetentionDays
	}
	if update.PrivacyLevel != nil {
		existingSettings.PrivacyLevel = *update.PrivacyLevel
	}
	if update.NotificationEnabled != nil {
		existingSettings.NotificationEnabled = *update.NotificationEnabled
	}
	if update.ThemePreference != nil {
		existingSettings.ThemePreference = update.ThemePreference
	}
	if update.DarkMode != nil {
		existingSettings.DarkMode = *update.DarkMode
	}
	if update.PreferredUnits != nil {
		existingSettings.PreferredUnits = update.PreferredUnits
	}

	return nil
}

// DeleteSettings deletes user settings
func (r *UserRepository) DeleteSettings(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.settings[userID]; !exists {
		return user.ErrUserSettingsNotFound
	}

	delete(r.settings, userID)
	return nil
}

// CreateAuth creates a new user authentication record
func (r *UserRepository) CreateAuth(ctx context.Context, auth *user.UserAuth) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if auth.UserID == "" {
		return user.ErrInvalidID
	}

	if _, exists := r.auths[auth.UserID]; exists {
		return user.ErrAuthAlreadyExists
	}

	r.auths[auth.UserID] = auth
	return nil
}

// GetAuthByUserID retrieves user authentication by user ID
func (r *UserRepository) GetAuthByUserID(ctx context.Context, userID string) (*user.UserAuth, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if userID == "" {
		return nil, user.ErrInvalidID
	}

	auth, exists := r.auths[userID]
	if !exists {
		return nil, user.ErrUserAuthNotFound
	}

	return auth, nil
}

// UpdateAuth updates user authentication password hash
func (r *UserRepository) UpdateAuth(ctx context.Context, userID string, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if userID == "" {
		return user.ErrInvalidID
	}

	auth, exists := r.auths[userID]
	if !exists {
		return user.ErrUserAuthNotFound
	}

	auth.PasswordHash = passwordHash
	return nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if userID == "" {
		return user.ErrInvalidID
	}

	auth, exists := r.auths[userID]
	if !exists {
		return user.ErrUserAuthNotFound
	}

	now := time.Now()
	auth.LastLoginAt = &now
	return nil
}

// DeactivateAuth deactivates user authentication
func (r *UserRepository) DeactivateAuth(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if userID == "" {
		return user.ErrInvalidID
	}

	auth, exists := r.auths[userID]
	if !exists {
		return user.ErrUserAuthNotFound
	}

	auth.IsActive = false
	return nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if email == "" {
		return false, user.ErrInvalidEmail
	}

	for _, u := range r.users {
		if u.Email == email {
			return true, nil
		}
	}

	return false, nil
}

// UsernameExists checks if a username already exists
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if username == "" {
		return false, user.ErrInvalidUsername
	}

	for _, u := range r.users {
		if u.Username == username {
			return true, nil
		}
	}

	return false, nil
}

// GetUserCount returns the total number of users
func (r *UserRepository) GetUserCount(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.users)), nil
}
