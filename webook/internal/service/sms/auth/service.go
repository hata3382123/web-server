package auth

import (
	"context"
	"errors"
	"webook/internal/service/sms"

	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

// send发送 biz必须是线下申请的一个代表业务方的 token
func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

	var tc Claims
	//如果这里能解析成功 说明对应的就是业务方
	//没有error就说明 token是我发的
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (any, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token不合法")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	jwt.Claims
	Tpl string
}
