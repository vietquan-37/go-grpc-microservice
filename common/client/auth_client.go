package client

import (
	"common/discovery"
	"common/discovery/consul"
	"common/pb"
	"context"
)

const (
	serviceName = "auth"
	resolver    = "consul"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func InitAuthClient(consulAddr string) (*AuthClient, error) {
	err := consul.RegisterConsulResolver(consulAddr)
	if err != nil {
		return nil, err
	}

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
