/**
User: cr-mao
Date: 2023/7/30
Time: 23:41
Desc: 业务服务器代理
*/
package gate

import (
	"fmt"
	"github.com/cr-mao/lorig/cluster/msg"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
	"sync"
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
	tcpClient := tcp.NewClient()
	tcpClient.OnConnect(func(conn network.Conn) {
		log.Infof("gateId:%d, connection node is opened,connId:%d,node remoteAddr:%s", np.gate.opts.id, conn.ID(), conn.RemoteAddr())
	})
	tcpClient.OnDisconnect(func(conn network.Conn) {
		log.Infof("gateId:%d, connection node is Disconnect,connId:%d,node remoteAddr:%s", np.gate.opts.id, conn.ID(), conn.RemoteAddr())
	})
	tcpClient.OnReceive(func(conn network.Conn, data []byte) {

		fmt.Println("网关收到业务服务器发来的的消息", string(data))
		// 从业务服读消息,这里还有 组播，广播逻辑 ....
		//innerMsg := &msg.InternalServerMsg{}
		//innerMsg.FromByteArray(data)
		//userId := innerMsg.
		//np.gate.session.Push()
		//fmt.Println(innerMsg.UserId)
	})

	return tcpClient.Dial()
}

// 网关投递消息到业务服务器
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
		MsgData: data, // message 结构体封包的数据
	}
	err = nodeConn.Push(innerMsg.ToByteArray())
	if err != nil {
		log.Errorf("nodeProxy push msg error %+v", err)
	}
}
