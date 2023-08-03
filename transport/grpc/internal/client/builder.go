package client

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/cr-mao/lorig/registry"
	"github.com/cr-mao/lorig/transport/grpc/internal/resolver/direct"
	"github.com/cr-mao/lorig/transport/grpc/internal/resolver/discovery"
)

type Builder struct {
	err      error
	opts     *Options
	dialOpts []grpc.DialOption
	pools    sync.Map
}

type Options struct {
	PoolSize   int
	CertFile   string
	ServerName string
	Discovery  registry.Discovery
	DialOpts   []grpc.DialOption
}

func NewBuilder(opts *Options) *Builder {
	b := &Builder{opts: opts}

	var creds credentials.TransportCredentials
	if opts.CertFile != "" && opts.ServerName != "" {
		creds, b.err = credentials.NewClientTLSFromFile(opts.CertFile, opts.ServerName)
		if b.err != nil {
			return b
		}
	} else {
		creds = insecure.NewCredentials()
	}

	resolvers := make([]resolver.Builder, 0, 2)
	resolvers = append(resolvers, direct.NewBuilder())
	if opts.Discovery != nil {
		resolvers = append(resolvers, discovery.NewBuilder(opts.Discovery))
	}

	b.dialOpts = make([]grpc.DialOption, 0, len(opts.DialOpts)+2)
	b.dialOpts = append(b.dialOpts, grpc.WithTransportCredentials(creds))
	b.dialOpts = append(b.dialOpts, grpc.WithResolvers(resolvers...))

	return b
}

// Build 构建连接
func (b *Builder) Build(target string) (*grpc.ClientConn, error) {
	if b.err != nil {
		return nil, b.err
	}

	val, ok := b.pools.Load(target)
	if ok {
		return val.(*Pool).Get(), nil
	}

	size := b.opts.PoolSize
	if size <= 0 {
		size = 10
	}

	pool, err := newPool(size, target, b.dialOpts...)
	if err != nil {
		return nil, err
	}

	b.pools.Store(target, pool)

	return pool.Get(), nil
}
