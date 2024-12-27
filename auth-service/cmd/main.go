package main

import (
	"common/interceptor"
	"github.com/rs/zerolog/log"
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
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file: ")
	}
	d := db.DbConn(c.DbSource)
	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen: ")
	}
	repo := initAuthRepo(d)
	jwtMaker, err := config.NewJwtWrapper(c.JwtSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load jwt secret: ")
	}
	validateInterceptor, err := interceptor.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create interceptor: ")
	}
	h := handler.NewAuthHandler(*jwtMaker, repo)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.GrpcLoggerInterceptor,
			validateInterceptor.ValidateInterceptor(),
		))
	pb.RegisterAuthServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	log.Info().Msgf("start  gRPC server server at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve: ")
	}
}
func initAuthRepo(db *gorm.DB) repository.IAuthRepo {
	return repository.NewAuthRepo(db)
}
