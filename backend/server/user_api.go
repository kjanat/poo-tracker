package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	user "github.com/kjanat/poo-tracker/backend/internal/domain/user"
	userDto "github.com/kjanat/poo-tracker/backend/internal/infrastructure/http/dto/user"
	"github.com/kjanat/poo-tracker/backend/internal/middleware"
	"github.com/kjanat/poo-tracker/backend/internal/service"
	"github.com/kjanat/poo-tracker/backend/internal/validation"
)

// getClientIP extracts the client's IP address from request headers
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (for reverse proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if parts := strings.Split(xff, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Check for X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if lastColon := strings.LastIndex(ip, ":"); lastColon != -1 {
		ip = ip[:lastColon]
	}
	return ip
}

// UserAPIHandlers handles user-related API endpoints with dependency injection
type UserAPIHandlers struct {
	AuthService  service.AuthService // Made public for middleware access
	loginLimiter *middleware.RateLimiter
}

// NewUserAPIHandlers creates a new UserAPIHandlers with the provided auth service
func NewUserAPIHandlers(authService service.AuthService) *UserAPIHandlers {
	// Create a stricter rate limiter for login attempts (5 attempts per minute)
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)

	return &UserAPIHandlers{
		AuthService:  authService,
		loginLimiter: loginLimiter,
	}
}

func (h *UserAPIHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req userDto.CreateUserRequest
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
	resp := userDto.LoginResponse{User: userDto.ToUserResponse(user), Token: token}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserAPIHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req userDto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Rate limiting to prevent brute force attacks
	clientIP := getClientIP(r)
	if !h.loginLimiter.Allow(clientIP) {
		http.Error(w, "too many login attempts, please try again later", http.StatusTooManyRequests)
		return
	}

	user, token, err := h.AuthService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	resp := userDto.LoginResponse{User: userDto.ToUserResponse(user), Token: token}

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

func getUserFromContext(ctx context.Context) *user.User {
	return middleware.UserFromContext(ctx)
}
