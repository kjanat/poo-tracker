package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/middleware"
	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/service"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

// UserAPIHandlers handles user-related API endpoints with dependency injection
type UserAPIHandlers struct {
	AuthService service.AuthService // Made public for middleware access
}

// NewUserAPIHandlers creates a new UserAPIHandlers with the provided auth service
func NewUserAPIHandlers(authService service.AuthService) *UserAPIHandlers {
	return &UserAPIHandlers{AuthService: authService}
}

// UserAPIHandler handles user-related API endpoints.
func (h *UserAPIHandlers) UserAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/register" {
			h.RegisterHandler(w, r)
			return
		}
		if r.URL.Path == "/login" {
			h.LoginHandler(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	case http.MethodGet:
		if r.URL.Path == "/profile" {
			h.ProfileHandler(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserAPIHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if err := validation.ValidateEmail(req.Email); err != nil {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	user, token, err := h.AuthService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		// Generic error response to avoid leaking internal details
		http.Error(w, "registration failed", http.StatusBadRequest)
		return
	}

	// Calculate actual token expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	resp := model.LoginResponse{User: *user, Token: token, ExpiresAt: expiresAt}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserAPIHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// TODO: Implement rate limiting here to prevent brute force attacks

	user, token, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Calculate actual token expiration time
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	resp := model.LoginResponse{User: *user, Token: token, ExpiresAt: expiresAt}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserAPIHandlers) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Authentication should be handled by middleware, not here
	// This handler assumes the user is already authenticated and set in context
	user := getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}

func getUserFromContext(ctx context.Context) *model.User {
	return middleware.UserFromContext(ctx)
}
