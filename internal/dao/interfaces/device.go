package dao

import (
	"context"
	"robot_scheduler/internal/model/entity"
)

// DeviceDAO 设备数据访问接口
type DeviceDAO interface {
	// Create 创建设备
	Create(ctx context.Context, device *entity.Device) error

	// Update 更新设备
	Update(ctx context.Context, device *entity.Device) error

	// Delete 删除设备(软删除)
	Delete(ctx context.Context, id uint) error

	// FindByID 根据ID查询设备
	FindByID(ctx context.Context, id uint) (*entity.Device, error)

	// FindAll 查询所有设备
	FindAll(ctx context.Context) ([]*entity.Device, error)

	// FindPage 分页查询设备
	FindPage(ctx context.Context, offset, limit int) ([]*entity.Device, int64, error)
}
