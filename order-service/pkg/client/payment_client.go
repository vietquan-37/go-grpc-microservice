package client

import (
	"common/discovery"
	"context"
	"github.com/vietquan-37/order-service/pkg/pb"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func InitPaymentClient(serviceName string) (*PaymentClient, error) {
	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	return &PaymentClient{client: pb.NewPaymentServiceClient(conn)}, nil
}
func (c *PaymentClient) CreatePaymentLink(ctx context.Context, req *pb.PaymentLinkRequest) (*pb.PaymentLinkResponse, error) {
	return c.client.CreatePaymentLink(ctx, req)
}
