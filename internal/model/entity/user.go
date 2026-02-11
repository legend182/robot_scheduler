package entity

import "gorm.io/gorm"

// RoleType 角色类型枚举
type RoleType string

const (
	RoleAdministrator RoleType = "administrator" // 超级管理员
	RoleManager       RoleType = "manager"       // 管理员
	RoleOperator      RoleType = "operator"      // 操作员
	RoleUser          RoleType = "user"          // 普通用户
)

// UserInfo 用户信息表
type UserInfo struct {
	gorm.Model
	UserName  string   `gorm:"type:text;not null;comment:用户名;uniqueIndex"`
	Password  string   `gorm:"type:text;comment:密码(DES加密)"`
	Role      RoleType `gorm:"type:text;comment:角色名称"`
	IsLocked  int      `gorm:"default:0;comment:是否锁定 0-不锁定 1-锁定"`
	ExtraInfo *string  `gorm:"type:text;comment:扩展信息(JSON)"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
