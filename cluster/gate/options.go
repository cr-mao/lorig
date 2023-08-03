package gate

import (
	"context"
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/locate"
	"github.com/cr-mao/lorig/transport"
	"github.com/cr-mao/lorig/utils/xuuid"
	"time"

	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/registry"
)

const (
	defaultName    = "gate"          // 默认名称
	defaultTimeout = 3 * time.Second // 默认超时时间
)

const (
	defaultIDKey      = "cluster.gate.id"
	defaultNameKey    = "cluster.gate.name"
	defaultTimeoutKey = "cluster.gate.timeout"
)

type Option func(o *options)

type options struct {
	id          string                // 实例ID
	name        string                // 实例名称
	ctx         context.Context       // 上下文
	timeout     time.Duration         // RPC调用超时时间
	server      network.Server        // 网关服务器
	locator     locate.Locator        // 用户定位器
	registry    registry.Registry     // 服务注册器
	transporter transport.Transporter // 消息传输器
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		timeout: defaultTimeout,
	}

	if id := conf.GetString(defaultIDKey); id != "" {
		opts.id = id
	} else if id, err := xuuid.UUID(); err == nil {
		opts.id = id
	}

	if name := conf.GetString(defaultNameKey); name != "" {
		opts.name = name
	}

	if timeout := conf.GetInt64(defaultTimeoutKey); timeout > 0 {
		opts.timeout = time.Duration(timeout) * time.Second
	}

	return opts
}

// WithID 设置实例ID
func WithID(id string) Option {
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

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置用户定位器
func WithLocator(locator locate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// WithTransporter 设置消息传输器
func WithTransporter(transporter transport.Transporter) Option {
	return func(o *options) { o.transporter = transporter }
}
