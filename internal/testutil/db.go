package testutil

import (
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/entity"
	"testing"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
// and automatically migrates all entities
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Initialize logger for tests (use nop logger to avoid file I/O)
	if logger.Logger == nil {
		logger.Logger = zap.NewNop()
	}

	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto-migrate all entities
	err = db.AutoMigrate(
		&entity.UserInfo{},
		&entity.Device{},
		&entity.PCDFile{},
		&entity.SemanticMap{},
		&entity.Task{},
		&entity.UserOperation{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TeardownTestDB closes the database connection
func TeardownTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
