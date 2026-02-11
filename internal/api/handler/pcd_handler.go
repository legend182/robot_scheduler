package handler

import (
	"strconv"

	"robot_scheduler/internal/api/middleware"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PCDFileHandler 点云地图处理器
type PCDFileHandler struct {
	pcdService *service.PCDFileService
}

func NewPCDFileHandler(pcdService *service.PCDFileService) *PCDFileHandler {
	return &PCDFileHandler{
		pcdService: pcdService,
	}
}

// GetPCDUploadToken 获取点云地图上传凭证
// @Summary 获取点云地图上传凭证
// @Description 返回 MinIO 预签名上传 URL
// @Tags 点云地图
// @Accept json
// @Produce json
// @Param request body dto.PCDFileUploadTokenRequest true "上传文件信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/pcd-files/upload-token [post]
// @Security BearerAuth
func (h *PCDFileHandler) GetPCDUploadToken(c *gin.Context) {
	logger.Info("handling get pcd upload token request")

	var req dto.PCDFileUploadTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid upload token request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 从 JWT 上下文中获取用户名
	userNameValue, _ := c.Get(string(middleware.UserNameKey))
	userName, _ := userNameValue.(string)

	resp, err := h.pcdService.GenerateUploadToken(c.Request.Context(), userName, &req)
	if err != nil {
		logger.Error("failed to generate pcd upload token", zap.Error(err))
		InternalServerError(c, "生成上传凭证失败: "+err.Error())
		return
	}

	Success(c, resp)
}

// CreatePCDFile 创建点云地图
// @Summary 创建点云地图
// @Description 创建新点云地图
// @Tags 点云地图
// @Accept json
// @Produce json
// @Param request body dto.PCDFileCreateRequest true "点云地图信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/pcd-files [post]
// @Security BearerAuth
func (h *PCDFileHandler) CreatePCDFile(c *gin.Context) {
	logger.Info("handling create pcd file request")

	var req dto.PCDFileCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	file, err := h.pcdService.CreatePCDFile(c.Request.Context(), &req)
	if err != nil {
		logger.Error("failed to create pcd file", zap.Error(err))
		InternalServerError(c, "创建点云地图失败: "+err.Error())
		return
	}

	Success(c, file)
}

// GetPCDFile 获取点云地图
// @Summary 获取点云地图
// @Description 根据ID获取点云地图信息
// @Tags 点云地图
// @Accept json
// @Produce json
// @Param id path int true "点云地图ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "点云地图不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/pcd-files/{id} [get]
// @Security BearerAuth
func (h *PCDFileHandler) GetPCDFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid pcd file id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的点云地图ID")
		return
	}

	logger.Info("handling get pcd file request", zap.Uint("id", uint(id)))

	file, err := h.pcdService.GetPCDFileByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("failed to get pcd file", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "获取点云地图失败: "+err.Error())
		return
	}

	if file == nil {
		logger.Warn("pcd file not found", zap.Uint("id", uint(id)))
		NotFound(c, "点云地图不存在")
		return
	}

	Success(c, file)
}

// UpdatePCDFile 更新点云地图
// @Summary 更新点云地图
// @Description 更新点云地图信息
// @Tags 点云地图
// @Accept json
// @Produce json
// @Param id path int true "点云地图ID"
// @Param request body dto.PCDFileUpdateRequest true "更新信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "点云地图不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/pcd-files/{id} [put]
// @Security BearerAuth
func (h *PCDFileHandler) UpdatePCDFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid pcd file id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的点云地图ID")
		return
	}

	var req dto.PCDFileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	logger.Info("handling update pcd file request", zap.Uint("id", uint(id)))

	if err := h.pcdService.UpdatePCDFile(c.Request.Context(), uint(id), &req); err != nil {
		logger.Error("failed to update pcd file", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "更新点云地图失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "更新成功"})
}

// DeletePCDFile 删除点云地图
// @Summary 删除点云地图
// @Description 删除点云地图（软删除）
// @Tags 点云地图
// @Accept json
// @Produce json
// @Param id path int true "点云地图ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/pcd-files/{id} [delete]
// @Security BearerAuth
func (h *PCDFileHandler) DeletePCDFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid pcd file id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的点云地图ID")
		return
	}

	logger.Info("handling delete pcd file request", zap.Uint("id", uint(id)))

	if err := h.pcdService.DeletePCDFile(c.Request.Context(), uint(id)); err != nil {
		logger.Error("failed to delete pcd file", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "删除点云地图失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "删除成功"})
}

// ListPCDFiles 查询点云地图列表
// @Summary 查询点云地图列表
// @Description 查询所有点云地图
// @Tags 点云地图
// @Accept json
// @Produce json
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/pcd-files [get]
// @Security BearerAuth
func (h *PCDFileHandler) ListPCDFiles(c *gin.Context) {
	logger.Info("handling list pcd files request")

	var pageReq dto.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		logger.Error("invalid pagination parameters", zap.Error(err))
		BadRequest(c, "无效的分页参数: "+err.Error())
		return
	}

	files, err := h.pcdService.ListPCDFiles(c.Request.Context(), pageReq)
	if err != nil {
		logger.Error("failed to list pcd files", zap.Error(err))
		InternalServerError(c, "查询点云地图列表失败: "+err.Error())
		return
	}

	Success(c, files)
}
