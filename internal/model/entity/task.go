package entity

import "gorm.io/gorm"

// Task 任务编排表
type Task struct {
	gorm.Model
	SemanticMapID uint        `gorm:"not null;comment:对应的语义地图id;index"`
	SemanticMap   SemanticMap `gorm:"foreignKey:SemanticMapID"`
	UserName      string      `gorm:"type:text;not null;comment:编辑人员"`
	TaskInfo      string      `gorm:"type:text;comment:任务信息"`
	Status        *TaskStatus `gorm:"type:text;default:'pending';comment:任务状态"`
	ExtraInfo     *string     `gorm:"type:text;comment:扩展信息(JSON)"`
}

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   // 待执行
	TaskStatusRunning   TaskStatus = "running"   // 执行中
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 失败
	TaskStatusCancelled TaskStatus = "cancelled" // 已取消
)

func (Task) TableName() string {
	return "task"
}
