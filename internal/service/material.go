package service

import (
	"fmt"

	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/util"
)

type MaterialService struct {
	material *material.Material
}

func NewMaterialService(wechatService *WechatService) *MaterialService {
	return &MaterialService{
		material: wechatService.GetOfficialAccount().GetMaterial(),
	}
}

func (s *MaterialService) UploadMaterialImage(filePath string) (string, error) {
	mediaType := material.MediaTypeImage
	mediaId, url, err := s.material.AddMaterial(mediaType, filePath)
	fmt.Printf("mediaId: %s, url: %s\n", mediaId, url)
	if err != nil {
		return "", err
	}
	return mediaId, nil
}

func (s *MaterialService) UploadNewsImage(filePath string) (string, error) {
	accessToken, err := s.material.GetAccessToken()
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s", accessToken)

	response, err := util.PostFile("media", filePath, uri)
	if err != nil {
		return "", err
	}

	var res struct {
		util.CommonError
		URL string `json:"url"`
	}

	err = util.DecodeWithError(response, &res, "UploadNewsImage")
	if err != nil {
		return "", err
	}

	return res.URL, nil
}
