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
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
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
	if err := registry.Register(instanceId, c.ServiceName, c.GrpcAddr); err != nil {
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
	authClient, err := commonclient.InitAuthClient()
	if err != nil {
		log.Fatal().Err(err).Msg("fail to init auth client")
	}
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
	pb.RegisterProductServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", c.GrpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to listen to server:")
	}
	waitGroup, ctx := errgroup.WithContext(ctx)
	waitGroup.Go(func() error {
		log.Info().Msgf("start  gRPC server server at %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return err
			}
			log.Error().Err(err).Msg("fail to serve  server:")
			return nil
		}
		return nil
	})
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msgf("shutting down gRPC server server at: %s", lis.Addr().String())

		grpcServer.GracefulStop()
		log.Info().Msgf("server stopped")
		return nil
	})
	if err := waitGroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("fail to wait for gRPC server:")
	}

}
func NewRepoInit(db *gorm.DB) repo.IProductRepo {
	return repo.NewProductRepo(db)
}
