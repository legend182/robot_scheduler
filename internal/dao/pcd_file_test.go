package impl

import (
	"context"
	"robot_scheduler/internal/testutil"
	"testing"
)

func TestPCDFileDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")

	if pcdFile.ID == 0 {
		t.Error("Expected PCD file ID to be set after creation")
	}
}

func TestPCDFileDAO_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewPCDFileDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	pcdFile.Size = 2048

	err := dao.Update(ctx, pcdFile)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, pcdFile.ID)
	if found.Size != 2048 {
		t.Errorf("Expected size 2048, got %d", found.Size)
	}
}

func TestPCDFileDAO_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewPCDFileDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")

	err := dao.Delete(ctx, pcdFile.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, pcdFile.ID)
	if found != nil {
		t.Error("Expected PCD file to be deleted")
	}
}

func TestPCDFileDAO_FindByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewPCDFileDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")

	found, err := dao.FindByID(ctx, pcdFile.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find PCD file")
	}
}

func TestPCDFileDAO_FindByName(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewPCDFileDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")

	found, err := dao.FindByName(ctx, "test.pcd")
	if err != nil {
		t.Fatalf("FindByName failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find PCD file by name")
	}

	if found.ID != pcdFile.ID {
		t.Errorf("Expected ID %d, got %d", pcdFile.ID, found.ID)
	}
}

func TestPCDFileDAO_FindPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewPCDFileDAO(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		testutil.CreateTestPCDFile(t, db, "test"+string(rune('0'+i))+".pcd")
	}

	files, total, err := dao.FindPage(ctx, 0, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}

func TestPCDFileDAO_FindAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewPCDFileDAO(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		testutil.CreateTestPCDFile(t, db, "test"+string(rune('0'+i))+".pcd")
	}

	files, err := dao.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}
