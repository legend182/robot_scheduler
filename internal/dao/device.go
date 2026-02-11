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

type DeviceDAOImpl struct {
	db *gorm.DB
}

func NewDeviceDAO(db *gorm.DB) dao.DeviceDAO {
	return &DeviceDAOImpl{db: db}
}

func (d *DeviceDAOImpl) Create(ctx context.Context, device *entity.Device) error {
	logger.Info("creating device", zap.String("type", string(device.Type)))

	if err := d.db.WithContext(ctx).Create(device).Error; err != nil {
		logger.Error("failed to create device", zap.Error(err))
		return err
	}

	logger.Info("device created successfully", zap.Uint("id", device.ID))
	return nil
}

func (d *DeviceDAOImpl) Update(ctx context.Context, device *entity.Device) error {
	logger.Info("updating device", zap.Uint("id", device.ID))

	result := d.db.WithContext(ctx).Save(device)
	if err := result.Error; err != nil {
		logger.Error("failed to update device", zap.Error(err), zap.Uint("id", device.ID))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("device not found for update", zap.Uint("id", device.ID))
		return errors.New("device not found")
	}

	logger.Info("device updated successfully", zap.Uint("id", device.ID))
	return nil
}

func (d *DeviceDAOImpl) Delete(ctx context.Context, id uint) error {
	logger.Info("deleting device", zap.Uint("id", id))

	result := d.db.WithContext(ctx).Delete(&entity.Device{}, id)
	if err := result.Error; err != nil {
		logger.Error("failed to delete device", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Warn("device not found for deletion", zap.Uint("id", id))
		return errors.New("device not found")
	}

	logger.Info("device deleted successfully", zap.Uint("id", id))
	return nil
}

func (d *DeviceDAOImpl) FindByID(ctx context.Context, id uint) (*entity.Device, error) {
	logger.Debug("finding device by id", zap.Uint("id", id))

	var device entity.Device
	err := d.db.WithContext(ctx).First(&device, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Debug("device not found", zap.Uint("id", id))
			return nil, nil
		}
		logger.Error("failed to find device by id", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}

	logger.Debug("device found", zap.Uint("id", id))
	return &device, nil
}

func (d *DeviceDAOImpl) FindAll(ctx context.Context) ([]*entity.Device, error) {
	logger.Debug("finding all devices")

	var devices []*entity.Device
	err := d.db.WithContext(ctx).Find(&devices).Error
	if err != nil {
		logger.Error("failed to find all devices", zap.Error(err))
		return nil, err
	}

	logger.Debug("found devices", zap.Int("count", len(devices)))
	return devices, nil
}

// FindPage 分页查询设备
func (d *DeviceDAOImpl) FindPage(ctx context.Context, offset, limit int) ([]*entity.Device, int64, error) {
	logger.Debug("finding devices with pagination", zap.Int("offset", offset), zap.Int("limit", limit))

	var (
		devices []*entity.Device
		total   int64
	)

	db := d.db.WithContext(ctx).Model(&entity.Device{})

	if err := db.Count(&total).Error; err != nil {
		logger.Error("failed to count devices for pagination", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []*entity.Device{}, 0, nil
	}

	if err := db.Offset(offset).Limit(limit).Find(&devices).Error; err != nil {
		logger.Error("failed to find devices with pagination", zap.Error(err))
		return nil, 0, err
	}

	logger.Debug("found devices with pagination", zap.Int("count", len(devices)), zap.Int64("total", total))
	return devices, total, nil
}