package service

import (
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/repository"
)

func TestJWTAuthService_RegisterAndLogin(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	svc := &JWTAuthService{UserRepo: repo, Secret: "secret", Expiry: time.Hour}

	user, token, err := svc.Register("test@example.com", "password", "Test User")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.ID == "" || token == "" {
		t.Fatal("expected user and token to be returned")
	}

	user2, token2, err := svc.Login("test@example.com", "password")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user2.ID != user.ID {
		t.Errorf("expected same user ID, got %s", user2.ID)
	}
	if token2 == "" {
		t.Error("expected token from login")
	}
}
