package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
)

// MockUserRepository is a mock for user.Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, update *user.UserUpdate) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserRepository) CreateAuth(ctx context.Context, auth *user.UserAuth) error {
	args := m.Called(ctx, auth)
	return args.Error(0)
}

func (m *MockUserRepository) GetAuthByUserID(ctx context.Context, userID string) (*user.UserAuth, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserAuth), args.Error(1)
}

func (m *MockUserRepository) UpdateAuth(ctx context.Context, userID string, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) DeactivateAuth(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) CreateSettings(ctx context.Context, settings *user.UserSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockUserRepository) GetSettingsByUserID(ctx context.Context, userID string) (*user.UserSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserSettings), args.Error(1)
}

func (m *MockUserRepository) UpdateSettings(ctx context.Context, userID string, update *user.UserSettingsUpdate) error {
	args := m.Called(ctx, userID, update)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetUserCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	input := &user.RegisterInput{
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("EmailExists", ctx, input.Email).Return(false, nil)
	mockRepo.On("UsernameExists", ctx, input.Username).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)
	mockRepo.On("CreateAuth", ctx, mock.AnythingOfType("*user.UserAuth")).Return(nil)
	mockRepo.On("CreateSettings", ctx, mock.AnythingOfType("*user.UserSettings")).Return(nil)

	result, err := service.Register(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Email, result.Email)
	assert.Equal(t, input.Username, result.Username)
	assert.Equal(t, input.Name, result.Name)
	assert.NotEmpty(t, result.ID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_EmailExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	input := &user.RegisterInput{
		Email:    "existing@example.com",
		Username: "testuser",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("EmailExists", ctx, input.Email).Return(true, nil)

	result, err := service.Register(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, user.ErrEmailAlreadyExists, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := "user123"
	expectedUser := &user.User{
		ID:       userID,
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}

	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	result, err := service.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := "nonexistent"

	mockRepo.On("GetByID", ctx, userID).Return(nil, user.ErrUserNotFound)

	result, err := service.GetByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, user.ErrUserNotFound))

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	input := &user.LoginInput{
		EmailOrUsername: "test@example.com",
		Password:        "password123",
	}

	// Create a test user
	testUser := &user.User{
		ID:       "user123",
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}

	// Create test auth with properly hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testAuth := &user.UserAuth{
		UserID:       testUser.ID,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	mockRepo.On("GetByEmail", ctx, input.EmailOrUsername).Return(testUser, nil)
	mockRepo.On("GetAuthByUserID", ctx, testUser.ID).Return(testAuth, nil)
	mockRepo.On("UpdateLastLogin", ctx, testUser.ID).Return(nil)

	result, err := service.Login(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result != nil {
		assert.Equal(t, testUser.ID, result.User.ID)
		assert.Equal(t, testUser.Email, result.User.Email)
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	input := &user.LoginInput{
		EmailOrUsername: "test@example.com",
		Password:        "wrongpassword",
	}

	testUser := &user.User{
		ID:       "user123",
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testAuth := &user.UserAuth{
		UserID:       testUser.ID,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	mockRepo.On("GetByEmail", ctx, input.EmailOrUsername).Return(testUser, nil)
	mockRepo.On("GetAuthByUserID", ctx, testUser.ID).Return(testAuth, nil)

	result, err := service.Login(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, user.ErrInvalidCredentials))

	mockRepo.AssertExpectations(t)
}
