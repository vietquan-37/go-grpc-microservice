package client

import (
	"common/discovery"
	"common/discovery/consul"
	"context"

	"github.com/vietquan-37/order-service/pkg/pb"
)

const (
	serviceName = "product"
	resolver    = "consul"
)

type ProductClient struct {
	client pb.ProductServiceClient
}

func InitProductClient(consulAddr string) (*ProductClient, error) {
	err := consul.RegisterConsulResolver(consulAddr)
	if err != nil {
		return nil, err
	}

	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	client := pb.NewProductServiceClient(conn)
	return &ProductClient{
		client: client,
	}, nil
}

func (c *ProductClient) FindOneProduct(ctx context.Context, productId int32) (*pb.ProductResponse, error) {

	return c.client.FindOneProduct(ctx, &pb.ProductRequest{
		Id: productId,
	})
}
