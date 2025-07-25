package client

import (
	"common/discovery"
	"context"

	"github.com/vietquan-37/order-service/pkg/pb"
)

const (
	resolver = "consul"
)

type ProductClient struct {
	client pb.ProductServiceClient
}

func InitProductClient(serviceName string) (*ProductClient, error) {

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
func (c *ProductClient) FindProducts(ctx context.Context, ids []int32) (*pb.ProductResponseList, error) {
	return c.client.GetProducts(ctx, &pb.GetProductsRequest{
		Id: ids,
	})
}
