package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

var AuthService service.AuthService // should be initialized in main

// UserAPIHandler handles user-related API endpoints.
func UserAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/register" {
			RegisterHandler(w, r)
			return
		}
		if r.URL.Path == "/login" {
			LoginHandler(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	case http.MethodGet:
		if r.URL.Path == "/profile" {
			ProfileHandler(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user, token, err := AuthService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := model.LoginResponse{User: *user, Token: token, ExpiresAt: 0}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user, token, err := AuthService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	resp := model.LoginResponse{User: *user, Token: token, ExpiresAt: 0}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}

func getUserFromContext(ctx context.Context) *model.User {
	user, _ := ctx.Value("user").(*model.User)
	return user
}
