package client

import (
	"context"
	"fmt"
	"github.com/vietquan-37/order-service/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	AuthClient pb.AuthServiceClient
}

func InitAuthClient(url string) AuthClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		fmt.Printf("Cannot connect to user client: %v", err)
	}
	a := AuthClient{
		AuthClient: pb.NewAuthServiceClient(conn),
	}
	return a
}
func (a *AuthClient) GetOneUser(id int32) (*pb.UserResponse, error) {
	req := &pb.GetOneUserRequest{
		Id: id,
	}
	return a.AuthClient.GetOneUser(context.Background(), req)
}
