package product

import (
	"fmt"
	"github.com/vietquan-37/gateway/pkg/product/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	Client pb.ProductServiceClient
}

func InitProductClient(addr string) ProductClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		fmt.Printf("cannot connect to product client: %v", err)
	}
	return ProductClient{
		Client: pb.NewProductServiceClient(conn),
	}
}
