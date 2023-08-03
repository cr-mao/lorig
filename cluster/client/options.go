package client

import (
	"context"
	"time"

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/crypto"
	_ "github.com/cr-mao/lorig/crypto/ecc"
	_ "github.com/cr-mao/lorig/crypto/rsa"
	"github.com/cr-mao/lorig/encoding"
	_ "github.com/cr-mao/lorig/encoding/json"
	_ "github.com/cr-mao/lorig/encoding/proto"
	_ "github.com/cr-mao/lorig/encoding/xml"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/utils/xuuid"
)

const (
	defaultName    = "client"        // 默认客户端名称
	defaultCodec   = "proto"         // 默认编解码器名称
	defaultTimeout = 3 * time.Second // 默认超时时间
)

const (
	defaultIDKey        = "cluster.client.id"
	defaultNameKey      = "cluster.client.name"
	defaultCodecKey     = "cluster.client.codec"
	defaultTimeoutKey   = "cluster.client.timeout"
	defaultEncryptorKey = "cluster.client.encryptor"
	defaultDecryptorKey = "cluster.client.decryptor"
)

type Option func(o *options)

type options struct {
	id        string           // 实例ID
	name      string           // 实例名称
	ctx       context.Context  // 上下文
	codec     encoding.Codec   // 编解码器
	client    network.Client   // 网络客户端
	timeout   time.Duration    // RPC调用超时时间
	encryptor crypto.Encryptor // 消息加密器
	decryptor crypto.Decryptor // 消息解密器
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		codec:   encoding.Invoke(defaultCodec),
		timeout: defaultTimeout,
	}

	if id := conf.Get(defaultIDKey, ""); id != "" {
		opts.id = id
	} else if id, err := xuuid.UUID(); err == nil {
		opts.id = id
	}

	if name := conf.Get(defaultNameKey, ""); name != "" {
		opts.name = name
	}

	if codec := conf.Get(defaultCodecKey, ""); codec != "" {
		opts.codec = encoding.Invoke(codec)
	}

	if timeout := conf.GetInt64(defaultTimeoutKey, 0); timeout > 0 {
		opts.timeout = time.Duration(timeout) * time.Second
	}

	if encryptor := conf.Get(defaultEncryptorKey, ""); encryptor != "" {
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

// WithClient 设置客户端
func WithClient(client network.Client) Option {
	return func(o *options) { o.client = client }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithEncryptor 设置消息加密器
func WithEncryptor(encryptor crypto.Encryptor) Option {
	return func(o *options) { o.encryptor = encryptor }
}

// WithDecryptor 设置消息解密器
func WithDecryptor(decryptor crypto.Decryptor) Option {
	return func(o *options) { o.decryptor = decryptor }
}
