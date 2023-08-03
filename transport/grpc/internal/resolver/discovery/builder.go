package discovery

import (
	"google.golang.org/grpc/resolver"

	"github.com/cr-mao/lorig/registry"
)

const scheme = "discovery"

type Builder struct {
	dis registry.Discovery
}

var _ resolver.Builder = &Builder{}

func NewBuilder(dis registry.Discovery) *Builder {
	return &Builder{dis: dis}
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return newResolver(b.dis, target.URL.Host, cc)
}

func (b *Builder) Scheme() string {
	return scheme
}
