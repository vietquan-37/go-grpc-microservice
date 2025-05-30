package main

import (
	commonclient "common/client"
	"common/discovery"
	"common/discovery/consul"
	"common/interceptor"
	"common/loggers"
	"common/mtdt"
	"common/routes"
	"common/timeout"
	"common/validate"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/config"
	"github.com/vietquan-37/order-service/pkg/db"
	"github.com/vietquan-37/order-service/pkg/handler"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repo"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {

	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load from config:")
	}
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()
	registry, err := consul.NewRegistry(c.ConsulAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to consul")
	}
	instanceId := discovery.GenerateInstanceID(c.ServiceName)
	if err := registry.Register(instanceId, c.ServiceName, c.GrpcServerAddress, c.Mode); err != nil {
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
	err = consul.RegisterConsulResolver(registry.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register consul resolver")
	}
	d := db.DbConn(c.DbSource)
	orderRepo := InitOrderRepo(d)
	productClient, err := client.InitProductClient(c.ProductServiceName)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot init product client")
	}
	authClient, err := commonclient.InitAuthClient(c.AuthServiceName)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot init auth client")
	}
	h := handler.NewOrderHandler(productClient, authClient, orderRepo)
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor:")
	}

	roles := routes.AccessibleRoles
	authInterceptor := interceptor.NewAuthInterceptor(authClient, roles())
	durations := time.Second * 6
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggers.GrpcLoggerInterceptor,
			mtdt.ForwardMetadataUnaryServerInterceptor(),
			authInterceptor.UnaryAuthInterceptor(),
			validateInterceptor.ValidateInterceptor(),
			timeout.UnaryTimeoutInterceptor(durations),
		),
	)
	pb.RegisterOrderServiceServer(grpcServer, h)
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start to server:")
	}

	go func() {
		log.Info().Msgf("start  gRPC server server at %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Error().Err(err).Msg("fail to serve server:")

		}
	}()
	<-ctx.Done()
	log.Info().Msgf("stop  gRPC server server at %s", lis.Addr().String())
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
	log.Info().Msg("stopped gRPC server server")
}
func InitOrderRepo(db *gorm.DB) repo.IOrderRepo {
	return repo.NewOrderRepo(db)
}
