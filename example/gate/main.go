/**
User: cr-mao
Date: 2023/7/31
Time: 15:09
Desc: main.go
*/
package main

import (
	"math/rand"
	"time"

	"github.com/cr-mao/lorig"
	"github.com/cr-mao/lorig/cluster/gate"
	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/example/gate/grpc_middleware"
	"github.com/cr-mao/lorig/locate/redis"
	"github.com/cr-mao/lorig/network/tcp"
	"github.com/cr-mao/lorig/registry/etcd"
	"github.com/cr-mao/lorig/transport/grpc"
	grpclib "google.golang.org/grpc"
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
	location := redis.NewLocator()
	serverOpts := make([]grpclib.ServerOption, 0, 1)
	serverOpts = append(serverOpts, grpclib.ChainUnaryInterceptor(grpc_middleware.UnaryCrashInterceptor))
	rpcServer := grpc.NewTransporter(grpc.WithServerOptions(serverOpts...))

	gateServer := gate.NewGate(
		gate.WithServer(tcp.NewServer()),
		gate.WithLocator(location),
		gate.WithTransporter(rpcServer),
		gate.WithRegistry(etcd.NewRegistry()),
	)
	// 添加网关组件, pprof分析
	contanier.Add(gateServer, component.NewPProf())
	// 启动容器
	contanier.Serve()
}
