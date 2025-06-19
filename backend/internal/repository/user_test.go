package repository

import (
	"sync"
	"testing"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/model"
)

func TestMemoryUserRepository(t *testing.T) {
	repo := NewMemoryUserRepository()

	t.Run("Create and Get User", func(t *testing.T) {
		user := &model.User{
			ID:        "1",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.CreateUser(user); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		got, err := repo.GetUserByID("1")
		if err != nil {
			t.Fatalf("GetUserByID failed: %v", err)
		}
		if got.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", got.Email)
		}
		if got.Name != "Test User" {
			t.Errorf("Expected name Test User, got %s", got.Name)
		}
	})

	t.Run("Get User by Email", func(t *testing.T) {
		user := &model.User{
			ID:        "2",
			Email:     "email@example.com",
			Name:      "Email User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.CreateUser(user); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		got, err := repo.GetUserByEmail("email@example.com")
		if err != nil {
			t.Fatalf("GetUserByEmail failed: %v", err)
		}
		if got.ID != "2" {
			t.Errorf("Expected ID 2, got %s", got.ID)
		}
	})

	t.Run("User Not Found Errors", func(t *testing.T) {
		_, err := repo.GetUserByID("nonexistent")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}

		_, err = repo.GetUserByEmail("nonexistent@example.com")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Create Duplicate Users", func(t *testing.T) {
		user1 := &model.User{
			ID:        "3",
			Email:     "duplicate@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		user2 := &model.User{
			ID:        "4",
			Email:     "duplicate@example.com", // Same email
			Name:      "User 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.CreateUser(user1); err != nil {
			t.Fatalf("First CreateUser failed: %v", err)
		}

		// This should succeed (we're not enforcing unique emails at repo level)
		// but when we get by email, we should get the latest one
		if err := repo.CreateUser(user2); err != nil {
			t.Fatalf("Second CreateUser failed: %v", err)
		}

		got, err := repo.GetUserByEmail("duplicate@example.com")
		if err != nil {
			t.Fatalf("GetUserByEmail failed: %v", err)
		}
		// Should get the latest user (user2)
		if got.ID != "4" {
			t.Errorf("Expected ID 4 (latest user), got %s", got.ID)
		}
	})

	t.Run("Update User", func(t *testing.T) {
		user := &model.User{
			ID:        "5",
			Email:     "update@example.com",
			Name:      "Original Name",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.CreateUser(user); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		// Update user
		user.Name = "Updated Name"
		user.UpdatedAt = time.Now()

		if err := repo.UpdateUser(user); err != nil {
			t.Fatalf("UpdateUser failed: %v", err)
		}

		got, err := repo.GetUserByID("5")
		if err != nil {
			t.Fatalf("GetUserByID after update failed: %v", err)
		}
		if got.Name != "Updated Name" {
			t.Errorf("Expected name Updated Name, got %s", got.Name)
		}
	})

	t.Run("Update Non-existent User", func(t *testing.T) {
		user := &model.User{
			ID:    "nonexistent",
			Email: "test@example.com",
			Name:  "Test",
		}

		err := repo.UpdateUser(user)
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Delete User", func(t *testing.T) {
		user := &model.User{
			ID:        "6",
			Email:     "delete@example.com",
			Name:      "Delete User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.CreateUser(user); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		// Also create auth for this user
		auth := &model.UserAuth{
			UserID:       "6",
			PasswordHash: "hashed_password",
			Provider:     "local",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := repo.CreateUserAuth(auth); err != nil {
			t.Fatalf("CreateUserAuth failed: %v", err)
		}

		// Delete user
		if err := repo.DeleteUser("6"); err != nil {
			t.Fatalf("DeleteUser failed: %v", err)
		}

		// Verify user is deleted
		_, err := repo.GetUserByID("6")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound after delete, got %v", err)
		}

		// Verify email mapping is deleted
		_, err = repo.GetUserByEmail("delete@example.com")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound for email after delete, got %v", err)
		}

		// Verify auth is deleted
		_, err = repo.GetUserAuthByUserID("6")
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound for auth after delete, got %v", err)
		}
	})

	t.Run("User Auth Operations", func(t *testing.T) {
		user := &model.User{
			ID:        "7",
			Email:     "auth@example.com",
			Name:      "Auth User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		auth := &model.UserAuth{
			UserID:       "7",
			PasswordHash: "hashed_password",
			Provider:     "local",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := repo.CreateUser(user); err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		if err := repo.CreateUserAuth(auth); err != nil {
			t.Fatalf("CreateUserAuth failed: %v", err)
		}

		// Test GetUserAuthByUserID
		gotAuth, err := repo.GetUserAuthByUserID("7")
		if err != nil {
			t.Fatalf("GetUserAuthByUserID failed: %v", err)
		}
		if gotAuth.PasswordHash != "hashed_password" {
			t.Errorf("Expected password hash hashed_password, got %s", gotAuth.PasswordHash)
		}

		// Test GetUserAuthByEmail
		gotAuth, err = repo.GetUserAuthByEmail("auth@example.com")
		if err != nil {
			t.Fatalf("GetUserAuthByEmail failed: %v", err)
		}
		if gotAuth.UserID != "7" {
			t.Errorf("Expected UserID 7, got %s", gotAuth.UserID)
		}

		// Test UpdateUserAuth
		auth.IsActive = false
		if err := repo.UpdateUserAuth(auth); err != nil {
			t.Fatalf("UpdateUserAuth failed: %v", err)
		}

		gotAuth, err = repo.GetUserAuthByUserID("7")
		if err != nil {
			t.Fatalf("GetUserAuthByUserID after update failed: %v", err)
		}
		if gotAuth.IsActive {
			t.Error("Expected IsActive to be false after update")
		}
	})

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				user := &model.User{
					ID:        string(rune('a' + id)),
					Email:     "concurrent" + string(rune('0'+id)) + "@example.com",
					Name:      "Concurrent User " + string(rune('0'+id)),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				if err := repo.CreateUser(user); err != nil {
					t.Errorf("Concurrent CreateUser failed: %v", err)
					return
				}

				// Read the user back
				got, err := repo.GetUserByID(user.ID)
				if err != nil {
					t.Errorf("Concurrent GetUserByID failed: %v", err)
					return
				}
				if got.Email != user.Email {
					t.Errorf("Concurrent read mismatch: expected %s, got %s", user.Email, got.Email)
				}
			}(i)
		}

		wg.Wait()
	})
}
