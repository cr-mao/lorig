/**
User: cr-mao
Date: 2023/7/30
Time: 21:13
Desc: 业务服务器
*/
package node

import (
	"context"

	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/session"
)

type Node struct {
	component.Base
	opts    *options
	ctx     context.Context
	cancel  context.CancelFunc
	session *session.Session
}

func NewNode(opts ...Option) *Node {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	g := &Node{}
	g.opts = o
	g.session = session.NewSession()
	g.ctx, g.cancel = context.WithCancel(o.ctx)

	return g
}

// Name 组件名称
func (g *Node) Name() string {
	return g.opts.name
}

// Init 初始化
func (g *Node) Init() {
	if g.opts.id == "" {
		log.Fatal("instance id can not be empty")
	}

	if g.opts.server == nil {
		log.Fatal("server component is not injected")
	}
}

//Start 启动组件
func (g *Node) Start() {
	g.startNetworkServer()

	//g.registerServiceInstance()
	//
	//g.proxy.watch(g.ctx)

	g.debugPrint()
}

// Destroy 销毁组件
func (g *Node) Destroy() {
	//g.deregisterServiceInstance()

	g.stopNetworkServer()

	//g.stopRPCServer()

	g.cancel()
}

// 启动网络服务器
func (g *Node) startNetworkServer() {
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnDisconnect(g.handleDisconnect)
	g.opts.server.OnReceive(g.handleReceive)

	if err := g.opts.server.Start(); err != nil {
		log.Fatalf("network server start failed: %v", err)
	}
}

// 停止网关服务器
func (g *Node) stopNetworkServer() {
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("network server stop failed: %v", err)
	}
}

// 处理连接打开
func (g *Node) handleConnect(conn network.Conn) {
	g.session.AddConn(conn)

	// 触发连接消息..... todo
	//cid, uid := conn.ID(), conn.UID()
	//ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//g.proxy.trigger(ctx, cluster.Connect, cid, uid)
	//cancel()
}

// 处理断开连接
func (g *Node) handleDisconnect(conn network.Conn) {
	g.session.RemConn(conn)

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
func (g *Node) handleReceive(conn network.Conn, data []byte) {
	//cid, uid := conn.ID(), conn.UID()

	// 投递消息给 node 节点...

	//ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	//g.proxy.deliver(ctx, cid, uid, data)
	//cancel()
}

func (g *Node) debugPrint() {
	log.Debugf("node server %s-%s startup successful", g.opts.id, g.opts.name)
	log.Debugf("%s server listen on %s", g.opts.server.Protocol(), g.opts.server.Addr())
}
