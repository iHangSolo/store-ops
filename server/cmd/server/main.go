package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"store-ops-server/internal/config"
	"store-ops-server/internal/handler"
	"store-ops-server/internal/middleware"
	"store-ops-server/internal/pkg/auth"
	"store-ops-server/internal/pkg/database"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/service"
	"store-ops-server/internal/ws"

	"github.com/gin-gonic/gin"
)

var hub *ws.Hub

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化 JWT
	auth.Init(&cfg.JWT)

	// 连接数据库
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer database.Close()

	db := database.GetDB()

	// 初始化仓库
	adminRepo := repository.NewAdminRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	commandLogRepo := repository.NewCommandLogRepository(db)
	softwareTaskRepo := repository.NewSoftwareTaskRepository(db)
	desktopSessionRepo := repository.NewDesktopSessionRepository(db)
	desktopQueueRepo := repository.NewDesktopQueueRepository(db)
	resourceHistoryRepo := repository.NewResourceHistoryRepository(db)

	// 初始化 WebSocket Hub
	hub = ws.NewHub()
	go hub.Run()

	// 初始化服务
	authService := service.NewAuthService(adminRepo, auditLogRepo)
	storeService := service.NewStoreService(storeRepo, auditLogRepo, hub)
	desktopService := service.NewDesktopService(desktopSessionRepo, desktopQueueRepo, storeRepo, auditLogRepo, hub)
	commandService := service.NewCommandService(commandLogRepo, storeRepo, auditLogRepo, hub)
	softwareService := service.NewSoftwareService(softwareTaskRepo, storeRepo, auditLogRepo, hub)
	processService := service.NewProcessService(storeRepo, auditLogRepo, hub)
	resourceService := service.NewResourceService(resourceHistoryRepo, storeRepo, hub)
	auditService := service.NewAuditService(auditLogRepo)

	// 初始化 WebSocket Handler
	wsHandler := ws.NewHandler(hub, storeRepo, auditLogRepo, resourceHistoryRepo, commandLogRepo, softwareTaskRepo, desktopSessionRepo, desktopQueueRepo)

	// 初始化 HTTP Handler
	authHandler := handler.NewAuthHandler(authService)
	storeHandler := handler.NewStoreHandler(storeService)
	desktopHandler := handler.NewDesktopHandler(desktopService)
	commandHandler := handler.NewCommandHandler(commandService)
	softwareHandler := handler.NewSoftwareHandler(softwareService)
	processHandler := handler.NewProcessHandler(processService)
	resourceHandler := handler.NewResourceHandler(resourceService)
	auditHandler := handler.NewAuditHandler(auditService)

	// 创建 Gin 引擎
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 路由
	setupRoutes(r, authHandler, storeHandler, desktopHandler, commandHandler, softwareHandler, processHandler, resourceHandler, auditHandler, wsHandler)

	// 启动服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("服务器启动，监听端口 %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭错误: %v", err)
	}

	log.Println("服务器已关闭")
}

func setupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	storeHandler *handler.StoreHandler,
	desktopHandler *handler.DesktopHandler,
	commandHandler *handler.CommandHandler,
	softwareHandler *handler.SoftwareHandler,
	processHandler *handler.ProcessHandler,
	resourceHandler *handler.ResourceHandler,
	auditHandler *handler.AuditHandler,
	wsHandler *ws.Handler,
) {
	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
		}

		// 受保护的路由
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// 门店管理
			stores := protected.Group("/stores")
			{
				stores.GET("", storeHandler.ListStores)
				stores.GET("/:id", storeHandler.GetStore)
				stores.POST("/:id/approve", storeHandler.ApproveStore)
				stores.POST("/:id/reject", storeHandler.RejectStore)
				stores.DELETE("/:id", storeHandler.DeleteStore)
				stores.POST("/:id/regenerate-token", storeHandler.RegenerateToken)

				// 远程桌面
				stores.GET("/:id/desktop/status", desktopHandler.GetDesktopStatus)
				stores.POST("/:id/desktop/connect", desktopHandler.ConnectDesktop)
				stores.POST("/:id/desktop/queue", desktopHandler.QueueDesktop)
				stores.DELETE("/:id/desktop/session/:session_id", desktopHandler.EndSession)
				stores.DELETE("/:id/desktop/queue/:queue_id", desktopHandler.LeaveQueue)

				// 命令执行
				stores.POST("/:id/command", commandHandler.ExecuteCommand)
				stores.GET("/:id/command/:commandId", commandHandler.GetCommandResult)

				// 软件推送
				stores.POST("/:id/software", softwareHandler.PushSoftware)
				stores.GET("/:id/software/:taskId", softwareHandler.GetSoftwareStatus)

				// 进程监控
				stores.GET("/:id/processes", processHandler.GetProcesses)
				stores.DELETE("/:id/processes", processHandler.KillProcess)
				stores.GET("/:id/services", processHandler.GetServices)
				stores.POST("/:id/services/:name/start", processHandler.StartService)
				stores.POST("/:id/services/:name/stop", processHandler.StopService)

				// 资源监控
				stores.GET("/:id/resource", resourceHandler.GetResource)
				stores.GET("/:id/resource/history", resourceHandler.GetResourceHistory)
			}

			// 审计日志
			audit := protected.Group("/audit-logs")
			{
				audit.GET("", auditHandler.ListAuditLogs)
				audit.GET("/:id", auditHandler.GetAuditLog)
			}
		}
	}

	// WebSocket
	r.GET("/ws", wsHandler.HandleWebSocket)
}