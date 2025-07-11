package app

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
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/order-service/pkg/client"
	"github.com/vietquan-37/order-service/pkg/config"
	"github.com/vietquan-37/order-service/pkg/db"
	"github.com/vietquan-37/order-service/pkg/handler"
	"github.com/vietquan-37/order-service/pkg/pb"
	"github.com/vietquan-37/order-service/pkg/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"time"
)

type ServiceDependencies struct {
	OrderHandler *handler.OrderHandler
}

func (s *Server) setupDependencies() (*ServiceDependencies, error) {
	repo := repository.NewOrderRepo(s.db)

	productClient, authClient, err := s.setupExternalService()
	if err != nil {
		return nil, err
	}
	h := handler.NewOrderHandler(productClient, authClient, repo)
	return &ServiceDependencies{
		h,
	}, nil
}
func (s *Server) setupExternalService() (*client.ProductClient, *commonclient.AuthClient, error) {
	productClient, err := client.InitProductClient(s.config.ProductServiceName)
	if err != nil {
		return nil, nil, err
	}
	authClient, err := commonclient.InitAuthClient(s.config.AuthServiceName)
	s.authClient = authClient
	if err != nil {
		return nil, nil, err
	}
	return productClient, authClient, nil
}

type Server struct {
	config     *config.Config
	db         *gorm.DB
	registry   *consul.Registry
	authClient *commonclient.AuthClient
	instanceId string
	grpcServer *grpc.Server
	healthSrv  *health.Server
}

func newService() *Server {
	c, err := config.LoadConfig("../")
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
	}

	return &Server{
		config: c,
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
	registry, err := consul.NewRegistry(s.config.ConsulAddr)
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
	interceptors := s.setupInterceptor()
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	pb.RegisterOrderServiceServer(s.grpcServer, dependencies.OrderHandler)
	reflection.Register(s.grpcServer)
	s.setupHealthCheck()
	return nil
}
func (s *Server) setupInterceptor() []grpc.UnaryServerInterceptor {
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor")
	}
	roles := routes.AccessibleRoles
	authInterceptor := interceptor.NewAuthInterceptor(s.authClient, roles())
	durations := time.Second * 6
	return []grpc.UnaryServerInterceptor{
		loggers.GrpcLoggerInterceptor,
		mtdt.ForwardMetadataUnaryServerInterceptor(),
		authInterceptor.UnaryAuthInterceptor(),
		validateInterceptor.ValidateInterceptor(),
		timeout.UnaryTimeoutInterceptor(durations),
	}
}

func (s *Server) setupHealthCheck() {
	s.healthSrv = health.NewServer()
	healthpb.RegisterHealthServer(s.grpcServer, s.healthSrv)
	s.healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
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

	log.Info().Msg("waiting for goroutines to finish")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	log.Info().Msg("Server stopped gracefully")
}

func (s *Server) cleanup() {
	if s.registry != nil && s.instanceId != "" {
		s.registry.Deregister(s.instanceId, s.config.ServiceName)
	}
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
