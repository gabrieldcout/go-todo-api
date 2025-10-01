package testutils

import (
	"go-todo-api/internal/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := testDB.AutoMigrate(&models.Task{}, &models.User{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return testDB
}
