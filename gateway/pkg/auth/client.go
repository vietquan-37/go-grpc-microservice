package auth

import (
	"common/discovery"
	"common/discovery/consul"
	"context"
	"github.com/vietquan-37/gateway/pkg/auth/pb"
)

type Client struct {
	Client pb.AuthServiceClient
}

func InitAuthClient(consulAddr, serviceName, resolver string) (*Client, error) {
	err := consul.RegisterConsulResolver(consulAddr)
	if err != nil {
		return nil, err
	}

	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &Client{
		Client: client,
	}, nil
}
