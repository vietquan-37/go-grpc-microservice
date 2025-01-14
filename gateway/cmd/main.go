package main

import (
	"common/discovery"
	"common/discovery/consul"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/gateway/pkg/auth"
	"github.com/vietquan-37/gateway/pkg/config"
	"github.com/vietquan-37/gateway/pkg/order"
	"github.com/vietquan-37/gateway/pkg/order/pb"
	"github.com/vietquan-37/gateway/pkg/product"
	productpb "github.com/vietquan-37/gateway/pkg/product/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"time"

	authpb "github.com/vietquan-37/gateway/pkg/auth/pb"
)

func main() {
	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("fail to load config:")
	}
	registry, err := consul.NewRegistry(c.ConsulAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to consul")
	}
	instanceId := discovery.GenerateInstanceID(c.ServiceName)
	if err := registry.Register(instanceId, c.ServiceName, c.GatewayPort); err != nil {
		log.Fatal().Err(err).Msg("failed to register service")
	}
	go func() {
		for {
			if err := registry.HealthCheck(instanceId); err != nil {
				log.Fatal().Err(err).Msg("failed to health check service")
			}
			time.Sleep(1 * time.Second)
		}

	}()

	authClient := auth.InitAuthClient(registry, c.AuthServiceName)
	productClient := product.InitProductClient(registry, c.ProductServiceName)
	orderClient := order.InitOrderClient(registry, c.OrderServiceName)
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(
		jsonOption,
	)

	if err = authpb.RegisterAuthServiceHandlerClient(context.Background(), mux, authClient.Client); err != nil {
		log.Fatal().Err(err).Msg("fail to register auth client: ")
	}
	if err = productpb.RegisterProductServiceHandlerClient(context.Background(), mux, productClient.Client); err != nil {
		log.Fatal().Err(err).Msg("fail to register product client:")
	}

	if err = pb.RegisterOrderServiceHandlerClient(context.Background(), mux, orderClient.Client); err != nil {
		log.Fatal().Err(err).Msg("fail to register order client:")
	}
	lis, err := net.Listen("tcp", c.GatewayPort)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to listen:")
	}
	if err = http.Serve(lis, mux); err != nil {
		log.Fatal().Err(err).Msg("gateway server closed abruptly: ")
	}

}
