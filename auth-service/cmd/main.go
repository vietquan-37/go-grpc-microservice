package main

import (
	"fmt"
	"log"
	"net"

	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/db"
	"github.com/vietquan-37/auth-service/pkg/handler"
	"github.com/vietquan-37/auth-service/pkg/pb"
	"github.com/vietquan-37/auth-service/pkg/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func main() {
	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("error while loading config: %v", err)
	}
	d := db.DbConn(c.DbSource)
	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatalf("ERROR STARTING THE SERVER: %v", err)
	}
	repo := initAuthRepo(d)
	jwtMaker, err := config.NewJwtWrapper(c.JwtSecret)
	if err != nil {
		log.Fatalf("error while creating jwt: %v", err)
	}
	h := handler.NewAuthHandler(*jwtMaker, repo)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	fmt.Println("Server is listening on port 5051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
func initAuthRepo(db *gorm.DB) repository.IAuthRepo {
	return repository.NewAuthRepo(db)
}
