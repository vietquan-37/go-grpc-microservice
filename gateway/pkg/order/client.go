package order

import (
	"fmt"
	"github.com/vietquan-37/gateway/pkg/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	Client pb.OrderServiceClient
}

func InitOrderClient(addr string) OrderClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		fmt.Printf("cannot connect to order client: %v", err)
	}
	return OrderClient{pb.NewOrderServiceClient(conn)}
}
