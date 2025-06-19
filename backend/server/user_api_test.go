package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

func TestUserAPIHandler(t *testing.T) {
	// Setup test dependencies
	userRepo := repository.NewMemoryUserRepository()
	authService := &service.JWTAuthService{
		UserRepo: userRepo,
		Secret:   "test_secret",
		Expiry:   24 * time.Hour,
	}
	originalAuthService := AuthService
	AuthService = authService
	defer func() { AuthService = originalAuthService }()

	t.Run("Register - Success", func(t *testing.T) {
		reqBody := model.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		RegisterHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp model.LoginResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatal("Failed to decode response")
		}

		if resp.User.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", resp.User.Email)
		}
		if resp.Token == "" {
			t.Error("Expected token to be set")
		}
		if resp.ExpiresAt == 0 {
			t.Error("Expected ExpiresAt to be set")
		}
	})

	t.Run("Register - Invalid Email", func(t *testing.T) {
		reqBody := model.CreateUserRequest{
			Email:    "invalid-email",
			Password: "password123",
			Name:     "Test User",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		RegisterHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Register - Missing Fields", func(t *testing.T) {
		reqBody := model.CreateUserRequest{
			Email: "test2@example.com",
			// Missing password and name
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		RegisterHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Login - Success", func(t *testing.T) {
		// First register a user
		_, _, err := authService.Register("login@example.com", "password123", "Login User")
		if err != nil {
			t.Fatal("Failed to register user for login test")
		}

		reqBody := model.LoginRequest{
			Email:    "login@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		LoginHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp model.LoginResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatal("Failed to decode response")
		}

		if resp.Token == "" {
			t.Error("Expected token to be set")
		}
	})

	t.Run("Login - Invalid Credentials", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()

		LoginHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Profile - Success", func(t *testing.T) {
		user := &model.User{
			ID:    "test-user-id",
			Email: "profile@example.com",
			Name:  "Profile User",
		}

		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user", user))
		w := httptest.NewRecorder()

		ProfileHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var respUser model.User
		if err := json.NewDecoder(w.Body).Decode(&respUser); err != nil {
			t.Fatal("Failed to decode response")
		}

		if respUser.Email != "profile@example.com" {
			t.Errorf("Expected email profile@example.com, got %s", respUser.Email)
		}
	})

	t.Run("Profile - Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()

		ProfileHandler(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}
