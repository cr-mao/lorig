/**
User: cr-mao
Date: 2023/7/31
Time: 15:09
Desc: main.go
*/
package main

import (
	"github.com/cr-mao/lorig"
	"github.com/cr-mao/lorig/cluster/gate"
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/network/tcp"
	"math/rand"
	"time"
)

func main() {
	conf.InitConfig("local")
	//随机数种子
	rand.Seed(time.Now().UnixNano())
	// 配置初始化，依赖命令行 --env 参数
	//全局设置时区
	var cstZone, _ = time.LoadLocation(conf.GetString("app.timezone"))
	time.Local = cstZone

	contanier := lorig.NewContainer()
	gateServer := gate.NewGate(gate.WithServer(
		tcp.NewServer(),
	))
	// 添加网关组件
	contanier.Add(gateServer)
	// 启动容器
	contanier.Serve()
}
