package handler

import (
	"strconv"

	"robot_scheduler/internal/config"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/dto"
	"robot_scheduler/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
	config      *config.Config
}

func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		config:      cfg,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body dto.UserCreateRequest true "用户信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	logger.Info("handling create user request")

	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	authConfig := h.config.Auth
	if authConfig == nil || authConfig.DESKey == "" {
		logger.Error("auth config or DES key not found")
		InternalServerError(c, "认证配置缺失或DES密钥未配置")
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req, authConfig.DESKey)
	if err != nil {
		logger.Error("failed to create user", zap.Error(err))
		InternalServerError(c, "创建用户失败: "+err.Error())
		return
	}

	Success(c, user)
}

// GetUser 获取用户
// @Summary 获取用户
// @Description 根据ID获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "用户不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid user id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的用户ID")
		return
	}

	logger.Info("handling get user request", zap.Uint("id", uint(id)))

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		logger.Error("failed to get user", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "获取用户失败: "+err.Error())
		return
	}

	if user == nil {
		logger.Warn("user not found", zap.Uint("id", uint(id)))
		NotFound(c, "用户不存在")
		return
	}

	Success(c, user)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body dto.UserUpdateRequest true "更新信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 404 {object} Response "用户不存在"
// @Failure 500 {object} Response "服务器错误"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid user id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的用户ID")
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	authConfig := h.config.Auth
	if authConfig == nil || authConfig.DESKey == "" {
		logger.Error("auth config or DES key not found")
		InternalServerError(c, "认证配置缺失或DES密钥未配置")
		return
	}

	logger.Info("handling update user request", zap.Uint("id", uint(id)))

	if err := h.userService.UpdateUser(c.Request.Context(), uint(id), &req, authConfig.DESKey); err != nil {
		logger.Error("failed to update user", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "更新用户失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "更新成功"})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户（软删除）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器错误"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("invalid user id", zap.String("id", idStr), zap.Error(err))
		BadRequest(c, "无效的用户ID")
		return
	}

	logger.Info("handling delete user request", zap.Uint("id", uint(id)))

	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		logger.Error("failed to delete user", zap.Error(err), zap.Uint("id", uint(id)))
		InternalServerError(c, "删除用户失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "删除成功"})
}

// ListUsers 查询用户列表
// @Summary 查询用户列表
// @Description 查询所有用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器错误"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	logger.Info("handling list users request")

	var pageReq dto.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		logger.Error("invalid pagination parameters", zap.Error(err))
		BadRequest(c, "无效的分页参数: "+err.Error())
		return
	}

	users, err := h.userService.ListUsers(c.Request.Context(), pageReq)
	if err != nil {
		logger.Error("failed to list users", zap.Error(err))
		InternalServerError(c, "查询用户列表失败: "+err.Error())
		return
	}

	Success(c, users)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录，返回JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "用户名或密码错误"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	logger.Info("handling login request")

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request parameters", zap.Error(err))
		BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 获取认证配置
	authConfig := h.config.Auth
	if authConfig == nil {
		logger.Error("auth config not found")
		InternalServerError(c, "认证配置未找到")
		return
	}

	// 调用服务层登录
	response, err := h.userService.Login(
		c.Request.Context(),
		req.UserName,
		req.Password,
		authConfig.DESKey,
		authConfig.JWTSecret,
		authConfig.JWTExpireHours,
	)
	if err != nil {
		logger.Warn("login failed", zap.Error(err), zap.String("username", req.UserName))
		Unauthorized(c, err.Error())
		return
	}

	Success(c, response)
}

// Logout 用户退出
// @Summary 用户退出
// @Description 用户退出登录
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response "成功"
// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	logger.Info("handling logout request")

	if err := h.userService.Logout(c.Request.Context()); err != nil {
		logger.Error("logout failed", zap.Error(err))
		InternalServerError(c, "退出失败: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "退出成功"})
}
