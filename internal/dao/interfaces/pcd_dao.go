package dao

import (
	"context"
	"robot_scheduler/internal/model/entity"
)

// PCDFileDAO 点云地图数据访问接口
type PCDFileDAO interface {
	// Create 创建点云地图
	Create(ctx context.Context, file *entity.PCDFile) error

	// Update 更新点云地图
	Update(ctx context.Context, file *entity.PCDFile) error

	// Delete 删除点云地图(软删除)
	Delete(ctx context.Context, id uint) error

	// FindByID 根据ID查询点云地图
	FindByID(ctx context.Context, id uint) (*entity.PCDFile, error)

	// FindAll 查询所有点云地图
	FindAll(ctx context.Context) ([]*entity.PCDFile, error)

	// FindByName 根据名称查询
	FindByName(ctx context.Context, name string) (*entity.PCDFile, error)

	// FindPage 分页查询点云地图
	FindPage(ctx context.Context, offset, limit int) ([]*entity.PCDFile, int64, error)
}
