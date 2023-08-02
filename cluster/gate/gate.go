/**
 * @Desc: 网关服务器
 */

package gate

import (
	"context"
	"github.com/cr-mao/lorig/cluster"

	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/session"
)

type Gate struct {
	component.Base
	opts     *options
	ctx      context.Context
	cancel   context.CancelFunc
	session  *session.Session
	proxy    *proxy
	provider *provider
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
	g.proxy = newProxy(g)
	g.provider = newProvider(g)
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

	if g.opts.location == nil {
		log.Fatal("user location is not injected")
	}
}

//Start 启动组件
func (g *Gate) Start() {
	g.startNetworkServer()
	//g.registerServiceInstance()
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

	//fmt.Println("add conn to session", conn.ID())
	g.session.AddConn(conn)
	// 无需通知node连接 ,相信大部分场景 是不用让node知道的
}

// 处理断开连接
func (g *Gate) handleDisconnect(conn network.Conn) {
	//fmt.Println("disconnect")
	g.session.RemConn(conn)
	// 断链推送 给 业务服务器....
	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	defer cancel()
	if cid, uid := conn.ID(), conn.UID(); uid != 0 {
		_ = g.proxy.unbindGate(ctx, cid, uid)
		g.proxy.PushMsg(g.opts.id, cid, uid, cluster.Disconnect, nil)
	}
	// 没有用户id的要不要通知， 和 上面handleConnect是同一个问题。
}

// 处理接收到的消息
func (g *Gate) handleReceive(conn network.Conn, data []byte) {
	// 接收到消息
	connId, userId := conn.ID(), conn.UID()
	g.proxy.PushMsg(g.opts.id, connId, userId, cluster.Send, data)
}

func (g *Gate) debugPrint() {
	log.Debugf("gate server %d-%s startup successful", g.opts.id, g.opts.name)
	log.Debugf("%s server listen on %s", g.opts.server.Protocol(), g.opts.server.Addr())
	//log.Debugf("%s server listen on %s", g.rpc.Scheme(), g.rpc.Addr())
}
