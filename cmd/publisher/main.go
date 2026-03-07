package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/api"
	"github.com/hankmor/wechat-publisher/internal/config"
	"github.com/hankmor/wechat-publisher/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

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
	router.StaticFile("/demo", "./demo.html")
	// 传递配置到路由设置
	api.SetupRoutes(router, controller, cfg)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
