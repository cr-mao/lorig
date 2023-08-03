package locate

import (
	"context"
	"github.com/cr-mao/lorig/cluster"
)

type Locator interface {
	// Get 获取用户定位
	Get(ctx context.Context, uid int64, insKind cluster.Kind) (string, error)
	// Set 设置用户定位
	Set(ctx context.Context, uid int64, insKind cluster.Kind, insID string) error
	// Rem 移除用户定位
	Rem(ctx context.Context, uid int64, insKind cluster.Kind, insID string) error
	// Watch 监听用户定位变化
	Watch(ctx context.Context, insKinds ...cluster.Kind) (Watcher, error)
}

type Watcher interface {
	// Next 返回用户位置列表
	Next() ([]*Event, error)
	// Stop 停止监听
	Stop() error
}

type Event struct {
	// 用户ID
	UID int64 `json:"uid"`
	// 事件类型
	Type EventType `json:"type"`
	// 实例ID
	InsID string `json:"ins_id"`
	// 实例类型
	InsKind cluster.Kind `json:"ins_kind"`
}

type EventType int

const (
	SetLocation EventType = iota // 设置定位
	RemLocation                  // 移除定位
)
