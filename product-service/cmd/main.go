package main

import (
	"github.com/vietquan-37/product-service/pkg/config"
	"github.com/vietquan-37/product-service/pkg/db"
	"github.com/vietquan-37/product-service/pkg/handler"
	"github.com/vietquan-37/product-service/pkg/pb"
	"github.com/vietquan-37/product-service/pkg/repo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"log"
	"net"
)

func main() {
	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("load config failed, err:%v", err)
	}
	d := db.DbConn(c.DbSource)
	lis, err := net.Listen("tcp", c.GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	r := NewRepoInit(d)
	h := handler.NewProductHandler(r)
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
func NewRepoInit(db *gorm.DB) repo.IProductRepo {
	return repo.NewProductRepo(db)
}
