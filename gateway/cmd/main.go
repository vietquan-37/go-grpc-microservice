package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vietquan-37/gateway/pkg/auth"
	"github.com/vietquan-37/gateway/pkg/config"
	"github.com/vietquan-37/gateway/pkg/middleware"
	"github.com/vietquan-37/gateway/pkg/order"
	"github.com/vietquan-37/gateway/pkg/order/pb"
	"github.com/vietquan-37/gateway/pkg/product"
	productpb "github.com/vietquan-37/gateway/pkg/product/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"

	authpb "github.com/vietquan-37/gateway/pkg/auth/pb"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("fail to load config: %v", err)
	}
	authClient := auth.InitAuthClient(c.AuthUrl)
	productClient := product.InitProductClient(c.ProductUrl)
	orderClient := order.InitOrderClient(c.OrderUrl)
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md, _ := metadata.FromOutgoingContext(ctx)
			log.Printf("Gateway - Passing metadata: %v", md)
			return md
		}), jsonOption,
	)

	if err = authpb.RegisterAuthServiceHandlerClient(context.Background(), mux, authClient.Client); err != nil {
		log.Fatalf("fail to register auth client: %v", err)
	}
	if err = productpb.RegisterProductServiceHandlerClient(context.Background(), mux, productClient.Client); err != nil {
		log.Fatalf("fail to register product client: %v", err)
	}
	if err = pb.RegisterOrderServiceHandlerClient(context.Background(), mux, orderClient.Client); err != nil {
		log.Fatalf("fail to register order client: %v", err)
	}
	authMiddleware := middleware.NewAuthMiddleWareConfig(authClient.Client)
	httpHandler := authMiddleware.AuthMiddleware(mux)
	lis, err := net.Listen("tcp", c.GatewayPort)
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	if err = http.Serve(lis, httpHandler); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}

}
