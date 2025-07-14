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
	"github.com/vietquan-37/product-service/pkg/pb"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/vietquan-37/product-service/pkg/repository"
	"time"

	"github.com/vietquan-37/product-service/pkg/config"
	"github.com/vietquan-37/product-service/pkg/db"
	"github.com/vietquan-37/product-service/pkg/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"gorm.io/gorm"
)

type ServiceDependencies struct {
	ProductHandler pb.ProductServiceServer
}

func (s *Server) setupDependencies() (*ServiceDependencies, error) {
	repo := repository.NewProductRepo(s.db)

	h := handler.NewProductHandler(repo)
	return &ServiceDependencies{
		h,
	}, nil
}

type Server struct {
	config     *config.Config
	db         *gorm.DB
	registry   *consul.Registry
	instanceId string
	grpcServer *grpc.Server
	healthSrv  *health.Server
}

func newService() *Server {
	c, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
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
func (s *Server) setupServiceRegistry() error {
	registry, err := consul.NewRegistry(s.config.ConsulAddr)
	if err != nil {
		return err
	}
	s.registry = registry
	s.instanceId = discovery.GenerateInstanceID(s.config.ServiceName)
	if err := s.registry.Register(s.instanceId, s.config.ServiceName, s.config.GrpcAddr); err != nil {
		return err
	}
	return consul.RegisterConsulResolver(s.registry.Client)
}
func (s *Server) setupDatabase() error {
	s.db = db.DbConn(s.config.DbSource)
	return nil
}
func (s *Server) setupGrpcServer() error {
	dependencies, err := s.setupDependencies()
	if err != nil {
		return err
	}
	interceptors := s.setupInterceptor()
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	pb.RegisterProductServiceServer(s.grpcServer, dependencies.ProductHandler)
	reflection.Register(s.grpcServer)
	s.setupHealthCheck()
	return nil
}
func (s *Server) setupInterceptor() []grpc.UnaryServerInterceptor {
	validateInterceptor, err := validate.NewValidationInterceptor()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create validator interceptor")
	}
	authClient, err := commonclient.InitAuthClient(s.config.AuthServiceName)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot init auth client")
	}
	roles := routes.AccessibleRoles
	authInterceptor := interceptor.NewAuthInterceptor(authClient, roles())
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
func (s *Server) cleanup() {
	if s.registry != nil && s.instanceId != "" {
		s.registry.Deregister(s.instanceId, s.config.ServiceName)
	}
}
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.GrpcAddr)
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
