package order

import (
	"common/discovery"
	"context"

	"github.com/vietquan-37/gateway/pkg/order/pb"
)

type Client struct {
	Client pb.OrderServiceClient
}

func InitOrderClient(serviceName, resolver string) (*Client, error) {

	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &Client{
		Client: client,
	}, nil
}
