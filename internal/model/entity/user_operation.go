package entity

import "time"

// OperationType 操作类型枚举
type OperationType string

const (
	OperationCreate OperationType = "create" // 创建
	OperationUpdate OperationType = "update" // 更新
	OperationDelete OperationType = "delete" // 删除
	OperationQuery  OperationType = "query"  // 查询
	OperationLogin  OperationType = "login"  // 登录
	OperationLogout OperationType = "logout" // 登出
)

// UserOperation 用户操作记录表
type UserOperation struct {
	ID         uint          `gorm:"primarykey;comment:主键ID"`
	UserName   string        `gorm:"type:text;not null;comment:操作人员"`
	Operation  OperationType `gorm:"type:text;not null;comment:操作类型"`
	Module     string        `gorm:"type:text;not null;comment:操作模块"`
	TargetID   *uint         `gorm:"comment:目标ID"`
	TargetName *string       `gorm:"type:text;comment:目标名称"`
	IP         *string       `gorm:"type:text;comment:操作IP"`
	UserAgent  *string       `gorm:"type:text;comment:用户代理"`
	ExtraInfo  string        `gorm:"type:text;not null;comment:操作信息(JSON)"`
	CreateTime *time.Time    `gorm:"type:datetime;autoCreateTime;comment:操作时间"`
}

func (UserOperation) TableName() string {
	return "user_operation"
}
