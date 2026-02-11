package handler

import (
	"strconv"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserOperationHandler 操作记录处理器
type UserOperationHandler struct {
	operationService *service.UserOperationService
}

func NewUserOperationHandler(operationService *service.UserOperationService) *UserOperationHandler {
	return &UserOperationHandler{
		operationService: operationService,
	}
}

// GetOperation 获取操作记录
// @Summary 获取操作记录
// @Description 根据ID获取操作记录信息
// @Tags 操作记录
// @Accept json
// @Produce json
// @Param id path int true "操作记录ID"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "操作记录不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /operations/{id} [get]
// @Security BearerAuth
func (h *UserOperationHandler) GetOperation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid operation id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的操作记录ID")
		return
	}

	logger.Info("handling get operation request", zap.Uint("id", uint(id)))

	operation, err := h.operationService.GetOperationByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("failed to get operation", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "获取操作记录失败: "+err.Error())
		return
	}

	if operation == nil {
		logger.Warn("operation not found", zap.Uint("id", uint(id)))
		NotFound(c, "操作记录不存在")
		return
	}

	Success(c, operation)
}

// ListOperations 查询操作记录列表
// @Summary 查询操作记录列表
// @Description 查询所有操作记录
// @Tags 操作记录
// @Accept json
// @Produce json
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器错误"
// @Router /operations [get]
// @Security BearerAuth
func (h *UserOperationHandler) ListOperations(c *gin.Context) {
	logger.Info("handling list operations request")

	var pageReq dto.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		logger.Error("invalid pagination parameters", zap.Error(err))
		BadRequest(c, "无效的分页参数: "+err.Error())
		return
	}

	operations, err := h.operationService.ListOperations(c.Request.Context(), pageReq)
	if err != nil {
		logger.Error("failed to list operations", zap.Error(err))
		InternalServerError(c, "查询操作记录列表失败: "+err.Error())
		return
	}

	Success(c, operations)
}
