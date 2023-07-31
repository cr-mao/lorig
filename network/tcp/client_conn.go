package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/network"
	"github.com/cr-mao/lorig/utils/xtime"
)

type clientConn struct {
	rw                sync.RWMutex
	id                int64         // 连接ID
	uid               int64         // 用户ID
	conn              net.Conn      // TCP源连接
	state             int32         // 连接状态
	client            *client       // 客户端
	chWrite           chan chWrite  // 写入队列
	lastHeartbeatTime int64         // 上次心跳时间
	done              chan struct{} // 写入完成信号
	close             chan struct{} // 关闭信号
	localAddr         string        // 当前链接的本地地址
	remoteAddr        string        // 当前链接的远程地址
}

var _ network.Conn = &clientConn{}

func newClientConn(client *client, id int64, conn net.Conn) network.Conn {
	c := &clientConn{
		id:                id,
		conn:              conn,
		state:             int32(network.ConnOpened),
		client:            client,
		chWrite:           make(chan chWrite, 1024),
		lastHeartbeatTime: xtime.Now().Unix(),
		done:              make(chan struct{}),
		close:             make(chan struct{}),
		localAddr:         conn.LocalAddr().String(),
		remoteAddr:        conn.RemoteAddr().String(),
	}

	go c.read()

	go c.write()

	if c.client.connectHandler != nil {
		c.client.connectHandler(c)
	}

	return c
}

// ID 获取连接ID
func (c *clientConn) ID() int64 {
	return c.id
}

// UID 获取用户ID
func (c *clientConn) UID() int64 {
	return atomic.LoadInt64(&c.uid)
}

// Bind 绑定用户ID
func (c *clientConn) Bind(uid int64) {
	atomic.StoreInt64(&c.uid, uid)
}

// Unbind 解绑用户ID
func (c *clientConn) Unbind() {
	atomic.StoreInt64(&c.uid, 0)
}

// Send 发送消息（同步）
func (c *clientConn) Send(msg []byte) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err = c.checkState(); err != nil {
		return
	}

	return write(c.conn, msg)
}

// Push 发送消息（异步）
func (c *clientConn) Push(msg []byte) (err error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	if err = c.checkState(); err != nil {
		return
	}

	c.chWrite <- chWrite{typ: dataPacket, msg: msg}

	return
}

// State 获取连接状态
func (c *clientConn) State() network.ConnState {
	return network.ConnState(atomic.LoadInt32(&c.state))
}

// Close 关闭连接
func (c *clientConn) Close(isForce ...bool) error {
	if len(isForce) > 0 && isForce[0] {
		return c.forceClose()
	} else {
		return c.graceClose()
	}
}

// LocalIP 获取本地IP
func (c *clientConn) LocalIP() string {
	addr := c.LocalAddr()
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "0.0.0.0"
	}
	return ip
}

// LocalAddr 获取本地地址
func (c *clientConn) LocalAddr() string {
	return c.localAddr
}

// RemoteIP 获取远端IP
func (c *clientConn) RemoteIP() string {
	addr := c.RemoteAddr()
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "0.0.0.0"
	}
	return ip
}

// RemoteAddr 获取远端地址
func (c *clientConn) RemoteAddr() string {
	return c.remoteAddr
}

// 检测连接状态
func (c *clientConn) checkState() error {
	switch network.ConnState(atomic.LoadInt32(&c.state)) {
	case network.ConnHanged:
		return network.ErrConnectionHanged
	case network.ConnClosed:
		return network.ErrConnectionClosed
	}

	return nil
}

// 优雅关闭
func (c *clientConn) graceClose() (err error) {
	c.rw.Lock()

	if err = c.checkState(); err != nil {
		c.rw.Unlock()
		return
	}

	atomic.StoreInt32(&c.state, int32(network.ConnHanged))
	c.chWrite <- chWrite{typ: closeSig}
	c.rw.Unlock()

	<-c.done

	c.rw.Lock()
	atomic.StoreInt32(&c.state, int32(network.ConnClosed))
	close(c.chWrite)
	close(c.close)
	close(c.done)
	c.conn.Close()
	c.conn = nil
	c.rw.Unlock()

	if c.client.disconnectHandler != nil {
		c.client.disconnectHandler(c)
	}

	return
}

// 强制关闭
func (c *clientConn) forceClose() (err error) {
	c.rw.Lock()

	if err = c.checkState(); err != nil {
		c.rw.Unlock()
		return
	}

	atomic.StoreInt32(&c.state, int32(network.ConnClosed))
	close(c.chWrite)
	close(c.close)
	close(c.done)
	c.conn.Close()
	c.conn = nil
	c.rw.Unlock()

	if c.client.disconnectHandler != nil {
		c.client.disconnectHandler(c)
	}

	return
}

// 读取消息
func (c *clientConn) read() {
	for {
		select {
		case <-c.close:
			return
		default:
			msg, err := read(c.conn)
			if err != nil {
				c.forceClose()
				return
			}

			if c.client.opts.heartbeatInterval > 0 {
				atomic.StoreInt64(&c.lastHeartbeatTime, xtime.Now().Unix())
			}

			switch c.State() {
			case network.ConnHanged:
				continue
			case network.ConnClosed:
				return
			}

			// ignore heartbeat packet
			if len(msg) == 0 {
				continue
			}

			if c.client.receiveHandler != nil {
				c.client.receiveHandler(c, msg)
			}
		}
	}
}

// 写入消息
func (c *clientConn) write() {
	var ticker *time.Ticker
	if c.client.opts.heartbeatInterval > 0 {
		ticker = time.NewTicker(c.client.opts.heartbeatInterval)
		defer ticker.Stop()
	} else {
		ticker = &time.Ticker{C: make(chan time.Time, 1)}
	}

	for {
		select {
		case r, ok := <-c.chWrite:
			if !ok {
				return
			}

			c.rw.RLock()
			if r.typ == closeSig {
				c.done <- struct{}{}
				c.rw.RUnlock()
				return
			}

			if atomic.LoadInt32(&c.state) == int32(network.ConnClosed) {
				c.rw.RUnlock()
				return
			}

			err := write(c.conn, r.msg)
			c.rw.RUnlock()

			if err != nil {
				log.Errorf("write message error: %v", err)
			}
		case <-ticker.C:
			deadline := xtime.Now().Add(-3 * c.client.opts.heartbeatInterval).Unix()
			if atomic.LoadInt64(&c.lastHeartbeatTime) < deadline {
				log.Debugf("connection heartbeat timeout")
				c.forceClose()
				return
			} else {
				c.rw.RLock()

				if atomic.LoadInt32(&c.state) == int32(network.ConnClosed) {
					c.rw.RUnlock()
					return
				}
				// send heartbeat packet
				err := write(c.conn, nil)
				c.rw.RUnlock()

				if err != nil {
					log.Errorf("send heartbeat packet failed: %v", err)
				}
			}
		}
	}
}
