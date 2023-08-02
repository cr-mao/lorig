package tcp_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
	"github.com/cr-mao/lorig/packet"
)

func TestNewClient_Dial(t *testing.T) {
	t.SkipNow()
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 1)
		go func(num int) {
			client := tcp.NewClient()
			client.OnConnect(func(conn network.Conn) {
				t.Log("connection is opened")
			})
			client.OnDisconnect(func(conn network.Conn) {
				t.Log("client connection is closed")
			})

			client.OnReceive(func(conn network.Conn, msg []byte) {
				message, err := packet.Unpack(msg)
				if err != nil {
					t.Error(err)
					return
				}

				t.Logf("receive msg from server, connection id: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))
			})

			defer wg.Done()

			conn, err := client.Dial()
			if err != nil {
				t.Fatal(err)
			}

			ticker := time.NewTicker(time.Second * 1)
			defer ticker.Stop()
			defer conn.Close()

			times := 0
			msg, _ := packet.Pack(&packet.Message{
				Seq:    1,
				Route:  1,
				Buffer: []byte("hello server~~"),
			})

			for {
				select {
				case <-ticker.C:
					if err = conn.Push(msg); err != nil {
						fmt.Println("push msg error", num)
						t.Error(err)
						return
					}
					times++
					if times >= 5 {
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
}

func Test_Benchmark(t *testing.T) {
	// 并发数
	concurrency := 6000
	// 消息量
	total := 12000
	// 总共发送的消息条数
	totalSent := int64(0)
	// 总共接收的消息条数
	totalRecv := int64(0)
	// 准备消息
	msg, err := packet.Pack(&packet.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte("hello server~~"),
	})
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	client := tcp.NewClient()
	client.OnReceive(func(conn network.Conn, msg []byte) {
		atomic.AddInt64(&totalRecv, 1)

		message, err := packet.Unpack(msg)
		if err != nil {
			fmt.Println(err)
		}
		if message.Seq != 1 {
			fmt.Println("seq error")
		}
		if message.Route != 1 {
			fmt.Println("Route error")
		}
		if string(message.Buffer) != "login ok" {
			fmt.Println("date error")
		}

		wg.Done()
	})

	wg.Add(total)

	chMsg := make(chan struct{}, total)

	// 准备连接
	conns := make([]network.Conn, concurrency)
	for i := 0; i < concurrency; i++ {
		conn, err := client.Dial()
		if err != nil {
			fmt.Println("connect failed", i, err)
			i--
			continue
		}

		conns[i] = conn
		time.Sleep(time.Millisecond * 2)
	}

	// 发送消息
	for _, conn := range conns {
		go func(conn network.Conn) {
			defer conn.Close(true)

			for {
				select {
				case _, ok := <-chMsg:
					if !ok {
						return
					}

					if err = conn.Push(msg); err != nil {
						t.Error(err)
						return
					}

					atomic.AddInt64(&totalSent, 1)
				}
			}
		}(conn)

	}

	startTime := time.Now().UnixNano()

	for i := 0; i < total; i++ {
		chMsg <- struct{}{}
	}

	wg.Wait()
	close(chMsg)

	totalTime := float64(time.Now().UnixNano()-startTime) / float64(time.Second)

	/*
		server               : tcp
		concurrency          : 1000
		latency              : 66.533924s
		sent     requests    : 1000000
		received requests    : 1000000
		throughput  (TPS)    : 15029
	*/

	fmt.Printf("server               : %s\n", "tcp")
	fmt.Printf("concurrency          : %d\n", concurrency)
	fmt.Printf("latency              : %fs\n", totalTime)
	fmt.Printf("sent     requests    : %d\n", totalSent)
	fmt.Printf("received requests    : %d\n", totalRecv)
	fmt.Printf("throughput  (TPS)    : %d\n", int64(float64(totalRecv)/totalTime))
}
