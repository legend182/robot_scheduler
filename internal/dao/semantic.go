package impl

import (
	"context"
	"errors"
	dao "robot_scheduler/internal/dao/interfaces"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SemanticMapDAOImpl struct {
	db *gorm.DB
}

func NewSemanticMapDAO(db *gorm.DB) dao.SemanticMapDAO {
	return &SemanticMapDAOImpl{db: db}
}

func (d *SemanticMapDAOImpl) Create(ctx context.Context, semanticMap *entity.SemanticMap) error {
	logger.Info("creating semantic map")

	if err := d.db.WithContext(ctx).Create(semanticMap).Error; err != nil {
		logger.Error("failed to create semantic map", zap.Error(err))
		return err
	}

	logger.Info("semantic map created successfully", zap.Uint("id", semanticMap.ID))
	return nil
}

func (d *SemanticMapDAOImpl) Update(ctx context.Context, semanticMap *entity.SemanticMap) error {
	logger.Info("updating semantic map", zap.Uint("id", semanticMap.ID))

	result := d.db.WithContext(ctx).Save(semanticMap)
	if err := result.Error; err != nil {
		logger.Error("failed to update semantic map", zap.Error(err), zap.Uint("id", semanticMap.ID))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("semantic map not found for update", zap.Uint("id", semanticMap.ID))
		return errors.New("semantic map not found")
	}

	logger.Info("semantic map updated successfully", zap.Uint("id", semanticMap.ID))
	return nil
}

func (d *SemanticMapDAOImpl) Delete(ctx context.Context, id uint) error {
	logger.Info("deleting semantic map", zap.Uint("id", id))

	result := d.db.WithContext(ctx).Delete(&entity.SemanticMap{}, id)
	if err := result.Error; err != nil {
		logger.Error("failed to delete semantic map", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("semantic map not found for deletion", zap.Uint("id", id))
		return errors.New("semantic map not found")
	}

	logger.Info("semantic map deleted successfully", zap.Uint("id", id))
	return nil
}

func (d *SemanticMapDAOImpl) FindByID(ctx context.Context, id uint) (*entity.SemanticMap, error) {
	logger.Debug("finding semantic map by id", zap.Uint("id", id))

	var semanticMap entity.SemanticMap
	err := d.db.WithContext(ctx).Preload("PCDFile").First(&semanticMap, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("semantic map not found", zap.Uint("id", id))
			return nil, nil
		}
		logger.Error("failed to find semantic map by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	logger.Debug("semantic map found", zap.Uint("id", id))
	return &semanticMap, nil
}

func (d *SemanticMapDAOImpl) FindAll(ctx context.Context) ([]*entity.SemanticMap, error) {
	logger.Debug("finding all semantic maps")

	var semanticMaps []*entity.SemanticMap
	err := d.db.WithContext(ctx).Preload("PCDFile").Find(&semanticMaps).Error
	if err != nil {
		logger.Error("failed to find all semantic maps", zap.Error(err))
		return nil, err
	}

	logger.Debug("found semantic maps", zap.Int("count", len(semanticMaps)))
	return semanticMaps, nil
}

// FindPage 分页查询语义地图
func (d *SemanticMapDAOImpl) FindPage(ctx context.Context, offset, limit int) ([]*entity.SemanticMap, int64, error) {
	logger.Debug("finding semantic maps with pagination", zap.Int("offset", offset), zap.Int("limit", limit))

	var (
		semanticMaps []*entity.SemanticMap
		total        int64
	)

	db := d.db.WithContext(ctx).Model(&entity.SemanticMap{}).Preload("PCDFile")

	if err := db.Count(&total).Error; err != nil {
		logger.Error("failed to count semantic maps for pagination", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.SemanticMap{}, 0, nil
	}

	if err := db.Offset(offset).Limit(limit).Find(&semanticMaps).Error; err != nil {
		logger.Error("failed to find semantic maps with pagination", zap.Error(err))
		return nil, 0, err
	}

	logger.Debug("found semantic maps with pagination", zap.Int("count", len(semanticMaps)), zap.Int64("total", total))
	return semanticMaps, total, nil
}