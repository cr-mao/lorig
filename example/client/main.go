/**
User: cr-mao
Date: 2023/7/31
Time: 15:38
Desc: 模拟app客户端
*/
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
	"github.com/cr-mao/lorig/packet"
)

func main() {

	conf.InitConfig("local")
	//随机数种子
	rand.Seed(time.Now().UnixNano())
	// 配置初始化，依赖命令行 --env 参数
	//全局设置时区
	var cstZone, _ = time.LoadLocation(conf.GetString("app.timezone"))
	time.Local = cstZone

	timeStart := time.Now().UnixMilli()
	for i := 0; i < 1; i++ {
		client := tcp.NewClient()
		client.OnConnect(func(conn network.Conn) {
			fmt.Println("connection is opened")
		})
		client.OnDisconnect(func(conn network.Conn) {
			fmt.Println("client connection is closed")
		})
		client.OnReceive(func(conn network.Conn, msg []byte) {
			message, err := packet.Unpack(msg)
			if err != nil {
				fmt.Println("receive err", err)
				return
			}
			fmt.Printf("receive msg from server, connection id: %d, seq: %d, route: %d, msg: %s\n", conn.ID(), message.Seq, message.Route, string(message.Buffer))
		})
		conn, err := client.Dial()
		if err != nil {
			fmt.Println("dial err", err)
		}

		defer conn.Close()
		msg, _ := packet.Pack(&packet.Message{
			Seq:    1,
			Route:  1,
			Buffer: []byte("hello server~~"),
		})
		num := 0
		for {
			conn.Push(msg)
			num++
			if num > 10000 {
				break
			}
		}
		//conn.Push(msg)
		//conn.Push(msg)
		//conn.Push(msg)
	}
	fmt.Println("耗时", time.Now().UnixMilli()-timeStart)

	time.Sleep(10 * time.Second)
}
