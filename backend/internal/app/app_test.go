package app

import (
	"testing"
)

func TestAppCanStart(t *testing.T) {
	// Set environment for SQLite testing
	t.Setenv("DB_TYPE", "sqlite")
	t.Setenv("SQLITE_DSN", ":memory:")
	t.Setenv("PORT", "0") // Use random port for testing

	// Create app
	app, err := New()
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Test that container is properly initialized
	if app.container == nil {
		t.Fatal("Container not initialized")
	}

	if app.container.Config == nil {
		t.Fatal("Config not initialized")
	}

	if app.container.Database == nil {
		t.Fatal("Database not initialized")
	}

	// Test database connection
	db := app.container.Database.GetDB()
	if db == nil {
		t.Fatal("Database connection not established")
	}

	// Test cleanup
	if err := app.container.Cleanup(); err != nil {
		t.Fatalf("Failed to cleanup: %v", err)
	}
}
