/**
User: cr-mao
Date: 2023/7/31
Time: 15:17
Desc: main.go
*/
package main

import (
	"github.com/cr-mao/lorig/cluster"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/packet"
	"math/rand"
	"time"

	"github.com/cr-mao/lorig"
	"github.com/cr-mao/lorig/cluster/node"
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
)

// 登录

func Login(conn network.Conn, req *cluster.InternalServerMsg, message *packet.Message) {
	message.Buffer = []byte("login ok")
	req.UserId = 1000
	req.MsgData, _ = packet.Pack(message)
	msgData, err := req.Pack()
	if err != nil {
		log.Errorf("login error: %v", err)
		return
	}
	// 找到网关对应到连接
	err = conn.Push(msgData)
	if err != nil {
		log.Errorf("login push msg error: %v", err)
		return
	}
}

func RegistRouter(node *node.Node) {
	node.AddRouter(1, Login)
}

func main() {
	conf.InitConfig("local")
	//随机数种子
	rand.Seed(time.Now().UnixNano())
	contanier := lorig.NewContainer()
	nodeServer := node.NewNode(node.WithServer(
		tcp.NewServer(),
	))
	RegistRouter(nodeServer)
	// 添加网关组件
	contanier.Add(nodeServer)
	// 启动容器
	contanier.Serve()
}
