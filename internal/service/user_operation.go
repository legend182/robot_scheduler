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

// UserOperationService 操作记录服务
type UserOperationService struct {
	operationDAO dao.UserOperationDAO
}

func NewUserOperationService(operationDAO dao.UserOperationDAO) *UserOperationService {
	return &UserOperationService{
		operationDAO: operationDAO,
	}
}

// CreateOperation 创建操作记录
func (s *UserOperationService) CreateOperation(ctx context.Context, req *dto.UserOperationCreateRequest) (*dto.UserOperationResponse, error) {
	logger.Info("creating operation record in service", zap.String("username", req.UserName))

	// 创建操作记录实体
	operation := &entity.UserOperation{
		UserName:   req.UserName,
		Operation:  req.Operation,
		Module:     req.Module,
		TargetID:   req.TargetID,
		TargetName: req.TargetName,
		IP:         req.IP,
		UserAgent:  req.UserAgent,
		ExtraInfo:  req.ExtraInfo,
	}

	// 保存到数据库
	if err := s.operationDAO.Create(ctx, operation); err != nil {
		logger.Error("failed to create operation record in service", zap.Error(err))
		return nil, err
	}

	logger.Info("operation record created successfully in service", zap.Uint("id", operation.ID))
	return dto.NewUserOperationResponseFromEntity(operation), nil
}

// UpdateOperation 更新操作记录
func (s *UserOperationService) UpdateOperation(ctx context.Context, id uint, req *dto.UserOperationUpdateRequest) error {
	logger.Info("updating operation record in service", zap.Uint("id", id))

	// 获取操作记录
	operation, err := s.operationDAO.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find operation record for update", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if operation == nil {
		logger.Warn("operation record not found for update", zap.Uint("id", id))
		return errors.New("operation record not found")
	}

	// 更新字段
	if req.UserName != nil {
		operation.UserName = *req.UserName
	}
	if req.Operation != nil {
		operation.Operation = *req.Operation
	}
	if req.Module != nil {
		operation.Module = *req.Module
	}
	if req.TargetID != nil {
		operation.TargetID = req.TargetID
	}
	if req.TargetName != nil {
		operation.TargetName = req.TargetName
	}
	if req.IP != nil {
		operation.IP = req.IP
	}
	if req.UserAgent != nil {
		operation.UserAgent = req.UserAgent
	}
	if req.ExtraInfo != nil {
		operation.ExtraInfo = *req.ExtraInfo
	}

	// 保存更新
	if err := s.operationDAO.Update(ctx, operation); err != nil {
		logger.Error("failed to update operation record in service", zap.Error(err), zap.Uint("id", id))
		return err
	}

	logger.Info("operation record updated successfully in service", zap.Uint("id", id))
	return nil
}

// DeleteOperation 删除操作记录
func (s *UserOperationService) DeleteOperation(ctx context.Context, id uint) error {
	logger.Info("deleting operation record in service", zap.Uint("id", id))
	return s.operationDAO.Delete(ctx, id)
}

// GetOperationByID 根据ID获取操作记录
func (s *UserOperationService) GetOperationByID(ctx context.Context, id uint) (*dto.UserOperationResponse, error) {
	logger.Debug("getting operation record by id in service", zap.Uint("id", id))
	operation, err := s.operationDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if operation == nil {
		return nil, nil
	}
	return dto.NewUserOperationResponseFromEntity(operation), nil
}

// ListOperations 分页获取操作记录列表
func (s *UserOperationService) ListOperations(ctx context.Context, req dto.PageRequest) (*dto.UserOperationListResponse, error) {
	logger.Debug("listing operation records in service with pagination", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	operations, total, err := s.operationDAO.FindPage(ctx, offset, req.PageSize)
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

	return dto.NewUserOperationListResponseFromEntities(operations, page), nil
}
