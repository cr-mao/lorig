/**
User: cr-mao
Desc: 业务服务器
*/
package node

import (
	"context"

	"github.com/cr-mao/lorig/cluster"
	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/packet"
)

// 用户ID， 网关连接id
type RequestHandler func(conn network.Conn, innerMsg *cluster.InternalServerMsg, message *packet.Message)

type Node struct {
	component.Base
	opts   *options
	ctx    context.Context
	cancel context.CancelFunc
	Route  map[int32]RequestHandler
	//session *session.Session
}

func NewNode(opts ...Option) *Node {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	node := &Node{}
	node.opts = o
	//node.session = session.NewSession()
	node.ctx, node.cancel = context.WithCancel(o.ctx)
	node.Route = make(map[int32]RequestHandler)
	return node
}

// 添加路由
func (n *Node) AddRouter(routerID int32, requestHandler RequestHandler) {
	n.Route[routerID] = requestHandler
}

// Name 组件名称
func (node *Node) Name() string {
	return node.opts.name
}

// Init 初始化
func (node *Node) Init() {
	if node.opts.id == 0 {
		log.Fatal("instance node id can not be empty")
	}
	if node.opts.server == nil {
		log.Fatal("node server component is not injected")
	}
}

//Start 启动组件
func (node *Node) Start() {
	node.startNetworkServer()

	//g.registerServiceInstance()
	//
	//g.proxy.watch(g.ctx)
	node.debugPrint()
}

// Destroy 销毁组件
func (node *Node) Destroy() {
	//g.deregisterServiceInstance()
	node.stopNetworkServer()
	//g.stopRPCServer()
	node.cancel()
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
	log.Infof("有连接进来 remoteAddr:%s,localAddr:%s", conn.RemoteAddr(), conn.LocalAddr())
}

// 处理断开连接
func (g *Node) handleDisconnect(conn network.Conn) {
	//todo  看看是否要报警处理.... 因为这个是内部的连接
	log.Infof("有连接进来 remoteAddr:%s,localAddr:%s", conn.RemoteAddr(), conn.LocalAddr())
}

// 处理接收到的消息
func (node *Node) handleReceive(conn network.Conn, data []byte) {
	innerMsg := &cluster.InternalServerMsg{}
	err := innerMsg.UnPack(data)
	if err != nil {
		log.Errorf("node handlerReceive error: %v", err)
		return
	}
	// 断连处理
	if innerMsg.EventType == int16(cluster.Disconnect) {
		return
	}
	realData := innerMsg.MsgData
	message, err := packet.Unpack(realData)
	if err != nil {
		log.Errorf("node handleReceive Unpack error: %v", err)
		return
	}
	requestHandle, ok := node.Route[message.Route]
	if !ok {
		log.Errorf("handleReceive routeId not exist %d", message.Route)
		return
	}

	// 处理消息
	requestHandle(conn, innerMsg, message)
}

func (node *Node) debugPrint() {
	log.Debugf("node server %d-%s startup successful", node.opts.id, node.opts.name)
	log.Debugf("%s server listen on %s", node.opts.server.Protocol(), node.opts.server.Addr())
}
