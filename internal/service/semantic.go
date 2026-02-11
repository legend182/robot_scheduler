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

// SemanticMapService 语义地图服务
type SemanticMapService struct {
	semanticDAO dao.SemanticMapDAO
}

func NewSemanticMapService(semanticDAO dao.SemanticMapDAO) *SemanticMapService {
	return &SemanticMapService{
		semanticDAO: semanticDAO,
	}
}

// CreateSemanticMap 创建语义地图
func (s *SemanticMapService) CreateSemanticMap(ctx context.Context, req *dto.SemanticMapCreateRequest) (*dto.SemanticMapResponse, error) {
	logger.Info("creating semantic map in service", zap.Uint("pcdFileID", req.PCDFileID))

	// 创建语义地图实体
	semanticMap := &entity.SemanticMap{
		PCDFileID:    req.PCDFileID,
		UserName:     req.UserName,
		SemanticInfo: req.SemanticInfo,
		ExtraInfo:    req.ExtraInfo,
	}

	// 保存到数据库
	if err := s.semanticDAO.Create(ctx, semanticMap); err != nil {
		logger.Error("failed to create semantic map in service", zap.Error(err))
		return nil, err
	}

	logger.Info("semantic map created successfully in service", zap.Uint("id", semanticMap.ID))
	return dto.NewSemanticMapResponseFromEntity(semanticMap), nil
}

// UpdateSemanticMap 更新语义地图
func (s *SemanticMapService) UpdateSemanticMap(ctx context.Context, id uint, req *dto.SemanticMapUpdateRequest) error {
	logger.Info("updating semantic map in service", zap.Uint("id", id))

	// 获取语义地图
	semanticMap, err := s.semanticDAO.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find semantic map for update", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if semanticMap == nil {
		logger.Warn("semantic map not found for update", zap.Uint("id", id))
		return errors.New("semantic map not found")
	}

	// 更新字段
	if req.PCDFileID != nil {
		semanticMap.PCDFileID = *req.PCDFileID
	}
	if req.UserName != nil {
		semanticMap.UserName = *req.UserName
	}
	if req.SemanticInfo != nil {
		semanticMap.SemanticInfo = *req.SemanticInfo
	}
	if req.ExtraInfo != nil {
		semanticMap.ExtraInfo = req.ExtraInfo
	}

	// 保存更新
	if err := s.semanticDAO.Update(ctx, semanticMap); err != nil {
		logger.Error("failed to update semantic map in service", zap.Error(err), zap.Uint("id", id))
		return err
	}

	logger.Info("semantic map updated successfully in service", zap.Uint("id", id))
	return nil
}

// DeleteSemanticMap 删除语义地图
func (s *SemanticMapService) DeleteSemanticMap(ctx context.Context, id uint) error {
	logger.Info("deleting semantic map in service", zap.Uint("id", id))
	return s.semanticDAO.Delete(ctx, id)
}

// GetSemanticMapByID 根据ID获取语义地图
func (s *SemanticMapService) GetSemanticMapByID(ctx context.Context, id uint) (*dto.SemanticMapResponse, error) {
	logger.Debug("getting semantic map by id in service", zap.Uint("id", id))
	semanticMap, err := s.semanticDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if semanticMap == nil {
		return nil, nil
	}
	return dto.NewSemanticMapResponseFromEntity(semanticMap), nil
}

// ListSemanticMaps 分页获取语义地图列表
func (s *SemanticMapService) ListSemanticMaps(ctx context.Context, req dto.PageRequest) (*dto.SemanticMapListResponse, error) {
	logger.Debug("listing semantic maps in service with pagination", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	maps, total, err := s.semanticDAO.FindPage(ctx, offset, req.PageSize)
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

	return dto.NewSemanticMapListResponseFromEntities(maps, page), nil
}
