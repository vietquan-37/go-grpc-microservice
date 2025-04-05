package client

import (
	"common/discovery"
	"common/pb"
	"context"
)

const (
	resolver = "consul"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func InitAuthClient(serviceName string) (*AuthClient, error) {

	conn, err := discovery.ServiceConnection(context.Background(), serviceName, resolver)
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{
		client: client,
	}, nil
}
func (a *AuthClient) Validate(ctx context.Context, token string) (*pb.ValidateRsp, error) {
	return a.client.Validate(ctx, &pb.ValidateReq{Token: token})
}
