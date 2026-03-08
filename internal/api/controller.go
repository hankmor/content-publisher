package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hankmor/wechat-publisher/internal/log"
	"github.com/hankmor/wechat-publisher/internal/service"
	"go.uber.org/zap"
)

type Controller struct {
	materialService *service.MaterialService
	draftService    *service.DraftService
}

func NewController(materialService *service.MaterialService, draftService *service.DraftService) *Controller {
	return &Controller{
		materialService: materialService,
		draftService:    draftService,
	}
}

func (c *Controller) UploadMaterialImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Error("获取上传文件失败", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tempPath := "/tmp/" + file.Filename
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		log.Error("保存上传文件失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mediaID, err := c.materialService.UploadMaterialImage(tempPath)
	if err != nil {
		log.Error("上传素材图片失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("上传素材图片成功", zap.String("media_id", mediaID))
	ctx.JSON(http.StatusOK, gin.H{"media_id": mediaID})
}

func (c *Controller) UploadNewsImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Error("获取上传文件失败", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tempPath := "/tmp/" + file.Filename
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		log.Error("保存上传文件失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url, err := c.materialService.UploadNewsImage(tempPath)
	if err != nil {
		log.Error("上传图文消息图片失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("上传图文消息图片成功", zap.String("url", url))
	ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func (c *Controller) CreateDraft(ctx *gin.Context) {
	var request struct {
		Type string `json:"type" default:"news"`
		Article struct {
			Title              string `json:"title"`
			Author             string `json:"author"`
			Digest             string `json:"digest"`
			Content            string `json:"content"`
			ContentSourceURL   string `json:"content_source_url"`
			ThumbMediaID       string `json:"thumb_media_id"`
			NeedOpenComment    uint   `json:"need_open_comment"`
			OnlyFansCanComment uint   `json:"only_fans_can_comment"`
		} `json:"article"`
		ImageArticle struct {
			Title        string   `json:"title"`
			Author       string   `json:"author"`
			Digest       string   `json:"digest"`
			Content      string   `json:"content"`
			MediaIDList  []string `json:"media_id_list"`
		} `json:"image_article"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("解析请求参数失败", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 默认类型为图文消息
	if request.Type == "" {
		request.Type = "news"
	}

	var mediaID string
	var err error

	if request.Type == "news" {
		// 创建图文草稿
		articleWithType := &service.ArticleWithType{
			ArticleType:        "news",
			Title:              request.Article.Title,
			Author:             request.Article.Author,
			Digest:             request.Article.Digest,
			Content:            request.Article.Content,
			ContentSourceURL:   request.Article.ContentSourceURL,
			ThumbMediaID:       request.Article.ThumbMediaID,
			NeedOpenComment:    request.Article.NeedOpenComment,
			OnlyFansCanComment: request.Article.OnlyFansCanComment,
		}
		mediaID, err = c.draftService.AddDraftWithType("news", articleWithType)
	} else if request.Type == "newspic" {
		// 创建图片草稿
		imageList := make([]service.ImageItem, 0)
		for _, mediaID := range request.ImageArticle.MediaIDList {
			imageList = append(imageList, service.ImageItem{ImageMediaID: mediaID})
		}
		articleWithType := &service.ArticleWithType{
			ArticleType:        "newspic",
			Title:              request.ImageArticle.Title,
			Author:             request.ImageArticle.Author,
			Content:            request.ImageArticle.Content,
			NeedOpenComment:    1, // 默认打开评论
			OnlyFansCanComment: 0, // 所有人都可以评论
			ImageInfo: &service.ImageInfo{
				ImageList: imageList,
			},
		}
		mediaID, err = c.draftService.AddDraftWithType("newspic", articleWithType)
	} else {
		log.Error("不支持的草稿类型", zap.String("type", request.Type))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不支持的草稿类型"})
		return
	}

	if err != nil {
		log.Error("创建草稿失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("创建草稿成功", zap.String("type", request.Type), zap.String("media_id", mediaID))
	ctx.JSON(http.StatusOK, gin.H{"media_id": mediaID})
}

func (c *Controller) PublishDraft(ctx *gin.Context) {
	var request struct {
		MediaID string `json:"media_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("解析请求参数失败", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publishID, err := c.draftService.PublishDraft(request.MediaID)
	if err != nil {
		log.Error("发布草稿失败", zap.Error(err), zap.String("media_id", request.MediaID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("发布草稿成功", zap.String("media_id", request.MediaID), zap.Int64("publish_id", publishID))
	ctx.JSON(http.StatusOK, gin.H{"publish_id": publishID})
}
