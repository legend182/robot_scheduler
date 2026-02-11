package dao

import (
	"context"
	"robot_scheduler/internal/model/entity"
)

// TaskDAO 任务数据访问接口
type TaskDAO interface {
	// Create 创建任务
	Create(ctx context.Context, task *entity.Task) error

	// Update 更新任务
	Update(ctx context.Context, task *entity.Task) error

	// Delete 删除任务(软删除)
	Delete(ctx context.Context, id uint) error

	// FindByID 根据ID查询任务
	FindByID(ctx context.Context, id uint) (*entity.Task, error)

	// FindAll 查询所有任务
	FindAll(ctx context.Context) ([]*entity.Task, error)

	// FindPage 分页查询任务
	FindPage(ctx context.Context, offset, limit int) ([]*entity.Task, int64, error)
}
