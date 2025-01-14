package main

import (
	commonclient "common/client"
	"common/discovery"
	"common/discovery/consul"
	"common/interceptor"
	"common/loggers"
	"common/mtdt"
	"common/routes"
	"common/validate"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/config"
	"github.com/vietquan-37/order-service/pkg/db"
	"github.com/vietquan-37/order-service/pkg/handler"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repo"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
)

func main() {

	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load from config:")
	}
	registry, err := consul.NewRegistry(c.ConsulAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to consul")
	}
	instanceId := discovery.GenerateInstanceID(c.ServiceName)
	if err := registry.Register(instanceId, c.ServiceName, c.GrpcServerAddress); err != nil {
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
	d := db.DbConn(c.DbSource)
	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start to server:")
	}
	orderRepo := InitOrderRepo(d)
	productClient := client.InitProductClient(registry)
	authClient := commonclient.InitAuthClient(registry)
	h := handler.NewOrderHandler(productClient, authClient, orderRepo)
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor:")
	}

	roles := routes.AccessibleRoles
	authInterceptor := interceptor.NewAuthInterceptor(authClient, roles())

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggers.GrpcLoggerInterceptor,
			mtdt.ForwardMetadataUnaryServerInterceptor(),
			authInterceptor.UnaryAuthInterceptor(),
			validateInterceptor.ValidateInterceptor(),
		),
	)
	pb.RegisterOrderServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	log.Info().Msgf("start  gRPC server server at %s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("fail to serve server:")
	}
}
func InitOrderRepo(db *gorm.DB) repo.IOrderRepo {
	return repo.NewOrderRepo(db)
}
