package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

func TestUserRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	// Test Create
	testUser := &user.User{
		ID:       "test-user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}

	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrieved.ID)
	assert.Equal(t, testUser.Email, retrieved.Email)
	assert.Equal(t, testUser.Username, retrieved.Username)

	// Test GetByEmail
	retrieved, err = repo.GetByEmail(ctx, testUser.Email)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrieved.ID)

	// Test GetByUsername
	retrieved, err = repo.GetByUsername(ctx, testUser.Username)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrieved.ID)

	// Test EmailExists
	exists, err := repo.EmailExists(ctx, testUser.Email)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.EmailExists(ctx, "nonexistent@example.com")
	require.NoError(t, err)
	assert.False(t, exists)

	// Test UsernameExists
	exists, err = repo.UsernameExists(ctx, testUser.Username)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.UsernameExists(ctx, "nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)

	// Test Update
	updateData := &user.UserUpdate{
		Name: stringPtr("Updated Name"),
		Age:  intPtr(25),
	}
	err = repo.Update(ctx, testUser.ID, updateData)
	require.NoError(t, err)

	retrieved, err = repo.GetByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.Name)
	assert.Equal(t, intPtr(25), retrieved.Age)

	// Test GetUserCount
	count, err := repo.GetUserCount(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Test Delete
	err = repo.Delete(ctx, testUser.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, testUser.ID)
	assert.Error(t, err)
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserRepository_Auth(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	// First create a user
	testUser := &user.User{
		ID:       "test-user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Test CreateAuth
	testAuth := &user.UserAuth{
		UserID:       testUser.ID,
		PasswordHash: "hashedpassword123",
		Provider:     "email",
		IsActive:     true,
	}

	err = repo.CreateAuth(ctx, testAuth)
	require.NoError(t, err)

	// Test GetAuthByUserID
	retrieved, err := repo.GetAuthByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testAuth.UserID, retrieved.UserID)
	assert.Equal(t, testAuth.PasswordHash, retrieved.PasswordHash)
	assert.Equal(t, testAuth.IsActive, retrieved.IsActive)

	// Test UpdateAuth
	newPasswordHash := "newhashed456"
	err = repo.UpdateAuth(ctx, testUser.ID, newPasswordHash)
	require.NoError(t, err)

	retrieved, err = repo.GetAuthByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, newPasswordHash, retrieved.PasswordHash)

	// Test UpdateLastLogin
	err = repo.UpdateLastLogin(ctx, testUser.ID)
	require.NoError(t, err)

	retrieved, err = repo.GetAuthByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved.LastLoginAt)

	// Test DeactivateAuth
	err = repo.DeactivateAuth(ctx, testUser.ID)
	require.NoError(t, err)

	retrieved, err = repo.GetAuthByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.False(t, retrieved.IsActive)
}

func TestUserRepository_Settings(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	// First create a user
	testUser := &user.User{
		ID:       "test-user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}
	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	// Test CreateSettings
	testSettings := &user.UserSettings{
		UserID:              testUser.ID,
		Timezone:            stringPtr("UTC"),
		ReminderEnabled:     true,
		DataRetentionDays:   90,
		PrivacyLevel:        3,
		NotificationEnabled: true,
		DarkMode:            false,
	}

	err = repo.CreateSettings(ctx, testSettings)
	require.NoError(t, err)

	// Test GetSettingsByUserID
	retrieved, err := repo.GetSettingsByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testSettings.UserID, retrieved.UserID)
	assert.Equal(t, testSettings.Timezone, retrieved.Timezone)
	assert.Equal(t, testSettings.ReminderEnabled, retrieved.ReminderEnabled)

	// Test UpdateSettings
	updateData := &user.UserSettingsUpdate{
		Timezone:          stringPtr("America/New_York"),
		ReminderEnabled:   boolPtr(false),
		DataRetentionDays: intPtr(180),
	}
	err = repo.UpdateSettings(ctx, testUser.ID, updateData)
	require.NoError(t, err)

	retrieved, err = repo.GetSettingsByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, stringPtr("America/New_York"), retrieved.Timezone)
	assert.Equal(t, false, retrieved.ReminderEnabled)
	assert.Equal(t, 180, retrieved.DataRetentionDays)
}

func TestUserRepository_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	repo := NewUserRepository()

	// Test duplicate email
	user1 := &user.User{
		ID:       "user1",
		Email:    "duplicate@example.com",
		Username: "user1",
		Name:     "User 1",
	}
	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	user2 := &user.User{
		ID:       "user2",
		Email:    "duplicate@example.com", // Same email
		Username: "user2",
		Name:     "User 2",
	}
	err = repo.Create(ctx, user2)
	assert.Error(t, err)
	assert.Equal(t, user.ErrEmailAlreadyExists, err)

	// Test duplicate username
	user3 := &user.User{
		ID:       "user3",
		Email:    "unique@example.com",
		Username: "user1", // Same username
		Name:     "User 3",
	}
	err = repo.Create(ctx, user3)
	assert.Error(t, err)
	assert.Equal(t, user.ErrUsernameAlreadyExists, err)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
