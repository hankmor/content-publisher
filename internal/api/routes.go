package api

import (
	"github.com/hankmor/wechat-publisher/internal/config"
	"github.com/hankmor/wechat-publisher/internal/log"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, controller *Controller, cfg *config.Config) {
	// 初始化频率限制存储
	var rateLimitStore *RateLimitStore
	if cfg.API.RateLimit.Enabled {
		var err error
		rateLimitStore, err = NewRateLimitStore("./data")
		if err != nil {
			log.Error("初始化频率限制存储失败", zap.Error(err))
			// 如果初始化失败，继续运行，但不启用频率限制
			cfg.API.RateLimit.Enabled = false
		}
	}

	// 应用访问控制中间件到 API 路由组
	api := router.Group("/api")
	api.Use(AccessControlMiddleware(cfg))
	
	// 添加频率限制中间件（如果启用）
	if cfg.API.RateLimit.Enabled && rateLimitStore != nil {
		rateLimitCfg := &RateLimitConfig{
			Enabled:        cfg.API.RateLimit.Enabled,
			MaxDraftPerDay: cfg.API.RateLimit.MaxDraftPerDay,
		}
		api.Use(RateLimitMiddleware(rateLimitStore, rateLimitCfg))
		log.Info("频率限制已启用", zap.Int("max_draft_per_day", cfg.API.RateLimit.MaxDraftPerDay))
	}
	
	{
		api.POST("/upload/material", controller.UploadMaterialImage)
		api.POST("/upload/news-image", controller.UploadNewsImage)
		api.POST("/draft/create", controller.CreateDraft)
		// 现在不允许个人账号发布了
		api.POST("/draft/publish", controller.PublishDraft)
	}
}
