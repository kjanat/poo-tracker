package repository

import (
	"poo-tracker/internal/model"
	"testing"
)

func TestMemoryUserRepository(t *testing.T) {
	repo := NewMemoryUserRepository()
	user := &model.User{ID: "1", Email: "test@example.com", Username: "testuser"}
	if err := repo.CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	got, err := repo.GetUserByID("1")
	if err != nil || got.Email != "test@example.com" {
		t.Errorf("GetUserByID failed: %v", err)
	}
}
