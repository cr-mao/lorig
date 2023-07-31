/**
User: cr-mao
Date: 2023/7/31
Time: 15:17
Desc: main.go
*/
package main

import (
	"fmt"
	"github.com/cr-mao/lorig"
	"github.com/cr-mao/lorig/cluster/node"
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
	"math/rand"
	"time"
)

// 登录
func Login(conn network.Conn, userId int64, gateWayConnId int64, data []byte) {

	fmt.Println("login", string(data))
}

// 移动
func Move(conn network.Conn, userId int64, gateWayConnId int64, data []byte) {
	fmt.Println("Move", string(data))
}

func RegistRouter(node *node.Node) {
	node.AddRouter(1, Login)
	node.AddRouter(2, Move)
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
