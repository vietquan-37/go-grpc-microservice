package auth

import (
	"fmt"
	"github.com/vietquan-37/gateway/pkg/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client pb.AuthServiceClient
}

func InitAuthClient(addr string) AuthClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		fmt.Printf("cannot connect to auth client: %v", err)
	}
	return AuthClient{Client: pb.NewAuthServiceClient(conn)}
}
