package auth

import (
	"common/discovery"
	"context"
	"github.com/vietquan-37/gateway/pkg/auth/pb"
)

type Client struct {
	Client pb.AuthServiceClient
}

func InitAuthClient(serviceName, resolver string) (*Client, error) {

	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &Client{
		Client: client,
	}, nil
}
