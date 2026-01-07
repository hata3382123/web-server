package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Set(ctx context.Context, u domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

// A用到了B , B一定是接口 => 这是保证面向接口
// A用到了B ,B一定是A的字段 =>规避包变量,包方法，都非常缺乏扩展性
// A用到了B A绝对不初始化B 而是外面注入(保持依赖注入和依赖反转)
func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 如果没有数据 要返回一个特定的err
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, ErrKeyNotExist
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}
func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
