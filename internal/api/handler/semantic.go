package handler

import (
	"strconv"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SemanticMapHandler 语义地图处理器
type SemanticMapHandler struct {
	semanticService *service.SemanticMapService
}

func NewSemanticMapHandler(semanticService *service.SemanticMapService) *SemanticMapHandler {
	return &SemanticMapHandler{
		semanticService: semanticService,
	}
}

// CreateSemanticMap 创建语义地图
// @Summary 创建语义地图
// @Description 创建新语义地图
// @Tags 语义地图
// @Accept json
// @Produce json
// @Param request body dto.SemanticMapCreateRequest true "语义地图信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/semantic-maps [post]
// @Security BearerAuth
func (h *SemanticMapHandler) CreateSemanticMap(c *gin.Context) {
	logger.Info("handling create semantic map request")

	var req dto.SemanticMapCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	semanticMap, err := h.semanticService.CreateSemanticMap(c.Request.Context(), &req)
	if err != nil {
		logger.Error("failed to create semantic map", zap.Error(err))
		InternalServerError(c, "创建语义地图失败: "+err.Error())
		return
	}

	Success(c, semanticMap)
}

// GetSemanticMap 获取语义地图
// @Summary 获取语义地图
// @Description 根据ID获取语义地图信息
// @Tags 语义地图
// @Accept json
// @Produce json
// @Param id path int true "语义地图ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "语义地图不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/semantic-maps/{id} [get]
// @Security BearerAuth
func (h *SemanticMapHandler) GetSemanticMap(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid semantic map id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的语义地图ID")
		return
	}

	logger.Info("handling get semantic map request", zap.Uint("id", uint(id)))

	semanticMap, err := h.semanticService.GetSemanticMapByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("failed to get semantic map", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "获取语义地图失败: "+err.Error())
		return
	}

	if semanticMap == nil {
		logger.Warn("semantic map not found", zap.Uint("id", uint(id)))
		NotFound(c, "语义地图不存在")
		return
	}

	Success(c, semanticMap)
}

// UpdateSemanticMap 更新语义地图
// @Summary 更新语义地图
// @Description 更新语义地图信息
// @Tags 语义地图
// @Accept json
// @Produce json
// @Param id path int true "语义地图ID"
// @Param request body dto.SemanticMapUpdateRequest true "更新信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "语义地图不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/semantic-maps/{id} [put]
// @Security BearerAuth
func (h *SemanticMapHandler) UpdateSemanticMap(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid semantic map id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的语义地图ID")
		return
	}

	var req dto.SemanticMapUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	logger.Info("handling update semantic map request", zap.Uint("id", uint(id)))

	if err := h.semanticService.UpdateSemanticMap(c.Request.Context(), uint(id), &req); err != nil {
		logger.Error("failed to update semantic map", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "更新语义地图失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "更新成功"})
}

// DeleteSemanticMap 删除语义地图
// @Summary 删除语义地图
// @Description 删除语义地图（软删除）
// @Tags 语义地图
// @Accept json
// @Produce json
// @Param id path int true "语义地图ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/semantic-maps/{id} [delete]
// @Security BearerAuth
func (h *SemanticMapHandler) DeleteSemanticMap(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid semantic map id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的语义地图ID")
		return
	}

	logger.Info("handling delete semantic map request", zap.Uint("id", uint(id)))

	if err := h.semanticService.DeleteSemanticMap(c.Request.Context(), uint(id)); err != nil {
		logger.Error("failed to delete semantic map", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "删除语义地图失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "删除成功"})
}

// ListSemanticMaps 查询语义地图列表
// @Summary 查询语义地图列表
// @Description 查询所有语义地图
// @Tags 语义地图
// @Accept json
// @Produce json
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器错误"
// @Router /maps/semantic-maps [get]
// @Security BearerAuth
func (h *SemanticMapHandler) ListSemanticMaps(c *gin.Context) {
	logger.Info("handling list semantic maps request")

	var pageReq dto.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		logger.Error("invalid pagination parameters", zap.Error(err))
		BadRequest(c, "无效的分页参数: "+err.Error())
		return
	}

	semanticMaps, err := h.semanticService.ListSemanticMaps(c.Request.Context(), pageReq)
	if err != nil {
		logger.Error("failed to list semantic maps", zap.Error(err))
		InternalServerError(c, "查询语义地图列表失败: "+err.Error())
		return
	}

	Success(c, semanticMaps)
}
