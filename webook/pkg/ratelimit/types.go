package ratelimit

import "context"

type Limiter interface {
	//limited 有没有出发限流 key就是限流对象
	//bool 代表是否限流 true就是要限流
	//err 限流器本身有没有错误
	Limited(ctx context.Context, key string) (bool, error)
}
