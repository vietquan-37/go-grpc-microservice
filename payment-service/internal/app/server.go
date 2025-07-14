package app

import (
	"common/discovery"
	"common/discovery/consul"
	"common/loggers"
	"common/mtdt"
	"common/timeout"
	"common/validate"
	"github.com/rs/zerolog/log"
	"github.com/vietquan-37/payment-service/pkg/config"
	"github.com/vietquan-37/payment-service/pkg/handler"
	"github.com/vietquan-37/payment-service/pkg/pb"
	"github.com/vietquan-37/payment-service/pkg/provider/stripe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Server struct {
	cfg          *config.Config
	registry     *consul.Registry
	instanceId   string
	grpcServer   *grpc.Server
	healthServer *health.Server
}

func newServer() *Server {
	c, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	return &Server{
		cfg: c,
	}
}
func (s *Server) initialize() error {
	if err := s.setupServiceRegistry(); err != nil {
		return err
	}
	return nil
}
func (s *Server) setupServiceRegistry() error {
	registry, err := consul.NewRegistry(s.cfg.ConsulAddress)
	if err != nil {
		return err
	}
	s.registry = registry
	s.instanceId = discovery.GenerateInstanceID(s.cfg.ServiceName)
	if err := registry.Register(s.instanceId, s.cfg.ServiceName, s.cfg.GrpcAddress); err != nil {
		return err
	}
	return consul.RegisterConsulResolver(registry.Client)
}
func (s *Server) setupGRPCServer() error {
	dependencies := s.setupDependencies()
	interceptor := s.setupInterceptors()
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor...),
	)
	pb.RegisterPaymentServiceServer(s.grpcServer, dependencies.PaymentHandler)
	reflection.Register(s.grpcServer)
	s.setupHealthCheck()
	return nil
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
	s.healthServer = health.NewServer()
	healthpb.RegisterHealthServer(s.grpcServer, s.healthServer)
	s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

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
func (s *Server) start() error {
	lis, err := net.Listen("tcp", s.cfg.GrpcAddress)
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

	if s.healthServer != nil {
		s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	log.Info().Msg("Server stopped gracefully")
}
func (s *Server) cleanup() {
	if s.registry != nil && s.instanceId != "" {
		s.registry.Deregister(s.instanceId, s.cfg.ServiceName)
	}
}

type ServiceDependencies struct {
	PaymentHandler pb.PaymentServiceServer
}

func (s *Server) setupDependencies() *ServiceDependencies {
	provider := stripe.NewPaymentProvider(s.cfg)
	h := handler.NewPaymentHandler(provider)
	return &ServiceDependencies{PaymentHandler: h}
}
