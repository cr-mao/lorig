package grpc

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/registry"
	"github.com/cr-mao/lorig/transport/grpc/internal/client"
	"github.com/cr-mao/lorig/transport/grpc/internal/server"

	"google.golang.org/grpc"
)

const (
	defaultServerAddr     = ":0" // 默认服务器地址
	defaultClientPoolSize = 10   // 默认客户端连接池大小
)

const (
	defaultServerAddrKey       = "transport.grpc.server.addr"
	defaultServerKeyFileKey    = "transport.grpc.server.keyFile"
	defaultServerCertFileKey   = "transport.grpc.server.certFile"
	defaultClientPoolSizeKey   = "transport.grpc.client.poolSize"
	defaultClientCertFileKey   = "transport.grpc.client.certFile"
	defaultClientServerNameKey = "transport.grpc.client.serverName"
)

type Option func(o *options)

type options struct {
	server server.Options
	client client.Options
}

func defaultOptions() *options {
	opts := &options{}
	opts.server.Addr = conf.Get(defaultServerAddrKey, defaultServerAddr)
	opts.server.KeyFile = conf.Get(defaultServerKeyFileKey)
	opts.server.CertFile = conf.Get(defaultServerCertFileKey)
	opts.client.PoolSize = conf.GetInt(defaultClientPoolSizeKey, defaultClientPoolSize)
	opts.client.CertFile = conf.GetString(defaultClientCertFileKey)
	opts.client.ServerName = conf.GetString(defaultClientServerNameKey)

	return opts
}

// WithServerListenAddr 设置服务器监听地址
func WithServerListenAddr(addr string) Option {
	return func(o *options) { o.server.Addr = addr }
}

// WithServerCredentials 设置服务器证书和秘钥
func WithServerCredentials(certFile, keyFile string) Option {
	return func(o *options) { o.server.KeyFile, o.server.CertFile = keyFile, certFile }
}

// WithServerOptions 设置服务器选项
func WithServerOptions(opts ...grpc.ServerOption) Option {
	return func(o *options) { o.server.ServerOpts = opts }
}

// WithClientPoolSize 设置客户端连接池大小
func WithClientPoolSize(size int) Option {
	return func(o *options) { o.client.PoolSize = size }
}

// WithClientCredentials 设置客户端证书和校验域名
func WithClientCredentials(certFile string, serverName string) Option {
	return func(o *options) { o.client.CertFile, o.client.ServerName = certFile, serverName }
}

// WithClientDiscovery 设置客户端服务发现组件
func WithClientDiscovery(discovery registry.Discovery) Option {
	return func(o *options) { o.client.Discovery = discovery }
}

// WithClientDialOptions 设置客户端拨号选项
func WithClientDialOptions(opts ...grpc.DialOption) Option {
	return func(o *options) { o.client.DialOpts = opts }
}
