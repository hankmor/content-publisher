package service

import (
	"fmt"

	"github.com/silenceper/wechat/v2/officialaccount/context"
	"github.com/silenceper/wechat/v2/officialaccount/draft"
	"github.com/silenceper/wechat/v2/officialaccount/freepublish"
	"github.com/silenceper/wechat/v2/util"
)

// ImageInfo 图片消息的图片信息
type ImageInfo struct {
	ImageList []ImageItem `json:"image_list"`
}

// ImageItem 图片列表项
type ImageItem struct {
	ImageMediaID string `json:"image_media_id"`
}

// CoverInfo 图片消息的封面信息
type CoverInfo struct {
	CropPercentList []CropInfo `json:"crop_percent_list"`
}

// CropInfo 封面裁剪信息
type CropInfo struct {
	Ratio string `json:"ratio"`
	X1    string `json:"x1"`
	Y1    string `json:"y1"`
	X2    string `json:"x2"`
	Y2    string `json:"y2"`
}

// ArticleWithType 支持类型的文章结构
type ArticleWithType struct {
	ArticleType string     `json:"article_type"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Digest      string     `json:"digest"`
	Content     string     `json:"content"`
	ContentSourceURL string `json:"content_source_url"`
	ThumbMediaID string    `json:"thumb_media_id"`
	NeedOpenComment uint   `json:"need_open_comment"`
	OnlyFansCanComment uint `json:"only_fans_can_comment"`
	ImageInfo   *ImageInfo `json:"image_info,omitempty"`
	CoverInfo   *CoverInfo `json:"cover_info,omitempty"`
}

type DraftService struct {
	draft       *draft.Draft
	freepublish *freepublish.FreePublish
	ctx         *context.Context
}

func NewDraftService(wechatService *WechatService) *DraftService {
	official := wechatService.GetOfficialAccount()
	return &DraftService{
		draft:       official.GetDraft(),
		freepublish: official.GetFreePublish(),
		ctx:         official.GetContext(),
	}
}

func (s *DraftService) AddDraft(article *draft.Article) (string, error) {
	result, err := s.draft.AddDraft([]*draft.Article{article})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *DraftService) AddDraftWithType(articleType string, article *ArticleWithType) (string, error) {
	accessToken, err := s.ctx.GetAccessToken()
	if err != nil {
		return "", err
	}

	var req struct {
		Articles []*ArticleWithType `json:"articles"`
	}
	req.Articles = []*ArticleWithType{article}

	uri := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/draft/add?access_token=%s", accessToken)
	// 打印调试信息
	fmt.Printf("发送给微信 API 的数据: %+v\n", req)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return "", err
	}
	// 打印微信 API 返回的响应
	fmt.Printf("微信 API 返回的响应: %s\n", string(response))

	var res struct {
		util.CommonError
		MediaID string `json:"media_id"`
	}
	err = util.DecodeWithError(response, &res, "AddDraftWithType")
	return res.MediaID, err
}

func (s *DraftService) PublishDraft(mediaID string) (int64, error) {
	publishID, err := s.freepublish.Publish(mediaID)
	if err != nil {
		return 0, err
	}
	return publishID, nil
}
