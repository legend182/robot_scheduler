package impl

import (
	"context"
	"robot_scheduler/internal/testutil"
	"testing"
)

func TestTaskDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	task := testutil.CreateTestTask(t, db, semanticMap.ID)

	if task.ID == 0 {
		t.Error("Expected task ID to be set after creation")
	}
}

func TestTaskDAO_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewTaskDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	task := testutil.CreateTestTask(t, db, semanticMap.ID)

	task.TaskInfo = "updated_task_info"

	err := dao.Update(ctx, task)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, task.ID)
	if found.TaskInfo != "updated_task_info" {
		t.Errorf("Expected task info 'updated_task_info', got %s", found.TaskInfo)
	}
}

func TestTaskDAO_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewTaskDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	task := testutil.CreateTestTask(t, db, semanticMap.ID)

	err := dao.Delete(ctx, task.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, task.ID)
	if found != nil {
		t.Error("Expected task to be deleted")
	}
}

func TestTaskDAO_FindByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewTaskDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	task := testutil.CreateTestTask(t, db, semanticMap.ID)

	found, err := dao.FindByID(ctx, task.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find task")
	}

	// Verify preload of SemanticMap and PCDFile
	if found.SemanticMap.ID == 0 {
		t.Error("Expected SemanticMap to be preloaded")
	}
}

func TestTaskDAO_FindPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewTaskDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	for i := 0; i < 5; i++ {
		testutil.CreateTestTask(t, db, semanticMap.ID)
	}

	tasks, total, err := dao.FindPage(ctx, 0, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

func TestTaskDAO_FindAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewTaskDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	for i := 0; i < 3; i++ {
		testutil.CreateTestTask(t, db, semanticMap.ID)
	}

	tasks, err := dao.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}
