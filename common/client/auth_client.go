package client

import (
	"common/discovery"
	"common/pb"
	"context"
	"github.com/rs/zerolog/log"
)

const (
	serviceName = "auth"
)

type AuthClient struct {
	registry discovery.Registry
}

func InitAuthClient(registry discovery.Registry) *AuthClient {
	return &AuthClient{
		registry: registry,
	}
}

func (a *AuthClient) Validate(ctx context.Context, token string) (*pb.ValidateRsp, error) {
	conn, err := discovery.ServiceConnection(ctx, serviceName, a.registry)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to dial to service: ")
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)
	return client.Validate(ctx, &pb.ValidateReq{
		Token: token,
	})
}
