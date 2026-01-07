package tencent

import (
	"context"
	"fmt"
	"webook/pkg/ratelimit"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

const key = "sms:tencent:"

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
}

func NewService(c *sms.Client, appId string, signature string, limiter ratelimit.Limiter) *Service {
	return &Service{
		client:   c,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signature),
		limiter:  limiter,
	}
}

// biz直接代表的就是tplId
func (s *Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	limit, err := s.limiter.Limited(ctx, key)
	if err != nil {
		//系统错误
		//可以限流 保守策略 下游很坑的时候
		//可以不限 你的下游很强 业务可用性很高 尽量容错策略
		//包一下这个错误
		return fmt.Errorf("短信服务判断是否限流出现问题 %w", err)
	}
	if limit {
		return fmt.Errorf("触发限流")
	}
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "OK" {
			return fmt.Errorf("发送失败，code:%s,原因:%s", *status.Code, *status.Message)
		}
	}
	return nil
}
func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
