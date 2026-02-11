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

type UserOperationDAOImpl struct {
	db *gorm.DB
}

func NewUserOperationDAO(db *gorm.DB) dao.UserOperationDAO {
	return &UserOperationDAOImpl{db: db}
}

func (d *UserOperationDAOImpl) Create(ctx context.Context, operation *entity.UserOperation) error {
	logger.Info("creating operation record", zap.String("username", operation.UserName))

	if err := d.db.WithContext(ctx).Create(operation).Error; err != nil {
		logger.Error("failed to create operation record", zap.Error(err))
		return err
	}

	logger.Info("operation record created successfully", zap.Uint("id", operation.ID))
	return nil
}

func (d *UserOperationDAOImpl) Update(ctx context.Context, operation *entity.UserOperation) error {
	logger.Info("updating operation record", zap.Uint("id", operation.ID))

	result := d.db.WithContext(ctx).Save(operation)
	if err := result.Error; err != nil {
		logger.Error("failed to update operation record", zap.Error(err), zap.Uint("id", operation.ID))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("operation record not found for update", zap.Uint("id", operation.ID))
		return errors.New("operation record not found")
	}

	logger.Info("operation record updated successfully", zap.Uint("id", operation.ID))
	return nil
}

func (d *UserOperationDAOImpl) Delete(ctx context.Context, id uint) error {
	logger.Info("deleting operation record", zap.Uint("id", id))

	result := d.db.WithContext(ctx).Delete(&entity.UserOperation{}, id)
	if err := result.Error; err != nil {
		logger.Error("failed to delete operation record", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("operation record not found for deletion", zap.Uint("id", id))
		return errors.New("operation record not found")
	}

	logger.Info("operation record deleted successfully", zap.Uint("id", id))
	return nil
}

func (d *UserOperationDAOImpl) FindByID(ctx context.Context, id uint) (*entity.UserOperation, error) {
	logger.Debug("finding operation record by id", zap.Uint("id", id))

	var operation entity.UserOperation
	err := d.db.WithContext(ctx).First(&operation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("operation record not found", zap.Uint("id", id))
			return nil, nil
		}
		logger.Error("failed to find operation record by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	logger.Debug("operation record found", zap.Uint("id", id))
	return &operation, nil
}

func (d *UserOperationDAOImpl) FindAll(ctx context.Context) ([]*entity.UserOperation, error) {
	logger.Debug("finding all operation records")

	var operations []*entity.UserOperation
	err := d.db.WithContext(ctx).Order("create_time desc").Find(&operations).Error
	if err != nil {
		logger.Error("failed to find all operation records", zap.Error(err))
		return nil, err
	}

	logger.Debug("found operation records", zap.Int("count", len(operations)))
	return operations, nil
}

// FindPage 分页查询操作记录
func (d *UserOperationDAOImpl) FindPage(ctx context.Context, offset, limit int) ([]*entity.UserOperation, int64, error) {
	logger.Debug("finding operation records with pagination", zap.Int("offset", offset), zap.Int("limit", limit))

	var (
		operations []*entity.UserOperation
		total      int64
	)

	db := d.db.WithContext(ctx).Model(&entity.UserOperation{})

	if err := db.Count(&total).Error; err != nil {
		logger.Error("failed to count operation records for pagination", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.UserOperation{}, 0, nil
	}

	if err := db.Order("create_time desc").Offset(offset).Limit(limit).Find(&operations).Error; err != nil {
		logger.Error("failed to find operation records with pagination", zap.Error(err))
		return nil, 0, err
	}

	logger.Debug("found operation records with pagination", zap.Int("count", len(operations)), zap.Int64("total", total))
	return operations, total, nil
}