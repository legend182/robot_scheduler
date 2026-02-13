package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"robot_scheduler/internal/api"
	"robot_scheduler/internal/config"
	impl "robot_scheduler/internal/dao"
	"robot_scheduler/internal/database"
	"robot_scheduler/internal/logger"
	"robot_scheduler/internal/minio_client"
	"robot_scheduler/internal/service"

	_ "robot_scheduler/docs"

	"go.uber.org/zap"
)

// @title 机器人调度系统 API
// @version 1.0
// @description 机器人调度系统后台管理接口
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	cfgPath := "configs/config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	if err := config.Init(cfgPath); err != nil {
		panic(fmt.Sprintf("failed to init config: %v", err))
	}

	cfg := config.Get()

	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("failed to sync logger", zap.Error(err))
		}
	}()

	logger.Info("starting robot scheduler",
		zap.String("version", cfg.App.Version),
		zap.String("mode", cfg.App.Mode),
	)

	// 初始化数据库
	if err := database.Init(cfg); err != nil {
		logger.Fatal("failed to init database", zap.Error(err))
	}

	// 初始化超级管理员用户
	if err := initSuperAdmin(cfg); err != nil {
		logger.Error("failed to init superAdmin", zap.Error(err))
		// 不阻止程序启动，只记录错误
	}

	// 初始化MinIO客户端
	if cfg.Minio != nil && cfg.Minio.Enabled {
		if err := minio_client.Init(cfg.Minio); err != nil {
			logger.Fatal("failed to init minio", zap.Error(err))
		}
	}

	// 初始化HTTP服务器
	server := api.NewServer(cfg)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info(fmt.Sprintf("server started on :%d", cfg.App.Port))
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}

	logger.Info("server exited")
}

// initSuperAdmin 初始化超级管理员用户
func initSuperAdmin(cfg *config.Config) error {
	// 检查认证配置
	if cfg.Auth == nil {
		return fmt.Errorf("auth config is nil")
	}

	if cfg.Auth.DESKey == "" {
		return fmt.Errorf("DES key is not configured")
	}

	// 创建用户服务
	userDAO := impl.NewUserDAO(database.DB)
	userService := service.NewUserService(userDAO)

	// 初始化超级管理员
	ctx := context.Background()
	return userService.InitSuperAdmin(ctx, cfg.Auth.DESKey)
}
