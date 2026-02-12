package impl

import (
	"context"
	"robot_scheduler/internal/model/entity"
	"robot_scheduler/internal/testutil"
	"testing"
)

func TestUserOperationDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	userOp := testutil.CreateTestUserOperation(t, db, "testuser", entity.OperationCreate)

	if userOp.ID == 0 {
		t.Error("Expected user operation ID to be set after creation")
	}
}

func TestUserOperationDAO_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserOperationDAO(db)
	ctx := context.Background()

	userOp := testutil.CreateTestUserOperation(t, db, "testuser", entity.OperationCreate)
	userOp.ExtraInfo = "{\"updated\": true}"

	err := dao.Update(ctx, userOp)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, userOp.ID)
	if found.ExtraInfo != "{\"updated\": true}" {
		t.Errorf("Expected extra info to be updated")
	}
}

func TestUserOperationDAO_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserOperationDAO(db)
	ctx := context.Background()

	userOp := testutil.CreateTestUserOperation(t, db, "testuser", entity.OperationCreate)

	err := dao.Delete(ctx, userOp.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, userOp.ID)
	if found != nil {
		t.Error("Expected user operation to be deleted")
	}
}

func TestUserOperationDAO_FindByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserOperationDAO(db)
	ctx := context.Background()

	userOp := testutil.CreateTestUserOperation(t, db, "testuser", entity.OperationCreate)

	found, err := dao.FindByID(ctx, userOp.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find user operation")
	}
}

func TestUserOperationDAO_FindPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserOperationDAO(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		testutil.CreateTestUserOperation(t, db, "testuser", entity.OperationCreate)
	}

	ops, total, err := dao.FindPage(ctx, 0, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(ops) != 3 {
		t.Errorf("Expected 3 operations, got %d", len(ops))
	}
}

func TestUserOperationDAO_FindAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewUserOperationDAO(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		testutil.CreateTestUserOperation(t, db, "testuser", entity.OperationCreate)
	}

	ops, err := dao.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(ops) != 3 {
		t.Errorf("Expected 3 operations, got %d", len(ops))
	}
}
