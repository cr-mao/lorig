package redis_location

import (
	"context"

	"github.com/cr-mao/lorig/cluster"
	"github.com/cr-mao/lorig/location"
	"github.com/cr-mao/lorig/utils/xconv"
	"github.com/go-redis/redis/v8"
)

const (
	userLocationsKeyPrefix = "loc:uid:" // hash
)

var _ location.Locator = &RedisLocation{}

type RedisLocation struct {
	opts *options
}

func NewRedisLocation(opts ...Option) *RedisLocation {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.client == nil {
		o.client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:      o.addrs,
			DB:         o.db,
			Username:   o.username,
			Password:   o.password,
			MaxRetries: o.maxRetries,
		})
	}

	l := &RedisLocation{
		opts: o,
	}
	return l
}

// Get 获取用户定位
func (l *RedisLocation) Get(ctx context.Context, uid int64, insKind cluster.Kind) (string, error) {
	keyStr := xconv.Int64ToString(uid)
	key := userLocationsKeyPrefix + keyStr
	val, err := l.opts.client.HGet(ctx, key, string(insKind)).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return val, nil
}

// Set 设置用户定位
func (l *RedisLocation) Set(ctx context.Context, uid int64, insKind cluster.Kind, insID string) error {
	keyStr := xconv.Int64ToString(uid)
	key := userLocationsKeyPrefix + keyStr
	err := l.opts.client.HSet(ctx, key, string(insKind), insID).Err()
	if err != nil {
		return err
	}
	return nil
}

// Rem 移除用户定位
func (l *RedisLocation) Rem(ctx context.Context, uid int64, insKind cluster.Kind, insID string) error {
	oldInsID, err := l.Get(ctx, uid, insKind)
	if err != nil {
		return err
	}
	// 在其他网关上了。
	if oldInsID == "" || oldInsID != insID {
		return nil
	}
	keyStr := xconv.Int64ToString(uid)
	key := userLocationsKeyPrefix + keyStr
	err = l.opts.client.HDel(ctx, key, string(insKind)).Err()
	if err != nil {
		return err
	}
	return nil
}
