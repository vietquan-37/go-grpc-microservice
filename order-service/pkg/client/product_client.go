package client

import (
	"common/discovery"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/order-service/pkg/pb"
)

const (
	serviceName = "product"
)

type ProductClient struct {
	registry discovery.Registry
}

func InitProductClient(registry discovery.Registry) *ProductClient {

	return &ProductClient{
		registry: registry,
	}
}
func (c *ProductClient) FindOneProduct(ctx context.Context, productId int32) (*pb.ProductResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to dial to service: ")
	}
	defer conn.Close()
	client := pb.NewProductServiceClient(conn)

	return client.FindOneProduct(ctx, &pb.ProductRequest{
		Id: productId,
	})
}
