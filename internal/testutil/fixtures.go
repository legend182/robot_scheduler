package testutil

import (
	"robot_scheduler/internal/model/entity"
	"testing"

	"gorm.io/gorm"
)

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *gorm.DB, username string, role entity.RoleType) *entity.UserInfo {
	t.Helper()

	user := &entity.UserInfo{
		UserName: username,
		Password: "encrypted_password",
		Role:     role,
		IsLocked: 0,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestDevice creates a test device in the database
func CreateTestDevice(t *testing.T, db *gorm.DB, deviceType entity.DeviceType) *entity.Device {
	t.Helper()

	status := entity.DeviceStatusOffline
	device := &entity.Device{
		Type:    deviceType,
		Company: entity.CompanyCyborg,
		Port:    8080,
		Status:  &status,
	}

	if err := db.Create(device).Error; err != nil {
		t.Fatalf("Failed to create test device: %v", err)
	}

	return device
}

// CreateTestPCDFile creates a test PCD file in the database
func CreateTestPCDFile(t *testing.T, db *gorm.DB, name string) *entity.PCDFile {
	t.Helper()

	pcdFile := &entity.PCDFile{
		Name:     name,
		Area:     "test_area",
		Path:     "/test/path/" + name,
		UserName: "test_user",
		Size:     1024,
	}

	if err := db.Create(pcdFile).Error; err != nil {
		t.Fatalf("Failed to create test PCD file: %v", err)
	}

	return pcdFile
}

// CreateTestSemanticMap creates a test semantic map in the database
func CreateTestSemanticMap(t *testing.T, db *gorm.DB, pcdFileID uint) *entity.SemanticMap {
	t.Helper()

	semanticMap := &entity.SemanticMap{
		PCDFileID:    pcdFileID,
		UserName:     "test_user",
		SemanticInfo: "test_semantic_info",
	}

	if err := db.Create(semanticMap).Error; err != nil {
		t.Fatalf("Failed to create test semantic map: %v", err)
	}

	return semanticMap
}

// CreateTestTask creates a test task in the database
func CreateTestTask(t *testing.T, db *gorm.DB, semanticMapID uint) *entity.Task {
	t.Helper()

	status := entity.TaskStatusPending
	task := &entity.Task{
		SemanticMapID: semanticMapID,
		UserName:      "test_user",
		TaskInfo:      "test_task_info",
		Status:        &status,
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	return task
}

// CreateTestUserOperation creates a test user operation in the database
func CreateTestUserOperation(t *testing.T, db *gorm.DB, username string, operation entity.OperationType) *entity.UserOperation {
	t.Helper()

	userOp := &entity.UserOperation{
		UserName:  username,
		Operation: operation,
		Module:    "test_module",
		ExtraInfo: "{}",
	}

	if err := db.Create(userOp).Error; err != nil {
		t.Fatalf("Failed to create test user operation: %v", err)
	}

	return userOp
}
