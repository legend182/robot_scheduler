package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"robot_scheduler/internal/config"
	dao "robot_scheduler/internal/dao/interfaces"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/minio_client"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/model/entity"

	"go.uber.org/zap"
)

// PCDFileService 点云地图服务
type PCDFileService struct {
	pcdDAO dao.PCDFileDAO
}

func NewPCDFileService(pcdDAO dao.PCDFileDAO) *PCDFileService {
	return &PCDFileService{
		pcdDAO: pcdDAO,
	}
}

// CreatePCDFile 创建点云地图
func (s *PCDFileService) CreatePCDFile(ctx context.Context, req *dto.PCDFileCreateRequest) (*dto.PCDFileResponse, error) {
	logger.Info("creating pcd file in service", zap.String("name", req.Name))

	// 检查名称是否已存在
	existingFile, err := s.pcdDAO.FindByName(ctx, req.Name)
	if err != nil {
		logger.Error("failed to check pcd file name existence", zap.Error(err), zap.String("name", req.Name))
		return nil, err
	}

	if existingFile != nil {
		logger.Warn("pcd file name already exists", zap.String("name", req.Name))
		return nil, errors.New("pcd file name already exists")
	}

	// 创建点云地图实体
	file := &entity.PCDFile{
		Name:      req.Name,
		Area:      req.Area,
		Path:      req.Path,
		UserName:  req.UserName,
		Size:      req.Size,
		MinioPath: req.MinioPath,
		ExtraInfo: req.ExtraInfo,
	}

	// 保存到数据库
	if err := s.pcdDAO.Create(ctx, file); err != nil {
		logger.Error("failed to create pcd file in service", zap.Error(err), zap.String("name", req.Name))
		return nil, err
	}

	logger.Info("pcd file created successfully in service", zap.String("name", req.Name), zap.Uint("id", file.ID))
	return dto.NewPCDFileResponseFromEntity(file), nil
}

// UpdatePCDFile 更新点云地图
func (s *PCDFileService) UpdatePCDFile(ctx context.Context, id uint, req *dto.PCDFileUpdateRequest) error {
	logger.Info("updating pcd file in service", zap.Uint("id", id))

	// 获取点云地图
	file, err := s.pcdDAO.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find pcd file for update", zap.Error(err), zap.Uint("id", id))
		return err
	}

	if file == nil {
		logger.Warn("pcd file not found for update", zap.Uint("id", id))
		return errors.New("pcd file not found")
	}

	// 更新字段
	if req.Name != nil {
		file.Name = *req.Name
	}
	if req.Area != nil {
		file.Area = *req.Area
	}
	if req.Path != nil {
		file.Path = *req.Path
	}
	if req.UserName != nil {
		file.UserName = *req.UserName
	}
	if req.Size != nil {
		file.Size = *req.Size
	}
	if req.MinioPath != nil {
		file.MinioPath = req.MinioPath
	}
	if req.ExtraInfo != nil {
		file.ExtraInfo = req.ExtraInfo
	}

	// 保存更新
	if err := s.pcdDAO.Update(ctx, file); err != nil {
		logger.Error("failed to update pcd file in service", zap.Error(err), zap.Uint("id", id))
		return err
	}

	logger.Info("pcd file updated successfully in service", zap.Uint("id", id))
	return nil
}

// DeletePCDFile 删除点云地图
func (s *PCDFileService) DeletePCDFile(ctx context.Context, id uint) error {
	logger.Info("deleting pcd file in service", zap.Uint("id", id))
	return s.pcdDAO.Delete(ctx, id)
}

// GetPCDFileByID 根据ID获取点云地图
func (s *PCDFileService) GetPCDFileByID(ctx context.Context, id uint) (*dto.PCDFileResponse, error) {
	logger.Debug("getting pcd file by id in service", zap.Uint("id", id))
	file, err := s.pcdDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, nil
	}
	return dto.NewPCDFileResponseFromEntity(file), nil
}

// ListPCDFiles 分页获取点云地图列表
func (s *PCDFileService) ListPCDFiles(ctx context.Context, req dto.PageRequest) (*dto.PCDFileListResponse, error) {
	logger.Debug("listing pcd files in service with pagination", zap.Int("page", req.Page), zap.Int("pageSize", req.PageSize))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	files, total, err := s.pcdDAO.FindPage(ctx, offset, req.PageSize)
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

	return dto.NewPCDFileListResponseFromEntities(files, page), nil
}

// GenerateUploadToken 生成点云地图上传凭证（预签名 PUT URL）
func (s *PCDFileService) GenerateUploadToken(ctx context.Context, userName string, req *dto.PCDFileUploadTokenRequest) (*dto.PCDFileUploadTokenResponse, error) {
	logger.Info("generating pcd upload token in service", zap.String("fileName", req.FileName), zap.String("userName", userName))

	cfg := config.Get()
	if cfg == nil || cfg.Minio == nil || !cfg.Minio.Enabled {
		logger.Error("minio not enabled in config")
		return nil, errors.New("minio is not enabled")
	}

	client := minio_client.Client()
	if client == nil {
		logger.Error("minio client not initialized")
		return nil, errors.New("minio client not initialized")
	}

	if userName == "" {
		userName = "unknown"
	}

	objectKey := fmt.Sprintf("pcd/%s/%d_%s", userName, time.Now().Unix(), req.FileName)
	expire := 10 * time.Minute

	url, err := client.PresignedPutObject(ctx, cfg.Minio.BucketName, objectKey, expire)
	if err != nil {
		logger.Error("failed to generate presigned put url", zap.Error(err))
		return nil, err
	}

	resp := &dto.PCDFileUploadTokenResponse{
		UploadURL: url.String(),
		Bucket:    cfg.Minio.BucketName,
		ObjectKey: objectKey,
		ExpireAt:  time.Now().Add(expire).Unix(),
	}

	logger.Info("pcd upload token generated successfully",
		zap.String("fileName", req.FileName),
		zap.String("userName", userName),
		zap.String("bucket", resp.Bucket),
		zap.String("objectKey", resp.ObjectKey),
	)

	return resp, nil
}
