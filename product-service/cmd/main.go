package main

import (
	commonclient "common/client"
	"common/interceptor"
	"common/loggers"
	"common/mtdt"
	"common/routes"
	"common/validate"
	"github.com/rs/zerolog/log"

	"github.com/vietquan-37/product-service/pkg/config"
	"github.com/vietquan-37/product-service/pkg/db"
	"github.com/vietquan-37/product-service/pkg/handler"
	"github.com/vietquan-37/product-service/pkg/pb"
	"github.com/vietquan-37/product-service/pkg/repo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"net"
)

func main() {
	c, err := config.LoadConfig("./")
	if err != nil {

		log.Fatal().Err(err).Msg("fail to load config file:")
	}
	d := db.DbConn(c.DbSource)
	lis, err := net.Listen("tcp", c.GrpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to listen to server:")
	}
	r := NewRepoInit(d)
	h := handler.NewProductHandler(r)
	authClient := commonclient.InitAuthClient(c.AuthUrl)
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor:")
	}
	roles := routes.AccessibleRoles
	authInterceptor := interceptor.NewAuthInterceptor(*authClient, roles())
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authInterceptor.UnaryAuthInterceptor(),
			loggers.GrpcLoggerInterceptor,
			mtdt.ForwardMetadataUnaryServerInterceptor(),

			validateInterceptor.ValidateInterceptor(),
		),
	)
	pb.RegisterProductServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	log.Info().Msgf("start  gRPC server server at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("fail to serve  server:")
	}

}
func NewRepoInit(db *gorm.DB) repo.IProductRepo {
	return repo.NewProductRepo(db)
}
