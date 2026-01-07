package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

const codeTplId = "187756"

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// biz 区别业务场景
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	//生成验证码
	code := svc.generateCode()
	//塞入redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//发送
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}
func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)

}
func (svc *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}
