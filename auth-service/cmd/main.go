package main

import (
	"common/discovery"
	"common/discovery/consul"
	"common/loggers"
	"common/mtdt"
	"common/timeout"
	"common/validate"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/db"
	"github.com/vietquan-37/auth-service/pkg/email"
	"github.com/vietquan-37/auth-service/pkg/handler"
	"github.com/vietquan-37/auth-service/pkg/pb"
	"github.com/vietquan-37/auth-service/pkg/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var interuptSignal = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {

	c, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignal...)
	defer stop()
	registry, err := consul.NewRegistry(c.ConsulAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to consul")
	}
	instanceId := discovery.GenerateInstanceID(c.ServiceName)
	if err := registry.Register(instanceId, c.ServiceName, c.GrpcServerAddress, c.Resolve); err != nil {
		log.Fatal().Err(err).Msg("failed to register service")
	}
	defer registry.Deregister(instanceId, c.ServiceName)
	go func() {
		for {
			if err := registry.HealthCheck(instanceId); err != nil {
				log.Fatal().Err(err).Msg("failed to health check service")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	d := db.DbConn(c.DbSource)
	repo := initAuthRepo(d, c.AdminUserName, c.AdminPassword)

	jwtMaker, err := config.NewJwtWrapper(c.JwtSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load jwt secret")
	}
	longTask := &sync.WaitGroup{}
	mail := email.NewMailService(c.SMTPHost, c.SMTPPort, c.EmailUsername, c.EmailPassword, c.EmailUsername)
	h := handler.NewAuthHandler(*jwtMaker, repo, *c, mail, longTask)

	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor")
	}
	durations := time.Second * 6
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggers.GrpcLoggerInterceptor,
			mtdt.ForwardMetadataUnaryServerInterceptor(),
			validateInterceptor.ValidateInterceptor(),
			timeout.UnaryTimeoutInterceptor(durations),
		),
	)
	pb.RegisterAuthServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", c.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	go func() {
		log.Info().Msgf("Starting gRPC server at %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Shutting down server...")
	longTask.Wait()
	log.Info().Msg("waiting for goroutines to finish")
	grpcServer.GracefulStop()
	log.Info().Msg("Server stopped gracefully")
}

func initAuthRepo(db *gorm.DB, adminUsername, adminPassword string) repository.IAuthRepo {
	return repository.NewAuthRepo(db, adminUsername, adminPassword)
}
