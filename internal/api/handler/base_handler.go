package handler

import (
	"net/http"

	"robot_scheduler/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	logger.Error("api error", zap.Int("code", code), zap.String("message", message))
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalServerError 服务器内部错误
func InternalServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}

// ForbiddenPermission 权限不足
func ForbiddenPermission(c *gin.Context, message string) {
	if message == "" {
		message = "权限不足：您没有执行此操作的权限"
	}
	Error(c, 403, message)
}
