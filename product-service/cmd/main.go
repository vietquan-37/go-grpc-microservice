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
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"os"
	"os/signal"
	"syscall"
	"time"

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

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	c, err := config.LoadConfig("./")
	if err != nil {

		log.Fatal().Err(err).Msg("fail to load config file:")
	}
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()
	registry, err := consul.NewRegistry(c.ConsulAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to consul")
	}
	instanceId := discovery.GenerateInstanceID(c.ServiceName)
	if err := registry.Register(instanceId, c.ServiceName, c.GrpcAddr, c.Mode); err != nil {
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
		log.Fatal().Err(err).Msg("fail to register consul resolver")
	}
	d := db.DbConn(c.DbSource)
	r := NewRepoInit(d)
	h := handler.NewProductHandler(r)
	authClient, err := commonclient.InitAuthClient(c.AuthServiceName)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to init auth client")
	}
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
	pb.RegisterProductServiceServer(grpcServer, h)
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", c.GrpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to listen to server:")
	}
	go func() {
		log.Info().Msgf("Starting gRPC server at %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()
	<-ctx.Done()
	log.Info().Msg("Shutting down server...")
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
	log.Info().Msg("Server stopped gracefully")
}
func NewRepoInit(db *gorm.DB) repo.IProductRepo {
	return repo.NewProductRepo(db)
}
