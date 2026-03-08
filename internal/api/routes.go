package api

import (
	"github.com/hankmor/wechat-publisher/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, controller *Controller, cfg *config.Config) {
	// 应用访问控制中间件到 API 路由组
	api := router.Group("/api")
	api.Use(AccessControlMiddleware(cfg))
	{
		api.POST("/upload/material", controller.UploadMaterialImage)
		api.POST("/upload/news-image", controller.UploadNewsImage)
		api.POST("/draft/create", controller.CreateDraft)
		// 现在不允许个人账号发布了
		api.POST("/draft/publish", controller.PublishDraft)
	}
}
