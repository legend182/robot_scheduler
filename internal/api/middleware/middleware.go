package middleware

import (
	"strings"
	"time"

	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/model/entity"
	"robot_scheduler/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserIDKey 用户ID上下文键
	UserIDKey ContextKey = "user_id"
	// UserNameKey 用户名上下文键
	UserNameKey ContextKey = "user_name"
	// UserRoleKey 用户角色上下文键
	UserRoleKey ContextKey = "user_role"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("panic recovered",
				zap.String("path", c.Request.URL.Path),
				zap.String("error", err),
			)
		}
		c.AbortWithStatusJSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
	})
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWTAuth JWT认证中间件
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("missing authorization header", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "未授权，请先登录",
			})
			return
		}

		// 提取token（格式：Bearer <token>）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("invalid authorization header format", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "无效的授权头格式",
			})
			return
		}

		tokenString := parts[1]

		// 验证token
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			logger.Warn("invalid token", zap.Error(err), zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "无效的token或token已过期",
			})
			return
		}

		// 将用户信息存入context
		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(UserNameKey), claims.UserName)
		c.Set(string(UserRoleKey), claims.Role)

		c.Next()
	}
}

// RequirePermission 权限检查中间件
// 要求用户拥有至少一个指定的权限
func RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从context中获取用户角色
		roleValue, exists := c.Get(string(UserRoleKey))
		if !exists {
			logger.Warn("user role not found in context", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(403, gin.H{
				"code":    403,
				"message": "权限不足：无法获取用户角色",
			})
			return
		}

		role, ok := roleValue.(entity.RoleType)
		if !ok {
			logger.Warn("invalid role type in context", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(403, gin.H{
				"code":    403,
				"message": "权限不足：无效的用户角色",
			})
			return
		}

		// 检查权限
		if !utils.HasAnyPermission(role, permissions...) {
			logger.Warn("permission denied",
				zap.String("path", c.Request.URL.Path),
				zap.String("role", string(role)),
				zap.Strings("required_permissions", permissions),
			)
			c.AbortWithStatusJSON(403, gin.H{
				"code":    403,
				"message": "权限不足：您没有执行此操作的权限",
			})
			return
		}

		c.Next()
	}
}
