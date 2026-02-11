package dao

import (
	"context"
	"robot_scheduler/internal/model/entity"
)

// UserOperationDAO 用户操作记录数据访问接口
type UserOperationDAO interface {
	// Create 创建操作记录
	Create(ctx context.Context, operation *entity.UserOperation) error

	// Update 更新操作记录
	Update(ctx context.Context, operation *entity.UserOperation) error

	// Delete 删除操作记录
	Delete(ctx context.Context, id uint) error

	// FindByID 根据ID查询操作记录
	FindByID(ctx context.Context, id uint) (*entity.UserOperation, error)

	// FindAll 查询所有操作记录
	FindAll(ctx context.Context) ([]*entity.UserOperation, error)

	// FindPage 分页查询操作记录
	FindPage(ctx context.Context, offset, limit int) ([]*entity.UserOperation, int64, error)
}
