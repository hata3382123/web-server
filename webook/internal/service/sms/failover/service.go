package failover

import (
	"context"
	"errors"
	"log"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}
func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		//正常输出日志
		//做好监控
		log.Println(err)
	}
	return errors.New("服务商全部发送失败")
}
