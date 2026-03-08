package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/api"
	"github.com/hankmor/wechat-publisher/internal/config"
	"github.com/hankmor/wechat-publisher/internal/log"
	"github.com/hankmor/wechat-publisher/internal/service"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 初始化日志
	if err := log.InitLogger(&cfg.Log); err != nil {
		panic(err)
	}
	defer log.Logger.Sync()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化服务
	wechatService := service.NewWechatService(cfg.Wechat)
	materialService := service.NewMaterialService(wechatService)
	draftService := service.NewDraftService(wechatService)

	// 初始化控制器
	controller := api.NewController(materialService, draftService)

	// 初始化路由
	router := gin.Default()
	// 提供静态文件服务
	router.StaticFile("/", "./API_GUIDE.md")
	// router.StaticFile("/demo", "./demo.html")
	// 传递配置到路由设置
	api.SetupRoutes(router, controller, cfg)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("服务器启动", zap.String("addr", addr))
	if err := router.Run(addr); err != nil {
		log.Fatal("启动服务器失败", zap.Error(err))
	}
}
