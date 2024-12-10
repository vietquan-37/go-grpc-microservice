package main

import (
	"fmt"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repo"

	"github.com/vietquan-37/order-service/pkg/handler"

	"github.com/vietquan-37/order-service/pkg/config"
	"github.com/vietquan-37/order-service/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"log"
	"net"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}
	d := db.DbConn(c.DbSource)
	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatalf("ERROR STARTING THE SERVER: %v", err)
	}
	orderRepo := InitOrderRepo(d)
	detailRepo := InitOrderDetailRepo(d)
	productClient := client.InitProductClient(c.ProductURL)
	authClient := client.InitAuthClient(c.AuthURL)
	h := handler.NewOrderHandler(productClient, authClient, orderRepo, detailRepo)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	fmt.Println("Server is listening on port 5054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
func InitOrderRepo(db *gorm.DB) repo.IOrderRepo {
	return repo.NewOrderRepo(db)
}
func InitOrderDetailRepo(db *gorm.DB) repo.IOrderDetailRepo {
	return repo.NewOrderDetailRepo(db)
}
