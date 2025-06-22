package service

import (
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

func setupAuthService() *JWTAuthService {
	userRepo := repository.NewMemoryUserRepository()
	return &JWTAuthService{
		UserRepo: userRepo,
		Secret:   "test_secret_key",
		Expiry:   24 * time.Hour,
	}
}

func TestJWTAuthService_Register(t *testing.T) {
	authService := setupAuthService()

	t.Run("Successful Registration", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		name := "Test User"

		user, token, err := authService.Register(email, password, name)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if user == nil {
			t.Fatal("Expected user to be returned")
		}
		if user.Email != email {
			t.Errorf("Expected email %s, got %s", email, user.Email)
		}
		if user.Name != name {
			t.Errorf("Expected name %s, got %s", name, user.Name)
		}
		if token == "" {
			t.Error("Expected token to be generated")
		}
		if user.ID == "" {
			t.Error("Expected user ID to be generated")
		}
	})

	t.Run("Duplicate Email Registration", func(t *testing.T) {
		email := "duplicate@example.com"
		password := "password123"
		name := "Test User"

		// Register first user
		_, _, err := authService.Register(email, password, name)
		if err != nil {
			t.Fatalf("First registration failed: %v", err)
		}

		// Try to register with same email
		_, _, err = authService.Register(email, password, "Another User")
		if err == nil {
			t.Error("Expected error for duplicate email registration")
		}
		if err.Error() != "email already registered" {
			t.Errorf("Expected 'email already registered' error, got %v", err)
		}
	})

	t.Run("Password Validation", func(t *testing.T) {
		email := "password@example.com"
		password := "test123"
		name := "Password User"

		user, _, err := authService.Register(email, password, name)
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}

		// Verify password was hashed (not stored in plain text)
		auth, err := authService.UserRepo.GetUserAuthByUserID(user.ID)
		if err != nil {
			t.Fatalf("Failed to get user auth: %v", err)
		}
		if auth.PasswordHash == password {
			t.Error("Password should be hashed, not stored in plain text")
		}
		if len(auth.PasswordHash) < 10 {
			t.Error("Password hash seems too short")
		}
	})
}

func TestJWTAuthService_Login(t *testing.T) {
	authService := setupAuthService()

	t.Run("Successful Login", func(t *testing.T) {
		email := "login@example.com"
		password := "password123"
		name := "Login User"

		// Register user first
		originalUser, _, err := authService.Register(email, password, name)
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}

		// Login with correct credentials
		user, token, err := authService.Login(email, password)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if user == nil {
			t.Fatal("Expected user to be returned")
		}
		if user.ID != originalUser.ID {
			t.Error("Login should return the same user")
		}
		if token == "" {
			t.Error("Expected token to be generated")
		}
	})

	t.Run("Invalid Email", func(t *testing.T) {
		_, _, err := authService.Login("nonexistent@example.com", "password123")
		if err == nil {
			t.Error("Expected error for non-existent email")
		}
		if err.Error() != "invalid credentials" {
			t.Errorf("Expected 'invalid credentials' error, got %v", err)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		email := "wrongpass@example.com"
		password := "password123"
		name := "Wrong Pass User"

		// Register user
		_, _, err := authService.Register(email, password, name)
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}

		// Login with wrong password
		_, _, err = authService.Login(email, "wrongpassword")
		if err == nil {
			t.Error("Expected error for wrong password")
		}
		if err.Error() != "invalid credentials" {
			t.Errorf("Expected 'invalid credentials' error, got %v", err)
		}
	})
}

func TestJWTAuthService_ValidateToken(t *testing.T) {
	authService := setupAuthService()

	t.Run("Valid Token", func(t *testing.T) {
		email := "token@example.com"
		password := "password123"
		name := "Token User"

		// Register and get token
		originalUser, token, err := authService.Register(email, password, name)
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}

		// Validate token
		user, err := authService.ValidateToken(token)
		if err != nil {
			t.Fatalf("Token validation failed: %v", err)
		}
		if user == nil {
			t.Fatal("Expected user to be returned")
		}
		if user.ID != originalUser.ID {
			t.Error("Token should validate to the same user")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := authService.ValidateToken("invalid.token.here")
		if err == nil {
			t.Error("Expected error for invalid token")
		}
	})

	t.Run("Malformed Token", func(t *testing.T) {
		_, err := authService.ValidateToken("malformed-token")
		if err == nil {
			t.Error("Expected error for malformed token")
		}
	})

	t.Run("Token with Wrong Algorithm", func(t *testing.T) {
		// Create a token with wrong signing method (this would need jwt library)
		// For now, just test that our validation rejects tokens with wrong algorithms
		// The actual implementation in ValidateToken checks for HS256
		invalidToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJ0ZXN0In0."
		_, err := authService.ValidateToken(invalidToken)
		if err == nil {
			t.Error("Expected error for token with wrong algorithm")
		}
	})

	t.Run("Empty Token", func(t *testing.T) {
		_, err := authService.ValidateToken("")
		if err == nil {
			t.Error("Expected error for empty token")
		}
	})
}
