/**
 * @Desc: 网关服务器
 */

package gate

import (
	"context"
	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/session"
)

type Gate struct {
	component.Base
	opts      *options
	ctx       context.Context
	cancel    context.CancelFunc
	session   *session.Session
	nodeProxy *nodeProxy
}

func NewGate(opts ...Option) *Gate {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	g := &Gate{}
	g.opts = o
	g.session = session.NewSession()
	g.ctx, g.cancel = context.WithCancel(o.ctx)
	g.nodeProxy = newNodeProxy(g)
	return g
}

// Name 组件名称
func (g *Gate) Name() string {
	return g.opts.name
}

// Init 初始化
func (g *Gate) Init() {
	if g.opts.id == 0 {
		log.Fatal("instance id can not be empty")
	}

	if g.opts.server == nil {
		log.Fatal("server component is not injected")
	}
}

//Start 启动组件
func (g *Gate) Start() {
	g.startNetworkServer()

	//g.registerServiceInstance()
	//
	//g.proxy.watch(g.ctx)

	g.debugPrint()
}

// Destroy 销毁组件
func (g *Gate) Destroy() {
	//g.deregisterServiceInstance()

	g.stopNetworkServer()

	//g.stopRPCServer()

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

	// 触发连接消息.....
	//cid, uid := conn.ID(), conn.UID()
	//ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//g.proxy.trigger(ctx, cluster.Connect, cid, uid)
	//cancel()
}

// 处理断开连接
func (g *Gate) handleDisconnect(conn network.Conn) {
	g.session.RemConn(conn)

	// 断链推送 给 业务服务器....

	//if cid, uid := conn.ID(), conn.UID(); uid != 0 {
	//	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//	_ = g.proxy.unbindGate(ctx, cid, uid)
	//	g.proxy.trigger(ctx, cluster.Disconnect, cid, uid)
	//	cancel()
	//} else {
	//	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//	g.proxy.trigger(ctx, cluster.Disconnect, cid, uid)
	//	cancel()
	//}
}

// 处理接收到的消息
func (g *Gate) handleReceive(conn network.Conn, data []byte) {
	// 接收到消息
	connId, userId := conn.ID(), conn.UID()
	g.nodeProxy.PushMsg(g.opts.id, connId, userId, data)
}

func (g *Gate) debugPrint() {
	log.Debugf("gate server %s-%s startup successful", g.opts.id, g.opts.name)
	log.Debugf("%s server listen on %s", g.opts.server.Protocol(), g.opts.server.Addr())
	//log.Debugf("%s server listen on %s", g.rpc.Scheme(), g.rpc.Addr())
}
