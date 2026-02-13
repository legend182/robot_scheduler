package impl

import (
	"context"
	"robot_scheduler/internal/model/entity"
	"robot_scheduler/internal/testutil"
	"testing"
)

func TestUserDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	user := &entity.UserInfo{
		UserName: "testuser",
		Password: "encrypted_password",
		Role:     entity.RoleUser,
		IsLocked: 0,
	}

	err := dao.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set after creation")
	}

	// Verify user was created
	found, err := dao.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find created user")
	}

	if found.UserName != user.UserName {
		t.Errorf("Expected username %s, got %s", user.UserName, found.UserName)
	}
}

func TestUserDAO_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Create a user first
	user := testutil.CreateTestUser(t, db, "testuser", entity.RoleUser)

	// Update the user
	user.Role = entity.RoleManager
	user.IsLocked = 1

	err := dao.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	found, err := dao.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.Role != entity.RoleManager {
		t.Errorf("Expected role %s, got %s", entity.RoleManager, found.Role)
	}

	if found.IsLocked != 1 {
		t.Errorf("Expected IsLocked 1, got %d", found.IsLocked)
	}
}

func TestUserDAO_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Create a user first
	user := testutil.CreateTestUser(t, db, "testuser", entity.RoleUser)

	// Delete the user
	err := dao.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion (soft delete)
	found, err := dao.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found != nil {
		t.Error("Expected user to be deleted")
	}
}

func TestUserDAO_FindByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Create a user
	user := testutil.CreateTestUser(t, db, "testuser", entity.RoleUser)

	// Find by ID
	found, err := dao.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find user")
	}

	if found.UserName != user.UserName {
		t.Errorf("Expected username %s, got %s", user.UserName, found.UserName)
	}
}

func TestUserDAO_FindByID_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Try to find non-existent user
	found, err := dao.FindByID(ctx, 999)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found != nil {
		t.Error("Expected nil for non-existent user")
	}
}

func TestUserDAO_FindByUserName(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Create a user
	user := testutil.CreateTestUser(t, db, "testuser", entity.RoleUser)

	// Find by username
	found, err := dao.FindByUserName(ctx, "testuser")
	if err != nil {
		t.Fatalf("FindByUserName failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find user")
	}

	if found.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, found.ID)
	}
}

func TestUserDAO_FindPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Create multiple users
	for i := 0; i < 5; i++ {
		testutil.CreateTestUser(t, db, "user"+string(rune('0'+i)), entity.RoleUser)
	}

	// Test pagination
	users, total, err := dao.FindPage(ctx, 0, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}

	// Test second page
	users, total, err = dao.FindPage(ctx, 3, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users on second page, got %d", len(users))
	}
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
}

func TestUserDAO_FindAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserDAO(db)
	ctx := context.Background()

	// Create multiple users
	for i := 0; i < 3; i++ {
		testutil.CreateTestUser(t, db, "user"+string(rune('0'+i)), entity.RoleUser)
	}

	// Find all
	users, err := dao.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}
