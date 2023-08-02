package session

import (
	"github.com/cr-mao/lorig/log"
	"sync"

	"github.com/cr-mao/lorig/errors"
	"github.com/cr-mao/lorig/network"
)

var (
	ErrNotFoundSession    = errors.New("not found session")
	ErrInvalidSessionKind = errors.New("invalid session kind")
)

const (
	Conn Kind = 1 // 连接SESSION
	User Kind = 2 // 用户SESSION
)

type Kind int

type Session struct {
	rw    sync.RWMutex           // 读写锁
	conns map[int64]network.Conn // 连接会话（连接ID -> network.Conn）
	users map[int64]network.Conn // 用户会话（用户ID -> network.Conn）
}

func NewSession() *Session {
	return &Session{
		conns: make(map[int64]network.Conn),
		users: make(map[int64]network.Conn),
	}
}

// AddConn 添加连接
func (s *Session) AddConn(conn network.Conn) {
	s.rw.Lock()
	defer s.rw.Unlock()

	cid, uid := conn.ID(), conn.UID()

	s.conns[cid] = conn

	if uid != 0 {
		s.users[uid] = conn
	}
}

// RemConn 移除连接
func (s *Session) RemConn(conn network.Conn) {
	s.rw.Lock()
	defer s.rw.Unlock()

	cid, uid := conn.ID(), conn.UID()

	delete(s.conns, cid)

	if uid != 0 {
		delete(s.users, uid)
	}
}

// Bind 绑定用户ID
func (s *Session) Bind(cid, uid int64) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	conn, err := s.conn(Conn, cid)
	if err != nil {
		return err
	}

	if oldUID := conn.UID(); oldUID != 0 {
		if uid == oldUID {
			return nil
		}
		delete(s.users, oldUID)
	}

	if oldConn, ok := s.users[uid]; ok {
		oldConn.Unbind()
	}

	conn.Bind(uid)
	s.users[uid] = conn

	return nil
}

// Unbind 解绑用户ID
func (s *Session) Unbind(uid int64) (int64, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	conn, err := s.conn(User, uid)
	if err != nil {
		return 0, err
	}

	conn.Unbind()
	delete(s.users, uid)

	return conn.ID(), nil
}

// LocalIP 获取本地IP
func (s *Session) LocalIP(kind Kind, target int64) string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, err := s.conn(kind, target)
	if err != nil {
		log.Errorf("session LocalIP error %+v", err)
		return ""
	}

	return conn.LocalIP()
}

// LocalAddr 获取本地地址
func (s *Session) LocalAddr(kind Kind, target int64) string {
	s.rw.RLock()
	defer s.rw.RUnlock()
	conn, err := s.conn(kind, target)
	if err != nil {
		log.Errorf("session LocalAddr error %+v", err)
		return ""
	}
	return conn.LocalAddr()
}

// RemoteIP 获取远端IP
func (s *Session) RemoteIP(kind Kind, target int64) string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, err := s.conn(kind, target)
	if err != nil {
		log.Errorf("session RemoteIP error %+v", err)
		return ""
	}

	return conn.RemoteIP()
}

// RemoteAddr 获取远端地址
func (s *Session) RemoteAddr(kind Kind, target int64) string {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, err := s.conn(kind, target)
	if err != nil {
		log.Errorf("session RemoteAddr error %+v", err)
		return ""
	}

	return conn.RemoteAddr()
}

// Close 关闭会话
func (s *Session) Close(kind Kind, target int64, isForce ...bool) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, err := s.conn(kind, target)
	if err != nil {
		return err
	}

	cid, uid := conn.ID(), conn.UID()

	err = conn.Close(isForce...)
	if err != nil {
		return err
	}

	delete(s.conns, cid)
	if uid != 0 {
		delete(s.users, uid)
	}

	return nil
}

// Send 发送消息（同步）
func (s *Session) Send(kind Kind, target int64, msg []byte) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, err := s.conn(kind, target)
	if err != nil {
		return err
	}

	return conn.Send(msg)
}

// Push 推送消息（异步）
func (s *Session) Push(kind Kind, target int64, msg []byte) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	conn, err := s.conn(kind, target)
	if err != nil {
		return err
	}

	return conn.Push(msg)
}

// Multicast 推送组播消息（异步）
func (s *Session) Multicast(kind Kind, targets []int64, msg []byte) (n int64, err error) {
	if len(targets) == 0 {
		return
	}

	s.rw.RLock()
	defer s.rw.RUnlock()

	var conns map[int64]network.Conn
	switch kind {
	case Conn:
		conns = s.conns
	case User:
		conns = s.users
	default:
		err = ErrInvalidSessionKind
		return
	}

	for _, target := range targets {
		conn, ok := conns[target]
		if !ok {
			continue
		}
		if conn.Push(msg) == nil {
			n++
		}
	}

	return
}

// Broadcast 推送广播消息（异步）
func (s *Session) Broadcast(kind Kind, msg []byte) (n int64, err error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	var conns map[int64]network.Conn
	switch kind {
	case Conn:
		conns = s.conns
	case User:
		conns = s.users
	default:
		err = ErrInvalidSessionKind
		return
	}

	for _, conn := range conns {
		if conn.Push(msg) == nil {
			n++
		}
	}

	return
}

// Stat 统计会话总数
func (s *Session) Stat(kind Kind) (int64, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	switch kind {
	case Conn:
		return int64(len(s.conns)), nil
	case User:
		return int64(len(s.users)), nil
	default:
		return 0, ErrInvalidSessionKind
	}
}

// 获取会话
func (s *Session) conn(kind Kind, target int64) (network.Conn, error) {

	switch kind {
	case Conn:
		conn, ok := s.conns[target]
		if !ok {
			return nil, ErrNotFoundSession
		}
		return conn, nil
	case User:
		conn, ok := s.users[target]
		if !ok {
			return nil, ErrNotFoundSession
		}
		return conn, nil
	default:
		return nil, ErrInvalidSessionKind
	}
}

// 获取会话
func (s *Session) Conn(kind Kind, target int64) (network.Conn, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	switch kind {
	case Conn:
		conn, ok := s.conns[target]
		if !ok {
			return nil, ErrNotFoundSession
		}
		return conn, nil
	case User:
		conn, ok := s.users[target]
		if !ok {
			return nil, ErrNotFoundSession
		}
		return conn, nil
	default:
		return nil, ErrInvalidSessionKind
	}
}
