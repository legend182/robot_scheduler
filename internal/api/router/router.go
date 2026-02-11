package router

import (
	"robot_scheduler/internal/api/handler"
	"robot_scheduler/internal/api/middleware"
	"robot_scheduler/internal/config"
	impl "robot_scheduler/internal/dao"
	"robot_scheduler/internal/database"
	"robot_scheduler/internal/service"
	"robot_scheduler/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 设置路由
func SetupRouter(router *gin.Engine, cfg *config.Config) {
	// 初始化DAO
	db := database.DB

	// 用户相关
	userDAO := impl.NewUserDAO(db)
	userService := service.NewUserService(userDAO)
	userHandler := handler.NewUserHandler(userService, cfg)

	// 点云地图相关
	pcdDAO := impl.NewPCDFileDAO(db)
	pcdService := service.NewPCDFileService(pcdDAO)
	pcdHandler := handler.NewPCDFileHandler(pcdService)

	// 语义地图相关
	semanticDAO := impl.NewSemanticMapDAO(db)
	semanticService := service.NewSemanticMapService(semanticDAO)
	semanticHandler := handler.NewSemanticMapHandler(semanticService)

	// 任务相关
	taskDAO := impl.NewTaskDAO(db)
	taskService := service.NewTaskService(taskDAO)
	taskHandler := handler.NewTaskHandler(taskService)

	// 设备相关
	deviceDAO := impl.NewDeviceDAO(db)
	deviceService := service.NewDeviceService(deviceDAO)
	deviceHandler := handler.NewDeviceHandler(deviceService)

	// 操作记录相关
	operationDAO := impl.NewUserOperationDAO(db)
	operationService := service.NewUserOperationService(operationDAO)
	operationHandler := handler.NewUserOperationHandler(operationService)

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "robot-scheduler",
			"version":   "1.0.0",
		})
	})

	// API路由组
	api := router.Group("/api/v1")
	{
		// 认证路由
		authConfig := cfg.Auth
		if authConfig == nil || authConfig.JWTSecret == "" {
			panic("JWT secret not configured")
		}

		// 登录路由（无需JWT认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}

		// 需要JWT认证的路由组
		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuth(authConfig.JWTSecret))
		{
			// 退出路由（需要认证）
			authAuthenticated := authenticated.Group("/auth")
			{
				authAuthenticated.POST("/logout", userHandler.Logout)
			}

			// 用户管理
			users := authenticated.Group("/users")
			{
				// 创建/编辑/删除需要用户管理权限
				users.POST("", middleware.RequirePermission(utils.PermissionUserManage), userHandler.CreateUser)
				users.PUT("/:id", middleware.RequirePermission(utils.PermissionUserManage), userHandler.UpdateUser)
				users.DELETE("/:id", middleware.RequirePermission(utils.PermissionUserManage), userHandler.DeleteUser)
				// 查看需要用户查看权限
				users.GET("/:id", middleware.RequirePermission(utils.PermissionUserView), userHandler.GetUser)
				users.GET("", middleware.RequirePermission(utils.PermissionUserView), userHandler.ListUsers)
			}

			// 地图管理模块
			maps := authenticated.Group("/maps")
			{
				// 点云地图管理
				pcds := maps.Group("/pcd-files")
				{
					// 上传凭证与创建/编辑/删除需要地图管理权限
					pcds.POST("/upload-token", middleware.RequirePermission(utils.PermissionMapManage), pcdHandler.GetPCDUploadToken)
					pcds.POST("", middleware.RequirePermission(utils.PermissionMapManage), pcdHandler.CreatePCDFile)
					pcds.PUT("/:id", middleware.RequirePermission(utils.PermissionMapManage), pcdHandler.UpdatePCDFile)
					pcds.DELETE("/:id", middleware.RequirePermission(utils.PermissionMapManage), pcdHandler.DeletePCDFile)
					// 查看需要地图查看权限
					pcds.GET("/:id", middleware.RequirePermission(utils.PermissionMapView), pcdHandler.GetPCDFile)
					pcds.GET("", middleware.RequirePermission(utils.PermissionMapView), pcdHandler.ListPCDFiles)
				}

				// 语义地图管理
				semantics := maps.Group("/semantic-maps")
				{
					// 创建/编辑/删除需要地图管理权限
					semantics.POST("", middleware.RequirePermission(utils.PermissionMapManage), semanticHandler.CreateSemanticMap)
					semantics.PUT("/:id", middleware.RequirePermission(utils.PermissionMapManage), semanticHandler.UpdateSemanticMap)
					semantics.DELETE("/:id", middleware.RequirePermission(utils.PermissionMapManage), semanticHandler.DeleteSemanticMap)
					// 查看需要地图查看权限
					semantics.GET("/:id", middleware.RequirePermission(utils.PermissionMapView), semanticHandler.GetSemanticMap)
					semantics.GET("", middleware.RequirePermission(utils.PermissionMapView), semanticHandler.ListSemanticMaps)
				}
			}

			// 任务管理
			tasks := authenticated.Group("/tasks")
			{
				// 创建/编辑/删除需要任务管理权限
				tasks.POST("", middleware.RequirePermission(utils.PermissionTaskManage), taskHandler.CreateTask)
				tasks.PUT("/:id", middleware.RequirePermission(utils.PermissionTaskManage), taskHandler.UpdateTask)
				tasks.DELETE("/:id", middleware.RequirePermission(utils.PermissionTaskManage), taskHandler.DeleteTask)
				// 查看需要任务查看权限
				tasks.GET("/:id", middleware.RequirePermission(utils.PermissionTaskView), taskHandler.GetTask)
				tasks.GET("", middleware.RequirePermission(utils.PermissionTaskView), taskHandler.ListTasks)
			}

			// 设备管理
			devices := authenticated.Group("/devices")
			{
				// 创建/编辑/删除需要设备管理权限
				devices.POST("", middleware.RequirePermission(utils.PermissionDeviceManage), deviceHandler.CreateDevice)
				devices.PUT("/:id", middleware.RequirePermission(utils.PermissionDeviceManage), deviceHandler.UpdateDevice)
				devices.DELETE("/:id", middleware.RequirePermission(utils.PermissionDeviceManage), deviceHandler.DeleteDevice)
				// 查看需要设备查看权限（普通用户也可以）
				devices.GET("/:id", middleware.RequirePermission(utils.PermissionDeviceView), deviceHandler.GetDevice)
				devices.GET("", middleware.RequirePermission(utils.PermissionDeviceView), deviceHandler.ListDevices)
			}

			// 操作记录查询（操作员及以上）
			operations := authenticated.Group("/operations")
			operations.Use(middleware.RequirePermission(utils.PermissionOperationView))
			{
				operations.GET("/:id", operationHandler.GetOperation)
				operations.GET("", operationHandler.ListOperations)
			}
		}
	}
}
