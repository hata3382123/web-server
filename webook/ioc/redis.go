package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:16379",
	})
	return redisClient
}
