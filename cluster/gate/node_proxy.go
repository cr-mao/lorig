/**
User: cr-mao
Date: 2023/7/30
Time: 23:41
Desc: 业务服务器代理
*/
package gate

import (
	"sync"
	"time"

	"github.com/cr-mao/lorig/cluster/msg"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
	"github.com/cr-mao/lorig/packet"
)

type nodeProxy struct {
	gate *Gate // 网关服
}

func newNodeProxy(gate *Gate) *nodeProxy {
	return &nodeProxy{
		gate: gate,
	}
}

var nodeServerConn network.Conn

var locker = &sync.Mutex{}

// todo 服务发现
func (np *nodeProxy) GetNodeServerConn() (network.Conn, error) {
	if nodeServerConn != nil {
		return nodeServerConn, nil
	}
	locker.Lock()
	defer locker.Unlock()
	if nil != nodeServerConn {
		return nodeServerConn, nil
	}
	tcpClient := tcp.NewClient(tcp.WithClientDialAddr("127.0.0.1:4001"), tcp.WithClientHeartbeatInterval(time.Second*2))
	tcpClient.OnConnect(func(conn network.Conn) {
		remoteNetAddr, err := conn.RemoteAddr()
		if err != nil {
			log.Infof("gateId:%d,connection node is opened,connId:%d", np.gate.opts.id, conn.ID())
			return
		}
		log.Infof("gateId:%d, connection node is opened,connId:%d,node remoteAddr:%s", np.gate.opts.id, conn.ID(), remoteNetAddr.String())
	})

	tcpClient.OnDisconnect(func(conn network.Conn) {
		log.Errorf("gateId:%d ,connection is Disconnect,connId:%d", np.gate.opts.id, conn.ID())
	})
	tcpClient.OnReceive(func(conn network.Conn, data []byte) {
		// 从业务服读消息,这里还有 组播，广播逻辑 ....

		//innerMsg := &msg.InternalServerMsg{}
		//innerMsg.FromByteArray(data)
		//userId := innerMsg.
		//np.gate.session.Push()
		//fmt.Println(innerMsg.UserId)
	})

	return tcpClient.Dial()
}

func (np *nodeProxy) PushMsg(gateId int32, connId int64, userId int64, data []byte) {
	nodeConn, err := np.GetNodeServerConn()
	if err != nil {
		log.Error("node conn error")
		return
	}
	innerMsg := &msg.InternalServerMsg{
		GateId:  gateId,
		ConnId:  connId,
		UserId:  userId,
		MsgData: data,
	}
	// 内部使用seq 0 ,route 0
	message, err := packet.Pack(&packet.Message{
		Seq:    0,
		Route:  0,
		Buffer: innerMsg.ToByteArray(),
	})
	if err != nil {
		log.Errorf("nodeProxy Pack msg error %+v", err)
		return
	}
	err = nodeConn.Push(message)
	if err != nil {
		log.Errorf("nodeProxy push msg error %+v", err)
	}
}