package app

import (
	"common/discovery"
	"common/discovery/consul"
	"common/kafka/producer"
	"common/loggers"
	"common/mtdt"
	"common/timeout"
	"common/validate"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/auth-service/pkg/config"
	"github.com/vietquan-37/auth-service/pkg/db"
	"github.com/vietquan-37/auth-service/pkg/handler"
	"github.com/vietquan-37/auth-service/pkg/pb"
	"github.com/vietquan-37/auth-service/pkg/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"sync"
	"time"
)

type Server struct {
	config     *config.Config
	db         *gorm.DB
	registry   *consul.Registry
	instanceId string
	producer   *producer.Producer
	grpcServer *grpc.Server
	healthSrv  *health.Server
	longTask   *sync.WaitGroup
}

func newServer() *Server {
	c, err := config.LoadConfig("../")
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
	}

	return &Server{
		config:   c,
		longTask: &sync.WaitGroup{},
	}
}

func (s *Server) initialize() error {
	if err := s.setupDatabase(); err != nil {
		return err
	}

	if err := s.setupServiceRegistry(); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupDatabase() error {
	s.db = db.DbConn(s.config.DbSource)
	return nil
}

func (s *Server) setupServiceRegistry() error {
	registry, err := consul.NewRegistry(s.config.ConsulAddress)
	if err != nil {
		return err
	}
	s.registry = registry
	s.instanceId = discovery.GenerateInstanceID(s.config.ServiceName)
	if err := s.registry.Register(s.instanceId, s.config.ServiceName, s.config.GrpcServerAddress); err != nil {
		return err
	}
	return consul.RegisterConsulResolver(s.registry.Client)
}

func (s *Server) setupGRPCServer() error {
	dependencies, err := s.setupDependencies()
	if err != nil {
		return err
	}

	interceptors := s.setupInterceptors()

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	pb.RegisterAuthServiceServer(s.grpcServer, dependencies.AuthHandler)
	reflection.Register(s.grpcServer)
	s.setupHealthCheck()

	return nil
}

func (s *Server) setupDependencies() (*ServiceDependencies, error) {
	repo := repository.NewAuthRepo(s.db, s.config.AdminUserName, s.config.AdminPassword)

	jwtMaker, err := config.NewJwtWrapper(s.config.JwtSecret)
	if err != nil {
		return nil, err
	}

	s.producer = producer.NewProducer(s.config.BrokerAddress)

	h := handler.NewAuthHandler(*jwtMaker, repo, *s.config, s.longTask, s.producer)

	return &ServiceDependencies{
		AuthHandler: h,
	}, nil
}

func (s *Server) setupInterceptors() []grpc.UnaryServerInterceptor {
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor")
	}

	durations := time.Second * 6
	return []grpc.UnaryServerInterceptor{
		loggers.GrpcLoggerInterceptor,
		mtdt.ForwardMetadataUnaryServerInterceptor(),
		validateInterceptor.ValidateInterceptor(),
		timeout.UnaryTimeoutInterceptor(durations),
	}
}

func (s *Server) setupHealthCheck() {
	s.healthSrv = health.NewServer()
	healthpb.RegisterHealthServer(s.grpcServer, s.healthSrv)
	s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
}

func (s *Server) start() error {
	lis, err := net.Listen("tcp", s.config.GrpcServerAddress)
	if err != nil {
		return err
	}

	go func() {
		log.Info().Msgf("Starting gRPC server at %s", lis.Addr().String())
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	log.Info().Msg("gRPC server started")
	return nil
}

func (s *Server) startHealthCheck() {
	go func() {
		for {
			if err := s.registry.HealthCheck(s.instanceId); err != nil {
				log.Fatal().Err(err).Msg("failed to health check service")
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (s *Server) gracefulShutdown() {
	log.Info().Msg("Shutting down server...")

	if s.healthSrv != nil {
		s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}

	s.longTask.Wait()
	log.Info().Msg("waiting for goroutines to finish")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	log.Info().Msg("Server stopped gracefully")
}

func (s *Server) cleanup() {
	if err := s.producer.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close producer connection")
	}
	if s.registry != nil && s.instanceId != "" {
		s.registry.Deregister(s.instanceId, s.config.ServiceName)
	}
}

type ServiceDependencies struct {
	AuthHandler pb.AuthServiceServer
}
