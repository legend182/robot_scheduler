package impl

import (
	"context"
	"robot_scheduler/internal/testutil"
	"testing"
)

func TestSemanticMapDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)

	if semanticMap.ID == 0 {
		t.Error("Expected semantic map ID to be set after creation")
	}
}

func TestSemanticMapDAO_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewSemanticMapDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	semanticMap.SemanticInfo = "updated_info"

	err := dao.Update(ctx, semanticMap)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, semanticMap.ID)
	if found.SemanticInfo != "updated_info" {
		t.Errorf("Expected semantic info 'updated_info', got %s", found.SemanticInfo)
	}
}

func TestSemanticMapDAO_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewSemanticMapDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)

	err := dao.Delete(ctx, semanticMap.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, semanticMap.ID)
	if found != nil {
		t.Error("Expected semantic map to be deleted")
	}
}

func TestSemanticMapDAO_FindByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewSemanticMapDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	semanticMap := testutil.CreateTestSemanticMap(t, db, pcdFile.ID)

	found, err := dao.FindByID(ctx, semanticMap.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find semantic map")
	}

	// Verify preload of PCDFile
	if found.PCDFile.ID == 0 {
		t.Error("Expected PCDFile to be preloaded")
	}
}

func TestSemanticMapDAO_FindPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewSemanticMapDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	for i := 0; i < 5; i++ {
		testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	}

	maps, total, err := dao.FindPage(ctx, 0, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(maps) != 3 {
		t.Errorf("Expected 3 maps, got %d", len(maps))
	}
}

func TestSemanticMapDAO_FindAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewSemanticMapDAO(db)
	ctx := context.Background()

	pcdFile := testutil.CreateTestPCDFile(t, db, "test.pcd")
	for i := 0; i < 3; i++ {
		testutil.CreateTestSemanticMap(t, db, pcdFile.ID)
	}

	maps, err := dao.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(maps) != 3 {
		t.Errorf("Expected 3 maps, got %d", len(maps))
	}
}
