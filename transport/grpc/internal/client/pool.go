package client

import (
	"sync/atomic"

	"google.golang.org/grpc"
)

type Pool struct {
	count uint64
	index uint64
	conns []*grpc.ClientConn
}

func newPool(count int, target string, opts ...grpc.DialOption) (*Pool, error) {
	p := &Pool{count: uint64(count), conns: make([]*grpc.ClientConn, count)}

	for i := 0; i < count; i++ {
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			return nil, err
		}
		p.conns[i] = conn
	}

	return p, nil
}

func (p *Pool) Get() *grpc.ClientConn {
	return p.conns[int(atomic.AddUint64(&p.index, 1)%p.count)]
}
