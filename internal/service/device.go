package service

import (
	"context"
	"errors"
	dao "robot_scheduler/internal/dao/interfaces"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/model/entity"

	"go.uber.org/zap"
)

// DeviceService 设备服务
type DeviceService struct {
	deviceDAO dao.DeviceDAO
}

func NewDeviceService(deviceDAO dao.DeviceDAO) *DeviceService {
	return &DeviceService{
		deviceDAO: deviceDAO,
	}
}

// CreateDevice 创建设备
func (s *DeviceService) CreateDevice(ctx context.Context, req *dto.DeviceCreateRequest) (*dto.DeviceResponse, error) {
	logger.Info("creating device in service", zap.String("type", string(req.Type)))

	// 创建设备实体
	device := &entity.Device{
		Type:      req.Type,
		Company:   req.Company,
		IP:        req.IP,
		Port:      req.Port,
		UserName:  req.UserName,
		Password:  req.Password,
		Status:    &[]entity.DeviceStatus{entity.DeviceStatusOffline}[0],
		ExtraInfo: req.ExtraInfo,
	}

	// 保存到数据库
	if err := s.deviceDAO.Create(ctx, device); err != nil {
		logger.Error("failed to create device in service", zap.Error(err))
		return nil, err
	}

	logger.Info("device created successfully in service", zap.Uint("id", device.ID))
	return dto.NewDeviceResponseFromEntity(device), nil
}

// UpdateDevice 更新设备
func (s *DeviceService) UpdateDevice(ctx context.Context, id uint, req *dto.DeviceUpdateRequest) error {
	logger.Info("updating device in service", zap.Uint("id", id))

	// 获取设备
	device, err := s.deviceDAO.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find device for update", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if device == nil {
		logger.Warn("device not found for update", zap.Uint("id", id))
		return errors.New("device not found")
	}

	// 更新字段
	if req.Type != nil {
		device.Type = *req.Type
	}
	if req.Company != nil {
		device.Company = *req.Company
	}
	if req.IP != nil {
		device.IP = req.IP
	}
	if req.Port != nil {
		device.Port = *req.Port
	}
	if req.UserName != nil {
		device.UserName = req.UserName
	}
	if req.Password != nil {
		device.Password = req.Password
	}
	if req.Status != nil {
		device.Status = req.Status
	}
	if req.ExtraInfo != nil {
		device.ExtraInfo = req.ExtraInfo
	}

	// 保存更新
	if err := s.deviceDAO.Update(ctx, device); err != nil {
		logger.Error("failed to update device in service", zap.Error(err), zap.Uint("id", id))
		return err
	}

	logger.Info("device updated successfully in service", zap.Uint("id", id))
	return nil
}

// DeleteDevice 删除设备
func (s *DeviceService) DeleteDevice(ctx context.Context, id uint) error {
	logger.Info("deleting device in service", zap.Uint("id", id))
	return s.deviceDAO.Delete(ctx, id)
}

// GetDeviceByID 根据ID获取设备
func (s *DeviceService) GetDeviceByID(ctx context.Context, id uint) (*dto.DeviceResponse, error) {
	logger.Debug("getting device by id in service", zap.Uint("id", id))
	device, err := s.deviceDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, nil
	}
	return dto.NewDeviceResponseFromEntity(device), nil
}

// ListDevices 分页获取设备列表
func (s *DeviceService) ListDevices(ctx context.Context, req dto.PageRequest) (*dto.DeviceListResponse, error) {
	logger.Debug("listing devices in service with pagination", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	devices, total, err := s.deviceDAO.FindPage(ctx, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	pages := 0
	if req.PageSize > 0 {
		pages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	}

	page := dto.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Pages:    pages,
	}

	return dto.NewDeviceListResponseFromEntities(devices, page), nil
}
