package product

import (
	"common/discovery"
	"context"

	"github.com/vietquan-37/gateway/pkg/product/pb"
)

type Client struct {
	Client pb.ProductServiceClient
}

func InitProductClient(serviceName, resolver string) (*Client, error) {

	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	client := pb.NewProductServiceClient(conn)
	return &Client{
		Client: client,
	}, nil
}
