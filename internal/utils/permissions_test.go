package utils

import (
	"robot_scheduler/internal/model/entity"
	"testing"
)

func TestGetRolePermissions_Administrator(t *testing.T) {
	permissions := GetRolePermissions(entity.RoleAdministrator)

	expectedCount := 9
	if len(permissions) != expectedCount {
		t.Errorf("Expected %d permissions, got %d", expectedCount, len(permissions))
	}

	// Administrator should have all permissions
	expectedPerms := []string{
		PermissionUserManage,
		PermissionUserView,
		PermissionMapManage,
		PermissionMapView,
		PermissionTaskManage,
		PermissionTaskView,
		PermissionDeviceManage,
		PermissionDeviceView,
		PermissionOperationView,
	}

	for _, expected := range expectedPerms {
		found := false
		for _, perm := range permissions {
			if perm == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected permission %s not found", expected)
		}
	}
}

func TestGetRolePermissions_Manager(t *testing.T) {
	permissions := GetRolePermissions(entity.RoleManager)

	// Manager should not have user:manage
	for _, perm := range permissions {
		if perm == PermissionUserManage {
			t.Error("Manager should not have user:manage permission")
		}
	}

	// Manager should have other manage permissions
	if !contains(permissions, PermissionMapManage) {
		t.Error("Manager should have map:manage permission")
	}
}

func TestGetRolePermissions_Operator(t *testing.T) {
	permissions := GetRolePermissions(entity.RoleOperator)

	// Operator should have task and device manage
	if !contains(permissions, PermissionTaskManage) {
		t.Error("Operator should have task:manage permission")
	}

	if !contains(permissions, PermissionDeviceManage) {
		t.Error("Operator should have device:manage permission")
	}

	// Operator should not have map:manage
	if contains(permissions, PermissionMapManage) {
		t.Error("Operator should not have map:manage permission")
	}
}

func TestGetRolePermissions_User(t *testing.T) {
	permissions := GetRolePermissions(entity.RoleUser)

	// User should only have view permissions
	for _, perm := range permissions {
		if perm == PermissionUserManage || perm == PermissionMapManage ||
			perm == PermissionTaskManage || perm == PermissionDeviceManage {
			t.Errorf("User should not have manage permission: %s", perm)
		}
	}
}

func TestGetRolePermissions_InvalidRole(t *testing.T) {
	permissions := GetRolePermissions(entity.RoleType("invalid"))

	if len(permissions) != 0 {
		t.Error("Invalid role should return empty permissions")
	}
}

func TestHasPermission_Success(t *testing.T) {
	if !HasPermission(entity.RoleAdministrator, PermissionUserManage) {
		t.Error("Administrator should have user:manage permission")
	}

	if !HasPermission(entity.RoleUser, PermissionUserView) {
		t.Error("User should have user:view permission")
	}
}

func TestHasPermission_Failure(t *testing.T) {
	if HasPermission(entity.RoleUser, PermissionUserManage) {
		t.Error("User should not have user:manage permission")
	}

	if HasPermission(entity.RoleManager, PermissionUserManage) {
		t.Error("Manager should not have user:manage permission")
	}
}

func TestHasAnyPermission_Success(t *testing.T) {
	// User has at least one of these permissions
	if !HasAnyPermission(entity.RoleUser, PermissionUserView, PermissionUserManage) {
		t.Error("User should have at least one permission")
	}

	// Empty permissions should return true
	if !HasAnyPermission(entity.RoleUser) {
		t.Error("Empty permissions should return true")
	}
}

func TestHasAnyPermission_Failure(t *testing.T) {
	// User has none of these permissions
	if HasAnyPermission(entity.RoleUser, PermissionUserManage, PermissionMapManage) {
		t.Error("User should not have any of these permissions")
	}
}

func TestHasAllPermissions_Success(t *testing.T) {
	// Administrator has all permissions
	if !HasAllPermissions(entity.RoleAdministrator, PermissionUserManage, PermissionMapManage) {
		t.Error("Administrator should have all permissions")
	}

	// Empty permissions should return true
	if !HasAllPermissions(entity.RoleUser) {
		t.Error("Empty permissions should return true")
	}
}

func TestHasAllPermissions_Failure(t *testing.T) {
	// User doesn't have all these permissions
	if HasAllPermissions(entity.RoleUser, PermissionUserView, PermissionUserManage) {
		t.Error("User should not have all these permissions")
	}

	// Manager doesn't have user:manage
	if HasAllPermissions(entity.RoleManager, PermissionMapManage, PermissionUserManage) {
		t.Error("Manager should not have user:manage permission")
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
