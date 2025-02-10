package consul

import (
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	resolverName = "consul"
)

type ResolverBuilder struct {
	client *api.Client
}

func NewResolverBuilder(client *api.Client) *ResolverBuilder {
	return &ResolverBuilder{client: client}
}

func (b *ResolverBuilder) Build(
	target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions,
) (resolver.Resolver, error) {

	return NewConsulResolver(target.Endpoint(), cc, b.client), nil
}

func (b *ResolverBuilder) Scheme() string {
	return resolverName
}

func RegisterConsulResolver(client *api.Client) error {
	builder := NewResolverBuilder(client)
	resolver.Register(builder)
	return nil
}
