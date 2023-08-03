package node

import (
	"context"
	"github.com/cr-mao/lorig/conf"
	"time"

	"github.com/cr-mao/lorig/crypto"
	"github.com/cr-mao/lorig/encoding"
	"github.com/cr-mao/lorig/locate"
	"github.com/cr-mao/lorig/registry"
	"github.com/cr-mao/lorig/transport"
	"github.com/cr-mao/lorig/utils/xuuid"
)

const (
	defaultName    = "node"          // 默认节点名称
	defaultCodec   = "proto"         // 默认编解码器名称
	defaultTimeout = 3 * time.Second // 默认超时时间
)

const (
	defaultIDKey        = "config.cluster.node.id"
	defaultNameKey      = "config.cluster.node.name"
	defaultCodecKey     = "config.cluster.node.codec"
	defaultTimeoutKey   = "config.cluster.node.timeout"
	defaultEncryptorKey = "config.cluster.node.encryptor"
	defaultDecryptorKey = "config.cluster.node.decryptor"
)

type Option func(o *options)

type options struct {
	id          string                // 实例ID
	name        string                // 实例名称
	ctx         context.Context       // 上下文
	codec       encoding.Codec        // 编解码器
	timeout     time.Duration         // RPC调用超时时间
	locator     locate.Locator        // 用户定位器
	registry    registry.Registry     // 服务注册器
	transporter transport.Transporter // 消息传输器
	encryptor   crypto.Encryptor      // 消息加密器
	decryptor   crypto.Decryptor      // 消息解密器
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		codec:   encoding.Invoke(defaultCodec),
		timeout: defaultTimeout,
	}

	if id := conf.GetString(defaultIDKey, ""); id != "" {
		opts.id = id
	} else if id, err := xuuid.UUID(); err == nil {
		opts.id = id
	}

	if name := conf.GetString(defaultNameKey, ""); name != "" {
		opts.name = name
	}

	if codec := conf.GetString(defaultCodecKey); codec != "" {
		opts.codec = encoding.Invoke(codec)
	}

	if timeout := conf.GetInt64(defaultTimeoutKey, 0); timeout > 0 {
		opts.timeout = time.Duration(timeout) * time.Second
	}

	if encryptor := conf.GetString(defaultEncryptorKey, ""); encryptor != "" {
		opts.encryptor = crypto.InvokeEncryptor(encryptor)
	}

	if decryptor := conf.Get(defaultDecryptorKey, ""); decryptor != "" {
		opts.decryptor = crypto.InvokeDecryptor(decryptor)
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

// WithCodec 设置编解码器
func WithCodec(codec encoding.Codec) Option {
	return func(o *options) { o.codec = codec }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置定位器
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

// WithEncryptor 设置消息加密器
func WithEncryptor(encryptor crypto.Encryptor) Option {
	return func(o *options) { o.encryptor = encryptor }
}

// WithDecryptor 设置消息解密器
func WithDecryptor(decryptor crypto.Decryptor) Option {
	return func(o *options) { o.decryptor = decryptor }
}
