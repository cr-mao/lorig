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

	for i := 0; i < 1; i++ {
		time.Sleep(time.Millisecond * 1)
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
			fmt.Printf("receive msg from server, connection id: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))
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
		conn.Push(msg)
		conn.Push(msg)
		conn.Push(msg)
		conn.Push(msg)
		conn.Push(msg)

	}

}
