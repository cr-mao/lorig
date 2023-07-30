package gate

import (
	"context"
	"math/rand"

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/network"
)

const (
	defaultName = "gate_name" // 默认名称
)

const (
	defaultIDKey   = "cluster.gate.id"
	defaultNameKey = "cluster.gate.name"
)

type Option func(o *options)

type options struct {
	id     int32           // 实例ID
	name   string          // 实例名称
	ctx    context.Context // 上下文
	server network.Server  // 网关服务器
}

func defaultOptions() *options {
	opts := &options{
		ctx:  context.Background(),
		name: conf.GetString(defaultNameKey, defaultName),
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
