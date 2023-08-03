package gate

import (
	"context"
	"time"

	"github.com/cr-mao/lorig/cluster"
	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/registry"
	"github.com/cr-mao/lorig/session"
	"github.com/cr-mao/lorig/transport"
)

type Gate struct {
	component.Base
	opts     *options
	ctx      context.Context
	cancel   context.CancelFunc
	proxy    *proxy
	instance *registry.ServiceInstance
	rpc      transport.Server
	session  *session.Session
}

func NewGate(opts ...Option) *Gate {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	g := &Gate{}
	g.opts = o
	g.proxy = newProxy(g)
	g.session = session.NewSession()
	g.ctx, g.cancel = context.WithCancel(o.ctx)

	return g
}

// Name 组件名称
func (g *Gate) Name() string {
	return g.opts.name
}

// Init 初始化
func (g *Gate) Init() {
	if g.opts.id == "" {
		log.Fatal("instance id can not be empty")
	}

	if g.opts.server == nil {
		log.Fatal("server component is not injected")
	}

	if g.opts.locator == nil {
		log.Fatal("locator component is not injected")
	}

	if g.opts.registry == nil {
		log.Fatal("registry component is not injected")
	}

	if g.opts.transporter == nil {
		log.Fatal("transporter component is not injected")
	}
}

// Start 启动组件
func (g *Gate) Start() {
	g.startNetworkServer()

	g.startRPCServer()

	g.registerServiceInstance()

	g.proxy.watch(g.ctx)

	g.debugPrint()
}

// Destroy 销毁组件
func (g *Gate) Destroy() {
	log.Infof("gate %s 停止服务", g.opts.server.Addr())
	g.deregisterServiceInstance()

	g.stopNetworkServer()

	g.stopRPCServer()

	g.cancel()
}

// 启动网络服务器
func (g *Gate) startNetworkServer() {
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnDisconnect(g.handleDisconnect)
	g.opts.server.OnReceive(g.handleReceive)

	if err := g.opts.server.Start(); err != nil {
		log.Fatalf("network server start failed: %v", err)
	}
}

// 停止网关服务器
func (g *Gate) stopNetworkServer() {
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("network server stop failed: %v", err)
	}
}

// 处理连接打开
func (g *Gate) handleConnect(conn network.Conn) {
	g.session.AddConn(conn)
	// 刚连上不处理
	//cid, uid := conn.ID(), conn.UID()
	//ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//g.proxy.trigger(ctx, cluster.Connect, cid, uid)
	//cancel()
}

// 处理断开连接
func (g *Gate) handleDisconnect(conn network.Conn) {
	g.session.RemConn(conn)

	if cid, uid := conn.ID(), conn.UID(); uid != 0 {
		ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
		_ = g.proxy.unbindGate(ctx, cid, uid)
		g.proxy.trigger(ctx, cluster.Disconnect, cid, uid)
		cancel()
	}
	// 刚连就断的，不处理.
	//else {
	//	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//	g.proxy.trigger(ctx, cluster.Disconnect, cid, uid)
	//	cancel()
	//}
}

// 处理接收到的消息
func (g *Gate) handleReceive(conn network.Conn, data []byte) {
	cid, uid := conn.ID(), conn.UID()
	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	g.proxy.deliver(ctx, cid, uid, data)
	cancel()
}

// 启动RPC服务器
func (g *Gate) startRPCServer() {
	var err error

	g.rpc, err = g.opts.transporter.NewGateServer(&provider{g})
	if err != nil {
		log.Fatalf("rpc server create failed: %v", err)
	}

	go func() {
		if err = g.rpc.Start(); err != nil {
			log.Fatalf("rpc server start failed: %v", err)
		}
	}()
}

// 停止RPC服务器
func (g *Gate) stopRPCServer() {
	if err := g.rpc.Stop(); err != nil {
		log.Errorf("rpc server stop failed: %v", err)
	}
}

// 注册服务实例
func (g *Gate) registerServiceInstance() {
	g.instance = &registry.ServiceInstance{
		ID:       g.opts.id,
		Name:     string(cluster.Gate),
		Kind:     cluster.Gate,
		Alias:    g.opts.name,
		State:    cluster.Work,
		Endpoint: g.rpc.Endpoint().String(),
	}

	ctx, cancel := context.WithTimeout(g.ctx, 10*time.Second)
	err := g.opts.registry.Register(ctx, g.instance)
	cancel()
	if err != nil {
		log.Fatalf("register dispatcher instance failed: %v", err)
	}
}

// 解注册服务实例
func (g *Gate) deregisterServiceInstance() {
	ctx, cancel := context.WithTimeout(g.ctx, 10*time.Second)
	defer cancel()
	err := g.opts.registry.Deregister(ctx, g.instance)
	if err != nil {
		log.Errorf("deregister dispatcher instance failed: %v", err)
	}
	log.Info("服务注册销毁成功")
}

func (g *Gate) debugPrint() {
	log.Debugf("gate server startup successful")
	log.Debugf("%s server listen on %s", g.opts.server.Protocol(), g.opts.server.Addr())
	log.Debugf("%s server listen on %s", g.rpc.Scheme(), g.rpc.Addr())
}
