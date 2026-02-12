package impl

import (
	"context"
	"robot_scheduler/internal/model/entity"
	"robot_scheduler/internal/testutil"
	"testing"
)

func TestDeviceDAO_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	device := testutil.CreateTestDevice(t, db, entity.DeviceTypeWheelRobot)

	if device.ID == 0 {
		t.Error("Expected device ID to be set after creation")
	}
}

func TestDeviceDAO_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewDeviceDAO(db)
	ctx := context.Background()

	device := testutil.CreateTestDevice(t, db, entity.DeviceTypeWheelRobot)
	newStatus := entity.DeviceStatusOnline
	device.Status = &newStatus

	err := dao.Update(ctx, device)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, device.ID)
	if *found.Status != entity.DeviceStatusOnline {
		t.Errorf("Expected status %s, got %s", entity.DeviceStatusOnline, *found.Status)
	}
}

func TestDeviceDAO_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewDeviceDAO(db)
	ctx := context.Background()

	device := testutil.CreateTestDevice(t, db, entity.DeviceTypeWheelRobot)

	err := dao.Delete(ctx, device.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	found, _ := dao.FindByID(ctx, device.ID)
	if found != nil {
		t.Error("Expected device to be deleted")
	}
}

func TestDeviceDAO_FindByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewDeviceDAO(db)
	ctx := context.Background()

	device := testutil.CreateTestDevice(t, db, entity.DeviceTypeWheelRobot)

	found, err := dao.FindByID(ctx, device.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find device")
	}
}

func TestDeviceDAO_FindPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewDeviceDAO(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		testutil.CreateTestDevice(t, db, entity.DeviceTypeWheelRobot)
	}

	devices, total, err := dao.FindPage(ctx, 0, 3)
	if err != nil {
		t.Fatalf("FindPage failed: %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}

	if len(devices) != 3 {
		t.Errorf("Expected 3 devices, got %d", len(devices))
	}
}

func TestDeviceDAO_FindAll(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(db)

	dao := NewDeviceDAO(db)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		testutil.CreateTestDevice(t, db, entity.DeviceTypeWheelRobot)
	}

	devices, err := dao.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(devices) != 3 {
		t.Errorf("Expected 3 devices, got %d", len(devices))
	}
}
