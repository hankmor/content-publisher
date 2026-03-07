package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, controller *Controller) {
	api := router.Group("/api")
	{
		api.POST("/upload/material", controller.UploadMaterialImage)
		api.POST("/upload/news-image", controller.UploadNewsImage)
		api.POST("/draft/create", controller.CreateDraft)
		// api.POST("/draft/publish", controller.PublishDraft)
	}
}
