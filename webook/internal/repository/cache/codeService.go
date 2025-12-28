package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数过多")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}
func (c *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, inputCode).Int64()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		//没有问题
		return true, nil
	case -1:
		//发送太频繁
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
	}
	return false, nil
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int64()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//没有问题
		return nil
	case -1:
		//发送太频繁
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}
func (c *CodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
