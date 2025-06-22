package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/user"
	userDto "github.com/kjanat/poo-tracker/backend/internal/infrastructure/http/dto/user"
	"github.com/kjanat/poo-tracker/backend/internal/middleware"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

// Test helper function to create user handlers
func createTestUserHandlers() *UserAPIHandlers {
	userRepo := repository.NewMemoryUserRepository()
	authService := &service.JWTAuthService{
		UserRepo: userRepo,
		Secret:   "test_secret",
		Expiry:   24 * time.Hour,
	}
	return NewUserAPIHandlers(authService)
}

func TestRegisterHandler_Success(t *testing.T) {
	userHandlers := createTestUserHandlers()

	reqBody := userDto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	userHandlers.RegisterHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp userDto.LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response")
	}

	if resp.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp.User.Email)
	}
	if resp.Token == "" {
		t.Error("Expected token to be set")
	}
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	userHandlers := createTestUserHandlers()

	reqBody := userDto.CreateUserRequest{
		Email:    "invalid-email",
		Password: "password123",
		Name:     "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	userHandlers.RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	userHandlers := createTestUserHandlers()

	reqBody := userDto.CreateUserRequest{
		Email: "test2@example.com",
		// Missing password and name
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	userHandlers.RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	userHandlers := createTestUserHandlers()

	// First register a user
	_, _, err := userHandlers.AuthService.Register("login@example.com", "password123", "Login User")
	if err != nil {
		t.Fatal("Failed to register user for login test")
	}

	reqBody := userDto.LoginRequest{
		Username: "login@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	userHandlers.LoginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp userDto.LoginResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response")
	}

	if resp.Token == "" {
		t.Error("Expected token to be set")
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	userHandlers := createTestUserHandlers()

	reqBody := userDto.LoginRequest{
		Username: "nonexistent@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	userHandlers.LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProfileHandler_Success(t *testing.T) {
	userHandlers := createTestUserHandlers()

	user := &user.User{
		ID:    "test-user-id",
		Email: "profile@example.com",
		Name:  "Profile User",
	}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(middleware.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()

	userHandlers.ProfileHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var respUser userDto.UserResponse
	if err := json.NewDecoder(w.Body).Decode(&respUser); err != nil {
		t.Fatal("Failed to decode response")
	}

	if respUser.Email != "profile@example.com" {
		t.Errorf("Expected email profile@example.com, got %s", respUser.Email)
	}
}

func TestProfileHandler_Unauthorized(t *testing.T) {
	userHandlers := createTestUserHandlers()

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()

	userHandlers.ProfileHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
