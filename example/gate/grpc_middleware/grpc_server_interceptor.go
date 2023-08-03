package grpc_middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"runtime/debug"

	"github.com/cr-mao/lorig/log"
)

// 防止panic crash 中间件
func UnaryCrashInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		log.Errorf("recovery method: %s, message: %+v\n \n %s", info.FullMethod, r, debug.Stack())
	})
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	if gstatus, ok := status.FromError(err); ok {
		errLog := "grpc error:method:%s, code:%v,message:%v"
		log.Errorf(errLog, info.FullMethod, gstatus.Code(), err.Error())
	} else {
		errLog := "not grpc error:method:%s,message:%v"
		log.Errorf(errLog, info.FullMethod, err.Error())
	}
	return
}

func handleCrash(hanlder func(interface{})) {
	if r := recover(); r != nil {
		hanlder(r)
	}
}
