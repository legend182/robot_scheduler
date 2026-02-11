package dao

import (
	"context"
	"robot_scheduler/internal/model/entity"
)

// UserDAO 用户数据访问接口
type UserDAO interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.UserInfo) error
	
	// Update 更新用户
	Update(ctx context.Context, user *entity.UserInfo) error
	
	// Delete 删除用户(软删除)
	Delete(ctx context.Context, id uint) error
	
	// FindByID 根据ID查询用户
	FindByID(ctx context.Context, id uint) (*entity.UserInfo, error)
	
	// FindAll 查询所有用户
	FindAll(ctx context.Context) ([]*entity.UserInfo, error)

	// FindPage 分页查询用户
	FindPage(ctx context.Context, offset, limit int) ([]*entity.UserInfo, int64, error)
	
	// FindByUserName 根据用户名查询用户
	FindByUserName(ctx context.Context, userName string) (*entity.UserInfo, error)
}