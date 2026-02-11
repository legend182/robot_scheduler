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

type PCDFileDAOImpl struct {
	db *gorm.DB
}

func NewPCDFileDAO(db *gorm.DB) dao.PCDFileDAO {
	return &PCDFileDAOImpl{db: db}
}

func (d *PCDFileDAOImpl) Create(ctx context.Context, file *entity.PCDFile) error {
	logger.Info("creating pcd file", zap.String("name", file.Name))

	if err := d.db.WithContext(ctx).Create(file).Error; err != nil {
		logger.Error("failed to create pcd file", zap.Error(err), zap.String("name", file.Name))
		return err
	}

	logger.Info("pcd file created successfully", zap.Uint("id", file.ID), zap.String("name", file.Name))
	return nil
}

func (d *PCDFileDAOImpl) Update(ctx context.Context, file *entity.PCDFile) error {
	logger.Info("updating pcd file", zap.Uint("id", file.ID))

	result := d.db.WithContext(ctx).Save(file)
	if err := result.Error; err != nil {
		logger.Error("failed to update pcd file", zap.Error(err), zap.Uint("id", file.ID))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("pcd file not found for update", zap.Uint("id", file.ID))
		return errors.New("pcd file not found")
	}

	logger.Info("pcd file updated successfully", zap.Uint("id", file.ID))
	return nil
}

func (d *PCDFileDAOImpl) Delete(ctx context.Context, id uint) error {
	logger.Info("deleting pcd file", zap.Uint("id", id))

	result := d.db.WithContext(ctx).Delete(&entity.PCDFile{}, id)
	if err := result.Error; err != nil {
		logger.Error("failed to delete pcd file", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("pcd file not found for deletion", zap.Uint("id", id))
		return errors.New("pcd file not found")
	}

	logger.Info("pcd file deleted successfully", zap.Uint("id", id))
	return nil
}

func (d *PCDFileDAOImpl) FindByID(ctx context.Context, id uint) (*entity.PCDFile, error) {
	logger.Debug("finding pcd file by id", zap.Uint("id", id))

	var file entity.PCDFile
	err := d.db.WithContext(ctx).First(&file, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("pcd file not found", zap.Uint("id", id))
			return nil, nil
		}
		logger.Error("failed to find pcd file by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	logger.Debug("pcd file found", zap.Uint("id", id))
	return &file, nil
}

func (d *PCDFileDAOImpl) FindAll(ctx context.Context) ([]*entity.PCDFile, error) {
	logger.Debug("finding all pcd files")

	var files []*entity.PCDFile
	err := d.db.WithContext(ctx).Find(&files).Error
	if err != nil {
		logger.Error("failed to find all pcd files", zap.Error(err))
		return nil, err
	}

	logger.Debug("found pcd files", zap.Int("count", len(files)))
	return files, nil
}

// FindByName 根据名称查询
func (d *PCDFileDAOImpl) FindByName(ctx context.Context, name string) (*entity.PCDFile, error) {
	logger.Debug("finding pcd file by name", zap.String("name", name))

	var file entity.PCDFile
	err := d.db.WithContext(ctx).Where("name = ?", name).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("pcd file not found by name", zap.String("name", name))
			return nil, nil
		}
		logger.Error("failed to find pcd file by name", zap.Error(err), zap.String("name", name))
		return nil, err
	}

	logger.Debug("pcd file found by name", zap.String("name", name))
	return &file, nil
}

// FindPage 分页查询点云地图
func (d *PCDFileDAOImpl) FindPage(ctx context.Context, offset, limit int) ([]*entity.PCDFile, int64, error) {
	logger.Debug("finding pcd files with pagination", zap.Int("offset", offset), zap.Int("limit", limit))

	var (
		files []*entity.PCDFile
		total int64
	)

	db := d.db.WithContext(ctx).Model(&entity.PCDFile{})

	if err := db.Count(&total).Error; err != nil {
		logger.Error("failed to count pcd files for pagination", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.PCDFile{}, 0, nil
	}

	if err := db.Offset(offset).Limit(limit).Find(&files).Error; err != nil {
		logger.Error("failed to find pcd files with pagination", zap.Error(err))
		return nil, 0, err
	}

	logger.Debug("found pcd files with pagination", zap.Int("count", len(files)), zap.Int64("total", total))
	return files, total, nil
}
