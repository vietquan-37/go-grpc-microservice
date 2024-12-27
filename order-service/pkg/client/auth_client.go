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
	// Use grpc.Dial to create a connection
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Cannot connect to Auth service: %v\n", err)
		return AuthClient{}
	}
	return AuthClient{
		AuthClient: pb.NewAuthServiceClient(conn),
	}
}

func (a *AuthClient) GetOneUser(id int32) (*pb.UserResponse, error) {
	req := &pb.GetOneUserRequest{
		Id: id,
	}
	return a.AuthClient.GetOneUser(context.Background(), req)
}
