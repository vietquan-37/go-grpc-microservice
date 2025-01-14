package auth

import (
	"common/discovery"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/gateway/pkg/auth/pb"
)

type Client struct {
	Client pb.AuthServiceClient
}

func InitAuthClient(registry discovery.Registry, serviceName string) Client {
	conn, err := discovery.ServiceConnection(context.Background(), serviceName, registry)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to dial to service: ")
	}
	return Client{Client: pb.NewAuthServiceClient(conn)}
}
