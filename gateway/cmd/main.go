package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vietquan-37/gateway/pkg/auth"
	"github.com/vietquan-37/gateway/pkg/config"
	"github.com/vietquan-37/gateway/pkg/order"
	"github.com/vietquan-37/gateway/pkg/order/pb"
	"github.com/vietquan-37/gateway/pkg/product"
	productpb "github.com/vietquan-37/gateway/pkg/product/pb"
	"log"
	"net/http"

	authpb "github.com/vietquan-37/gateway/pkg/auth/pb"
)

func main() {
	c, err := config.LoadConfig("../")
	if err != nil {
		log.Fatalf("fail to load config: %v", err)
	}
	authClient := auth.InitAuthClient(c.AuthUrl)
	productClient := product.InitProductClient(c.ProductUrl)
	orderClient := order.InitOrderClient(c.OrderUrl)
	mux := runtime.NewServeMux()
	if err = authpb.RegisterAuthServiceHandlerClient(context.Background(), mux, authClient.Client); err != nil {
		log.Fatalf("fail to register auth client: %v", err)
	}
	if err = productpb.RegisterProductServiceHandlerClient(context.Background(), mux, productClient.Client); err != nil {
		log.Fatalf("fail to register auth client: %v", err)
	}
	if err = pb.RegisterOrderServiceHandlerClient(context.Background(), mux, orderClient.Client); err != nil {
		log.Fatalf("fail to register auth client: %v", err)
	}
	if err = http.ListenAndServe(c.GatewayPort, mux); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}

}
