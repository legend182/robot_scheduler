package utils

import "robot_scheduler/internal/model/entity"

// 权限常量定义
const (
	PermissionUserManage    = "user:manage"    // 用户管理
	PermissionUserView      = "user:view"      // 用户查看
	PermissionMapManage     = "map:manage"     // 地图管理（创建/编辑/删除）
	PermissionMapView       = "map:view"       // 地图查看
	PermissionTaskManage    = "task:manage"    // 任务管理
	PermissionTaskView      = "task:view"      // 任务查看
	PermissionDeviceManage  = "device:manage"  // 设备管理（创建/编辑/删除）
	PermissionDeviceView    = "device:view"    // 设备查看
	PermissionOperationView = "operation:view" // 操作记录查看
)

// GetRolePermissions 获取角色对应的权限列表
func GetRolePermissions(role entity.RoleType) []string {
	switch role {
	case entity.RoleAdministrator:
		// 超级管理员：所有权限
		return []string{
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
	case entity.RoleManager:
		// 管理员：除用户管理外的所有权限
		return []string{
			PermissionUserView,
			PermissionMapManage,
			PermissionMapView,
			PermissionTaskManage,
			PermissionTaskView,
			PermissionDeviceManage,
			PermissionDeviceView,
			PermissionOperationView,
		}
	case entity.RoleOperator:
		// 操作员：地图查看、任务管理、设备管理、操作记录查看
		return []string{
			PermissionUserView,
			PermissionMapView,
			PermissionTaskManage,
			PermissionTaskView,
			PermissionDeviceManage,
			PermissionDeviceView,
			PermissionOperationView,
		}
	case entity.RoleUser:
		// 普通用户：仅查看
		return []string{
			PermissionUserView,
			PermissionDeviceView,
			PermissionTaskView,
			PermissionOperationView,
			PermissionMapView,
		}
	default:
		return []string{}
	}
}

// HasPermission 检查角色是否拥有指定权限
func HasPermission(role entity.RoleType, permission string) bool {
	permissions := GetRolePermissions(role)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission 检查角色是否拥有任意一个指定权限
func HasAnyPermission(role entity.RoleType, requiredPermissions ...string) bool {
	if len(requiredPermissions) == 0 {
		return true
	}

	rolePermissions := GetRolePermissions(role)
	permissionMap := make(map[string]bool)
	for _, p := range rolePermissions {
		permissionMap[p] = true
	}

	for _, required := range requiredPermissions {
		if permissionMap[required] {
			return true
		}
	}

	return false
}

// HasAllPermissions 检查角色是否拥有所有指定权限
func HasAllPermissions(role entity.RoleType, requiredPermissions ...string) bool {
	if len(requiredPermissions) == 0 {
		return true
	}

	rolePermissions := GetRolePermissions(role)
	permissionMap := make(map[string]bool)
	for _, p := range rolePermissions {
		permissionMap[p] = true
	}

	for _, required := range requiredPermissions {
		if !permissionMap[required] {
			return false
		}
	}

	return true
}
