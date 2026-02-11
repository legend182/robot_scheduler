package dao

import (
	"context"
	"robot_scheduler/internal/model/entity"
)

// SemanticMapDAO 语义地图数据访问接口
type SemanticMapDAO interface {
	// Create 创建语义地图
	Create(ctx context.Context, semanticMap *entity.SemanticMap) error

	// Update 更新语义地图
	Update(ctx context.Context, semanticMap *entity.SemanticMap) error

	// Delete 删除语义地图(软删除)
	Delete(ctx context.Context, id uint) error

	// FindByID 根据ID查询语义地图
	FindByID(ctx context.Context, id uint) (*entity.SemanticMap, error)

	// FindAll 查询所有语义地图
	FindAll(ctx context.Context) ([]*entity.SemanticMap, error)

	// FindPage 分页查询语义地图
	FindPage(ctx context.Context, offset, limit int) ([]*entity.SemanticMap, int64, error)
}
