package client

import (
	"common/pb"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	AuthClient pb.AuthServiceClient
}

func InitAuthClient(url string) *AuthClient {
	// Use grpc.Dial to create a connection
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Cannot connect to Auth service: %v\n", err)
		return nil
	}
	return &AuthClient{
		AuthClient: pb.NewAuthServiceClient(conn),
	}
}

func (a *AuthClient) GetOneUser(id int32) (*pb.User, error) {
	req := &pb.GetOneUseReq{
		Id: id,
	}
	return a.AuthClient.GetOneUser(context.Background(), req)
}

func (a *AuthClient) Validate(token string) (*pb.ValidateRsp, error) {
	return a.AuthClient.Validate(context.Background(), &pb.ValidateReq{
		Token: token,
	})
}
