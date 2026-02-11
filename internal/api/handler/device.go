package handler

import (
	"strconv"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeviceHandler 设备处理器
type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler(deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// CreateDevice 创建设备
// @Summary 创建设备
// @Description 创建新设备
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param request body dto.DeviceCreateRequest true "设备信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /devices [post]
// @Security BearerAuth
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	logger.Info("handling create device request")

	var req dto.DeviceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	device, err := h.deviceService.CreateDevice(c.Request.Context(), &req)
	if err != nil {
		logger.Error("failed to create device", zap.Error(err))
		InternalServerError(c, "创建设备失败: "+err.Error())
		return
	}

	Success(c, device)
}

// GetDevice 获取设备
// @Summary 获取设备
// @Description 根据ID获取设备信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "设备ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "设备不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /devices/{id} [get]
// @Security BearerAuth
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid device id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的设备ID")
		return
	}

	logger.Info("handling get device request", zap.Uint("id", uint(id)))

	device, err := h.deviceService.GetDeviceByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("failed to get device", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "获取设备失败: "+err.Error())
		return
	}

	if device == nil {
		logger.Warn("device not found", zap.Uint("id", uint(id)))
		NotFound(c, "设备不存在")
		return
	}

	Success(c, device)
}

// UpdateDevice 更新设备
// @Summary 更新设备
// @Description 更新设备信息
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "设备ID"
// @Param request body dto.DeviceUpdateRequest true "更新信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "设备不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /devices/{id} [put]
// @Security BearerAuth
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid device id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的设备ID")
		return
	}

	var req dto.DeviceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	logger.Info("handling update device request", zap.Uint("id", uint(id)))

	if err := h.deviceService.UpdateDevice(c.Request.Context(), uint(id), &req); err != nil {
		logger.Error("failed to update device", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "更新设备失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "更新成功"})
}

// DeleteDevice 删除设备
// @Summary 删除设备
// @Description 删除设备（软删除）
// @Tags 设备管理
// @Accept json
// @Produce json
// @Param id path int true "设备ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /devices/{id} [delete]
// @Security BearerAuth
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid device id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的设备ID")
		return
	}

	logger.Info("handling delete device request", zap.Uint("id", uint(id)))

	if err := h.deviceService.DeleteDevice(c.Request.Context(), uint(id)); err != nil {
		logger.Error("failed to delete device", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "删除设备失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "删除成功"})
}

// ListDevices 查询设备列表
// @Summary 查询设备列表
// @Description 查询所有设备
// @Tags 设备管理
// @Accept json
// @Produce json
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器错误"
// @Router /devices [get]
// @Security BearerAuth
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	logger.Info("handling list devices request")

	var pageReq dto.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		logger.Error("invalid pagination parameters", zap.Error(err))
		BadRequest(c, "无效的分页参数: "+err.Error())
		return
	}

	devices, err := h.deviceService.ListDevices(c.Request.Context(), pageReq)
	if err != nil {
		logger.Error("failed to list devices", zap.Error(err))
		InternalServerError(c, "查询设备列表失败: "+err.Error())
		return
	}

	Success(c, devices)
}
