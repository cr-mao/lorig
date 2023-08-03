package gate

import (
	"context"
	
	"github.com/cr-mao/lorig/cluster"
	"github.com/cr-mao/lorig/packet"
	"github.com/cr-mao/lorig/session"
)

type provider struct {
	gate *Gate
}

// Bind 绑定用户与网关间的关系
func (p *provider) Bind(ctx context.Context, cid, uid int64) error {
	if cid <= 0 || uid <= 0 {
		return ErrInvalidArgument
	}

	err := p.gate.session.Bind(cid, uid)
	if err != nil {
		return err
	}

	err = p.gate.proxy.bindGate(ctx, cid, uid)
	if err != nil {
		_, _ = p.gate.session.Unbind(uid)
	}

	return err
}

// Unbind 解绑用户与网关间的关系
func (p *provider) Unbind(ctx context.Context, uid int64) error {
	if uid == 0 {
		return ErrInvalidArgument
	}

	cid, err := p.gate.session.Unbind(uid)
	if err != nil {
		return err
	}

	err = p.gate.proxy.unbindGate(ctx, cid, uid)
	if err != nil {
		return err
	}

	return nil
}

// GetIP 获取客户端IP地址
func (p *provider) GetIP(ctx context.Context, kind session.Kind, target int64) (string, error) {
	return p.gate.session.RemoteIP(kind, target)
}

// Push 发送消息
func (p *provider) Push(ctx context.Context, kind session.Kind, target int64, message *packet.Message) error {
	msg, err := packet.Pack(message)
	if err != nil {
		return err
	}

	err = p.gate.session.Push(kind, target, msg)
	if kind == session.User && err == session.ErrNotFoundSession {
		err = p.gate.opts.locator.Rem(ctx, target, cluster.Gate, p.gate.opts.id)
		if err != nil {
			return err
		}
	}

	return err
}

// Multicast 推送组播消息
func (p *provider) Multicast(ctx context.Context, kind session.Kind, targets []int64, message *packet.Message) (int64, error) {
	if len(targets) == 0 {
		return 0, nil
	}

	msg, err := packet.Pack(message)
	if err != nil {
		return 0, err
	}

	return p.gate.session.Multicast(kind, targets, msg)
}

// Broadcast 推送广播消息
func (p *provider) Broadcast(ctx context.Context, kind session.Kind, message *packet.Message) (int64, error) {
	msg, err := packet.Pack(message)
	if err != nil {
		return 0, err
	}

	return p.gate.session.Broadcast(kind, msg)
}

// Disconnect 断开连接
func (p *provider) Disconnect(ctx context.Context, kind session.Kind, target int64, isForce bool) error {
	return p.gate.session.Close(kind, target, isForce)
}

// Stat 统计会话总数
func (p *provider) Stat(ctx context.Context, kind session.Kind) (int64, error) {
	return p.gate.session.Stat(kind)
}
