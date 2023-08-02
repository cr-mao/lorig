/**
User: cr-mao
Date: 2023/7/30
Time: 23:41
Desc: 业务服务器代理
*/
package gate

import (
	"context"
	"sync"

	"github.com/cr-mao/lorig/cluster"
	"github.com/cr-mao/lorig/location"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/network/tcp"
	"github.com/cr-mao/lorig/session"
	"github.com/cr-mao/lorig/utils/xconv"
)

type proxy struct {
	gate     *Gate // 网关服
	location location.Locator
}

func newProxy(gate *Gate) *proxy {
	return &proxy{
		gate:     gate,
		location: gate.opts.location,
	}
}

var nodeServerConn network.Conn

var locker = &sync.Mutex{}

// todo 服务发现
func (np *proxy) GetNodeServerConn() (network.Conn, error) {
	if nodeServerConn != nil {
		return nodeServerConn, nil
	}
	locker.Lock()
	defer locker.Unlock()
	if nodeServerConn != nil {
		return nodeServerConn, nil
	}

	tcpClient := tcp.NewClient()
	tcpClient.OnConnect(func(conn network.Conn) {
		log.Infof("gateId:%d, connection node is opened,connId:%d,node remoteAddr:%s", np.gate.opts.id, conn.ID(), conn.RemoteAddr())
	})

	tcpClient.OnDisconnect(func(conn network.Conn) {
		// 发往飞书 要......
		log.Infof("gateId:%d, connection node is Disconnect,connId:%d,node remoteAddr:%s", np.gate.opts.id, conn.ID(), conn.RemoteAddr())
	})

	tcpClient.OnReceive(func(conn network.Conn, data []byte) {
		innerMsg := &cluster.InternalServerMsg{}
		innerMsg.UnPack(data)
		//log.Info("tcp client 收到 connid :", innerMsg.ConnId)
		//log.Info("tcp client 收到 userid:", innerMsg.UserId)
		gateConn, err := np.gate.session.Conn(session.Conn, innerMsg.ConnId)
		if err != nil {
			log.Errorf("get conn by connid err:%+v", err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), np.gate.opts.timeout)
		defer cancel()
		// 第一次绑定用户
		if gateConn.UID() <= 0 && innerMsg.UserId > 0 {
			//log.Infof("userId:%d, 绑定session", innerMsg.UserId)
			err = np.gate.provider.Bind(ctx, innerMsg.ConnId, innerMsg.UserId)
			if err != nil {
				log.Error(err)
				return
			}
		}
		//  innerMsg.UserId 理论上 都>0
		err = np.gate.session.Push(session.Conn, innerMsg.ConnId, innerMsg.MsgData)
		if err != nil {
			log.Errorf("push msg to user err; userId:%d,connId:%d,err:%v", innerMsg.UserId, innerMsg.ConnId, err)
		}
		// 从业务服读消息,这里还有 组播，广播逻辑 ....
	})
	var err error
	nodeServerConn, err = tcpClient.Dial()
	return nodeServerConn, err
}

// 网关投递消息到业务服务器
func (np *proxy) PushMsg(gateId int32, connId int64, userId int64, eventType cluster.Event, data []byte) {
	nodeConn, err := np.GetNodeServerConn()
	if err != nil {
		log.Errorf("proxy conn error:%+v", err)
		return
	}
	innerMsg := &cluster.InternalServerMsg{
		GateId:    gateId,
		ConnId:    connId,
		UserId:    userId,
		EventType: int16(eventType),
		MsgData:   data, // message 结构体封包的数据
	}
	newData, err := innerMsg.Pack()
	if err != nil {
		log.Errorf("proxy Pack error %+v", err)
		return
	}
	err = nodeConn.Push(newData)
	if err != nil {
		log.Errorf("proxy push msg error %+v", err)
	}
}

// 取消绑定网关
func (np *proxy) unbindGate(ctx context.Context, connId int64, userId int64) error {
	err := np.location.Rem(ctx, userId, cluster.Gate, xconv.Int32ToString(np.gate.opts.id))
	if err != nil {
		log.Errorf("user unbind failed, gid: %d, cid: %d, uid: %d, err: %v", np.gate.opts.id, connId, userId, err)
	}
	return err
}

// 绑定用户与网关间的关系 todo 确定下放在哪在哪里 合适。
func (p *proxy) bindGate(ctx context.Context, cid, uid int64) error {
	// 这里有性能问题
	err := p.gate.opts.location.Set(ctx, uid, cluster.Gate, xconv.Int32ToString(p.gate.opts.id))
	if err != nil {
		return err
	}
	//p.trigger(ctx, cluster.Reconnect, cid, uid)
	return nil
}
