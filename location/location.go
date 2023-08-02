/**
User: cr-mao
Date: 2023/8/1 16:13
Email: crmao@qq.com
Desc: location.go
*/
package location

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
}
