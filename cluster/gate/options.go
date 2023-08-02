package gate

import (
	"context"
	"math/rand"
	"time"

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/location"
	"github.com/cr-mao/lorig/network"
)

const (
	defaultName    = "gate_name" // 默认名称
	defaultTimeOut = 3
)

const (
	defaultIDKey      = "cluster.gate.id"
	defaultNameKey    = "cluster.gate.name"
	defaultTimeOutKey = "cluster.gate.timeout"
)

type Option func(o *options)

type options struct {
	id       int32           // 实例ID
	name     string          // 实例名称
	ctx      context.Context // 上下文
	server   network.Server  // 网关服务器
	location location.Locator
	timeout  time.Duration // 用户定位器, redis超时时间, rpc 请求超时时间。
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    conf.GetString(defaultNameKey, defaultName),
		timeout: time.Duration(conf.GetInt64(defaultTimeOutKey, defaultTimeOut)) * time.Second,
	}
	opts.id = conf.GetInt32(defaultIDKey, rand.Int31())
	return opts
}

// WithID 设置实例ID
func WithID(id int32) Option {
	return func(o *options) { o.id = id }
}

// WithName 设置实例名称
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithServer 设置服务器
func WithServer(server network.Server) Option {
	return func(o *options) { o.server = server }
}

func WithLocation(location location.Locator) Option {
	return func(o *options) { o.location = location }
}
