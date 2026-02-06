package wechat

import (
	"context"
	"fmt"

	uuid "github.com/lithammer/shortuuid/v4"
)

const redirectURL = "daming.com"

type Service interface {
	Auth2URL(ctx context.Context) (string, error)
}
type service struct {
	appID string
}

func NewService(appID string) Service {
	return &service{
		appID: appID,
	}
}

func (s *service) Auth2URL(ctx context.Context) (string, error) {
	const urlPattern = "http://open.weixin.come%s123%s321%345"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appID, redirectURL, state), nil
}
