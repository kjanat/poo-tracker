package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/model"
	"github.com/kjanat/poo-tracker/backend/internal/repository"
	"github.com/kjanat/poo-tracker/backend/internal/service"
)

func TestAuthMiddleware(t *testing.T) {
	// Setup test dependencies
	userRepo := repository.NewMemoryUserRepository()
	authService := &service.JWTAuthService{
		UserRepo: userRepo,
		Secret:   "test_secret",
		Expiry:   24 * time.Hour,
	}

	// Create a test user and get a valid token
	var testUser *model.User
	var validToken string
	var err error
	testUser, validToken, err = authService.Register("test@example.com", "password123", "Test User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create middleware instance (assuming you have one - this is a placeholder)
	middleware := func(authSvc service.AuthService) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Extract token from Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					http.Error(w, "missing authorization header", http.StatusUnauthorized)
					return
				}

				// Check for "Bearer " prefix
				const bearerPrefix = "Bearer "
				if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
					http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
					return
				}

				token := authHeader[len(bearerPrefix):]
				user, err := authSvc.ValidateToken(token)
				if err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}

				// Add user to context
				ctx := ContextWithUser(r.Context(), user)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}
	}

	// Test handler that checks if user is in context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil {
			http.Error(w, "user not found in context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("success"))
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		handler := middleware(authService)(testHandler)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
		if w.Body.String() != "success" {
			t.Errorf("Expected 'success', got %s", w.Body.String())
		}
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		handler := middleware(authService)(testHandler)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Invalid Authorization Header Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		handler := middleware(authService)(testHandler)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		handler := middleware(authService)(testHandler)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Create auth service with very short expiry
		shortExpiryAuth := &service.JWTAuthService{
			UserRepo: userRepo,
			Secret:   "test_secret",
			Expiry:   time.Millisecond, // Very short expiry
		}

		// Generate token
		_, expiredToken, err := shortExpiryAuth.Register("expired@example.com", "password123", "Expired User")
		if err != nil {
			t.Fatalf("Failed to create user with short expiry: %v", err)
		}

		// Wait for token to expire
		time.Sleep(time.Millisecond * 10)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		w := httptest.NewRecorder()

		handler := middleware(authService)(testHandler)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for expired token, got %d", w.Code)
		}
	})

	t.Run("User Context Injection", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		// Handler that checks user context
		contextTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				t.Error("User should be in context")
				http.Error(w, "no user in context", http.StatusInternalServerError)
				return
			}
			if user.ID != testUser.ID {
				t.Errorf("Expected user ID %s, got %s", testUser.ID, user.ID)
			}
			if user.Email != testUser.Email {
				t.Errorf("Expected user email %s, got %s", testUser.Email, user.Email)
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware(authService)(contextTestHandler)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Middleware Chain Integration", func(t *testing.T) {
		// Test that middleware works properly in a chain
		called := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()

		handler := middleware(authService)(nextHandler)
		handler.ServeHTTP(w, req)

		if !called {
			t.Error("Next handler should have been called")
		}
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}
