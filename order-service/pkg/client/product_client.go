package client

import (
	"context"
	"github.com/vietquan-37/order-service/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ProductClient struct {
	Client pb.ProductServiceClient
}

func InitProductClient(url string) ProductClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := ProductClient{
		Client: pb.NewProductServiceClient(conn),
	}
	return c
}
func (c *ProductClient) FindOneProduct(productId int32) (*pb.ProductResponse, error) {
	req := &pb.ProductRequest{
		Id: productId,
	}
	return c.Client.FindOneProduct(context.Background(), req)
}
