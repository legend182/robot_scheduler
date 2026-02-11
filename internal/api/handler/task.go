package handler

import (
	"strconv"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 创建新任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param request body dto.TaskCreateRequest true "任务信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /tasks [post]
// @Security BearerAuth
func (h *TaskHandler) CreateTask(c *gin.Context) {
	logger.Info("handling create task request")

	var req dto.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), &req)
	if err != nil {
		logger.Error("failed to create task", zap.Error(err))
		InternalServerError(c, "创建任务失败: "+err.Error())
		return
	}

	Success(c, task)
}

// GetTask 获取任务
// @Summary 获取任务
// @Description 根据ID获取任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "任务不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /tasks/{id} [get]
// @Security BearerAuth
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid task id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的任务ID")
		return
	}

	logger.Info("handling get task request", zap.Uint("id", uint(id)))

	task, err := h.taskService.GetTaskByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("failed to get task", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "获取任务失败: "+err.Error())
		return
	}

	if task == nil {
		logger.Warn("task not found", zap.Uint("id", uint(id)))
		NotFound(c, "任务不存在")
		return
	}

	Success(c, task)
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Param request body dto.TaskUpdateRequest true "更新信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "任务不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /tasks/{id} [put]
// @Security BearerAuth
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid task id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的任务ID")
		return
	}

	var req dto.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	logger.Info("handling update task request", zap.Uint("id", uint(id)))

	if err := h.taskService.UpdateTask(c.Request.Context(), uint(id), &req); err != nil {
		logger.Error("failed to update task", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "更新任务失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "更新成功"})
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 删除任务（软删除）
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /tasks/{id} [delete]
// @Security BearerAuth
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid task id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的任务ID")
		return
	}

	logger.Info("handling delete task request", zap.Uint("id", uint(id)))

	if err := h.taskService.DeleteTask(c.Request.Context(), uint(id)); err != nil {
		logger.Error("failed to delete task", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "删除任务失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "删除成功"})
}

// ListTasks 查询任务列表
// @Summary 查询任务列表
// @Description 查询所有任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器错误"
// @Router /tasks [get]
// @Security BearerAuth
func (h *TaskHandler) ListTasks(c *gin.Context) {
	logger.Info("handling list tasks request")

	var pageReq dto.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		logger.Error("invalid pagination parameters", zap.Error(err))
		BadRequest(c, "无效的分页参数: "+err.Error())
		return
	}

	tasks, err := h.taskService.ListTasks(c.Request.Context(), pageReq)
	if err != nil {
		logger.Error("failed to list tasks", zap.Error(err))
		InternalServerError(c, "查询任务列表失败: "+err.Error())
		return
	}

	Success(c, tasks)
}
