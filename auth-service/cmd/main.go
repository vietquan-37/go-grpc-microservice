package main

import (
	"common/discovery"
	"common/discovery/consul"
	"common/loggers"
	"common/mtdt"
	"common/validate"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/db"
	"github.com/vietquan-37/auth-service/pkg/handler"
	"github.com/vietquan-37/auth-service/pkg/pb"
	"github.com/vietquan-37/auth-service/pkg/repository"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var interuptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {

	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file: ")
	}
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignal...)
	defer stop()
	registry, err := consul.NewRegistry(c.ConsulAddress)
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
	defer registry.Deregister(instanceId, c.ServiceName)
	d := db.DbConn(c.DbSource)

	repo := initAuthRepo(d)
	jwtMaker, err := config.NewJwtWrapper(c.JwtSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load jwt secret: ")
	}
	h := handler.NewAuthHandler(*jwtMaker, repo, *c)
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor:")
	}

	grpcServer := grpc.NewServer(

		grpc.ChainUnaryInterceptor(
			loggers.GrpcLoggerInterceptor,
			mtdt.ForwardMetadataUnaryServerInterceptor(),
			validateInterceptor.ValidateInterceptor(),
		))
	pb.RegisterAuthServiceServer(grpcServer, h)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen: ")
	}
	waitGroup, ctx := errgroup.WithContext(ctx)
	waitGroup.Go(func() error {
		log.Info().Msgf("start  gRPC server server at %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			if !errors.Is(err, grpc.ErrServerStopped) {
				log.Error().Err(err).Msg("failed to serve: ")
				return err
			}

		}
		return nil
	})
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msgf("stop  gRPC server server at %s", lis.Addr().String())
		//important
		grpcServer.GracefulStop()
		log.Info().Msg("stopped gRPC server server")
		return nil
	})
	//
	if err := waitGroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("cannot wait for server:")
	}

}
func initAuthRepo(db *gorm.DB) repository.IAuthRepo {
	return repository.NewAuthRepo(db)
}
