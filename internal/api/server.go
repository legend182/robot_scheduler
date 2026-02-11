package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"robot_scheduler/internal/api/middleware"
	"robot_scheduler/internal/api/router"
	"robot_scheduler/internal/config"
	"robot_scheduler/internal/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
}

func NewServer(cfg *config.Config) *Server {
	// 设置Gin模式
	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	engine := gin.New()

	// 注册中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 注册路由
	router.SetupRouter(engine, cfg)

	// 注册Swagger
	if cfg.App.Mode != "release" {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}


	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.App.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.App.WriteTimeout) * time.Second,
	}

	return &Server{
		config: cfg,
		router: engine,
		server: server,
	}
}

func (s *Server) Start() error {
	logger.Info(fmt.Sprintf("starting HTTP server on %s", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
