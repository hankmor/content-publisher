package service

import (
	"github.com/hankmor/wechat-publisher/internal/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	officialaccountconfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

type WechatService struct {
	app      *wechat.Wechat
	official *officialaccount.OfficialAccount
}

func NewWechatService(cfg config.WechatConfig) *WechatService {
	app := wechat.NewWechat()
	officialConfig := &officialaccountconfig.Config{
		AppID:     cfg.AppID,
		AppSecret: cfg.AppSecret,
		Token:     cfg.Token,
		Cache:     cache.NewMemory(),
	}
	official := app.GetOfficialAccount(officialConfig)

	return &WechatService{
		app:      app,
		official: official,
	}
}

func (s *WechatService) GetOfficialAccount() *officialaccount.OfficialAccount {
	return s.official
}
