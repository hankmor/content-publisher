package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/api"
	"github.com/hankmor/wechat-publisher/internal/config"
	"github.com/hankmor/wechat-publisher/internal/log"
	"github.com/hankmor/wechat-publisher/internal/service"
	"go.uber.org/zap"
)

// getStaticFileDir 获取静态文件所在目录
// 优先从当前工作目录查找，如果不存在则从可执行文件目录查找
func getStaticFileDir() string {
	// 首先尝试从当前工作目录查找
	wd, err := os.Getwd()
	if err == nil {
		// 检查当前目录是否有 API_GUIDE.md
		if _, err := os.Stat(filepath.Join(wd, "API_GUIDE.md")); err == nil {
			return wd
		}
	}

	// 如果当前目录没有，则从可执行文件目录查找
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	execDir := filepath.Dir(ex)

	// 检查可执行文件目录是否有 API_GUIDE.md
	if _, err := os.Stat(filepath.Join(execDir, "API_GUIDE.md")); err == nil {
		return execDir
	}

	// 如果是 go run 产生的临时文件，尝试找到项目根目录
	// go run 的临时文件通常在 /tmp/go-build*/exe/main 或类似路径
	if strings.Contains(execDir, os.TempDir()) || strings.Contains(execDir, "go-build") {
		// 尝试从当前工作目录返回
		if wd != "" {
			return wd
		}
	}

	return execDir
}

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

	// 获取静态文件所在目录
	staticDir := getStaticFileDir()
	log.Info("静态文件目录", zap.String("path", staticDir))

	// 初始化路由
	router := gin.Default()
	// 提供静态文件服务 - 使用绝对路径
	apiGuidePath := filepath.Join(staticDir, "API_GUIDE.md")
	// demoPath := filepath.Join(staticDir, "demo.html")

	// 检查文件是否存在
	if _, err := os.Stat(apiGuidePath); err == nil {
		router.StaticFile("/api", apiGuidePath)
		log.Info("API文档服务已启用", zap.String("path", apiGuidePath))
	} else {
		log.Warn("API文档文件不存在", zap.String("path", apiGuidePath))
	}

	// if _, err := os.Stat(demoPath); err == nil {
	// 	router.StaticFile("/demo", demoPath)
	// 	log.Info("Demo页面服务已启用", zap.String("path", demoPath))
	// } else {
	// 	log.Warn("Demo页面文件不存在", zap.String("path", demoPath))
	// }

	// 传递配置到路由设置
	api.SetupRoutes(router, controller, cfg)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("服务器启动", zap.String("addr", addr))
	if err := router.Run(addr); err != nil {
		log.Fatal("启动服务器失败", zap.Error(err))
	}
}
