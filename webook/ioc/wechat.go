package ioc

import "webook/internal/service/oauth2/wechat"

func InitWechatService() wechat.Service {
	appID := "123"
	return wechat.NewService(appID)
}
