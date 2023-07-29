package tcp

import (
	"time"

	"github.com/cr-mao/lorig/conf"
)

const (
	defaultServerAddr              = ":3553"
	defaultServerMaxConnNum        = 5000
	defaultServerHeartbeatInterval = 10
	defaultHandlerMsgAsync         = false
)

const (
	defaultServerAddrKey              = "network.tcp.server.addr"
	defaultServerMaxConnNumKey        = "network.tcp.server.maxConnNum"
	defaultServerHeartbeatIntervalKey = "network.tcp.server.heartbeatInterval"
	defaultHandlerMsgAsyncKey         = "network.tcp.server.handlerMsgAsync"
)

type ServerOption func(o *serverOptions)

type serverOptions struct {
	addr              string        // 监听地址，默认0.0.0.0:3553
	maxConnNum        int           // 最大连接数，默认5000
	heartbeatInterval time.Duration // 心跳检测间隔时间，默认10s
	handleMsgAsync    bool
}

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		addr:              conf.GetString(defaultServerAddrKey, defaultServerAddr),
		maxConnNum:        conf.GetInt(defaultServerMaxConnNumKey, defaultServerMaxConnNum),
		heartbeatInterval: time.Duration(conf.GetInt(defaultServerHeartbeatIntervalKey, defaultServerHeartbeatInterval)) * time.Second,
		handleMsgAsync:    conf.GetBool(defaultHandlerMsgAsyncKey, defaultHandlerMsgAsync),
	}
}

// WithServerListenAddr 设置监听地址
func WithServerListenAddr(addr string) ServerOption {
	return func(o *serverOptions) { o.addr = addr }
}

// WithServerMaxConnNum 设置连接的最大连接数
func WithServerMaxConnNum(maxConnNum int) ServerOption {
	return func(o *serverOptions) { o.maxConnNum = maxConnNum }
}

// WithServerHeartbeatInterval 设置心跳检测间隔时间
func WithServerHeartbeatInterval(heartbeatInterval time.Duration) ServerOption {
	return func(o *serverOptions) { o.heartbeatInterval = heartbeatInterval }
}

func WithServerHandlerMsgAsync(handlerMsgAsync bool) ServerOption {
	return func(o *serverOptions) { o.handleMsgAsync = handlerMsgAsync }
}
