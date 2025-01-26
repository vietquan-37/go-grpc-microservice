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

func NewConsulResolverBuilder(consulAddr string) (*ResolverBuilder, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ResolverBuilder{client: client}, nil
}

func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return NewConsulResolver(target.Endpoint(), cc, b.client), nil
}

func (b *ResolverBuilder) Scheme() string {
	return resolverName
}

func RegisterConsulResolver(consulAddr string) error {
	builder, err := NewConsulResolverBuilder(consulAddr)
	if err != nil {
		return err
	}
	resolver.Register(builder)
	return nil
}
