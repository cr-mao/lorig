package etcd

import (
	"context"
	"time"

	"github.com/cr-mao/lorig/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	defaultAddr          = "127.0.0.1:2379"
	defaultDialTimeout   = 5
	defaultNamespace     = "services"
	defaultTimeout       = 3
	defaultRetryTimes    = 3
	defaultRetryInterval = 10
)

const (
	defaultAddrsKey         = "registry.etcd.addrs"
	defaultDialTimeoutKey   = "registry.etcd.dialTimeout"
	defaultNamespaceKey     = "registry.etcd.namespace"
	defaultTimeoutKey       = "registry.etcd.timeout"
	defaultRetryTimesKey    = "registry.etcd.retryTimes"
	defaultRetryIntervalKey = "registry.etcd.retryInterval"
)

type Option func(o *options)

type options struct {
	// 客户端连接地址
	// 内建客户端配置，默认为[]string{"localhost:2379"}
	addrs []string

	// 客户端拨号超时时间
	// 内建客户端配置，默认为5秒
	dialTimeout time.Duration

	// 外部客户端
	// 外部客户端配置，存在外部客户端时，优先使用外部客户端，默认为nil
	client *clientv3.Client

	// 上下文
	// 默认context.Background
	ctx context.Context

	// 命名空间
	// 默认为services
	namespace string

	// 上下文超时时间
	// 默认为3秒
	timeout time.Duration

	// 心跳重试次数
	// 默认为3次
	retryTimes int

	// 心跳重试间隔
	// 默认为10秒
	retryInterval time.Duration
}

func defaultOptions() *options {
	return &options{
		ctx:           context.Background(),
		addrs:         conf.GetStrings(defaultAddrsKey, []string{defaultAddr}),
		dialTimeout:   time.Duration(conf.GetInt64(defaultDialTimeoutKey, defaultDialTimeout)) * time.Second,
		namespace:     conf.GetString(defaultNamespaceKey, defaultNamespace),
		timeout:       time.Duration(conf.GetInt64(defaultTimeoutKey, defaultTimeout)) * time.Second,
		retryTimes:    conf.GetInt(defaultRetryTimesKey, defaultRetryTimes),
		retryInterval: time.Duration(conf.GetInt64(defaultRetryIntervalKey, defaultRetryInterval)) * time.Second,
	}
}

// WithAddrs 设置客户端连接地址
func WithAddrs(addrs ...string) Option {
	return func(o *options) { o.addrs = addrs }
}

// WithDialTimeout 设置客户端拨号超时时间
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *options) { o.dialTimeout = dialTimeout }
}

// WithClient 设置外部客户端
func WithClient(client *clientv3.Client) Option {
	return func(o *options) { o.client = client }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithNamespace 设置命名空间
func WithNamespace(namespace string) Option {
	return func(o *options) { o.namespace = namespace }
}

// WithTimeout 设置上下文超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithRetryTimes 设置心跳重试次数
func WithRetryTimes(retryTimes int) Option {
	return func(o *options) { o.retryTimes = retryTimes }
}

// WithRetryInterval 设置心跳重试间隔时间
func WithRetryInterval(retryInterval time.Duration) Option {
	return func(o *options) { o.retryInterval = retryInterval }
}
